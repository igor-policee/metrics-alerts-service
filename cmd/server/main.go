package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type gauge float64
type counter int64

var (
	// Use maps to store multiple gauge and counter metrics identified by their names.
	gaugeMetrics   = make(map[string]gauge)
	counterMetrics = make(map[string]counter)
	mu             sync.Mutex
)

func main() {
	// Register handlers with pattern matching for metric names and values.
	http.HandleFunc("/update/gauge/", updateGaugeHandler)
	// http.HandleFunc("/update/counter/", updateCounterHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// updateGaugeHandler handles both setting and getting gauge metrics based on the request path and method.
func updateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	// Expected paths:
	// POST /update/gauge/<metric name>/<metric value>
	// GET  /update/gauge/<metric name>

	pathParts := strings.Split(r.URL.Path, "/")
	// Path should be ['', 'update', 'gauge', '<metric name>', '<metric value>']
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL. Expected format: /update/gauge/<metric name>/<metric value>", http.StatusBadRequest)
		return
	}

	metricName := pathParts[3]

	switch r.Method {
	case http.MethodPost:
		if len(pathParts) != 5 {
			http.Error(w, "Invalid URL for POST. Expected format: /update/gauge/<metric name>/<metric value>", http.StatusBadRequest)
			return
		}
		handlePostGauge(w, r, metricName, pathParts[4])
	case http.MethodGet:
		if len(pathParts) != 4 {
			http.Error(w, "Invalid URL for GET. Expected format: /update/gauge/<metric name>", http.StatusBadRequest)
			return
		}
		handleGetGauge(w, r, metricName)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// handlePostGauge sets the gauge metric for a given metric name.
func handlePostGauge(w http.ResponseWriter, r *http.Request, metricName string, metricValue string) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	parsedNumber, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please send a valid float64.", http.StatusBadRequest)
		return
	}

	mu.Lock()
	gaugeMetrics[metricName] = gauge(parsedNumber)
	mu.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Gauge metric '%s' set to %.2f successfully.", metricName, parsedNumber)))
}

// handleGetGauge retrieves the gauge metric for a given metric name.
func handleGetGauge(w http.ResponseWriter, r *http.Request, metricName string) {
	mu.Lock()
	value, exists := gaugeMetrics[metricName]
	mu.Unlock()

	if !exists {
		http.Error(w, fmt.Sprintf("Gauge metric '%s' not found.", metricName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%.2f", value)))
}
