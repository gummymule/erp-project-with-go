package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

func NewCustomer(name, email, phone, address string) *Customer {
	return &Customer{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
