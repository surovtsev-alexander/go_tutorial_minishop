package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Order struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

var orders = map[string]Order{}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if order.ID == "" || order.ProductID == "" || order.Quantity <= 0 {
		http.Error(w, "Invalid order data", http.StatusBadRequest)
		return
	}

	// Вызов Product Service
	productResp, err := http.Get("http://product-service:8080/products/" + order.ProductID)
	if err != nil || productResp.StatusCode != http.StatusOK {
		http.Error(w, "Product not found", http.StatusBadRequest)
		return
	}

	var product struct {
		Price float64 `json:"price"`
	}
	json.NewDecoder(productResp.Body).Decode(&product)

	order.Total = float64(order.Quantity) * product.Price

	orders[order.ID] = order

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/orders", getOrders).Methods("GET")
	r.HandleFunc("/orders", createOrder).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Order service запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
