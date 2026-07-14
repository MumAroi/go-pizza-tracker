package main

import (
	"log/slog"
	"os"

	"pizza-tracker/internal/app"
	"pizza-tracker/internal/config"
	"pizza-tracker/internal/order"
	"pizza-tracker/internal/route"
	"pizza-tracker/internal/session"
	"pizza-tracker/internal/shared/util"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		slog.Error("Error loading config", "error", err)
		os.Exit(1)
	}

	app, err := app.NewApp(cfg.DBPath)
	if err != nil {
		slog.Error("failed to create app", "error", err)
		os.Exit(1)
	}
	defer app.Close()

	order.RegisterCustomValidators()

	router := gin.Default()

	util.LoadTemplate(router)

	sessionStore := session.NewSessionStore(app.DB, []byte(cfg.SessionSecret))

	route.SetupRoutes(router, sessionStore, app)

	router.Run(":" + cfg.Port)
}
