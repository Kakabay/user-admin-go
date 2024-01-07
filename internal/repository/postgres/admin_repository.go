package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"user-admin/internal/domain"
	"user-admin/pkg/lib/utils"

	"golang.org/x/crypto/bcrypt"
)

type PostgresAdminRepository struct {
	DB *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{DB: db}
}

// TODO: GET ALL ADMINS
// TODO: SEARCH ADMINS
// TODO: UPDATE ADMINS

func (r *PostgresAdminRepository) GetAdminByID(id int32) (*domain.CommonAdminResponse, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, username, role
		FROM admins
		WHERE id = $1
	`)
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(context.TODO(), id)

	var admin domain.CommonAdminResponse

	err = row.Scan(
		&admin.ID,
		&admin.Username,
		&admin.Role,
	)
	if err != nil {
		slog.Error("error scanning admin row: %v", utils.Err(err))
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminRepository) CreateAdmin(request *domain.CreateAdminRequest) (*domain.CommonAdminResponse, error) {
	if request.Username == "" || request.Password == "" || request.Role == "" {
		return nil, fmt.Errorf("username, password, and role are required fields")
	}

	var existingUsername string
	err := r.DB.QueryRow("SELECT username FROM admins WHERE username = $1 LIMIT 1", request.Username).Scan(&existingUsername)
	if err == sql.ErrNoRows {
	} else if err != nil {
		slog.Error("error checking admin existence: %v", utils.Err(err))
		return nil, err
	} else {
		return nil, domain.ErrAdminAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("error hashing password: %v", utils.Err(err))
		return nil, err
	}

	stmt, err := r.DB.Prepare(`
		INSERT INTO admins (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, username, password, role
	`)
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	var admin domain.CommonAdminResponse

	err = stmt.QueryRow(
		request.Username,
		hashedPassword,
		request.Role,
	).Scan(
		&admin.ID,
		&admin.Username,
		&hashedPassword,
		&admin.Role,
	)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return nil, err
	}

	return &admin, nil
}

/*
func (r *PostgresAdminRepository) DeleteAdmin(id int) error {
	var exists bool
	err := r.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM admins WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		slog.Error("error checking admin existence: %v", utils.Err(err))
		return err
	}

	if !exists {
		return fmt.Errorf("admin with ID %d not found", id)
	}

	stmt, err := r.DB.Prepare(`DELETE FROM admins WHERE id = $1`)
	if err != nil {
		slog.Error("error preparing query: %v", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return err
	}

	return nil
}
*/
