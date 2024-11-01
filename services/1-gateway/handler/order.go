package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	base_url string
}

func NewOrderHandler(base_url string) *OrderHandler {
	return &OrderHandler{
		base_url: base_url,
	}
}

func (oh *OrderHandler) HealthCheck(c *fiber.Ctx) error {
	route := oh.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Auth health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})
}

func (oh *OrderHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprint("/api/v1/orders/payment-intents/create")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Create Payment Intent error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) ConfirmPayment(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/payment-intents/%s/confirm", c.Params("paymentId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Confirm Payment error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) HandleStripeWebhook(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/stripe/webhook")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Payment Confirm Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) StripeTOSAcceptance(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/stripe/tos-acceptance")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Payment Confirm Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) RequestDeadlineExtension(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/deadline/extension/%s/request", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Request Deadline Extension Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) BuyerDeadlineExtensionResponse(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/deadline/extension/%s/response", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Deadline Extension Response Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) FindOrderByID(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/%s", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Find Order By ID Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) FindOrdersByBuyerID(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/buyer/my-orders")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Find My Orders As Buyer Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) FindOrdersBySellerID(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/seller/my-orders")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Find My Orders As Seller Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) OrderComplete(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/%s/complete", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Order Completion Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) CancelOrder(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/%s/cancel", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Order Cancelation Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) OrderRefund(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/%s/refund", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Order Refunds Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) AcknowledgeOrder(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/%s/acknowledge", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Acknowledging Order Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) BuyerResponseForDeliveredOrder(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/deliver/%s/response", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Buyer Responding Delivered Order Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) DeliveringOrder(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/deliver/%s", c.Params("orderId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Seller Delivering Order Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (oh *OrderHandler) FindMyOrdersNotifications(c *fiber.Ctx) error {
	route := oh.base_url + fmt.Sprintf("/api/v1/orders/buyer/my-orders-notifications")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Find My Orders On Notifications Bar Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
