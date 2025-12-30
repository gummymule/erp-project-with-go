package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitDB() error {
	var err error
	var driverName string

	// Check if we're on Railway (PostgreSQL) or local (SQLite)
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL != "" && strings.Contains(dbURL, "postgresql://") {
		// Railway PostgreSQL
		driverName = "postgres"
		DB, err = sql.Open(driverName, dbURL)
		log.Println("Using PostgreSQL (Railway)")
	} else {
		// Local SQLite development
		driverName = "sqlite"
		dbPath := "erp.db"

		// You can also use in-memory for testing
		if os.Getenv("IN_MEMORY_DB") == "true" {
			dbPath = ":memory:"
			log.Println("Using in-memory SQLite")
		} else {
			log.Printf("Using SQLite file: %s", dbPath)
		}

		DB, err = sql.Open(sqlite.DriverName, dbPath)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return createTables()
}

func createTables() error {
	// Determine database type for SQL dialect
	dbURL := os.Getenv("DATABASE_URL")
	isPostgreSQL := dbURL != "" && strings.Contains(dbURL, "postgresql://")

	var (
		productTable, customerTable, orderTable, orderItemTable string
		supplierTable, productSupplierTable                     string
		warehouseTable, warehouseLocationTable, inventoryTable  string
	)

	if isPostgreSQL {
		// PostgreSQL syntax
		productTable = `
		CREATE TABLE IF NOT EXISTS products (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			sku VARCHAR(255) UNIQUE NOT NULL,
			price DECIMAL(10,2) NOT NULL,
			quantity INTEGER NOT NULL,
			category VARCHAR(100),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		customerTable = `
		CREATE TABLE IF NOT EXISTS customers (
			id VARCHAR(255) PRIMARY KEY, 
			name VARCHAR(255) NOT NULL, 
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(50),
			address TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		orderTable = `
		CREATE TABLE IF NOT EXISTS orders (
			id VARCHAR(255) PRIMARY KEY,
			customer_id VARCHAR(255) NOT NULL,
			total_amount DECIMAL(10,2) NOT NULL,
			status VARCHAR(50) DEFAULT 'pending',
			order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE
		);`

		orderItemTable = `
		CREATE TABLE IF NOT EXISTS order_items (
			id VARCHAR(255) PRIMARY KEY,
			order_id VARCHAR(255) NOT NULL,
			product_id VARCHAR(255) NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price DECIMAL(10,2) NOT NULL,
			total_price DECIMAL(10,2) NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
		);`

		supplierTable = `
		CREATE TABLE IF NOT EXISTS suppliers (
			id VARCHAR(255) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			code VARCHAR(50) UNIQUE NOT NULL,
			contact_person VARCHAR(255),
			email VARCHAR(255),
			phone VARCHAR(50),
			address TEXT, 
			tax_id VARCHAR(100),
			payment_terms TEXT,
			status VARCHAR(50) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		productSupplierTable = `
		CREATE TABLE IF NOT EXISTS product_suppliers (
			id VARCHAR(255) PRIMARY KEY,
			product_id VARCHAR(255) NOT NULL,
			supplier_id VARCHAR(255) NOT NULL,
			supplier_sku VARCHAR(100),
			cost_price DECIMAL(10,2),
			lead_time_days INTEGER,
			is_primary BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY (supplier_id) REFERENCES suppliers(id) ON DELETE CASCADE,
			UNIQUE(product_id, supplier_id)
		);`

		warehouseTable = `
		CREATE TABLE IF NOT EXISTS warehouses (
			id VARCHAR(255) PRIMARY KEY,
			code VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			location TEXT,
			manager_name VARCHAR(255),
			phone VARCHAR(50),
			email VARCHAR(255),
			capacity INTEGER,
			status VARCHAR(50) DEFAULT 'active', --active, inactive
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		warehouseLocationTable = `
		CREATE TABLE IF NOT EXISTS warehouse_locations (
			id VARCHAR(255) PRIMARY KEY,
			warehouse_id VARCHAR(255) NOT NULL,
			location_code VARCHAR(100) NOT NULL,
			location_name VARCHAR(255),
			zone VARCHAR(50),
			row_number INTEGER,
			shelf_number INTEGER,
			max_capacity INTEGER,
			current_quantity INTEGER DEFAULT 0,
			status VARCHAR(50) DEFAULT 'available', -- available, full, inactive, ocupied, reserved, maintanance
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
			UNIQUE(warehouse_id, location_code)
		);`

		inventoryTable = `
		CREATE TABLE IF NOT EXISTS inventory (
			id VARCHAR(255) PRIMARY KEY,
			product_id VARCHAR(255) NOT NULL,
			warehouse_id VARCHAR(255) NOT NULL,
			location_id VARCHAR(255),
			quantity INTEGER NOT NULL DEFAULT 0,
			reserved_quantity INTEGER DEFAULT 0,
			min_quantity INTEGER DEFAULT 0,
			max_quantity INTEGER,
			last_restocked TIMESTAMP,
			last_checked TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id) ON DELETE CASCADE,
			FOREIGN KEY (location_id) REFERENCES warehouse_locations(id) ON DELETE SET NULL,
			UNIQUE(product_id, warehouse_id, location_id)
		);`
	} else {
		// SQLite syntax (your original)
		productTable = `
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
		);`

		customerTable = `
		CREATE TABLE IF NOT EXISTS customers (
			id TEXT PRIMARY KEY, 
			name TEXT NOT NULL, 
			email TEXT UNIQUE NOT NULL,
			phone TEXT,
			address TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		orderTable = `
		CREATE TABLE IF NOT EXISTS orders (
			id TEXT PRIMARY KEY,
			customer_id TEXT NOT NULL,
			total_amount REAL NOT NULL,
			status TEXT DEFAULT 'pending',
			order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (customer_id) REFERENCES customers(id)
		);`

		orderItemTable = `
		CREATE TABLE IF NOT EXISTS order_items (
			id TEXT PRIMARY KEY,
			order_id TEXT NOT NULL,
			product_id TEXT NOT NULL,
			quantity INTEGER NOT NULL,
			unit_price REAL NOT NULL,
			total_price REAL NOT NULL,
			FOREIGN KEY (order_id) REFERENCES orders(id),
			FOREIGN KEY (product_id) REFERENCES products(id)
		);`

		supplierTable = `
		CREATE TABLE IF NOT EXISTS suppliers (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			code TEXT UNIQUE NOT NULL,
			contact_person TEXT,
			email TEXT,
			phone TEXT,
			address TEXT, 
			tax_id TEXT,
			payment_terms TEXT,
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		productSupplierTable = `
		CREATE TABLE IF NOT EXISTS product_suppliers (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			supplier_id TEXT NOT NULL,
			supplier_sku TEXT,
			cost_price REAL,
			lead_time_days INTEGER,
			is_primary INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (supplier_id) REFERENCES suppliers(id),
			UNIQUE(product_id, supplier_id)
		);`

		warehouseTable = `
		CREATE TABLE IF NOT EXISTS warehouses (
			id TEXT PRIMARY KEY,
			code TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			location TEXT,
			manager_name TEXT,
			phone TEXT,
			email TEXT,
			capacity INTEGER,
			status TEXT DEFAULT 'active', --active, inactive
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`

		warehouseLocationTable = `
		CREATE TABLE IF NOT EXISTS warehouse_locations (
			id TEXT PRIMARY KEY,
			warehouse_id TEXT NOT NULL,
			location_code TEXT NOT NULL,
			location_name TEXT,
			zone TEXT,
			row_number INTEGER,
			shelf_number INTEGER,
			max_capacity INTEGER,
			current_quantity INTEGER DEFAULT 0,
			status TEXT DEFAULT 'available', -- available, full, inactive, ocupied, reserved, maintanance
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
			UNIQUE(warehouse_id, location_code)
		);`

		inventoryTable = `
		CREATE TABLE IF NOT EXISTS inventory (
			id TEXT PRIMARY KEY,
			product_id TEXT NOT NULL,
			warehouse_id TEXT NOT NULL,
			location_id TEXT,
			quantity INTEGER NOT NULL DEFAULT 0,
			reserved_quantity INTEGER DEFAULT 0,
			min_quantity INTEGER DEFAULT 0,
			max_quantity INTEGER,
			last_restocked TIMESTAMP,
			last_checked TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (product_id) REFERENCES products(id),
			FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
			FOREIGN KEY (location_id) REFERENCES warehouse_locations(id),
			UNIQUE(product_id, warehouse_id, location_id)
		);`
	}

	tables := []string{
		productTable, customerTable, orderTable, orderItemTable,
		supplierTable, productSupplierTable, warehouseTable,
		warehouseLocationTable, inventoryTable,
	}

	for _, table := range tables {
		_, err := DB.Exec(table)
		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("✅ Tables created successfully")
	return nil
}
