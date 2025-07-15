package event

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EventPayload struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}

type TopicPayload struct {
	Name  string       `json:"name"`
	Event EventPayload `json:"event"`
}

func ConnectToRabbit() *amqp.Connection {
	var conn *amqp.Connection
	var counts int64
	var backoff time.Duration

	for {
		connection, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Println("Error trying to connect. Trying again...")
			counts++
		} else {

			conn = connection
			break
		}

		if counts > 10 {
			log.Println("Error trying to connect.")
		}

		backoff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backoff)
		continue
	}

	return conn
}

func SendToListener(w http.ResponseWriter, exchange string, topicPayload TopicPayload) error {
	conn := ConnectToRabbit()
	e := NewEmitter(conn, exchange)

	err := e.Push(w, topicPayload)
	if err != nil {
		return err
	}

	return nil
}
