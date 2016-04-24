package util

import (
	"os"
	"path/filepath"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}

func SaveCurrentDir() string {
	prevDir, _ := filepath.Abs(".")
	return prevDir
}
