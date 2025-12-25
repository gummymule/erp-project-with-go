package repositories

import (
	"database/sql"
	"erp-project/models"
	"erp-project/utils"
	"fmt"
	"log"
	"strings"
	"time"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

// Create product - FIXED: Changed ? to $1, $2, etc.
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, sku, price, quantity, category, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.DB.Exec(
		query,
		product.ID,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Quantity,
		product.Category,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating product: %v", err)
		return err
	}
	return nil
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, name, description, sku, price, quantity, category, created_at, updated_at FROM products`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.SKU,
			&p.Price,
			&p.Quantity,
			&p.Category,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning product:", err)
			continue
		}
		products = append(products, p)
	}

	return products, nil
}

// Get products with pagination - FIXED: Parameter placeholders
func (r *ProductRepository) GetProductsWithPagination(page, pageSize int, search, category string) ([]models.Product, int, error) {
	// Build WHERE clause with PostgreSQL placeholders
	var whereClauses []string
	var args []interface{}
	argCounter := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name LIKE $%d OR description LIKE $%d OR sku LIKE $%d)",
			argCounter, argCounter+1, argCounter+2))
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		argCounter += 3
	}

	if category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argCounter))
		args = append(args, category)
		argCounter++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM products"
	if whereClause != "" {
		countQuery += " " + whereClause
	}

	var total int
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := utils.CalculateOffset(page, pageSize)
	query := fmt.Sprintf(`
		SELECT id, name, description, sku, price, quantity, category, created_at, updated_at
		FROM products %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCounter, argCounter+1)

	// Add pagination parameters
	args = append(args, pageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.SKU,
			&p.Price,
			&p.Quantity,
			&p.Category,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning product: %v\n", err)
			continue
		}
		products = append(products, p)
	}

	return products, total, nil
}

// Get all Products
func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	products, _, err := r.GetProductsWithPagination(1, 1000, "", "")
	return products, err
}

// Get product by ID - FIXED: Changed ? to $1
func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	query := `SELECT id, name, description, sku, price, quantity, category, created_at, updated_at FROM products WHERE id = $1`
	row := r.DB.QueryRow(query, id)
	product := &models.Product{}
	err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SKU,
		&product.Price,
		&product.Quantity,
		&product.Category,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		log.Printf("Error scanning product: %v", err)
		return nil, err
	}
	return product, nil
}

// Update product - FIXED: Changed ? to $1, $2, etc.
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = $1, description = $2, sku = $3, price = $4, quantity = $5, category = $6, updated_at = $7 
		WHERE id = $8
	`
	// Update timestamp
	product.UpdatedAt = time.Now()

	_, err := r.DB.Exec(
		query,
		product.Name,
		product.Description,
		product.SKU,
		product.Price,
		product.Quantity,
		product.Category,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return err
	}
	return nil
}

// Delete product - FIXED: Changed ? to $1
func (r *ProductRepository) DeleteProduct(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return err
	}
	return nil
}

// Update product quantity - FIXED: PostgreSQL CURRENT_TIMESTAMP syntax
func (r *ProductRepository) UpdateProductQuantity(id string, quantity int) error {
	query := `UPDATE products SET quantity = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.DB.Exec(query, quantity, id)
	if err != nil {
		log.Printf("Error updating product quantity: %v", err)
		return err
	}
	return nil
}
