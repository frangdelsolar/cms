package builder

import (
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrDBNotInitialized    = errors.New("database not initialized")
	ErrDBConfigNotProvided = errors.New("database config not provided")
)

// Database represents a database connection managed by GORM.
type Database struct {
	DB *gorm.DB // Embedded GORM DB instance for database access
}

// FindById retrieves a single record from the database that matches the provided ID.
// It allows for an optional query extension to refine the search criteria.
//
// Parameters:
//   - id: the unique identifier of the record to be retrieved.
//   - entity: the destination where the result will be stored.
//   - queryExtension: an optional additional query condition.
//
// Returns:
//   - *gorm.DB: the result of the database query, which can be used to check for errors.
func (db *Database) FindById(id string, entity interface{}, queryExtension string) *gorm.DB {
	q := "id = '" + id + "'"

	if queryExtension != "" {
		q += " AND " + queryExtension
	}

	return db.DB.Where(q).First(entity)
}

// Find retrieves records from the database based on the provided query.
// If pagination is provided, the query will be limited to the specified number of records
// and offset to the correct page.
//
// Parameters:
//   - entity: the destination where the result will be stored.
//   - query: the query to be executed, it can be a raw SQL query or a GORM query.
//   - pagination: optional pagination information.
//
// Returns:
//   - *gorm.DB: the result of the database query, which can be used to check for errors.
func (db *Database) Find(entity interface{}, query string, pagination *Pagination) *gorm.DB {

	if pagination == nil {
		return db.DB.Where(query).Find(entity)
	}

	// Retrieve total number of records
	db.DB.Model(entity).Where(query).Count(&pagination.Total)

	// Apply pagination
	filtered := db.DB.Where(query)
	limit := pagination.Limit
	offset := (pagination.Page - 1) * pagination.Limit

	return filtered.Limit(limit).Offset(offset).Find(entity)
}

// Create creates a new record in the database.
//
// Parameters:
//   - entity: the model instance to be created.
//
// Returns:
//   - *gorm.DB: the result of the database query, which can be used to check for errors.
func (db *Database) Create(entity interface{}) *gorm.DB {
	return db.DB.Create(entity)
}

// Delete deletes the record in the database.
//
// Parameters:
//   - entity: the model instance to be deleted.
//
// Returns:
//   - *gorm.DB: the result of the database query, which can be used to check for errors.
func (db *Database) Delete(entity interface{}) *gorm.DB {
	return db.DB.Delete(entity)
}

// Save updates a record in the database if it already exists, or creates a new one if it does not.
//
// Parameters:
//   - entity: the model instance to be saved.
//
// Returns:
//   - *gorm.DB: the result of the database query, which can be used to check for errors.
func (db *Database) Save(entity interface{}) *gorm.DB {
	return db.DB.Save(entity)
}

// DBConfig defines the configuration options for connecting to a database.
type DBConfig struct {
	// URL: Used for connecting to a PostgreSQL database.
	// Provide a complete connection string (e.g., "postgres://user:password@host:port/database").
	URL string
	// Path: Used for connecting to a SQLite database.
	// Provide the path to the SQLite database file.
	Path string
}

// LoadDB establishes a connection to the database based on the provided configuration.
//
// It takes a pointer to a DBConfig struct as input, which specifies the connection details.
// On successful connection, it returns a pointer to a Database instance encapsulating the GORM DB object.
// Otherwise, it returns an error indicating the connection failure.
func LoadDB(config *DBConfig) (*Database, error) {

	if config == nil || (config.URL == "" && config.Path == "") {
		return nil, ErrDBConfigNotProvided
	}

	var db *Database

	if config.URL != "" {
		// Connect to PostgreSQL
		gormDB, err := gorm.Open(postgres.Open(config.URL), &gorm.Config{})
		if err != nil {
			return db, err
		}
		return &Database{
			gormDB,
		}, nil
	}

	if config.Path != "" {
		// Connect to SQLite
		gormDB, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{})
		if err != nil {
			return db, err
		}
		return &Database{
			gormDB,
		}, nil
	}

	return db, ErrDBConfigNotProvided // Should never be reached, but added for completeness
}

// Migrate calls the AutoMigrate method on the GORM DB instance.
func (db *Database) Migrate(model interface{}) error {
	if db == nil {
		return ErrDBNotInitialized
	}
	db.DB.AutoMigrate(model)
	return nil
}
