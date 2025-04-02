package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v82"
	"github.com/vp2305/common"
	"github.com/vp2305/common/broker"
	"github.com/vp2305/common/discovery"
	"github.com/vp2305/common/discovery/consul"
	"github.com/vp2305/payments/processor"
	"google.golang.org/grpc"
)

var (
	serviceName = "payments"
	grpcAddr    = common.EnvString("GRPC_ADDR", "localhost:2001")
	consulAddr  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	amqpUser    = common.EnvString("RABBITMQ_USER", "guest")
	amqpPass    = common.EnvString("RABBITMQ_PASS", "guest")
	amqpHost    = common.EnvString("RABBITMQ_HOST", "localhost")
	amqpPort    = common.EnvString("RABBITMQ_PORT", "5672")
	stripeKey   = common.EnvString("STRIPE_KEY", "")
)

func main() {
	// Register a service
	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}

	go func() {
		for {
			if err := registry.HealthCheck(instanceID, serviceName); err != nil {
				log.Fatal("failed to health check")
			}

			time.Sleep(time.Second * 1)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	// stripe setup
	stripe.Key = stripeKey

	// Broker connection
	ch, close := broker.Connect(amqpUser, amqpPass, amqpHost, amqpPort)
	defer func() {
		close()
		ch.Close()
	}()

	paymentProcessor := processor.NewPaymentProcessor()
	svc := NewService(paymentProcessor)
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	// gRPC server
	grpcServer := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	log.Println("GRPC Server Started at", grpcAddr)
	if err := grpcServer.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
