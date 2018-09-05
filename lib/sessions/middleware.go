package sessions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type guardian func(sess NoobSession) bool

func guard(g guardian, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := Default(c)
		defer session.Save()

		if !g(session) {
			session.AddFlash(message)
			c.Redirect(http.StatusSeeOther, "/")
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func LoggedIn() gin.HandlerFunc {
	return guard(
		NoobSession.IsLoggedIn,
		"Not logged in.",
	)
}

func Admin() gin.HandlerFunc {
	return guard(
		NoobSession.IsAdmin,
		"Insufficient permissions.",
	)
}
