package utils

import (
	"fmt"
	"io"
)

const CLEAR_LINE = "\r\u001b[K"

// / ClearLine fprintf
func Cfprintf(w io.Writer, format string, a ...any) {
	fmt.Fprintf(w, "\r\u001b[K%s", fmt.Sprintf(format, a...))
}
