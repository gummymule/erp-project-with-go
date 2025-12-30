package repositories

import (
	"database/sql"
	"erp-project/models"
	"fmt"
	"log"
	"strings"
	"time"
)

type SupplierRepository struct {
	DB *sql.DB
}

func NewSupplierRepository(db *sql.DB) *SupplierRepository {
	return &SupplierRepository{DB: db}
}

func (r *SupplierRepository) CreateSupplier(supplier *models.Supplier) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `INSERT INTO suppliers (id, name, code, contact_person, email, phone, address, tax_id, payment_terms, status, created_at, updated_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	} else {
		query = `INSERT INTO suppliers (id, name, code, contact_person, email, phone, address, tax_id, payment_terms, status, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.DB.Exec(
		query,
		supplier.ID,
		supplier.Name,
		supplier.Code,
		supplier.ContactPerson,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.TaxID,
		supplier.PaymentTerms,
		supplier.Status,
		supplier.CreatedAt,
		supplier.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating supplier: %v", err)
		return err
	}
	return nil
}

func (r *SupplierRepository) GetAllSuppliers() ([]models.Supplier, error) {
	query := `SELECT id, name, code, contact_person, email, phone, address, tax_id, payment_terms, status, created_at, updated_at 
	          FROM suppliers ORDER BY name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppliers []models.Supplier
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Code,
			&s.ContactPerson,
			&s.Email,
			&s.Phone,
			&s.Address,
			&s.TaxID,
			&s.PaymentTerms,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning supplier: %v", err)
			continue
		}
		suppliers = append(suppliers, s)
	}

	return suppliers, nil
}

func (r *SupplierRepository) GetSupplierByID(id string) (*models.Supplier, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT id, name, code, contact_person, email, phone, address, tax_id, payment_terms, status, created_at, updated_at 
		         FROM suppliers WHERE id = $1`
	} else {
		query = `SELECT id, name, code, contact_person, email, phone, address, tax_id, payment_terms, status, created_at, updated_at 
		         FROM suppliers WHERE id = ?`
	}

	row := r.DB.QueryRow(query, id)

	var s models.Supplier
	err := row.Scan(
		&s.ID,
		&s.Name,
		&s.Code,
		&s.ContactPerson,
		&s.Email,
		&s.Phone,
		&s.Address,
		&s.TaxID,
		&s.PaymentTerms,
		&s.Status,
		&s.CreatedAt,
		&s.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting supplier: %v", err)
		return nil, err
	}

	return &s, nil
}

func (r *SupplierRepository) UpdateSupplier(supplier *models.Supplier) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `UPDATE suppliers SET name = $1, code = $2, contact_person = $3, email = $4, phone = $5, 
		         address = $6, tax_id = $7, payment_terms = $8, status = $9, updated_at = $10 
		         WHERE id = $11`
	} else {
		query = `UPDATE suppliers SET name = ?, code = ?, contact_person = ?, email = ?, phone = ?, 
		         address = ?, tax_id = ?, payment_terms = ?, status = ?, updated_at = ? 
		         WHERE id = ?`
	}

	_, err := r.DB.Exec(
		query,
		supplier.Name,
		supplier.Code,
		supplier.ContactPerson,
		supplier.Email,
		supplier.Phone,
		supplier.Address,
		supplier.TaxID,
		supplier.PaymentTerms,
		supplier.Status,
		time.Now(),
		supplier.ID,
	)

	if err != nil {
		log.Printf("Error updating supplier: %v", err)
		return err
	}
	return nil
}

func (r *SupplierRepository) DeleteSupplier(id string) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `DELETE FROM suppliers WHERE id = $1`
	} else {
		query = `DELETE FROM suppliers WHERE id = ?`
	}

	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting supplier: %v", err)
		return err
	}
	return nil
}

// ProductSupplier methods
func (r *SupplierRepository) AddProductSupplier(ps *models.ProductSupplier) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `INSERT INTO product_suppliers (id, product_id, supplier_id, supplier_sku, cost_price, lead_time_days, is_primary, created_at, updated_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	} else {
		query = `INSERT INTO product_suppliers (id, product_id, supplier_id, supplier_sku, cost_price, lead_time_days, is_primary, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.DB.Exec(
		query,
		ps.ID,
		ps.ProductID,
		ps.SupplierID,
		ps.SupplierSKU,
		ps.CostPrice,
		ps.LeadTimeDays,
		ps.IsPrimary,
		ps.CreatedAt,
		ps.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error adding product supplier: %v", err)
		return err
	}
	return nil
}

