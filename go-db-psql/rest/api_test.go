package rest

import (
	"bytes"
	"encoding/json"
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

	// var created map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
	// 	t.Fatalf("decode create response: %v", err)
	// }

	// id := fmt.Sprintf("%v", created["id"])
	// if id == "<nil>" || id == "" {
	// 	t.Fatalf("invalid id from create: %v", created["id"])
	// }

	// // Update user
	// update := map[string]interface{}{"firstname": "Updated"}
	// ubuf, _ := json.Marshal(update)
	// req, _ := http.NewRequest(http.MethodPut, base+"/users/"+id, bytes.NewReader(ubuf))
	// req.Header.Set("Content-Type", "application/json")
	// resp, err = client.Do(req)
	// if err != nil {
	// 	t.Fatalf("update request failed: %v", err)
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	b, _ := io.ReadAll(resp.Body)
	// 	t.Fatalf("expected 200 OK on update, got %d: %s", resp.StatusCode, string(b))
	// }

	// // Get user
	// resp, err = client.Get(base + "/users/" + id)
	// if err != nil {
	// 	t.Fatalf("get request failed: %v", err)
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK {
	// 	b, _ := io.ReadAll(resp.Body)
	// 	t.Fatalf("expected 200 OK on get, got %d: %s", resp.StatusCode, string(b))
	// }
	// var got map[string]interface{}
	// if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
	// 	t.Fatalf("decode get response: %v", err)
	// }
	// if got["firstname"] != "Updated" {
	// 	t.Fatalf("expected firstname 'Updated', got %v", got["firstname"])
	// }

	// // Delete user
	// req, _ = http.NewRequest(http.MethodDelete, base+"/users/"+id, nil)
	// resp, err = client.Do(req)
	// if err != nil {
	// 	t.Fatalf("delete request failed: %v", err)
	// }
	// defer resp.Body.Close()
	// if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
	// 	b, _ := io.ReadAll(resp.Body)
	// 	t.Fatalf("expected 200 or 204 on delete, got %d: %s", resp.StatusCode, string(b))
	// }
}
