package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		env := os.Getenv("ENV")
		if env == "" {
			env = "dev"
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "env": env})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"message": "Backend Go funcionando"})
	})

	http.ListenAndServe(":8080", nil)
}
