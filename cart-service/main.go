package main

import (
	"log"
	"net/http"
	"os"

	"cart/internal/handler"
	"cart/internal/repository"
	"cart/internal/usecase"

	"github.com/go-redis/redis/v8"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	defer client.Close()

	repo := repository.NewCartRedisRepository(client)
	usecase := usecase.NewCartUsecase(repo)
	handler := handler.NewCartHandler(usecase)

	http.HandleFunc("/cart/add", handler.AddItem)
	http.HandleFunc("/cart", handler.GetCart)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("Cart service запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
