package handlers

import (
	"log"
	"strings"
	"time"

	"erp-project/models"
	"erp-project/repositories"
	"erp-project/utils"

	"github.com/gin-gonic/gin"
)

type SupplierHandler struct {
	repo *repositories.SupplierRepository
}

func NewSupplierHandler(repo *repositories.SupplierRepository) *SupplierHandler {
	return &SupplierHandler{repo: repo}
}

// Request structs
type CreateSupplierRequest struct {
	Name          string `json:"name" binding:"required,min=2,max=255"`
	Code          string `json:"code" binding:"required,min=2,max=50"`
	ContactPerson string `json:"contact_person" binding:"max=255"`
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone" binding:"omitempty,max=50"`
	Address       string `json:"address" binding:"max=500"`
	TaxID         string `json:"tax_id" binding:"max=100"`
	PaymentTerms  string `json:"payment_terms" binding:"max=500"`
}

type UpdateSupplierRequest struct {
	Name          string `json:"name" binding:"omitempty,min=2,max=255"`
	Code          string `json:"code" binding:"omitempty,min=2,max=50"`
	ContactPerson string `json:"contact_person" binding:"omitempty,max=255"`
	Email         string `json:"email" binding:"omitempty,email"`
	Phone         string `json:"phone" binding:"omitempty,max=50"`
	Address       string `json:"address" binding:"omitempty,max=500"`
	TaxID         string `json:"tax_id" binding:"omitempty,max=100"`
	PaymentTerms  string `json:"payment_terms" binding:"omitempty,max=500"`
	Status        string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type AddProductSupplierRequest struct {
	ProductID    string  `json:"product_id" binding:"required"`
	SupplierSKU  string  `json:"supplier_sku" binding:"max=100"`
	CostPrice    float64 `json:"cost_price" binding:"required,gt=0"`
	LeadTimeDays int     `json:"lead_time_days" binding:"omitempty,min=0"`
	IsPrimary    bool    `json:"is_primary"`
}

// Supplier CRUD handlers
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	supplier := models.NewSupplier(
		req.Name,
		req.Code,
		req.ContactPerson,
		req.Email,
		req.Phone,
		req.Address,
		req.TaxID,
		req.PaymentTerms,
	)

	if err := h.repo.CreateSupplier(supplier); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.DuplicateErrorResponse(c, "Duplicate supplier code", "A supplier with this code already exists")
			return
		}

		log.Printf("CreateSupplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create supplier", "Database error")
		return
	}

	utils.CreatedResponse(c, "Supplier created successfully", supplier)
}

func (h *SupplierHandler) GetAllSuppliers(c *gin.Context) {
	suppliers, err := h.repo.GetAllSuppliers()
	if err != nil {
		log.Printf("GetAllSuppliers error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve suppliers", "Database error")
		return
	}

	utils.SuccessResponse(c, "Suppliers retrieved successfully", suppliers)
}

func (h *SupplierHandler) GetSupplierByID(c *gin.Context) {
	id := c.Param("id")
	supplier, err := h.repo.GetSupplierByID(id)
	if err != nil {
		log.Printf("GetSupplierByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve supplier", "Database error")
		return
	}
	if supplier == nil {
		utils.NotFoundResponse(c, "Supplier not found")
		return
	}
	utils.SuccessResponse(c, "Supplier retrieved successfully", supplier)
}

