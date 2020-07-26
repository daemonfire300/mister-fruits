package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type ProductGroup struct {
	Product  *Product `json:"product"`
	Quantity int      `json:"quantity"`
}

type ProductGroupRule struct {
	Product  *Product
	Quantity int
}

type ProductSet struct {
	Products []Product `json:"-"`
}

func (s ProductSet) In(items []*ProductGroup) bool {
	in := make(map[string]int)
	for _, p := range s.Products {
		in[p.ID] = 0
	}
	for id := range in {
		for _, grp := range items {
			if grp.Product.ID != id {
				continue
			}
			in[id]++
		}
	}
	for _, count := range in {
		if count < 1 {
			return false
		}
	}
	return true
}

// ExactlyIn is like ProductSet.In but returns false if the s.Products ∪ items ≠ s.Products
func (s ProductSet) ExactlyIn(items []*ProductGroup) bool {
	in := make(map[string]int)
	for _, p := range s.Products {
		in[p.ID] = 0
	}
	for id := range in {
		for _, grp := range items {
			if _, ok := in[grp.Product.ID]; !ok {
				return false
			}
			if grp.Product.ID != id {
				continue
			}
			in[id]++
		}
	}
	for _, count := range in {
		if count < 1 {
			return false
		}
	}
	return true
}

type Coupon struct {
	ID                          string                      `json:"id"`
	Name                        string                      `json:"name"`
	ApplicableProductSets       []ProductSet                `json:"-"`
	RequiredProductQuantityList RequiredProductQuantityList `json:"-"`
	IsGlobal                    bool                        `json:"-"`
	Expiration                  time.Time                   `json:"-"`
	Discount                    decimal.Decimal             `json:"discount"`
	OwnerID                     string                      `json:"-"`
	AppliedToCart               string                      `json:"-"`
	AppliedAt                   time.Time                   `json:"-"`
}

type RequiredProductQuantityList []RequiredProductQuantity

func (l RequiredProductQuantityList) GetQuantityByID(id string) (ok bool, minQuanity, maxQuantity int) {
	for _, pq := range l {
		if pq.Product.ID == id {
			return true, pq.MinQuantity, pq.MaxQuantity
		}
	}
	return false, 0, 0
}

type RequiredProductQuantity struct {
	Product     Product `json:"-"`
	MinQuantity int     `json:"-"`
	MaxQuantity int     `json:"-"`
}

type CartSet struct {
	Items              []*ProductGroup `json:"items"`
	Coupons            []*Coupon       `json:"coupons"`
	ComputedGrossPrice decimal.Decimal `json:"grossPrice"`
}

func (s CartSet) GetGroupByName(name string) (ok bool, quantity int) {
	// could be optimized, i.e., using maps for items or using a search index
	for _, grp := range s.Items {
		if grp.Product.Name != name {
			continue
		}
		return true, grp.Quantity
	}
	return false, 0
}

func (s CartSet) GetGroupByID(id string) (ok bool, quantity int) {
	// could be optimized, i.e., using maps for items or using a search index
	for _, grp := range s.Items {
		if grp.Product.ID != id {
			continue
		}
		return true, grp.Quantity
	}
	return false, 0
}

type Cart struct {
	ID      string          `json:"id"`
	OwnerID string          `json:"-"`
	Items   []*ProductGroup `json:"items"`
	Coupons []Coupon        `json:"coupons"`
}
