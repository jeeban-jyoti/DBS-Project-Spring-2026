package orders

import "database/sql"

type cartItem struct {
	bookID int
	qty    int
}

type OrderItemDetail struct {
	BookTitle    string `json:"book_title"`
	Quantity     int    `json:"quantity"`
	PurchaseType string `json:"purchase_type"`
}

type OrderDetail struct {
	OrderID        int               `json:"order_id"`
	StudentID      int               `json:"student_id"`
	StudentEmail   string            `json:"student_email,omitempty"` // Added for admin view
	CreatedDate    sql.NullString    `json:"created_date"`
	FulfilledDate  sql.NullString    `json:"fulfilled_date,omitempty"`
	ShippingType   sql.NullString    `json:"shipping_type,omitempty"`
	CardHolderName sql.NullString    `json:"card_holder_name,omitempty"`
	Status         string            `json:"status"`
	Items          []OrderItemDetail `json:"items"`
}
