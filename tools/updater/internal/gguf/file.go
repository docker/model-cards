package gguf

import (
	"fmt"
	"strconv"
	"strings"

	parser "github.com/gpustack/gguf-parser-go"
)

// FieldNotFoundError represents an error when a required field is not found in the GGUF file
type FieldNotFoundError struct {
	Field string
}

// Error implements the error interface
func (e *FieldNotFoundError) Error() string {
	return fmt.Sprintf("field not found: %s", e.Field)
}

// NewFieldNotFoundError creates a new FieldNotFoundError
func NewFieldNotFoundError(field string) *FieldNotFoundError {
	return &FieldNotFoundError{Field: field}
}

// File implements the GGUFFile interface
type File struct {
	file *parser.GGUFFile
}

// GetParameters returns the model parameters (raw count, formatted string, error)
func (g *File) GetParameters() (float64, string, error) {
	if g.file == nil {
		return 0, "", fmt.Errorf("file is nil")
	}

	// size_label is the human-readable size of the model
	sizeLabel, found := g.file.Header.MetadataKV.Get("general.size_label")
	if found {
		formattedValue := sizeLabel.ValueString()
		// Parse the formatted value to get the raw value
		rawValue := parseParameters(formattedValue)
		if rawValue != 0 { // Skip non-numeric size labels (e.g. "large" in mxbai-embed-large-v1)
			return rawValue, formattedValue, nil
		}
	}

	// If no size label is found, use the parameters which is the exact number of parameters in the model
	paramsStr := g.file.Metadata().Parameters.String()
	if paramsStr == "" {
		return 0, "", NewFieldNotFoundError("parameters")
	}

	formattedValue := strings.TrimSpace(g.file.Metadata().Parameters.String())
	rawValue := parseParameters(formattedValue)
	return rawValue, formattedValue, nil
}

// GetArchitecture returns the model architecture (raw string, formatted string, error)
func (g *File) GetArchitecture() (string, string, error) {
	if g.file == nil {
		return "", "", fmt.Errorf("file is nil")
	}

	architecture := g.file.Metadata().Architecture
	if architecture == "" {
		return "", "", NewFieldNotFoundError("architecture")
	}

	rawValue := architecture
	formattedValue := strings.TrimSpace(rawValue)
	return rawValue, formattedValue, nil
}

// GetQuantization returns the model quantization (raw string, formatted string, error)
func (g *File) GetQuantization() (string, string, error) {
	if g.file == nil {
		return "", "", fmt.Errorf("file is nil")
	}

	fileTypeStr := g.file.Metadata().FileType.String()
	if fileTypeStr == "" {
		return "", "", NewFieldNotFoundError("file_type")
	}

	rawValue := fileTypeStr
	formattedValue := strings.TrimSpace(rawValue)
	return rawValue, formattedValue, nil
}

// GetSize returns the model size (raw bytes, formatted string, error)
func (g *File) GetSize() (int64, string, error) {
	if g.file == nil {
		return 0, "", fmt.Errorf("file is nil")
	}

	sizeStr := g.file.Metadata().Size.String()
	if sizeStr == "" {
		return 0, "", NewFieldNotFoundError("size")
	}

	// Parse the size string to get the raw value in bytes
	// The size string is typically in the format "123.45 MB" or similar
	rawValue := int64(0)
	formattedValue := sizeStr

	// Extract the numeric part and convert to bytes
	parts := strings.Fields(sizeStr)
	if len(parts) >= 2 {
		value, err := strconv.ParseFloat(parts[0], 64)
		if err == nil {
			unit := strings.ToUpper(parts[1])
			switch {
			case strings.HasPrefix(unit, "B"):
				rawValue = int64(value)
			case strings.HasPrefix(unit, "KB") || strings.HasPrefix(unit, "K"):
				rawValue = int64(value * 1024)
			case strings.HasPrefix(unit, "MB") || strings.HasPrefix(unit, "M"):
				rawValue = int64(value * 1024 * 1024)
			case strings.HasPrefix(unit, "GB") || strings.HasPrefix(unit, "G"):
				rawValue = int64(value * 1024 * 1024 * 1024)
			case strings.HasPrefix(unit, "TB") || strings.HasPrefix(unit, "T"):
				rawValue = int64(value * 1024 * 1024 * 1024 * 1024)
			}
		}
	}

	return rawValue, formattedValue, nil
}

