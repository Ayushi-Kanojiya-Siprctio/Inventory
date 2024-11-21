package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ID          string   `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string   `gorm:"size:255" json:"product_name"`
	Price       int      `gorm:"column:price" json:"price"`
	Currency    string   `gorm:"size:10" json:"currency"`
	Discount    int      `gorm:"column:discount" json:"discount"`
	Vendor      string   `gorm:"size:255" json:"vendor"`
	Accessories []string `gorm:"type:text" json:"accessories"` // This stores an array of strings as JSON

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
