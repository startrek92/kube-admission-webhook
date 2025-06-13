package db

import (
	"context"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


var MongoClient *mongo.Client;

func Connect(dbConnectionString string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second);
	defer cancel()

	clientOptions := options.Client().ApplyURI(dbConnectionString);

	client, err := mongo.Connect(ctx, clientOptions);

	if err != nil {
		log.Fatal("db connection failed", err);
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("failed to ping db", err);
	}

	log.Printf("db connection success")
	MongoClient = client;
	
}
