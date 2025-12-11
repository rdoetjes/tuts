package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func testServersBaseURL() string {
	return "http://localhost:3000/v1"
}

func authServerBaseURL() string {
	return "http://localhost:3000"
}

// Helper function to get an authenticated token
func loginAsAdmin(t *testing.T, client *http.Client, base string) string {
	loginReq := map[string]string{
		"email":    "test@localhost",
		"password": "Testing123",
	}
	loginBuf, _ := json.Marshal(loginReq)
	resp, err := client.Post(base+"/auth/login", "application/json", bytes.NewReader(loginBuf))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on login, got %d: %s", resp.StatusCode, string(b))
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("invalid token in login response: %v", loginResp)
	}

	return token
}

// Helper function to login with email/password
func login(t *testing.T, client *http.Client, base, email, password string) (string, bool) {
	loginReq := map[string]string{
		"email":    email,
		"password": password,
	}
	loginBuf, _ := json.Marshal(loginReq)
	resp, err := client.Post(base+"/auth/login", "application/json", bytes.NewReader(loginBuf))
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return "", false // login failed as expected
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on login, got %d: %s", resp.StatusCode, string(b))
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("invalid token in login response: %v", loginResp)
	}

	return token, true
}

// Helper function to create a user
func createUser(t *testing.T, client *http.Client, base, token string, userData map[string]interface{}) string {
	buf, _ := json.Marshal(userData)
	req, _ := http.NewRequest(http.MethodPost, base+"/users", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("create request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 201 Created, got %d: %s", resp.StatusCode, string(b))
	}

	var created map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}

	idVal := created["id"]
	if idVal == nil {
		t.Fatalf("invalid id from create: %v", created["id"])
	}

	return fmt.Sprintf("%v", idVal)
}

// Helper function to update a user
func updateUser(t *testing.T, client *http.Client, base, token, userID string, updateData map[string]interface{}) {
	buf, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(http.MethodPut, base+"/users/"+userID, bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on update, got %d: %s", resp.StatusCode, string(b))
	}
}

// Helper function to get a user
func getUser(t *testing.T, client *http.Client, base, token, userID string) map[string]interface{} {
	req, _ := http.NewRequest(http.MethodGet, base+"/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on get, got %d: %s", resp.StatusCode, string(b))
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	return user
}

// Helper function to delete a user
func deleteUser(t *testing.T, client *http.Client, base, token, userID string) {
	req, _ := http.NewRequest(http.MethodDelete, base+"/users/"+userID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 or 204 on delete, got %d: %s", resp.StatusCode, string(b))
	}
}

// TestCreateUserWithPassword tests creating a user and logging in with the password
func TestCreateUserWithPassword(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServersBaseURL()

	// Login as admin
	token := loginAsAdmin(t, client, authServerBaseURL())

	// Create user
	userID := createUser(t, client, base, token, map[string]interface{}{
		"firstname": "Create",
		"lastname":  "User",
		"dob":       "1990-01-01",
		"email":     "create@email.com",
		"password":  "CreatePassword123",
	})

	// Login with newly created user
	newToken, ok := login(t, client, authServerBaseURL(), "create@email.com", "CreatePassword123")
	if !ok || newToken == "" {
		t.Fatal("failed to login with newly created user")
	}

	// Cleanup
	deleteUser(t, client, base, token, userID)
}

// TestUpdateUserPassword tests updating a user's password and logging in with the new password
func TestUpdateUserPassword(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServersBaseURL()

	// Login as admin
	token := loginAsAdmin(t, client, authServerBaseURL())

	// Create user
	userID := createUser(t, client, base, token, map[string]interface{}{
		"firstname": "Update",
		"lastname":  "User",
		"dob":       "1990-01-01",
		"email":     "update@email.com",
		"password":  "OldPassword123",
	})

	// Update user password
	updateUser(t, client, base, token, userID, map[string]interface{}{
		"firstname": "UpdatedFirst",
		"lastname":  "UpdatedLast",
		"email":     "updated@email.com",
		"dob":       "1973-12-12",
		"password":  "NewPassword456",
	})

	// Login with new password
	newToken, ok := login(t, client, authServerBaseURL(), "updated@email.com", "NewPassword456")
	if !ok || newToken == "" {
		t.Fatal("failed to login with updated password")
	}

	// Cleanup
	deleteUser(t, client, base, token, userID)
}

// TestLoginWithWrongPassword tests that login fails with wrong password
func TestLoginWithWrongPassword(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServersBaseURL()

	// Login as admin
	token := loginAsAdmin(t, client, authServerBaseURL())

	// Create user
	userID := createUser(t, client, base, token, map[string]interface{}{
		"firstname": "Wrong",
		"lastname":  "Pass",
		"dob":       "1990-01-01",
		"email":     "wrongpass@email.com",
		"password":  "CorrectPassword123",
	})

	// Attempt login with wrong password
	_, ok := login(t, client, authServerBaseURL(), "wrongpass@email.com", "WrongPassword123")
	if ok {
		t.Fatal("login should have failed with wrong password")
	}

	// Cleanup
	deleteUser(t, client, base, token, userID)
}

// TestCompleteUserCRUD tests the full user lifecycle: create, read, update, delete
func TestCompleteUserCRUD(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServersBaseURL()

	// Login as admin
	token := loginAsAdmin(t, client, authServerBaseURL())

	// Create user
	userID := createUser(t, client, base, token, map[string]interface{}{
		"firstname": "Complete",
		"lastname":  "Test",
		"dob":       "1990-01-01",
		"email":     "complete@email.com",
		"password":  "CompletePass123",
	})
	defer deleteUser(t, client, base, token, userID)

	// Login with new user
	_, ok := login(t, client, authServerBaseURL(), "complete@email.com", "CompletePass123")
	if !ok {
		t.Fatal("failed to login with newly created user")
	}

	// Update user
	updateUser(t, client, base, token, userID, map[string]interface{}{
		"firstname": "CompletedFirst",
		"lastname":  "CompletedLast",
		"email":     "completed@email.com",
		"dob":       "1973-12-12",
		"password":  "UpdatedPass456",
	})

	// Verify old password no longer works
	_, ok = login(t, client, authServerBaseURL(), "completed@email.com", "CompletePass123")
	if ok {
		t.Fatal("old password should not work after update")
	}

	// Verify new password works
	_, ok = login(t, client, authServerBaseURL(), "completed@email.com", "UpdatedPass456")
	if !ok {
		t.Fatal("failed to login with updated password")
	}

	// Get user and verify fields
	user := getUser(t, client, base, token, userID)
	if user["FirstName"] != "CompletedFirst" {
		t.Fatalf("expected FirstName 'CompletedFirst', got %v", user["FirstName"])
	}
	if user["LastName"] != "CompletedLast" {
		t.Fatalf("expected LastName 'CompletedLast', got %v", user["LastName"])
	}
	if user["Email"] != "completed@email.com" {
		t.Fatalf("expected Email 'completed@email.com', got %v", user["Email"])
	}
	if user["DOB"] != "1973-12-12T00:00:00Z" {
		t.Fatalf("expected DOB 1973-12-12T00:00:00Z, got %v", user["DOB"])
	}
	var p = fmt.Sprintf("%v", user["Password"])
	if strings.Contains(p, "$2a$10$") == false {
		t.Fatalf("expected Password to have $2a$10$, got %v", user["Password"])
	}

	deleteUser(t, client, base, token, userID)

	// Verify user is deleted (login should fail)
	_, ok = login(t, client, authServerBaseURL(), "completed@email.com", "UpdatedPass456")
	if ok {
		t.Fatal("login should fail after user deletion")
	}
}
