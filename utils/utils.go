// Package Utils provides basic utilities for the errand server.
package utils

import (
	"time"

	"github.com/polygon-io/errands-server/schemas"
)

func GetTimestamp() int64 {
	return time.Now().UnixNano() / 1_000_000
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func ContainsStatus(slice []schemas.Status, status schemas.Status) bool {
	for _, s := range slice {
		if s == status {
			return true
		}
	}

	return false
}
