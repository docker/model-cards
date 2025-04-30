package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/model-cards/tools/build-tables/markdown"
	"github.com/docker/model-cards/tools/build-tables/registry"
)

func main() {
	fmt.Println("🔍 Finding all model readme files in ai/ folder...")
	fmt.Println("")

	// Find all markdown files in the ai/ directory
	// Use the correct path relative to the current working directory
	files, err := filepath.Glob("../../ai/*.md")
	if err != nil {
		fmt.Printf("Error finding model files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d model files\n", len(files))

	// Count total models for progress tracking
	totalModels := len(files)
	current := 0

	// Process each markdown file in the ai/ directory
	for _, file := range files {
		// Extract the model name from the filename (remove path and extension)
		modelName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		// Increment counter
		current++

		// Display progress
		fmt.Println("===============================================")
		fmt.Printf("🔄 Processing model %d/%d: ai/%s\n", current, totalModels, modelName)
		fmt.Println("===============================================")

		// Process the model file
		err := processModelFile(file)
		if err != nil {
			fmt.Printf("Error processing model %s: %v\n", modelName, err)
			continue
		} else {
			fmt.Printf("Successfully processed model %s\n", modelName)
		}

		fmt.Println("")
		fmt.Printf("✅ Completed ai/%s\n", modelName)
		fmt.Println("")
	}

	fmt.Println("===============================================")
	fmt.Println("🎉 All model tables have been updated!")
	fmt.Println("===============================================")
}

// processModelFile processes a single model markdown file
func processModelFile(filePath string) error {
	// Extract the repository name from the file path
	// Convert the path to be relative to the project root
	relPath := strings.TrimPrefix(filePath, "../../")
	repoName := strings.TrimSuffix(relPath, filepath.Ext(filePath))

	fmt.Printf("📄 Using readme file: %s\n", filePath)

	// Check if the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("readme file '%s' does not exist", filePath)
	}

	// List all tags for the repository
	fmt.Printf("📦 Listing tags for repository: %s\n", repoName)
	tags, err := registry.ListTags(repoName)
	if err != nil {
		return fmt.Errorf("error listing tags: %v", err)
	}

	// Process each tag and collect model variants
	variants, err := registry.ProcessTags(repoName, tags)
	if err != nil {
		return fmt.Errorf("error processing tags: %v", err)
	}

	// Update the markdown file with the new table
	err = markdown.UpdateModelTable(filePath, variants)
	if err != nil {
		return fmt.Errorf("error updating markdown file: %v", err)
	}

	return nil
}
