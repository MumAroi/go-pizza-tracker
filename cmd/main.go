package main

import (
	"log"

	"pizza-tracker/internal/app"

	"github.com/gin-gonic/gin"
)

func main() {
	app, err := app.NewApp("pizza.db")
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	// orderHandler := order.NewHandler(app.OrderRepo)

	router := gin.Default()

	router.Run(":3003")
}
