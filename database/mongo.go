package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDatabase struct {
	Ctx      context.Context
	Cancel   context.CancelFunc
	database *mongo.Database
}

func (db *MongoDatabase) Initialize(uri string, database string) error {
	clientOptions := options.Client().ApplyURI(uri)
	defer db.Cancel()

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping mongodb")
	}

	data := client.Database(database)

	fmt.Println("connected")

	db.database = data
	return nil
}

func (db MongoDatabase) DeleteOne(table string, query map[string]interface{}) (bool, error) {
	col := db.database.Collection(table)

	res, err := col.DeleteOne(context.Background(), query)

	if err != nil {
		return false, err
	}

	if res.DeletedCount == 0 {
		return false, nil
	}

	return true, nil
}

func (db MongoDatabase) FindOne(table string, query map[string]interface{}) (interface{}, bool) {
	col := db.database.Collection(table)

	res := col.FindOne(context.Background(), query)

	if res.Err() != nil {
		return map[string]interface{}{}, false
	}

	var result map[string]interface{}
	res.Decode(&result)

	convertArraysToSlices(result)

	return result, true
}

func (db MongoDatabase) Find(table string, query map[string]interface{}) ([]interface{}, bool) {
	col := db.database.Collection(table)

	res, err := col.Find(context.Background(), query)

	if err != nil {
		return []interface{}{}, false
	}

	var results []map[string]interface{}
	res.All(context.Background(), &results)

	finalResults := make([]interface{}, len(results))

	for i, result := range results {
		convertArraysToSlices(result)
		finalResults[i] = result
	}

	return finalResults, true
}

func convertArraysToSlices(data interface{}) {
	switch val := data.(type) {
	case map[string]interface{}:
		for key, value := range val {
			convertArraysToSlices(value)
			if isArray(value) {
				val[key] = convertToSlice(value)
			}
		}
	}
}

func isArray(value interface{}) bool {
	switch value.(type) {
	case primitive.A, []interface{}:
		return true
	default:
		return false
	}
}

func convertToSlice(value interface{}) []interface{} {
	switch v := value.(type) {
	case primitive.A:
		return convertPrimitiveArrayToSlice(v)
	case []interface{}:
		return v
	default:
		return nil
	}
}

func convertPrimitiveArrayToSlice(array primitive.A) []interface{} {
	slice := make([]interface{}, len(array))
	for i, item := range array {
		slice[i] = item
	}
	return slice
}

func (db MongoDatabase) InsertOne(table string, body map[string]interface{}) (interface{}, error) {
	col := db.database.Collection(table)

	res, err := col.InsertOne(context.Background(), body)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (db MongoDatabase) UpdateOne(table string, query map[string]interface{}, body map[DatabaseUpdateActions]map[string]interface{}) (interface{}, error) {
	col := db.database.Collection(table)

	bsonBody := bson.M{}
	var useAction DatabaseUpdateActions
	for action, fields := range body {
		useAction = action
		bsonFields := bson.M{}
		for key, value := range fields {
			bsonFields[key] = value
		}
		bsonBody = bsonFields
	}

	var res *mongo.UpdateResult
	var err error

	switch useAction {
	case PUSH:
		res, err = col.UpdateOne(context.Background(), query, bson.M{"$push": bsonBody})
		break
	case SET:
		res, err = col.UpdateOne(context.Background(), query, bson.M{"$set": bsonBody})
		break
	}
	if err != nil {
		return nil, err
	}

	return res, nil
}
