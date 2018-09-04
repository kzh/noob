package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	noobsess "github.com/kzh/noob/lib/sessions"
)

func main() {
	log.Println("Noob: Frontend MS is starting...")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	r.Use(noobsess.Sessions())

	log.Println("Connected to Redis.")

	// Load templates
	r.LoadHTMLGlob("templates/*.tmpl")

	// Get / Handler
	r.GET("/", func(c *gin.Context) {
		session := noobsess.Default(c)

		data := struct {
			User    string
			Admin   bool
			Message string
		}{}
		if session.IsLoggedIn() {
			data.User = session.Username()
			log.Println("logged in")
		}
		messages := session.Flashes()
		if len(messages) > 0 {
			data.Message = messages[0].(string)
		}
		data.Admin = session.IsAdmin()

		session.Save()
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	// Get /login/ Handler
	r.GET("/login/", func(c *gin.Context) {
		session := noobsess.Default(c)

		data := struct {
			Message string
		}{}
		messages := session.Flashes()
		if len(messages) > 0 {
			data.Message = messages[0].(string)
		}

		session.Save()
		c.HTML(http.StatusOK, "login.tmpl", data)
	})

	// Get /register/ Handler
	r.GET("/register/", func(c *gin.Context) {
		session := noobsess.Default(c)

		data := struct {
			Message string
		}{}
		messages := session.Flashes()
		if len(messages) > 0 {
			data.Message = messages[0].(string)
		}

		session.Save()
		c.HTML(http.StatusOK, "register.tmpl", data)
	})

	r.GET("/create/", func(c *gin.Context) {
		session := noobsess.Default(c)

		data := struct {
			Message string
		}{}
		messages := session.Flashes()
		if len(messages) > 0 {
			data.Message = messages[0].(string)
		}

		session.Save()
		c.HTML(http.StatusOK, "create.tmpl", data)
	})

	r.Run()
}
