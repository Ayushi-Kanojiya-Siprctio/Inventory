package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// type:uuid;default:gen_random_uuid();primaryKey"
type Inventory struct {
	ID       string `gorm:"column:id;bson:"_id" json:"id"`
	Name     string `gorm:"size:255;column:product_name" bson:"product_name" json:"product_name"`
	Price    int    `gorm:"column:price" bson:"price" json:"price"`
	Currency string `gorm:"size:10;column:currency" bson:"currency" json:"currency"`
	Discount int    `gorm:"column:discount" bson:"discount" json:"discount"`
	Vendor   string `gorm:"size:255;column:vendor" bson:"vendor" json:"vendor"`
}

func (i *Inventory) SetMongoDB() {
	if i.ID == "" {
		i.ID = primitive.NewObjectID().Hex()
	}
}

func (i *Inventory) GenerateUUID() {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
}
