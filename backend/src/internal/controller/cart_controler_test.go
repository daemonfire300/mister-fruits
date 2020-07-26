//nolint
package controller

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/daemonfire/mister-fruits/internal/model"
)

//			{ID: "a699bd1b-bf8e-4569-be13-f0b22b589e31", Name: "Apple", GrossPrice: decimal.NewFromFloat(1.37)},
//			{ID: "d787ca6c-0056-48e1-bbc6-b856090cf051", Name: "Banana", GrossPrice: decimal.NewFromFloat(0.99)},
//			{ID: "633f1b20-fc3d-413c-aab8-ec9f80ccc99b", Name: "Pear", GrossPrice: decimal.NewFromFloat(0.37)},
//			{ID: "101d4008-2761-4d6a-bdc6-6946efdafd5f", Name: "Orange", GrossPrice: decimal.NewFromFloat(2.01)},
func Test_constructCardView(t *testing.T) {
	cart := model.Cart{
		ID:      "",
		OwnerID: "",
		Items: []*model.ProductGroup{
			{model.Banana, 5},
			{model.Apple, 1},
			{model.Pear, 9},
			{model.Orange, 2},
		},
		Coupons: []model.Coupon{},
	}
	view := constructCardView(&cart)
	require.NotNil(t, view)
}

func Test_applyNoCouponsToView(t *testing.T) {
	cart := model.Cart{
		ID:      "",
		OwnerID: "",
		Items: []*model.ProductGroup{
			{model.Apple, 1},
		},
		Coupons: []model.Coupon{},
	}
	view := constructCardView(&cart)
	require.Equal(t, decimal.NewFromFloat(1.37), view.GetGrossComputedPrice())
	coupons := []model.Coupon{
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
	}
	applyCouponsToView(view, coupons)
	require.Equal(t, decimal.NewFromFloat(1.37), view.GetGrossComputedPrice())
}

func Test_applyGlobalAppleCouponsToView(t *testing.T) {
	cart := model.Cart{
		ID:      "",
		OwnerID: "",
		Items: []*model.ProductGroup{
			{model.Apple, 7},
		},
		Coupons: []model.Coupon{},
	}
	view := constructCardView(&cart)
	require.Equal(t, decimal.NewFromFloat(9.59), view.GetGrossComputedPrice())
	coupons := []model.Coupon{
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
	}
	applyCouponsToView(view, coupons)
	require.Equal(t, decimal.NewFromFloat(8.631), view.GetGrossComputedPrice())
}

func Test_applyPearBananaCouponsToView(t *testing.T) {
	cart := model.Cart{
		ID:      "",
		OwnerID: "",
		Items: []*model.ProductGroup{
			{model.Pear, 4},
			{model.Banana, 2},
		},
		Coupons: []model.Coupon{},
	}
	view := constructCardView(&cart)
	require.Equal(t, decimal.NewFromFloat(3.46), view.GetGrossComputedPrice())
	coupons := []model.Coupon{
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
	applyCouponsToView(view, coupons)
	require.Equal(t, decimal.NewFromFloat(2.422), view.GetGrossComputedPrice())
}
