package responses


type InventoryResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"product_name"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	Discount    int      `json:"discount"`
	Vendor      string   `json:"vendor"`
}
