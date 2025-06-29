package db

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/startrek92/kube-admission-webhook/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient   *mongo.Client
	MongoDatabase *mongo.Database
)

func Connect(dbConnectionString string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logrus.Infof("Connecting to MongoDB at %s", dbConnectionString)

	clientOptions := options.Client().ApplyURI(dbConnectionString)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Fatalf("DB connection failed: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		logrus.Fatalf("Failed to ping MongoDB: %v", err)
	}

	cfg := config.GetConfig()
	MongoClient = client
	MongoDatabase = client.Database(cfg.Database.DBName)

	logrus.Infof("MongoDB connection established to DB: %s", cfg.Database.DBName)
}
