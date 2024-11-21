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

	item.SetMongoDB()

	result, err := config.InventoryCollection.InsertOne(ctx, item)
	if err != nil {
		log.Printf("Error inserting inventory item: %v", err)
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		item.ID = objectID.Hex()
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

	if len(vendors) > 0 {
		filter = bson.D{{Key: "vendor", Value: bson.D{{Key: "$in", Value: vendors}}}}
	}

	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	cursor, err := config.InventoryCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Printf("Error fetching inventory items from MongoDB: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		log.Printf("Error reading cursor data: %v", err)
		return nil, 0, err
	}

	if len(items) == 0 && pageNumber > 1 {
		return nil, 0, errors.New("no records available for the requested page")
	}

	totalCount, err := config.InventoryCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Printf("Error counting inventory documents in MongoDB: %v", err)
		return nil, 0, err
	}

	for _, item := range items {
		item.ID = item.ID
	}

	return items, totalCount, nil
}

func GetItemByID(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory
	err := config.InventoryCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error fetching inventory item by ID: %v", err)
		return nil, err
	}

	return &item, nil
}

func UpdateItem(ctx context.Context, id string, item *models.Inventory) (*models.Inventory, error) {
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

	result := config.InventoryCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update)
	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			log.Printf("Item with ID %s not found", id)
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error updating inventory item: %v", result.Err())
		return nil, fmt.Errorf("failed to update item: %v", result.Err())
	}

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
