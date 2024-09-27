package main

import (
	"context"
	"log"
	"os"

	"github.com/auth-service/data"
	"github.com/jackc/pgx/v5/pgxpool"
	amqp "github.com/rabbitmq/amqp091-go"
)

const PORT = ":9090"

type Config struct{
	DB *pgxpool.Pool
	Models data.Models
	RabbitConn *amqp.Connection
}

func main() {
	connPool, err := pgxpool.New(context.Background(), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Panic("error connecting to database")
	}
	defer connPool.Close()
	
	var msg string
	err = connPool.QueryRow(context.Background(), "SELECT 'Database successfully connected'").Scan(&msg)
	if err != nil {
		log.Panic(err)
	}

	conn, err := connectRabbitMQ()
	if err != nil {
		log.Panicf("Error connecting to rabbitmq: %s", err)
	}
	defer conn.Close()

	app := Config{
		DB: connPool,
		Models: *data.New(connPool),
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