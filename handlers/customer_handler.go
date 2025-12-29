package handlers

import (
	"log"
	"strings"

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
