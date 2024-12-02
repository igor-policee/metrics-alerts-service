package main

import (
	"fmt"
	"github.com/igor-policee/metrics-alerts-service/cmd/server/gauge"
	"log"
	"net/http"
)

type counter int64

func main() {
	// Register handlers with pattern matching for metric names and values.
	http.HandleFunc("/update/gauge/", gauge.UpdateGaugeHandler)
	// http.HandleFunc("/update/gauge/", updateCounterHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
