package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/model-cards/tools/build-tables/internal/domain"
	"github.com/docker/model-cards/tools/build-tables/internal/logger"
	"github.com/docker/model-cards/tools/build-tables/internal/markdown"
	"github.com/docker/model-cards/tools/build-tables/internal/registry"
	"github.com/docker/model-cards/tools/build-tables/internal/utils"
	"github.com/sirupsen/logrus"
)

// Application encapsulates the main application logic
type Application struct {
	registryClient  domain.RegistryClient
	markdownUpdater domain.MarkdownUpdater
	modelDir        string
	modelFile       string
}

// NewApplication creates a new application instance
func NewApplication(registryClient domain.RegistryClient, markdownUpdater domain.MarkdownUpdater, modelDir string, modelFile string) *Application {
	return &Application{
		registryClient:  registryClient,
		markdownUpdater: markdownUpdater,
		modelDir:        modelDir,
		modelFile:       modelFile,
	}
}

// Run executes the main application logic
func (a *Application) Run() error {
	var files []string
	var err error

	// Check if a specific model file is requested
	if a.modelFile != "" {
		// Process only the specified model file
		modelFilePath := filepath.Join(a.modelDir, a.modelFile)
		if !utils.FileExists(modelFilePath) {
			err := fmt.Errorf("model file '%s' does not exist", modelFilePath)
			logger.WithField("file", modelFilePath).Error("model file does not exist")
			return err
		}

		logger.Infof("🔍 Processing single model file: %s", a.modelFile)
		files = []string{modelFilePath}
	} else {
		// Process all model files in the directory
		logger.Info("🔍 Finding all model readme files in ai/ folder...")

		// Find all markdown files in the model directory
		files, err = markdown.FindMarkdownFiles(a.modelDir)
		if err != nil {
			logger.WithError(err).Error("error finding model files")
			return err
		}

		logger.Infof("Found %d model files", len(files))
	}

	// Count total models for progress tracking
	totalModels := len(files)
	current := 0

	// Process each markdown file
	for _, file := range files {
		// Extract the model name from the filename
		modelName := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))

		// Increment counter
		current++

		// Display progress
		logger.Info("===============================================")
		logger.Infof("🔄 Processing model %d/%d: %s/%s", current, totalModels, filepath.Base(a.modelDir), modelName)
		logger.Info("===============================================")

		// Process the model file
		err := a.processModelFile(file)
		if err != nil {
			logger.WithFields(logger.Fields{
				"model": modelName,
				"error": err,
			}).Error("Error processing model")
			continue
		} else {
			logger.WithField("model", modelName).Info("Successfully processed model")
		}

		logger.Infof("✅ Completed %s/%s", filepath.Base(a.modelDir), modelName)
	}

	logger.Info("===============================================")
	logger.Info("🎉 All model tables have been updated!")
	logger.Info("===============================================")

	return nil
}

// processModelFile processes a single model markdown file
func (a *Application) processModelFile(filePath string) error {
	// Extract the repository name from the file path
	repoName := utils.GetRepositoryName(filePath, filepath.Dir(a.modelDir))

	logger.WithField("file", filePath).Info("📄 Using readme file")

	// Check if the file exists
	if !utils.FileExists(filePath) {
		err := fmt.Errorf("readme file '%s' does not exist", filePath)
		logger.WithField("file", filePath).Error("readme file does not exist")
		return err
	}

	// List all tags for the repository
	logger.WithField("repository", repoName).Info("📦 Listing tags for repository")
	tags, err := a.registryClient.ListTags(repoName)
	if err != nil {
		logger.WithFields(logger.Fields{
			"repository": repoName,
			"error":      err,
		}).Error("error listing tags")
		return fmt.Errorf("error listing tags: %v", err)
	}

	// Process each tag and collect model variants
	variants, err := a.registryClient.ProcessTags(repoName, tags)
	if err != nil {
		logger.WithFields(logger.Fields{
			"repository": repoName,
			"error":      err,
		}).Error("error processing tags")
		return fmt.Errorf("error processing tags: %v", err)
	}

	// Update the markdown file with the new table
	err = a.markdownUpdater.UpdateModelTable(filePath, variants)
	if err != nil {
		logger.WithFields(logger.Fields{
			"file":  filePath,
			"error": err,
		}).Error("error updating markdown file")
		return fmt.Errorf("error updating markdown file: %v", err)
	}

	return nil
}

// ModelInspector encapsulates the model inspection logic
type ModelInspector struct {
	registryClient domain.RegistryClient
	repository     string
	tag            string
	showAll        bool
	showParams     bool
	showArch       bool
	showQuant      bool
	showSize       bool
	showContext    bool
	showVRAM       bool
	formatJSON     bool
}

