package utils

import (
	"database/sql"
	"log/slog"
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