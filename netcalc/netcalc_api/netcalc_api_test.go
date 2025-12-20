package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	expectedResponse := ApiResponse{
		IPAddress:        "192.168.1.1",
		Netmask:          "255.255.252.0",
		NetworkAddress:   "192.168.0.0",
		BroadcastAddress: "192.168.3.255",
		UsableHosts:      1022,
	}

	req, err := http.NewRequest("GET", "/calculate?ip=192.168.1.1&netmask=255.255.252.0", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(calculateHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("1:handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var actualResponse ApiResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &actualResponse); err != nil {
		t.Fatalf("Could not unmarshal response: %v", err)
	}

	if actualResponse != expectedResponse {
		t.Errorf("1: handler returned unexpected body: got %v want %v", actualResponse, expectedResponse)
	}

	req, err = http.NewRequest("GET", "/calculate?cidr=192.168.1.1/23", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("2: handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if actualResponse != expectedResponse {
		t.Errorf("2: handler returned unexpected body: got %v want %v", actualResponse, expectedResponse)
	}
}
