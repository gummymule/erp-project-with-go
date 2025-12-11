package handlers

import (
	"net/http"

	"erp-project/models"
	"erp-project/repositories"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderRepo    *repositories.OrderRepository
	productRepo  *repositories.ProductRepository
	customerRepo *repositories.CustomerRepository
}

func NewOrderHandler(
	orderRepo *repositories.OrderRepository,
	productRepo *repositories.ProductRepository,
	customerRepo *repositories.CustomerRepository) *OrderHandler {
	return &OrderHandler{
		orderRepo:    orderRepo,
		productRepo:  productRepo,
		customerRepo: customerRepo,
	}
}

type CreateOrderRequest struct {
	CustomerID string             `json:"customer_id" binding:"required"`
	Items      []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify costumer exists
	customer, err := h.customerRepo.GetCustomerByID(req.CustomerID)
	if err != nil || customer == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer ID"})
		return
	}

	var totalAmount float64
	var orderItems []*models.OrderItem

	// process each item
	for _, itemReq := range req.Items {
		// get product
		product, err := h.productRepo.GetProductByID(itemReq.ProductID)
		if err != nil || product == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID: " + itemReq.ProductID})
			return
		}

		// check stock
		if product.Quantity < itemReq.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + product.Name})
			return
		}

		// calculate total amount
		itemTotal := product.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		// create order item
		orderItem := models.NewOrderItem("", product.ID, itemReq.Quantity, product.Price)
		orderItems = append(orderItems, orderItem)
	}

	// create order
	order := models.NewOrder(req.CustomerID, totalAmount)

	// update order items with order ID
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}

	// save order with items (transaction)
	if err := h.orderRepo.CreateOrderWithItems(order, orderItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// get order details with items for response
	orderWithItems := gin.H{
		"order": order,
		"items": orderItems,
	}

	c.JSON(http.StatusCreated, orderWithItems)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderRepo.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	orderID := c.Param("id")

	items, err := h.orderRepo.GetOrderItems(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch order items"})
		return
	}

	c.JSON(http.StatusOK, items)
}
