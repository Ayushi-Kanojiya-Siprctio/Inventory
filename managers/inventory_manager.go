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

func (m *InventoryManager) GetItems(ctx context.Context, flag bool) ([]*models.Inventory, int64, error) {

	switch flag {
	case true:
		items, totalCount, err := service.GetItems(ctx)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil

	case false:
		items, totalCount, err := service.GetItemsPostgres(ctx)
		if err != nil {
			return nil, 0, err
		}
		return items, totalCount, nil

	default:
		return nil, 0, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) CreateItem(ctx context.Context, flag bool, item *models.Inventory) (*models.Inventory, error) {
	log.Println("flag--------->", flag)
	log.Println("manager item-------->", item.Name)

	switch flag {
	case true:
		return service.CreateItem(ctx, item)

	case false:
		return service.CreateItemPostgres(ctx, item)

	default:
		return nil, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) GetItemByID(ctx context.Context, flag bool, id string) (*models.Inventory, error) {
	switch flag {
	case true:
		return service.GetItemByID(ctx, id)
	case false:
		return service.GetItemByIDPostgres(ctx, id)
	default:
		return nil, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) UpdateItem(ctx context.Context, flag bool, id string, item *models.Inventory) (*models.Inventory, error) {
	log.Printf("Updating item with flag: %v, ID: %v", flag, id)

	switch flag {
	case true:
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Printf("Invalid ID format: %v", err)
			return nil, err
		}

		updatedItem, err := service.UpdateItem(ctx, objectID, item)
		if err != nil {
			log.Printf("Error updating item in service: %v", err)
			return nil, fmt.Errorf("failed to update item: %v", err)
		}
		return updatedItem, nil

	case false:
		return service.UpdateItemPostgres(ctx, id, item)

	default:
		return nil, errors.New("invalid flag type")
	}
}

func (m *InventoryManager) DeleteItem(ctx context.Context, flag bool, id string) error {
	switch flag {
	case true:

		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return err
		}

		return service.DeleteItem(ctx, objectID.Hex())

	case false:
		return service.DeleteItemPostgres(ctx, id)

	default:
		return errors.New("invalid flag type")
	}
}

