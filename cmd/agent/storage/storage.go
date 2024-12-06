package storage

import (
	"reflect"
	"runtime"
	"time"
)

func CollectGaugeMetrics() map[string]map[string]uint64 {
	metricStorage := make(map[string]map[string]uint64)
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	v := reflect.ValueOf(memStats)
	typeOfMemStats := v.Type()

	now := time.Now().Format(time.RFC3339)
	metricStorage[now] = make(map[string]uint64)

	for i := 0; i < v.NumField(); i++ {
		field := typeOfMemStats.Field(i).Name
		fieldValue := v.Field(i)

		if fieldValue.Kind() == reflect.Uint64 {
			metricStorage[now][field] = fieldValue.Uint()
		}
	}

	return metricStorage
}
