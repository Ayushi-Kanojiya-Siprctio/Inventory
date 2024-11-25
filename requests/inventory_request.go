package requests

//gorm:"column:id;type:uuid;default:gen_random_uuid()"

type InventoryRequest struct {
	Name        string   `json:"product_name" binding:"required"`
	Price       int      `json:"price" binding:"required"`
	Currency    string   `json:"currency" binding:"required"`
	Discount    int      `json:"discount" binding:"required"`
	Vendor      string   `json:"vendor" binding:"required"`
}
