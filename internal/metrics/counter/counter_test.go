// internal/metrics/counter/counter_test.go

package counter

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateCounterHandler(t *testing.T) {
	// Initialization
	req, err := http.NewRequest("POST", "/update/counter/testCounter/10", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	// Call the handler
	UpdateCounterHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	// Check the response content
	expected := "Counter metric 'testCounter' updated to 10 successfully."
	if rr.Body.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestGetCounterHandler(t *testing.T) {
	// Pre-set counter value
	mu.Lock()
	counterMetrics["testCounter"] = 15
	mu.Unlock()

	// Create the request
	req, err := http.NewRequest("GET", "/value/counter/testCounter", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	GetCounterHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	// Check the response content
	expected := "15"
	if rr.Body.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestGetCounterHandler_NotFound(t *testing.T) {
	// Create a request for a non-existent counter
	req, err := http.NewRequest("GET", "/value/counter/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	GetCounterHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}
