package main

import (
	"encoding/json"
	"net/http"
	"os"
)

const filePath = "message.txt"

type Message struct {
	Content string `json:"content"`
}

func getMessage(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(filePath)
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "Could not read message", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(Message{Content: string(data)})
}

func updateMessage(w http.ResponseWriter, r *http.Request) {
	var msg Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := os.WriteFile(filePath, []byte(msg.Content), 0644); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Message updated successfully"})
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
