package database

import (
	"database/sql"
	"log"

	"github.com/glebarez/sqlite"
)

var DB *sql.DB

func InitDB() error {
	var err error

	// This one uses "sqlite" as driver name
	dbPath := "erp.db"
	DB, err = sql.Open(sqlite.DriverName, dbPath)
	if err != nil {
		return err
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Connected to SQLite database")
	return createTables()
}

func createTables() error {
	// Products table
	productTable := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		sku TEXT UNIQUE NOT NULL,
		price REAL NOT NULL,
		quantity INTEGER NOT NULL,
		category TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Customers table
	customerTable := `
	CREATE TABLE IF NOT EXISTS customers (
		id TEXT PRIMARY KEY, 
		name TEXT NOT NULL, 
		email TEXT UNIQUE NOT NULL,
		phone TEXT,
		address TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	// Orders table
	orderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		customer_id TEXT NOT NULL,
		total_amount REAL NOT NULL,
		status TEXT DEFAULT 'pending',
		order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (customer_id) REFERENCES customers(id)
	);
	`

	// Order Items table (many-to-many relationship)
	orderItemTable := `
	CREATE TABLE IF NOT EXISTS order_items (
		id TEXT PRIMARY KEY,
		order_id TEXT NOT NULL,
		product_id TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		unit_price REAL NOT NULL,
		total_price REAL NOT NULL,
		FOREIGN KEY (order_id) REFERENCES orders(id),
		FOREIGN KEY (product_id) REFERENCES products(id)
	);
	`

	tables := []string{productTable, customerTable, orderTable, orderItemTable}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return err
		}
	}

	log.Println("Tables created successfully")
	return nil
}
