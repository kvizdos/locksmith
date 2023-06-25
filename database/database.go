package database

import (
	"fmt"

	"github.com/google/uuid"
)

type DatabaseUpdateActions string

const (
	SET  DatabaseUpdateActions = "set"
	PUSH DatabaseUpdateActions = "push"
)

type DatabaseAccessor interface {
	InsertOne(table string, body map[string]interface{}) (interface{}, error)
	UpdateOne(table string, query map[string]interface{}, body map[DatabaseUpdateActions]map[string]interface{}) (interface{}, error)
	FindOne(table string, query map[string]interface{}) (interface{}, bool)
}

type TestDatabase struct {
	Tables map[string]map[string]interface{}
}

func (db TestDatabase) GetValue(table string, query map[string]interface{}) (interface{}, error) {
	if tableData, ok := db.Tables[table]; ok {
		for _, value := range tableData {
			match := true
			for queryKey, queryValue := range query {
				if tableValue, ok := value.(map[string]interface{}); ok {
					if fieldValue, ok := tableValue[queryKey]; ok && fieldValue != queryValue {
						match = false
						break
					}
				} else {
					match = false
					break
				}
			}
			if match {
				return value, nil
			}
		}
	}

	return nil, fmt.Errorf("value not found")
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
