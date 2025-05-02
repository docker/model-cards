package main

import (
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

func main() {
	// Parse command line flags
	logLevel := flag.String("log-level", "info", "Log level (debug, info, warn, error)")
	modelDir := flag.String("model-dir", "../../ai", "Directory containing model markdown files")
	modelFile := flag.String("model-file", "", "Specific model markdown file to update (without path)")
	flag.Parse()

	// Configure logger
	switch *logLevel {
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

	logger.Info("Starting model-cards updater")
	logger.Debugf("Log level set to: %s", *logLevel)

	// Create dependencies
	registryClient := registry.NewClient()
	markdownUpdater := markdown.NewUpdater()

	// Create the application
	app := NewApplication(registryClient, markdownUpdater, *modelDir, *modelFile)

	// Run the application
	if err := app.Run(); err != nil {
		logger.WithError(err).Errorf("Application failed: %v", err)
		os.Exit(1)
	}

	logger.Info("Application completed successfully")
}
