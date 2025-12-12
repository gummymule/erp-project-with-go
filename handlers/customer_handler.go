package handlers

import (
	"net/http"

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
	Phone   string `json:"phone" binding:"omitempty,phone"`
	Address string `json:"address" binding:"max=200"`
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := models.NewCustomer(
		req.Name,
		req.Email,
		req.Phone,
		req.Address,
	)

	if err := h.repo.CreateCustomer(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *CustomerHandler) GetAllCustomers(c *gin.Context) {
	page, pageSize := utils.GetPaginationParams(c)
	search := c.Query("search")
	email := c.Query("email")

	customers, total, err := h.repo.GetCustomerWithPagination(page, pageSize, search, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customers"})
		return
	}

	totalPages := utils.CalculateTotalPages(total, pageSize)

	response := utils.PaginatedResponse{
		Data: customers,
		Pagination: utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *CustomerHandler) GetCustomerByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.repo.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve customer"})
		return
	}
	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}
