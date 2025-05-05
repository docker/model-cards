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
