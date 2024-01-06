package status

import "net/http"

const (
	BadRequest   = http.StatusBadRequest
	Unauthorized = http.StatusUnauthorized
	NotFound     = http.StatusNotFound
	OK           = http.StatusOK
)
