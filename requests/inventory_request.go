package requests

//gorm:"column:id;type:uuid;default:gen_random_uuid()"

type InventoryRequest struct {
	Name        string   `json:"product_name" validate:"required" binding:"required"`
	Price       int      `json:"price" validate:"required" binding:"required"`
	Currency    string   `json:"currency" validate:"required" binding:"required"`
	Discount    int      `json:"discount" validate:"required" binding:"required"`
	Vendor      string   `json:"vendor" validate:"required" binding:"required"`
}
