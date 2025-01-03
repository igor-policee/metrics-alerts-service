// internal/metrics/counter/counter.go

// Package counter handles counter metric operations, including incrementing and retrieving counter values.
package counter

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/igor-policee/metrics-alerts-service/internal/utils"
)

// Counter represents a counter metric.
type Counter int64

// counterMetrics stores the current values of all counter metrics identified by their names.
var (
	counterMetrics = make(map[string]Counter)
	mu             sync.Mutex
)

// UpdateCounterHandler processes POST requests to update counter metrics.
func UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into segments.
	segments := utils.SplitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 5 {
		http.Error(w, "Bad Request. Expected format: /update/counter/<metricName>/<value>", http.StatusNotFound)
		return
	}

	metricName := segments[3]
	valueStr := segments[4]

	// Ensure the Content-Type header is set to text/plain.
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	// Parse the metric value from string to int64.
	parsedValue, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please provide a valid integer.", http.StatusBadRequest)
		return
	}

	// Update the counter metric safely.
	mu.Lock()
	counterMetrics[metricName] += Counter(parsedValue)
	currentValue := counterMetrics[metricName]
	mu.Unlock()

	// Respond with a success message.
	response := fmt.Sprintf("Counter metric '%s' updated to %d successfully.", metricName, currentValue)
	w.Header().Set("Content-Type", "text/plain")
	_, err = w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		return
	}
}

// GetCounterHandler processes GET requests to retrieve counter metrics.
func GetCounterHandler(w http.ResponseWriter, r *http.Request) {
	// Split the URL path into segments.
	segments := utils.SplitPath(r.URL.Path)

	// Validate the URL structure.
	if len(segments) != 4 {
		http.Error(w, "Bad Request. Expected format: /value/counter/<metricName>", http.StatusNotFound)
		return
	}

	metricName := segments[3]

	// Retrieve the counter metric safely.
	mu.Lock()
	value, exists := counterMetrics[metricName]
	mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	// Respond with the current value of the counter metric.
	response := fmt.Sprintf("%d", value)
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(response))
	if err != nil {
		http.Error(w, "Internal Server Error.", http.StatusInternalServerError)
		return
	}
}
