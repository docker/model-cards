package domain

import (
	"context"
	"github.com/docker/model-cards/tools/build-tables/types"
)

// ModelVariant represents a single model variant with its properties
type ModelVariant struct {
	RepoName      string
	Tag           string
	Parameters    string
	Quantization  string
	Size          string
	IsLatest      bool
	ContextLength uint32
	VRAM          float64
	GGUF          types.GGUFFile
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
