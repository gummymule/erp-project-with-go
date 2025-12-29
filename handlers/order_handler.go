package handlers

import (
	"log"

	"erp-project/models"
	"erp-project/repositories"
	"erp-project/utils"

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
		utils.ValidationErrorResponse(c, "Validation error", err.Error())
		return
	}

	// Verify customer exists
	customer, err := h.customerRepo.GetCustomerByID(req.CustomerID)
	if err != nil {
		log.Printf("CreateOrder - GetCustomerByID error: %v", err)
		utils.InternalErrorResponse(c, "Failed to validate customer", "Database error")
		return
	}

	if customer == nil {
		utils.BadRequestResponse(c, "Invalid customer ID", "Customer not found")
		return
	}

	var totalAmount float64
	var orderItems []*models.OrderItem

	// process each item
	for _, itemReq := range req.Items {
		// get product
		product, err := h.productRepo.GetProductByID(itemReq.ProductID)
		if err != nil {
			log.Printf("CreateOrder - GetProductByID error: %v", err)
			utils.InternalErrorResponse(c, "Failed to validate product", "Database error")
			return
		}

		if product == nil {
			utils.BadRequestResponse(c, "Invalid product ID", map[string]interface{}{
				"product_id": itemReq.ProductID,
				"message":    "Product not found",
			})
			return
		}

		// check stock
		if product.Quantity < itemReq.Quantity {
			utils.BadRequestResponse(c, "Insufficient stock", map[string]interface{}{
				"product_id":   itemReq.ProductID,
				"product_name": product.Name,
				"available":    product.Quantity,
				"requested":    itemReq.Quantity,
			})
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
		log.Printf("CreateOrder - CreateOrderWithItems error: %v", err)
		utils.InternalErrorResponse(c, "Failed to create order", "Database error")
		return
	}

	// Prepare response data
	responseData := map[string]interface{}{
		"order": map[string]interface{}{
			"id":            order.ID,
			"customer_id":   order.CustomerID,
			"customer_name": customer.Name,
			"total_amount":  order.TotalAmount,
			"status":        order.Status,
			"order_date":    order.OrderDate,
		},
		"items": orderItems,
		"summary": map[string]interface{}{
			"total_items":  len(orderItems),
			"total_amount": totalAmount,
		},
	}

	utils.CreatedResponse(c, "Order created successfully", responseData)
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderRepo.GetOrders()
	if err != nil {
		log.Printf("GetOrders error: %v", err)
		utils.InternalErrorResponse(c, "Failed to fetch orders", "Database error")
		return
	}

	utils.SuccessResponse(c, "Orders retrieved successfully", orders)
}

func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	orderID := c.Param("id")

	items, err := h.orderRepo.GetOrderItems(orderID)
	if err != nil {
		log.Printf("GetOrderItems error: %v", err)
		utils.InternalErrorResponse(c, "Failed to fetch order items", "Database error")
		return
	}

	utils.SuccessResponse(c, "Order items retrieved successfully", items)
}
