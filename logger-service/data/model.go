package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbTimeOut = 15 * time.Second

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Data      string    `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// POST
func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("log").Collection("log")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	if err != nil {
		log.Println("Error inserting into log", err)
		return err
	}

	return nil
}

// GET
func (l *LogEntry) GetAll() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()
	collection := client.Database("log").Collection("log")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var logs []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)

		if err != nil {
			return nil, err
		} else {
			logs = append(logs, &item)
		}
	}
	return logs, nil
}

func (l *LogEntry) getOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)

	defer cancel()
	collection := client.Database("log").Collection("log")
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var entry LogEntry

	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// UPDATE
func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)

	defer cancel()

	collection := client.Database("log").Collection("log")

	docId, err := primitive.ObjectIDFromHex(l.ID)

	if err != nil {
		return nil, err
	}

	res, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": docId},
		bson.D{
			{
				Key: "$set", Value: bson.D{
					{Key: "name", Value: l.Name},
					{Key: "data", Value: l.Data},
					{Key: "updated_at", Value: time.Now()},
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return res, nil

}

// DELETE
func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeOut)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		return err
	}
	return nil
}
