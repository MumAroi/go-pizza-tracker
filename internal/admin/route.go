package admin

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h Handler) {
	router.GET("/dashboard", h.Dashboard)
}
