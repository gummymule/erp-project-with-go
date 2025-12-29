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

type CustomerHandler struct {
	repo *repositories.CustomerRepository
}

func NewCustomerHandler(repo *repositories.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

type CreateCustomerRequest struct {
	Name    string `json:"name" binding:"required,min=2,max=100"`
	Email   string `json:"email" binding:"required,email"`
	Phone   string `json:"phone" binding:"omitempty,min=10,max=20"`
	Address string `json:"address" binding:"max=200"`
}

type UpdateCustomerRequest struct {
	Name    string `json:"name" binding:"omitempty,min=2,max=100"`
	Email   string `json:"email" binding:"omitempty,email"`
	Phone   string `json:"phone" binding:"omitempty,min=10,max=20"`
	Address string `json:"address" binding:"omitempty,max=200"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	customer := models.NewCustomer(
		req.Name,
		req.Email,
		req.Phone,
		req.Address,
	)

	if err := h.repo.CreateCustomer(customer); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate email", "A customer with this email already exists")
			return
		}

		log.Printf("CreateCustomer error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create customer", "Database error")
		return
	}

	utils.CreatedResponse(c, "Customer created successfully", customer)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	search := c.Query("search")
	email := c.Query("email")

	customers, total, err := h.repo.GetCustomerWithPagination(page, pageSize, search, email)
	if err != nil {
		log.Printf("GetAllCustomers error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve customers", "Database error")
		return
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	paginationData := map[string]interface{}{
		"customers": customers,
		"pagination": utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    totalPages,
		},
	}

	utils.SuccessResponse(c, "Customers retrieved successfully", paginationData)
}

func (h *CustomerHandler) GetListCustomers(c *gin.Context) {
	customers, err := h.repo.GetAllCustomers()
	if err != nil {
		log.Printf("GetListCustomers error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve customers", "Database error")
		return
	}

	utils.SuccessResponse(c, "Customers list retrieved successfully", customers)
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("GetCustomerByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve customer", "Database error")
		return
	}
	if customer == nil {
		utils.NotFoundResponse(c, "Customer not found")
		return
	}
	utils.SuccessResponse(c, "Customer retrieved successfully", customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// get existing customer
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("UpdateCustomer - GetCustomerByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve customer", "Database error")
		return
	}

	if customer == nil {
		utils.NotFoundResponse(c, "Customer not found")
		return
	}

	// update fields only if they are provided in request
	updatedFields := []string{}
	if req.Name != "" && req.Name != customer.Name {
		customer.Name = req.Name
		updatedFields = append(updatedFields, "name")
	}
	if req.Email != "" && req.Email != customer.Email {
		customer.Email = req.Email
		updatedFields = append(updatedFields, "email")
	}
	if req.Phone != "" && req.Phone != customer.Phone {
		customer.Phone = req.Phone
		updatedFields = append(updatedFields, "phone")
	}
	if req.Address != "" && req.Address != customer.Address {
		customer.Address = req.Address
		updatedFields = append(updatedFields, "address")
	}

	// If no fields were updated
	if len(updatedFields) == 0 {
		utils.SuccessResponse(c, "No changes detected", customer)
		return
	}

	// Update timestamp
	customer.UpdatedAt = time.Now()

	if err := h.repo.UpdateCustomer(customer); err != nil {
		// Check for duplicate email error
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.DuplicateErrorResponse(c, "Duplicate email", "A customer with this email already exists")
			return
		}

		log.Printf("UpdateCustomer - UpdateCustomer error: %v", err)
		utils.InternalErrorResponse(c, "Failed to update customer", "Database error")
		return
	}

	// Get updated customer to return fresh data
	updatedCustomer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("UpdateCustomer - Get updated customer error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve updated customer", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"customer":       updatedCustomer,
		"updated_fields": updatedFields,
	}

	utils.SuccessResponse(c, "Customer updated successfully", responseData)
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		log.Printf("DeleteCustomer - GetCustomerByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve customer", "Database error")
		return
	}
	if customer == nil {
		utils.NotFoundResponse(c, "Customer not found")
		return
	}

	if err := h.repo.DeleteCustomer(id); err != nil {
		log.Printf("DeleteCustomer - DeleteCustomer error: %v", err)
		utils.InternalErrorResponse(c, "Failed to delete customer", "Database error")
		return
	}

	utils.SuccessResponse(c, "Customer deleted successfully", nil)
}
