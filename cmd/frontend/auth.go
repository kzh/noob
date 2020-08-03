package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	noobsess "github.com/kzh/noob/pkg/sessions"
)

func handleLogin(c *gin.Context) {
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
}

func handleRegister(c *gin.Context) {
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
}
