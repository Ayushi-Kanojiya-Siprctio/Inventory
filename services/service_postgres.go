// package service

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"main/config"
// 	"main/models"

// 	"gorm.io/gorm"
// )

// func CreateTableIfNotExists(db *gorm.DB) error {
//     // SQL query to create table if not exists
//     query := `
// 	  DROP TABLE IF EXISTS "inventories";
//     CREATE TABLE "inventories" (
//         "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
//         "name" varchar(255),
//         "price" bigint,
//         "currency" varchar(10),
//         "discount" bigint,
//         "vendor" varchar(255),
//         "accessories" jsonb
//     );
//     `
//     // Execute the query
//     err := db.Exec(query).Error
//     if err != nil {
//         log.Printf("Error executing table creation query: %v", err)
//         return fmt.Errorf("failed to create table: %v", err)
//     }

//     log.Println("Table 'inventories' checked/created successfully.")
//     return nil
// }

// func CreateItemPostgres(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
// 	log.Println("postgres create service--------->")
// 	log.Println("PG-------->", config.PG)

// 	if config.PG == nil {
// 		log.Println("Error: PostgreSQL database connection is not initialized.")
// 		return nil, errors.New("PostgreSQL database connection is not initialized")
// 	}

// 	err := config.PG.AutoMigrate(&models.Inventory{}) // Create the table if it doesn't exist
// 	if err != nil {
// 		log.Println("Error during AutoMigrate:", err)
// 		return &models.Inventory{}, fmt.Errorf("error ensuring table exists: %w", err)
// 	}

// 	// Now insert the student record
// 	err = config.PG.Create(&item).Error
// 	if err != nil {
// 		log.Println("Error inserting item:", err)
// 		return &models.Inventory{}, fmt.Errorf("error inserting item: %w", err)
// 	}

// 	return item, nil
// }

// func GetItemsPostgres(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
// 	var items []*models.Inventory
// 	var totalCount int64

// 	if config.PG == nil {
// 		log.Println("Error: PostgreSQL database connection is not initialized.")
// 		return nil, 0, errors.New("PostgreSQL database connection is not initialized")
// 	}

// 	query := config.PG.WithContext(ctx).Model(&models.Inventory{})
// 	if len(vendors) > 0 {
// 		query = query.Where("vendor IN ?", vendors)
// 	}

// 	if err := query.Count(&totalCount).Error; err != nil {
// 		log.Printf("Error counting inventory items in PostgreSQL: %v", err)
// 		return nil, 0, err
// 	}

// 	offset := (pageNumber - 1) * pageSize
// 	if err := query.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
// 		log.Printf("Error fetching inventory items from PostgreSQL: %v", err)
// 		return nil, 0, err
// 	}

// 	return items, totalCount, nil
// }

// func GetItemByIDPostgres(ctx context.Context, id string) (*models.Inventory, error) {
// 	var item models.Inventory

// 	if config.PG == nil {
// 		log.Println("Error: PostgreSQL database connection is not initialized.")
// 		return nil, errors.New("PostgreSQL database connection is not initialized")
// 	}

// 	if err := config.PG.WithContext(ctx).First(&item, "id = ?", id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("inventory item not found")
// 		}
// 		log.Printf("Error fetching inventory item by ID from PostgreSQL: %v", err)
// 		return nil, err
// 	}
// 	return &item, nil
// }

// func UpdateItemPostgres(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
// 	existingItem := &models.Inventory{}

// 	if config.PG == nil {
// 		log.Println("Error: PostgreSQL database connection is not initialized.")
// 		return nil, errors.New("PostgreSQL database connection is not initialized")
// 	}

// 	if err := config.PG.WithContext(ctx).First(existingItem, "id = ?", id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("inventory item not found")
// 		}
// 		log.Printf("Error finding inventory item for update in PostgreSQL: %v", err)
// 		return nil, err
// 	}

// 	if err := config.PG.WithContext(ctx).Model(existingItem).Updates(item).Error; err != nil {
// 		log.Printf("Error updating inventory item in PostgreSQL: %v", err)
// 		return nil, err
// 	}

// 	item.ID = existingItem.ID
// 	return item, nil
// }

// func DeleteItemPostgres(ctx context.Context, id string) error {
// 	if config.PG == nil {
// 		log.Println("Error: PostgreSQL database connection is not initialized.")
// 		return errors.New("PostgreSQL database connection is not initialized")
// 	}

// 	if err := config.PG.WithContext(ctx).Delete(&models.Inventory{}, "id = ?", id).Error; err != nil {
// 		log.Printf("Error deleting inventory item from PostgreSQL: %v", err)
// 		return err
// 	}
// 	return nil
// }



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
	// SQL query to drop and create table
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

	// Execute the query
	err := db.Exec(query).Error
	if err != nil {
		log.Printf("Error executing table creation query: %v", err)
		return fmt.Errorf("failed to create table: %v", err)
	}

	log.Println("Table 'inventories' checked/created successfully.")
	return nil
}

// Insert an inventory item into the database
func CreateItemPostgres(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	log.Println("postgres create service--------->")
	log.Println("PG-------->", config.PG)

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	// Insert query
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

// Get all inventory items with pagination and filtering by vendor
func GetItemsPostgres(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	var items []*models.Inventory
	var totalCount int64

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, 0, errors.New("PostgreSQL database connection is not initialized")
	}

	// Get All items with pagination and optional vendor filtering
	query := config.PG.WithContext(ctx).Model(&models.Inventory{})
	if len(vendors) > 0 {
		query = query.Where("vendor IN ?", vendors)
	}

	// Count total items for pagination
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

// Get an inventory item by ID
func GetItemByIDPostgres(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory

	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	// Get Item By ID query
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

// Update an inventory item by ID
func UpdateItemPostgres(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return nil, errors.New("PostgreSQL database connection is not initialized")
	}

	// Update Item by ID query
	query := `UPDATE inventories SET product_name = ?, price = ?, currency = ?, discount = ?, vendor = ?, accessories = ? WHERE id = ?`
	err := config.PG.Exec(query, item.Name, item.Price, item.Currency, item.Discount, item.Vendor, item.Accessories, id).Error
	if err != nil {
		log.Printf("Error updating inventory item in PostgreSQL: %v", err)
		return nil, fmt.Errorf("error updating item: %w", err)
	}

	// Retrieve the updated item
	updatedItem, err := GetItemByIDPostgres(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching updated item: %w", err)
	}

	log.Println("Item updated successfully:", updatedItem)
	return updatedItem, nil
}

// Delete an inventory item by ID
func DeleteItemPostgres(ctx context.Context, id string) error {
	if config.PG == nil {
		log.Println("Error: PostgreSQL database connection is not initialized.")
		return errors.New("PostgreSQL database connection is not initialized")
	}

	// Delete Item by ID query
	query := `DELETE FROM inventories WHERE id = ?`
	err := config.PG.Exec(query, id).Error
	if err != nil {
		log.Printf("Error deleting inventory item from PostgreSQL: %v", err)
		return fmt.Errorf("error deleting item: %w", err)
	}

	log.Println("Item deleted successfully with ID:", id)
	return nil
}
