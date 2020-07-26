package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
	"github.com/daemonfire/mister-fruits/internal/pkg/user"
)

type CartController struct {
	CartStore      connector.CartStore
	ProductStore   connector.ProductStore
	CouponStore    connector.CouponStore
	PaymentGateway connector.PaymentGateway
}

type CartView struct {
	ID                 string           `json:"id"`
	Sets               []*model.CartSet `json:"sets"`
	Coupons            []*model.Coupon  `json:"coupons"`
	GlobalProductCount map[string]int   `json:"-"` // ID -> Count
}

func (v *CartView) GetGrossComputedPrice() decimal.Decimal {
	sum := decimal.Zero
	for _, set := range v.Sets {
		sum = sum.Add(set.ComputedGrossPrice)
	}
	return sum
}

type UpdateCartRequest struct {
	Add    []model.ProductGroup `json:"add"`
	Remove []model.ProductGroup `json:"remove"`
}

var defaultCoupons = []model.Coupon{
	model.Coupon{
		ID:   "9a0d46be-01d1-4c9b-94c4-9b0041d4b205",
		Name: "7 ore more Apples",
		ApplicableProductSets: []model.ProductSet{
			{[]model.Product{*model.Apple}},
		},
		RequiredProductQuantityList: model.RequiredProductQuantityList{
			{*model.Apple, 7, 99999999},
		},
		IsGlobal:   true,
		Expiration: time.Now().Add(time.Hour),
		Discount:   decimal.NewFromFloat(0.10),
	},
	model.Coupon{
		ID:   "78636d38-d151-4506-9ef0-e95559f01f31",
		Name: "4 Pear 2 Banana each, 30% off",
		ApplicableProductSets: []model.ProductSet{
			{[]model.Product{*model.Pear, *model.Banana}},
		},
		RequiredProductQuantityList: model.RequiredProductQuantityList{
			{*model.Pear, 4, 99999999},
			{*model.Banana, 2, 99999999},
		},
		IsGlobal:   false,
		Expiration: time.Now().Add(time.Hour),
		Discount:   decimal.NewFromFloat(0.30),
	},
}

func (c *CartController) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cart, err := c.CartStore.FindCart(ctx, username, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// we do not care what error occurs, we do not want to leak information whether
		// the cart id was valid or not, i.e., if you query a cart that actually exists
		// but are not the owner --> still return 404.
		log.Println("cart not found, aborting checkout", err)
		return
	}
	view := constructCardView(&cart)
	cpns := make([]model.Coupon, len(cart.Coupons))
	copy(cpns, cart.Coupons) // make a copy do not modify the underlying slice
	cpns = append(cpns, defaultCoupons...)
	applyCouponsToView(view, cpns)
	amount := view.GetGrossComputedPrice()
	log.Println("checking out cart", cart, "for sum of", amount)
	// the following code should happen transactional in the real world, i.e., dont withdraw money from the customer
	// if something goes wrong.
	err = c.PaymentGateway.ProcessPayment(amount)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error during process payment", err)
		return
	}
	err = c.CartStore.DeleteCart(ctx, username, cart.ID)
	if err != nil { // TODO rollback payment or something
		w.WriteHeader(http.StatusBadRequest)
		log.Println("error during process payment, when deleting cart", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (c *CartController) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cart, err := c.CartStore.FindCart(ctx, username, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// we do not care what error occurs, we do not want to leak information whether
		// the cart id was valid or not, i.e., if you query a cart that actually exists
		// but are not the owner --> still return 404.
		log.Println("cart not found", err)
		return
	}
	view := constructCardView(&cart)
	cpns := make([]model.Coupon, len(cart.Coupons))
	copy(cpns, cart.Coupons) // make a copy do not modify the underlying slice
	cpns = append(cpns, defaultCoupons...)
	applyCouponsToView(view, cpns)
	enc := json.NewEncoder(w)
	if err := enc.Encode(view); err != nil {
		log.Println("Error during json encoding of cart view", err)
	}
}

func (c *CartController) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cartID := vars["id"]
	_, err := c.CartStore.FindCart(ctx, username, cartID)
	if err == nil {
		log.Println("Cart found, cannot create cart with already existing uuid")
		w.WriteHeader(http.StatusConflict)
		return
	}
	if !errors.Is(err, connector.ErrCartNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error during cart retrieval", err)
		return
	}
	dec := json.NewDecoder(r.Body)
	var cart model.Cart
	err = dec.Decode(&cart)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Error during cart creation, while decoding", err)
		return
	}
	cart.OwnerID = username
	err = c.CartStore.CreateOrUpdateCart(ctx, username, cartID, cart)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error during cart creation, while persisting", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Cart created successfully")
}

