package types

import parser "github.com/gpustack/gguf-parser-go"

// ModelDescriptor represents the data of a Model
type ModelDescriptor interface {
	// GetParameters returns the model parameters (raw count, formatted string, error)
	GetParameters() (float64, string, error)

	// GetArchitecture returns the model architecture
	GetArchitecture() string

	// GetQuantization returns the model quantization
	GetQuantization() parser.GGUFFileType

	// GetSize returns the model size (raw bytes, formatted string, error)
	GetSize() (int64, string, error)

	// GetContextLength returns the model context length (raw length, formatted string, error)
	GetContextLength() (uint32, string, error)

	// GetVRAM returns the estimated VRAM requirements (bytes, error)
	GetVRAM() (float64, error)

	// GetMetadata returns the model metadata (map[string]string)
	GetMetadata() map[string]string
}
