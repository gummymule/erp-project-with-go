package models

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Code          string    `json:"code"`
	ContactPerson string    `json:"contact_person"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	Address       string    `json:"address"`
	TaxID         string    `json:"tax_id"`
	PaymentTerms  string    `json:"payment_terms"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ProductSupplier struct {
	ID           string    `json:"id"`
	ProductID    string    `json:"product_id"`
	SupplierID   string    `json:"supplier_id"`
	SupplierSKU  string    `json:"supplier_sku"`
	CostPrice    float64   `json:"cost_price"`
	LeadTimeDays int       `json:"lead_time_days"`
	IsPrimary    bool      `json:"is_primary"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	ProductName  *Product  `json:"product_name,omitempty"`
	SupplierName *Supplier `json:"supplier_name,omitempty"`
}

func NewSupplier(name, code, contactPerson, email, phone, address, taxID, paymentTerms string) *Supplier {
	now := time.Now()

	return &Supplier{
		ID:            uuid.New().String(),
		Name:          name,
		Code:          code,
		ContactPerson: contactPerson,
		Email:         email,
		Phone:         phone,
		Address:       address,
		TaxID:         taxID,
		PaymentTerms:  paymentTerms,
		Status:        "active",
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func NewProductSupplier(productID, supplierID, supplierSKU string, costPrice float64, leadTimeDays int, isPrimary bool) *ProductSupplier {
	now := time.Now()

	return &ProductSupplier{
		ID:           uuid.New().String(),
		ProductID:    productID,
		SupplierID:   supplierID,
		SupplierSKU:  supplierSKU,
		CostPrice:    costPrice,
		LeadTimeDays: leadTimeDays,
		IsPrimary:    isPrimary,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
