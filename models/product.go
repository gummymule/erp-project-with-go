package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Category    string    `json:"category"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProduct(name, description, sku, category string, price float64, quantity int) *Product {
	return &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		SKU:         sku,
		Price:       price,
		Quantity:    quantity,
		Category:    category,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
