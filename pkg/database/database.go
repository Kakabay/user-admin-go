package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"user-admin/internal/config"
	"user-admin/pkg/lib/utils"

	_ "github.com/lib/pq" // init postgresql driver
)

type Database struct {
	db *sql.DB
}

func InitDB(cfg *config.Config) (*Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBname, cfg.Sslmode)
	
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		slog.Error("failed to initialize database: %v", utils.Err(err))
		return nil, err
	}

	if err := db.Ping(); err != nil {
		slog.Error("failed to initialize database: %v", utils.Err(err))
		return nil, err
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