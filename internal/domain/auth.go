package domain

import "errors"

type Admin struct {
	ID       int32  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Roles    []string `json:"roles"`
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

var (
    ErrAdminNotFound = errors.New("admin not found")
)