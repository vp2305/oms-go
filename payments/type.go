package main

import (
	"context"

	pb "github.com/vp2305/common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
