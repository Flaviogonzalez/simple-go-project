package event

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type handler map[string]func(event json.RawMessage) ([]byte, error)

type Consumer struct {
	conn     *amqp.Connection
	exchange string
	handlers handler
}

func NewConsumer(conn *amqp.Connection, Exchange string, handlers handler) *Consumer {
	return &Consumer{
		conn:     conn,
		exchange: Exchange,
		handlers: handlers,
	}
}

func (c *Consumer) Setup() error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	return ch.ExchangeDeclare(
		c.exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
}

func (c *Consumer) Listen(topics []string) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = ch.QueueBind(
			q.Name,
			topic,
			c.exchange,
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	forever := make(chan bool)

	log.Println("escuchando peticiones!")
	for d := range messages {
		go func(d amqp.Delivery) {
			log.Println("nuevo mensaje de: ", d.Exchange)
			var eventPayload EventPayload
			err = json.Unmarshal(d.Body, &eventPayload)
			if err != nil {
				_ = d.Nack(false, false)
				return
			}

			err = c.handlePayload(eventPayload, ch, d)
			if err != nil {
				log.Printf("error handling payload: %v", err)
				_ = d.Nack(false, false)
				return
			}
			d.Ack(true)
		}(d)
	}
	<-forever

	return nil
}
