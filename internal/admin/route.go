package admin

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h Handler) {
	router.GET("/notifications", h.GetNotification)
	router.GET("/dashboard", h.Dashboard)
	router.POST("/orders/:id/update", h.OrderPut)
	router.POST("/orders/:id/delete", h.OrderDelete)
}
