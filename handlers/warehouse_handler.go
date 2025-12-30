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

type WarehouseHandler struct {
	repo *repositories.WarehouseRepository
}

func NewWarehouseHandler(repo *repositories.WarehouseRepository) *WarehouseHandler {
	return &WarehouseHandler{repo: repo}
}

// Request structs
type CreateWarehouseRequest struct {
	Code        string `json:"code" binding:"required,min=2,max=50"`
	Name        string `json:"name" binding:"required,min=2,max=255"`
	Location    string `json:"location" binding:"max=500"`
	ManagerName string `json:"manager_name" binding:"max=255"`
	Phone       string `json:"phone" binding:"omitempty,max=50"`
	Email       string `json:"email" binding:"omitempty,email"`
	Capacity    int    `json:"capacity" binding:"omitempty,min=0"`
}

type UpdateWarehouseRequest struct {
	Code        string `json:"code" binding:"omitempty,min=2,max=50"`
	Name        string `json:"name" binding:"omitempty,min=2,max=255"`
	Location    string `json:"location" binding:"omitempty,max=500"`
	ManagerName string `json:"manager_name" binding:"omitempty,max=255"`
	Phone       string `json:"phone" binding:"omitempty,max=50"`
	Email       string `json:"email" binding:"omitempty,email"`
	Capacity    int    `json:"capacity" binding:"omitempty,min=0"`
	Status      string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type CreateLocationRequest struct {
	LocationCode string `json:"location_code" binding:"required,min=1,max=50"`
	LocationName string `json:"location_name" binding:"max=255"`
	Zone         string `json:"zone" binding:"max=50"`
	RowNumber    int    `json:"row_number" binding:"omitempty,min=0"`
	ShelfNumber  int    `json:"shelf_number" binding:"omitempty,min=0"`
	MaxCapacity  int    `json:"max_capacity" binding:"omitempty,min=0"`
}

type UpdateLocationRequest struct {
	LocationCode string `json:"location_code" binding:"omitempty,min=1,max=50"`
	LocationName string `json:"location_name" binding:"omitempty,max=255"`
	Zone         string `json:"zone" binding:"omitempty,max=50"`
	RowNumber    int    `json:"row_number" binding:"omitempty,min=0"`
	ShelfNumber  int    `json:"shelf_number" binding:"omitempty,min=0"`
	MaxCapacity  int    `json:"max_capacity" binding:"omitempty,min=0"`
	Status       string `json:"status" binding:"omitempty,oneof=available occupied reserved maintenance"`
}

type CreateInventoryRequest struct {
	ProductID   string `json:"product_id" binding:"required"`
	LocationID  string `json:"location_id" binding:""`
	Quantity    int    `json:"quantity" binding:"required,min=0"`
	MinQuantity int    `json:"min_quantity" binding:"omitempty,min=0"`
}

type UpdateInventoryRequest struct {
	Quantity    int `json:"quantity" binding:"omitempty,min=0"`
	MinQuantity int `json:"min_quantity" binding:"omitempty,min=0"`
}

// Warehouse CRUD handlers
func (h *WarehouseHandler) CreateWarehouse(c *gin.Context) {
	var req CreateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	warehouse := models.NewWarehouse(
		req.Code,
		req.Name,
		req.Location,
		req.ManagerName,
		req.Phone,
		req.Email,
		req.Capacity,
	)

	if err := h.repo.CreateWarehouse(warehouse); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			utils.DuplicateErrorResponse(c, "Duplicate warehouse code", "A warehouse with this code already exists")
			return
		}

		log.Printf("CreateWarehouse error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create warehouse", "Database error")
		return
	}

	utils.CreatedResponse(c, "Warehouse created successfully", warehouse)
}

func (h *WarehouseHandler) GetAllWarehouses(c *gin.Context) {
	warehouses, err := h.repo.GetAllWarehouses()
	if err != nil {
		log.Printf("GetAllWarehouses error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouses", "Database error")
		return
	}

	utils.SuccessResponse(c, "Warehouses retrieved successfully", warehouses)
}

