package connector

import (
	"errors"

	"github.com/daemonfire/mister-fruits/internal/model"
)

type UserStore interface {
	FindUser(username string) (model.User, error)
	StoreUser(user model.User) error
}

var (
	ErrUserAlreadyExists  = errors.New("error: user already exists")
	ErrUserNotFound       = errors.New("error: user not found")
	ErrInvalidCredentials = errors.New("error: invalid credentials")
)
