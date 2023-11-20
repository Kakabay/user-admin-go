package database

import (
	"database/sql"
	"fmt"
	"user-admin/internal/config"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(cfg *config.Config) (*sql.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode)

	var err error 
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping database: %v", err)
	}

	return db, nil
}

func GetDB() *sql.DB {
	return db
}