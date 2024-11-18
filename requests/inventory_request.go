package requests

type InventoryRequest struct {
	Name        string   `json:"product_name" bson:"product_name" validate:"required,max=20"`
	Price       int      `json:"price" bson:"price" validate:"required,max=2000"`
	Currency    string   `json:"currency" bson:"currency" validate:"required,len=3"`
	Discount    int      `json:"discount" bson:"discount"`
	Vendor      string   `json:"vendor" bson:"vendor" validate:"required"`
	Accessories []string `json:"accessories,omitempty" bson:"accessories,omitempty"`
}
