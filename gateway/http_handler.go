package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct {
	// gateway
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.HealthCheckHandler)
		r.Post("/api/customers/{customerID}/orders", h.HandleCreateOrder)
	})

	return r
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	log.Print("handle create order")
}

func (h *handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     "dev",
		"version": "0.0.1",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}
