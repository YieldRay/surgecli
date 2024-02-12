package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// home directory path
var home string

func GetHome() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}

// given a path, return a file path in HOME directory
func AtHome(elem ...string) string {
	if len(home) == 0 {
		home = GetHome()
	}
	e := []string{home}
	e = append(e, elem...)
	return filepath.Join(e...)
}
