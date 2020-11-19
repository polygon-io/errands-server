// Package Utils provides basic utilities for the errand server.
//nolint:golint,stylecheck // TODO
package utils

import (
	"time"
)

func GetTimestamp() int64 {
	return (time.Now().UnixNano() / 1_000_000)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
