package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/daemonfire/mister-fruits/internal/controller"
	"github.com/daemonfire/mister-fruits/internal/middleware"
)

// Middleware could be handled in a more abstract way
func Router(controllers *controller.Controllers, authenticationMiddleware *middleware.TokenGuardMiddleware) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/signup", controllers.UserController.Signup)
	r.HandleFunc("/login", controllers.UserController.Login)
	protectedResources := r.PathPrefix("/").Subrouter()
	protectedResources.Use(authenticationMiddleware.Handler())
	protectedResources.HandleFunc("/products/{id}", controllers.ProductController.View)
	protectedResources.HandleFunc("/products", controllers.ProductController.ListAll)
	protectedResources.HandleFunc("/carts/{id}/checkout", controllers.CartController.Checkout).Methods(http.MethodPost)
	protectedResources.HandleFunc("/carts/{id}/coupon/{couponId}", controllers.CartController.ApplyCoupon).Methods(http.MethodPost)
	protectedResources.HandleFunc("/carts/{id}", controllers.CartController.View).Methods(http.MethodGet)
	protectedResources.HandleFunc("/carts/{id}", controllers.CartController.Create).Methods(http.MethodPost)
	protectedResources.HandleFunc("/carts/{id}", controllers.CartController.Update).Methods(http.MethodPatch)
	protectedResources.HandleFunc("/coupons/generatedummy", controllers.CartController.DummyCoupon).Methods(http.MethodPost)
	frontendDir := "../../../../frontend"
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendDir)))
	return r
}
