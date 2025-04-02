package main

import (
	"context"

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
			PriceID:  "price_1R9XiWRnokJdA0odrCdpCJoa",
			ID:       i.ID,
			Quantity: i.Quantity,
		})
	}

	return itemsWithPrice, nil
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
