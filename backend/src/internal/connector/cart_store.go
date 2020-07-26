package connector

import (
	"context"
	"errors"

	"github.com/daemonfire/mister-fruits/internal/model"
)

type CartStore interface {
	FindCart(ctx context.Context, ownerID, cartID string) (model.Cart, error)
	CreateOrUpdateCart(ctx context.Context, ownerID, cartID string, cart model.Cart) error
	DeleteCart(ctx context.Context, ownerID, cartID string) error
}

var (
	ErrCartNotFound = errors.New("error: cart not found")
)
