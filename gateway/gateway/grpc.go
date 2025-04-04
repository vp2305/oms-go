package gateway

import (
	"context"

	pb "github.com/vp2305/common/api"
	"github.com/vp2305/common/discovery"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway {
	return &gateway{registry: registry}
}

func (g *gateway) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)

	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: p.CustomerID,
		Items:      p.Items,
	})
}

func (g *gateway) GetOrder(ctx context.Context, orderID, customerID string) (*pb.Order, error) {
	conn, err := discovery.ServiceConnection(ctx, "orders", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)

	return c.GetOrder(ctx, &pb.GetOrderRequest{
		OrderID:    orderID,
		CustomerID: customerID,
	})
}
