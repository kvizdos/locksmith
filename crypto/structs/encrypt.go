package structs

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kvizdos/locksmith/crypto"
)

// TODO: Support nested structures

// First bool says whether or not it matches,
// Second bool says whether or not it is a map
func isSupportedType(tp reflect.Type) bool {
	if strings.Contains(tp.String(), ".") {
		return true
	}

	supported := map[string]bool{
		"string":            true,
		"[]string":          true,
		"map[string]string": true,
	}

	_, yes := supported[tp.String()]

	return yes
}

// This will take in a struct with `encrypt` tags assigned to it,
// and return the encrypted map version
// Encryption only supports string and []string
func EncryptStructIntoMap(s interface{}, useSerializer Serializer, engine crypto.CryptoEngineInterface) (map[string]interface{}, error) {
	// ValueOf returns a Value representing the run-time data
	v := reflect.ValueOf(s)
	finalizedMap := make(map[string]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		tag := field.Tag.Get("encrypt")
		serializer := field.Tag.Get(string(useSerializer))
		fieldValue := v.Field(i)

		if serializer == "" || serializer == "-" {
			continue
		}

		if tag == "" && !strings.Contains(field.Type.String(), ".") {
			finalizedMap[serializer] = fieldValue.Interface()
			continue
		}

		if !isSupportedType(field.Type) {
			fmt.Printf("%s has unsupported type for encryption\n", field.Name)
			finalizedMap[serializer] = fieldValue.Interface()
			continue
		}

		key, err := engine.GetKeyForOperation("structs")

		if err != nil {
			return map[string]interface{}{}, err
		}

		switch fieldValue.Type().String() {
		case "[]string":
			slice := fieldValue.Interface().([]string)

			newValue := make([]string, len(slice))
			for i, value := range slice {
				encrypted, err := engine.Encrypt(key, value)
				if err != nil {
					return map[string]interface{}{}, fmt.Errorf("failed to encrypt array value")
				}
				newValue[i] = encrypted
			}

			finalizedMap[serializer] = newValue
		case "string":
			encrypted, err := engine.Encrypt(key, fieldValue.String())
			if err != nil {
				return map[string]interface{}{}, fmt.Errorf("failed to encrypt array value")
			}
			finalizedMap[serializer] = encrypted
		case "map[string]string":
			// handle maps
			mapValue := fieldValue.Interface().(map[string]string)

			for k, v := range mapValue {
				encrypted, err := engine.Encrypt(key, v)
				if err != nil {
					return map[string]interface{}{}, fmt.Errorf("failed to encrypt array value")
				}
				mapValue[k] = encrypted
			}

			finalizedMap[serializer] = mapValue
		}
	}
	return finalizedMap, nil
}
