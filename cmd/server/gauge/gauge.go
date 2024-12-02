package gauge

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type gauge float64

var gaugeMetrics = make(map[string]gauge)
var mu sync.Mutex

func UpdateGaugeHandler(w http.ResponseWriter, r *http.Request) {
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
	defer mu.Unlock()
	gaugeMetrics[metricName] = gauge(parsedNumber)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Gauge metric '%s' set to %.2f successfully.", metricName, parsedNumber)))
}

// handleGetGauge retrieves the gauge metric for a given metric name.
func handleGetGauge(w http.ResponseWriter, r *http.Request, metricName string) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := gaugeMetrics[metricName]

	if !exists {
		http.Error(w, fmt.Sprintf("Gauge metric '%s' not found.", metricName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%.2f", value)))
}
