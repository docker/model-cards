package markdown

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/docker/model-cards/tools/build-tables/registry"
)

// UpdateModelTable updates the "Available model variants" table in a markdown file
func UpdateModelTable(filePath string, variants []registry.ModelVariant) error {
	// Read the markdown file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read markdown file: %v", err)
	}

	// Find the "Available model variants" section
	sectionRegex := regexp.MustCompile(`(?m)^## Available model variants\s*$`)
	sectionMatch := sectionRegex.FindIndex(content)
	if sectionMatch == nil {
		return fmt.Errorf("could not find the 'Available model variants' section")
	}

	// Find the next section after "Available model variants"
	nextSectionRegex := regexp.MustCompile(`(?m)^##\s+[^#]`)
	nextSectionMatch := nextSectionRegex.FindIndex(content[sectionMatch[1]:])

	var endOfTableSection int
	if nextSectionMatch != nil {
		endOfTableSection = sectionMatch[1] + nextSectionMatch[0]
	} else {
		endOfTableSection = len(content)
	}

	// Extract the content before and after the table section
	beforeTable := content[:sectionMatch[1]]
	afterTable := content[endOfTableSection:]

	// Generate the new table
	var tableBuilder strings.Builder
	tableBuilder.WriteString("\n\n")
	tableBuilder.WriteString("| Model variant | Parameters | Quantization | Context window | VRAM | Size |\n")
	tableBuilder.WriteString("|---------------|------------|--------------|----------------|------|-------|\n")

	// Add all the rows
	var latestTag string
	for _, variant := range variants {
		// Format the model variant
		var modelVariant string
		if variant.IsLatest {
			modelVariant = fmt.Sprintf("`%s:latest`<br><br>`%s:%s`", variant.RepoName, variant.RepoName, variant.Tag)
			latestTag = variant.Tag
		} else {
			modelVariant = fmt.Sprintf("`%s:%s`", variant.RepoName, variant.Tag)
		}

		// Format the parameters
		formattedParams := registry.FormatParameters(variant.Parameters)

		// Format the size
		formattedSize := registry.FormatSize(variant.SizeMB)

		// Create the table row
		row := fmt.Sprintf("| %s | %s | %s | - | - | %s |\n", modelVariant, formattedParams, variant.Quantization, formattedSize)
		tableBuilder.WriteString(row)
	}

	// Add the footnote for VRAM estimation
	tableBuilder.WriteString("\n¹: VRAM estimation.\n")

	// Add the latest tag mapping note if we found a match
	if latestTag != "" {
		tableBuilder.WriteString(fmt.Sprintf("\n> `:latest` → `%s`\n", latestTag))
	}

	// Combine the parts
	newContent := append(beforeTable, []byte(tableBuilder.String())...)
	newContent = append(newContent, afterTable...)

	// Write the updated content back to the file
	err = os.WriteFile(filePath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated markdown file: %v", err)
	}

	fmt.Printf("✅ Successfully updated %s with all variants for %s\n", filePath, variants[0].RepoName)
	return nil
}
