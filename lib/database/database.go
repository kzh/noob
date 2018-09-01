package db

import (
	"log"
	"os"

	"github.com/globalsign/mgo"
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
