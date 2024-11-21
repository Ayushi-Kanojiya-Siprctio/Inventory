package routes

import (
	"main/controllers"
	manager "main/managers"

	"github.com/labstack/echo/v4"
)

func RegisterInventoryRoutes(e *echo.Echo, inventoryController *controllers.InventoryController) {
	inventoryManager := &manager.InventoryManager{}
	inventoryController.InventoryManager = inventoryManager
	e.POST("/inventory", inventoryController.CreateItemHandler)
	e.GET("/inventory", inventoryController.GetItemsHandler)
	e.GET("/inventory/:id", inventoryController.GetItemByIDHandler)
	e.PUT("/inventory/:id", inventoryController.UpdateItemHandler)
	e.DELETE("/inventory/:id", inventoryController.DeleteItemHandler)
}