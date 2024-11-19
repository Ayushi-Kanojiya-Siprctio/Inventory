package main

import (
	"log"
	"main/config"
	"main/controllers"
	"main/routes"
	_ "github.com/lib/pq" 
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	config.PostgresConnect()
	config.InitMongoDB()
	e := echo.New()
	inventoryController := &controllers.InventoryController{
		Validate: validator.New(),
	}
	routes.RegisterInventoryRoutes(e, inventoryController)

	log.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
