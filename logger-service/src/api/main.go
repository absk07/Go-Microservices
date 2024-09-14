package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	PORT     = ":3000"
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
}

func main() {
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(mongoURL))
	if err != nil {
		log.Panic("error connecting to database")
	}

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
