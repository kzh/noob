package sessions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type guardian func(sess NoobSession) bool

func guard(g guardian, message string, redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := Default(c)
		defer session.Save()

		if g(session) {
			c.Next()
			return
		}

		if redirect {
			session.AddFlash(message)
			c.Redirect(http.StatusSeeOther, "/")
		} else {
			c.JSON(
				http.StatusUnauthorized,
				gin.H{"error": message},
			)
		}

		c.Abort()
	}
}

func LoggedIn(redirect bool) gin.HandlerFunc {
	return guard(
		NoobSession.IsLoggedIn,
		"Not logged in.",
		redirect,
	)
}

func Admin(redirect bool) gin.HandlerFunc {
	return guard(
		NoobSession.IsAdmin,
		"Insufficient permissions.",
		redirect,
	)
}
