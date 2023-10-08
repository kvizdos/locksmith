package magic

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MagicFromMap(input map[string]interface{}) MagicAuthentication {
	var permissions []string
	switch input["permissions"].(type) {
	case primitive.A:
		for _, item := range input["permissions"].(primitive.A) {
			permissions = append(permissions, item.(string))
		}
	case []string:
		permissions = input["permissions"].([]string)
	case []interface{}:
		for _, item := range input["permissions"].([]interface{}) {
			permissions = append(permissions, item.(string))
		}
	}
	return MagicAuthentication{
		Code:               input["code"].(string),
		AllowedPermissions: permissions,
		ExpiresAt:          input["expires"].(int64),
	}
}

func MagicsFromMap(input []interface{}) MagicAuthentications {
	magics := make([]MagicAuthentication, len(input))

	for i, rm := range input {
		magics[i] = MagicFromMap(rm.(map[string]interface{}))
	}

	return magics
}
