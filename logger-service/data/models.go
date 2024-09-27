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
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Method       string             `bson:"method" json:"method"`
	Path         string             `bson:"path" json:"path"`
	RemoteAddr   string             `bson:"remote_addr" json:"remote_addr"`
	ResponseTime string             `bson:"response_time" json:"response_time"`
	StartTime    string             `bson:"start_time" json:"start_time"`
	StatusCode   string             `bson:"status_code" json:"status_code"`
	Severity     string             `bson:"severity" json:"severity"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.Background(), LogEntry{
		Method:       entry.Method,
		Path:         entry.Path,
		RemoteAddr:   entry.RemoteAddr,
		ResponseTime: entry.ResponseTime,
		StartTime:    entry.StartTime,
		StatusCode:   entry.StatusCode,
		Severity:     entry.Severity,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
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
	opts.SetSort(bson.D{{"created_at", -1}})

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
	err = collection.FindOne(ctx, bson.M{"_id": _id}).Decode(&entry)
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

// func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	collection := client.Database("logs").Collection("logs")

// 	_id, err := primitive.ObjectIDFromHex(l.ID)
// 	if err != nil {
// 		log.Println("Error converting _id", err)
// 		return nil, err
// 	}

// 	result, err := collection.UpdateOne(
// 		ctx,
// 		bson.M{ "_id": _id },
// 		bson.D{
// 			{
// 				"$set", bson.D{
// 					{ "name", l.Name },
// 					{ "data", l.Data },
// 					{ "updated_at", time.Now() },
// 				},
// 			},
// 		},
// 	)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return result, nil
// }
