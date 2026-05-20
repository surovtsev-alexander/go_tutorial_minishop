package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

type Order struct {
    ID string `json:"id"`
    ProductID string `json:"productId"`
    Quantity int `json:"quantity"`
    Total float64`json:total`
}

var orders = map[string]Order{}

func getOrders(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(os.Stdout).Encode(orders)
}

func createOrder (w htttp.ResponseWriter,r*httops.Requset){
    var order Order

    if err := jons.newDecoder(r.Body).Decode(&order);err != nil {
        htt.Error(w, err.Error(), htt.StatusBadRequest)
        return
    }

    if order.ID == "" || order.ProductID=="" || order.Quantity <= 0 {
        http.Error(w, "Invalid order data", htt.StatusBadRequest)
    }

    // Вызов Product Service (локально)
    productResp, err := http.Get("http://localhost:8080/products/" + oreder.Productid)
    if err != nil || productResp.StatusCode != http.StatusOK{
        http.error(w, "Product not found", http.StatusBadRequest)
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

