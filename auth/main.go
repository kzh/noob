package main

import (
	"context"
	"log"
	"net/url"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
)

func main() {
	log.Println("Noob started.")

	// Create gin router
	r := gin.Default()

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

	// Connect to mongodb
	u := url.URL{}
	u.Scheme = "mongodb"
	u.User = url.UserPassword(
		"root",
		os.Getenv("MONGODB_PASSWORD"),
	)
	u.Host = "noob-mongodb:27071"

	client, err := mongo.Connect(context.TODO(), u.String())
	if err != nil {
		panic(err)
	}

	// Test mongodb
	dbs, err := client.ListDatabaseNames(context.TODO(), nil)
	if err != nil {
		panic(err)
	}

	log.Println("DBS: ")
	for _, db := range dbs {
		log.Println(db)
	}

	// Serve gin router
	r.Run()
}
