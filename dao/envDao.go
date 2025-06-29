package dao

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/startrek92/kube-admission-webhook/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetWorkloadEnv(collectionName string, workLoadId string) (map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := db.MongoDatabase.Collection(collectionName)

	var result map[string]interface{}
	err := collection.FindOne(ctx, bson.M{"_id": workLoadId}).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			logrus.Infof("No document found for workload ID: %s in collection: %s", workLoadId, collectionName)
			return nil, nil
		}
		logrus.Errorf("Error fetching workload ID %s from collection %s: %v", workLoadId, collectionName, err)
		return nil, err
	}

	logrus.Debugf("Fetched workload env: %+v", result)
	return result, nil
}