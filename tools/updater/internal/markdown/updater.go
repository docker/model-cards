package markdown

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/docker/model-cards/tools/build-tables/internal/domain"
)

// Updater implements the domain.MarkdownUpdater interface
type Updater struct{}

// NewUpdater creates a new markdown updater
func NewUpdater() *Updater {
	return &Updater{}
}

// UpdateModelTable updates the "Available model variants" table in a markdown file
func (u *Updater) UpdateModelTable(filePath string, variants []domain.ModelVariant) error {
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
		formattedParams := domain.FormatParameters(variant.Parameters)

		// Format the context window
		contextWindow := "-"
		if variant.ContextLength > 0 {
			contextWindow = fmt.Sprintf("%d tokens", variant.ContextLength)
		}

		// Create the table row
		row := fmt.Sprintf("| %s | %s | %s | %s | - | %s |\n",
			modelVariant,
			formattedParams,
			variant.Quantization,
			contextWindow,
			variant.Size)
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

	return nil
}
