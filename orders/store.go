package main

import (
	"context"
	"fmt"

	pb "github.com/vp2305/common/api"
)

var orders = make([]*pb.Order, 0)

type store struct {
	// Add here our mongoDB
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, newOrder *pb.Order) error {
	orders = append(orders, newOrder)

	return nil
}

func (s *store) Get(ctx context.Context, id, customerID string) (*pb.Order, error) {
	for _, o := range orders {
		if o.ID == id && o.CustomerID == customerID {
			return o, nil
		}
	}

	return nil, fmt.Errorf("order not found by id %s", id)
}
