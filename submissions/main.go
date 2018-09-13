package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	//	noobdb "github.com/kzh/noob/lib/database"
	"github.com/kzh/noob/lib/model"
	_ "github.com/kzh/noob/lib/queue"
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
}

func main() {
	log.Println("Noob: Problems MS is starting...")

	r := gin.Default()

	log.Println("Connecting to Redis...")
	r.Use(noobsess.Sessions())
	log.Println("Connected to Redis.")

	r.Use(noobsess.LoggedIn(false))

	r.POST("/submit", handleSubmit)

	r.Run()
}
