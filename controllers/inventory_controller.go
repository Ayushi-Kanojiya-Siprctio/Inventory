package controllers

import (
	"log"
	manager "main/managers"
	"main/models"
	"main/requests"
	"main/responses"

	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type InventoryController struct {
	Validate         *validator.Validate
	InventoryManager *manager.InventoryManager
}

func parseFlag(flag string) (bool, error) {
	if flag == "" {
		return false, nil
	}
	return strconv.ParseBool(flag)
}

func (c *InventoryController) CreateItemHandler(ctx echo.Context) error {
	flagStr := ctx.QueryParam("flag")
	flag, err := parseFlag(flagStr)
	log.Println("flag--------->", flag)

	var req requests.InventoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := c.Validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	item := &models.Inventory{
		Name:     req.Name,
		Price:    req.Price,
		Currency: req.Currency,
		Discount: req.Discount,
		Vendor:   req.Vendor,
	}
	log.Println("controller item-------->", item.Name)

	createdItem, err := c.InventoryManager.CreateItem(ctx.Request().Context(), flag, item)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create inventory item"})
	}

	response := responses.InventoryResponse{
		// ID:       createdItem.ID,
		Name:     createdItem.Name,
		Price:    createdItem.Price,
		Currency: createdItem.Currency,
		Discount: createdItem.Discount,
		Vendor:   createdItem.Vendor,
	}

	return ctx.JSON(http.StatusCreated, response)
}


func (c *InventoryController) GetItemsHandler(ctx echo.Context) error {
	flagStr := ctx.QueryParam("flag")
	flag, err := parseFlag(flagStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	items, totalCount, err := c.InventoryManager.GetItems(ctx.Request().Context(), flag)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if len(items) == 0 {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": "No items found",
		})
	}

	var itemResponses []responses.InventoryResponse
	for _, item := range items {
		itemResponses = append(itemResponses, responses.InventoryResponse{
			// ID:       item.ID,
			Name:     item.Name,
			Price:    item.Price,
			Currency: item.Currency,
			Discount: item.Discount,
			Vendor:   item.Vendor,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"items":        itemResponses,
		"totalRecords": totalCount,
	})
}

func (c *InventoryController) GetItemByIDHandler(ctx echo.Context) error {
	flagStr := ctx.QueryParam("flag")
	flag, err := parseFlag(flagStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid flag"})
	}

	id := ctx.Param("id")

	item, err := c.InventoryManager.GetItemByID(ctx.Request().Context(), flag, id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Item not found"})
	}

	return ctx.JSON(http.StatusOK, responses.InventoryResponse{
		// ID:       item.ID,
		Name:     item.Name,
		Price:    item.Price,
		Currency: item.Currency,
		Discount: item.Discount,
		Vendor:   item.Vendor,
	})
}

func (c *InventoryController) UpdateItemHandler(ctx echo.Context) error {
	flagStr := ctx.QueryParam("flag")
	flag, err := parseFlag(flagStr)
	id := ctx.Param("id")
	var req requests.InventoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	item := &models.Inventory{
		Name:     req.Name,
		Price:    req.Price,
		Currency: req.Currency,
		Discount: req.Discount,
		Vendor:   req.Vendor,
	}

	updatedItem, err := c.InventoryManager.UpdateItem(ctx.Request().Context(), flag, id, item)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update item"})
	}

	return ctx.JSON(http.StatusOK, responses.InventoryResponse{
		// ID:       updatedItem.ID,
		Name:     updatedItem.Name,
		Price:    updatedItem.Price,
		Currency: updatedItem.Currency,
		Discount: updatedItem.Discount,
		Vendor:   updatedItem.Vendor,
	})
}
func (c *InventoryController) DeleteItemHandler(ctx echo.Context) error {
	flagStr := ctx.QueryParam("flag")
	flag, err := parseFlag(flagStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid flag value"})
	}

	id := ctx.Param("id")

	err = c.InventoryManager.DeleteItem(ctx.Request().Context(), flag, id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete item"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Item deleted successfully"})
}


