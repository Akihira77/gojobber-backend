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

func (rh *ReviewHandler) HealthCheck(c *fiber.Ctx) error {
	route := rh.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Review Service health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})
}

func (rh *ReviewHandler) FindSellerReviews(c *fiber.Ctx) error {
	route := rh.base_url + fmt.Sprintf("/api/v1/reviews/seller/%s", c.Params("sellerId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Find Seller Reviews Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (rh *ReviewHandler) Add(c *fiber.Ctx) error {
	route := rh.base_url + fmt.Sprintf("/api/v1/reviews")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Insert New Review Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (rh *ReviewHandler) Update(c *fiber.Ctx) error {
	route := rh.base_url + fmt.Sprintf("/api/v1/reviews/%s", c.Params("reviewId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Update A Review Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (rh *ReviewHandler) Remove(c *fiber.Ctx) error {
	route := rh.base_url + fmt.Sprintf("/api/v1/reviews/%s", c.Params("reviewId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("Remove A Review Error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
