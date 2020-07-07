package main

import (
	"log"
	"net/http"

	"github.com/daemonfire/mister-fruits/internal/api"
	"github.com/daemonfire/mister-fruits/internal/component/user"
	"github.com/daemonfire/mister-fruits/internal/controller"
)

func main() {
	log.Println("Creating controllers")
	userTokenStoreAndAuthenticator := user.NewStore()
	controllers := &controller.Controllers{
		ProductController: nil,
		CartController:    nil,
		UserController: &controller.UserController{
			UserStore:     userTokenStoreAndAuthenticator,
			Authenticator: userTokenStoreAndAuthenticator,
		},
	}
	log.Println("Registering controllers")
	router := api.Router(controllers)
	log.Println("Starting server")
	if err := http.ListenAndServe(":8181", router); err != nil {
		log.Fatalln(err)
	}
	log.Println("Done")
}
