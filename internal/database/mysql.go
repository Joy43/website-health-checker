package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ssjoy/website-health-checker/internal/config"
)

// DB holds the MySQL database connection pool.
type DB struct {
	*sql.DB
}

// Connect initializes a MySQL database connection using the provided configuration.
func Connect(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %w", err)
	}

	// Set connection pool parameters as requested
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(time.Hour)

	// Verify the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping mysql: %w", err)
	}

	return &DB{db}, nil
}

// CheckHealth verifies if the database connection is still alive.
func (d *DB) CheckHealth() error {
	return d.Ping()
}
