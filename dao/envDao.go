package dao

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
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
			logrus.Infof("No document found for workload ID: %s in collection: %s", workloadID, collectionName)
			return nil, nil // nil result is valid, just not found
		}
		logrus.Errorf("Error fetching workload ID %s from collection %s: %v", workloadID, collectionName, err)
		return nil, err
	}

	logrus.Infof("Fetched typed workload config for '%s': %+v", workloadID, result)
	return &result, nil
}
