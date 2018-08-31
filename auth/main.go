package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	noobdb "github.com/kzh/noob/lib/database"
	noobsess "github.com/kzh/noob/lib/sessions"
)

type Credential struct {
	Username string `form:"username" binding: "required"`
	Password string `form:"password" binding: "required"`
}

func handleLogin(ctx *gin.Context) {
	var redirect string
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	if session.IsLoggedIn() {
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

	rec, err := noobdb.Authenticate(creds.Username, creds.Password)
	if err != nil {
		session.AddFlash(err.Error())
		redirect = "/login/"
		return
	}

	session.Login(gin.H{
		"username": creds.Username,
		"role":     rec["role"].(string),
	})
	session.AddFlash("Success!")
	redirect = "/"
}

func handleRegister(ctx *gin.Context) {
	var redirect string
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	if session.IsLoggedIn() {
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

	rec := bson.M{"role": "user"}
	if err := noobdb.Register(creds.Username, creds.Password, rec); err != nil {
		session.AddFlash(err.Error())
		redirect = "/register/"
	}

	session.AddFlash("Success!")
	redirect = "/"
}

func handleLogout(ctx *gin.Context) {
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, "/")
	}()

	if !session.IsLoggedIn() {
		session.AddFlash("Not logged in.")
		return
	}

	session.Clear()
	session.AddFlash("Success!")
}

func main() {
	log.Println("Noob: Auth MS is starting...")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")
	r.Use(noobsess.Sessions())
	log.Println("Connected to Redis.")

	r.POST("/login", handleLogin)
	r.POST("/register", handleRegister)
	r.POST("/logout", handleLogout)

	// Serve gin router
	r.Run()
}
