package repositories

import (
	"database/sql"
	"erp-project/models"
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

// Get all Products
func (r *ProductRepository) GetAllProducts() ([]*models.Product, error) {
	query := `SELECT id, name, description, sku, price, quantity, category, created_at, updated_at FROM products`
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting products: %v", err)
		return nil, err
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		product := &models.Product{}
		err := rows.Scan(
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
		products = append(products, product)
	}
	return products, nil
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
