package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/kzh/noob/lib/message"
	"github.com/kzh/noob/lib/model"
	noobsess "github.com/kzh/noob/lib/sessions"
)

var upgrader = websocket.Upgrader{}

func handleSubmit(ctx *gin.Context) {
	c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer c.Close()

	var submission model.Submission
	if err := c.ReadJSON(&submission); err != nil {
		return
	}
	submission.ID = uuid.New().String()

	if err := message.Schedule(submission); err != nil {
		return
	}

	results, err := message.Subscribe(submission.ID)
	if err != nil {
		return
	}

	log.Println("Results:")
	for result := range results {
		log.Println(string(result.Body))
		err := c.WriteMessage(websocket.TextMessage, result.Body)
		if err != nil {
			return
		}
	}
}

func main() {
	log.Println("Noob: Submissions MS is starting...")

	r := gin.Default()
	r.Use(noobsess.Sessions())
	r.Use(noobsess.LoggedIn(false))

	r.GET("/submit", handleSubmit)

	if err := r.Run(); err != nil {
		log.Println(err)
	}
}
