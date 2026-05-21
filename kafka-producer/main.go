package main

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"producer/docs"
)

var kafkaWriter *kafka.Writer

func init() {
	kafkaWriter = &kafka.Writer{
		Addr:                   kafka.TCP("kafka:29092"),
		Topic:                  "messages",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true, // Разрешаем продюсеру запрашивать создание топика
		MaxAttempts:            5,    // Продюсер сделает до 5 попыток при ошибке Unknown Topic
	}
}

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

type SendMessageRequest struct {
	Message string `json:"message"`
}

// @Summary Send message
// @Description Send message to Kafka
// @Tags messages
// @Accept json
// @Produce json
// @Param request body SendMessageRequest true "Message request"
// @Success 200 {object} map[string]string
// @Router /send [post]
func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req SendMessageRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		http.Error(w, "Message cannot be empty", http.StatusBadRequest)
		return
	}
	message := kafka.Message{
		Value: []byte(req.Message),
	}
	err := kafkaWriter.WriteMessages(context.Background(), message)

	if err != nil {
		log.Printf("Failed to write message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "message sent",
		"message": req.Message,
	})
}
