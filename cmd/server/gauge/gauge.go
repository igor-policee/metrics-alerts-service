// Package gauge handles gauge metric operations, including setting and retrieving gauge values.
package gauge

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
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
	segments := splitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 5 {
		http.Error(w, "Bad Request. Expected format: /update/gauge/<metricName>/<value>", http.StatusBadRequest)
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
	defer mu.Unlock()
	gaugeMetrics[metricName] = Gauge(parsedValue)

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
	segments := splitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 4 {
		http.Error(w, "Bad Request. Expected format: /value/gauge/<metricName>", http.StatusBadRequest)
		return
	}

	metricName := segments[3]

	// Retrieve the gauge metric safely.
	mu.Lock()
	defer mu.Unlock()
	value, exists := gaugeMetrics[metricName]

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

// splitPath is a helper function to split the URL path and remove empty segments.
func splitPath(path string) []string {
	// Remove any trailing slash and split the path.
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}
