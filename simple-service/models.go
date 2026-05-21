package main

// Cart represents a user's shopping cart
type Cart struct {
	UserID string `json:"userId"`
	Items  []Item `json:"items"`
}

// Item represents an item in the cart
type Item struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
