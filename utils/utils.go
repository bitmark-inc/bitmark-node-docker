package utils

import (
	"os"
)

// check if file exists
func EnsureFileExists(name string) bool {
	_, err := os.Stat(name)
	return nil == err
}
