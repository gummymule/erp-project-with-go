package repositories

import (
	"database/sql"
	"erp-project/models"
	"log"
	"time"
)

type WarehouseRepository struct {
	DB *sql.DB
}

func NewWarehouseRepository(db *sql.DB) *WarehouseRepository {
	return &WarehouseRepository{DB: db}
}

// Warehouse CRUD
func (r *WarehouseRepository) CreateWarehouse(warehouse *models.Warehouse) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `INSERT INTO warehouses (id, code, name, location, manager_name, phone, email, capacity, status, created_at, updated_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	} else {
		query = `INSERT INTO warehouses (id, code, name, location, manager_name, phone, email, capacity, status, created_at, updated_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.DB.Exec(
		query,
		warehouse.ID,
		warehouse.Code,
		warehouse.Name,
		warehouse.Location,
		warehouse.ManagerName,
		warehouse.Phone,
		warehouse.Email,
		warehouse.Capacity,
		warehouse.Status,
		warehouse.CreatedAt,
		warehouse.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating warehouse: %v", err)
		return err
	}
	return nil
}

func (r *WarehouseRepository) GetAllWarehouses() ([]models.Warehouse, error) {
	query := `SELECT id, code, name, location, manager_name, phone, email, capacity, status, created_at, updated_at 
	          FROM warehouses ORDER BY name`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []models.Warehouse
	for rows.Next() {
		var w models.Warehouse
		err := rows.Scan(
			&w.ID,
			&w.Code,
			&w.Name,
			&w.Location,
			&w.ManagerName,
			&w.Phone,
			&w.Email,
			&w.Capacity,
			&w.Status,
			&w.CreatedAt,
			&w.UpdatedAt,
		)
		if err != nil {
			log.Printf("Error scanning warehouse: %v", err)
			continue
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

func (r *WarehouseRepository) GetWarehouseByID(id string) (*models.Warehouse, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT id, code, name, location, manager_name, phone, email, capacity, status, created_at, updated_at 
		         FROM warehouses WHERE id = $1`
	} else {
		query = `SELECT id, code, name, location, manager_name, phone, email, capacity, status, created_at, updated_at 
		         FROM warehouses WHERE id = ?`
	}

	row := r.DB.QueryRow(query, id)

	var w models.Warehouse
	err := row.Scan(
		&w.ID,
		&w.Code,
		&w.Name,
		&w.Location,
		&w.ManagerName,
		&w.Phone,
		&w.Email,
		&w.Capacity,
		&w.Status,
		&w.CreatedAt,
		&w.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		log.Printf("Error getting warehouse: %v", err)
		return nil, err
	}

	return &w, nil
}

func (r *WarehouseRepository) UpdateWarehouse(warehouse *models.Warehouse) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `UPDATE warehouses SET code = $1, name = $2, location = $3, manager_name = $4, 
		         phone = $5, email = $6, capacity = $7, status = $8, updated_at = $9 
		         WHERE id = $10`
	} else {
		query = `UPDATE warehouses SET code = ?, name = ?, location = ?, manager_name = ?, 
		         phone = ?, email = ?, capacity = ?, status = ?, updated_at = ? 
		         WHERE id = ?`
	}

	_, err := r.DB.Exec(
		query,
		warehouse.Code,
		warehouse.Name,
		warehouse.Location,
		warehouse.ManagerName,
		warehouse.Phone,
		warehouse.Email,
		warehouse.Capacity,
		warehouse.Status,
		time.Now(),
		warehouse.ID,
	)

	if err != nil {
		log.Printf("Error updating warehouse: %v", err)
		return err
	}
	return nil
}

func (r *WarehouseRepository) DeleteWarehouse(id string) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `DELETE FROM warehouses WHERE id = $1`
	} else {
		query = `DELETE FROM warehouses WHERE id = ?`
	}

	_, err := r.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting warehouse: %v", err)
		return err
	}
	return nil
}

// WarehouseLocation CRUD
func (r *WarehouseRepository) CreateLocation(location *models.WarehouseLocation) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `INSERT INTO warehouse_locations (id, warehouse_id, location_code, location_name, zone, 
		         row_number, shelf_number, max_capacity, current_quantity, status, created_at) 
		         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	} else {
		query = `INSERT INTO warehouse_locations (id, warehouse_id, location_code, location_name, zone, 
		         row_number, shelf_number, max_capacity, current_quantity, status, created_at) 
		         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := r.DB.Exec(
		query,
		location.ID,
		location.WarehouseID,
		location.LocationCode,
		location.LocationName,
		location.Zone,
		location.RowNumber,
		location.ShelfNumber,
		location.MaxCapacity,
		location.CurrentQuantity,
		location.Status,
		location.CreatedAt,
	)

	if err != nil {
		log.Printf("Error creating warehouse location: %v", err)
		return err
	}
	return nil
}

func (r *WarehouseRepository) GetLocationsByWarehouse(warehouseID string) ([]models.WarehouseLocation, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT wl.id, wl.warehouse_id, wl.location_code, wl.location_name, wl.zone, 
		                wl.row_number, wl.shelf_number, wl.max_capacity, wl.current_quantity, 
		                wl.status, wl.created_at, w.name as warehouse_name
		         FROM warehouse_locations wl
		         JOIN warehouses w ON wl.warehouse_id = w.id
		         WHERE wl.warehouse_id = $1
		         ORDER BY wl.zone, wl.row_number, wl.shelf_number`
	} else {
		query = `SELECT wl.id, wl.warehouse_id, wl.location_code, wl.location_name, wl.zone, 
		                wl.row_number, wl.shelf_number, wl.max_capacity, wl.current_quantity, 
		                wl.status, wl.created_at, w.name as warehouse_name
		         FROM warehouse_locations wl
		         JOIN warehouses w ON wl.warehouse_id = w.id
		         WHERE wl.warehouse_id = ?
		         ORDER BY wl.zone, wl.row_number, wl.shelf_number`
	}

	rows, err := r.DB.Query(query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.WarehouseLocation
	for rows.Next() {
		var loc models.WarehouseLocation
		err := rows.Scan(
			&loc.ID,
			&loc.WarehouseID,
			&loc.LocationCode,
			&loc.LocationName,
			&loc.Zone,
			&loc.RowNumber,
			&loc.ShelfNumber,
			&loc.MaxCapacity,
			&loc.CurrentQuantity,
			&loc.Status,
			&loc.CreatedAt,
			&loc.WarehouseName,
		)
		if err != nil {
			log.Printf("Error scanning warehouse location: %v", err)
			continue
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

func (r *WarehouseRepository) GetAvailableLocations(warehouseID string) ([]models.WarehouseLocation, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT wl.id, wl.warehouse_id, wl.location_code, wl.location_name, wl.zone, 
		                wl.row_number, wl.shelf_number, wl.max_capacity, wl.current_quantity, 
		                wl.status, wl.created_at, w.name as warehouse_name
		         FROM warehouse_locations wl
		         JOIN warehouses w ON wl.warehouse_id = w.id
		         WHERE wl.warehouse_id = $1 AND wl.status = 'available' 
		               AND (wl.max_capacity IS NULL OR wl.current_quantity < wl.max_capacity)
		         ORDER BY wl.zone, wl.row_number, wl.shelf_number`
	} else {
		query = `SELECT wl.id, wl.warehouse_id, wl.location_code, wl.location_name, wl.zone, 
		                wl.row_number, wl.shelf_number, wl.max_capacity, wl.current_quantity, 
		                wl.status, wl.created_at, w.name as warehouse_name
		         FROM warehouse_locations wl
		         JOIN warehouses w ON wl.warehouse_id = w.id
		         WHERE wl.warehouse_id = ? AND wl.status = 'available' 
		               AND (wl.max_capacity IS NULL OR wl.current_quantity < wl.max_capacity)
		         ORDER BY wl.zone, wl.row_number, wl.shelf_number`
	}

	rows, err := r.DB.Query(query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.WarehouseLocation
	for rows.Next() {
		var loc models.WarehouseLocation
		err := rows.Scan(
			&loc.ID,
			&loc.WarehouseID,
			&loc.LocationCode,
			&loc.LocationName,
			&loc.Zone,
			&loc.RowNumber,
			&loc.ShelfNumber,
			&loc.MaxCapacity,
			&loc.CurrentQuantity,
			&loc.Status,
			&loc.CreatedAt,
			&loc.WarehouseName,
		)
		if err != nil {
			log.Printf("Error scanning available location: %v", err)
			continue
		}
		locations = append(locations, loc)
	}

	return locations, nil
}

func (r *WarehouseRepository) UpdateLocation(location *models.WarehouseLocation) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `UPDATE warehouse_locations SET location_code = $1, location_name = $2, zone = $3, 
		         row_number = $4, shelf_number = $5, max_capacity = $6, current_quantity = $7, 
		         status = $8 WHERE id = $9`
	} else {
		query = `UPDATE warehouse_locations SET location_code = ?, location_name = ?, zone = ?, 
		         row_number = ?, shelf_number = ?, max_capacity = ?, current_quantity = ?, 
		         status = ? WHERE id = ?`
	}

	_, err := r.DB.Exec(
		query,
		location.LocationCode,
		location.LocationName,
		location.Zone,
		location.RowNumber,
		location.ShelfNumber,
		location.MaxCapacity,
		location.CurrentQuantity,
		location.Status,
		location.ID,
	)

	if err != nil {
		log.Printf("Error updating warehouse location: %v", err)
		return err
	}
	return nil
}

// Inventory CRUD
func (r *WarehouseRepository) CreateInventory(inventory *models.Inventory) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	var err error

	if inventory.LocationID != nil {
		if isPostgreSQL {
			query = `INSERT INTO inventory (id, product_id, warehouse_id, location_id, quantity, 
			         reserved_quantity, min_quantity, max_quantity, last_restocked, last_checked, 
			         created_at, updated_at) 
			         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
			_, err = r.DB.Exec(
				query,
				inventory.ID,
				inventory.ProductID,
				inventory.WarehouseID,
				inventory.LocationID,
				inventory.Quantity,
				inventory.ReservedQuantity,
				inventory.MinQuantity,
				inventory.MaxQuantity,
				inventory.LastRestocked,
				inventory.LastChecked,
				inventory.CreatedAt,
				inventory.UpdatedAt,
			)
		} else {
			query = `INSERT INTO inventory (id, product_id, warehouse_id, location_id, quantity, 
			         reserved_quantity, min_quantity, max_quantity, last_restocked, last_checked, 
			         created_at, updated_at) 
			         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
			_, err = r.DB.Exec(
				query,
				inventory.ID,
				inventory.ProductID,
				inventory.WarehouseID,
				inventory.LocationID,
				inventory.Quantity,
				inventory.ReservedQuantity,
				inventory.MinQuantity,
				inventory.MaxQuantity,
				inventory.LastRestocked,
				inventory.LastChecked,
				inventory.CreatedAt,
				inventory.UpdatedAt,
			)
		}
	} else {
		if isPostgreSQL {
			query = `INSERT INTO inventory (id, product_id, warehouse_id, quantity, 
			         reserved_quantity, min_quantity, max_quantity, last_restocked, last_checked, 
			         created_at, updated_at) 
			         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
			_, err = r.DB.Exec(
				query,
				inventory.ID,
				inventory.ProductID,
				inventory.WarehouseID,
				inventory.Quantity,
				inventory.ReservedQuantity,
				inventory.MinQuantity,
				inventory.MaxQuantity,
				inventory.LastRestocked,
				inventory.LastChecked,
				inventory.CreatedAt,
				inventory.UpdatedAt,
			)
		} else {
			query = `INSERT INTO inventory (id, product_id, warehouse_id, quantity, 
			         reserved_quantity, min_quantity, max_quantity, last_restocked, last_checked, 
			         created_at, updated_at) 
			         VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
			_, err = r.DB.Exec(
				query,
				inventory.ID,
				inventory.ProductID,
				inventory.WarehouseID,
				inventory.Quantity,
				inventory.ReservedQuantity,
				inventory.MinQuantity,
				inventory.MaxQuantity,
				inventory.LastRestocked,
				inventory.LastChecked,
				inventory.CreatedAt,
				inventory.UpdatedAt,
			)
		}
	}

	if err != nil {
		log.Printf("Error creating inventory: %v", err)
		return err
	}
	return nil
}

func (r *WarehouseRepository) GetInventoryByProduct(productID string) ([]models.Inventory, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT i.id, i.product_id, i.warehouse_id, i.location_id, i.quantity, 
		                i.reserved_quantity, i.min_quantity, i.max_quantity, i.last_restocked, 
		                i.last_checked, i.created_at, i.updated_at,
		                p.name as product_name, p.sku, w.name as warehouse_name, 
		                wl.location_code
		         FROM inventory i
		         JOIN products p ON i.product_id = p.id
		         JOIN warehouses w ON i.warehouse_id = w.id
		         LEFT JOIN warehouse_locations wl ON i.location_id = wl.id
		         WHERE i.product_id = $1
		         ORDER BY w.name`
	} else {
		query = `SELECT i.id, i.product_id, i.warehouse_id, i.location_id, i.quantity, 
		                i.reserved_quantity, i.min_quantity, i.max_quantity, i.last_restocked, 
		                i.last_checked, i.created_at, i.updated_at,
		                p.name as product_name, p.sku, w.name as warehouse_name, 
		                wl.location_code
		         FROM inventory i
		         JOIN products p ON i.product_id = p.id
		         JOIN warehouses w ON i.warehouse_id = w.id
		         LEFT JOIN warehouse_locations wl ON i.location_id = wl.id
		         WHERE i.product_id = ?
		         ORDER BY w.name`
	}

	rows, err := r.DB.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inv models.Inventory
		var locationID sql.NullString
		var locationCode sql.NullString
		var maxQuantity sql.NullInt64
		var lastRestocked, lastChecked sql.NullTime
		var sku string

		err := rows.Scan(
			&inv.ID,
			&inv.ProductID,
			&inv.WarehouseID,
			&locationID,
			&inv.Quantity,
			&inv.ReservedQuantity,
			&inv.MinQuantity,
			&maxQuantity,
			&lastRestocked,
			&lastChecked,
			&inv.CreatedAt,
			&inv.UpdatedAt,
			&inv.ProductName,
			&sku,
			&inv.WarehouseName,
			&locationCode,
		)

		if err != nil {
			log.Printf("Error scanning inventory: %v", err)
			continue
		}

		// Handle nullable fields
		if locationID.Valid {
			locID := locationID.String
			inv.LocationID = &locID
		}
		if locationCode.Valid {
			inv.LocationCode = locationCode.String
		}
		if maxQuantity.Valid {
			maxQty := int(maxQuantity.Int64)
			inv.MaxQuantity = &maxQty
		}
		if lastRestocked.Valid {
			inv.LastRestocked = &lastRestocked.Time
		}
		if lastChecked.Valid {
			inv.LastChecked = &lastChecked.Time
		}

		inv.SKU = sku
		inv.AvailableQuantity = inv.Quantity - inv.ReservedQuantity

		inventories = append(inventories, inv)
	}

	return inventories, nil
}

func (r *WarehouseRepository) GetInventoryByWarehouse(warehouseID string) ([]models.Inventory, error) {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `SELECT i.id, i.product_id, i.warehouse_id, i.location_id, i.quantity, 
		                i.reserved_quantity, i.min_quantity, i.max_quantity, i.last_restocked, 
		                i.last_checked, i.created_at, i.updated_at,
		                p.name as product_name, p.sku, w.name as warehouse_name, 
		                wl.location_code
		         FROM inventory i
		         JOIN products p ON i.product_id = p.id
		         JOIN warehouses w ON i.warehouse_id = w.id
		         LEFT JOIN warehouse_locations wl ON i.location_id = wl.id
		         WHERE i.warehouse_id = $1
		         ORDER BY p.name`
	} else {
		query = `SELECT i.id, i.product_id, i.warehouse_id, i.location_id, i.quantity, 
		                i.reserved_quantity, i.min_quantity, i.max_quantity, i.last_restocked, 
		                i.last_checked, i.created_at, i.updated_at,
		                p.name as product_name, p.sku, w.name as warehouse_name, 
		                wl.location_code
		         FROM inventory i
		         JOIN products p ON i.product_id = p.id
		         JOIN warehouses w ON i.warehouse_id = w.id
		         LEFT JOIN warehouse_locations wl ON i.location_id = wl.id
		         WHERE i.warehouse_id = ?
		         ORDER BY p.name`
	}

	rows, err := r.DB.Query(query, warehouseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var inventories []models.Inventory
	for rows.Next() {
		var inv models.Inventory
		var locationID sql.NullString
		var locationCode sql.NullString
		var maxQuantity sql.NullInt64
		var lastRestocked, lastChecked sql.NullTime
		var sku string

		err := rows.Scan(
			&inv.ID,
			&inv.ProductID,
			&inv.WarehouseID,
			&locationID,
			&inv.Quantity,
			&inv.ReservedQuantity,
			&inv.MinQuantity,
			&maxQuantity,
			&lastRestocked,
			&lastChecked,
			&inv.CreatedAt,
			&inv.UpdatedAt,
			&inv.ProductName,
			&sku,
			&inv.WarehouseName,
			&locationCode,
		)

		if err != nil {
			log.Printf("Error scanning inventory: %v", err)
			continue
		}

		// Handle nullable fields
		if locationID.Valid {
			locID := locationID.String
			inv.LocationID = &locID
		}
		if locationCode.Valid {
			inv.LocationCode = locationCode.String
		}
		if maxQuantity.Valid {
			maxQty := int(maxQuantity.Int64)
			inv.MaxQuantity = &maxQty
		}
		if lastRestocked.Valid {
			inv.LastRestocked = &lastRestocked.Time
		}
		if lastChecked.Valid {
			inv.LastChecked = &lastChecked.Time
		}

		inv.SKU = sku
		inv.AvailableQuantity = inv.Quantity - inv.ReservedQuantity

		inventories = append(inventories, inv)
	}

	return inventories, nil
}

func (r *WarehouseRepository) UpdateInventoryQuantity(id string, quantity, reservedQuantity int) error {
	isPostgreSQL := isPostgreSQL(r.DB)

	var query string
	if isPostgreSQL {
		query = `UPDATE inventory SET quantity = $1, reserved_quantity = $2, updated_at = $3 WHERE id = $4`
	} else {
		query = `UPDATE inventory SET quantity = ?, reserved_quantity = ?, updated_at = ? WHERE id = ?`
	}

	_, err := r.DB.Exec(
		query,
		quantity,
		reservedQuantity,
		time.Now(),
		id,
	)

	if err != nil {
		log.Printf("Error updating inventory quantity: %v", err)
		return err
	}
	return nil
}
