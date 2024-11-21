package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/models"

	"gorm.io/gorm"
)

func CreateTableIfNotExists(db *gorm.DB) error {
	query := `
	DROP TABLE IF EXISTS "inventories";
	CREATE TABLE "inventories" (
		"id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		"product_name" varchar(255),
		"price" bigint,
		"currency" varchar(10),
		"discount" bigint,
		"vendor" varchar(255),
		"accessories" jsonb
	);
	`

	err := db.Exec(query).Error
	if err != nil {
		log.Printf("Error executing table creation query: %v", err)
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Table 'inventories' checked/created successfully.")
	return nil
}

func CreateItemPostgres(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	log.Println("postgres create service--------->")
	log.Println("PG-------->", config.PG)

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	query := `INSERT INTO inventories (product_name, price, currency, discount, vendor, accessories)
				VALUES (?, ?, ?, ?, ?, ?) RETURNING id`
	err := config.PG.Raw(query, item.Name, item.Price, item.Currency, item.Discount, item.Vendor, item.Accessories).Scan(&item.ID).Error
	if err != nil {
		log.Println("Error inserting item:", err)
		return nil, fmt.Errorf("error inserting item: %w", err)
	}

	log.Println("Item created successfully:", item)
	return item, nil
}

func GetItemsPostgres(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	var items []*models.Inventory
	var totalCount int64

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, 0, errors.New("PostgreSQL database connection is not initialized")
	}

	query := config.PG.WithContext(ctx).Model(&models.Inventory{})
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

	log.Println("Items fetched successfully:", items)
	return items, totalCount, nil
}

func GetItemByIDPostgres(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	query := `SELECT * FROM inventories WHERE id = ?`
	err := config.PG.Raw(query, id).Scan(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error fetching inventory item by ID from PostgreSQL: %v", err)
		return nil, err
	}

	log.Println("Item fetched by ID:", item)
	return &item, nil
}

func UpdateItemPostgres(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	query := `UPDATE inventories SET product_name = ?, price = ?, currency = ?, discount = ?, vendor = ?, accessories = ? WHERE id = ?`
	err := config.PG.Exec(query, item.Name, item.Price, item.Currency, item.Discount, item.Vendor, item.Accessories, id).Error
	if err != nil {
		log.Printf("Error updating inventory item in PostgreSQL: %v", err)
		return nil, fmt.Errorf("error updating item: %w", err)
	}

	updatedItem, err := GetItemByIDPostgres(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated item: %w", err)
	}

	log.Println("Item updated successfully:", updatedItem)
	return updatedItem, nil
}

func DeleteItemPostgres(ctx context.Context, id string) error {
	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return errors.New("PostgreSQL database connection is not initialized")
	}

	query := `DELETE FROM inventories WHERE id = ?`
	err := config.PG.Exec(query, id).Error
	if err != nil {
		log.Printf("Error deleting inventory item from PostgreSQL: %v", err)
		return fmt.Errorf("error deleting item: %w", err)
	}

	log.Println("Item deleted successfully with ID:", id)
	return nil
}
