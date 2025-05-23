package gateway

import (
	"context"

	pb "github.com/vp2305/common/api"
)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, orderID string, customerID string) (*pb.Order, error)
}
