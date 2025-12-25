package repositories

import (
	"database/sql"
	"erp-project/models"
	"erp-project/utils"
	"fmt"
	"log"
	"strings"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) CreateCustomer(customer *models.Customer) error {
	// FIXED: Changed ? to $1, $2, etc.
	query := `
		INSERT INTO customers (id, name, email, phone, address, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.DB.Exec(
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.Phone,
		customer.Address,
		customer.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating customer: %v", err)
		return err
	}
	return nil
}

func (r *CustomerRepository) GetCustomerWithPagination(page, pageSize int, search, email string) ([]models.Customer, int, error) {
	// Build WHERE clause with PostgreSQL placeholders
	var whereClauses []string
	var args []interface{}
	argCounter := 1

	if search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("(name LIKE $%d OR email LIKE $%d OR phone LIKE $%d)",
			argCounter, argCounter+1, argCounter+2))
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		argCounter += 3
	}

	if email != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("email = $%d", argCounter))
		args = append(args, email)
		argCounter++
	}

	whereClause := ""
	if len(whereClauses) > 0 {
		whereClause = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Get total count
	countQuery := "SELECT COUNT(*) FROM customers"
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
		SELECT id, name, email, phone, address, created_at
		FROM customers %s
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

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.Phone,
			&c.Address,
			&c.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning customer: %v", err)
			continue
		}
		customers = append(customers, c)
	}

	return customers, total, nil
}

func (r *CustomerRepository) GetAllCustomers() ([]*models.Customer, error) {
	query := `SELECT id, name, email, phone, address, created_at FROM customers`
	rows, err := r.DB.Query(query)
	if err != nil {
		log.Printf("Error getting customers: %v", err)
		return nil, err
	}
	defer rows.Close()

	var customers []*models.Customer
	for rows.Next() {
		customer := &models.Customer{}
		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.Phone,
			&customer.Address,
			&customer.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning customer: %v", err)
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *CustomerRepository) GetCustomerByID(id string) (*models.Customer, error) {
	// FIXED: Changed ? to $1
	query := `SELECT id, name, email, phone, address, created_at FROM customers WHERE id = $1`
	row := r.DB.QueryRow(query, id)

	customer := &models.Customer{}
	err := row.Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.Phone,
		&customer.Address,
		&customer.CreatedAt,
	)
	if err != nil {
		log.Printf("Error scanning customer: %v", err)
		return nil, err
	}
	return customer, nil
}
