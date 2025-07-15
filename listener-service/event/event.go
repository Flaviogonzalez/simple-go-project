package event

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPayload struct { // los tipos de datos que se envian y reciben tienen que tener LOS MISMOS NOMBRES TANTO EN JSON COMO EN EL STRUCT
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

func (c *Consumer) handlePayload(payload EventPayload, ch *amqp.Channel, msg amqp.Delivery) error {
	var response []byte
	function, ok := c.handlers[payload.Name]
	if !ok {
		log.Println("error trying to execute a function method for the event")
		return nil
	}

	log.Println("Event detected, executing function")
	r, err := function(payload.Data)
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
