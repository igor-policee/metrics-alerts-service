package main

import (
	"fmt"
	"github.com/igor-policee/metrics-alerts-service/cmd/server/counter"
	"github.com/igor-policee/metrics-alerts-service/cmd/server/gauge"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/update/gauge/", gauge.UpdateGaugeHandler)
	http.HandleFunc("/update/counter/", counter.UpdateCounterHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
