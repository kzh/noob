package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	noobdb "github.com/kzh/noob/pkg/database"
	"github.com/kzh/noob/pkg/model"
	noobsess "github.com/kzh/noob/pkg/sessions"
)

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

	var cred model.Credential
	if err := ctx.ShouldBind(&cred); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/login/"
		return
	}

	rec, err := noobdb.Authenticate(cred)
	if err != nil {
		session.AddFlash(err.Error())
		redirect = "/login/"
		return
	}

	session.SetM(gin.H{
		"username": cred.Username,
		"role":     rec["role"],
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

	var cred model.Credential
	if err := ctx.ShouldBind(&cred); err != nil {
		session.AddFlash("Invalid username or password.")
		redirect = "/register/"
		return
	}

	rec := bson.M{"role": "user"}
	if err := noobdb.Register(cred, rec); err != nil {
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

	r.Use(noobsess.Sessions())

	r.POST("/login", handleLogin)
	r.POST("/register", handleRegister)
	r.POST("/logout", handleLogout)

	// Serve gin router
	if err := r.Run(); err != nil {
		log.Println(err)
	}
}
