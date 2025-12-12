package handlers

import (
	"net/http"
	"strings"
	"time"

	"erp-project/models"
	"erp-project/repositories"
	"erp-project/utils"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	repo *repositories.ProductRepository
}

func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"max=500"`
	SKU         string  `json:"sku" binding:"required,sku"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int     `json:"quantity" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"max=50"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=2,max=100"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	SKU         string  `json:"sku" binding:"omitempty,sku"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity    int     `json:"quantity" binding:"omitempty,gte=0"`
	Category    string  `json:"category" binding:"omitempty,max=50"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, utils.NewAppError(
			http.StatusBadRequest,
			"Validation Error",
			err.Error(),
		))
		return
	}

	product := models.NewProduct(
		req.Name,
		req.Description,
		req.SKU,
		req.Category,
		req.Price,
		req.Quantity,
	)

	if err := h.repo.CreateProduct(product); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.ErrorResponse(c, utils.NewAppError(
				http.StatusBadRequest,
				"Duplicated SKU",
				"A product with this SKU already exists",
			))
			return
		}

		utils.ErrorResponse(c, utils.NewAppError(
			http.StatusInternalServerError,
			"Failed to create product",
			err.Error(),
		))
		return
	}

	utils.CreatedResponse(c, product)

}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	// Get pagination parameters
	page, pageSize := utils.GetPaginationParams(c)

	// get search and filter parameters
	search := c.Query("search")
	category := c.Query("category")

	// get products with pagination
	products, total, err := h.repo.GetProductsWithPagination(page, pageSize, search, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}

	// calculate pagination details
	totalPages := utils.CalculateTotalPages(total, pageSize)

	// create paginated response
	response := utils.PaginatedResponse{
		Data: products,
		Pagination: utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    totalPages,
		},
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get existing product
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// update fields
	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.SKU != "" {
		product.SKU = req.SKU
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Quantity >= 0 {
		product.Quantity = req.Quantity
	}
	if req.Category != "" {
		product.Category = req.Category
	}
	product.UpdatedAt = time.Now()

	if err := h.repo.UpdateProduct(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Check if product exists
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := h.repo.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
