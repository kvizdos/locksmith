//go:build e2e
// +build e2e

package database

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

var db MongoDatabase

func TestMain(m *testing.M) {
	ctx, timeout := context.WithTimeout(context.Background(), 20*time.Second)
	db = MongoDatabase{
		Ctx:    ctx,
		Cancel: timeout,
	}
	err := db.Initialize("mongodb://localhost:27017", "testingdb")

	if err != nil {
		fmt.Println(err)
		return
	}

	code := m.Run()

	// Clean up
	db.database.Drop(context.Background())

	os.Exit(code)
}

func TestMongoDeleteValue(t *testing.T) {
	_, err := db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"hello":    "world",
	})

	if err != nil {
		t.Error(err)
		return
	}

	deleted, err := db.DeleteOne("test", map[string]interface{}{
		"username": "kenton",
		"hello":    "world",
	})

	if err != nil {
		t.Error(err)
		return
	}

	if !deleted {
		t.Error("item not delted.")
	}
}

func TestMongoInsertValue(t *testing.T) {
	_, err := db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"hello":    "world",
	})

	if err != nil {
		t.Error(err)
		return
	}
}

func TestMongoInsertStructValue(t *testing.T) {
	type test struct {
		Value string `bson:"value"`
	}
	_, err := db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"hello": test{
			Value: "world",
		},
	})

	if err != nil {
		t.Error(err)
		return
	}

	t.Cleanup(func() {
		db.database.Drop(context.Background())
	})
}

func TestMongoFindOne(t *testing.T) {
	type test struct {
		Value string `bson:"value"`
	}

	_, err := db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"hello": test{
			Value: "world",
		},
		"arraytest": []string{"testing", "1", "2", "3"},
		"itsint":    1,
	})

	if err != nil {
		t.Errorf("failed insertion prereq: %s", err.Error())
		return
	}

	_, found := db.FindOne("test", map[string]interface{}{
		"username": "kenton",
	})

	if found == false {
		t.Errorf("didnt find items")
		return
	}

	t.Cleanup(func() {
		db.database.Drop(context.Background())
	})
}

func TestMongoFindMany(t *testing.T) {
	type test struct {
		Value string `bson:"value"`
	}

	_, err := db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"hello": test{
			Value: "world",
		},
		"arraytest": []string{"testing", "1", "2", "3"},
		"itsint":    1,
	})

	if err != nil {
		t.Errorf("failed insertion prereq: %s", err.Error())
		return
	}

	_, err = db.InsertOne("test", map[string]interface{}{
		"username": "bob",
		"hello": test{
			Value: "world",
		},
	})

	if err != nil {
		t.Errorf("failed insertion prereq: %s", err.Error())
		return
	}

	res, found := db.Find("test", map[string]interface{}{})

	if found == false {
		t.Errorf("didnt find items")
		return
	}

	expecting := 2
	if len(res) != 2 {
		t.Errorf("expected %d users got %d", expecting, len(res))
	}

	t.Cleanup(func() {
		db.database.Drop(context.Background())
	})
}

func TestMongoUpdateOnePUSH(t *testing.T) {
	db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"sessions": []string{"exists already"},
	})

	_, err := db.UpdateOne("test", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		PUSH: {
			"sessions": "i'm a session!",
		},
	})

	if err != nil {
		t.Errorf("failed to update: %s", err.Error())
		return
	}

	v, found := db.FindOne("test", map[string]interface{}{
		"username": "kenton",
	})

	if found == false {
		t.Errorf("didnt fidn items")
		return
	}

	value := v.(map[string]interface{})

	if len(value["sessions"].([]interface{})) != 2 {
		t.Error("did not push anything")
		return
	}

	t.Cleanup(func() {
		db.database.Drop(context.Background())
	})
}

func TestMongoUpdateOneSET(t *testing.T) {
	db.InsertOne("test", map[string]interface{}{
		"username": "kenton",
		"sessions": []string{"exists already"},
	})

	_, err := db.UpdateOne("test", map[string]interface{}{
		"username": "kenton",
	}, map[DatabaseUpdateActions]map[string]interface{}{
		SET: {
			"username": "bob",
		},
	})

	if err != nil {
		t.Errorf("failed to update: %s", err.Error())
		return
	}

	v, found := db.FindOne("test", map[string]interface{}{
		"username": "bob",
	})

	if found == false {
		t.Errorf("didnt find items")
		return
	}

	value := v.(map[string]interface{})

	if value["username"].(string) != "bob" {
		t.Error("did not set")
	}

	t.Cleanup(func() {
		db.database.Drop(context.Background())
	})
}
