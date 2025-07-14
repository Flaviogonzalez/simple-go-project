package event

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPayload struct {
	Name  string
	Event json.RawMessage
}

func (c *Consumer) handlePayload(payload EventPayload, ch *amqp.Channel, msg amqp.Delivery) error {
	var response []byte
	function, ok := c.handlers[payload.Name]
	if !ok {
		log.Println("error trying to execute a function method for the event")
		return nil
	}

	log.Println("Event detected, executing function")
	r, err := function(payload.Event)
	if err != nil {
		log.Println("error executing event: ", err)
	}
	response = r

	if msg.ReplyTo == "" {
		log.Printf("No ReplyTo queue specified for event: %s, correlationId: %s", payload.Name, msg.CorrelationId)
		return nil
	}

	err = ch.Publish(
		"",
		msg.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: msg.CorrelationId,
			Body:          response,
		},
	)

	if err != nil {
		log.Printf("Failed to publish response to %s: %v", msg.ReplyTo, err)
		return err
	}

	log.Printf("Successfully published response to %s for event: %s, correlationId: %s", msg.ReplyTo, payload.Name, msg.CorrelationId)
	return nil
}
