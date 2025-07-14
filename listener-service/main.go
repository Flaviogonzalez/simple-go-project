package main

import (
	"encoding/json"
	"listener-service/event"
	"listener-service/handlers"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var authHandlers = map[string]func(event json.RawMessage) ([]byte, error){
	"RegisterEvent": handlers.HandleRegister,
}

func main() {
	conn := connect()

	consumer := event.NewConsumer(conn, "AuthenticationService", authHandlers)
	err := consumer.Setup()
	if err != nil {
		log.Println(err.Error())
	}
	log.Println("bien. creó el consumidor!")

	err = consumer.Listen([]string{"user.register", "user.login"})
	if err != nil {
		log.Println(err.Error())
	}
}

func connect() *amqp.Connection {
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
