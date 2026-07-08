package route

import (
	"pizza-tracker/internal/app"
	"pizza-tracker/internal/order"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, app *app.App) {
	order.RegisterRoutes(router.Group("/orders"), app.OrderRepo)

	orderHandler := order.NewHandler(app.OrderRepo)
	router.GET("/customers/:id", orderHandler.ServeCustomer)
	router.Static("/static", "/templates/static")
}
