package domain

import "context"

// ModelVariant represents a single model variant with its properties
type ModelVariant struct {
	RepoName      string
	Tag           string
	Parameters    string
	Quantization  string
	Size          string
	IsLatest      bool
	ContextLength uint32
}

// RegistryClient defines the interface for interacting with model registries
type RegistryClient interface {
	// ListTags lists all tags for a repository
	ListTags(repoName string) ([]string, error)

	// ProcessTags processes all tags for a repository and returns model variants
	ProcessTags(repoName string, tags []string) ([]ModelVariant, error)

	// GetModelVariant gets information about a specific model tag
	GetModelVariant(ctx context.Context, repoName, tag string) (ModelVariant, error)
}

// MarkdownUpdater defines the interface for updating markdown files
type MarkdownUpdater interface {
	// UpdateModelTable updates the "Available model variants" table in a markdown file
	UpdateModelTable(filePath string, variants []ModelVariant) error
}

// ModelProcessor defines the interface for processing model files
type ModelProcessor interface {
	// ProcessModelFile processes a single model markdown file
	ProcessModelFile(filePath string) error
}

// GGUFParser defines the interface for parsing GGUF files
type GGUFParser interface {
	// ParseRemote parses a remote GGUF file
	ParseRemote(ctx context.Context, url, token string) (GGUFFile, error)
}

// GGUFFile represents the metadata from a GGUF file
type GGUFFile interface {
	// GetParameters returns the model parameters
	GetParameters() string

	// GetArchitecture returns the model architecture
	GetArchitecture() string

	// GetQuantization returns the model quantization
	GetQuantization() string

	// GetSize returns the model size
	GetSize() string

	// GetContextLength returns the model context length
	GetContextLength() uint32
}
