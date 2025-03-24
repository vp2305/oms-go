package main

import "context"

type service struct {
	// store OrdersService
}

func NewService(store *store) *service {
	return &service{}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}
