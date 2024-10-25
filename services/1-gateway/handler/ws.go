package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Akihira77/gojobber/services/1-gateway/types"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type client struct {
	sync.Mutex
	wsConn *websocket.Conn
}

type wsReg struct {
	wsConn *websocket.Conn
	userId string
}

var (
	clients    = make(map[string]*client)
	register   = make(chan wsReg)
	unregister = make(chan wsReg)
	broadcast  = make(chan string)
)

func runHub() {
	for {
		select {
		case connection := <-register:
			clients[connection.userId] = &client{
				wsConn: connection.wsConn,
			}
			log.Println("connection registered")

		case connection := <-unregister:
			delete(clients, connection.userId)

			log.Println("connection unregistered")
		}
	}
}

func WsUpgrade(app fiber.Router) {
	app.Use("/ws", func(c *fiber.Ctx) error {
		log.Println("Client make a websocket upgrade request")

		if websocket.IsWebSocketUpgrade(c) {
			userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
			if !ok {
				return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
			}

			c.Locals("current_user", userInfo)
			return c.Next()
		}

		return fiber.NewError(http.StatusUpgradeRequired, "Can't establish Websocket connection")
	})

	go runHub()

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		u, ok := c.Locals("current_user").(*types.JWTClaims)
		if !ok {
			log.Println("current_user is invalid", u)
			return
		}

		defer func() {
			unregister <- wsReg{
				userId: u.UserID,
				wsConn: c,
			}
			c.Close()
		}()

		register <- wsReg{
			userId: u.UserID,
			wsConn: c,
		}

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return
			}

			type Message struct {
				ReceiverID string `json:"receiverId"`
				Message    string `json:"message"`
			}
			var msg Message
			_ = json.Unmarshal(message, &msg)

			targetConn, err := getClient(msg.ReceiverID)
			if err != nil {
				log.Printf("Sending message to [%s] error:\n+%v", msg.ReceiverID, err)
				err = c.Conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Failed sending message to [%s]", msg.ReceiverID)))
				if err != nil {
					log.Println("WebSocket writing message error", err)
				}

				continue
			}

			err = targetConn.wsConn.WriteMessage(websocket.BinaryMessage, []byte(msg.Message))
			if err != nil {
				log.Println("WebSocket writing message error", err)
			}
			log.Println("message from ws connection", message)
		}
	}))
}

func getClient(id string) (*client, error) {
	clientConn, ok := clients[id]
	if !ok {
		return nil, fmt.Errorf("Websocket Client with id [%s] did not exists", id)
	}

	return clientConn, nil
}

func SendMessage(senderId, receiverId string, data []byte) {
	receiverConn, err := getClient(receiverId)
	if err != nil {
		log.Println("SendMessage error", err)
		return
	}

	receiverConn.Mutex.Lock()
	defer receiverConn.Mutex.Unlock()

	if err := receiverConn.wsConn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		log.Printf("Sender [%s] sending data to [%s] through out ws is failed \n+%v", senderId, receiverId, err)

		// receiverConn.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
		// receiverConn.wsConn.Close()
		// unregister <- wsReg{
		// 	wsConn: receiverConn.wsConn,
		// 	userId: receiverId,
		// }
	}
}