func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Get existing supplier
	supplier, err := h.repo.GetSupplierByID(id)
	if err != nil {
		log.Printf("UpdateSupplier - GetSupplierByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve supplier", "Database error")
		return
	}

	if supplier == nil {
		utils.NotFoundResponse(c, "Supplier not found")
		return
	}

	// Update fields only if provided
	updatedFields := []string{}
	if req.Name != "" && req.Name != supplier.Name {
		supplier.Name = req.Name
		updatedFields = append(updatedFields, "name")
	}
	if req.Code != "" && req.Code != supplier.Code {
		supplier.Code = req.Code
		updatedFields = append(updatedFields, "code")
	}
	if req.ContactPerson != "" && req.ContactPerson != supplier.ContactPerson {
		supplier.ContactPerson = req.ContactPerson
		updatedFields = append(updatedFields, "contact_person")
	}
	if req.Email != "" && req.Email != supplier.Email {
		supplier.Email = req.Email
		updatedFields = append(updatedFields, "email")
	}
	if req.Phone != "" && req.Phone != supplier.Phone {
		supplier.Phone = req.Phone
		updatedFields = append(updatedFields, "phone")
	}
	if req.Address != "" && req.Address != supplier.Address {
		supplier.Address = req.Address
		updatedFields = append(updatedFields, "address")
	}
	if req.TaxID != "" && req.TaxID != supplier.TaxID {
		supplier.TaxID = req.TaxID
		updatedFields = append(updatedFields, "tax_id")
	}
	if req.PaymentTerms != "" && req.PaymentTerms != supplier.PaymentTerms {
		supplier.PaymentTerms = req.PaymentTerms
		updatedFields = append(updatedFields, "payment_terms")
	}
	if req.Status != "" && req.Status != supplier.Status {
		supplier.Status = req.Status
		updatedFields = append(updatedFields, "status")
	}

	// If no fields were updated
	if len(updatedFields) == 0 {
		utils.SuccessResponse(c, "No changes detected", supplier)
		return
	}

	supplier.UpdatedAt = time.Now()

	if err := h.repo.UpdateSupplier(supplier); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate supplier code", "A supplier with this code already exists")
			return
		}

		log.Printf("UpdateSupplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to update supplier", "Database error")
		return
	}

	// Get updated supplier
	updatedSupplier, err := h.repo.GetSupplierByID(id)
	if err != nil {
		log.Printf("UpdateSupplier - Get updated supplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve updated supplier", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"supplier":       updatedSupplier,
		"updated_fields": updatedFields,
	}

	utils.SuccessResponse(c, "Supplier updated successfully", responseData)
}

func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	supplier, err := h.repo.GetSupplierByID(id)
	if err != nil {
		log.Printf("DeleteSupplier - GetSupplierByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve supplier", "Database error")
		return
	}
	if supplier == nil {
		utils.NotFoundResponse(c, "Supplier not found")
		return
	}

	if err := h.repo.DeleteSupplier(id); err != nil {
		log.Printf("DeleteSupplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to delete supplier", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"deleted_supplier_id":   id,
		"deleted_supplier_name": supplier.Name,
	}

	utils.SuccessResponse(c, "Supplier deleted successfully", responseData)
}

// Product-Supplier relationship handlers
func (h *SupplierHandler) AddProductSupplier(c *gin.Context) {
	supplierID := c.Param("id")

	var req AddProductSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Verify supplier exists
	supplier, err := h.repo.GetSupplierByID(supplierID)
	if err != nil || supplier == nil {
		utils.NotFoundResponse(c, "Supplier not found")
		return
	}

	productSupplier := models.NewProductSupplier(
		req.ProductID,
		supplierID,
		req.SupplierSKU,
		req.CostPrice,
		req.LeadTimeDays,
		req.IsPrimary,
	)

	if err := h.repo.AddProductSupplier(productSupplier); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate product-supplier", "This product is already linked to this supplier")
			return
		}

		log.Printf("AddProductSupplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to link product to supplier", "Database error")
		return
	}

	utils.CreatedResponse(c, "Product linked to supplier successfully", productSupplier)
}

func (h *SupplierHandler) GetSupplierProducts(c *gin.Context) {
	supplierID := c.Param("id")

	// Verify supplier exists
	supplier, err := h.repo.GetSupplierByID(supplierID)
	if err != nil || supplier == nil {
		utils.NotFoundResponse(c, "Supplier not found")
		return
	}

	products, err := h.repo.GetSupplierProducts(supplierID)
	if err != nil {
		log.Printf("GetSupplierProducts error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve supplier products", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"supplier": supplier,
		"products": products,
		"count":    len(products),
	}

	utils.SuccessResponse(c, "Supplier products retrieved successfully", responseData)
}

func (h *SupplierHandler) RemoveProductSupplier(c *gin.Context) {
	productSupplierID := c.Param("product_supplier_id")

	if err := h.repo.RemoveProductSupplier(productSupplierID); err != nil {
		log.Printf("RemoveProductSupplier error: %v", err)
		utils.InternalErrorResponse(c, "Failed to remove product from supplier", "Database error")
		return
	}

	utils.SuccessResponse(c, "Product removed from supplier successfully", nil)
}
