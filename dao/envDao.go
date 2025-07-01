package dao

import (
	"context"
	"time"

	log "log/slog"

	"github.com/startrek92/kube-admission-webhook/db"
	mongomodels "github.com/startrek92/kube-admission-webhook/mongoModels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWorkloadEnv(collectionName string, workloadID string) (*mongomodels.WorkloadConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoDatabase.Collection(collectionName)

	var result mongomodels.WorkloadConfig
	err := collection.FindOne(ctx, bson.M{"_id": workloadID}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Info("No document found for workload", "workload_id", workloadID, "collection", collectionName)
			return nil, nil
		}
		log.Error("Error fetching workload from MongoDB", "workload_id", workloadID, "collection", collectionName, "error", err)
		return nil, err
	}

	log.Info("Fetched workload config", "workload_id", workloadID, "config", result)
	return &result, nil
}
