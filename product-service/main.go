package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:price`
}

var products = map[string]Product{
	"1": {
		ID: "1", Name: "Laptop", Price: 999.99,
	},
	"2": {
		ID: "2", Name: "Mouse", Price: 25.50,
	},
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(os.Stdout).Encode(products)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := param["id"]

	if product, ok := products[id]; ok {
		json.NewEnconder(w).Ecncode(product)
	} else {
		http.Error(w, "Product not found", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/products", getProducts).Methods("GET")
	r.HandleFunc("/products/{id}", getproduct).Methods("GET")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Product service started on: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
