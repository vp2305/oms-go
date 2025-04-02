package main

import (
	"context"
	"testing"

	pb "github.com/vp2305/common/api"
	"github.com/vp2305/payments/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInmem()
	svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &pb.Order{})
		if err != nil {
			t.Errorf("CreatePayment() error = %v, want nil", err)
		}

		if link == "" {
			t.Errorf("CreatePayment() link is empty")
		}
	})
}
