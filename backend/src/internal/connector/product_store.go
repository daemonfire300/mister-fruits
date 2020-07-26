package connector

import (
	"context"
	"errors"

	"github.com/daemonfire/mister-fruits/internal/model"
)

type ProductStore interface {
	FindProduct(ctx context.Context, id string) (model.Product, error)
	AllProducts(ctx context.Context) ([]model.Product, error)
}

var (
	ErrProductNotFound = errors.New("error: product not found")
)
