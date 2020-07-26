package payment

import (
	"log"

	"github.com/shopspring/decimal"

	"github.com/daemonfire/mister-fruits/internal/connector"
)

type DummyGateway struct {
}

func (d *DummyGateway) ProcessPayment(amount decimal.Decimal) error {
	log.Println("Processed payment: ", amount.String(), "successfully")
	return nil
}

var _ connector.PaymentGateway = (*DummyGateway)(nil)
