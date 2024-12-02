// Package main is the entry point for the Metrics Alerts Service server application.
package main

import (
	"github.com/igor-policee/metrics-alerts-service/cmd/server/counter"
	"github.com/igor-policee/metrics-alerts-service/cmd/server/gauge"
	"log"
	"net/http"
	"strings"
)

// supportedMetrics maps the supported metric types to their corresponding handlers.
var supportedMetrics = map[string]http.HandlerFunc{
	"counter": counter.UpdateCounterHandler,
	"gauge":   gauge.UpdateGaugeHandler,
}

// updateHandler routes the update requests to the appropriate metric handler.
func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is POST.
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Example URL: /update/unknown/testCounter/100
	// Split the URL path into segments.
	segments := splitPath(r.URL.Path)

	// Expecting exactly 5 segments: "", "update", "metricType", "metricName", "value".
	if len(segments) != 5 {
		http.Error(w, "Bad Request. Expected format: /update/<metricType>/<metricName>/<value>", http.StatusBadRequest)
		return
	}

	metricType := segments[2]

	// Retrieve the handler for the metric type.
	handler, exists := supportedMetrics[metricType]
	if !exists {
		http.Error(w, "Bad Request. Unknown <metricType>", http.StatusBadRequest)
		return
	}

	// Delegate to the specific metric handler.
	handler(w, r)
}

// splitPath is a helper function to split the URL path and remove empty segments.
func splitPath(path string) []string {
	// Remove any trailing slash and split the path.
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}

func main() {
	// Register the updateHandler for the /update/ path.
	http.HandleFunc("/update/", updateHandler)

	// Register handlers for retrieving metric values.
	http.HandleFunc("/value/counter/", counter.GetCounterHandler)
	http.HandleFunc("/value/gauge/", gauge.GetGaugeHandler)

	// Start the HTTP server.
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
