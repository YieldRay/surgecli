package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return strconv.FormatInt(bytes, 10) + " B"
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func JSONStringify(v any) string {
	j, _ := json.MarshalIndent(v, "", "    ")
	return (string(j))
}
