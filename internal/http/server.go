package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/HUARACHE7/orders-service-1/internal/cache"
	"github.com/HUARACHE7/orders-service-1/internal/config"
)

func StartServer(cache *cache.AppCache) {
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
			return
		}

		order, found := cache.Get(id)
		if !found {
			http.Error(w, "Order not found in cache", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println("HTTP-сервер запущен на http://localhost" + config.HTTPServerPort)
	if err := http.ListenAndServe(config.HTTPServerPort, nil); err != nil {
		log.Fatalf("HTTP-сервер упал: %v", err)
	}
}
