package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

var cart Cart

func addItemHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	productID := r.URL.Query().Get("productId")
	quantity, _ := strconv.Atoi(r.URL.Query().Get("quantity"))
	price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

	if userID == "" || productID == "" || quantity <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Initialize cart if empty
	if cart.UserID == "" {
		cart.UserID = userID
	}

	// Add item to cart
	cart.Items = append(cart.Items, Item{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
	})

	w.WriteHeader(http.StatusCreated)
}

func getCartHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}

func main() {
	http.HandleFunc("/cart/add", addItemHandler)
	http.HandleFunc("/cart", getCartHandler)

	log.Println("Simple cart service запущен на :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}