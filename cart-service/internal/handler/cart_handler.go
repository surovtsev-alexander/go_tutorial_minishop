package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cart/internal/usecase"
)

type CartHandler struct {
	usecase *usecase.CartUsecase
}

func NewCartHandler(uc *usecase.CartUsecase) *CartHandler {
	return &CartHandler{usecase: uc}
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	productID := r.URL.Query().Get("productId")
	quantity, _ := strconv.Atoi(r.URL.Query().Get("quantity"))
	price, _ := strconv.ParseFloat(r.URL.Query().Get("price"), 64)

	if userID == "" || productID == "" || quantity <= 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	err := h.usecase.AddItem(userID, productID, quantity, price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	cart, err := h.usecase.GetCart(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cart)
}
