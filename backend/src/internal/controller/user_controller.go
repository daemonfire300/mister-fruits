package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/daemonfire/mister-fruits/internal/connector"
	"github.com/daemonfire/mister-fruits/internal/model"
)

type UserController struct {
	UserStore     connector.UserStore
	Authenticator connector.Authenticator
}

func (c *UserController) Signup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	var user model.User
	err := dec.Decode(&user)
	if err != nil {
		log.Println(fmt.Sprintf("error during decoding of user input: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = c.UserStore.FindUser(user.Username)
	if err != nil && !errors.Is(err, connector.ErrUserNotFound) {
		log.Println(fmt.Sprintf("error during user retrieval: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.UserStore.StoreUser(user)
	if err != nil {
		log.Println(fmt.Sprintf("error while persisting user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	var user model.User
	err := dec.Decode(&user)
	if err != nil {
		log.Println(fmt.Sprintf("error during decoding of user input: %v", err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = c.UserStore.FindUser(user.Username)
	if err != nil && !errors.Is(err, connector.ErrUserNotFound) {
		log.Println(fmt.Sprintf("error during user retrieval: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = c.Authenticator.Check(user.Username, user.Password)
	if err != nil {
		log.Println(fmt.Sprintf("could not validate user credentials: %v", err))
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tok, err := c.Authenticator.GenerateToken(user.Username)
	if err != nil {
		log.Println(fmt.Sprintf("could not generate token for user: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	if err := enc.Encode(tok); err != nil {
		log.Println(fmt.Sprintf("error while writing json to response writer: %v", err))
	}
}
