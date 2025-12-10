package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func testServersBaseURL() string {
	return "http://localhost:3000"
}

// TestUserCRUD_ServerIntegration performs an integration test against a running server.
// Ensure your server is running (default http://localhost:3000) before running.
func TestUserCRUD_ServerIntegration(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServersBaseURL()

	// Step 1: Login to get JWT token
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

	// Helper function to add auth header
	authHeader := func(req *http.Request) {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Step 2: Create user
	user := map[string]interface{}{
		"firstname": "Test",
		"lastname":  "User",
		"dob":       "1990-01-01",
		"email":     "test@email.com",
	}

	buf, _ := json.Marshal(user)
	req := &http.Request{}
	req, _ = http.NewRequest(http.MethodPost, base+"/users", bytes.NewReader(buf))
	req.Header.Set("Content-Type", "application/json")
	authHeader(req)
	resp, err = client.Do(req)
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
	id := fmt.Sprintf("%v", idVal)

	// Step 3: Update user
	update := map[string]interface{}{
		"email":     "updatetest@email.com",
		"firstname": "UpdateFirstName",
		"dob":       "1973-12-12",
		"lastname":  "UpdateLastName",
	}
	ubuf, _ := json.Marshal(update)
	req, _ = http.NewRequest(http.MethodPut, base+"/users/"+id, bytes.NewReader(ubuf))
	req.Header.Set("Content-Type", "application/json")
	authHeader(req)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on update, got %d: %s", resp.StatusCode, string(b))
	}

	// Step 4: Get user
	req, _ = http.NewRequest(http.MethodGet, base+"/users/"+id, nil)
	authHeader(req)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("get request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on get, got %d: %s", resp.StatusCode, string(b))
	}

	var got map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("decode get response: %v", err)
	}

	if got["FirstName"] != "UpdateFirstName" {
		t.Fatalf("expected firstname 'UpdateFirstName', got %v", got["FirstName"])
	}

	if got["LastName"] != "UpdateLastName" {
		t.Fatalf("expected lastname 'UpdateLastName', got %v", got["LastName"])
	}

	if got["DOB"] != "1973-12-12T00:00:00Z" {
		t.Fatalf("expected dob '1973-12-12T00:00:00', got %v", got["DOB"])
	}

	if got["Email"] != "updatetest@email.com" {
		t.Fatalf("expected email 'updatetest@email.com', got %v", got["Email"])
	}

	// Step 5: Delete user
	req, _ = http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
	authHeader(req)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 or 204 on delete, got %d: %s", resp.StatusCode, string(b))
	}
}
