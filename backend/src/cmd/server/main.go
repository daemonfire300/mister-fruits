package main

import (
	"log"
	"net/http"

	"github.com/daemonfire/mister-fruits/internal/api"
	"github.com/daemonfire/mister-fruits/internal/component/cart"
	"github.com/daemonfire/mister-fruits/internal/component/coupon"
	"github.com/daemonfire/mister-fruits/internal/component/payment"
	"github.com/daemonfire/mister-fruits/internal/component/product"
	"github.com/daemonfire/mister-fruits/internal/component/user"
	"github.com/daemonfire/mister-fruits/internal/controller"
	"github.com/daemonfire/mister-fruits/internal/middleware"
)

func main() {
	log.Println("Creating controllers")
	userTokenStoreAndAuthenticator := user.NewStore()
	pStore := product.NewStore()
	controllers := &controller.Controllers{
		ProductController: &controller.ProductController{ProductStore: pStore},
		CartController: &controller.CartController{
			CartStore:      cart.NewStore(),
			ProductStore:   pStore,
			CouponStore:    coupon.NewStore(),
			PaymentGateway: &payment.DummyGateway{},
		},
		UserController: &controller.UserController{
			UserStore:     userTokenStoreAndAuthenticator,
			Authenticator: userTokenStoreAndAuthenticator,
		},
	}
	log.Println("Registering controllers")
	router := api.Router(controllers, &middleware.TokenGuardMiddleware{Authenticator: userTokenStoreAndAuthenticator})
	log.Println("Starting server")
	if err := http.ListenAndServe(":8181", router); err != nil {
		log.Fatalln(err)
	}
	log.Println("Done")
}
