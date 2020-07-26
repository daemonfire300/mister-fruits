package connector

import "github.com/daemonfire/mister-fruits/internal/model"

type Authenticator interface {
	Check(username, password string) error
	CheckToken(token string) (string, error)
	GenerateToken(username string) (model.Token, error)
}
