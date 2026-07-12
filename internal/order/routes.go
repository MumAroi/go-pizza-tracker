package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.RouterGroup, repo OrderRepository) {
	handler := NewHandler(repo)
	router.GET("/:id", handler.ServeInfo)
	router.GET("/", handler.ServeNewOrderForm)
	router.POST("/", handler.HandleNewOrderPost)

}
