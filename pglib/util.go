package pglib

import "os"

func IsFileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsDirExist(path string) bool {
	return IsFileExist(path)
}
