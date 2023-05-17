package utils

import "os"

func IsDir(f string) bool {
	fi, err := os.Stat(f)
	return !os.IsNotExist(err) && fi.IsDir()
}

func IsFileExist(path string) bool {
	if s, err := os.Stat(path); err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	} else {
		return s.Mode().IsRegular()
	}
}
