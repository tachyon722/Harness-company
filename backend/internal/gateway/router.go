package gateway

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", healthCheck)
	})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
		log.Printf("health check write error: %v", err)
	}
}
