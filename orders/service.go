package main

import (
	"context"
	"errors"

	"github.com/vp2305/common"
	pb "github.com/vp2305/common/api"
)

type service struct {
	// store OrdersService
}

func NewService(store *store) *service {
	return &service{}
}

func (s *service) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := s.ValidateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	o := &pb.Order{
		ID:         "42",
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items:      items,
	}

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
	if p.OrderID == "" {
		return nil, errors.New("order id cannot be empty")
	}

	if p.CustomerID == "" {
		return nil, errors.New("customer id cannot be empty")
	}

	order := &pb.Order{
		ID:         p.OrderID,
		CustomerID: p.CustomerID,
		Status:     "pending",
		Items: []*pb.Item{
			{
				ID:       "1",
				Quantity: 2,
				PriceID:  "price_1R9XiWRnokJdA0odrCdpCJoa",
			},
			{
				ID:       "2",
				Quantity: 5,
				PriceID:  "price_1R9XiWRnokJdA0odrCdpCJoa",
			},
		},
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
