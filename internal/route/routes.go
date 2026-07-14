package route

import (
	"pizza-tracker/internal/admin"
	"pizza-tracker/internal/app"
	"pizza-tracker/internal/middleware"
	"pizza-tracker/internal/order"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, store sessions.Store, app *app.App) {

	router.Use(sessions.Sessions("pizza-tracker", store))

	orderH := order.NewHandler(order.OrderDeps{
		OrderRepo: app.OrderRepo,
	})
	order.RegisterRoutes(router.Group("/orders"), orderH)

	adminH := admin.NewHandler(admin.AdminDeps{
		UserRepo: app.UserRepo,
	})
	adminRouter := router.Group("/admin", middleware.AuthMiddleware(app.UserRepo))
	admin.RegisterRoutes(adminRouter, adminH)

	router.GET("/login", adminH.RenderLogin)
	router.POST("/login", adminH.Login)
	router.POST("/logout", adminH.Logout)

	router.Static("/static", "./templates/static")
}
