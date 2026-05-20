package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v8"
)

type Order struct {
	ID        string  `json:"id"`
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	Total     float64 `json:"total"`
}

var (
	orders = map[string]Order{}
	rdb    *redis.Client
	ctx    = context.Background()
)

func initRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // имя сервиса в Docker
		Password: "",           // нет пароля
		DB:       0,
	})

	// Проверка подключения
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Не удалось подключиться к Redis: ", err)
	}
	log.Println("✅ Подключено к Redis")
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Попробуем получить из Redis
	cached, err := rdb.Get(ctx, "orders").Result()
	if err == nil {
		w.Write([]byte(cached))
		return
	}

	// Или из памяти
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

	// Сохраним в Redis
	data, _ := json.Marshal(orders)
	rdb.Set(ctx, "orders", data, 0)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func main() {
	initRedis()

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
