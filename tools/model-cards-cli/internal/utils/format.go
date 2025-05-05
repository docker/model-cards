package utils

import (
	"fmt"
	"math"
)

// FormatVRAM converts bytes to GB and returns a formatted string
// The value is rounded to 2 decimal places
func FormatVRAM(bytes float64) string {
	// Convert bytes to GB (1 GB = 1024^3 bytes)
	gb := bytes / (1024 * 1024 * 1024)

	// Round to 2 decimal places
	rounded := math.Round(gb*100) / 100

	return fmt.Sprintf("%.2f GB", rounded)
}

// FormatContextLength formats a token count with K/M/B suffixes
// For example: 1000 -> "1K", 1500 -> "1.5K", 1000000 -> "1M"
func FormatContextLength(tokens uint32) string {
	const (
		K = 1000
		M = K * 1000
		B = M * 1000
	)

	switch {
	case tokens >= B:
		return fmt.Sprintf("%dB", int(math.Round(float64(tokens)/float64(B))))
	case tokens >= M:
		return fmt.Sprintf("%dM tokens", int(math.Round(float64(tokens)/float64(M))))
	case tokens >= K:
		return fmt.Sprintf("%dK tokens", int(math.Round(float64(tokens)/float64(K))))
	default:
		return fmt.Sprintf("%d tokens", tokens)
	}
}
