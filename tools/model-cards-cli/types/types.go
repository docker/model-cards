package types

// GGUFFile represents the metadata from a GGUF file
type GGUFFile interface {
	// GetParameters returns the model parameters (raw count, formatted string, error)
	GetParameters() (float64, string, error)

	// GetArchitecture returns the model architecture (raw string, formatted string, error)
	GetArchitecture() (string, string, error)

	// GetQuantization returns the model quantization (raw string, formatted string, error)
	GetQuantization() (string, string, error)

	// GetSize returns the model size (raw bytes, formatted string, error)
	GetSize() (int64, string, error)

	// GetContextLength returns the model context length (raw length, formatted string, error)
	GetContextLength() (uint32, string, error)

	// GetVRAM returns the estimated VRAM requirements (raw GB, formatted string, error)
	GetVRAM() (float64, string, error)
}