// GetContextLength returns the model context length (raw length, formatted string, error)
func (g *File) GetContextLength() (uint32, string, error) {
	if g.file == nil {
		return 0, "", fmt.Errorf("file is nil")
	}

	architecture, found := g.file.Header.MetadataKV.Get("general.architecture")
	if !found {
		return 0, "", NewFieldNotFoundError("general.architecture")
	}

	contextLength, found := g.file.Header.MetadataKV.Get(architecture.ValueString() + ".context_length")
	if !found {
		return 0, "", NewFieldNotFoundError(architecture.ValueString() + ".context_length")
	}

	rawValue := contextLength.ValueUint32()
	formattedValue := fmt.Sprintf("%d", rawValue)
	return rawValue, formattedValue, nil
}

// GetVRAM returns the estimated VRAM requirements (raw GB, formatted string, error)
func (g *File) GetVRAM() (float64, string, error) {
	if g.file == nil {
		return 0, "", fmt.Errorf("file is nil")
	}

	// Get parameter count
	params, _, err := g.GetParameters()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get parameters: %w", err)
	}

	// Determine quantization
	_, quantFormatted, err := g.GetQuantization()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get quantization: %w", err)
	}

	var bytesPerParam float64
	switch {
	case strings.Contains(quantFormatted, "F16"):
		bytesPerParam = 2
	case strings.Contains(quantFormatted, "Q8"):
		bytesPerParam = 1
	case strings.Contains(quantFormatted, "Q5"):
		bytesPerParam = 0.68
	case strings.Contains(quantFormatted, "Q4"):
		bytesPerParam = 0.6
	default:
		// Fail if we don't know the bytes per parameter
		return 0, "", fmt.Errorf("unknown quantization: %s", quantFormatted)
	}

	// Get architecture prefix for metadata lookups
	_, archFormatted, err := g.GetArchitecture()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get architecture: %w", err)
	}

	// Extract KV cache dimensions
	nLayer, found := g.file.Header.MetadataKV.Get(archFormatted + ".block_count")
	if !found {
		return 0, "", NewFieldNotFoundError(archFormatted + ".block_count")
	}
	nEmb, found := g.file.Header.MetadataKV.Get(archFormatted + ".embedding_length")
	if !found {
		return 0, "", NewFieldNotFoundError(archFormatted + ".embedding_length")
	}

	// Get context length
	contextLength, _, err := g.GetContextLength()
	if err != nil {
		return 0, "", fmt.Errorf("failed to get context length: %w", err)
	}

	// Calculate model weights size
	modelSizeGB := (params * bytesPerParam) / (1024 * 1024 * 1024)
	// Calculate KV cache size
	kvCacheBytes := contextLength * nLayer.ValueUint32() * nEmb.ValueUint32() * 2 * 2
	kvCacheGB := float64(kvCacheBytes) / (1024 * 1024 * 1024)

	// Total VRAM estimate with 20% overhead
	totalVRAM := (modelSizeGB + kvCacheGB) * 1.2
	formattedValue := fmt.Sprintf("%.2f GB", totalVRAM)
	return totalVRAM, formattedValue, nil
}

// parseParameters converts parameter string to float64
func parseParameters(paramStr string) float64 {
	// Remove any non-numeric characters except decimal point
	toParse := strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || r == '.' {
			return r
		}
		return -1
	}, paramStr)

	// Parse the number
	params, err := strconv.ParseFloat(toParse, 64)
	if err != nil {
		return 0
	}

	// Convert to actual number of parameters (e.g., 1.24B -> 1.24e9)
	if strings.Contains(strings.ToUpper(paramStr), "B") {
		params *= 1e9
	} else if strings.Contains(strings.ToUpper(paramStr), "M") {
		params *= 1e6
	}

	return params
}
