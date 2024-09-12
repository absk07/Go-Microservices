package main

import (
	"context"
	"log"
	"os"

	"github.com/auth-service/data"
	"github.com/jackc/pgx/v5/pgxpool"
)

const PORT = ":9090"

type Config struct{
	DB *pgxpool.Pool
	Models data.Models
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

	app := Config{
		DB: connPool,
		Models: *data.New(connPool),
	}

	server := app.Routes()

	err = server.Run(PORT)
	if err != nil {
		log.Panic(err)
	}
}