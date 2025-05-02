package gguf

import (
	"context"
	"fmt"
	"strings"

	parser "github.com/gpustack/gguf-parser-go"

	"github.com/docker/model-cards/tools/build-tables/internal/domain"
)

// Parser implements the domain.GGUFParser interface
type Parser struct{}

// NewParser creates a new GGUF parser
func NewParser() *Parser {
	return &Parser{}
}

// ParseRemote parses a remote GGUF file
func (p *Parser) ParseRemote(ctx context.Context, url, token string) (domain.GGUFFile, error) {
	gf, err := parser.ParseGGUFFileRemote(ctx, url, parser.UseBearerAuth(token))
	if err != nil {
		return nil, fmt.Errorf("failed to parse GGUF: %w", err)
	}

	return &File{
		file: gf,
	}, nil
}

// ParseLocal parses a local GGUF file
func (p *Parser) ParseLocal(path string) (domain.GGUFFile, error) {
	gf, err := parser.ParseGGUFFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GGUF: %w", err)
	}

	return &File{
		file: gf,
	}, nil
}

// File implements the domain.GGUFFile interface
type File struct {
	file *parser.GGUFFile
}

// GetParameters returns the model parameters
func (g *File) GetParameters() string {
	if g.file == nil {
		return ""
	}
	// size_label is the human-readable size of the model
	sizeLabel, found := g.file.Header.MetadataKV.Get("general.size_label")
	if found {
		return sizeLabel.ValueString()
	}

	// If no size label is found, use the parameters which is the exact number of parameters in the model
	return strings.TrimSpace(g.file.Metadata().Parameters.String())
}

// GetArchitecture returns the model architecture
func (g *File) GetArchitecture() string {
	if g.file == nil {
		return ""
	}
	return strings.TrimSpace(g.file.Metadata().Architecture)
}

// GetQuantization returns the model quantization
func (g *File) GetQuantization() string {
	if g.file == nil {
		return ""
	}
	return strings.TrimSpace(g.file.Metadata().FileType.String())
}

// GetSize returns the model size
func (g *File) GetSize() string {
	return g.file.Metadata().Size.String()
}

// GetContextLength returns the model context length
func (g *File) GetContextLength() uint32 {
	if g.file == nil {
		return 0
	}

	architecture, ok := g.file.Header.MetadataKV.Get("general.architecture")
	if !ok {
		return 0
	}

	contextLength, ok := g.file.Header.MetadataKV.Get(architecture.ValueString() + ".context_length")
	if !ok {
		return 0
	}

	return contextLength.ValueUint32()
}
