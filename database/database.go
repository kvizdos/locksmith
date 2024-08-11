package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DatabaseUpdateActions string

const (
	SET  DatabaseUpdateActions = "set"
	PUSH DatabaseUpdateActions = "push"
)

type HealthCheckInterface interface {
	SetMongoDown()
	SetMongoUp()
	IsMongoUp() bool
}

type TransactionOptions struct {
	SessionOptions     *options.SessionOptions
	TransactionOptions *options.TransactionOptions
}

type DatabaseAccessor interface {
	InsertOne(table string, body map[string]interface{}) (interface{}, error)
	UpdateOne(table string, query map[string]interface{}, body map[DatabaseUpdateActions]map[string]interface{}) (interface{}, error)
	FindOne(table string, query map[string]interface{}) (interface{}, bool)
	Find(table string, query map[string]interface{}) ([]interface{}, bool)
	FindPaginated(table string, query map[string]interface{}, maxPages int64, lastID string) ([]map[string]interface{}, bool)
	DeleteOne(table string, query map[string]interface{}) (bool, error)
	CreateTextIndex(table string, keys []string) error
	CreateRegularIndex(table string, keys map[string]Direction, unique bool) error
	Drop(table string) error
	Aggregate(table string, pipeline []map[string]interface{}) ([]map[string]interface{}, error)
	GetUTCTimestampFromID(dbID primitive.ObjectID) (time.Time, error)
	MonitorConnection(heartbeat time.Duration, health HealthCheckInterface)
	Transact(ctx context.Context, opts *TransactionOptions, transaction func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error)
}

type TestDatabase struct {
	Tables map[string]map[string]interface{}
}

func (db TestDatabase) GetUTCTimestampFromID(dbID primitive.ObjectID) (time.Time, error) {
	return time.Now(), nil
}

func (db TestDatabase) MonitorConnection(heartbeat time.Duration, health HealthCheckInterface) {
	// do nothing, nothing to monitor
}

func (db TestDatabase) CreateTextIndex(table string, keys []string) error {
	return nil
}

func (db TestDatabase) CreateRegularIndex(table string, keys map[string]Direction, unique bool) error {
	return nil
}

func (db TestDatabase) Transact(ctx context.Context, opts *TransactionOptions, transaction func(sessCtx mongo.SessionContext) (interface{}, error)) (interface{}, error) {
	log.Println("WARNING: Transactions not supported in a testing environment.")
	return nil, nil
}

func (db TestDatabase) Aggregate(table string, pipeline []map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

func (db TestDatabase) Drop(table string) error {
	db.Tables[table] = map[string]interface{}{}
	return nil
}

func (db TestDatabase) FindPaginated(table string, query map[string]interface{}, maxPages int64, lastID string) ([]map[string]interface{}, bool) {
	// TODO: Probably stub this
	return []map[string]interface{}{}, false
}

func (db TestDatabase) InsertOne(table string, body map[string]interface{}) (interface{}, error) {
	if _, ok := db.Tables[table]; !ok {
		db.Tables[table] = make(map[string]interface{})
	}

	// Generate a unique ID for the new row
	id := uuid.New().String()

	// Assign the generated ID to the body
	body["id"] = id

	// Insert the new row into the table
	db.Tables[table][id] = body

	return id, nil
}

func (db TestDatabase) UpdateOne(table string, query map[string]interface{}, body map[DatabaseUpdateActions]map[string]interface{}) (interface{}, error) {
	if tableData, ok := db.Tables[table]; ok {
		for _, row := range tableData {
			match := true
			for queryKey, queryValue := range query {
				if rowData, ok := row.(map[string]interface{}); ok {
					if fieldValue, ok := rowData[queryKey]; ok && fieldValue != queryValue {
						match = false
						break
					}
				} else {
					match = false
					break
				}
			}

			if match {
				for action, updateBody := range body {
					if action == SET {
						for key, value := range updateBody {
							row.(map[string]interface{})[key] = value
						}
					} else if action == PUSH {
						for key, value := range updateBody {
							if arrayField, ok := row.(map[string]interface{})[key].([]interface{}); ok {
								// Check if the arrayField is an empty interface
								if arrayField == nil {
									// Initialize as an empty slice of strings
									arrayField = []interface{}{}
								}

								// Append the value to the arrayField
								arrayField = append(arrayField, value)

								// Update the arrayField in the row
								row.(map[string]interface{})[key] = arrayField
							} else {
								match = false
								break
							}
						}
					} else {
						return nil, fmt.Errorf("unsupported update action")
					}
				}

				if match {
					return row, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("row not found")
}

func (db TestDatabase) FindOne(table string, query map[string]interface{}) (interface{}, bool) {
	if tableData, ok := db.Tables[table]; ok {
		for _, row := range tableData {
			match := true
			for queryKey, queryValue := range query {
				if rowData, ok := row.(map[string]interface{}); ok {
					if fieldValue, ok := rowData[queryKey]; ok && fieldValue != queryValue {
						match = false
						break
					}
				} else {
					match = false
					break
				}
			}

			if match {
				return row, true
			}
		}
	}

	return nil, false
}

func (db TestDatabase) Find(table string, query map[string]interface{}) ([]interface{}, bool) {
	if tableData, ok := db.Tables[table]; ok {
		var results []interface{}

		// Special case for $or query
		if orQueries, ok := query["$or"]; ok {
			if orQueriesList, ok := orQueries.([]map[string]interface{}); ok {
				for _, orQuery := range orQueriesList {
					if result, ok := db.Find(table, orQuery); ok {
						results = append(results, result...)
					}
				}

				if len(results) > 0 {
					return results, true
				}
			}
		} else {
			// Regular query
			for _, row := range tableData {
				match := true
				for queryKey, queryValue := range query {
					if rowData, ok := row.(map[string]interface{}); ok {
						if fieldValue, ok := rowData[queryKey]; ok && fieldValue != queryValue {
							match = false
							break
						}
					} else {
						match = false
						break
					}
				}

				if match {
					results = append(results, row)
				}
			}

			if len(results) > 0 {
				return results, true
			}
		}
	}

	return nil, false
}

func (db TestDatabase) DeleteOne(table string, query map[string]interface{}) (bool, error) {
	if tableData, ok := db.Tables[table]; ok {
		for key, row := range tableData {
			match := true
			for queryKey, queryValue := range query {
				if rowData, ok := row.(map[string]interface{}); ok {
					if fieldValue, ok := rowData[queryKey]; ok && fieldValue != queryValue {
						match = false
						break
					}
				} else {
					match = false
					break
				}
			}

			if match {
				// Create a new map without the matching item
				newTableData := make(map[string]interface{})
				for k, v := range tableData {
					if k != key {
						newTableData[k] = v
					}
				}

				// Update the map in the database with the modified map
				db.Tables[table] = newTableData
				return true, nil
			}
		}
	}

	return false, nil
}
