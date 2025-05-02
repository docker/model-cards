# Model Cart Updater

Automatically updates the "Available model variants" tables in model card markdown files based on the characteristics of OCI Model Artifacts.

## Features

- Scans the `ai/` directory for markdown files
- For each model, fetches OCI manifest information
- Locates GGUF files in the manifest via mediaType
- Extracts metadata from GGUF files without downloading the entire file
- Updates the "Available model variants" table in each markdown file

## Installation

```bash
go mod tidy
make build
```

## Usage

You can use the provided Makefile to build and run the application:

```bash
# Build the Go application
make build

# Run the application
make run

# Clean up the binary
make clean
```

Or you can run the binary directly if it's already built:

```bash
./bin/build-tables
```

This will scan all markdown files in the `ai/` directory and update their "Available model variants" tables.

## Implementation Details

### Domain Models and Interfaces

The application uses a clean architecture approach with well-defined interfaces:

- `RegistryClient`: Interacts with OCI registries to fetch model information
- `MarkdownUpdater`: Updates markdown files with model information
- `GGUFParser`: Parses GGUF files to extract metadata
- `ModelProcessor`: Processes model files

### OCI Registry Interaction

The application uses `github.com/google/go-containerregistry` to:
- List tags for a repository
- Fetch manifests
- Identify layers by mediaType
- Access layer content without downloading the entire file

### GGUF Metadata Extraction

The application uses `github.com/gpustack/gguf-parser-go` to:
- Parse GGUF headers and metadata without downloading the entire file
- Extract parameters, quantization, and other relevant information

### Markdown File Processing

The application:
- Finds the "Available model variants" section
- Generates a new table with the extracted information
- Updates the file with the new table
- Preserves the rest of the file content

## License

Same as the parent project.
