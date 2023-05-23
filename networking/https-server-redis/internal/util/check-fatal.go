// Package util exposes utilities
package util

import (
	"fmt"
	"os"
)

// CheckFatal checks if provided err is non-nil; if so, prints out the error,
// a custom message and exits the application.
func CheckFatal(err error, msg string) {
	if err != nil {
		fmt.Printf("Fatal error situation: %s. Error: %v\n", msg, err)
		os.Exit(1)
	}
}