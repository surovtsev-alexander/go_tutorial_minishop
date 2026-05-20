// go_tutorial_minishop/user-service/main.go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var users = map[string]User{
	"1": {
		ID: "1", Name: "Alice", Email: "alice@example.com",
	},
	"2": {
		ID: "2", Name: "Bob", Email: "bob@example.com",
	},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if user, ok := users[id]; ok {
		json.NewEncoder(w).Encode(user)
	} else {
		http.Error(w, "User not found", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id}", getUser).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("User service запущен на :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
