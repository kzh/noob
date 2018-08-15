package main

import (
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
)

func main() {
	log.Println("Noob started.")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	store, err := redis.NewStore(
		10,
		"tcp",
		"noob-redis-master:6379",
		os.Getenv("REDIS_PASSWORD"),
		// []byte(os.Getenv("SESSION_SECRET")),
	)
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("noob", store))

	log.Println("Connected to Redis.")

	// Connect to mongodb
	db := &mgo.DialInfo{
		Addrs:    []string{"noob-mongodb:27017"},
		Database: "admin",
		Username: "root",
		Password: os.Getenv("MONGODB_PASSWORD"),
	}

	log.Println("Connecting to MongoDB...")

	session, err := mgo.DialWithInfo(db)
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB.")
	log.Println("Inserting into MongoDB...")

	users := session.DB("noob").C("users")
	err = users.Insert(map[string]string{"hello": "world"})
	if err != nil {
		panic(err)
	}

	log.Println("Finished DB connections.")

	// Serve gin router
	r.Run()
}
