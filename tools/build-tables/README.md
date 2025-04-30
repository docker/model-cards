# Model Tables Builder

This Go script automatically updates the "Available model variants" tables in model card markdown files based on the characteristics of OCI Artifacts that represent Large Language Models.

## Features

- Scans the `ai/` directory for markdown files
- For each model, fetches OCI manifest information using `go-containerregistry`
- Locates GGUF files in the manifest via mediaType
- Extracts metadata from GGUF files using `gguf-parser-go` without downloading the entire file
- Updates the "Available model variants" table in each markdown file

## Requirements

- Go 1.18 or higher
- Access to the OCI registry containing the model artifacts

## Installation

```bash
go mod tidy
go build -o build-tables
```

## Usage

```bash
./build-tables
```

This will scan all markdown files in the `ai/` directory and update their "Available model variants" tables.

## Implementation Details

### OCI Registry Interaction

The script uses `github.com/google/go-containerregistry` to:
- List tags for a repository
- Fetch manifests
- Identify layers by mediaType
- Access layer content without downloading the entire file

### GGUF Metadata Extraction

The script uses `github.com/gpustack/gguf-parser-go` to:
- Parse GGUF headers and metadata without downloading the entire file
- Extract parameters, quantization, and other relevant information

### Markdown File Processing

The script:
- Finds the "Available model variants" section
- Generates a new table with the extracted information
- Updates the file with the new table
- Preserves the rest of the file content

## License

Same as the parent project.
