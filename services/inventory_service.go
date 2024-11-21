package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"main/config"
	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateItem(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	log.Println("mongo create service--------->")

	// Ensure ID is generated before inserting into MongoDB
	item.GenerateUUID()

	result, err := config.InventoryCollection.InsertOne(ctx, item)
	if err != nil {
		log.Printf("Error inserting inventory item: %v", err)
		return nil, err
	}

	// MongoDB will assign an ObjectID by default to `_id` field.
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		item.ID = objectID.Hex() // Store the MongoDB ObjectID as string if needed
	}

	return item, nil
}
func GetItems(ctx context.Context, pageNumber, pageSize int, vendors []string) ([]*models.Inventory, int64, error) {
	if pageSize == -1 {
		pageNumber = 1
		pageSize = 100
	}

	skip := (pageNumber - 1) * pageSize
	if skip < 0 {
		skip = 0
	}

	var items []*models.Inventory
	filter := bson.D{}

	// Filter by vendor if vendors are provided
	if len(vendors) > 0 {
		filter = bson.D{{Key: "vendor", Value: bson.D{{Key: "$in", Value: vendors}}}}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))      // Set pagination skip
	findOptions.SetLimit(int64(pageSize)) // Set pagination limit

	cursor, err := config.InventoryCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("Error fetching inventory items from MongoDB: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	// Deserialize cursor data into `items`
	if err := cursor.All(ctx, &items); err != nil {
		log.Printf("Error reading cursor data: %v", err)
		return nil, 0, err
	}

	// Return empty result with error if no items are found
	if len(items) == 0 && pageNumber > 1 {
		return nil, 0, errors.New("no records available for the requested page")
	}

	// Count the total number of inventory items in the collection
	totalCount, err := config.InventoryCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Printf("Error counting inventory documents in MongoDB: %v", err)
		return nil, 0, err
	}

	// Convert each item's ID to a string if needed (MongoDB returns ObjectID)
	for _, item := range items {
		item.ID = item.ID // Ensure it's a string (should already be in string format)
	}

	return items, totalCount, nil
}

func GetItemByID(ctx context.Context, id string) (*models.Inventory, error) {
	// Query MongoDB directly with the UUID (string) as _id
	var item models.Inventory
	err := config.InventoryCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error fetching inventory item by ID: %v", err)
		return nil, err
	}

	// Set the ID as string (it should already be a string as stored in MongoDB)
	return &item, nil
}

func UpdateItem(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
    // Directly use the ID as a string (no need for ObjectID conversion for UUIDs)
    update := bson.M{
        "$set": bson.M{
            "product_name": item.Name,
            "price":        item.Price,
            "currency":     item.Currency,
            "discount":     item.Discount,
            "vendor":       item.Vendor,
            "accessories":  item.Accessories,
        },
    }

    // Update the item in MongoDB using UUID as the string ID
    result := config.InventoryCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update)
    if result.Err() != nil {
        if result.Err().Error() == "mongo: no documents in result" {
            log.Printf("Item with ID %s not found", id)
            return nil, errors.New("inventory item not found")
        }
        log.Printf("Error updating inventory item: %v", result.Err())
        return nil, fmt.Errorf("failed to update item: %v", result.Err())
    }

    // Return the updated item with the same ID
    item.ID = id
    return item, nil
}


func DeleteItem(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return err
	}

	result, err := config.InventoryCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		log.Printf("Error deleting inventory item: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("inventory item not found")
	}

	return nil
}
