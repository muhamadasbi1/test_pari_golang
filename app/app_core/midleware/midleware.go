package routes

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"
)

// Middleware for basic authentication
func BasicAuthMiddleware(db *sql.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the base64-encoded credentials
		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Split username and password
		credentials := strings.SplitN(string(decodedCredentials), ":", 2)
		if len(credentials) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		username, password := credentials[0], credentials[1]

		// Check credentials against the database
		if !checkCredentials(db, username, password) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Call the next handler if credentials are valid
		next.ServeHTTP(w, r)
	}
}

// checkCredentials verifies user credentials against the database
func checkCredentials(db *sql.DB, username, password string) bool {
	var storedPassword string
	query := "SELECT password FROM auth_users WHERE name = ?"
	err := db.QueryRow(query, username).Scan(&storedPassword)

	if err != nil {
		return false // User not found or other error
	}

	return storedPassword == password
}
