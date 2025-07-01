package db

import (
	"context"
	"time"

	log "log/slog"

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

	log.Info("Connecting to MongoDB", "uri", dbConnectionString)

	clientOptions := options.Client().ApplyURI(dbConnectionString)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Error("DB connection failed", "error", err)
		panic(err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Error("Failed to ping MongoDB", "error", err)
		panic(err)
	}

	cfg := config.GetConfig()
	MongoClient = client
	MongoDatabase = client.Database(cfg.Database.DBName)

	log.Info("MongoDB connection established", "database", cfg.Database.DBName)
}
