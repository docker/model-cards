package markdown

import (
	"fmt"
	"strconv"
	"strings"
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

// FormatContextLength formats a token count to a human-readable format
// Examples: 1000 -> "1K tokens", 1500000 -> "1.5M tokens"
func FormatContextLength(length uint32) string {
	if length == 0 {
		return "-"
	}

	switch {
	case length >= 1000000:
		return fmt.Sprintf("%.1fM tokens", float64(length)/1000000)
	case length >= 1000:
		return fmt.Sprintf("%.1fK tokens", float64(length)/1000)
	default:
		return fmt.Sprintf("%d tokens", length)
	}
}
