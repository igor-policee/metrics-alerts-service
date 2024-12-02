// Package counter provides functionality to handle counter metrics,
// allowing incrementing and retrieving counter values.
package counter

import (
	"fmt"      // Format package for formatted I/O.
	"net/http" // HTTP package for handling HTTP requests and responses.
	"strconv"  // String conversion package for parsing strings to integers.
	"strings"  // Strings package for string manipulation.
	"sync"     // Synchronization package for mutex to handle concurrent access.
)

// counter is a custom type based on int64 to represent counter metrics.
type counter int64

// counterMetrics stores the current values of all counter metrics identified by their names.
var counterMetrics = make(map[string]counter)

// mu is a mutex to ensure thread-safe access to the counterMetrics map.
var mu sync.Mutex

// UpdateCounterHandler is the HTTP handler for updating and retrieving counter metrics.
// It handles both POST requests to update a counter and GET requests to retrieve a counter value.
func UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into parts to extract metric name and value.
	pathParts := strings.Split(r.URL.Path, "/")

	// Validate the URL structure.
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL. Expected format: /update/counter/<metric name>/<metric value>",
			http.StatusBadRequest)
		return
	}

	// Extract the metric name from the URL.
	metricName := pathParts[3]

	// Determine the HTTP method and handle accordingly.
	switch r.Method {
	case http.MethodPost:
		// For POST requests, expect both metric name and value in the URL.
		if len(pathParts) != 5 {
			http.Error(w, "Invalid URL for POST. Expected format: /update/counter/<metric name>/<metric value>", http.StatusBadRequest)
			return
		}
		// Delegate to the handler for POST requests.
		handlePostCounter(w, r, metricName, pathParts[4])
	case http.MethodGet:
		// For GET requests, expect only the metric name in the URL.
		if len(pathParts) != 4 {
			http.Error(w, "Invalid URL for GET. Expected format: /update/counter/<metric name>", http.StatusBadRequest)
			return
		}
		// Delegate to the handler for GET requests.
		handleGetCounter(w, metricName)
	default:
		// Respond with Method Not Allowed for unsupported HTTP methods.
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePostCounter processes POST requests to update a counter metric.
// It parses the metric value and updates the corresponding counter in the map.
func handlePostCounter(w http.ResponseWriter, r *http.Request, metricName string, metricValue string) {
	// Ensure the Content-Type header is set to text/plain.
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	// Parse the metric value from string to int64.
	parsedNumber, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please send a valid int64.", http.StatusBadRequest)
		return
	}

	// Lock the mutex to ensure thread-safe access to counterMetrics.
	mu.Lock()
	defer mu.Unlock()

	// Increment the counter metric by the parsed value.
	counterMetrics[metricName] += counter(parsedNumber)

	// Set the response headers and write a success message.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(fmt.Sprintf("Counter metric '%s' set to %d successfully.", metricName, parsedNumber)))
	if err != nil {
		http.Error(w, "Unknown error.", http.StatusInternalServerError)
		return
	}
}

// handleGetCounter processes GET requests to retrieve the value of a counter metric.
// It fetches the current value of the specified counter and returns it.
func handleGetCounter(w http.ResponseWriter, metricName string) {
	// Lock the mutex to ensure thread-safe access to counterMetrics.
	mu.Lock()
	defer mu.Unlock()

	// Retrieve the counter value from the map.
	value, exists := counterMetrics[metricName]

	// If the counter does not exist, respond with a 404 error.
	if !exists {
		http.Error(w, fmt.Sprintf("Counter metric '%s' not found.", metricName), http.StatusNotFound)
		return
	}

	// Set the response headers and write the counter value.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf("%d", value)))
	if err != nil {
		http.Error(w, "Unknown error.", http.StatusInternalServerError)
		return
	}
}
