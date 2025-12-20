package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/calculate?ip=192.168.1.1&netmask=255.255.255.0", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(calculateHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expectedResponse := ApiResponse{
		IPAddress:        "192.168.1.1",
		Netmask:          "255.255.255.0",
		NetworkAddress:   "192.168.1.0",
		BroadcastAddress: "192.168.1.255",
		UsableHosts:      254,
	}

	var actualResponse ApiResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}

	if actualResponse != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", actualResponse, expectedResponse)
	}
}
