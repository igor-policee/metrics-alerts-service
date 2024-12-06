// cmd/agent/main.go

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/igor-policee/metrics-alerts-service/cmd/agent/sender"
	"github.com/igor-policee/metrics-alerts-service/cmd/agent/storage"
)

func main() {
	fmt.Println("--> Hello Gophers! <--")

	const serverAddress = "http://localhost:8080"
	const metricOperation = "update"
	const pollInterval = 5

	pollCount := 0

	for {
		// Collect and send gauge metrics
		metricType := "gauge"
		memStats := storage.GetMemStats()
		for _, metrics := range memStats {
			for metricName, metricValue := range metrics {
				endpointSlice := []string{serverAddress, metricOperation, metricType, metricName, metricValue}
				endpointString := strings.Join(endpointSlice, "/")
				sender.SendPostRequest(endpointString)
			}
		}

		// Send a random gauge metric
		metricName := "RandomValue"
		randomNumber := rand.Int()
		metricValue := strconv.Itoa(randomNumber)
		endpointSlice := []string{serverAddress, metricOperation, metricType, metricName, metricValue}
		endpointString := strings.Join(endpointSlice, "/")
		fmt.Println(endpointString)
		sender.SendPostRequest(endpointString)

		// Send a counter metric
		metricType = "counter"
		metricName = "PollCount"
		pollCount++
		metricValue = strconv.Itoa(pollCount)
		endpointSlice = []string{serverAddress, metricOperation, metricType, metricName, metricValue}
		endpointString = strings.Join(endpointSlice, "/")
		sender.SendPostRequest(endpointString)

		time.Sleep(time.Duration(pollInterval) * time.Second)

		fmt.Println("--> Good Bye Gophers! <--")
	}
}
