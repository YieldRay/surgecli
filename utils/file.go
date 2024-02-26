package utils

import (
	"os"
	"strings"
)

func IsDir(f string) bool {
	fi, err := os.Stat(f)
	return err == nil && fi.IsDir()
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

func ReadTextFile(name string) (string, error) {
	b, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ReadTextFileTrim(name string) (string, error) {
	str, err := ReadTextFile(name)
	if err != nil {
		return "", err
	}
	return strings.Trim(str, " "), nil
}
