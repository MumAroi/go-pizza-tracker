package main

import (
	"log"

	"pizza-tracker/internal/app"
	"pizza-tracker/internal/config"
	"pizza-tracker/internal/order"
	"pizza-tracker/internal/route"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Error loading config")
	}

	app, err := app.NewApp(cfg.DBPath)
	if err != nil {
		log.Fatal("Error creating app")
	}
	defer app.Close()

	order.RegisterCustomValidators()

	router := gin.Default()
	route.SetupRoutes(router, app)

	router.Run(":" + cfg.Port)
}
