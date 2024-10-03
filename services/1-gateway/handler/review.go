package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ReviewHandler struct {
	base_url string
}

func NewReviewHandler(base_url string) *ReviewHandler {
	return &ReviewHandler{
		base_url: base_url,
	}
}

func (ah *ReviewHandler) HealthCheck(c *fiber.Ctx) error {
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
