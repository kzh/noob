package db

import (
	"log"
	"os"
	"strconv"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var db *mgo.Database

func init() {
	log.Println("Connecting to MongoDB...")

	// Connect to mongodb
	dbInfo := &mgo.DialInfo{
		Addrs:    []string{"noob-mongodb:27017"},
		Database: "admin",
		Username: "root",
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	session, err := mgo.DialWithInfo(dbInfo)
	if err != nil {
		panic(err)
	}

	db = session.DB("noob")
	log.Println("Connected to MongoDB.")
}

func count(name string) (string, error) {
	counters := db.C("counters")

	query := bson.M{"_id": name}
	update := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"count": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	var doc bson.M
	_, err := counters.Find(query).Apply(update, &doc)
	if err != nil {
		return "", err
	}

	id := strconv.Itoa(doc["count"].(int))
	return id, nil
}
