package main

import (
	"fmt"
	"github.com/igor-policee/metrics-alerts-service/cmd/agent/sender"
	"github.com/igor-policee/metrics-alerts-service/cmd/agent/storage"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("--> Hello Gophers! <--")

	serverAddress := "http://localhost:8080"
	metricOperation := "update"
	metricType := "gauge"

	for _, metrics := range storage.CollectGaugeMetrics() {
		for metricName, metricValue := range metrics {
			endpointSlice := []string{serverAddress, metricOperation, metricType, metricName, strconv.FormatUint(metricValue, 10)}
			endpointString := strings.Join(endpointSlice, "/")
			sender.SendPostRequest(endpointString)
		}
	}

	fmt.Println("--> Good Bye Gophers! <--")
}
