package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	TotalAmount  float64   `json:"total_amount"`
	Status       string    `json:"status"`
	OrderDate    time.Time `json:"order_date"`
	CustomerName string    `json:"customer_name,omitempty"` // For joins
}

type OrderItem struct {
	ID          string  `json:"id"`
	OrderID     string  `json:"order_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	ProductName string  `json:"product_name,omitempty"` // For joins
}

func NewOrder(customerID string, totalAmount float64) *Order {
	return &Order{
		ID:          uuid.New().String(),
		CustomerID:  customerID,
		TotalAmount: totalAmount,
		Status:      "pending",
		OrderDate:   time.Now(),
	}
}

func NewOrderItem(orderID, productID string, quantity int, unitPrice float64) *OrderItem {
	return &OrderItem{
		ID:         uuid.New().String(),
		OrderID:    orderID,
		ProductID:  productID,
		Quantity:   quantity,
		UnitPrice:  unitPrice,
		TotalPrice: float64(quantity) * unitPrice,
	}
}
