package cart

import (
	"context"
	"log"
	"sync"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
)

var _ connector.CartStore = (*InMemoryStore)(nil)

type InMemoryStore struct {
	sync.RWMutex

	carts []model.Cart
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		carts: make([]model.Cart, 0),
	}
}

func (i *InMemoryStore) DeleteCart(ctx context.Context, ownerID, cartID string) error {
	i.Lock()
	defer i.Unlock()
	index := -1
	for idx, p := range i.carts {
		if p.ID == cartID && p.OwnerID == ownerID {
			index = idx
			break
		}
	}
	if index < 0 {
		return connector.ErrCartNotFound
	}
	i.carts = append(i.carts[:index], i.carts[index+1:]...)
	return nil
}

func (i *InMemoryStore) CreateOrUpdateCart(_ context.Context, ownerID, cartID string, cart model.Cart) error {
	i.Lock()
	defer i.Unlock()
	for idx, p := range i.carts {
		if p.ID == cartID && p.OwnerID == ownerID {
			i.carts[idx] = cart
			return nil
		}
	}
	i.carts = append(i.carts, cart)
	return nil
}

func (i *InMemoryStore) FindCart(_ context.Context, ownerID, cartID string) (model.Cart, error) {
	i.RLock()
	defer i.RUnlock()
	for _, p := range i.carts {
		if p.ID == cartID && p.OwnerID == ownerID {
			return p, nil
		}
	}
	log.Println("Debug, listing all carts", i.carts)
	return model.Cart{}, connector.ErrCartNotFound
}
