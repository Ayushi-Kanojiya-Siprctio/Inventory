package main

import (
	"log"
	"main/config"
	"main/controllers"
	"main/routes"
	service "main/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func main() {
	// Initialize PostgreSQL connection
	config.PostgresConnect() // no value expected, just initialization
	if err := service.CreateTableIfNotExists(config.PG); err != nil {
		log.Fatal("Error creating tables:", err)
	}
	// Initialize MongoDB connection
	config.InitMongoDB()

	// Create a new Echo instance
	e := echo.New()

	// Create an InventoryController

	inventoryController := &controllers.InventoryController{
		Validate: validator.New(),
		
	}

	// Register routes
	routes.RegisterInventoryRoutes(e, inventoryController)

	// Start the server
	log.Println("Starting server on port 8080...")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
