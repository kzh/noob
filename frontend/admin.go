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
