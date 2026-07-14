package middleware

import (
	"net/http"
	"pizza-tracker/internal/shared/util"
	"pizza-tracker/internal/user"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(userRepo user.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := util.GetSessionString(c, "userID")

		if userID == "" {
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		_, err := userRepo.GetByID(userID)
		if err != nil {
			util.ClearSession(c)
			c.Redirect(http.StatusSeeOther, "/login")
			c.Abort()
			return
		}

		c.Next()
	}
}
