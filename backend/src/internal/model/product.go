package model

import "github.com/shopspring/decimal"

type Product struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	GrossPrice decimal.Decimal `json:"grossPrice"`
	// Imagine a currency and a tax rate and so on...
}

var (
	Apple  = &Product{ID: "a699bd1b-bf8e-4569-be13-f0b22b589e31", Name: "Apple", GrossPrice: decimal.NewFromFloat(1.37)}
	Banana = &Product{ID: "d787ca6c-0056-48e1-bbc6-b856090cf051", Name: "Banana", GrossPrice: decimal.NewFromFloat(0.99)}
	Pear   = &Product{ID: "633f1b20-fc3d-413c-aab8-ec9f80ccc99b", Name: "Pear", GrossPrice: decimal.NewFromFloat(0.37)}
	Orange = &Product{ID: "101d4008-2761-4d6a-bdc6-6946efdafd5f", Name: "Orange", GrossPrice: decimal.NewFromFloat(2.01)}
)
