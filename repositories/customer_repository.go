package repositories

import (
	"database/sql"
	"erp-project/models"
	"log"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (r *CustomerRepository) CreateCustomer(customer *models.Customer) error {
	query := `
		INSERT INTO customers (id, name, email, phone, address, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)
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
	query := `SELECT id, name, email, phone, address, created_at FROM customers WHERE id = ?`
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
