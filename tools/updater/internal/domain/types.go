package domain

import (
	"fmt"
	"strconv"
	"strings"
)

// Constants for formatting and display
const (
	// MediaTypeGGUF is the media type for GGUF files in OCI manifests
	MediaTypeGGUF = "application/vnd.docker.ai.gguf.v3"
)

// FormatParameters formats the parameters to match the table format
func FormatParameters(params string) string {
	// If already formatted with M or B suffix, return as is
	if strings.HasSuffix(params, "M") || strings.HasSuffix(params, "B") {
		return params
	}

	// Try to parse as a number
	num, err := strconv.ParseFloat(params, 64)
	if err != nil {
		return params
	}

	// Format based on size
	if num >= 1000000000 {
		return fmt.Sprintf("%.1fB", num/1000000000)
	} else if num >= 1000000 {
		return fmt.Sprintf("%.0fM", num/1000000)
	}

	return params
}
