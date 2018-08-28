package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/crypto/bcrypt"
)

type Credential struct {
	Username string `form:"username" binding: "required"`
	Password string `form:"password" binding: "required"`
}

type Users struct {
	c *mgo.Collection
}

func (u *Users) handleLogin(ctx *gin.Context) {
	var redirect string
	session := sessions.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	username, ok := session.Get("username").(string)
	if ok && username != "" {
		session.AddFlash("Already logged in.")
		redirect = "/"
		return
	}

	var creds Credential
	if err := ctx.ShouldBind(&creds); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/login/"
		return
	}

	filter := bson.M{"username": creds.Username}

	var rec bson.M
	if err := u.c.Find(filter).One(&rec); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/login/"
		return
	}

	hash := rec["password"].(string)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/login/"
		return
	}

	session.Set("username", creds.Username)
	session.AddFlash("Success!")

	redirect = "/"
}

func (u *Users) handleRegister(ctx *gin.Context) {
	var redirect string
	session := sessions.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	username, ok := session.Get("username").(string)
	if ok && username != "" {
		session.AddFlash("Already logged in.")
		redirect = "/"
		return
	}

	var creds Credential
	if err := ctx.ShouldBind(&creds); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/register/"
		return
	}

	filter := bson.M{"username": creds.Username}

	var rec bson.M
	if err := u.c.Find(filter).One(&rec); err == nil {
		session.AddFlash("Username taken.")
		redirect = "/register/"
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		session.AddFlash("Internal server error.")
		redirect = "/register/"
		return
	}

	rec = bson.M{"username": creds.Username, "password": string(hash)}
	if err := u.c.Insert(rec); err != nil {
		session.AddFlash("Internal server error.")
		redirect = "/register/"
		return
	}

	session.AddFlash("Success!")

	redirect = "/"
}

func handleLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	defer session.Save()
	defer ctx.Redirect(http.StatusSeeOther, "/")

	if session.Get("username") == nil {
		session.AddFlash("Not logged in.")
		return
	}

	session.Clear()
	session.AddFlash("Success!")
}

func main() {
	log.Println("Noob Auth started.")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	store, err := redis.NewStore(
		10,
		"tcp",
		"noob-redis-master:6379",
		os.Getenv("REDIS_PASSWORD"),
		[]byte("NOOB_SESSION_SECRET"),
		//[]byte(os.Getenv("SESSION_SECRET")),
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

	users := &Users{session.DB("noob").C("users")}
	r.POST("/login", users.handleLogin)
	r.POST("/register", users.handleRegister)

	r.POST("/logout", handleLogout)

	// Serve gin router
	r.Run()
}
