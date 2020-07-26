package coupon

import (
	"context"
	"sync"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
)

var _ connector.CouponStore = (*InMemoryStore)(nil)

type InMemoryStore struct {
	sync.RWMutex

	coupons []model.Coupon
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		coupons: make([]model.Coupon, 0),
	}
}

func (i *InMemoryStore) CreateOrUpdateCoupon(_ context.Context, ownerID, cartID string, coupon model.Coupon) error {
	i.Lock()
	defer i.Unlock()
	for idx, p := range i.coupons {
		if p.ID == cartID && p.OwnerID == ownerID {
			i.coupons[idx] = coupon
			return nil
		}
	}
	i.coupons = append(i.coupons, coupon)
	return nil
}

func (i *InMemoryStore) FindCoupon(_ context.Context, ownerID, couponID string) (model.Coupon, error) {
	i.RLock()
	defer i.RUnlock()
	for _, p := range i.coupons {
		if p.ID == couponID && p.OwnerID == ownerID {
			return p, nil
		}
	}
	return model.Coupon{}, connector.ErrCouponNotFound
}
