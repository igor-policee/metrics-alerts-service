// Package main is the entry point for the Metrics Alerts Service server application.
package main

import (
	"fmt"      // Format package for formatted I/O.
	"log"      // Log package for logging errors.
	"net/http" // HTTP package for creating the server and handling requests.

	// Importing the counter and gauge sub-packages for handling respective metrics.
	"github.com/igor-policee/metrics-alerts-service/cmd/server/counter"
	"github.com/igor-policee/metrics-alerts-service/cmd/server/gauge"
)

// main function initializes the HTTP server and sets up route handlers.
func main() {
	// Register the gauge update handler for the "/update/gauge/" route.
	http.HandleFunc("/update/gauge/", gauge.UpdateGaugeHandler)

	// Register the counter update handler for the "/update/counter/" route.
	http.HandleFunc("/update/counter/", counter.UpdateCounterHandler)

	// Informational message indicating the server is starting.
	fmt.Println("Server is running on port 8080...")

	// Start the HTTP server on port 8080 and log any fatal errors.
	log.Fatal(http.ListenAndServe(":8080", nil))
}