func (r *SupplierRepository) GetProductSuppliers(productID string) ([]models.ProductSupplier, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT ps.id, ps.product_id, ps.supplier_id, ps.supplier_sku, ps.cost_price, 
		                ps.lead_time_days, ps.is_primary, ps.created_at, ps.updated_at,
		                p.name as product_name, s.name as supplier_name
		         FROM product_suppliers ps
		         JOIN products p ON ps.product_id = p.id
		         JOIN suppliers s ON ps.supplier_id = s.id
		         WHERE ps.product_id = $1
		         ORDER BY ps.is_primary DESC, s.name`
	} else {
		query = `SELECT ps.id, ps.product_id, ps.supplier_id, ps.supplier_sku, ps.cost_price, 
		                ps.lead_time_days, ps.is_primary, ps.created_at, ps.updated_at,
		                p.name as product_name, s.name as supplier_name
		         FROM product_suppliers ps
		         JOIN products p ON ps.product_id = p.id
		         JOIN suppliers s ON ps.supplier_id = s.id
		         WHERE ps.product_id = ?
		         ORDER BY ps.is_primary DESC, s.name`
	}

	rows, err := r.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productSuppliers []models.ProductSupplier
	for rows.Next() {
		var ps models.ProductSupplier
		err := rows.Scan(
			&ps.ID,
			&ps.ProductID,
			&ps.SupplierID,
			&ps.SupplierSKU,
			&ps.CostPrice,
			&ps.LeadTimeDays,
			&ps.IsPrimary,
			&ps.CreatedAt,
			&ps.UpdatedAt,
			&ps.ProductName,
			&ps.SupplierName,
		)
		if err != nil {
			log.Printf("Error scanning product supplier: %v", err)
			continue
		}
		productSuppliers = append(productSuppliers, ps)
	}

	return productSuppliers, nil
}

func (r *SupplierRepository) GetSupplierProducts(supplierID string) ([]models.ProductSupplier, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT ps.id, ps.product_id, ps.supplier_id, ps.supplier_sku, ps.cost_price, 
		                ps.lead_time_days, ps.is_primary, ps.created_at, ps.updated_at,
		                p.name as product_name, s.name as supplier_name
		         FROM product_suppliers ps
		         JOIN products p ON ps.product_id = p.id
		         JOIN suppliers s ON ps.supplier_id = s.id
		         WHERE ps.supplier_id = $1
		         ORDER BY p.name`
	} else {
		query = `SELECT ps.id, ps.product_id, ps.supplier_id, ps.supplier_sku, ps.cost_price, 
		                ps.lead_time_days, ps.is_primary, ps.created_at, ps.updated_at,
		                p.name as product_name, s.name as supplier_name
		         FROM product_suppliers ps
		         JOIN products p ON ps.product_id = p.id
		         JOIN suppliers s ON ps.supplier_id = s.id
		         WHERE ps.supplier_id = ?
		         ORDER BY p.name`
	}

	rows, err := r.DB.Query(query, supplierID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productSuppliers []models.ProductSupplier
	for rows.Next() {
		var ps models.ProductSupplier
		err := rows.Scan(
			&ps.ID,
			&ps.ProductID,
			&ps.SupplierID,
			&ps.SupplierSKU,
			&ps.CostPrice,
			&ps.LeadTimeDays,
			&ps.IsPrimary,
			&ps.CreatedAt,
			&ps.UpdatedAt,
			&ps.ProductName,
			&ps.SupplierName,
		)
		if err != nil {
			log.Printf("Error scanning supplier product: %v", err)
			continue
		}
		productSuppliers = append(productSuppliers, ps)
	}

	return productSuppliers, nil
}

func (r *SupplierRepository) RemoveProductSupplier(id string) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `DELETE FROM product_suppliers WHERE id = $1`
	} else {
		query = `DELETE FROM product_suppliers WHERE id = ?`
	}

	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error removing product supplier: %v", err)
		return err
	}
	return nil
}

// Helper function to check database type
func isPostgreSQL(db *sql.DB) bool {
	driver := fmt.Sprintf("%T", db.Driver())
	return strings.Contains(driver, "postgres")
}
