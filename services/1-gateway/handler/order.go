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

func (ah *OrderHandler) HealthCheck(c *fiber.Ctx) error {
	route := ah.base_url + "/health-check"
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
