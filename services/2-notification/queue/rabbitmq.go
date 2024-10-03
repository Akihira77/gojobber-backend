package queue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Akihira77/gojobber/services/2-notification/helper"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

// Connection is the connection created
type Connection struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	err      chan error
	validate *validator.Validate
}

func NewConnection() *Connection {
	amqpURI := os.Getenv("RABBITMQ_URI")
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		log.Fatalf("Error in creating rabbitmq connection with %s : %s", amqpURI, err.Error())
	}

	go func() {
		<-conn.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		log.Println("Connection closed")
	}()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Channel: %s", err)
	}

	log.Println("RabbitMQ connection established", conn.LocalAddr())
	return &Connection{
		conn:     conn,
		channel:  channel,
		err:      make(chan error),
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (c *Connection) ConsumeFromAuthService() {
	routing, exchange := "auth_to_notification", "auth_to_notification"
	err := c.channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring exchange :\n%+v", err)
		return
	}

	q, err := c.channel.QueueDeclare(routing, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring queue :\n%+v", err)
		return
	}

	err = c.channel.QueueBind(q.Name, routing, exchange, false, nil)
	if err != nil {
		fmt.Printf("error binding queue :\n%+v", err)
		return
	}

	msgs, err := c.channel.Consume(q.Name, routing, false, false, false, false, nil)
	if err != nil {
		fmt.Printf("error consuming message :\n%+v", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("receive message from [%s] route, body: %v\n", msg.RoutingKey, string(msg.Body))

			topic := msg.Headers["Topic"]
			errCh := make(chan error, 1)
			switch topic {
			case "user-verifying-email":
				var val struct {
					ReceiverEmail string `validate:"required,email"`
					Template      string `validate:"required"`
					VerifyLink    string `validate:"required"`
				}

				err := json.Unmarshal(msg.Body, &val)
				if err != nil {
					fmt.Printf("userVerifyingEmail error:\n%+v", err)
					break
				}

				err = c.validate.Struct(&val)
				if err != nil {
					fmt.Printf("userVerifyingEmail error:\n%+v", err)
					break
				}

				go helper.VerifyAccountURLMail(errCh, val.ReceiverEmail, "Verify Account URL", val.VerifyLink)
				if err = <-errCh; err != nil {
					fmt.Printf("userVerifyingEmail error:\n%+v", err)
					break
				}

				break
			case "user-forgot-password":
				var val struct {
					ReceiverEmail string `validate:"required,email"`
					Template      string `validate:"required"`
					Username      string `validate:"required"`
					ResetLink     string `validate:"required"`
				}

				err := json.Unmarshal(msg.Body, &val)
				if err != nil {
					fmt.Printf("userForgotPassword error:\n%+v", err)
					break
				}

				err = c.validate.Struct(&val)
				if err != nil {
					fmt.Printf("userForgotPassword error:\n%+v", err)
					break
				}

				go helper.ForgotPasswordMail(errCh, val.ReceiverEmail, "Reset Password URL", val.ResetLink, val.Username)
				if err = <-errCh; err != nil {
					fmt.Printf("userForgotPassword error:\n%+v", err)
					break
				}

				break
			case "user-reset-password-success":
				var val struct {
					ReceiverEmail string `validate:"required,email"`
					Template      string `validate:"required"`
					Username      string `validate:"required"`
				}

				err := json.Unmarshal(msg.Body, &val)
				if err != nil {
					fmt.Printf("userResetPassword error:\n%+v", err)
					break
				}

				err = c.validate.Struct(&val)
				if err != nil {
					fmt.Printf("userResetPassword error:\n%+v", err)
					break
				}

				go helper.ResetPasswordSuccessMail(errCh, val.ReceiverEmail, "Success Reseting Your Password", val.Username)
				if err = <-errCh; err != nil {
					fmt.Printf("userResetPassword error:\n%+v", err)
					break
				}

				break
			default:
				log.Println("unhandled message topic", topic, msg.RoutingKey)
				break
			}

			if err = msg.Ack(false); err != nil {
				fmt.Printf("Message from route [%s] error:\n%+v", msg.RoutingKey, err)
			}
		}
	}()

	log.Printf(" [*] Waiting for logs from Auth Service. To exit press CTRL+C")
	<-forever
}

func (c *Connection) ConsumeFromChatService() {
	routing, exchange := "chat_to_notification", "chat_to_notification"
	err := c.channel.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring exchange :\n%+v", err)
		return
	}

	q, err := c.channel.QueueDeclare(routing, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring queue :\n%+v", err)
		return
	}

	err = c.channel.QueueBind(q.Name, routing, exchange, false, nil)
	if err != nil {
		fmt.Printf("error binding queue :\n%+v", err)
		return
	}

	msgs, err := c.channel.Consume(q.Name, routing, false, false, false, false, nil)
	if err != nil {
		fmt.Printf("error consuming message :\n%+v", err)
		return
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {
			log.Printf("receive message from [%s] route, body: %v\n", msg.RoutingKey, string(msg.Body))

			topic := msg.Headers["Topic"]
			errCh := make(chan error, 1)
			switch topic {
			case "email-chat-notification":
				var val struct {
					ReceiverEmail string `validate:"required,email"`
					SenderEmail   string `validate:"required,email"`
					Message       string `validate:"required"`
				}

				err := json.Unmarshal(msg.Body, &val)
				if err != nil {
					fmt.Printf("emailChatNotification error:\n%+v", err)
					close(errCh)
					break
				}

				err = c.validate.Struct(&val)
				if err != nil {
					fmt.Printf("emailChatNotification error:\n%+v", err)
					close(errCh)
					break
				}

				go helper.SendMail(val.ReceiverEmail, fmt.Sprintf("You receive message from user: %s", val.SenderEmail), val.Message)
				if err = <-errCh; err != nil {
					fmt.Printf("emailChatNotification error:\n%+v", err)
					close(errCh)
					break
				}

				close(errCh)
				break
			default:
				log.Println("unhandled message topic", topic, msg.RoutingKey)
				close(errCh)
				break
			}

			if err = msg.Ack(false); err != nil {
				fmt.Printf("Message from route [%s] error:\n%+v", msg.RoutingKey, err)
			}
		}
	}()
	log.Printf(" [*] Waiting for logs from Chat Service. To exit press CTRL+C")
	<-forever
}
