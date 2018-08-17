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

type Credentials struct {
	Username string `form:"username" json:"username" binding: "required"`
	Password string `form:"password" json:"password" binding: "required"`
}

type Users struct {
	c *mgo.Collection
}

func (u *Users) handleLogin(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username, ok := session.Get("username").(string)
	if ok && username != "" {
		ctx.JSON(http.StatusOK, gin.H{"message": "already logged in"})
		return
	}

	var creds Credentials
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username and password"})
		return
	}

	filter := bson.M{"username": creds.Username}

	var rec bson.M
	if err := u.c.Find(filter).One(&rec); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username and password"})
		return
	}

	hash := rec["password"].(string)
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(creds.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username and password"})
		return
	}

	session.Set("username", creds.Username)
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (u *Users) handleRegister(ctx *gin.Context) {
	session := sessions.Default(ctx)
	username, ok := session.Get("username").(string)
	if ok && username != "" {
		ctx.JSON(http.StatusOK, gin.H{"message": "already logged in"})
		return
	}

	var creds Credentials
	if err := ctx.ShouldBindJSON(&creds); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username and password"})
		return
	}

	filter := bson.M{"username": creds.Username}

	var rec bson.M
	if err := u.c.Find(filter).One(&rec); err == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "username taken"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	rec = bson.M{"username": creds.Username, "password": string(hash)}
	if err := u.c.Insert(rec); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func handleLogout(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if session.Get("username") == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return
	}

	session.Clear()
	session.Save()

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

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

	users := &Users{session.DB("noob").C("users")}
	r.POST("/login", users.handleLogin)
	r.POST("/register", users.handleRegister)

	r.POST("/logout", handleLogout)

	// Serve gin router
	r.Run()
}
