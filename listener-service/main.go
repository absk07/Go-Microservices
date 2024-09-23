package main

import (
	"log"

	"github.com/listener-service/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	conn, err := connect()
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
	}
	defer conn.Close()

	// start listening for mesages
	log.Println("Listenibg and consuming for RabbitMQ messages")

	// create consumer
	consumer, err := events.NewConsumer(conn)
	if err != nil {
		log.Panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
		return nil, err
	}

	return conn, nil
}
