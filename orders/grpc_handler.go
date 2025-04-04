package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/vp2305/common/api"
	"github.com/vp2305/common/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer

	service OrdersService
	channel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrdersService, ch *amqp.Channel) {
	handler := &grpcHandler{
		service: service,
		channel: ch,
	}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order received! Order %+v", p)

	o, err := h.service.CreateOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	q, err := h.channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	h.channel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})

	return o, nil
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	log.Printf("Get order information! %+v", p)

	order, err := h.service.GetOrder(ctx, p)
	if err != nil {
		return nil, err
	}

	return order, nil
}
