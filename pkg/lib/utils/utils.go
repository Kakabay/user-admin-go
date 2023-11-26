package utils

import (
	"database/sql"
	"log/slog"
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