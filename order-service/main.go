package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/go-redis/redis/v8"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "order-service/docs" // импорт для Swagger
)

// Order представляет заказ
// swagger:model
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
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Не удалось подключиться к Redis: ", err)
	}
	log.Println("✅ Подключено к Redis")
}

func getOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cached, err := rdb.Get(ctx, "orders").Result()
	if err == nil {
		w.Write([]byte(cached))
		return
	}

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

	data, _ := json.Marshal(orders)
	rdb.Set(ctx, "orders", data, 0)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// @title Order Service API
// @version 1.0
// @description API для управления заказами
// @host localhost:8081
// @BasePath /api/v1
func main() {
	initRedis()

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// @Summary Получить все заказы
	// @Description Возвращает список всех заказов
	// @Tags orders
	// @Produce json
	// @Success 200 {array} Order
	// @Router /orders [get]
	api.HandleFunc("/orders", getOrders).Methods("GET")

	// @Summary Создать заказ
	// @Description Создаёт новый заказ
	// @Tags orders
	// @Accept json
	// @Produce json
	// @Param order body Order true "Данные заказа"
	// @Success 201 {object} Order
	// @Failure 400 {object} map[string]string
	// @Router /orders [post]
	api.HandleFunc("/orders", createOrder).Methods("POST")

	// Подключаем Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Order service запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
