package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/vp2305/common"
	pb "github.com/vp2305/common/api"
)

type service struct {
	// store OrdersService
	store *store
}

func NewService(store *store) *service {
	return &service{store}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	o := &pb.Order{
		ID:         "1",
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}

	s.store.Create(ctx, o)
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, p *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(p.Items) == 0 {
		return nil, common.ErrNoItems
	}

	mergedItems := mergeItemQuantities(p.Items)

	// TODO: validate with the stock service

	// TODO: Remove Temporary
	var itemsWithPrice []*pb.Item
	for _, i := range mergedItems {
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			ID:       i.ID,
			Quantity: i.Quantity,
			PriceID:  "price_1R9XiWRnokJdA0odrCdpCJoa",
		})
	}

	return itemsWithPrice, nil
}

func (s *service) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	if p == nil {
		return nil, errors.New("request cannot be nil")
	}

	if p.OrderID == "" {
		return nil, errors.New("invalid request: order ID is required")
	}

	orderID, err := strconv.ParseInt(p.OrderID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID format")
	}
	if orderID <= 0 {
		return nil, errors.New("invalid request: order ID must be a positive integer")
	}

	if p.CustomerID == "" {
		return nil, errors.New("invalid request: customer ID is required")
	}

	customerID, err := strconv.ParseInt(p.CustomerID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid customer ID format")
	}
	if customerID <= 0 {
		return nil, errors.New("invalid request: customer ID must be a positive integer")
	}

	order, err := s.store.Get(ctx, p.OrderID, p.CustomerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

func mergeItemQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, item := range items {
		found := false

		for _, finalItem := range merged {
			if finalItem.ID == item.ID {
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, item)
		}
	}

	return merged
}
