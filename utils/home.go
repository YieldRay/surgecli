package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// home directory path
var home string

func getHome() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	} else {
		return os.Getenv("HOME")
	}
}

// given a path, return a file path in HOME directory
func AtHome(elem ...string) string {
	if len(home) == 0 {
		home = getHome()
	}
	e := []string{home}
	e = append(e, elem...)
	return filepath.Join(e...)
}
