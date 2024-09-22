package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := connect()
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
	}
	defer conn.Close()
}

func connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
		return nil, err
	}

	return conn, nil
}
