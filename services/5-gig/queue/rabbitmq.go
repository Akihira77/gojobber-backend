package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Akihira77/gojobber/services/5-gig/types"
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
	conn     *amqp.Connection
	channel  *amqp.Channel
	err      chan error
	db       *gorm.DB
	validate *validator.Validate
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
		conn:     conn,
		channel:  channel,
		err:      make(chan error),
		db:       db,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}
}

func (c *Connection) NewChannel() (*amqp.Channel, error) {
	return c.conn.Channel()
}

func (c *Connection) PublishDirect(ctx context.Context, errCh chan<- error, exchangeName, routingKey, topic string, msg []byte, corrId ...string) error {
	err := c.channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		err = fmt.Errorf("publishing error declaring exchange :\n%+v", err)
		errCh <- err
		return err
	}

	err = c.channel.PublishWithContext(ctx, exchangeName, routingKey, false, false, amqp.Publishing{
		Headers: amqp.Table{
			"Topic": topic,
		},
		ContentType:   "application/json",
		Body:          msg,
		CorrelationId: corrId[0],
		ReplyTo:       routingKey,
	})
	if err != nil {
		err = fmt.Errorf("error publishing message :\n%+v", err)
		errCh <- err
		return err
	}

	errCh <- nil
	log.Println("message published correlationId:", corrId[0])
	return nil
}

func (c *Connection) ConsumeFromReviewService() {
	routing, exchange := "gig_from_review", "gig_from_review"
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
			ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
			defer cancel()

			topic := msg.Headers["Topic"]
			log.Printf("receive message from [%s] route, [%s] topic, body: %v\n", msg.RoutingKey, topic, string(msg.Body))

			//TODO: SET THE REVIEW PUBLISH AND CONSUME (RPC)
			switch topic {
			case "upsert-gig-review":
				gigId := string(msg.Body)
				res := types.RabbitMQResponse[string]{
					Success: true,
					Data:    "",
				}

				err := c.db.
					Debug().
					WithContext(ctx).
					Model(&types.Gig{}).
					Where("id = ?", gigId).
					First(&types.Gig{}).
					Error
				if err != nil {
					res.Success = false
					res.Data = err.Error()
				}

				b, _ := json.Marshal(res)
				err = c.channel.PublishWithContext(ctx, "gig_to_review", "gig_to_review", false, false, amqp.Publishing{
					CorrelationId: msg.CorrelationId,
					ReplyTo:       msg.ReplyTo,
					ContentType:   "application/json",
					Body:          b,
				})
				break
			default:
				fmt.Println("unknown message")
				break
			}

			if err := msg.Ack(false); err != nil {
				fmt.Println("Error acknowledging message", msg.RoutingKey, topic, string(msg.Body))
			}
		}
	}()

	fmt.Println(" [*] Waiting for logs from Review Service. To exit press CTRL+C")
	<-forever
}

func (c *Connection) ConsumeFromUserRPC(ctx context.Context, ch *amqp.Channel) <-chan amqp.Delivery {
	routing, exchange := "gig_from_user_RPC", "gig_from_user_RPC"
	err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		fmt.Printf("consuming error declaring exchange :\n%+v", err)
		return nil
	}

	q, err := ch.QueueDeclare(routing, true, false, false, false, nil)
	if err != nil {
		fmt.Printf("error declaring queue :\n%+v", err)
		return nil
	}

	err = ch.QueueBind(q.Name, routing, exchange, false, nil)
	if err != nil {
		fmt.Printf("error binding queue :\n%+v", err)
		return nil
	}

	msgs, err := ch.ConsumeWithContext(ctx, q.Name, routing, false, false, false, false, nil)
	if err != nil {
		fmt.Printf("error consuming message :\n%+v", err)
		return nil
	}

	return msgs
}
