.PHONY: upload-overview

# Define variables for upload-overview command
FILE ?=
REPO ?=
USERNAME ?=
TOKEN ?=

upload-overview:
	@if [ -z "$(FILE)" ]; then \
		echo "Error: FILE parameter is required."; \
		echo "Usage: make upload-overview FILE=<overview-file> REPO=<namespace/repository> USERNAME=<username> TOKEN=<token>"; \
		exit 1; \
	fi
	@if [ -z "$(REPO)" ]; then \
		echo "Error: REPO parameter is required."; \
		echo "Usage: make upload-overview FILE=<overview-file> REPO=<namespace/repository> USERNAME=<username> TOKEN=<token>"; \
		exit 1; \
	fi
	@if [ -z "$(USERNAME)" ]; then \
		echo "Error: USERNAME parameter is required."; \
		echo "Usage: make upload-overview FILE=<overview-file> REPO=<namespace/repository> USERNAME=<username> TOKEN=<token>"; \
		exit 1; \
	fi
	@if [ -z "$(TOKEN)" ]; then \
		echo "Error: TOKEN parameter is required."; \
		echo "Usage: make upload-overview FILE=<overview-file> REPO=<namespace/repository> USERNAME=<username> TOKEN=<token>"; \
		exit 1; \
	fi
	@echo "Building model-cards-cli..."
	@$(MAKE) -C tools/model-cards-cli build
	@echo "Uploading overview from $(FILE) to $(REPO)..."
	@tools/model-cards-cli/bin/model-cards-cli upload-overview --file="$(FILE)" --repository="$(REPO)" --username="$(USERNAME)" --token="$(TOKEN)"

help:
	@echo "Available targets:"
	@echo "  upload-overview  - Upload an overview to Docker Hub (Usage: make upload-overview FILE=<overview-file> REPO=<namespace/repository> USERNAME=<username> TOKEN=<token>)"
	@echo "                     Example: make upload-overview FILE=ai/llama3.1.md REPO=ai/llama3 USERNAME=your_username TOKEN=your_pat_here"
	@echo "  help             - Show this help message"