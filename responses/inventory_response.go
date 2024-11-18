package responses

// type Response struct {
// 	Status  int         `json:"status"`
// 	Message string      `json:"message"`
// 	Data    interface{} `json:"data,omitempty"`
// }

// func NewResponse(message string, data interface{}) *Response {
// 	return &Response{
// 		Status:  http.StatusOK,
// 		Message: message,
// 		Data:    data,
// 	}
// }

// func ErrorResponse(message string) *Response {
// 	return &Response{
// 		Status:  http.StatusBadRequest,
// 		Message: message,
// 	}
// }

type InventoryResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"product_name"`
	Price       int      `json:"price"`
	Currency    string   `json:"currency"`
	Discount    int      `json:"discount,omitempty"`
	Vendor      string   `json:"vendor"`
	Accessories []string `json:"accessories,omitempty"`
}