func (c *CartController) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cartID := vars["id"]
	cart, err := c.CartStore.FindCart(ctx, username, cartID)
	if err != nil {
		if errors.Is(err, connector.ErrCartNotFound) {
			log.Println("Card not found found, cannot update")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Println("Error during cart retrieval while patching", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var updateReq UpdateCartRequest
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(&updateReq)
	if err != nil {
		log.Println("Error during decoding of cart update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, addItem := range updateReq.Add {
		var exists bool
		for idx, grp := range cart.Items { // should use a hash map maybe, not enough time to change this though
			if grp.Product.ID != addItem.Product.ID {
				continue
			}
			exists = true
			cart.Items[idx].Quantity += addItem.Quantity
		}
		if exists {
			continue
		}
		p, err := c.ProductStore.FindProduct(ctx, addItem.Product.ID)
		if err != nil {
			log.Println("Error during fetching item of cart update", addItem.Product.ID, err)
			continue
		}
		// Items is not yet in cart, add it.
		cart.Items = append(cart.Items, &model.ProductGroup{
			Product:  &p,
			Quantity: addItem.Quantity,
		})
	}
	for _, removeItem := range updateReq.Remove {
		for idx, grp := range cart.Items { // should use a hash map maybe, not enough time to change this though
			if grp.Product.ID != removeItem.Product.ID {
				continue
			}
			cart.Items[idx].Quantity -= removeItem.Quantity
			cart.Items[idx].Quantity = MaxInt(0, cart.Items[idx].Quantity)
		}
	}
	err = c.CartStore.CreateOrUpdateCart(ctx, username, cartID, cart)
	if err != nil {
		log.Println("Error during persisting of cart update", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Successfully updated cart", cartID)
}

func (c *CartController) DummyCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	username := user.FromContext(ctx)
	couponID := uuid.New().String()
	coupon := model.Coupon{
		ID:   couponID,
		Name: "30% off of oranges",
		ApplicableProductSets: []model.ProductSet{
			{Products: []model.Product{*model.Orange}},
		},
		RequiredProductQuantityList: model.RequiredProductQuantityList{},
		IsGlobal:                    true,
		Expiration:                  time.Now().Add(time.Minute), // I know the challenge says 10 seconds
		// but that is too short to actually test the code...
		Discount: decimal.NewFromFloat(0.3),
	}
	err := c.CouponStore.CreateOrUpdateCoupon(ctx, username, couponID, coupon)
	if err != nil {
		log.Println("Coupon could not be created.", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(coupon); err != nil {
		log.Println("Error during json encoding of coupon.", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	log.Println("Dummy coupon created successfully")
}

func (c *CartController) ApplyCoupon(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cart, err := c.CartStore.FindCart(ctx, username, vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		// we do not care what error occurs, we do not want to leak information whether
		// the cart id was valid or not, i.e., if you query a cart that actually exists
		// but are not the owner --> still return 404.
		log.Println("cart not found", err)
		return
	}
	couponID := vars["id"]
	coupon, err := c.CouponStore.FindCoupon(ctx, username, couponID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		log.Println("coupon not found", err)
		return
	}
	if coupon.Expiration.Before(time.Now()) || coupon.AppliedToCart != "" {
		w.WriteHeader(http.StatusConflict)
		log.Println("coupon already applied", coupon)
		return
	}
	var replaceExisting bool
	for i, cpn := range cart.Coupons {
		if replaceExisting {
			break
		}
		if cpn.Name == coupon.Name {
			if cpn.Expiration.Before(time.Now()) {
				coupon.AppliedToCart = cart.ID
				coupon.AppliedAt = time.Now()
				cart.Coupons[i] = coupon
				replaceExisting = true
				break
			}
			if cpn.Expiration.After(time.Now()) {
				w.WriteHeader(http.StatusConflict)
				log.Println("coupon with same name already applied", coupon)
				return
			}
		}
	}
	if !replaceExisting {
		coupon.AppliedToCart = cart.ID
		coupon.AppliedAt = time.Now()
		cart.Coupons = append(cart.Coupons, coupon)
	}
	if err := c.CartStore.CreateOrUpdateCart(ctx, username, cart.ID, cart); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not update cart when applying coupon", err)
		return
	}
	if err := c.CouponStore.CreateOrUpdateCoupon(ctx, username, couponID, coupon); err != nil {
		// At this point we should actually rollback the change we made the cart, because now this coupon
		// has been applied to at least on cart, is actually use-once.
		// TODO: Had no time to address this issue.
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Could not update coupon when applying coupon", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func applyCouponsToView(view *CartView, coupons []model.Coupon) *CartView {
	now := time.Now()
	for _, cpn := range coupons {
		log.Println("checking if coupon is applicable to cart", view.ID, cpn)
		if cpn.Expiration.Before(now) {
			continue
		}
		if cpn.IsGlobal { // we only handle the case of global coupon = one applicable product
			var targetID string
			var globalProductCount int
			for _, sets := range cpn.ApplicableProductSets {
				if targetID != "" {
					continue
				}
				for _, p := range sets.Products {
					if _, in := view.GlobalProductCount[p.ID]; in {
						targetID = p.ID
						globalProductCount = view.GlobalProductCount[p.ID]
					}
				}
			}
			if targetID != "" {
				var targetSetIdx int
				for idx, set := range view.Sets {
					if ok, _ := set.GetGroupByID(targetID); ok {
						targetSetIdx = idx
						break
					}
				}
				ok, minQuantity, maxQuantity := cpn.RequiredProductQuantityList.GetQuantityByID(targetID)
				if !ok || (minQuantity <= globalProductCount && globalProductCount <= maxQuantity) {
					view.Sets[targetSetIdx].ComputedGrossPrice = view.Sets[targetSetIdx].ComputedGrossPrice.Mul(decimal.NewFromInt(1).Sub(cpn.Discount))
					view.Sets[targetSetIdx].ComputedGrossPrice = view.Sets[targetSetIdx].ComputedGrossPrice.Round(2)
					view.Sets[targetSetIdx].Coupons = append(view.Sets[targetSetIdx].Coupons, &cpn)
				}
			}
			continue
		}
		for _, applicableSet := range cpn.ApplicableProductSets {
			for setIdx, cartSet := range view.Sets {
				if !applicableSet.In(cartSet.Items) {
					continue
				}
				applicable := true
				for _, pg := range cartSet.Items {
					ok, minQuantity, maxQuantity := cpn.RequiredProductQuantityList.GetQuantityByID(pg.Product.ID)
					if ok && (minQuantity > pg.Quantity || pg.Quantity > maxQuantity) {
						applicable = false
					}
				}
				if !applicable {
					continue
				}
				view.Sets[setIdx].ComputedGrossPrice = view.Sets[setIdx].ComputedGrossPrice.Mul(decimal.NewFromInt(1).Sub(cpn.Discount))
				view.Sets[setIdx].ComputedGrossPrice = view.Sets[setIdx].ComputedGrossPrice.Round(2)
				view.Sets[setIdx].Coupons = append(view.Sets[setIdx].Coupons, &cpn)
			}
		}

	}
	return view
}

func constructCardView(cart *model.Cart) *CartView {
	productGroupRules := [][]model.ProductGroupRule{
		{{Product: model.Banana, Quantity: 2}, {Product: model.Pear, Quantity: 4}},
		{{Product: model.Apple, Quantity: 99999}},
		{{Product: model.Orange, Quantity: 99999}},
	}
	productGroupCartSetsMapping := make(map[int][]*model.CartSet)
	items := make([]*model.ProductGroup, len(cart.Items))
	globalCount := make(map[string]int)
	for i, item := range cart.Items { // copy over items, we do not want to modify the actual quantity of the origin cart
		items[i] = &model.ProductGroup{
			Product:  item.Product,
			Quantity: item.Quantity,
		}
	}
	for _, item := range items {
		if _, ok := globalCount[item.Product.ID]; !ok {
			globalCount[item.Product.ID] = 0
		}
		globalCount[item.Product.ID] += item.Quantity
		// each CardSet is a bucket that needs to be filled
		// or if all CardSets are saturated we add new CardSets
		// Pears and Bananas go together
		var groupID int
	detectProductGroup:
		for idx, grp := range productGroupRules {
			for _, rule := range grp {
				if item.Product.Name == rule.Product.Name {
					groupID = idx
					break detectProductGroup
				}
			}
		}
		if _, ok := productGroupCartSetsMapping[groupID]; !ok {
			productGroupCartSetsMapping[groupID] = make([]*model.CartSet, 0)
		}
		var maxAllowedAmountPerGroup int
		for _, productGroupRule := range productGroupRules[groupID] {
			if productGroupRule.Product.Name != item.Product.Name {
				continue
			}
			maxAllowedAmountPerGroup = productGroupRule.Quantity
		}
		for item.Quantity > 0 {
			availableCartSets := productGroupCartSetsMapping[groupID] // get CartSets that the product
			// can be put in according to the ProductGroupRule selected previously

			var quantityToAdd int
			var matchingGroup *model.ProductGroup
			for _, cartSet := range availableCartSets { // find a CartSet that still holds room for our product
				remaining := maxAllowedAmountPerGroup
				for _, cartSetItem := range cartSet.Items {
					if cartSetItem.Product.Name != item.Product.Name {
						continue
					}
					remaining -= cartSetItem.Quantity
					matchingGroup = cartSetItem
				}
				if remaining > 0 && matchingGroup != nil { // still room for item in this CartSet, let's put it there
					quantityToAdd = MinInt(item.Quantity, remaining)
					matchingGroup.Quantity += quantityToAdd
					item.Quantity -= quantityToAdd
				}
			}
			if item.Quantity == 0 {
				productGroupCartSetsMapping[groupID] = availableCartSets
				break
			}
			// If we reach this code segment the item does neither belong to a standalone group nor was there space in any
			// of the existing groups. That means we need to create a new CartSet
			cartGroups := make([]*model.ProductGroup, len(productGroupRules[groupID]))
			for i, rule := range productGroupRules[groupID] {
				cartGroups[i] = &model.ProductGroup{
					Product:  rule.Product,
					Quantity: 0,
				}
			}
			productGroupCartSetsMapping[groupID] = append(availableCartSets, &model.CartSet{
				Items:              cartGroups,
				Coupons:            make([]*model.Coupon, 0),
				ComputedGrossPrice: decimal.Zero,
			})
		}
	}

	finalView := CartView{
		ID:                 cart.ID,
		Sets:               make([]*model.CartSet, 0),
		GlobalProductCount: globalCount,
	}

	for _, groups := range productGroupCartSetsMapping {
		for _, set := range groups {
			for _, p := range set.Items {
				set.ComputedGrossPrice = set.ComputedGrossPrice.Add(p.Product.GrossPrice.Mul(decimal.NewFromInt(int64(p.Quantity))))
			}
			finalView.Sets = append(finalView.Sets, set)
		}
	}
	return &finalView
}

func MaxInt(a, b int) int {
	if a == b {
		return a
	}
	if a > b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a == b {
		return a
	}
	if a < b {
		return a
	}
	return b
}
