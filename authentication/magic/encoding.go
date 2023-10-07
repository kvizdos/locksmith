package magic

func MagicFromMap(input map[string]interface{}) MagicAuthentication {
	return MagicAuthentication{
		Code:               input["code"].(string),
		AllowedPermissions: input["permissions"].([]string),
		ExpiresAt:          input["expires"].(int64),
	}
}

func MagicsFromMap(input []map[string]interface{}) MagicAuthentications {
	magics := make([]MagicAuthentication, len(input))

	for i, rm := range input {
		magics[i] = MagicFromMap(rm)
	}

	return magics
}
