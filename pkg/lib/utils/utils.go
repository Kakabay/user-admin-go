package utils

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
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

func RespondWithError(w http.ResponseWriter, code int, message string) {
    w.WriteHeader(code)
    w.Write([]byte(message))
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        slog.Error("Error encoding JSON: ", Err(err))
        RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
    }
}