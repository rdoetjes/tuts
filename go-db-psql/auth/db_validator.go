package auth

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// DBCredentialValidator validates credentials against a database
type DBCredentialValidator struct {
	db *sqlx.DB
}

// NewDBCredentialValidator creates a new database credential validator
func NewDBCredentialValidator(db *sqlx.DB) *DBCredentialValidator {
	return &DBCredentialValidator{db: db}
}

// ValidateCredentials checks if email and password match a user in the database
// Returns (userID, isValid, error)
func (v *DBCredentialValidator) ValidateCredentials(email, password string) (int, bool, error) {
	var userID int
	var hashedPassword string

	// Query the database for the user's hashed password
	// Adjust this query to match your users table schema
	err := v.db.QueryRow(
		"SELECT id, password FROM users WHERE email = $1",
		email,
	).Scan(&userID, &hashedPassword)

	if err == sql.ErrNoRows {
		// User not found
		return 0, false, nil
	}

	if err != nil {
		// Database error
		return 0, false, err
	}

	// Verify the password
	isValid := VerifyPassword(hashedPassword, password)
	return userID, isValid, nil
}
