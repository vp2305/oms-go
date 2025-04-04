package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vp2305/common"
	pb "github.com/vp2305/common/api"
	"github.com/vp2305/gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	gateway gateway.OrdersGateway
}

func NewHandler(gateway gateway.OrdersGateway) *handler {
	return &handler{gateway: gateway}
}

func (h *handler) registerRoutes() http.Handler {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", h.healthCheckHandler)

		r.Route("/customers", func(r chi.Router) {
			r.Post("/{customerID}/orders", h.handleCreateOrder)
			r.Get("/{customerID}/orders/{orderID}", h.handleGetOrder)
		})
	})

	return r
}

func (h *handler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateItems(items); err != nil {
		common.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	order, err := h.gateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerID: customerID,
		Items:      items,
	})

	errStatus := status.Convert(err)

	if errStatus != nil {
		if errStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusCreated, order)
}

func (h *handler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customerID")
	orderID := r.PathValue("orderID")

	order, err := h.gateway.GetOrder(r.Context(), orderID, customerID)

	errStatus := status.Convert(err)

	if errStatus != nil {
		if errStatus.Code() != codes.InvalidArgument {
			common.WriteError(w, http.StatusBadRequest, errStatus.Message())
			return
		}

		common.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	common.WriteJSON(w, http.StatusOK, order)
}

func (h *handler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     "dev",
		"version": "0.0.1",
	}

	common.WriteJSON(w, http.StatusOK, data)
}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, i := range items {
		if i.ID == "" {
			return errors.New("item ID is required")
		}

		if i.Quantity <= 0 {
			return errors.New("items must have a valid quantity")
		}
	}

	return nil
}
