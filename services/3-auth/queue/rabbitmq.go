package queue

import (
	"context"
	"fmt"
	"log"
	"os"

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
	conn    *amqp.Connection
	channel *amqp.Channel
	err     chan error
	db      *gorm.DB
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
		conn:    conn,
		channel: channel,
		err:     make(chan error),
		db:      db,
	}
}

func (c *Connection) PublishDirect(ctx context.Context, errCh chan<- error, exchangeName, routingKey, topic string, msg []byte, corrId ...string) error {
	err := c.channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		err = fmt.Errorf("error declaring exchange :\n%+v", err)
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
	log.Println("message published")
	return nil
}
