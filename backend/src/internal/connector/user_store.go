package connector

import (
	"context"
	"errors"

	"github.com/daemonfire/mister-fruits/internal/model"
)

type UserStore interface {
	FindUser(ctx context.Context, username string) (model.User, error)
	StoreUser(ctx context.Context, user model.User) error
}

var (
	ErrUserAlreadyExists  = errors.New("error: user already exists")
	ErrUserNotFound       = errors.New("error: user not found")
	ErrInvalidCredentials = errors.New("error: invalid credentials")
)
