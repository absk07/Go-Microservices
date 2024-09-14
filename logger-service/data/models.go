package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) *Models {
	client = mongo
	return &Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID string `bson:"_id,omitempty" json:"id,omitempty"`
	Name string `bson:"name" json:"name"`
	Data string `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.Background(), LogEntry{
		Name: entry.Name,
		Data: entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error Inserting logs", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{ "created_at", -1 }})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error finding logs", err)
		return nil, err
	}

	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry
		if err := cursor.Decode(&item); err != nil {
			log.Println("Error decoding logs", err)
			return nil, err
		}
		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting _id", err)
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{ "_id": _id }).Decode(&entry)
	if err != nil {
		log.Println("Error finding log by _id", err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	_id, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting _id", err)
		return nil, err
	}

	result, err := collection.UpdateOne(
		ctx,
		bson.M{ "_id": _id },
		bson.D{
			{ 
				"$set", bson.D{
					{ "name", l.Name },
					{ "data", l.Data },
					{ "updated_at", time.Now() },
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return result, nil
}