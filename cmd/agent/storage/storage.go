package storage

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func GetGaugeMetrics() map[string]map[string]string {
	// Define the metrics to collect as a set for efficient lookup
	collectMetrics := map[string]struct{}{
		"HeapAlloc":     {},
		"HeapInuse":     {},
		"StackSys":      {},
		"MCacheInuse":   {},
		"MCacheSys":     {},
		"OtherSys":      {},
		"TotalAlloc":    {},
		"HeapSys":       {},
		"HeapReleased":  {},
		"HeapObjects":   {},
		"MSpanInuse":    {},
		"NextGC":        {},
		"Mallocs":       {},
		"StackInuse":    {},
		"GCSys":         {},
		"PauseTotalNs":  {},
		"Alloc":         {},
		"Sys":           {},
		"Lookups":       {},
		"Frees":         {},
		"HeapIdle":      {},
		"MSpanSys":      {},
		"BuckHashSys":   {},
		"LastGC":        {},
		"GCCPUFraction": {},
		"NumForcedGC":   {},
		"NumGC":         {},
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	v := reflect.ValueOf(memStats)
	typeOfMemStats := v.Type()

	metricStorage := make(map[string]map[string]string)
	now := time.Now().Format(time.RFC3339)
	metricStorage[now] = make(map[string]string)

	// Iterate through the fields once and collect the desired metrics
	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfMemStats.Field(i).Name
		if _, exists := collectMetrics[fieldName]; exists {
			field := v.Field(i)
			var valueStr string

			// Convert the field value to a string based on its kind
			switch field.Kind() {
			case reflect.Uint64, reflect.Uint32, reflect.Uint:
				valueStr = fmt.Sprintf("%d", field.Uint())
			case reflect.Float64:
				valueStr = fmt.Sprintf("%f", field.Float())
			default:
				// Skip unsupported types without panicking
				continue
			}

			metricStorage[now][fieldName] = valueStr
		}
	}

	return metricStorage
}
