package main

import (
	"github.com/swaggo/http-swagger"
	"producer/docs"
	"log"
	"net/http"
	"github.com/segmentio/kafka-go"
	"context"
)

var kafkaWriter *kafka.Writer

func init() {
    kafkaWriter = &kafka.Writer{
        Addr: kafka.TCP("kafka:29092"),
        Topic: "messages",
        Balancer:&kafka.LeastBytes{},
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

// @Summary Send message
// @Description Send message to Kafka
// @Tags messages
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /send [post]
func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	message := kafka.Message{
	    Value: []byte(`{"message": "test message"}`),
	}
	err := kafkaWriter.WriteMessages(context.Background(), message)

	if err != nil {
        log.Printf("Failed to write message: %v", err)
        http.Error(w, "Failed to send message", http.StatusInternalServerError)
        return
    }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "message received (stub)"}`))
}
