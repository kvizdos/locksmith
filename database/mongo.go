package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Direction type to specify the index direction
type Direction int

const (
	Ascending  Direction = 1  // Ascending order
	Descending Direction = -1 // Descending order
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

	fmt.Println("Connected to", uri, database)

	data := client.Database(database)

	db.database = data
	return nil
}
func (db MongoDatabase) MonitorConnection(heartbeat time.Duration, health HealthCheckInterface) {
	ticker := time.NewTicker(heartbeat) // Adjust the interval as needed.
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, timeout := context.WithTimeout(context.Background(), heartbeat)
			defer timeout()

			err := db.database.Client().Ping(ctx, readpref.Primary())
			if err != nil {
				fmt.Printf("MongoDB failed to ping health: %v\n", err)
				health.SetMongoDown()
			} else if !health.IsMongoUp() {
				health.SetMongoUp()
			}
		}
	}
}

func (db MongoDatabase) GetUTCTimestampFromID(dbID primitive.ObjectID) (time.Time, error) {
	timestamp := dbID.Timestamp()

	return timestamp, nil
}

func (db MongoDatabase) Drop(table string) error {
	col := db.database.Collection(table)

	err := col.Drop(context.Background())

	if err != nil {
		return err
	}

	return nil
}

func (db MongoDatabase) Aggregate(table string, pipeline []map[string]interface{}) ([]map[string]interface{}, error) {
	col := db.database.Collection(table)

	res, err := col.Aggregate(context.TODO(), pipeline)

	if err != nil {
		return []map[string]interface{}{}, err
	}

	var results []map[string]interface{}
	res.All(context.Background(), &results)

	finalResults := make([]map[string]interface{}, len(results))

	for i, result := range results {
		convertArraysToSlices(result)
		finalResults[i] = result
	}

	return finalResults, nil
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
		fmt.Println(err)
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

func (db MongoDatabase) CreateTextIndex(table string, keys []string) error {
	// Convert []string to bson.D
	var keysPrim bson.D
	for _, fieldName := range keys {
		keysPrim = append(keysPrim, bson.E{Key: fieldName, Value: "text"})
	}

	model := mongo.IndexModel{Keys: keysPrim}
	_, err := db.database.Collection(table).Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}

	return err
}

func (db MongoDatabase) CreateRegularIndex(table string, keys map[string]Direction, unique bool) error {
	// Convert map[string]Direction to bson.D
	var keysPrim bson.D
	for fieldName, direction := range keys {
		keysPrim = append(keysPrim, bson.E{Key: fieldName, Value: direction})
	}

	model := mongo.IndexModel{
		Keys:    keysPrim,
		Options: options.Index().SetUnique(unique),
	}
	_, err := db.database.Collection(table).Indexes().CreateOne(context.TODO(), model)
	if err != nil {
		fmt.Println("Error creating index:", err)
		return err
	}

	return nil
}

func (db MongoDatabase) Count(table string, filter map[string]interface{}) (int64, error) {
	return db.database.Collection(table).CountDocuments(context.TODO(), filter)
}

func (db MongoDatabase) FindPaginated(table string, query map[string]interface{}, maxPages int64, lastID string) ([]map[string]interface{}, bool) {
	// Set the maximum and sort by ID
	opts := options.Find().SetLimit(maxPages).SetSort(bson.D{{"_id", 1}})

	if len(lastID) > 0 {
		objID, err := primitive.ObjectIDFromHex(lastID)

		if err != nil {
			return []map[string]interface{}{}, false
		}

		query["_id"] = map[string]interface{}{"$gt": objID}
	}

	res, err := db.database.Collection(table).Find(context.TODO(), query, opts)

	if err != nil {
		return []map[string]interface{}{}, false
	}

	var results []map[string]interface{}
	res.All(context.Background(), &results)

	return results, true
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
