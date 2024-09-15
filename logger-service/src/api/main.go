package main

import (
	"context"
	"log"
	"time"

	"github.com/logger-service/data"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	PORT     = ":7070"
	mongoURL = "mongodb://mongo:27017"
)

type Config struct {
	DB     *mongo.Client
	Models data.Models
}

func main() {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Panic("error connecting to database")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		DB:     client,
		Models: *data.New(client),
	}

	// go app.startGinServer()
	server := app.Routes()

	err = server.Run(PORT)
	if err != nil {
		log.Panic("error starting gin server")
	}
}

// func (app *Config) startGinServer() {
// 	server := app.Routes()

// 	err := server.Run(PORT)
// 	if err != nil {
// 		log.Panic("error starting gin server")
// 	}
// }
