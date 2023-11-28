package utils

import (
	"database/sql"
	"log/slog"
	"strings"
	"user-admin/internal/domain"
)

// Use slog to handle errors
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func HandleNullString(valid sql.NullString) string {
	if valid.Valid {
		return valid.String
	}
	return ""
}

func NullIfEmptyStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func NullIfEmptyDate(d domain.Date) sql.NullTime {
	if d == (domain.Date{}) {
		return sql.NullTime{}
	}
	return sql.NullTime{Valid: true}
}

func IsValidPhoneNumber(phoneNumber string) bool {
	// Check if the phone number consists of 12 digits and starts with "+993"
	const validPrefix = "+993"
	return len(phoneNumber) == 12 && strings.HasPrefix(phoneNumber, validPrefix)
}