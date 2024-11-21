package main

import (
	"log"
	"main/config"
	"main/controllers"
	"main/routes"
	service "main/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	service.DB = config.PostgresDB
	//config.ConnectPostgres()
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