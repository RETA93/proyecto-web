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

	// Cambiamos este endpoint para validar el nuevo Hola Mundo en el servidor
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "¡Hola Mundo desde el Backend de Go en la ThinkCentre!",
			"status":  "healthy",
		})
	})

	http.ListenAndServe(":8080", nil)
}
