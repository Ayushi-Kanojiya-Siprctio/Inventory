package service

import (
	"context"
	"errors"
	"log"
	"main/config"
	"main/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateItem(ctx context.Context, item *models.Inventory) (*models.Inventory, error) {
	log.Println("mongo create service--------->")

	// item.SetMongoDB()
	// item.GenerateUUID()

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

func GetItems(ctx context.Context) ([]*models.Inventory, int64, error) {
	var items []*models.Inventory
	var totalCount int64

	cursor, err := config.InventoryCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}

	totalCount, err = config.InventoryCollection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return items, totalCount, nil
}

func GetItemByID(ctx context.Context, id string) (*models.Inventory, error) {
	var item models.Inventory

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	err = config.InventoryCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error fetching inventory item by ID: %v", err)
		return nil, err
	}

	return &item, nil
}

func UpdateItem(ctx context.Context, id primitive.ObjectID, item *models.Inventory) (*models.Inventory, error) {
	update := bson.M{"$set": item}

	result := config.InventoryCollection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update)

	if result.Err() != nil {
		if result.Err().Error() == "mongo: no documents in result" {
			return nil, errors.New("inventory item not found")
		}
		log.Printf("Error updating inventory item: %v", result.Err())
		return nil, result.Err()
	}

	item.ID = id.Hex()
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
