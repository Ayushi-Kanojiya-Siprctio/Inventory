package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" gorm:"primaryKey"`
	Name        string             `bson:"product_name" gorm:"size:255"`
	Price       int                `bson:"price"`
	Currency    string             `bson:"currency" gorm:"size:10"`
	Discount    int                `bson:"discount"`
	Vendor      string             `bson:"vendor" gorm:"size:255"`
	Accessories []string           `bson:"accessories" gorm:"-"`
}
