package main

import (
	"strings"
)

func inSlice(key string, slice []string) bool {
	for _, item := range slice {
		if strings.ToLower(key) == strings.ToLower(item) {
			return true
		}
	}
	return false
}