// NewModelInspector creates a new model inspector
func NewModelInspector(registryClient domain.RegistryClient, repository, tag string, options map[string]bool) *ModelInspector {
	return &ModelInspector{
		registryClient: registryClient,
		repository:     repository,
		tag:            tag,
		showAll:        options["all"],
		showParams:     options["parameters"],
		showArch:       options["architecture"],
		showQuant:      options["quantization"],
		showSize:       options["size"],
		showContext:    options["context"],
		showVRAM:       options["vram"],
		formatJSON:     options["json"],
	}
}

// Run executes the model inspection
func (m *ModelInspector) Run() error {
	// If no specific options are selected, show all
	if !m.showParams && !m.showArch && !m.showQuant && !m.showSize && !m.showContext && !m.showVRAM {
		m.showAll = true
	}

	// If showAll is true, enable all options
	if m.showAll {
		m.showParams = true
		m.showArch = true
		m.showQuant = true
		m.showSize = true
		m.showContext = true
		m.showVRAM = true
	}

	// If a specific tag is provided, inspect only that tag
	if m.tag != "" {
		return m.inspectTag(m.repository, m.tag)
	}

	// Otherwise, list all tags and inspect each one
	tags, err := m.registryClient.ListTags(m.repository)
	if err != nil {
		return fmt.Errorf("failed to list tags: %v", err)
	}

	logger.Infof("Found %d tags for repository %s", len(tags), m.repository)

	// If JSON output is requested, collect all results in a map
	if m.formatJSON {
		results := make(map[string]interface{})
		for _, tag := range tags {
			variant, err := m.registryClient.GetModelVariant(context.Background(), m.repository, tag)
			if err != nil {
				logger.Warnf("Failed to get info for %s:%s: %v", m.repository, tag, err)
				continue
			}
			results[tag] = m.variantToMap(variant)
		}

		// Output as JSON
		jsonData, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Otherwise, output in text format
	for _, tag := range tags {
		if err := m.inspectTag(m.repository, tag); err != nil {
			logger.Warnf("Failed to inspect %s:%s: %v", m.repository, tag, err)
		}
		fmt.Println("----------------------------------------")
	}

	return nil
}

// inspectTag inspects a specific tag and outputs the requested information
func (m *ModelInspector) inspectTag(repository, tag string) error {
	logger.Infof("Inspecting %s:%s", repository, tag)

	// Get model variant information
	variant, err := m.registryClient.GetModelVariant(context.Background(), repository, tag)
	if err != nil {
		return fmt.Errorf("failed to get model variant: %v", err)
	}

	// If JSON output is requested, output as JSON
	if m.formatJSON {
		result := m.variantToMap(variant)
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err)
		}
		fmt.Println(string(jsonData))
		return nil
	}

	// Otherwise, output in text format
	fmt.Printf("🔍 Model: %s:%s\n", repository, tag)

	if m.showParams {
		fmt.Printf("   • Parameters   : %s\n", variant.Parameters)
	}

	if m.showArch {
		// Architecture is not directly stored in the variant, but we can try to infer it
		fmt.Printf("   • Architecture : %s\n", inferArchitecture(repository))
	}

	if m.showQuant {
		fmt.Printf("   • Quantization : %s\n", variant.Quantization)
	}

	if m.showSize {
		fmt.Printf("   • Size         : %s\n", variant.Size)
	}

	if m.showContext {
		if variant.ContextLength > 0 {
			fmt.Printf("   • Context      : %d tokens\n", variant.ContextLength)
		} else {
			fmt.Printf("   • Context      : Unknown\n")
		}
	}

	if m.showVRAM {
		if variant.VRAM > 0 {
			fmt.Printf("   • VRAM         : %.2f GB\n", variant.VRAM)
		} else {
			fmt.Printf("   • VRAM         : Unknown\n")
		}
	}

	return nil
}

// variantToMap converts a ModelVariant to a map for JSON output
func (m *ModelInspector) variantToMap(variant domain.ModelVariant) map[string]interface{} {
	result := make(map[string]interface{})

	if m.showParams {
		result["parameters"] = variant.Parameters
	}

	if m.showArch {
		result["architecture"] = inferArchitecture(variant.RepoName)
	}

	if m.showQuant {
		result["quantization"] = variant.Quantization
	}

	if m.showSize {
		result["size"] = variant.Size
	}

	if m.showContext {
		if variant.ContextLength > 0 {
			result["context_length"] = variant.ContextLength
		} else {
			result["context_length"] = nil
		}
	}

	if m.showVRAM {
		if variant.VRAM > 0 {
			result["vram_gb"] = variant.VRAM
		} else {
			result["vram_gb"] = nil
		}
	}

	return result
}

