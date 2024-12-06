// internal/metrics/gauge/gauge_test.go

package gauge

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateGaugeHandler(t *testing.T) {
	// Initialization
	req, err := http.NewRequest("POST", "/update/gauge/testGauge/75.5", strings.NewReader(""))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	// Call the handler
	UpdateGaugeHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	// Check the response content
	expected := "Gauge metric 'testGauge' set to 75.50 successfully."
	if rr.Body.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestGetGaugeHandler(t *testing.T) {
	// Pre-set gauge value
	mu.Lock()
	gaugeMetrics["testGauge"] = 80.25
	mu.Unlock()

	// Create the request
	req, err := http.NewRequest("GET", "/value/gauge/testGauge", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	GetGaugeHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, status)
	}

	// Check the response content
	expected := "80.25"
	if rr.Body.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestGetGaugeHandler_NotFound(t *testing.T) {
	// Create a request for a non-existent gauge
	req, err := http.NewRequest("GET", "/value/gauge/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Call the handler
	GetGaugeHandler(rr, req)

	// Check the response status
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Expected status %v, got %v", http.StatusNotFound, status)
	}
}
