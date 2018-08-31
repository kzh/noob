package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Println("Noob: Frontend MS is starting...")

	// Create gin router
	r := gin.Default()

	log.Println("Connecting to Redis...")

	// Use redis sessions middleware
	store, err := redis.NewStore(
		10,
		"tcp",
		"noob-redis-master:6379",
		os.Getenv("REDIS_PASSWORD"),
		[]byte("NOOB_SESSION_SECRET"),
		//[]byte(os.Getenv("SESSION_SECRET")),
	)
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("noob", store))

	log.Println("Connected to Redis.")

	// Load templates
	r.LoadHTMLGlob("templates/*")

	// Get / Handler
	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)

		data := struct {
			User    string
			Message string
		}{}
		if username, ok := session.Get("username").(string); ok {
			data.User = username
		}
		messages := session.Flashes()
		if len(messages) > 0 {
			data.Message = messages[0].(string)
			log.Println("Message: " + data.Message)
		}

		session.Save()
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	// Get /login/ Handler
	r.GET("/login/", func(c *gin.Context) {
		session := sessions.Default(c)

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
		session := sessions.Default(c)

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

	r.Run()
}
