// internal/metrics/gauge/gauge.go

// Package gauge handles gauge metric operations, including setting and retrieving gauge values.
package gauge

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/igor-policee/metrics-alerts-service/internal/utils"
)

// Gauge represents a gauge metric.
type Gauge float64

// gaugeMetrics stores the current values of all gauge metrics identified by their names.
var (
	gaugeMetrics = make(map[string]Gauge)
	mu           sync.Mutex
)

// UpdateGaugeHandler processes POST requests to set gauge metrics.
func UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into segments.
	segments := utils.SplitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 5 {
		http.Error(w, "Bad Request. Expected format: /update/gauge/<metricName>/<value>", http.StatusNotFound)
		return
	}

	metricName := segments[3]
	valueStr := segments[4]

	// Ensure the Content-Type header is set to text/plain.
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	// Parse the metric value from string to float64.
	parsedValue, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please provide a valid float.", http.StatusBadRequest)
		return
	}

	// Set the gauge metric safely.
	mu.Lock()
	gaugeMetrics[metricName] = Gauge(parsedValue)
	mu.Unlock()

	// Respond with a success message.
	response := fmt.Sprintf("Gauge metric '%s' set to %.2f successfully.", metricName, parsedValue)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		return
	}
}

// GetGaugeHandler processes GET requests to retrieve gauge metrics.
func GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into segments.
	segments := utils.SplitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 4 {
		http.Error(w, "Bad Request. Expected format: /value/gauge/<metricName>", http.StatusNotFound)
		return
	}

	metricName := segments[3]

	// Retrieve the gauge metric safely.
	mu.Lock()
	value, exists := gaugeMetrics[metricName]
	mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	// Respond with the current value of the gauge metric.
	response := fmt.Sprintf("%.2f", value)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		return
	}
}
