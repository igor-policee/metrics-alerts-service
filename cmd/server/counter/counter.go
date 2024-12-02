package counter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type counter int64

var counterMetrics = make(map[string]counter)
var mu sync.Mutex

func UpdateCounterHandler(w http.ResponseWriter, r *http.Request) {

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL. Expected format: /update/counter/<metric name>/<metric value>",
			http.StatusBadRequest)
		return
	}

	metricName := pathParts[3]

	switch r.Method {
	case http.MethodPost:
		if len(pathParts) != 5 {
			http.Error(w, "Invalid URL for POST. Expected format: /update/counter/<metric name>/<metric value>", http.StatusBadRequest)
			return
		}
		handlePostCounter(w, r, metricName, pathParts[4])
	case http.MethodGet:
		if len(pathParts) != 4 {
			http.Error(w, "Invalid URL for GET. Expected format: /update/counter/<metric name>", http.StatusBadRequest)
			return
		}
		handleGetCounter(w, r, metricName)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handlePostCounter(w http.ResponseWriter, r *http.Request, metricName string, metricValue string) {
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Unsupported Media Type. Use text/plain.", http.StatusUnsupportedMediaType)
		return
	}

	parsedNumber, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		http.Error(w, "Invalid metric value. Please send a valid int64.", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()
	counterMetrics[metricName] += counter(parsedNumber)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Counter metric '%s' set to %d successfully.", metricName, parsedNumber)))
}

func handleGetCounter(w http.ResponseWriter, r *http.Request, metricName string) {
	mu.Lock()
	defer mu.Unlock()
	value, exists := counterMetrics[metricName]

	if !exists {
		http.Error(w, fmt.Sprintf("Counter metric '%s' not found.", metricName), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%d", value)))
}
