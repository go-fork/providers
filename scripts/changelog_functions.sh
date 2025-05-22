# Function to find previous tag for a module
function find_previous_tag {
  local module=$1
  local current_tag=$2
  
  # For module tags
  if [[ -n "$module" && "$module" != "." ]]; then
    # Find all tags for this module, sort by version, and get the previous one
    git tag -l "${module}-v*" | sort -V | grep -B 1 "${current_tag}$" | head -n 1
  else
    # For repository tags, find all vX.Y.Z tags
    git tag -l "v*" | grep -v "-" | sort -V | grep -B 1 "${current_tag}$" | head -n 1
  fi
}

# Function to generate changelog from commits
function generate_changelog {
  local module=$1
  local current_tag=$2
  local previous_tag=$(find_previous_tag "$module" "$current_tag")
  local changelog_file=$(mktemp)
  
  echo "Generating changelog for $current_tag"
  if [[ -n "$previous_tag" ]]; then
    echo "Found previous tag: $previous_tag"
    
    # Generate changelog header
    {
      echo "## Changes since $previous_tag"
      echo ""
      
      # Add compare URL
      echo "**Full Changelog**: https://github.com/go-fork/providers/compare/${previous_tag}...${current_tag}"
      echo ""
      
      # Add commits
      if [[ -n "$module" && "$module" != "." ]]; then
        # For module, only include commits that touch the module directory
        echo "### Commits"
        echo ""
        git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" -- "$module/" | grep -v "Merge" | sort
      else
        # For repository, include all commits
        echo "### Commits"
        echo ""
        git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" | grep -v "Merge" | sort
      fi
      
      # Group changes by type using conventional commits
      echo ""
      echo "### Features"
      echo ""
      git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" | grep -E "^feat(\([^)]+\))?:" | sed 's/feat(\([^)]*\)):/\1:/g' | sort
      
      echo ""
      echo "### Bug Fixes"
      echo ""
      git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" | grep -E "^fix(\([^)]+\))?:" | sed 's/fix(\([^)]*\)):/\1:/g' | sort
      
      echo ""
      echo "### Documentation"
      echo ""
      git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" | grep -E "^docs(\([^)]+\))?:" | sed 's/docs(\([^)]*\)):/\1:/g' | sort
      
      echo ""
      echo "### Other Changes"
      echo ""
      git log --pretty=format:"* %s (%h)" "${previous_tag}..HEAD" | grep -v -E "^(feat|fix|docs|test|build|ci|chore|style|refactor|perf|test)(\([^)]+\))?:" | grep -v "Merge" | sort
    } > "$changelog_file"
  else
    echo "No previous tag found, generating changelog from all commits"
    
    # Generate changelog header
    {
      echo "## Initial Release"
      echo ""
      
      # Add commits
      if [[ -n "$module" && "$module" != "." ]]; then
        echo "### Commits"
        echo ""
        git log --pretty=format:"* %s (%h)" -- "$module/" | grep -v "Merge" | sort
      else
        echo "### Commits"
        echo ""
        git log --pretty=format:"* %s (%h)" | grep -v "Merge" | sort
      fi
    } > "$changelog_file"
  fi
  
  echo "$changelog_file"
}

# Function to update CHANGELOG.md with generated content
function update_changelog_file {
  local module_path=$1
  local version=$2
  local generated_changelog=$3
  
  # Determine path to CHANGELOG.md
  local changelog_path="${module_path}/CHANGELOG.md"
  if [[ "$module_path" == "." ]]; then
    changelog_path="CHANGELOG.md"
  fi
  
  # Create CHANGELOG.md if it doesn't exist
  if [[ ! -f "$changelog_path" ]]; then
    echo "Creating new $changelog_path"
    echo "# Changelog" > "$changelog_path"
    echo "" >> "$changelog_path"
    echo "All notable changes to this project will be documented in this file." >> "$changelog_path"
    echo "" >> "$changelog_path"
  fi
  
  # Get current date
  CURRENT_DATE=$(date +"%Y-%m-%d")
  
  # Create temporary file for new changelog
  local temp_file=$(mktemp)
  
  # Write header and first three lines of existing changelog
  head -n 3 "$changelog_path" > "$temp_file"
  
  # Add new release header
  echo "" >> "$temp_file"
  echo "## $version - $CURRENT_DATE" >> "$temp_file"
  echo "" >> "$temp_file"
  
  # Add generated changelog content
  cat "$generated_changelog" >> "$temp_file"
  echo "" >> "$temp_file"
  
  # Add existing changelog content (skip first three lines)
  if [[ $(wc -l < "$changelog_path") -gt 3 ]]; then
    tail -n +4 "$changelog_path" >> "$temp_file"
  fi
  
  # Replace original file with new content
  mv "$temp_file" "$changelog_path"
  
  # Commit changelog update
  git add "$changelog_path"
  if [[ "$module_path" == "." ]]; then
    git commit -m "docs: update CHANGELOG for $version"
  else
    git commit -m "docs($module_path): update CHANGELOG for $version"
  fi
  
  echo "Updated $changelog_path with generated changelog"
}
