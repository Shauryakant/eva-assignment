package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func ConnectDB(uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping mongodb: %w", err)
	}

	db := client.Database(dbName)

	usersColl := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err = usersColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		log.Printf("warning: could not create unique index on users.email: %v", err)
	}

	return &DB{
		Client:   client,
		Database: db,
	}, nil
}

func (db *DB) UsersCollection() *mongo.Collection {
	return db.Database.Collection("users")
}

func (db *DB) TicketsCollection() *mongo.Collection {
	return db.Database.Collection("tickets")
}
