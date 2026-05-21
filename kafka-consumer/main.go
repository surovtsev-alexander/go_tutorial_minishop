package main

import (
	"github.com/swaggo/http-swagger"
	"consumer/docs"
	"log"
	"net/http"
)

func main() {
	log.Println("Consumer stated")

	docs.SwaggerInfo.Title = "Kafka Consumer API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API for sending messages to Kafka"
	docs.SwaggerInfo.BasePath = "/"

	http.HandleFunc("/receive", addMessageHandler)

	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8086/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8086", nil))
}

// @Summary Receive message
// @Description Receive message from Kafka
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /receive [get]
func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "message received (stub)"}`))
}
