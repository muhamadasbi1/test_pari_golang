package utils

import (
	"net/http"
)

// CekError checks for an error and sends an HTTP error response if there's one
func CekError(w http.ResponseWriter, err error) bool {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
	return false
}