// inferArchitecture tries to infer the architecture from the repository name
func inferArchitecture(repository string) string {
	repoName := filepath.Base(repository)

	switch {
	case strings.Contains(repoName, "llama"):
		return "llama"
	case strings.Contains(repoName, "mistral"):
		return "mistral"
	case strings.Contains(repoName, "phi"):
		return "phi"
	case strings.Contains(repoName, "gemma"):
		return "gemma"
	case strings.Contains(repoName, "qwen"):
		return "qwen"
	case strings.Contains(repoName, "deepseek"):
		return "deepseek"
	case strings.Contains(repoName, "smollm"):
		return "smollm"
	default:
		return "unknown"
	}
}

func main() {
	// Define command flags
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	inspectCmd := flag.NewFlagSet("inspect-model", flag.ExitOnError)

	// Update command flags
	updateLogLevel := updateCmd.String("log-level", "info", "Log level (debug, info, warn, error)")
	updateModelDir := updateCmd.String("model-dir", "../../ai", "Directory containing model markdown files")
	updateModelFile := updateCmd.String("model-file", "", "Specific model markdown file to update (without path)")

	// Inspect command flags
	inspectLogLevel := inspectCmd.String("log-level", "info", "Log level (debug, info, warn, error)")
	inspectTag := inspectCmd.String("tag", "", "Specific tag to inspect")
	inspectAll := inspectCmd.Bool("all", false, "Show all metadata")
	inspectParams := inspectCmd.Bool("parameters", false, "Show parameters")
	inspectArch := inspectCmd.Bool("architecture", false, "Show architecture")
	inspectQuant := inspectCmd.Bool("quantization", false, "Show quantization")
	inspectSize := inspectCmd.Bool("size", false, "Show size")
	inspectContext := inspectCmd.Bool("context", false, "Show context length")
	inspectVRAM := inspectCmd.Bool("vram", false, "Show VRAM requirements")
	inspectJSON := inspectCmd.Bool("json", false, "Output in JSON format")

	// Check if a command is provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'update' or 'inspect-model' subcommand")
		fmt.Println("Usage:")
		fmt.Println("  model-cards-cli update [options]")
		fmt.Println("  model-cards-cli inspect-model [options] REPOSITORY")
		os.Exit(1)
	}

	// Configure logger based on the command
	var logLevel string

	// Parse the appropriate command
	switch os.Args[1] {
	case "update":
		updateCmd.Parse(os.Args[2:])
		logLevel = *updateLogLevel
	case "inspect-model":
		inspectCmd.Parse(os.Args[2:])
		logLevel = *inspectLogLevel
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Expected 'update' or 'inspect-model' subcommand")
		os.Exit(1)
	}

	// Configure logger
	switch logLevel {
	case "debug":
		logger.Log.SetLevel(logrus.DebugLevel)
	case "info":
		logger.Log.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.Log.SetLevel(logrus.WarnLevel)
	case "error":
		logger.Log.SetLevel(logrus.ErrorLevel)
	default:
		logger.Log.SetLevel(logrus.InfoLevel)
	}

	logger.Debugf("Log level set to: %s", logLevel)

	// Create dependencies
	registryClient := registry.NewClient()

	// Execute the appropriate command
	if updateCmd.Parsed() {
		logger.Info("Starting model-cards updater")

		markdownUpdater := markdown.NewUpdater()
		app := NewApplication(registryClient, markdownUpdater, *updateModelDir, *updateModelFile)

		if err := app.Run(); err != nil {
			logger.WithError(err).Errorf("Application failed: %v", err)
			os.Exit(1)
		}

		logger.Info("Application completed successfully")
	} else if inspectCmd.Parsed() {
		logger.Info("Starting model inspector")

		// Check if a repository is provided
		args := inspectCmd.Args()
		if len(args) < 1 {
			fmt.Println("Error: Repository argument is required")
			fmt.Println("Usage: updater inspect-model [options] REPOSITORY")
			os.Exit(1)
		}

		repository := args[0]

		// Create options map
		options := map[string]bool{
			"all":          *inspectAll,
			"parameters":   *inspectParams,
			"architecture": *inspectArch,
			"quantization": *inspectQuant,
			"size":         *inspectSize,
			"context":      *inspectContext,
			"vram":         *inspectVRAM,
			"json":         *inspectJSON,
		}

		inspector := NewModelInspector(registryClient, repository, *inspectTag, options)

		if err := inspector.Run(); err != nil {
			logger.WithError(err).Errorf("Inspection failed: %v", err)
			os.Exit(1)
		}

		logger.Info("Inspection completed successfully")
	}
}
