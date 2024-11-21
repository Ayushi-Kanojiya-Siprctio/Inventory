package managers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/models"
	service "main/services"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryManager struct{}

func (m *InventoryManager) GetItems(ctx context.Context, flag string, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	switch flag {
	case "0":
		items, totalCount, err := service.GetItems(ctx, pageNumber, pageSize, vendors)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil

	case "1":
		items, totalCount, err := service.GetItemsPostgres(ctx, pageNumber, pageSize, vendors)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil

	default:
		return nil, 0, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) CreateItem(ctx context.Context, flag string, item *models.Inventory) (*models.Inventory, error) {
	log.Println("flag--------->", flag)

	switch flag {
	case "0":
		return service.CreateItem(ctx, item)

	case "1":
		return service.CreateItemPostgres(ctx, item)

	default:
		return nil, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) GetItemByID(ctx context.Context, flag string, id string) (*models.Inventory, error) {
	switch flag {
	case "0":
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		return service.GetItemByID(ctx, objectID.Hex())

	case "1":
		return service.GetItemByIDPostgres(ctx, id)

	default:
		return nil, errors.New("invalid flag type")
	}
}
func (m *InventoryManager) UpdateItem(ctx context.Context, flag string, id string, item *models.Inventory) (*models.Inventory, error) {
	switch flag {
	case "0":
		updatedItem, err := service.UpdateItem(ctx, id, item)
		if err != nil {
			log.Printf("Error updating item in service: %v", err)
			return nil, fmt.Errorf("failed to update item: %v", err)
		}
		return updatedItem, nil

	case "1":
		return service.UpdateItemPostgres(ctx, id, item)

	default:
		return nil, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) DeleteItem(ctx context.Context, flag string, id string) error {
	switch flag {
	case "0":

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		return service.DeleteItem(ctx, objectID.Hex())

	case "1":
		return service.DeleteItemPostgres(ctx, id)

	default:
		return errors.New("invalid flag type")
	}
}
