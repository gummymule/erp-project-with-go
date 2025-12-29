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

type ProductHandler struct {
	repo *repositories.ProductRepository
}

func NewProductHandler(repo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2,max=100"`
	Description string  `json:"description" binding:"max=500"`
	SKU         string  `json:"sku" binding:"required,min=3,max=50"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    int     `json:"quantity" binding:"required,gte=0"`
	Category    string  `json:"category" binding:"max=50"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name" binding:"omitempty,min=2,max=100"`
	Description string  `json:"description" binding:"omitempty,max=500"`
	SKU         string  `json:"sku" binding:"omitempty,min=3,max=50"`
	Price       float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity    int     `json:"quantity" binding:"omitempty,gte=0"`
	Category    string  `json:"category" binding:"omitempty,max=50"`
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
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
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") ||
			strings.Contains(err.Error(), "already exists") {
			utils.DuplicateErrorResponse(c, "Duplicate SKU", "A product with this SKU already exists")
			return
		}

		log.Printf("CreateProduct error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create product", "Database error occurred")
		return
	}

	utils.CreatedResponse(c, "Product created successfully", product)
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
		log.Printf("GetAllProducts error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve products", "Database error")
		return
	}

	// calculate pagination details
	totalPages := utils.CalculateTotalPages(total, pageSize)

	// create paginated response
	paginationData := map[string]interface{}{
		"products": products,
		"pagination": utils.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
			Pages:    totalPages,
		},
	}

	utils.SuccessResponse(c, "Products retrieved successfully", paginationData)
}

func (h *ProductHandler) GetListProducts(c *gin.Context) {
	products, err := h.repo.GetAll()
	if err != nil {
		log.Printf("GetListProducts error: %v", err)
		utils.InternalErrorResponse(c, "Failed to fetch products", "Database error")
		return
	}

	utils.SuccessResponse(c, "Products list retrieved successfully", products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		log.Printf("GetProductByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve product", "Database error")
		return
	}

	if product == nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	utils.SuccessResponse(c, "Product retrieved successfully", product)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// get existing product
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		log.Printf("UpdateProduct - GetProductByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve product", "Database error")
		return
	}

	if product == nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	// update fields
	updatedFields := []string{}
	if req.Name != "" && req.Name != product.Name {
		product.Name = req.Name
		updatedFields = append(updatedFields, "name")
	}
	if req.Description != "" && req.Description != product.Description {
		product.Description = req.Description
		updatedFields = append(updatedFields, "description")
	}
	if req.SKU != "" && req.SKU != product.SKU {
		product.SKU = req.SKU
		updatedFields = append(updatedFields, "sku")
	}
	if req.Price > 0 && req.Price != product.Price {
		product.Price = req.Price
		updatedFields = append(updatedFields, "price")
	}
	if req.Quantity >= 0 && req.Quantity != product.Quantity {
		product.Quantity = req.Quantity
		updatedFields = append(updatedFields, "quantity")
	}
	if req.Category != "" && req.Category != product.Category {
		product.Category = req.Category
		updatedFields = append(updatedFields, "category")
	}

	// If no fields were updated
	if len(updatedFields) == 0 {
		utils.SuccessResponse(c, "No changes detected", product)
		return
	}

	product.UpdatedAt = time.Now()

	if err := h.repo.UpdateProduct(product); err != nil {
		log.Printf("UpdateProduct error: %v", err)
		utils.InternalErrorResponse(c, "Failed to update product", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"product":        product,
		"updated_fields": updatedFields,
	}

	utils.SuccessResponse(c, "Product updated successfully", responseData)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Check if product exists
	product, err := h.repo.GetProductByID(id)
	if err != nil {
		log.Printf("DeleteProduct - GetProductByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to fetch product", "Database error")
		return
	}

	if product == nil {
		utils.NotFoundResponse(c, "Product not found")
		return
	}

	if err := h.repo.DeleteProduct(id); err != nil {
		log.Printf("DeleteProduct error: %v", err)
		utils.InternalErrorResponse(c, "Failed to delete product", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"deleted_product_id":   id,
		"deleted_product_name": product.Name,
	}

	utils.SuccessResponse(c, "Product deleted successfully", responseData)
}
