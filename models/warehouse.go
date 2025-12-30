package models

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	ManagerName string    `json:"manager_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email"`
	Capacity    int       `json:"capacity"`
	Status      string    `json:"status"` // active, inactive
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type WarehouseLocation struct {
	ID              string    `json:"id"`
	WarehouseID     string    `json:"warehouse_id"`
	LocationCode    string    `json:"location_code"`
	LocationName    string    `json:"location_name"`
	Zone            string    `json:"zone"`
	RowNumber       int       `json:"row_number"`
	ShelfNumber     int       `json:"shelf_number"`
	MaxCapacity     int       `json:"max_capacity"`
	CurrentQuantity int       `json:"current_quantity"`
	Status          string    `json:"status"` // available, occupied, reserved, maintenance
	CreatedAt       time.Time `json:"created_at"`

	// For joins
	WarehouseName string `json:"warehouse_name,omitempty"`
}

type Inventory struct {
	ID                string     `json:"id"`
	ProductID         string     `json:"product_id"`
	WarehouseID       string     `json:"warehouse_id"`
	LocationID        *string    `json:"location_id,omitempty"`
	Quantity          int        `json:"quantity"`
	ReservedQuantity  int        `json:"reserved_quantity"`
	AvailableQuantity int        `json:"available_quantity"` // Calculated: quantity - reserved_quantity
	MinQuantity       int        `json:"min_quantity"`
	MaxQuantity       *int       `json:"max_quantity,omitempty"`
	LastRestocked     *time.Time `json:"last_restocked,omitempty"`
	LastChecked       *time.Time `json:"last_checked,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`

	// For joins
	ProductName   string `json:"product_name,omitempty"`
	WarehouseName string `json:"warehouse_name,omitempty"`
	LocationCode  string `json:"location_code,omitempty"`
	SKU           string `json:"sku,omitempty"`
}

func NewWarehouse(code, name, location, managerName, phone, email string, capacity int) *Warehouse {
	now := time.Now()
	return &Warehouse{
		ID:          uuid.New().String(),
		Code:        code,
		Name:        name,
		Location:    location,
		ManagerName: managerName,
		Phone:       phone,
		Email:       email,
		Capacity:    capacity,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func NewWarehouseLocation(warehouseID, locationCode, locationName, zone string, rowNumber, shelfNumber, maxCapacity int) *WarehouseLocation {
	now := time.Now()
	return &WarehouseLocation{
		ID:              uuid.New().String(),
		WarehouseID:     warehouseID,
		LocationCode:    locationCode,
		LocationName:    locationName,
		Zone:            zone,
		RowNumber:       rowNumber,
		ShelfNumber:     shelfNumber,
		MaxCapacity:     maxCapacity,
		CurrentQuantity: 0,
		Status:          "available",
		CreatedAt:       now,
	}
}

func NewInventory(productID, warehouseID string, locationID *string, quantity, minQuantity int) *Inventory {
	now := time.Now()
	availableQuantity := quantity

	inventory := &Inventory{
		ID:                uuid.New().String(),
		ProductID:         productID,
		WarehouseID:       warehouseID,
		Quantity:          quantity,
		ReservedQuantity:  0,
		AvailableQuantity: availableQuantity,
		MinQuantity:       minQuantity,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if locationID != nil {
		inventory.LocationID = locationID
	}

	return inventory
}
