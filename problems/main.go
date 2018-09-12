package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	noobdb "github.com/kzh/noob/lib/database"
	noobsess "github.com/kzh/noob/lib/sessions"
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

	log.Println("Connecting to Redis...")
	r.Use(noobsess.Sessions())
	log.Println("Connected to Redis.")

	r.GET("/list", handleList)
	r.GET("/get/:id", handleSelect)

	admin := r.Group("/")
	admin.Use(noobsess.LoggedIn())
	admin.Use(noobsess.Admin())

	admin.POST("/create", handleCreate)
	admin.POST("/edit", handleEdit)
	admin.POST("/delete", handleDelete)
	admin.GET("/get/:id/io", handleSelectIO)

	r.Run()
}
