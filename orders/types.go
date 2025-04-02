package main

import (
	"context"

	pb "github.com/vp2305/common/api"
)

type OrdersService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
}

type OrdersStore interface {
	Create(context.Context) error
}
