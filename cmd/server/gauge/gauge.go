// Package gauge provides functionality to handle gauge metrics,
// allowing setting and retrieving gauge values.
package gauge

import (
	"fmt"      // Format package for formatted I/O.
	"net/http" // HTTP package for handling HTTP requests and responses.
	"strconv"  // String conversion package for parsing strings to floats.
	"strings"  // Strings package for string manipulation.
	"sync"     // Synchronization package for mutex to handle concurrent access.
)

// gauge is a custom type based on float64 to represent gauge metrics.
type gauge float64

// gaugeMetrics stores the current values of all gauge metrics identified by their names.
var gaugeMetrics = make(map[string]gauge)

// mu is a mutex to ensure thread-safe access to the gaugeMetrics map.
var mu sync.Mutex

// UpdateGaugeHandler is the HTTP handler for updating and retrieving gauge metrics.
// It handles both POST requests to set a gauge value and GET requests to retrieve a gauge value.
func UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into parts to extract metric name and value.
	pathParts := strings.Split(r.URL.Path, "/")

	// Validate the URL structure.
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL. Expected format: /update/gauge/<metric name>/<metric value>", http.StatusBadRequest)
		return
	}

	// Extract the metric name from the URL.
	metricName := pathParts[3]

	// Determine the HTTP method and handle accordingly.
	switch r.Method {
	case http.MethodPost:
		// For POST requests, expect both metric name and value in the URL.
		if len(pathParts) != 5 {
			http.Error(w, "Invalid URL for POST. Expected format: /update/gauge/<metric name>/<metric value>", http.StatusBadRequest)
			return
		}
		// Delegate to the handler for POST requests.
		handlePostGauge(w, r, metricName, pathParts[4])
	case http.MethodGet:
		// For GET requests, expect only the metric name in the URL.
		if len(pathParts) != 4 {
			http.Error(w, "Invalid URL for GET. Expected format: /update/gauge/<metric name>", http.StatusBadRequest)
			return
		}
		// Delegate to the handler for GET requests.
		handleGetGauge(w, r, metricName)
	default:
		// Respond with Method Not Allowed for unsupported HTTP methods.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePostGauge processes POST requests to set a gauge metric.
// It parses the metric value and updates the corresponding gauge in the map.
func handlePostGauge(w http.ResponseWriter, r *http.Request, metricName string, metricValue string) {
	// Ensure the Content-Type header is set to text/plain.
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	// Parse the metric value from string to float64.
	parsedNumber, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please send a valid float64.", http.StatusBadRequest)
		return
	}

	// Lock the mutex to ensure thread-safe access to gaugeMetrics.
	mu.Lock()
	defer mu.Unlock()

	// Set the gauge metric to the parsed value.
	gaugeMetrics[metricName] = gauge(parsedNumber)

	// Set the response headers and write a success message.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Gauge metric '%s' set to %.2f successfully.", metricName, parsedNumber)))
}

// handleGetGauge processes GET requests to retrieve the value of a gauge metric.
// It fetches the current value of the specified gauge and returns it.
func handleGetGauge(w http.ResponseWriter, r *http.Request, metricName string) {
	// Lock the mutex to ensure thread-safe access to gaugeMetrics.
	mu.Lock()
	defer mu.Unlock()

	// Retrieve the gauge value from the map.
	value, exists := gaugeMetrics[metricName]

	// If the gauge does not exist, respond with a 404 error.
	if !exists {
		http.Error(w, fmt.Sprintf("Gauge metric '%s' not found.", metricName), http.StatusNotFound)
		return
	}

	// Set the response headers and write the gauge value.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%.2f", value)))
}
