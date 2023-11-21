package repository

import (
	"context"
	"database/sql"
	"log"
	"user-admin/internal/domain"
)

type PostgresUserRepository struct {
	DB *sql.DB
}

func (r *PostgresUserRepository) GetAllUsers() (*domain.UsersList, error) {
	rows, err := r.DB.QueryContext(context.TODO(), `
	SELECT id, first_name, last_name, phone_number, blocked, 
	registration_date, gender, date_of_birth, location, 
	email, profile_photo_url
	FROM users`)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		return nil, err
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
