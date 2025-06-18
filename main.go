package main

import (
	"encoding/json"
	"net/http"
	"os"
)

const messageFilePath = "message.txt"
const userFilePath = "user.txt"

type Message struct {
	Content string `json:"content"`
	User    string `json:"user"`
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(messageFilePath)
	user, userErr := os.ReadFile(userFilePath)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Could not read message", http.StatusInternalServerError)
		return
	}
	if userErr != nil && !os.IsNotExist(userErr) {
		http.Error(w, "Could not read user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(Message{Content: string(data), User: string(user)})
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if msg.Content == "" {
		http.Error(w, "Your request body must include a 'user' field and a 'content' field", http.StatusBadRequest)
		return
	}
	if err := os.WriteFile(messageFilePath, []byte(msg.Content), 0644); err != nil {
		http.Error(w, "Failed to write message file", http.StatusInternalServerError)
		return
	}
	if userErr := os.WriteFile(userFilePath, []byte(msg.User), 0644); userErr != nil {
		http.Error(w, "Failed to write user file", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Message and user updated successfully"})
}

func main() {
	http.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			getMessage(w, r)
		case http.MethodPost:
			updateMessage(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
