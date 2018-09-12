package db

import (
	"log"

	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"

	. "github.com/kzh/noob/lib/model"
)

func Authenticate(cred Credential) (bson.M, error) {
	users := db.C("users")

	var rec bson.M
	query := bson.M{"username": cred.Username}
	if err := users.Find(query).One(&rec); err != nil {
		return nil, ErrInvalidCredential
	}

	hash := rec["password"].(string)
	if err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(cred.Password),
	); err != nil {
		return nil, ErrInvalidCredential
	}

	delete(rec, "password")
	return rec, nil
}

func Register(cred Credential, data bson.M) error {
	users := db.C("users")

	var rec bson.M
	query := bson.M{"username": cred.Username}
	if err := users.Find(query).One(&rec); err == nil {
		return ErrUnavailableUsername
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(cred.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Println(err)
		return ErrInternalServer
	}

	rec = bson.M{
		"username": cred.Username,
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
