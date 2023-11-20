package user

import "time"

type DateOfBirth struct {
	Day int32 `json:"day"`
	Month int32 `json:"month"`
	Year int32 `json:"year"`
}

type Empty struct{}

type UserID struct {
	ID int32 `json:"id"`
}

// ADD PAGINATION(!)
type UsersList struct {
	Users        []GetUserResponse `json:"users"`
}

type GetUserResponse struct {
	ID               int32       `json:"id"`
	FirstName        string      `json:"first_name"`
	LastName         string      `json:"last_name"`
	PhoneNumber      string      `json:"phone_number"`
	Blocked          bool        `json:"blocked"`
	Gender           string      `json:"gender"`
	RegistrationDate time.Time   `json:"registration_date"`
	DateOfBirth      DateOfBirth `json:"date_of_birth"`
	Location         string      `json:"location"`
	Email            string      `json:"email"`
	ProfilePhotoURL  string      `json:"profile_photo_url"`
}

type CreateUserRequest struct {
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	PhoneNumber     string      `json:"phone_number"`
	Gender          string      `json:"gender"`
	DateOfBirth     DateOfBirth `json:"date_of_birth"`
	Location        string      `json:"location"`
	Email           string      `json:"email"`
	ProfilePhotoURL string      `json:"profile_photo_url"`
}

type CreateUserResponse struct {
	ID               int32           `json:"id"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	PhoneNumber      string          `json:"phone_number"`
	Blocked          bool            `json:"blocked"`
	Gender           string          `json:"gender"`
	RegistrationDate time.Time `json:"registration_date"`
	DateOfBirth      DateOfBirth     `json:"date_of_birth"`
	Location         string          `json:"location"`
	Email            string          `json:"email"`
	ProfilePhotoURL  string          `json:"profile_photo_url"`
}

type UpdateUserRequest struct {
	ID              int32       `json:"id"`
	FirstName       string      `json:"first_name"`
	LastName        string      `json:"last_name"`
	PhoneNumber     string      `json:"phone_number"`
	Gender          string      `json:"gender"`
	DateOfBirth     DateOfBirth `json:"date_of_birth"`
	Location        string      `json:"location"`
	Email           string      `json:"email"`
	ProfilePhotoURL string      `json:"profile_photo_url"`
}

type UpdateUserResponse struct {
	ID               int32           `json:"id"`
	FirstName        string          `json:"first_name"`
	LastName         string          `json:"last_name"`
	PhoneNumber      string          `json:"phone_number"`
	Blocked          bool            `json:"blocked"`
	Gender           string          `json:"gender"`
	RegistrationDate time.Time `json:"registration_date"`
	DateOfBirth      DateOfBirth     `json:"date_of_birth"`
	Location         string          `json:"location"`
	Email            string          `json:"email"`
	ProfilePhotoURL  string          `json:"profile_photo_url"`
}