#!/usr/bin/env bash
set -euo pipefail

# Accept repository name as input
REPO="${1:-}"
if [ -z "$REPO" ]; then
  echo "Usage: $0 <repository-name>"
  echo "Example: $0 ai/smollm2"
  exit 1
fi

# Extract model name and namespace
MODEL_NAME=${REPO##*/}
NAMESPACE=${REPO%/*}
README_FILE="${NAMESPACE}/${MODEL_NAME}.md"

echo "📄 Using readme file: $README_FILE"
if [ ! -f "$README_FILE" ]; then
  echo "Error: Readme file '$README_FILE' does not exist."
  exit 1
fi

# List all tags for the repository
echo "📦 Listing tags for repository: $REPO"
TAGS=$(crane ls "$REPO")

# Create an array to store all rows
declare -a TABLE_ROWS

# Find which tag corresponds to latest
LATEST_TAG=""
LATEST_QUANT=""
LATEST_PARAMS=""

for TAG in $TAGS; do
  if [ "$TAG" = "latest" ]; then
    # Get info for the latest tag
    LATEST_INFO=$(./tools/inspect-model.sh "${REPO}:latest")
    LATEST_PARAMS=$(echo "$LATEST_INFO" | grep "Parameters" | sed -E 's/.*: (.+)$/\1/' | tr -d ' ')
    LATEST_QUANT=$(echo "$LATEST_INFO" | grep "Quantization" | sed -E 's/.*: (.+)$/\1/' | tr -d ' ')
    break
  fi
done

# Process each tag
for TAG in $TAGS; do
  # Skip the latest tag - we'll handle it specially
  if [ "$TAG" = "latest" ]; then
    continue
  fi
  
  MODEL_REF="${REPO}:${TAG}"
  echo "🔍 Processing tag: $TAG"
  
  # Run inspect-model.sh to get model information
  MODEL_INFO=$(./tools/inspect-model.sh "$MODEL_REF")
  
  # Extract information from the output
  PARAMETERS=$(echo "$MODEL_INFO" | grep "Parameters" | sed -E 's/.*: (.+)$/\1/' | tr -d ' ')
  QUANTIZATION=$(echo "$MODEL_INFO" | grep "Quantization" | sed -E 's/.*: (.+)$/\1/' | tr -d ' ')
  
  # Extract both MB and GB sizes from the output
  MB_SIZE=$(echo "$MODEL_INFO" | grep "Artifact Size" | sed -E 's/.*: .* \((.+) MB \/ .+\)$/\1/' | tr -d ' ')
  GB_SIZE=$(echo "$MODEL_INFO" | grep "Artifact Size" | sed -E 's/.*: .* \(.+ MB \/ (.+) GB\)$/\1/' | tr -d ' ')
  
  # Decide which unit to use based on the size
  if (( $(echo "$MB_SIZE >= 1000" | bc -l) )); then
    FORMATTED_SIZE="${GB_SIZE} GB"
  else
    FORMATTED_SIZE="${MB_SIZE} MB"
  fi
  
  # Format the parameters to match the table format
  if [[ "$TAG" == *"360M"* ]]; then
    # For 360M models, use "360M" for consistency
    FORMATTED_PARAMS="360M"
  elif [[ "$TAG" == *"135M"* ]]; then
    # For 135M models, use "135M" for consistency
    FORMATTED_PARAMS="135M"
  elif [[ "$PARAMETERS" == *"M"* ]]; then
    FORMATTED_PARAMS="$PARAMETERS"
  elif [[ "$PARAMETERS" == *"B"* ]]; then
    FORMATTED_PARAMS="$PARAMETERS"
  else
    FORMATTED_PARAMS="$PARAMETERS"
  fi
  
  # Check if this tag matches the latest tag
  if [ -n "$LATEST_QUANT" ] && [ "$QUANTIZATION" = "$LATEST_QUANT" ] && [ "$PARAMETERS" = "$LATEST_PARAMS" ]; then
    # This is the tag that matches latest - create a special row
    MODEL_VARIANT="\`${REPO}:latest\`<br><br>\`${REPO}:${TAG}\`"
    # Save this tag for the latest mapping note
    LATEST_TAG="$TAG"
  else
    # Regular tag
    MODEL_VARIANT="\`${REPO}:${TAG}\`"
  fi
  
  # Create the table row
  ROW="| $MODEL_VARIANT | $FORMATTED_PARAMS | $QUANTIZATION | - | - | $FORMATTED_SIZE |"
  
  # Add the row to our array
  TABLE_ROWS+=("$ROW")
done

# Find the "Available model variants" section in the readme file
TABLE_SECTION_LINE=$(grep -n "^## Available model variants" "$README_FILE" | cut -d: -f1)
if [ -z "$TABLE_SECTION_LINE" ]; then
  echo "Error: Could not find the 'Available model variants' section in $README_FILE."
  exit 1
fi

# Create a temporary file for the updated content
TMP_FILE=$(mktemp)

# First part: Content before the table
sed -n "1,${TABLE_SECTION_LINE}p" "$README_FILE" > "$TMP_FILE"
echo "" >> "$TMP_FILE"  # Add a newline after the section header

# Add the table header and separator
echo "| Model Variant | Parameters | Quantization | Context window | VRAM | Size |" >> "$TMP_FILE"
echo "|---------------|------------|--------------|----------------|------|-------|" >> "$TMP_FILE"

# Add all the rows
for ROW in "${TABLE_ROWS[@]}"; do
  echo "$ROW" >> "$TMP_FILE"
done

# Add the footnote for VRAM estimation
echo "" >> "$TMP_FILE"
echo "¹: VRAM estimation." >> "$TMP_FILE"

# Add the latest tag mapping note if we found a match
if [ -n "$LATEST_TAG" ]; then
  echo "" >> "$TMP_FILE"
  echo "> \`:latest\` → \`${LATEST_TAG}\`" >> "$TMP_FILE"
fi

# Find the next section after "Available model variants"
NEXT_SECTION_LINE=$(tail -n +$((TABLE_SECTION_LINE + 1)) "$README_FILE" | grep -n "^##" | head -1 | cut -d: -f1)
if [ -n "$NEXT_SECTION_LINE" ]; then
  NEXT_SECTION_LINE=$((TABLE_SECTION_LINE + NEXT_SECTION_LINE))
  
  # Add the content after the table
  echo "" >> "$TMP_FILE"  # Add a newline after the table
  sed -n "${NEXT_SECTION_LINE},\$p" "$README_FILE" >> "$TMP_FILE"
fi

# Replace the original file with the updated content
mv "$TMP_FILE" "$README_FILE"

echo "✅ Successfully updated $README_FILE with all variants for $REPO"
