package main

import (
	"consumer/docs"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"os"
	"time"
)

var kafkaReader *kafka.Reader

func init() {
	broker := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if broker == "" {
		broker = "kafka:29092"
	}

	kafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    "messages",
		GroupID:  "consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})
}

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
	// Создаем контекст с таймаутом 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var messages []string

	for i := 0; i < 3; i++ {
		m, err := kafkaReader.ReadMessage(ctx)
		if err != nil {
			// Если контекст отменен (таймаут), выходим из цикла
			if errors.Is(err, context.DeadlineExceeded) {
				log.Printf("Timeout reading messages: %v", err)
				break
			}
			log.Printf("Failed to read message: %v", err)
			break
		}
		messages = append(messages, string(m.Value))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "success",
		"messages": messages,
		"count":    len(messages),
	})
}