func (h *WarehouseHandler) GetWarehouseByID(c *gin.Context) {
	id := c.Param("id")
	warehouse, err := h.repo.GetWarehouseByID(id)
	if err != nil {
		log.Printf("GetWarehouseByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouse", "Database error")
		return
	}
	if warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}
	utils.SuccessResponse(c, "Warehouse retrieved successfully", warehouse)
}

func (h *WarehouseHandler) UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")

	var req UpdateWarehouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Get existing warehouse
	warehouse, err := h.repo.GetWarehouseByID(id)
	if err != nil {
		log.Printf("UpdateWarehouse - GetWarehouseByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouse", "Database error")
		return
	}

	if warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	// Update fields only if provided
	updatedFields := []string{}
	if req.Code != "" && req.Code != warehouse.Code {
		warehouse.Code = req.Code
		updatedFields = append(updatedFields, "code")
	}
	if req.Name != "" && req.Name != warehouse.Name {
		warehouse.Name = req.Name
		updatedFields = append(updatedFields, "name")
	}
	if req.Location != "" && req.Location != warehouse.Location {
		warehouse.Location = req.Location
		updatedFields = append(updatedFields, "location")
	}
	if req.ManagerName != "" && req.ManagerName != warehouse.ManagerName {
		warehouse.ManagerName = req.ManagerName
		updatedFields = append(updatedFields, "manager_name")
	}
	if req.Phone != "" && req.Phone != warehouse.Phone {
		warehouse.Phone = req.Phone
		updatedFields = append(updatedFields, "phone")
	}
	if req.Email != "" && req.Email != warehouse.Email {
		warehouse.Email = req.Email
		updatedFields = append(updatedFields, "email")
	}
	if req.Capacity > 0 && req.Capacity != warehouse.Capacity {
		warehouse.Capacity = req.Capacity
		updatedFields = append(updatedFields, "capacity")
	}
	if req.Status != "" && req.Status != warehouse.Status {
		warehouse.Status = req.Status
		updatedFields = append(updatedFields, "status")
	}

	// If no fields were updated
	if len(updatedFields) == 0 {
		utils.SuccessResponse(c, "No changes detected", warehouse)
		return
	}

	warehouse.UpdatedAt = time.Now()

	if err := h.repo.UpdateWarehouse(warehouse); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate warehouse code", "A warehouse with this code already exists")
			return
		}

		log.Printf("UpdateWarehouse error: %v", err)
		utils.InternalErrorResponse(c, "Failed to update warehouse", "Database error")
		return
	}

	// Get updated warehouse
	updatedWarehouse, err := h.repo.GetWarehouseByID(id)
	if err != nil {
		log.Printf("UpdateWarehouse - Get updated warehouse error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve updated warehouse", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"warehouse":      updatedWarehouse,
		"updated_fields": updatedFields,
	}

	utils.SuccessResponse(c, "Warehouse updated successfully", responseData)
}

func (h *WarehouseHandler) DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	warehouse, err := h.repo.GetWarehouseByID(id)
	if err != nil {
		log.Printf("DeleteWarehouse - GetWarehouseByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouse", "Database error")
		return
	}
	if warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	if err := h.repo.DeleteWarehouse(id); err != nil {
		log.Printf("DeleteWarehouse error: %v", err)
		utils.InternalErrorResponse(c, "Failed to delete warehouse", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"deleted_warehouse_id":   id,
		"deleted_warehouse_name": warehouse.Name,
	}

	utils.SuccessResponse(c, "Warehouse deleted successfully", responseData)
}

// Location handlers
func (h *WarehouseHandler) CreateLocation(c *gin.Context) {
	warehouseID := c.Param("id")

	var req CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Verify warehouse exists
	warehouse, err := h.repo.GetWarehouseByID(warehouseID)
	if err != nil || warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	location := models.NewWarehouseLocation(
		warehouseID,
		req.LocationCode,
		req.LocationName,
		req.Zone,
		req.RowNumber,
		req.ShelfNumber,
		req.MaxCapacity,
	)

	if err := h.repo.CreateLocation(location); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate location code", "A location with this code already exists in this warehouse")
			return
		}

		log.Printf("CreateLocation error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create location", "Database error")
		return
	}

	utils.CreatedResponse(c, "Location created successfully", location)
}

