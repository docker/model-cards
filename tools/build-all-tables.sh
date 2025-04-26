#!/usr/bin/env bash
set -euo pipefail

# Script to build tables for all models in the ai/ folder

echo "🔍 Finding all model readme files in ai/ folder..."
echo ""

# Check if the force flag is provided
FORCE=""
if [ "${1:-}" = "--force" ]; then
  FORCE="--force"
  echo "⚠️ Force mode enabled - will overwrite existing tables without prompting"
  echo ""
fi

# Count total models for progress tracking
TOTAL_MODELS=$(ls -1 ai/*.md | wc -l)
CURRENT=0

# Process each markdown file in the ai/ directory
for file in ai/*.md; do
  # Extract the model name from the filename (remove path and extension)
  model_name=$(basename "$file" .md)
  
  # Increment counter
  CURRENT=$((CURRENT + 1))
  
  # Display progress
  echo "==============================================="
  echo "🔄 Processing model $CURRENT/$TOTAL_MODELS: ai/$model_name"
  echo "==============================================="
  
  # Run the build-model-table script for this model
  if [ -n "$FORCE" ]; then
    ./tools/build-model-table.sh --force "ai/$model_name"
  else
    ./tools/build-model-table.sh "ai/$model_name"
  fi
  
  echo ""
  echo "✅ Completed ai/$model_name"
  echo ""
done

echo "==============================================="
echo "🎉 All model tables have been updated!"
echo "==============================================="
