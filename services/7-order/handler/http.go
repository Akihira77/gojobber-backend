package handler

import "github.com/Akihira77/gojobber/services/7-order/service"

type OrderHttpHandler struct {
	orderSvc service.OrderServiceImpl
}

func NewOrderHttpHandler(orderSvc service.OrderServiceImpl) *OrderHttpHandler {
	return &OrderHttpHandler{
		orderSvc: orderSvc,
	}
}
