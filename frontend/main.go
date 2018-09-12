package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	noobsess "github.com/kzh/noob/lib/sessions"
)

func handleHome(c *gin.Context) {
	session := noobsess.Default(c)

	data := struct {
		User    string
		Admin   bool
		Message string
	}{}
	if session.IsLoggedIn() {
		data.User = session.Username()
	}
	messages := session.Flashes()
	if len(messages) > 0 {
		data.Message = messages[0].(string)
	}
	data.Admin = session.IsAdmin()

	session.Save()
	c.HTML(http.StatusOK, "index.tmpl", data)
}

func main() {
	log.Println("Noob: Frontend MS is starting...")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	r.Use(noobsess.Sessions())

	log.Println("Connected to Redis.")

	r.Static("/static", "./static/")
	r.LoadHTMLGlob("templates/*.tmpl")

	r.GET("/", handleHome)

	problems := r.Group("/")
	problems.GET("/problems/", handleProblems)
	problems.GET("/problem/:id/", handleProblem)

	// Auth
	auth := r.Group("/")
	auth.GET("/login/", handleLogin)
	auth.GET("/register/", handleRegister)

	// Admin
	admin := r.Group("/")
	admin.Use(noobsess.LoggedIn(true))
	admin.Use(noobsess.Admin(true))
	admin.GET("/create/", handleCreate)
	admin.GET("/problem/:id/edit/", handleEdit)

	if err := r.Run(); err != nil {
		log.Println(err)
	}
}
