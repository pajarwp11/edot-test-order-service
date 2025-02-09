package rabbitmq

import (
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

var (
	RabbitConn  *amqp091.Connection
	exhangeName = "stock_events"
)

func Connect() {
	var err error
	RabbitConn, err = amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
}
