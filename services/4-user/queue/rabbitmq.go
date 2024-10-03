package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Akihira77/gojobber/services/4-user/service"
	svc "github.com/Akihira77/gojobber/services/4-user/service"
	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/go-playground/validator/v10"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// MessageBody is the struct for the body passed in the AMQP message. The type will be set on the Request header
type MessageBody struct {
	Data []byte
	Type string
}

// Connection is the connection created
type Connection struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	err       chan error
	buyerSvc  svc.BuyerServiceImpl
	sellerSvc svc.SellerServiceImpl
	validate  *validator.Validate
}

func NewConnection(db *gorm.DB) *Connection {
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
		conn:      conn,
		channel:   channel,
		err:       make(chan error),
		buyerSvc:  service.NewBuyerService(db),
		sellerSvc: service.NewSellerService(db),
		validate:  validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (c *Connection) ConsumeFromAuthService() {
	routing, exchange := "user_from_auth", "user_from_auth"
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
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			topic := msg.Headers["Topic"]
			log.Printf("receive message from [%s] route, [%s] topic, body: %v\n", msg.RoutingKey, topic, string(msg.Body))

			switch topic {
			case "registration-account":
				var buyer types.Buyer
				res := types.RabbitMQResponse[string]{
					Success: false,
					Data:    "",
				}
				err := json.Unmarshal(msg.Body, &buyer)
				respBody := amqp.Publishing{
					CorrelationId: msg.CorrelationId,
					ReplyTo:       msg.ReplyTo,
				}
				if err != nil {
					log.Printf("Error unmarshal buyer registration data")

					res.Data = "Error unmarshal buyer registration data"
					b, _ := json.Marshal(&res)
					respBody.Body = b
					err = c.channel.PublishWithContext(ctx, "auth_from_user_RPC", "auth_from_user_RPC", false, false, respBody)
					if err != nil {
						log.Println(err)
					}

					break
				}

				err = c.validate.Struct(&buyer)
				if err != nil {
					log.Println("validation result:", err)

					res.Data = "Error validating auth data"
					b, _ := json.Marshal(&res)
					respBody.Body = b
					err = c.channel.PublishWithContext(ctx, "auth_from_user_RPC", "auth_from_user_RPC", false, false, respBody)
					if err != nil {
						log.Println(err)
					}

					break
				}

				err = c.buyerSvc.Create(ctx, buyer)
				if err != nil {
					log.Printf("Error saving buyer data")

					res.Data = "Error saving buyer data"
					b, _ := json.Marshal(&res)
					respBody.Body = b
					err = c.channel.PublishWithContext(ctx, "auth_from_user_RPC", "auth_from_user_RPC", false, false, respBody)
					if err != nil {
						log.Println(err)
					}

					break
				}

				res.Success = true
				res.Data = "Done"
				b, _ := json.Marshal(&res)
				respBody.Body = b
				err = c.channel.PublishWithContext(ctx, "auth_from_user_RPC", "auth_from_user_RPC", false, false, respBody)
				if err != nil {
					log.Println(err)
					break
				}

				log.Println("message published")
				break
			case "rollback-registration-account":
				type Request struct {
					UserID string `json:"userId"`
				}
				var val Request
				err := json.Unmarshal(msg.Body, &val)
				if err != nil {
					fmt.Println("Error unmarshaling on event rollback-registration-account")
					break
				}

				c.buyerSvc.Delete(ctx, val.UserID)
				break
			default:
				log.Printf("unknown message from route [%s] and topic [%s]\n", msg.RoutingKey, topic)
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

func (c *Connection) ConsumeFromGigService() {
	routing, exchange := "user_from_gig", "user_from_gig"
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
			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			topic := msg.Headers["Topic"]
			log.Printf("receive message from [%s] route, correlationId [%s], [%s] topic, body: %v\n", msg.RoutingKey, msg.CorrelationId, topic, string(msg.Body))

			switch topic {
			case "retrieving-seller":
				type Request struct {
					SellerID string `json:"sellerId"`
				}
				var val Request
				_ = json.Unmarshal(msg.Body, &val)

				res := types.RabbitMQResponse[*types.SellerOverview]{
					Success: false,
					Data:    nil,
				}
				r := amqp.Publishing{
					CorrelationId: msg.CorrelationId,
					ContentType:   "application/json",
					ReplyTo:       msg.ReplyTo,
				}
				s, err := c.sellerSvc.FindSellerOverviewByID(ctx, val.SellerID)
				if err != nil {
					b, _ := json.Marshal(&res)
					r.Body = b
					c.channel.PublishWithContext(ctx, "gig_from_user_RPC", "gig_from_user_RPC", false, false, r)
					break
				}

				res.Success = true
				res.Data = &types.SellerOverview{
					FullName:         s.FullName,
					RatingsCount:     s.RatingsCount,
					RatingSum:        s.RatingSum,
					RatingCategories: s.RatingCategories,
				}
				b, _ := json.Marshal(&res)
				r.Body = b
				c.channel.PublishWithContext(ctx, "gig_from_user_RPC", "gig_from_user_RPC", false, false, r)

				log.Println("Message published")
				break
			default:
				log.Printf("Unknown message route [%s]", msg.RoutingKey)
				break
			}

			if err := msg.Ack(false); err != nil {
				fmt.Println("Error acknowledging message", msg.RoutingKey, topic, string(msg.Body))
			}

		}
	}()

	log.Printf(" [*] Waiting for logs from Gig Service. To exit press CTRL+C")
	<-forever
}
