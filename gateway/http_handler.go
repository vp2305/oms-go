package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vp2305/common"
	pb "github.com/vp2305/common/api"
)

type handler struct {
	// gateway

	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client}
}

func (h *handler) registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.HealthCheckHandler)

		r.Post("/customers/{customerID}/orders", h.HandleCreateOrder)
	})

	return r
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})
}

func (h *handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     "dev",
		"version": "0.0.1",
	}

	common.WriteJSON(w, http.StatusOK, data)
}
