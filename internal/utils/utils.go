package utils

import (
	"os"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}