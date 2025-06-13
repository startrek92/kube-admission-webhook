package dao

import (
	"context"

	"github.com/startrek92/kube-admission-webhook/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	dataBase   = "env_db"
	collection = "env_global"
)

func getEnvVars(docId string) (any, error) {
	collection := db.MongoClient.Database(dataBase).Collection(collection)
	var result any
	err := collection.FindOne(context.TODO(), bson.M{"_id": docId}).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	return &result, err
}
