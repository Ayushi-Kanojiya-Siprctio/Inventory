package controllers

import (
	"fmt"
	manager "main/managers"
	"main/models"
	"main/requests"
	"main/responses"

	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type InventoryController struct {
	Validate         *validator.Validate
	InventoryManager *manager.InventoryManager
}

func (c *InventoryController) CreateItemHandler(ctx echo.Context) error {
	var req requests.InventoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	if err := c.Validate.Struct(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	item := &models.Inventory{
		Name:        req.Name,
		Price:       req.Price,
		Currency:    req.Currency,
		Discount:    req.Discount,
		Vendor:      req.Vendor,
		Accessories: req.Accessories,
	}

	createdItem, err := c.InventoryManager.CreateItem(ctx.Request().Context(), item)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create inventory item"})
	}

	// Convert `primitive.ObjectID` to string for the response
	response := responses.InventoryResponse{
		ID:          createdItem.ID.Hex(),
		Name:        createdItem.Name,
		Price:       createdItem.Price,
		Currency:    createdItem.Currency,
		Discount:    createdItem.Discount,
		Vendor:      createdItem.Vendor,
		Accessories: createdItem.Accessories,
	}

	return ctx.JSON(http.StatusCreated, response)
}


func (c *InventoryController) GetItemsHandler(ctx echo.Context) error {
	pageParam := ctx.QueryParam("page")
	pageSizeParam := ctx.QueryParam("pageSize")
	vendorParam := ctx.QueryParam("vendor")

	page := 1
	pageSize := 100

	if pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil || page < 0 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "'page' must be a positive integer"})
		}
	}

	if pageSizeParam != "" {
		var err error
		pageSize, err = strconv.Atoi(pageSizeParam)
		if err != nil || pageSize <= -2 {
			return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "'pageSize' must be a positive integer"})
		}
	}

	if pageSize == -1 {
		page = 1
		pageSize = 100
	}

	var vendors []string
	if vendorParam != "" {
		vendors = strings.Split(vendorParam, ",")
	}

	items, totalCount, err := c.InventoryManager.GetItems(ctx.Request().Context(), page, pageSize, vendors)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if len(items) == 0 {
		return ctx.JSON(http.StatusNotFound, map[string]string{
			"message": fmt.Sprintf("No record found with these vendors: %s", vendorParam),
		})
	}

	var itemResponses []responses.InventoryResponse
	for _, item := range items {
		itemResponses = append(itemResponses, responses.InventoryResponse{
			ID:          item.ID.Hex(),
			Name:        item.Name,
			Price:       item.Price,
			Currency:    item.Currency,
			Discount:    item.Discount,
			Vendor:      item.Vendor,
			Accessories: item.Accessories,
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"items":        itemResponses,
		"totalRecords": totalCount,
		"currentPage":  page,
	})
}

func (c *InventoryController) GetItemByIDHandler(ctx echo.Context) error {
	id := ctx.Param("id")
	item, err := c.InventoryManager.GetItemByID(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Item not found"})
	}

	return ctx.JSON(http.StatusOK, responses.InventoryResponse{
		ID:          item.ID.Hex(),
		Name:        item.Name,
		Price:       item.Price,
		Currency:    item.Currency,
		Discount:    item.Discount,
		Vendor:      item.Vendor,
		Accessories: item.Accessories,
	})
}

func (c *InventoryController) UpdateItemHandler(ctx echo.Context) error {
	id := ctx.Param("id")
	var req requests.InventoryRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	item := &models.Inventory{
		Name:        req.Name,
		Price:       req.Price,
		Currency:    req.Currency,
		Discount:    req.Discount,
		Vendor:      req.Vendor,
		Accessories: req.Accessories,
	}

	updatedItem, err := c.InventoryManager.UpdateItem(ctx.Request().Context(), id, item)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update item"})
	}

	return ctx.JSON(http.StatusOK, responses.InventoryResponse{
		ID:          updatedItem.ID.Hex(),
		Name:        updatedItem.Name,
		Price:       updatedItem.Price,
		Currency:    updatedItem.Currency,
		Discount:    updatedItem.Discount,
		Vendor:      updatedItem.Vendor,
		Accessories: updatedItem.Accessories,
	})
}

func (c *InventoryController) DeleteItemHandler(ctx echo.Context) error {
	id := ctx.Param("id")
	err := c.InventoryManager.DeleteItem(ctx.Request().Context(), id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete item"})
	}

	return ctx.JSON(http.StatusOK, map[string]string{"message": "Item deleted successfully"})
}
