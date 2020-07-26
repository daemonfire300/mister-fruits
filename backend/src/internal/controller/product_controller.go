package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/daemonfire/mister-fruits/internal/connector"
)

type ProductController struct {
	ProductStore connector.ProductStore
}

func (c *ProductController) View(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	product, err := c.ProductStore.FindProduct(ctx, vars["id"])
	if err != nil {
		if errors.Is(err, connector.ErrProductNotFound) {
			w.WriteHeader(404)
			log.Println("Product not found")
			return
		}
		w.WriteHeader(500)
		log.Println("Could not fetch product", err)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(product); err != nil {
		log.Println("Error during encoding of product response", err)
	}
}

func (c *ProductController) ListAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := c.ProductStore.AllProducts(ctx)
	if err != nil {
		w.WriteHeader(500)
		log.Println("Could not fetch products", err)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(products); err != nil {
		log.Println("Error during encoding of products response", err)
	}
}
