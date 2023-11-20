package database

import (
	"database/sql"
	"fmt"
	"user-admin/internal/config"

	_ "github.com/lib/pq" // init postgresql driver
)

type Database struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) (*Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode)
	
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}