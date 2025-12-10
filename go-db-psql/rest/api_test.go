package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

func testServerBaseURL() string {
	if u := os.Getenv("TEST_SERVER_URL"); u != "" {
		return u
	}
	return "http://localhost:3000"
}

// TestUserCRUD_Server performs an integration test against a running server.
// Ensure your server is running (default http://localhost:3000) before running.
func TestUserCRUD_Server(t *testing.T) {
	client := &http.Client{Timeout: 5 * time.Second}
	base := testServerBaseURL()

	// Create user
	user := map[string]interface{}{
		"firstname": "Test",
		"lastname":  "User",
		"dob":       "1990-01-01",
		"email":     "test@email.com",
	}

	buf, _ := json.Marshal(user)
	resp, err := client.Post(base+"/users", "application/json", bytes.NewReader(buf))
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

	// Update user - send ID as integer, not string
	update := map[string]interface{}{
		//"id":        int(idVal.(float64)), // JSON numbers decode to float64
		"email":     "updatetest@email.com",
		"firstname": "UpdateFirstName",
		"dob":       "1973-12-12",
		"lastname":  "UpdateLastName",
	}
	ubuf, _ := json.Marshal(update)
	req, _ := http.NewRequest(http.MethodPut, base+"/users/"+id, bytes.NewReader(ubuf))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected 200 OK on update, got %d: %s", resp.StatusCode, string(b))
	}

	// Get user
	resp, err = client.Get(base + "/users/" + id)
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
	fmt.Println(resp.Body)
	fmt.Println(got)

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

	// Delete user
	req, _ = http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
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
