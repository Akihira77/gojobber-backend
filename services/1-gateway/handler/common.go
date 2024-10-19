package handler

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func sendHttpReqToAnotherService(c *fiber.Ctx, url string) (int, []byte, []error) {
	a := fiber.AcquireAgent()
	a.Debug()

	req := a.Request()
	req.Header.SetMethod(c.Method())
	req.SetRequestURI(url)
	gatewayToken, _ := c.UserContext().Value("gatewayToken").(string)
	req.Header.Add("gatewayToken", gatewayToken)

	tokenStr := c.Cookies("token", "")
	if tokenStr == "" {
		authHeader := c.Get("Authorization", "")
		if authHeader != "" && len(strings.Split(authHeader, " ")) > 1 {
			tokenStr = strings.Split(authHeader, " ")[1]
		}
	}
	a.Cookie("token", tokenStr)

	if len(c.Body()) > 0 {
		a.ContentType(c.Get("Content-Type", "application/json"))
		req.Header.Add("Accept", c.Get("Accept", "*"))
		a.Body(c.Body())
	}

	if err := a.Parse(); err != nil {
		log.Printf("send http request to another service error:\n+%v", err)
		panic(err)
	}

	return a.Bytes()
}