func (h *WarehouseHandler) GetWarehouseLocations(c *gin.Context) {
	warehouseID := c.Param("id")

	// Verify warehouse exists
	warehouse, err := h.repo.GetWarehouseByID(warehouseID)
	if err != nil || warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	locations, err := h.repo.GetLocationsByWarehouse(warehouseID)
	if err != nil {
		log.Printf("GetWarehouseLocations error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouse locations", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"warehouse": warehouse,
		"locations": locations,
		"count":     len(locations),
	}

	utils.SuccessResponse(c, "Warehouse locations retrieved successfully", responseData)
}

func (h *WarehouseHandler) GetAvailableLocations(c *gin.Context) {
	warehouseID := c.Param("id")

	// Verify warehouse exists
	warehouse, err := h.repo.GetWarehouseByID(warehouseID)
	if err != nil || warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	locations, err := h.repo.GetAvailableLocations(warehouseID)
	if err != nil {
		log.Printf("GetAvailableLocations error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve available locations", "Database error")
		return
	}

	responseData := map[string]interface{}{
		"warehouse": warehouse,
		"locations": locations,
		"count":     len(locations),
	}

	utils.SuccessResponse(c, "Available locations retrieved successfully", responseData)
}

// Inventory handlers
func (h *WarehouseHandler) CreateInventory(c *gin.Context) {
	warehouseID := c.Param("id")

	var req CreateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Verify warehouse exists
	warehouse, err := h.repo.GetWarehouseByID(warehouseID)
	if err != nil || warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	// Handle location ID
	var locationID *string
	if req.LocationID != "" {
		locationID = &req.LocationID
	}

	// Set default min quantity if not provided
	minQuantity := req.MinQuantity
	if minQuantity == 0 {
		minQuantity = 10 // Default value
	}

	inventory := models.NewInventory(
		req.ProductID,
		warehouseID,
		locationID,
		req.Quantity,
		minQuantity,
	)

	if err := h.repo.CreateInventory(inventory); err != nil {
		if strings.Contains(err.Error(), "duplicate key value") ||
			strings.Contains(err.Error(), "violates unique constraint") {
			utils.DuplicateErrorResponse(c, "Duplicate inventory", "Inventory for this product already exists in this location")
			return
		}

		log.Printf("CreateInventory error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create inventory", "Database error")
		return
	}

	utils.CreatedResponse(c, "Inventory created successfully", inventory)
}

func (h *WarehouseHandler) GetWarehouseInventory(c *gin.Context) {
	warehouseID := c.Param("id")

	// Verify warehouse exists
	warehouse, err := h.repo.GetWarehouseByID(warehouseID)
	if err != nil || warehouse == nil {
		utils.NotFoundResponse(c, "Warehouse not found")
		return
	}

	inventory, err := h.repo.GetInventoryByWarehouse(warehouseID)
	if err != nil {
		log.Printf("GetWarehouseInventory error: %v", err)
		utils.InternalErrorResponse(c, "Failed to retrieve warehouse inventory", "Database error")
		return
	}

	// Calculate totals
	var totalQuantity, totalReserved, totalAvailable int
	for _, inv := range inventory {
		totalQuantity += inv.Quantity
		totalReserved += inv.ReservedQuantity
		totalAvailable += inv.AvailableQuantity
	}

	responseData := map[string]interface{}{
		"warehouse": warehouse,
		"inventory": inventory,
		"summary": map[string]interface{}{
			"total_items":     len(inventory),
			"total_quantity":  totalQuantity,
			"total_reserved":  totalReserved,
			"total_available": totalAvailable,
		},
	}

	utils.SuccessResponse(c, "Warehouse inventory retrieved successfully", responseData)
}

func (h *WarehouseHandler) UpdateInventory(c *gin.Context) {
	inventoryID := c.Param("inventory_id")

	var req UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Note: For simplicity, we're just updating quantity here
	// In a real app, you'd have more complex inventory management

	utils.SuccessResponse(c, "Inventory update endpoint (to be implemented)", map[string]interface{}{
		"inventory_id": inventoryID,
		"request":      req,
	})
}
