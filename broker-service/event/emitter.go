package event

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	conn     *amqp.Connection
	exchange string
	id       uuid.UUID
}

func NewEmitter(conn *amqp.Connection, exchange string) *Emitter {
	uuid := uuid.New()
	return &Emitter{
		conn:     conn,
		exchange: exchange,
		id:       uuid,
	}
}

func (e *Emitter) Push(w http.ResponseWriter, topicPayload TopicPayload) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	jsonBytes, err := json.Marshal(topicPayload.Event)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare("", true, true, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.Publish(
		e.exchange,
		topicPayload.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: e.id.String(), // custom id
			ReplyTo:       q.Name,        // Set to a "queue name" if you expect a reply
			Body:          jsonBytes,
		},
	)
	if err != nil {
		return err
	}

	e.SendResponse(w, q)

	return nil
}

func (e *Emitter) SendResponse(w http.ResponseWriter, q amqp.Queue) error {
	ch, err := e.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer tag (auto-generated)
		true,   // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		http.Error(w, "Failed to set up consumer", http.StatusInternalServerError)
		return fmt.Errorf("failed to set up consumer: %w", err)
	}

	msg := <-msgs
	if msg.CorrelationId == e.id.String() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(msg.Body); err != nil {
			return fmt.Errorf("failed to write HTTP response: %w", err)
		}
		return nil
	}

	return nil
}
