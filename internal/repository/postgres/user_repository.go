package repository

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"
	"strings"
	"user-admin/internal/domain"
	"user-admin/pkg/lib/utils"
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
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(context.TODO())
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
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
			slog.Error("Error scanning user row: %v", utils.Err(err))
			return nil, err
		}

		user.FirstName = utils.HandleNullString(firstName)
		user.LastName = utils.HandleNullString(lastName)
		user.Gender = utils.HandleNullString(gender)
		user.Location = utils.HandleNullString(location)
		user.Email = utils.HandleNullString(email)
		user.ProfilePhotoURL = utils.HandleNullString(profilePhotoURL)

        if dateOfBirth.Valid {
			// Extract year, month, and day from the Date of Birth
			user.DateOfBirth.Year = int32(dateOfBirth.Time.Year())
			user.DateOfBirth.Month = int32(dateOfBirth.Time.Month())
			user.DateOfBirth.Day = int32(dateOfBirth.Time.Day())
		}

		userList.Users = append(userList.Users, user)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Error iterating  over user rows: %v", utils.Err(err))
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
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(context.TODO(), id)

	var user domain.GetUserResponse
        var firstName, lastName,  gender, location, email, profilePhotoURL sql.NullString
        var dateOfBirth sql.NullTime
	err = row.Scan(
		&user.ID, &firstName, &lastName, &user.PhoneNumber, &user.Blocked,
        &user.RegistrationDate, &gender, &dateOfBirth, &location,
        &email, &profilePhotoURL,
	)
	if err != nil {
		slog.Error("error scanning user row: %v", utils.Err(err))
		return nil, err
	}

	user.FirstName = utils.HandleNullString(firstName)
	user.LastName = utils.HandleNullString(lastName)
	user.Gender = utils.HandleNullString(gender)
	user.Location = utils.HandleNullString(location)
	user.Email = utils.HandleNullString(email)
	user.ProfilePhotoURL = utils.HandleNullString(profilePhotoURL)

    if dateOfBirth.Valid {
		// Extract year, month, and day from the Date of Birth
		user.DateOfBirth.Year = int32(dateOfBirth.Time.Year())
		user.DateOfBirth.Month = int32(dateOfBirth.Time.Month())
		user.DateOfBirth.Day = int32(dateOfBirth.Time.Day())
	}

	return &user, nil
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
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	var user domain.CreateUserResponse
	err = stmt.QueryRow(
		utils.NullIfEmptyStr(request.FirstName), utils.NullIfEmptyStr(request.LastName), request.PhoneNumber,
		request.Gender, utils.NullIfEmptyDate(request.DateOfBirth), utils.NullIfEmptyStr(request.Location), // if you are able to show registration date in response, try it
		utils.NullIfEmptyStr(request.Email), utils.NullIfEmptyStr(request.ProfilePhotoURL),
	).Scan(
        &user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, 
        &user.Gender, &user.DateOfBirth, &user.Location,
        &user.Email, &user.ProfilePhotoURL,
    )
	if err != nil {
		slog.Error("error executing query: %v", utils.Err(err))
		return nil, err
	}

	return &user, nil
}

func (r PostgresUserRepository) UpdateUser(request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	updateQuery := "UPDATE users SET"
	var queryParams []interface{}
	var queryArgs []string

	if request.FirstName != "" {
		queryArgs = append(queryArgs, "first_name = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.FirstName)
	}

	if request.LastName != "" {
		queryArgs = append(queryArgs, "last_name = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.LastName)
	}

	if request.PhoneNumber != "" {
		queryArgs = append(queryArgs, "phone_number = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.PhoneNumber)
	}

	if request.Gender != "" {
		queryArgs = append(queryArgs, "gender = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.Gender)
	}

	if request.DateOfBirth.Year != 0 || request.DateOfBirth.Month != 0 || request.DateOfBirth.Day != 0 {
		queryArgs = append(queryArgs, "date_of_birth = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.DateOfBirth)
	}

	if request.Location != "" {
		queryArgs = append(queryArgs, "location = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.Location)
	}

	if request.Email != "" {
		queryArgs = append(queryArgs, "email = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.Email)
	}

	if request.ProfilePhotoURL != "" {
		queryArgs = append(queryArgs, "profile_photo_url = $"+strconv.Itoa(len(queryParams)+1))
		queryParams = append(queryParams, request.ProfilePhotoURL)
	}

	updateQuery += " " + strings.Join(queryArgs, ", ") + " WHERE id = $" + strconv.Itoa(len(queryParams)+1)
	queryParams = append(queryParams, request.ID)

	updateQuery += " RETURNING id, first_name, last_name, phone_number, gender, date_of_birth, location, email, profile_photo_url"

	stmt, err := r.DB.Prepare(updateQuery)
	if err != nil {
		slog.Error("error preparing query: %v", utils.Err(err))
		return nil, err
	}
	defer stmt.Close()

	var user domain.UpdateUserResponse
	err = stmt.QueryRow(queryParams...).Scan(
        &user.ID, &user.FirstName, &user.LastName, &user.PhoneNumber, &user.Gender, &user.DateOfBirth, &user.Location, &user.Email, &user.ProfilePhotoURL,
    )
	if err != nil {
		slog.Error("error executing  query: %v", utils.Err(err))
		return nil, err
	}

	return &user, nil
}

func (r PostgresUserRepository) DeleteUser(id int32) error {
	stmt, err := r.DB.Prepare(`DELETE FROM users WHERE id = $1`)
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

func (r *PostgresUserRepository) BlockUser(id int32) error {
    stmt, err := r.DB.Prepare("UPDATE users SET blocked = true WHERE id = $1")
    if err != nil {
        slog.Error("error preparing query: %v", utils.Err(err))
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

func (r *PostgresUserRepository) UnblockUser(id int32) error {
    stmt, err := r.DB.Prepare("UPDATE users SET blocked = false WHERE id = $1")
    if err != nil {
        slog.Error("error preparing query: %v", utils.Err(err))
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
