package connector

import (
	"context"
	"errors"

	"github.com/daemonfire/mister-fruits/internal/model"
)

type CouponStore interface {
	FindCoupon(ctx context.Context, ownerID, couponID string) (model.Coupon, error)
	CreateOrUpdateCoupon(ctx context.Context, ownerID, couponID string, coupon model.Coupon) error
}

var (
	ErrCouponNotFound = errors.New("error: coupon not found")
)
