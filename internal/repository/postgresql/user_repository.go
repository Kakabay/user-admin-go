package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"user-admin/internal/domain"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) GetAllUsers() (*domain.UsersList, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, first_name, last_name, phone_number, blocked, 
		registration_date, gender, date_of_birth, location, 
		email, profile_photo_url
		FROM users`)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var userList domain.UsersList
	for rows.Next() {
		var user domain.CommonUserResponse
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Blocked,
			&user.RegistrationDate, &user.Gender, &user.DateOfBirth, &user.Location,
			&user.Email, &user.ProfilePhotoURL,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, err
		}
		userList.Users = append(userList.Users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating  over user rows: %v", err)
		return nil, err
	}

	return &userList, nil
}
