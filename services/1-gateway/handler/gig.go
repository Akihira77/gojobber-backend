package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type GigHandler struct {
	base_url string
}

func NewGigHandler(base_url string) *GigHandler {
	return &GigHandler{
		base_url: base_url,
	}
}

func (gh *GigHandler) HealthCheck(c *fiber.Ctx) error {
	route := gh.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})
}

func (gh *GigHandler) FindGigByID(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/id/%s", c.Params("id"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - find gig by id error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) FindSellerActiveGigs(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/seller/%s/%s/%s", c.Params("sellerId"), c.Params("page"), c.Params("size"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - find seller active gigs error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) FindSellerInactiveGigs(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/seller/inactive/%s/%s/%s", c.Params("sellerId"), c.Params("page"), c.Params("size"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - find seller inactive gigs error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) FindGigsByCategory(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/category/%s/%s/%s", c.Params("category"), c.Params("page"), c.Params("size"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - find gigs by category error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) GetPopularGigs(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/popular/category/%s", c.Params("category"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - get popular gigs by category error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) FindSimilarGigs(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/similar/%s/%s/%s", c.Params("gigId"), c.Params("page"), c.Params("size"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - find similar gigs error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) GigQuerySearch(c *fiber.Ctx) error {
	query := fmt.Sprintf("query=%v&max=%v&delivery_time=%v", c.Query("query"), c.QueryInt("max"), c.QueryInt("delivery_time"))
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/search/%s/%s?%s", c.Params("page"), c.Params("size"), query)
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - get gigs by query error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) CreateGig(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - creating gig error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) UpdateGig(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/%s/%s", c.Params("sellerId"), c.Params("gigId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - updating gig error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) ActivateGigStatus(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/update-status/%s/%s", c.Params("sellerId"), c.Params("gigId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - activate gig status error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (gh *GigHandler) DeactivateGigStatus(c *fiber.Ctx) error {
	route := gh.base_url + fmt.Sprintf("/api/v1/gig/%s/%s", c.Params("sellerId"), c.Params("gigId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("GIG - deactivate gig status error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
