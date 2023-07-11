package structs

import (
	"testing"

	"github.com/kvizdos/locksmith/crypto"
)

type nestedEncryptedStruct struct {
	Value   string `json:"value" bson:"value" encrypt:"true"`
	NonEncr string `json:"nc" bson:"nc"`
}

type testStruct struct {
	DontInclude string
	Username    string                `json:"username" bson:"user"`
	Password    string                `json:"password" bson:"pass" encrypt:"true"`
	ArrayTest   []string              `json:"array" bson:"arrtest" encrypt:"true"`
	SkipMe      int                   `json:"skip" bson:"skipmee" encrypt:"true"` // should print an error
	Test        map[string]string     `json:"test" bson:"test" encrypt:"true"`    // should print an error
	Test2       map[int]string        `json:"test2" bson:"test2" encrypt:"true"`  // should print an error
	Test3       map[int]int           `json:"test3" bson:"test3" encrypt:"true"`  // should print an error
	Nested      nestedEncryptedStruct `json:"nested" bson:"nested"`               // encrypt flag does nothing on these, nested are automatically scanned for security
}

func getTestStruct() testStruct {
	return testStruct{
		DontInclude: "me",
		Username:    "kvizdos",
		Password:    "securepass123",
		ArrayTest:   []string{"value1", "value2"},
		SkipMe:      4,
		Test: map[string]string{
			"valueencrypted": "yes",
		},
		Test2: map[int]string{
			0: "no",
		},
		Test3: map[int]int{
			1: 2,
		},
		Nested: nestedEncryptedStruct{
			Value:   "yes",
			NonEncr: "no",
		},
	}
}

func getTestEngine() crypto.CryptoTestEngine {
	return crypto.CryptoTestEngine{
		Keys: map[string]string{
			"structs": "structureKey",
		},
	}
}

func TestEncryptStructJSON(t *testing.T) {
	testStr := getTestStruct()
	testEngine := getTestEngine()

	encrypted, error := EncryptStructIntoMap(testStr, JSON, testEngine)

	if error != nil {
		t.Errorf(error.Error())
		return
	}

	expect := map[string]interface{}{
		"username": "kvizdos",
		"password": "encrypted:securepass123",
		"array":    []string{"encrypted:value1", "encrypted:value2"},
		"skip":     4,
		"test": map[string]string{
			"valueencrypted": "encrypted:yes",
		},
		"test2": map[int]string{
			0: "no",
		},
		"test3": map[int]int{
			1: 2,
		},
	}

	if len(encrypted) != 7 {
		t.Errorf("expected length of 7, got %d", len(encrypted))
		return
	}

	for k, v := range encrypted {
		switch k {
		case "array":
			if expect[k].([]string)[0] != v.([]string)[0] || expect[k].([]string)[1] != v.([]string)[1] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test":
			if expect[k].(map[string]string)["valueencrypted"] != v.(map[string]string)["valueencrypted"] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test2":
			if expect[k].(map[int]string)[0] != v.(map[int]string)[0] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test3":
			if expect[k].(map[int]int)[0] != v.(map[int]int)[0] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		default:
			if expect[k] != v {
				t.Errorf("invalid: %s - %s", k, v)
			}
		}
	}
}

func TestEncryptStructBSON(t *testing.T) {
	testStr := getTestStruct()
	testEngine := getTestEngine()

	encrypted, error := EncryptStructIntoMap(testStr, BSON, testEngine)

	if error != nil {
		t.Errorf(error.Error())
		return
	}

	expect := map[string]interface{}{
		"user":    "kvizdos",
		"pass":    "encrypted:securepass123",
		"arrtest": []string{"encrypted:value1", "encrypted:value2"},
		"skipmee": 4,
		"test": map[string]string{
			"valueencrypted": "encrypted:yes",
		},
		"test2": map[int]string{
			0: "no",
		},
		"test3": map[int]int{
			1: 2,
		},
	}

	if len(encrypted) != 7 {
		t.Errorf("expected length of 7, got %d", len(encrypted))
		return
	}

	for k, v := range encrypted {
		switch k {
		case "arrtest":
			if expect[k].([]string)[0] != v.([]string)[0] || expect[k].([]string)[1] != v.([]string)[1] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test":
			if expect[k].(map[string]string)["valueencrypted"] != v.(map[string]string)["valueencrypted"] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test2":
			if expect[k].(map[int]string)[0] != v.(map[int]string)[0] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		case "test3":
			if expect[k].(map[int]int)[0] != v.(map[int]int)[0] {
				t.Errorf("invalid: %s - %+v", k, v)
			}
		default:
			if expect[k] != v {
				t.Errorf("invalid: %s - %v != %v", k, v, expect[k])
			}
		}
	}
}
