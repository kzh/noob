package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	noobdb "github.com/kzh/noob/lib/database"
	noobsess "github.com/kzh/noob/lib/sessions"
)

func handleCreate(ctx *gin.Context) {
	var redirect string
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	if !session.IsAdmin() {
		redirect = "/"
		session.AddFlash("Insufficient permissions.")
		return
	}

	var prob noobdb.ProblemData
	if err := ctx.ShouldBind(&prob); err != nil {
		redirect = "/"
		session.AddFlash("Invalid form data format")
		log.Println(err)
		return
	}

	problem, err := noobdb.CreateProblem(prob)
	if err != nil {
		redirect = "/"
	}
	log.Println(problem)

	session.AddFlash("Success!")
	redirect = "/"
}

func handleEdit(ctx *gin.Context) {
	var redirect string
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	if !session.IsAdmin() {
		redirect = "/"
		session.AddFlash("Insufficient permissions.")
		return
	}

	var prob noobdb.Problem
	if err := ctx.ShouldBind(&prob); err != nil {
		redirect = "/"
		session.AddFlash("Invalid form data format")
		log.Println(err)
		return
	}

	err := noobdb.EditProblem(prob)
	if err != nil {
		redirect = "/"
	}
	log.Println(prob)

	session.AddFlash("Success!")
	redirect = "/"
}

func handleDelete(ctx *gin.Context) {
	var redirect string
	session := noobsess.Default(ctx)

	defer func() {
		session.Save()
		ctx.Redirect(http.StatusSeeOther, redirect)
	}()

	if !session.IsAdmin() {
		redirect = "/"
		session.AddFlash("Insufficient permissions.")
		return
	}

	var prob noobdb.ProblemID
	if err := ctx.ShouldBind(&prob); err != nil {
		redirect = "/"
		session.AddFlash("Invalid form data format")
		log.Println(err)
		return
	}

	err := noobdb.DeleteProblem(prob)
	if err != nil {
		redirect = "/"
	}
	log.Println(prob)

	session.AddFlash("Success!")
	redirect = "/"

}

func main() {
	log.Println("Noob: Admin MS is starting...")

	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	r.Use(noobsess.Sessions())
	r.Use(noobsess.LoggedIn())

	log.Println("Connected to Redis.")

	r.POST("/create", handleCreate)
	r.POST("/edit", handleEdit)
	r.POST("/delete", handleDelete)

	r.Run()
}
