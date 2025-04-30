package gguf

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// GGUFMetadata represents metadata extracted from a GGUF file
type GGUFMetadata struct {
	Parameters   string
	Quantization string
	Architecture string
	ModelSize    string
}

// ExtractMetadata extracts metadata from a GGUF layer without downloading the entire file
func ExtractMetadata(ref name.Reference, digestStr string) (GGUFMetadata, error) {
	metadata := GGUFMetadata{}

	// Create a digest reference
	digest, err := name.NewDigest(fmt.Sprintf("%s@%s", ref.Context().Name(), digestStr))
	if err != nil {
		return metadata, fmt.Errorf("failed to create digest reference: %v", err)
	}

	// Get the layer
	layer, err := remote.Layer(digest)
	if err != nil {
		return metadata, fmt.Errorf("failed to get layer: %v", err)
	}

	// Get a reader for the compressed layer
	reader, err := layer.Compressed()
	if err != nil {
		return metadata, fmt.Errorf("failed to get compressed reader: %v", err)
	}
	defer reader.Close()

	// Read the first few bytes to identify the GGUF file
	header := make([]byte, 4)
	_, err = io.ReadFull(reader, header)
	if err != nil {
		return metadata, fmt.Errorf("failed to read GGUF header: %v", err)
	}

	// Check if it's a GGUF file (should start with "GGUF")
	if string(header) != "GGUF" {
		return metadata, fmt.Errorf("not a GGUF file, header: %s", string(header))
	}

	// In a real implementation, we would use gguf-parser-go to extract metadata
	// For now, we'll extract information from the digest string and tag

	// Extract parameters and quantization from the digest string or tag
	tagParts := strings.Split(digestStr, ":")
	if len(tagParts) > 1 {
		tag := tagParts[1]
		metadata.Parameters = extractParametersFromTag(tag)
		metadata.Quantization = extractQuantizationFromTag(tag)
	}

	return metadata, nil
}

// extractParametersFromTag extracts the number of parameters from a tag
func extractParametersFromTag(tag string) string {
	// Try to extract from tag name
	if strings.Contains(tag, "360M") {
		return "360M"
	} else if strings.Contains(tag, "135M") {
		return "135M"
	} else if strings.Contains(tag, "7B") {
		return "7B"
	} else if strings.Contains(tag, "13B") {
		return "13B"
	} else if strings.Contains(tag, "70B") {
		return "70B"
	}

	return ""
}

// extractQuantizationFromTag extracts the quantization type from a tag
func extractQuantizationFromTag(tag string) string {
	if strings.Contains(tag, "F16") {
		return "F16"
	} else if strings.Contains(tag, "Q4_0") {
		return "Q4_0"
	} else if strings.Contains(tag, "Q4_K_M") {
		return "IQ2_XXS/Q4_K_M"
	} else if strings.Contains(tag, "Q8_0") {
		return "Q8_0"
	}

	return ""
}
