package responses

import "go.mongodb.org/mongo-driver/bson/primitive"

type InventoryResponse struct {
	ID          primitive.ObjectID `json:"id"`
	Name        string             `json:"product_name"`
	Price       int                `json:"price"`
	Currency    string             `json:"currency"`
	Discount    int                `json:"discount,omitempty"`
	Vendor      string             `json:"vendor"`
	Accessories []string           `json:"accessories,omitempty"`
}
