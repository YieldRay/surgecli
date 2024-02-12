package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var termWidth int

func init() {
	fd := int(os.Stdin.Fd())
	var err error

	termWidth, _, err = term.GetSize(fd)
	if err != nil {
		termWidth = 100
	}
}

func ClearLine() {
	fmt.Print("\r")
	fmt.Print(strings.Repeat(" ", termWidth))
	fmt.Print("\r")
}
