#!/bin/bash
# filepath: /Users/cluster/dev/go/github.com/go-fork/providers/scripts/update_main_changelog.sh

set -e

# This script updates the main CHANGELOG.md with information from all module-specific changelogs
# It creates a comprehensive summary by extracting information from each module's CHANGELOG.md

# Ensure we're at the repo root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$REPO_ROOT"

VERSION=$1
if [[ -z "$VERSION" ]]; then
  echo "Error: Version parameter is required"
  echo "Usage: $0 <version> (e.g., v0.0.3)"
  exit 1
fi

# Get current date
CURRENT_DATE=$(date +"%Y-%m-%d")

# Get all modules
MODULES=$(find . -name "go.mod" -not -path "*/vendor/*" -not -path "*/.git/*" | xargs dirname | sed 's/^\.\///' | sort)

# Create a temporary file for the new CHANGELOG entry
TEMP_ENTRY=$(mktemp)

# Write header
echo "## $VERSION - $CURRENT_DATE" > "$TEMP_ENTRY"
echo "" >> "$TEMP_ENTRY"
echo "### Module Updates" >> "$TEMP_ENTRY"
echo "" >> "$TEMP_ENTRY"

# Process each module's CHANGELOG
for module in $MODULES; do
  if [[ -f "$module/CHANGELOG.md" ]]; then
    echo "Extracting changes from $module/CHANGELOG.md"
    
    # Check if this module has the current version
    if grep -q "## \[$VERSION\]\\|\## $VERSION" "$module/CHANGELOG.md"; then
      echo "- **$module**: Updated to $VERSION" >> "$TEMP_ENTRY"
      
      # Extract version details
      DETAILS=$(mktemp)
      
      # Extract all sections between version header and next version header
      awk -v ver="$VERSION" '
        BEGIN { p=0 }
        $0 ~ "## \\[?"ver"\\]?|## "ver { p=1; next }
        $0 ~ /^## / { p=0 }
        p == 1 && $0 ~ /^### / { section=$0; gsub(/^### /, "", section); current_section=section; next }
        p == 1 && $0 ~ /^- / && length(current_section) > 0 { 
          items[current_section] = items[current_section] ? items[current_section] "\n    " $0 : "    " $0
        }
        END {
          for (section in items) {
            print "    * " section ":";
            print items[section];
          }
        }
      ' "$module/CHANGELOG.md" > "$DETAILS"
      
      if [[ -s "$DETAILS" ]]; then
        cat "$DETAILS" >> "$TEMP_ENTRY"
        echo "" >> "$TEMP_ENTRY"
      fi
      rm "$DETAILS"
    else
      echo "- **$module**: No changes in this release" >> "$TEMP_ENTRY"
    fi
  fi
done

echo "" >> "$TEMP_ENTRY"
echo "### General Changes" >> "$TEMP_ENTRY"
echo "" >> "$TEMP_ENTRY"
echo "- Updated dependencies to latest versions" >> "$TEMP_ENTRY"
echo "- Improved documentation" >> "$TEMP_ENTRY"
echo "- Bug fixes and performance improvements" >> "$TEMP_ENTRY"
echo "" >> "$TEMP_ENTRY"
echo "See individual module CHANGELOGs for detailed information:" >> "$TEMP_ENTRY"
echo "" >> "$TEMP_ENTRY"

# Add links to module changelogs
for module in $MODULES; do
  if [[ -f "$module/CHANGELOG.md" ]]; then
    echo "- [$module](./$module/CHANGELOG.md)" >> "$TEMP_ENTRY"
  fi
done

echo "" >> "$TEMP_ENTRY"

# Update main CHANGELOG.md
if [[ -f "CHANGELOG.md" ]]; then
  # Create a temporary file with the new content
  TEMP_CHANGELOG=$(mktemp)
  
  # Write the header (first 3 lines)
  head -n 3 "CHANGELOG.md" > "$TEMP_CHANGELOG"
  
  # Add the new entry
  echo "" >> "$TEMP_CHANGELOG"
  cat "$TEMP_ENTRY" >> "$TEMP_CHANGELOG"
  
  # Add the rest of the original content (skip first 3 lines)
  tail -n +4 "CHANGELOG.md" >> "$TEMP_CHANGELOG"
  
  # Replace the original file
  mv "$TEMP_CHANGELOG" "CHANGELOG.md"
  echo "Main CHANGELOG.md updated with comprehensive information"
else
  # Create new CHANGELOG.md
  echo "# Changelog" > "CHANGELOG.md"
  echo "" >> "CHANGELOG.md"
  echo "All notable changes to this project will be documented in this file." >> "CHANGELOG.md"
  echo "" >> "CHANGELOG.md"
  cat "$TEMP_ENTRY" >> "CHANGELOG.md"
  echo "Main CHANGELOG.md created with initial version information"
fi

# Clean up
rm "$TEMP_ENTRY"

echo "Process completed successfully!"
