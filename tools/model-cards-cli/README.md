# Model Cards CLI

A command-line tool for working with model cards. It can update the "Available model variants" tables in model card markdown files and inspect model repositories to extract metadata.

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

The Model Cards CLI provides three main commands:

1. `update` - Updates the "Available model variants" tables in model card markdown files
2. `inspect-model` - Inspects a model repository and displays metadata about the model variants
3. `upload-overview` - Uploads an overview to Docker Hub for a specified repository

### Update Command

You can use the provided Makefile to build and run the application:

```bash
# Build the Go application
make build

# Run the application to update all model files
make run

# Run the application to update a single model file
make run-single MODEL=<model-file.md>

# Clean up the binary
make clean
```

Or you can run the binary directly if it's already built:

```bash
# Update all model files
./bin/model-cards-cli update

# Update a specific model file
./bin/model-cards-cli update --model-file=<model-file.md>
```

By default, the tool will scan all markdown files in the `ai/` directory and update their "Available model variants" tables. If you specify a model file with the `--model-file` flag or the `MODEL` parameter, it will only update that specific file.

#### Update Command Options

- `--model-dir`: Directory containing model markdown files (default: "../../ai")
- `--model-file`: Specific model markdown file to update (without path)
- `--log-level`: Log level (debug, info, warn, error) (default: "info")

### Inspect Model Command

The `inspect-model` command allows you to inspect a model repository and display metadata about the model variants. This is useful for getting information about a model without having to update the markdown files.

You can use the provided Makefile to run the inspect command:

```bash
# Inspect all tags in a repository
make inspect REPO=ai/smollm2

# Inspect a specific tag
make inspect REPO=ai/smollm2 TAG=360M-Q4_K_M

# Inspect with metadata
make inspect REPO=ai/smollm2 OPTIONS="--all"
```

Or you can run the binary directly if it's already built:

```bash
# Inspect all tags in a repository
./bin/model-cards-cli inspect-model ai/smollm2

# Inspect a specific tag
./bin/model-cards-cli inspect-model --tag=360M-Q4_K_M ai/smollm2

# Inspect with specific options
./bin/model-cards-cli inspect-model --all ai/smollm2
```

### Upload Overview Command

The `upload-overview` command allows you to upload a model overview (markdown content) to Docker Hub. This is useful for updating the repository description that appears on Docker Hub.

You can use the provided Makefile to run the upload-overview command:

```bash
# Upload an overview to Docker Hub
make upload-overview FILE=../../ai/llama3.1.md REPO=ai/llama3 TOKEN=your_token_here
```

Or you can run the binary directly if it's already built:

```bash
# Upload an overview to Docker Hub
./bin/model-cards-cli upload-overview --file=../../ai/llama3.1.md --repository=ai/llama3 --token=your_token_here
```

The command requires three parameters:
- `FILE` or `--file`: Path to the markdown file containing the overview content
- `REPO` or `--repository`: Repository to upload the overview to (format: namespace/repository)
- `TOKEN` or `--token`: Authentication token with repo:admin scope

#### Upload Overview Command Options

- `--file`: Path to the overview file to upload (required)
- `--repository`: Repository to upload the overview to in the format namespace/repository (required)
- `--token`: Authentication token with repo:admin scope (required)
- `--log-level`: Log level (debug, info, warn, error) (default: "info")

The command will read the specified markdown file and upload its content as the full description for the specified repository on Docker Hub. The API endpoint used is `https://api.docker.team/v2/namespaces/{namespace}/repositories/{repository}`.
