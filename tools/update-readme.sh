#!/usr/bin/env bash
set -euo pipefail

# Parse arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    *)
      if [ -z "${MODEL_REF:-}" ]; then
        MODEL_REF="$1"
      elif [ -z "${CONTEXT_WINDOW:-}" ]; then
        CONTEXT_WINDOW="$1"
      elif [ -z "${VRAM:-}" ]; then
        VRAM="$1"
      else
        echo "❌ Unexpected argument: $1"
        echo "Usage: $0 <model-reference> [context-window] [vram]"
        exit 1
      fi
      shift
      ;;
  esac
done

# Check if the required arguments are provided
if [ -z "${MODEL_REF:-}" ]; then
  echo "Usage: $0 <model-reference> [context-window] [vram]"
  echo "Example: $0 ai/smollm2:360M-Q4_0 8K 220"
  exit 1
fi

# Set default values for optional parameters
CONTEXT_WINDOW="${CONTEXT_WINDOW:-}"
VRAM="${VRAM:-}"

# Validate model reference format
if [[ ! "$MODEL_REF" == *":"* ]]; then
  echo "❌ Error: Model reference must include a tag (e.g., ai/modelname:tag)"
  exit 1
fi

if [[ ! "$MODEL_REF" == *"/"* ]]; then
  echo "❌ Error: Model reference must include a namespace (e.g., ai/modelname:tag)"
  exit 1
fi

# Extract repository part (before the colon)
REPO_PART=${MODEL_REF%%:*}

# Extract model name (after the last slash)
MODEL_NAME=${REPO_PART##*/}

# Extract namespace (before the last slash)
NAMESPACE=${REPO_PART%/*}

# Construct readme path
README_FILE="${NAMESPACE}/${MODEL_NAME}.md"

echo "📄 Using readme file: $README_FILE"

# Check if the readme file exists
if [ ! -f "$README_FILE" ]; then
  echo "Error: Readme file '$README_FILE' does not exist."
  exit 1
fi

echo "🔍 Running inspect-model.sh for $MODEL_REF..."
MODEL_INFO=$(./tools/inspect-model.sh "$MODEL_REF")

# Extract information from the output
MODEL_VARIANT=$(echo "$MODEL_INFO" | grep "Image" | sed -E 's/.*: (.+)$/\1/' | tr -d ' ')
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
if [[ "$PARAMETERS" == *"M"* ]]; then
  # Already in M format
  FORMATTED_PARAMS="$PARAMETERS"
elif [[ "$PARAMETERS" == *"B"* ]]; then
  # Already in B format
  FORMATTED_PARAMS="$PARAMETERS"
else
  # Try to convert to a readable format
  FORMATTED_PARAMS="$PARAMETERS"
fi

# Set default values for optional parameters if not provided
if [ -z "$CONTEXT_WINDOW" ]; then
  CONTEXT_WINDOW="N/A"
else
  CONTEXT_WINDOW="${CONTEXT_WINDOW} tokens"
fi

if [ -z "$VRAM" ]; then
  VRAM="N/A"
else
  VRAM="${VRAM} MB¹"
fi

# Create the new table row
NEW_ROW="| \`$MODEL_VARIANT\` | $FORMATTED_PARAMS | $QUANTIZATION | $CONTEXT_WINDOW | $VRAM | $FORMATTED_SIZE |"

echo "📝 Adding the following row to $README_FILE:"
echo "$NEW_ROW"

# Check if the model variant already exists in the file
# Use a more precise pattern to avoid partial matches
if grep -q "\`$MODEL_VARIANT\`" "$README_FILE"; then
  echo "Model variant $MODEL_VARIANT already exists. Updating entry."
  
  # Remove the existing line with this model variant
  TMP_FILE=$(mktemp)
  grep -v "$MODEL_VARIANT" "$README_FILE" > "$TMP_FILE"
  mv "$TMP_FILE" "$README_FILE"
  echo "Removed existing entry for $MODEL_VARIANT."
fi

# Find the "Available model variants" section and the table within it
echo "🔍 Finding the model variants table..."

# Create a temporary file for the updated content
TMP_FILE=$(mktemp)

# Find the line number of the "Available model variants" section
TABLE_SECTION_LINE=$(grep -n "^## Available model variants" "$README_FILE" | cut -d: -f1)

if [ -z "$TABLE_SECTION_LINE" ]; then
  echo "Error: Could not find the 'Available model variants' section in $README_FILE."
  exit 1
fi

echo "📊 Found model variants section at line $TABLE_SECTION_LINE"

# First pass: Find the last line of the table
LINE_NUM=0
IN_TABLE=false
LAST_TABLE_LINE=0

while IFS= read -r line; do
  LINE_NUM=$((LINE_NUM + 1))
  
  # Check if we're in the "Available model variants" section
  if [ $LINE_NUM -ge $TABLE_SECTION_LINE ] && [[ "$line" =~ ^## && ! "$line" =~ ^"## Available model variants" ]]; then
    # We've reached the next section, so we're no longer in the table section
    IN_TABLE=false
  fi
  
  # If we're in the table section and the line starts with "|", update the last table line
  if [ $LINE_NUM -ge $TABLE_SECTION_LINE ] && $IN_TABLE && [[ "$line" =~ \| ]]; then
    LAST_TABLE_LINE=$LINE_NUM
  fi
  
  # If we've found the "Available model variants" section, we're in the table section
  if [ $LINE_NUM -eq $TABLE_SECTION_LINE ]; then
    IN_TABLE=true
  fi
done < "$README_FILE"

echo "📊 Found last table line at line $LAST_TABLE_LINE"

# Second pass: Create the updated file with the new row
LINE_NUM=0

while IFS= read -r line; do
  LINE_NUM=$((LINE_NUM + 1))
  
  # Print the current line to the temporary file
  echo "$line" >> "$TMP_FILE"
  
  # If we've just processed the last line of the table, add the new row
  if [ $LINE_NUM -eq $LAST_TABLE_LINE ]; then
    echo "$NEW_ROW" >> "$TMP_FILE"
    echo "📝 Added new row after line $LAST_TABLE_LINE"
  fi
done < "$README_FILE"

# If we didn't find any table lines, append the row at the end of the file
if [ $LAST_TABLE_LINE -eq 0 ]; then
  echo "⚠️ Could not find the end of the table. Appending the row at the end of the file."
  echo "$NEW_ROW" >> "$TMP_FILE"
fi

# Replace the original file with the updated content
mv "$TMP_FILE" "$README_FILE"

echo "✅ Successfully updated $README_FILE with information for $MODEL_REF."
