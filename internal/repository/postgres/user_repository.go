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

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository{
	return &PostgresUserRepository{DB: db}
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
        var firstName, lastName,  gender, location, email, profilePhotoURL sql.NullString
        var dateOfBirth sql.NullTime
		err := rows.Scan(
			&user.ID, &firstName, &lastName, &user.PhoneNumber, &user.Blocked,
			&user.RegistrationDate, &gender, &dateOfBirth, 
			&location, &email, &profilePhotoURL,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			return nil, err
		}

        if dateOfBirth.Valid {
			// Extract year, month, and day from the Date of Birth
			user.DateOfBirth.Year = int32(dateOfBirth.Time.Year())
			user.DateOfBirth.Month = int32(dateOfBirth.Time.Month())
			user.DateOfBirth.Day = int32(dateOfBirth.Time.Day())
		}

        if email.Valid {
            user.Email = email.String
        }

        if profilePhotoURL.Valid {
            user.ProfilePhotoURL = profilePhotoURL.String
        }

		userList.Users = append(userList.Users, user)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating  over user rows: %v", err)
		return nil, err
	}

	return &userList, nil
}

func (r *PostgresUserRepository) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, first_name, last_name, phone_number, blocked, registration_date, gender, date_of_birth, location, email, profile_photo_url
		FROM users 
		WHERE id = $1
	`)

	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(context.TODO(), id)

	var user domain.GetUserResponse
	err = row.Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Blocked,
        &user.RegistrationDate, &user.Gender, &user.DateOfBirth, &user.Location,
        &user.Email, &user.ProfilePhotoURL,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning user row: %v", err)
	}

	return &user, err
}

func (r *PostgresUserRepository) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	stmt, err := r.DB.Prepare(`
		INSERT INTO users (first_name, last_name, phone_number, blocked,
			registration_date, gender, date_of_birth, location,
			email, profile_photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, first_name, last_name, phone_number, blocked,
			registration_date, gender, date_of_birth, location,
			email, profile_photo_url
	`)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()
	
	var user domain.CreateUserResponse
	err = stmt.QueryRow(
		request.FirstName, request.LastName, request.PhoneNumber,
		request.Gender, request.DateOfBirth, request.Location, // if you are able to show registration date in response, try it
        request.Email, request.ProfilePhotoURL,
    ).Scan(
        &user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, 
        &user.Gender, &user.DateOfBirth, &user.Location,
        &user.Email, &user.ProfilePhotoURL,
    )
	if err != nil {
		return nil, fmt.Errorf("error executing query: %v", err)
	}

	return &user, nil
}

func (r PostgresUserRepository) UpdateUser(request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	stmt, err := r.DB.Prepare(`
		UPDATE users
		SET first_name = $2, last_name = $3, phone_number = $4, blocked = $5,
			registration_date = $6, gender = $7, date_of_birth = $8, location = $9,
			email = $10, profile_photo_url = $11
		WHERE id = $1
		RETURNING id, first_name, last_name, phone_number, blocked,
			registration_date, gender, date_of_birth, location,
			email, profile_photo_url
	`)
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	var user domain.UpdateUserResponse
	err = stmt.QueryRow(
        request.ID, request.FirstName, request.LastName, request.PhoneNumber, 
		request.Gender, request.DateOfBirth, request.Location,
        request.Email, request.ProfilePhotoURL,
    ).Scan(
        &user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Blocked,
        &user.RegistrationDate, &user.Gender, &user.DateOfBirth, &user.Location,
        &user.Email, &user.ProfilePhotoURL,
    )
	if err != nil {
		return nil, fmt.Errorf("error executing  query: %v", err)
	}

	return &user, nil
}

func (r PostgresUserRepository) DeleteUser(id int32) error {
	stmt, err := r.DB.Prepare(`DELETE FROM users WHERE id = $1`)
    if err != nil {
        return fmt.Errorf("error preparing query: %v", err)
    }
    defer stmt.Close()
	
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}

	return nil
}

func (r *PostgresUserRepository) BlockUser(id int32) error {
    stmt, err := r.DB.Prepare("UPDATE users SET blocked = true WHERE id = $1")
    if err != nil {
        return fmt.Errorf("error preparing query: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        return fmt.Errorf("error executing query: %v", err)
    }

    return nil
}

func (r *PostgresUserRepository) UnblockUser(id int32) error {
    stmt, err := r.DB.Prepare("UPDATE users SET blocked = false WHERE id = $1")
    if err != nil {
        return fmt.Errorf("error preparing query: %v", err)
    }
    defer stmt.Close()

    _, err = stmt.Exec(id)
    if err != nil {
        return fmt.Errorf("error executing query: %v", err)
    }

    return nil
}
