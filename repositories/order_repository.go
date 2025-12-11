package repositories

import (
	"database/sql"
	"erp-project/models"
	"log"
)

type OrderRepository struct {
	DB *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{DB: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}
	query := `
		INSERT INTO orders (id, customer_id, total_amount, status, order_date) 
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(
		query,
		order.ID,
		order.CustomerID,
		order.TotalAmount,
		order.Status,
		order.OrderDate,
	)

	if err != nil {
		tx.Rollback()
		log.Printf("Error creating order: %v", err)
		return err
	}

	return tx.Commit()
}

func (r *OrderRepository) CreateOrderWithItems(order *models.Order, items []*models.OrderItem) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		return err
	}

	// Insert order
	orderQuery := `
		INSERT INTO orders (id, customer_id, total_amount, status, order_date) 
		VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.Exec(
		orderQuery,
		order.ID,
		order.CustomerID,
		order.TotalAmount,
		order.Status,
		order.OrderDate,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Error creating order: %v", err)
		return err
	}

	// Insert order items
	itemQuery := `
		INSERT INTO order_items (id, order_id, product_id, quantity, unit_price, total_price) 
		VALUES (?, ?, ?, ?, ?, ?)
	`
	for _, item := range items {
		_, err = tx.Exec(
			itemQuery,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.UnitPrice,
			item.TotalPrice,
		)
		if err != nil {
			tx.Rollback()
			log.Printf("Error creating order item: %v", err)
			return err
		}

		// Update product quantity
		updateQuery := `
			UPDATE products SET quantity = quantity - ? WHERE id = ?
		`
		_, err = tx.Exec(updateQuery, item.Quantity, item.ProductID)
		if err != nil {
			tx.Rollback()
			log.Printf("Error updating product quantity: %v", err)
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrders() ([]*models.Order, error) {
	query := `
		SELECT o.id, o.customer_id, o.total_amount, o.status, o.order_date, c.name 
		FROM orders o
		LEFT JOIN customers c ON o.customer_id = c.id
		ORDER BY o.order_date DESC
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting orders: %v", err)
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		order := &models.Order{}
		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.TotalAmount,
			&order.Status,
			&order.OrderDate,
			&order.CustomerName,
		)
		if err != nil {
			log.Printf("Error scanning order: %v", err)
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *OrderRepository) GetOrderItems(orderID string) ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.unit_price, oi.total_price, p.name
		FROM order_items oi
		LEFT JOIN products p ON oi.product_id = p.id
		WHERE oi.order_id = ?
	`
	rows, err := r.DB.Query(query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.OrderItem{}
	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&item.ProductName,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
