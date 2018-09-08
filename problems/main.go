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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, problems)
}

func handleSelect(ctx *gin.Context) {
	id := ctx.Param("id")
	problem, err := noobdb.ProblemFromID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, problem)
}

func main() {
	log.Println("Noob: Problems MS is starting...")

	r := gin.Default()

	log.Println("Connecting to Redis...")
	r.Use(noobsess.Sessions())
	log.Println("Connected to Redis.")

	r.GET("/list", handleList)
	r.GET("/get/:id", handleSelect)

	r.Run()
}
