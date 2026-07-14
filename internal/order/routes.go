package order

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, h Handler) {
	router.GET("/:id", h.ServeInfo)
	router.GET("/", h.ServeNewOrderForm)
	router.POST("/", h.HandleNewOrderPost)
}
