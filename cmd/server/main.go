// cmd/server/main.go

// Package main is the entry point for the Metrics Alerts Service server application.
package main

import (
	"log"
	"net/http"

	"github.com/igor-policee/metrics-alerts-service/internal/metrics/counter"
	"github.com/igor-policee/metrics-alerts-service/internal/metrics/gauge"
	"github.com/igor-policee/metrics-alerts-service/internal/utils"
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
	segments := utils.SplitPath(r.URL.Path)

	// Expecting exactly 5 segments: "", "update", "metricType", "metricName", "value".
	if len(segments) != 5 {
		http.Error(w, "Bad Request. Expected format: /update/<metricType>/<metricName>/<value>", http.StatusNotFound)
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
