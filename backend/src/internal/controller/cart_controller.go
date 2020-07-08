package controller

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
	"github.com/daemonfire/mister-fruits/internal/pkg/user"
)

type CartController struct {
	CartStore    connector.CartStore
	ProductStore connector.ProductStore
}

type CartView struct {
	ID   string          `json:"id"`
	Sets []model.CartSet `json:"sets"`
}

func (c *CartController) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	username := user.FromContext(ctx)
	cart, err := c.CartStore.FindCart(ctx, username, vars["id"])
	if err != nil {
		w.WriteHeader(404)
		// we do not care what error occurs, we do not want to leak information whether
		// the cart id was valid or not, i.e., if you query a cart that actually exists
		// but are not the owner --> still return 404.
		log.Println("cart not found", err)
		return
	}
	finalView := CartView{
		ID:   cart.ID,
		Sets: make([]model.CartSet, 0),
	}
	standaloneProducts := map[string]string{"Apple": "Apple", "Orange": "Orange"} // should
	// use IDs in the future in case products have the same literal name but different IDs
	standaloneProductCartSets := make(map[string]*model.CartSet, 0)
	productGroups := [][]model.ProductGroup{
		{{"Banana", 2}, {"Pear", 4}},
		{{"Apple", 99999}},
		{{"Orange", 99999}},
	}
	productGroupCartSets := make(map[int][]*model.CartSet, 0)

	var isGroupedProduct bool
	for _, item := range cart.Items {
		// each CardSet is a bucket that needs to be filled
		// or if all CardSets are saturated we add new CardSets
		// Pears and Bananas go together
		if standaloneProductName, ok := standaloneProducts[item.Name]; ok {
			standaloneProductCartSet, ok := standaloneProductCartSets[standaloneProductName]
			if !ok {
				standaloneProductCartSet = &model.CartSet{
					Items:              make([]model.Product, 0),
					Coupons:            make([]model.Coupon, 0),
					ComputedGrossPrice: decimal.Zero,
				}
			}
			standaloneProductCartSet.Items = append(standaloneProductCartSet.Items, item)
			continue
		}

		var groupID int
	detectedGroupdProduct:
		for idx, grp := range productGroups {
			for _, product := range grp {
				if item.Name == product.Name {
					isGroupedProduct = true
					groupID = idx
					break detectedGroupdProduct
				}
			}
		}
		if !isGroupedProduct {
			continue
		}
		if _, ok := productGroupCartSets[groupID]; !ok {
			productGroupCartSets[groupID] = make([]*model.CartSet, 0)
		}
		groupCartSets := productGroupCartSets[groupID]
		for _, groupCartSet := range groupCartSets {
			fillables := make(map[string]int, 0)
			for _, groupItem := range productGroups[groupID] {
				fillables[groupItem.Name] = groupItem.Quantity
			}
			for _, groupCartItem := range groupCartSet.Items {
				fillables[groupCartItem.Name]--
			}
			if fillables[item.Name] > 0 { // still room for item in this CartSet, let's put it there
				groupCartSet.Items = append(groupCartSet.Items, item)
				break
			}
		}
		// If we reach this code segment the item does neither belong to a standalone group nor was there space in any
		// of the existing groups. That means we need to create a new CartSet
		groupCartSets = append(groupCartSets, &model.CartSet{
			Items:              []model.Product{item},
			Coupons:            make([]model.Coupon, 0),
			ComputedGrossPrice: decimal.Zero,
		})
	}
}

func (c *CartController) Create(w http.ResponseWriter, r *http.Request) {

}

func (c *CartController) Update(w http.ResponseWriter, r *http.Request) {

}
