package main

import (
	"fmt"
	"github.com/igor-policee/metrics-alerts-service/cmd/agent/sender"
	"github.com/igor-policee/metrics-alerts-service/cmd/agent/storage"
	"math/rand/v2"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("--> Hello Gophers! <--")

	const serverAddress = "http://localhost:8080"
	const metricOperation = "update"
	const pollInterval = 5

	for {

		metricType := "gauge"
		for _, metrics := range storage.GetMemStats() {
			for metricName, metricValue := range metrics {
				endpointSlice := []string{serverAddress, metricOperation, metricType, metricName, metricValue}
				endpointString := strings.Join(endpointSlice, "/")
				sender.SendPostRequest(endpointString)
			}
		}

		metricType = "gauge"
		metricName := "RandomValue"
		randomNumber := rand.Int()
		metricValue := strconv.Itoa(randomNumber)
		endpointSlice := []string{serverAddress, metricOperation, metricType, metricName, metricValue}
		endpointString := strings.Join(endpointSlice, "/")
		fmt.Println(endpointString)
		sender.SendPostRequest(endpointString)

		metricType = "counter"
		metricName = "PollCount"
		metricValue = strconv.Itoa(1)
		endpointSlice = []string{serverAddress, metricOperation, metricType, metricName, metricValue}
		endpointString = strings.Join(endpointSlice, "/")
		sender.SendPostRequest(endpointString)

		time.Sleep(pollInterval * time.Second)

		fmt.Println("--> Good Bye Gophers! <--")

	}
}
