package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Producer started")

	http.HandleFunc("/send", addMessageHandler)

	log.Fatal(http.ListenAndServe(":8085", nil))
}

func addMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(http.StatusOk)
	w.Write([]byte(`{"status": "message received (stub)"}`))
}
