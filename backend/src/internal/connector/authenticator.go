package connector

import "github.com/daemonfire/mister-fruits/internal/model"

type Authenticator interface {
	Check(username, password string) error
	CheckToken(username, token string) error
	GenerateToken(username string) (model.Token, error)
}
