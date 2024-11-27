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
	config.PostgresConnect()
	if err := service.CreateTableIfNotExists(config.PG); err != nil {
		log.Fatal("Error creating tables:", err)
	}
	config.InitMongoDB()

	e := echo.New()

	inventoryController := &controllers.InventoryController{
		Validate: validator.New(),
	}

	routes.RegisterInventoryRoutes(e, inventoryController)
	var mongo config.MongoConfig
	port := mongo.MongoPort
	if port == "" {
		port = ":8080"
	}

	if err := e.Start(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
