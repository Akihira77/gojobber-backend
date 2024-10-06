package handler

import (
	"encoding/json"
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

	type InsertMsgRes struct {
		SenderID string `json:"senderId"`
		Receiver struct {
			ID             string `json:"id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profilePicture"`
		} `json:"receiver"`
		UnreadMessages int `json:"unreadMessages"`
	}

	var res InsertMsgRes
	_ = json.Unmarshal(body, &res)
	msg := struct {
		Topic    string `json:"topic"`
		SenderID string `json:"senderId"`
		Receiver struct {
			ID             string `json:"id"`
			Username       string `json:"username"`
			ProfilePicture string `json:"profilePicture"`
		} `json:"receiver"`
		UnreadMessages int `json:"unreadMessages"`
	}{
		Topic:          "Chat-Notification",
		SenderID:       res.SenderID,
		Receiver:       res.Receiver,
		UnreadMessages: res.UnreadMessages,
	}
	b, _ := json.Marshal(msg)
	go SendMessage(res.SenderID, res.Receiver.ID, b)

	return c.Status(statusCode).Send(body)
}
