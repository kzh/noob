package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/kzh/noob/lib/model"
	"github.com/kzh/noob/lib/queue"
	noobsess "github.com/kzh/noob/lib/sessions"
)

func handleSubmit(ctx *gin.Context) {
	var submission model.Submission
	if err := ctx.ShouldBind(&submission); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := queue.Schedule(submission); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success!",
	})
}

func main() {
	log.Println("Noob: Problems MS is starting...")

	r := gin.Default()
	r.Use(noobsess.Sessions())
	r.Use(noobsess.LoggedIn(false))

	r.POST("/submit", handleSubmit)

	if err := r.Run(); err != nil {
		log.Println(err)
	}
}
