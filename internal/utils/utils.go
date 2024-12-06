// internal/utils/utils.go

// Package utils provides common utility functions for the Metrics Alerts Service.
package utils

import "strings"

// SplitPath splits the URL path into segments and removes any trailing slash.
func SplitPath(path string) []string {
	// Remove any trailing slash and split the path.
	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	return strings.Split(path, "/")
}
