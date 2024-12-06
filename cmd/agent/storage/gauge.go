// cmd/agent/storage/gauge.go
package storage

import (
	"fmt"
	"runtime"
	"time"
)

func GetMemStats() map[string]map[string]string {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	metricStorage := make(map[string]map[string]string)
	now := time.Now().Format(time.RFC3339)
	metricStorage[now] = make(map[string]string)

	collectMetrics := map[string]interface{}{
		"HeapAlloc":     memStats.HeapAlloc,
		"HeapInuse":     memStats.HeapInuse,
		"StackSys":      memStats.StackSys,
		"MCacheInuse":   memStats.MCacheInuse,
		"MCacheSys":     memStats.MCacheSys,
		"OtherSys":      memStats.OtherSys,
		"TotalAlloc":    memStats.TotalAlloc,
		"HeapSys":       memStats.HeapSys,
		"HeapReleased":  memStats.HeapReleased,
		"HeapObjects":   memStats.HeapObjects,
		"MSpanInuse":    memStats.MSpanInuse,
		"NextGC":        memStats.NextGC,
		"Mallocs":       memStats.Mallocs,
		"StackInuse":    memStats.StackInuse,
		"GCSys":         memStats.GCSys,
		"PauseTotalNs":  memStats.PauseTotalNs,
		"Alloc":         memStats.Alloc,
		"Sys":           memStats.Sys,
		"Lookups":       memStats.Lookups,
		"Frees":         memStats.Frees,
		"HeapIdle":      memStats.HeapIdle,
		"MSpanSys":      memStats.MSpanSys,
		"BuckHashSys":   memStats.BuckHashSys,
		"LastGC":        memStats.LastGC,
		"GCCPUFraction": memStats.GCCPUFraction,
		"NumForcedGC":   memStats.NumForcedGC,
		"NumGC":         memStats.NumGC,
	}

	for metricName, metricValue := range collectMetrics {
		var valueStr string
		switch v := metricValue.(type) {
		case uint64, uint32, uint:
			valueStr = fmt.Sprintf("%d", v)
		case float64:
			valueStr = fmt.Sprintf("%f", v)
		default:
			continue
		}
		metricStorage[now][metricName] = valueStr
	}

	return metricStorage
}
