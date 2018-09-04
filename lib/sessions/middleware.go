package sessions

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := Default(c)
		defer session.Save()

		if !session.IsLoggedIn() {
			session.AddFlash("Not logged in.")
			c.Redirect(http.StatusSeeOther, "/")
		} else {
			c.Next()
		}
	}
}
