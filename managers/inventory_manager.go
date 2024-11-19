package manager

import (
	"context"
	"main/models"
	service "main/services"

	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryManager struct{}

func (m *InventoryManager) GetItems(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	items, totalCount, err := service.GetItems(ctx, pageNumber, pageSize, vendors)
	if err != nil {
		log.Printf("Error fetching inventory items: %v", err)
		return nil, 0, err
	}
	return items, totalCount, nil
}

func (m *InventoryManager) CreateItem(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	createdItem, err := service.CreateItem(ctx, item)
	if err != nil {
		log.Printf("Error creating inventory item: %v", err)
		return nil, err
	}
	return createdItem, nil
}

func (m *InventoryManager) GetItemByID(ctx context.Context, id string) (*models.Inventory, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid inventory ID format: %v", err)
		return nil, err
	}

	item, err := service.GetItemByID(ctx, objectID)
	if err != nil {
		log.Printf("Error fetching inventory item by ID: %v", err)
		return nil, err
	}
	return item, nil
}

func (m *InventoryManager) UpdateItem(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
	updatedItem, err := service.UpdateItem(ctx, id, item)
	if err != nil {
		log.Printf("Error updating inventory item: %v", err)
		return nil, err
	}
	return updatedItem, nil
}

func (m *InventoryManager) DeleteItem(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid inventory ID format: %v", err)
		return err
	}

	err = service.DeleteItem(ctx, objectID)
	if err != nil {
		log.Printf("Error deleting inventory item: %v", err)
		return err
	}
	return nil
}