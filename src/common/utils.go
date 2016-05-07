// Place any kind of application-wide utility here

package common

import (
	"os"
	"strconv"
)

// StringInSlice searches a slice for a string
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// ConvertToInt converts s to an int and ignores errors
func ConvertToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// FileExists returns whether a file exists or not
func FileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
