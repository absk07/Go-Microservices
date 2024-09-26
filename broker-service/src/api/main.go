package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = ":8080"

type Config struct {
	RabbitConn *amqp.Connection
}

func main() {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
	}
	defer conn.Close()

	app := Config{
		RabbitConn: conn,
	}

	server := app.Routes()

	err = server.Run(PORT)
	if err != nil {
		log.Panic(err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
		return nil, err
	}

	return conn, nil
}