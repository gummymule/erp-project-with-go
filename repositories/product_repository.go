package repositories

import (
	"database/sql"
	"erp-project/models"
	"erp-project/utils"
	"fmt"
	"log"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

// Create product
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, sku, price, quantity, category, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
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

// Get products with pagination
func (r *ProductRepository) GetProductsWithPagination(page, pageSize int, search, category string) ([]models.Product, int, error) {
	// build where clause
	var whereClause string
	var args []interface{}

	if search != "" {
		whereClause = "WHERE (name LIKE ? OR description LIKE ? OR sku LIKE ?)"
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if category != "" {
		if whereClause != "" {
			whereClause += " AND category = ?"
		} else {
			whereClause = "WHERE category = ?"
		}
		args = append(args, category)
	}

	// get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products %s", whereClause)
	var total int
	err := r.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// get paginated data
	offset := utils.CalculateOffset(page, pageSize)
	query := fmt.Sprintf(`
		SELECT id, name, description, sku, price, quantity, category, created_at, updated_at
		FROM products %s
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)

	// add pagination parameters to args
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

// Get product by ID
func (r *ProductRepository) GetProductByID(id string) (*models.Product, error) {
	query := `SELECT id, name, description, sku, price, quantity, category, created_at, updated_at FROM products WHERE id = ?`
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

// Update product
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := `
		UPDATE products 
		SET name = ?, description = ?, sku = ?, price = ?, quantity = ?, category = ?, updated_at = ? 
		WHERE id = ?
	`
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

// Delete product
func (r *ProductRepository) DeleteProduct(id string) error {
	query := `DELETE FROM products WHERE id = ?`
	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting product: %v", err)
		return err
	}
	return nil
}

// Update product quantity
func (r *ProductRepository) UpdateProductQuantity(id string, quantity int) error {
	query := `UPDATE products SET quantity = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.DB.Exec(query, quantity, id)
	if err != nil {
		log.Printf("Error updating product quantity: %v", err)
		return err
	}
	return nil

}
