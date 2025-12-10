package auth

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"phonax.com/db/db"
	"phonax.com/db/metrics"
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

	start := time.Now()
	// Query the database for the user's hashed password
	// Adjust this query to match your users table schema

	db.CheckConnection(v.db)

	err := v.db.QueryRow(
		"SELECT id, password FROM users WHERE email = $1",
		email,
	).Scan(&userID, &hashedPassword)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			metrics.RecordFailedLogin()
			metrics.RecordQuery("SELECT_VALIDATE_CREDENTIALS", time.Since(start), false)
		} else {
			fmt.Println("ERROR VALIDATE_CREDENTIAL: ", err)
			metrics.RecordQuery("SELECT_VALIDATE_CREDENTIALS", time.Since(start), true)
		}
	} else {
		metrics.RecordQuery("SELECT_VALIDATE_CREDENTIALS", time.Since(start), false)
	}

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
