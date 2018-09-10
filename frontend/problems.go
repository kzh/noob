package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	noobsess "github.com/kzh/noob/lib/sessions"
)

func handleCreate(c *gin.Context) {
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
}

func handleProblems(c *gin.Context) {
	session := noobsess.Default(c)

	data := struct {
		Admin   bool
		Message string
	}{}
	messages := session.Flashes()
	if len(messages) > 0 {
		data.Message = messages[0].(string)
	}
	data.Admin = session.IsAdmin()

	session.Save()
	c.HTML(http.StatusOK, "problems.tmpl", data)
}

func handleProblem(c *gin.Context) {
	session := noobsess.Default(c)

	data := struct {
		Message string
	}{}
	messages := session.Flashes()
	if len(messages) > 0 {
		data.Message = messages[0].(string)
	}

	session.Save()
	c.HTML(http.StatusOK, "problem.tmpl", data)
}
