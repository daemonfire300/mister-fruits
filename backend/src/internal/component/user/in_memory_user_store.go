package user

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
)

var _ connector.UserStore = (*InMemoryStore)(nil)
var _ connector.Authenticator = (*InMemoryStore)(nil) // implement everything in one struct, dummy application, why bother

type InMemoryStore struct {
	sync.RWMutex

	users  []model.User
	tokens map[string]model.Token // 1 user 1 token for simplicity
}

func (i *InMemoryStore) Check(username, password string) error {
	i.RLock()
	defer i.RUnlock()
	for _, u := range i.users {
		if u.Username != username {
			continue
		}
		if u.Password != password {
			return connector.ErrInvalidCredentials
		}
		return nil
	}
	return connector.ErrUserNotFound
}

func (i *InMemoryStore) CheckToken(token string) (string, error) {
	i.RLock()
	defer i.RUnlock()
	for username, tok := range i.tokens {
		if tok.Value != token {
			continue
		}
		return username, nil
	}
	return "", connector.ErrInvalidCredentials
}

func (i *InMemoryStore) GenerateToken(username string) (model.Token, error) {
	i.Lock()
	defer i.Unlock()
	i.tokens[username] = model.Token{
		Value: uuid.New().String(),
	}
	return i.tokens[username], nil
}

func NewStore() *InMemoryStore {
	return &InMemoryStore{
		users:  make([]model.User, 0),
		tokens: make(map[string]model.Token),
	}
}

func (i *InMemoryStore) FindUser(username string) (model.User, error) {
	i.RLock()
	defer i.RUnlock()
	for _, u := range i.users {
		if u.Username == username {
			return u, nil
		}
	}
	return model.User{}, connector.ErrUserNotFound
}

func (i *InMemoryStore) StoreUser(user model.User) error {
	i.Lock()
	defer i.Unlock()
	for _, u := range i.users {
		if u.Username == user.Username {
			return connector.ErrUserAlreadyExists
		}
	}
	i.users = append(i.users, user)
	return nil
}
