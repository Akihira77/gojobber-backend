package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	base_url string
}

func NewUserHandler(base_url string) *UserHandler {
	return &UserHandler{
		base_url: base_url,
	}
}

func (uh *UserHandler) HealthCheck(c *fiber.Ctx) error {
	route := uh.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})
}

func (uh *UserHandler) GetMyBuyerInfo(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/buyers/my-info")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - get my buyer info error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) FindBuyerByUsername(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/buyers/%s", c.Params("username"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - find buyer by username error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) GetMySellerInfo(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers/my-info")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - get my seller info error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) FindSellerByID(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers/id/%s", c.Params("id"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - find seller by id error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) FindSellerByUsername(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers/username/%s", c.Params("username"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - find seller by username error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) GetRandomSellers(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers/random/%s", c.Params("count"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - get random sellers error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) Create(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - creating seller error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (uh *UserHandler) Update(c *fiber.Ctx) error {
	route := uh.base_url + fmt.Sprintf("/api/v1/users/sellers")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("USER - updating seller error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
