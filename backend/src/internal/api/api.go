package api

import (
	"github.com/gorilla/mux"

	"github.com/daemonfire/mister-fruits/internal/controller"
)

func Router(controllers *controller.Controllers) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controllers.UserController.Signup)
	r.HandleFunc("/login", controllers.UserController.Login)
	/*r.HandleFunc("/products/{id}", controllers.ProductController.View)
	r.HandleFunc("/products", controllers.ProductController.ListAll)
	r.HandleFunc("/carts/{id}/checkout", controllers.CartController.Checkout).Methods(http.MethodPost)
	r.HandleFunc("/carts/{id}/coupon/{couponId}", controllers.CartController.Checkout).Methods(http.MethodPost)
	r.HandleFunc("/carts/{id}", controllers.CartController.View).Methods(http.MethodGet)
	r.HandleFunc("/carts/{id}", controllers.CartController.Create).Methods(http.MethodPost)
	r.HandleFunc("/carts/{id}", controllers.CartController.Update).Methods(http.MethodPatch)*/
	return r
}
