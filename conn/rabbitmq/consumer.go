package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type StockHandler interface {
	UpdateStatusByConsumer(data interface{}) error
}

type RabbitConsumer struct {
	rabbitConn   *amqp091.Connection
	stockHandler StockHandler
}

func NewRabbitConsumer(rabbitConn *amqp091.Connection, stockHandler StockHandler) *RabbitConsumer {
	return &RabbitConsumer{
		rabbitConn:   rabbitConn,
		stockHandler: stockHandler,
	}
}

func (r *RabbitConsumer) ConsumeEvents() {
	ch, err := r.rabbitConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exhangeName,
		"topic",
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	topics := map[string]string{
		"order.update_status": "queue_order_status",
	}

	for routingKey, queueName := range topics {
		go r.startConsumer(ch, queueName, routingKey)
	}

	log.Println("consumer started")
	select {}
}

func (r *RabbitConsumer) startConsumer(ch *amqp091.Channel, queueName, routingKey string) {
	q, err := ch.QueueDeclare(
		queueName,
		true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, routingKey, exhangeName, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {
			var event Event
			json.Unmarshal(d.Body, &event)
			log.Printf("Received from %s: %+v\n", queueName, event)
			err := r.handleEvent(event)
			if err != nil {
				log.Printf("Error processing event with data %s: %v\n", event, err)
				d.Nack(false, true)
			} else {
				log.Printf("Succes processing event with data: %v\n", event)
				d.Ack(false)
			}
		}
	}()
}

func (r *RabbitConsumer) handleEvent(event Event) error {
	switch event.Type {
	case "order.update_status":
		return r.stockHandler.UpdateStatusByConsumer(event.Data)
	default:
		fmt.Println("Unknown event:", event.Type)
		return nil
	}
}
