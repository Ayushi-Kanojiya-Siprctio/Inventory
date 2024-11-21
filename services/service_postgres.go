package service

import (
	"context"
	"errors"
	"log"
	"main/models"

	"gorm.io/gorm"
)

var DB *gorm.DB 
func CreateItemPostgres(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	if err := DB.WithContext(ctx).Create(&item).Error; err != nil {
		log.Printf("Error inserting inventory item into PostgreSQL: %v", err)
		return nil, err
	}
	return item, nil
}

func GetItemsPostgres(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	var items []*models.Inventory
	var totalCount int64

	query := DB.WithContext(ctx).Model(&models.Inventory{})
	if len(vendors) > 0 {
		query = query.Where("vendor IN ?", vendors)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting inventory items in PostgreSQL: %v", err)
		return nil, 0, err
	}

	offset := (pageNumber - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		log.Printf("Error fetching inventory items from PostgreSQL: %v", err)
		return nil, 0, err
	}

	return items, totalCount, nil
}

func GetItemByIDPostgres(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory
	if err := DB.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error fetching inventory item by ID from PostgreSQL: %v", err)
		return nil, err
	}
	return &item, nil
}

func UpdateItemPostgres(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
	existingItem := &models.Inventory{}
	if err := DB.WithContext(ctx).First(existingItem, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error finding inventory item for update in PostgreSQL: %v", err)
		return nil, err
	}

	if err := DB.WithContext(ctx).Model(existingItem).Updates(item).Error; err != nil {
		log.Printf("Error updating inventory item in PostgreSQL: %v", err)
		return nil, err
	}

	item.ID = existingItem.ID
	return item, nil
}

func DeleteItemPostgres(ctx context.Context, id string) error {
	if err := DB.WithContext(ctx).Delete(&models.Inventory{}, "id = ?", id).Error; err != nil {
		log.Printf("Error deleting inventory item from PostgreSQL: %v", err)
		return err
	}
	return nil
}
