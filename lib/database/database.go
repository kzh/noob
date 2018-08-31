package db

import (
	"errors"
	"log"
	"os"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
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

	log.Println("Connecting to MongoDB...")

	session, err := mgo.DialWithInfo(dbInfo)
	if err != nil {
		panic(err)
	}

	db = session.DB("noob")
	log.Println("Connected to MongoDB.")
}

var (
	ErrInvalidCredential   = errors.New("Invalid username or password.")
	ErrUnavailableUsername = errors.New("Username taken.")
	ErrInternalServer      = errors.New("Internal server error.")
)

func Authenticate(username, password string) (bson.M, error) {
	users := db.C("users")

	var rec bson.M
	query := bson.M{"username": username}
	if err := users.Find(query).One(&rec); err != nil {
		return nil, ErrInvalidCredential
	}

	hash := rec["password"].(string)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return nil, ErrInvalidCredential
	}

	delete(rec, "password")
	return rec, nil
}

func Register(username, password string, data bson.M) error {
	users := db.C("users")

	var rec bson.M
	query := bson.M{"username": username}
	if err := users.Find(query).One(&rec); err == nil {
		return ErrUnavailableUsername
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	rec = bson.M{
		"username": username,
		"password": string(hash),
	}

	for k, v := range data {
		rec[k] = v
	}

	if err := users.Insert(rec); err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	return nil
}
