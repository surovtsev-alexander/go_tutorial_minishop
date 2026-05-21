package main

import (
	"github.com/swaggo/http-swagger"
	"producer/docs"
	"log"
	"net/http"
)

func main() {
	log.Println("Producer started")

	docs.SwaggerInfo.Title = "Kafka Producer API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Description = "API for sending messages to Kafka"
	docs.SwaggerInfo.BasePath = "/"

	http.HandleFunc("/send", addMessageHandler)

	http.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8085/swagger/doc.json"),
	))

	log.Fatal(http.ListenAndServe(":8085", nil))
}

// @Summary Send message
// @Description Send message to Kafka
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /send [post]
func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "message received (stub)"}`))
}
