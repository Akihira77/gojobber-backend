package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ChatHandler struct {
	base_url string
}

func NewChatHandler(base_url string) *ChatHandler {
	return &ChatHandler{
		base_url: base_url,
	}
}

func (ch *ChatHandler) HealthCheck(c *fiber.Ctx) error {
	route := ch.base_url + "/health-check"
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("CHAT - health check error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"response": string(body),
	})
}

func (ch *ChatHandler) GetAllMyConversations(c *fiber.Ctx) error {
	route := ch.base_url + fmt.Sprintf("/api/v1/chat/my-conversations")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("CHAT - find all my conversations error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ch *ChatHandler) GetMessagesInsideConversation(c *fiber.Ctx) error {
	route := ch.base_url + fmt.Sprintf("/api/v1/chat/id/%s", c.Params("conversationId"))
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("CHAT - get messages inside conversation error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}

func (ch *ChatHandler) InsertMessage(c *fiber.Ctx) error {
	route := ch.base_url + fmt.Sprintf("/api/v1/chat")
	statusCode, body, errs := sendHttpReqToAnotherService(c, route)
	if len(errs) > 0 {
		fmt.Println("CHAT - inserting message error", errs)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errs": errs,
		})
	}

	return c.Status(statusCode).Send(body)
}
