package domain

import (
	"github.com/docker/model-cards/tools/build-tables/types"
)

// ModelVariant represents a single model variant with its properties
type ModelVariant struct {
	RepoName      string
	Tag           string
	Architecture  string
	Parameters    string
	Quantization  string
	Size          string
	IsLatest      bool
	ContextLength uint32
	VRAM          float64
	Descriptor    types.ModelDescriptor
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
