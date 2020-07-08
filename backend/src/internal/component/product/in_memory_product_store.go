package product

import (
	"context"
	"sync"

	"github.com/shopspring/decimal"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
)

var _ connector.ProductStore = (*InMemoryStore)(nil)

type InMemoryStore struct {
	sync.RWMutex

	products []model.Product
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		products: []model.Product{
			{ID: "a699bd1b-bf8e-4569-be13-f0b22b589e31", Name: "Apple", GrossPrice: decimal.NewFromFloat(1.37)},
			{ID: "d787ca6c-0056-48e1-bbc6-b856090cf051", Name: "Banana", GrossPrice: decimal.NewFromFloat(0.99)},
			{ID: "633f1b20-fc3d-413c-aab8-ec9f80ccc99b", Name: "Pear", GrossPrice: decimal.NewFromFloat(0.37)},
			{ID: "101d4008-2761-4d6a-bdc6-6946efdafd5f", Name: "Orange", GrossPrice: decimal.NewFromFloat(2.01)},
		},
	}
}

func (i *InMemoryStore) FindProduct(_ context.Context, id string) (model.Product, error) {
	i.RLock()
	defer i.RUnlock()
	for _, p := range i.products {
		if p.ID == id {
			return p, nil
		}
	}
	return model.Product{}, connector.ErrProductNotFound
}

func (i *InMemoryStore) AllProducts(ctx context.Context) ([]model.Product, error) {
	i.RLock()
	defer i.RUnlock()
	out := make([]model.Product, len(i.products))
	copy(out, i.products) // although this is a dummy application we do not want to expose the underlying slice to the
	// caller, since the caller would then be able to manipulate the slice
	// although slices are passed by value, the struct that is a slice in go internals speak contains a header, i.e.,
	// pointer to where the slice starts, i.e., comparison to C/C++ pointer to the first entry of an array
	// therefore slices are "kinda" passed by reference.
	return out, nil
}
