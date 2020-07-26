package connector

import "github.com/shopspring/decimal"

type PaymentGateway interface {
	// ProcessPayment processes a payment.
	// Could add some more arguments, like list of items, etc.
	ProcessPayment(amount decimal.Decimal) error
}
