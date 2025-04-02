package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/vp2305/common/api"
	"github.com/vp2305/common/broker"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment: %v", err)
				continue
			}

			log.Printf("Payment link created %s", paymentLink)
		}
	}()

	<-forever
}
