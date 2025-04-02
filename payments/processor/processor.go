package processor

import (
	pb "github.com/vp2305/common/api"
	stripeProcessor "github.com/vp2305/payments/processor/stripe"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}

func NewPaymentProcessor() PaymentProcessor {
	processor := stripeProcessor.NewProcessor()
	return processor
}
