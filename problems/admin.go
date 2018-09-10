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
