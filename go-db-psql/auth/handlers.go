package auth

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// LoginRequest represents a login request payload
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents a successful login response
type LoginResponse struct {
	Token string `json:"token"`
}

// CredentialValidator is an interface for validating user credentials
// Implement this with your database logic
type CredentialValidator interface {
	ValidateCredentials(email string, password string) (UserID int, valid bool, err error)
}

var credentialValidator CredentialValidator

// SetCredentialValidator sets the credential validator for the login handler
func SetCredentialValidator(validator CredentialValidator) {
	credentialValidator = validator
}

// LoginHandler handles POST /auth/login requests
// Validates email and password against the credential validator.
// If no validator is set, it uses a demo mode that accepts any email/password.
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	var UserID int
	var valid bool

	// If a credential validator is set, use it; otherwise use demo mode
	if credentialValidator != nil {
		var err error
		UserID, valid, err = credentialValidator.ValidateCredentials(req.Email, req.Password)
		if err != nil {
			http.Error(w, "credential validation error", http.StatusInternalServerError)
			return
		}
	} else {
		// Demo mode: accept any email/password
		UserID = 1
		valid = true
	}

	if !valid {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate token with a 24-hour expiration
	token, err := GenerateToken(UserID, req.Email, 24)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// HashPassword hashes a password using bcrypt (utility for credential validators)
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// VerifyPassword checks if a password matches a bcrypt hash
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
