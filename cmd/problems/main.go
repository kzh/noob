package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	noobdb "github.com/kzh/noob/pkg/database"
	noobsess "github.com/kzh/noob/pkg/sessions"
)

func handleList(ctx *gin.Context) {
	problems, err := noobdb.Problems()
	if err == nil {
		ctx.JSON(http.StatusOK, problems)
		return
	}

	status := http.StatusInternalServerError
	ctx.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func handleSelect(ctx *gin.Context) {
	id := ctx.Param("id")
	problem, err := noobdb.SnapProblem(id)
	if err == nil {
		ctx.JSON(http.StatusOK, problem)
		return
	}

	status := http.StatusInternalServerError
	if err == noobdb.ErrNoSuchProblem {
		status = http.StatusNotFound
	}

	ctx.JSON(status, gin.H{
		"error": err.Error(),
	})
}

func main() {
	log.Println("Noob: Problems MS is starting...")

	r := gin.Default()

	r.Use(noobsess.Sessions())

	r.GET("/list", handleList)
	r.GET("/get/:id", handleSelect)

	admin := r.Group("/")
	admin.Use(noobsess.LoggedIn(true))
	admin.Use(noobsess.Admin(true))

	admin.POST("/create", handleCreate)
	admin.POST("/edit", handleEdit)
	admin.POST("/delete", handleDelete)

	adminNR := r.Group("/")
	adminNR.Use(noobsess.LoggedIn(false))
	adminNR.Use(noobsess.Admin(false))

	adminNR.GET("/get/:id/io", handleSelectIO)

	if err := r.Run(); err != nil {
		log.Println(err)
	}
}
