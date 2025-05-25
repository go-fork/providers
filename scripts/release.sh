#!/bin/bash
# filepath: /Users/cluster/dev/go/github.com/go-fork/providers/scripts/release.sh

set -e

# Source the changelog functions
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
source "$SCRIPT_DIR/changelog_functions.sh"

# Display help message
function show_help {
  echo "Usage: ./release.sh [options]"
  echo ""
  echo "Options:"
  echo "  -m, --module MODULE   Specify the module to release (e.g., cache, log)"
  echo "  -v, --version VERSION Specify the version to release (e.g., v0.1.0)"
  echo "  -a, --all VERSION     Release all modules with the same version"
  echo "  -r, --repo VERSION    Release the entire repository with given version"
  echo "  -c, --create-release  Create GitHub release in addition to tags"
  # Add a new parameter --push-only to allow pushing tags without creating them
  echo "  -f, --force           Skip checking for uncommitted changes"
  echo "  -o, --overwrite       Overwrite existing tags if they already exist" 
  echo "  -p, --push-only       Only push existing tags, don't create new ones"
  echo "  -g, --generate-changelog Generate changelog from commit messages"
  echo "  -h, --help            Display this help message"
  echo ""
  echo "Examples:"
  echo "  ./release.sh --module cache --version v0.1.0  # Release cache module v0.1.0"
  echo "  ./release.sh --all v0.1.0                     # Release all modules with v0.1.0"
  echo "  ./release.sh --repo v0.1.0                    # Release the entire repository as v0.1.0"
  echo "  ./release.sh --module cache --version v0.1.0 --create-release  # Also create GitHub release"
  echo "  ./release.sh --repo v0.1.0 --force            # Skip checking for uncommitted changes"
  echo "  ./release.sh --module config --version v0.1.0 --overwrite  # Overwrite existing tag"
  echo "  ./release.sh --module config --version v0.1.0 --push-only  # Push existing tag only"
  echo "  ./release.sh --module config --version v0.1.0 --generate-changelog  # Auto-generate changelog"
}

# Ensure we're at the repo root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$REPO_ROOT"

# Process command line arguments
while (( "$#" )); do
  case "$1" in
    -m|--module)
      MODULE="$2"
      shift 2
      ;;
    -v|--version)
      VERSION="$2"
      shift 2
      ;;
    -a|--all)
      RELEASE_ALL=true
      VERSION="$2"
      shift 2
      ;;
    -r|--repo)
      RELEASE_REPO=true
      VERSION="$2"
      shift 2
      ;;
    -c|--create-release)
      CREATE_RELEASE=true
      shift
      ;;
    -f|--force)
      FORCE_RELEASE=true
      shift
      ;;
    -o|--overwrite)
      OVERWRITE_TAG=true
      shift
      ;;
    -p|--push-only)
      PUSH_ONLY=true
      shift
      ;;
    -g|--generate-changelog)
      GENERATE_CHANGELOG=true
      shift
      ;;
    -h|--help)
      show_help
      exit 0
      ;;
    --) # end argument parsing
      shift
      break
      ;;
    -*|--*=) # unsupported flags
      echo "Error: Unsupported flag $1" >&2
      show_help
      exit 1
      ;;
    *)
      echo "Error: Unknown parameter $1" >&2
      show_help
      exit 1
      ;;
  esac
done

# Validate version format
if [[ -n "$VERSION" && ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Error: Version must be in format vX.Y.Z (e.g., v0.1.0)" >&2
  exit 1
fi

# Validate options
if [[ "$RELEASE_ALL" == true && "$RELEASE_REPO" == true ]]; then
  echo "Error: Cannot use --all and --repo together" >&2
  exit 1
fi

if [[ "$RELEASE_ALL" == true && -n "$MODULE" ]]; then
  echo "Error: Cannot use --all and --module together" >&2
  exit 1
fi

if [[ "$RELEASE_REPO" == true && -n "$MODULE" ]]; then
  echo "Error: Cannot use --repo and --module together" >&2
  exit 1
fi

if [[ -z "$RELEASE_ALL" && -z "$RELEASE_REPO" && -z "$MODULE" ]]; then
  echo "Error: Must specify --module, --all, or --repo" >&2
  show_help
  exit 1
fi

if [[ -z "$VERSION" ]]; then
  echo "Error: Version is required" >&2
  show_help
  exit 1
fi

# Get all modules
MODULES=$(find . -name "go.mod" -not -path "*/vendor/*" -not -path "*/.git/*" | xargs dirname | sed 's/^\.\///' | sort)

# Function to create GitHub release for a module
function create_github_release_for_module {
  local module=$1
  local version=$2
  local tag_name=$3  # This should now be in the format "module/vX.Y.Z"
  
  # Generate release notes file if CHANGELOG exists
  local release_notes=""
  if [[ -n "$RELEASE_NOTES_FILE" && -f "$RELEASE_NOTES_FILE" ]]; then
    # Use already generated changelog
    release_notes="$RELEASE_NOTES_FILE"
  elif [[ -f "$module/CHANGELOG.md" ]]; then
    release_notes=$(mktemp)
    # Extract the latest release notes from CHANGELOG
    awk -v ver="$version" 'BEGIN{p=0} $0 ~ "^## " ver {p=1} $0 ~ "^## " && $0 !~ ver && p==1 {p=0} p==1 {print}' "$module/CHANGELOG.md" > "$release_notes"
    
    # If we didn't extract anything meaningful (just the header), enhance with features from README
    if [[ $(grep -v "^## $version" "$release_notes" | grep -v "^$" | wc -l) -eq 0 ]]; then
      echo "Enhancing release notes with information from the README.md"
      if [[ -f "$module/README.md" ]]; then
        echo "## Key Features" >> "$release_notes"
        echo "" >> "$release_notes"
        # Extract features from README.md - typically in a section called "Features", "Tính năng", etc.
        grep -A 20 "## Tính năng\|## Features" "$module/README.md" | grep "^-\|^\*" | head -10 >> "$release_notes"
        echo "" >> "$release_notes"
      fi
    fi
  fi
  
  create_github_release "$tag_name" "Release $module $version" "$release_notes"
  
  # Clean up temporary file
  if [[ -n "$release_notes" && -f "$release_notes" ]]; then
    rm "$release_notes"
  fi
}

# Function to create GitHub release for repository
function create_github_release_for_repo {
  local version=$1
  
  # Generate release notes file if CHANGELOG exists
  local release_notes=""
  if [[ -n "$RELEASE_NOTES_FILE" && -f "$RELEASE_NOTES_FILE" ]]; then
    # Use already generated changelog
    release_notes="$RELEASE_NOTES_FILE"
  elif [[ -f "CHANGELOG.md" ]]; then
    release_notes=$(mktemp)
    # Extract the latest release notes from CHANGELOG
    awk -v ver="$version" 'BEGIN{p=0} $0 ~ "^## " ver {p=1} $0 ~ "^## " && $0 !~ ver && p==1 {p=0} p==1 {print}' "CHANGELOG.md" > "$release_notes"
  fi
  
  create_github_release "$version" "Release $version" "$release_notes"
  
  # Clean up temporary file
  if [[ -n "$release_notes" && -f "$release_notes" ]]; then
    rm "$release_notes"
  fi
}

# Function to create a GitHub release
function create_github_release {
  local tag=$1
  local title=$2
  local release_notes=$3
  
  echo "Creating GitHub release for tag: $tag"
  
  # Check if gh CLI is installed
  if ! command -v gh &> /dev/null; then
    echo "Warning: GitHub CLI (gh) is not installed. Cannot create GitHub release."
    echo "To install GitHub CLI, visit: https://cli.github.com/"
    echo "Skipping GitHub release creation..."
    return 1
  fi
  
  # Check if authenticated with GitHub
  if ! gh auth status &> /dev/null; then
    echo "Warning: Not authenticated with GitHub CLI. Cannot create GitHub release."
    echo "Please run 'gh auth login' to authenticate."
    echo "Skipping GitHub release creation..."
    return 1
  fi
  
  # Create GitHub release
  if [ -n "$release_notes" ] && [ -f "$release_notes" ]; then
    gh release create "$tag" --title "$title" --notes-file "$release_notes"
  else
    gh release create "$tag" --title "$title" --notes "Release $tag"
  fi
  
  if [ $? -eq 0 ]; then
    echo "GitHub release created successfully"
  else
    echo "Failed to create GitHub release"
    return 1
  fi
  
  return 0
}

# Function to release a single module
function release_module {
  local module=$1
  local version=$2
  
  local tag_name="${module}/${version}"
  
  # If push-only, just push the tag and exit
  if [[ "$PUSH_ONLY" == true ]]; then
    if git tag | grep -q "^$tag_name$"; then
      echo "Pushing existing tag $tag_name to remote"
      git push origin "$tag_name"
      
      # Create GitHub release if requested
      if [[ "$CREATE_RELEASE" == true ]]; then
        create_github_release_for_module "$module" "$version" "$tag_name"
      fi
      
      return 0
    else
      echo "Error: Tag $tag_name does not exist locally. Cannot push." >&2
      exit 1
    fi
  fi
  
  echo "Preparing release for module: $module, version: $version"
  
  # Verify module exists
  if [[ ! -d "$module" || ! -f "$module/go.mod" ]]; then
    echo "Error: Module directory $module does not exist or doesn't contain go.mod" >&2
    exit 1
  fi
  
  # Check for uncommitted changes
  if [[ "$FORCE_RELEASE" != true && -n "$(git status --porcelain "$module")" ]]; then
    echo "Error: Module $module has uncommitted changes" >&2
    echo "Please commit or stash changes before releasing, or use --force to skip this check" >&2
    exit 1
  fi
  
  # Update CHANGELOG.md if it exists or create it if it doesn't
  if [[ -f "$module/CHANGELOG.md" || "$GENERATE_CHANGELOG" == true || ! -f "$module/CHANGELOG.md" ]]; then
    if [[ ! -f "$module/CHANGELOG.md" ]]; then
      echo "Creating $module/CHANGELOG.md"
      # Creating basic CHANGELOG structure
      echo "# Changelog" > "$module/CHANGELOG.md"
      echo "" >> "$module/CHANGELOG.md"
      echo "All notable changes to this module will be documented in this file." >> "$module/CHANGELOG.md"
      echo "" >> "$module/CHANGELOG.md"
    else
      echo "Updating $module/CHANGELOG.md"
    fi
    
    if [[ "$GENERATE_CHANGELOG" == true ]]; then
      # Generate changelog from commit messages
      local generated_changelog=$(generate_changelog "$module" "$tag_name")
      
      # Update changelog file with generated content
      update_changelog_file "$module" "$version" "$generated_changelog"
      
      # Save generated changelog for release description
      RELEASE_NOTES_FILE="$generated_changelog"
    else
      # Traditional changelog update
      # Get current date
      CURRENT_DATE=$(date +"%Y-%m-%d")
      
      # Create changelog entry
      if [[ -f "$module/README.md" ]]; then
        # Extract feature highlights from the README.md
        echo "Extracting features from README.md for the changelog"
        FEATURE_HIGHLIGHTS=$(mktemp)
        grep -A 20 "## Tính năng\|## Features" "$module/README.md" | grep "^-\|^\*" | head -10 > "$FEATURE_HIGHLIGHTS"
        
        if [[ -s "$FEATURE_HIGHLIGHTS" ]]; then
          # Create proper newlines using actual newlines instead of \n
          {
            echo "## $version - $CURRENT_DATE"
            echo ""
            echo "### Added"
            echo ""
            cat "$FEATURE_HIGHLIGHTS"
            echo ""
            echo ""
          } > "$FEATURE_HIGHLIGHTS.formatted"
          
          # Read the formatted content into the CHANGELOG_ENTRY variable
          CHANGELOG_ENTRY=$(cat "$FEATURE_HIGHLIGHTS.formatted")
          
          rm "$FEATURE_HIGHLIGHTS" "$FEATURE_HIGHLIGHTS.formatted"
        else
          CHANGELOG_ENTRY=$(echo -e "## $version - $CURRENT_DATE\n\n* See GitHub release notes\n\n")
        fi
      else
        CHANGELOG_ENTRY=$(echo -e "## $version - $CURRENT_DATE\n\n* See GitHub release notes\n\n")
      fi
      
      # Add entry to the top of the changelog (after the header)
      if [[ -f "$module/CHANGELOG.md" ]]; then
        # Create a temporary file with the new content
        TEMP_CHANGELOG=$(mktemp)
        
        # Write the header (first 3 lines)
        head -n 3 "$module/CHANGELOG.md" > "$TEMP_CHANGELOG"
        
        # Add the new entry
        echo "" >> "$TEMP_CHANGELOG"
        echo "$CHANGELOG_ENTRY" >> "$TEMP_CHANGELOG"
        
        # Add the rest of the original content (skip first 3 lines)
        tail -n +4 "$module/CHANGELOG.md" >> "$TEMP_CHANGELOG"
        
        # Replace the original file
        mv "$TEMP_CHANGELOG" "$module/CHANGELOG.md"
      else
        # Create new file with header and entry
        echo "# Changelog" > "$module/CHANGELOG.md"
        echo "" >> "$module/CHANGELOG.md"
        echo "All notable changes to the $module module will be documented in this file." >> "$module/CHANGELOG.md"
        echo "" >> "$module/CHANGELOG.md"
        echo "$CHANGELOG_ENTRY" >> "$module/CHANGELOG.md"
      fi
      
      # Commit changelog update
      git add "$module/CHANGELOG.md"
      git commit -m "docs($module): update CHANGELOG for $version"
    fi
  fi
  
  # Create tag for module (using the standard Go modules convention with slash separator)
  local tag_name="${module}/${version}"
  
  # Check if tag exists and handle it
  if git tag | grep -q "^$tag_name$"; then
    if [[ "$OVERWRITE_TAG" == true ]]; then
      echo "Tag $tag_name already exists, removing it (--overwrite specified)"
      git tag -d "$tag_name"
      
      # Check if tag exists in remote
      if git ls-remote --tags origin | grep -q "refs/tags/$tag_name$"; then
        echo "Removing tag from remote as well"
        git push origin ":refs/tags/$tag_name" || {
          echo "Warning: Failed to delete remote tag. You may need to delete it manually."
        }
      fi
    else
      echo "Error: Tag $tag_name already exists. Use --overwrite to replace it." >&2
      exit 1
    fi
  fi
  
  echo "Creating tag: $tag_name"
  git tag -a "$tag_name" -m "Release $module $version"
  
  echo "Module $module $version prepared for release"
  
  # Automatically push tag to remote if creating a GitHub release
  if [[ "$CREATE_RELEASE" == true ]]; then
    echo "Pushing tag to remote repository (required for GitHub release)"
    git push origin "$tag_name"
  else
    echo "To push this release, run: git push origin $tag_name"
  fi
  
  # Create GitHub release if requested
  if [[ "$CREATE_RELEASE" == true ]]; then
    create_github_release_for_module "$module" "$version" "$tag_name"
  fi
}

# Function to release the entire repository
function release_repo {
  local version=$1
  
  # If push-only, just push the tag and exit
  if [[ "$PUSH_ONLY" == true ]]; then
    if git tag | grep -q "^$version$"; then
      echo "Pushing existing tag $version to remote"
      git push origin "$version"
      
      # Create GitHub release if requested
      if [[ "$CREATE_RELEASE" == true ]]; then
        create_github_release_for_repo "$version"
      fi
      
      return 0
    else
      echo "Error: Tag $version does not exist locally. Cannot push." >&2
      exit 1
    fi
  fi
  
  echo "Preparing release for entire repository, version: $version"
  
  # Check for uncommitted changes
  if [[ "$FORCE_RELEASE" != true && -n "$(git status --porcelain)" ]]; then
    echo "Error: Repository has uncommitted changes" >&2
    echo "Please commit or stash changes before releasing, or use --force to skip this check" >&2
    exit 1
  fi
  
  # Update root CHANGELOG.md if it exists
  if [[ -f "CHANGELOG.md" || "$GENERATE_CHANGELOG" == true ]]; then
    echo "Updating CHANGELOG.md"
    
    if [[ "$GENERATE_CHANGELOG" == true ]]; then
      # Generate changelog from commit messages
      local generated_changelog=$(generate_changelog "." "$version")
      
      # Update changelog file with generated content
      update_changelog_file "." "$version" "$generated_changelog"
      
      # Save generated changelog for release description
      RELEASE_NOTES_FILE="$generated_changelog"
    else
      # Traditional changelog update
      # Get current date
      CURRENT_DATE=$(date +"%Y-%m-%d")
      
      # Create changelog entry
      CHANGELOG_ENTRY="## $version - $CURRENT_DATE\n\n* See GitHub release notes\n\n"
      
      # Add entry to the top of the changelog (after the header)
      sed -i '' -e "4i\\
$CHANGELOG_ENTRY
" "CHANGELOG.md"
      
      # Commit changelog update
      git add "CHANGELOG.md"
      git commit -m "docs: update CHANGELOG for $version"
    fi
  fi
  
  # Create tag for repository
  echo "Creating tag: $version"
  
  # Check if tag exists and handle it
  if git tag | grep -q "^$version$"; then
    if [[ "$OVERWRITE_TAG" == true ]]; then
      echo "Tag $version already exists, removing it (--overwrite specified)"
      git tag -d "$version"
      
      # Check if tag exists in remote
      if git ls-remote --tags origin | grep -q "refs/tags/$version$"; then
        echo "Removing tag from remote as well"
        git push origin ":refs/tags/$version" || {
          echo "Warning: Failed to delete remote tag. You may need to delete it manually."
        }
      fi
    else
      echo "Error: Tag $version already exists. Use --overwrite to replace it." >&2
      exit 1
    fi
  fi
  
  git tag -a "$version" -m "Release $version"
  
  echo "Repository $version prepared for release"
  
  # Automatically push tag to remote if creating a GitHub release
  if [[ "$CREATE_RELEASE" == true ]]; then
    echo "Pushing tag to remote repository (required for GitHub release)"
    git push origin "$version"
  else
    echo "To push this release, run: git push origin $version"
  fi
  
  # Create GitHub release if requested
  if [[ "$CREATE_RELEASE" == true ]]; then
    create_github_release_for_repo "$version"
  fi
}

# Check module compatibility if releasing multiple modules
function check_module_compatibility {
  if [[ -f "$SCRIPT_DIR/check_compatibility.sh" ]]; then
    echo "Checking module compatibility..."
    bash "$SCRIPT_DIR/check_compatibility.sh"
    
    # Ask user if they want to continue despite any issues
    if [[ $? -ne 0 ]]; then
      read -p "Do you want to continue with the release anyway? (y/n): " confirm
      if [[ ! "$confirm" =~ ^[yY]$ ]]; then
        echo "Release cancelled."
        exit 1
      fi
    fi
  else
    echo "Warning: check_compatibility.sh not found, skipping compatibility check"
  fi
}

# Execute the requested action
if [[ "$RELEASE_ALL" == true ]]; then
  # Check module compatibility
  check_module_compatibility
  echo "Releasing all modules with version $VERSION"
  
  # Update main CHANGELOG.md with comprehensive information from all module changelogs
  if [[ -f "$SCRIPT_DIR/update_main_changelog.sh" ]]; then
    echo "Updating main CHANGELOG.md with comprehensive information from all module changelogs"
    bash "$SCRIPT_DIR/update_main_changelog.sh" "$VERSION"
    
    # Commit changelog update
    git add "CHANGELOG.md"
    git commit -m "docs: update main CHANGELOG for $VERSION with comprehensive module information"
  else
    echo "Warning: update_main_changelog.sh not found, using simple changelog update"
    
    if [[ -f "CHANGELOG.md" ]]; then
      echo "Updating main CHANGELOG.md with links to module changelogs"
      CURRENT_DATE=$(date +"%Y-%m-%d")
      
      # Create changelog entry with links to module changelogs
      MODULE_LINKS=""
      for module in $MODULES; do
        if [[ -f "$module/CHANGELOG.md" ]]; then
          MODULE_LINKS+="* [$module](./$module/CHANGELOG.md)\n"
        fi
      done
      
      CHANGELOG_ENTRY="## $VERSION - $CURRENT_DATE\n\n### Module Updates\nSee individual module changelogs for details:\n\n$MODULE_LINKS\n"
      
      # Add entry to the top of the changelog (after the header)
      sed -i '' -e "4i\\
$CHANGELOG_ENTRY
" "CHANGELOG.md"
      
      # Commit changelog update
      git add "CHANGELOG.md"
      git commit -m "docs: update main CHANGELOG for $VERSION with module links"
    fi
  fi
  
  # Process each module
  for module in $MODULES; do
    release_module "$module" "$VERSION"
  done
  echo "All modules prepared for release with version $VERSION"
  echo "To push all releases, run: git push origin --tags"
  
  # Create GitHub releases if requested
  if [[ "$CREATE_RELEASE" == true ]]; then
    echo "Creating GitHub releases for all modules..."
    for module in $MODULES; do
      local tag_name="${module}/${VERSION}"
      create_github_release_for_module "$module" "$VERSION" "$tag_name"
    done
  fi
elif [[ "$RELEASE_REPO" == true ]]; then
  release_repo "$VERSION"
else
  release_module "$MODULE" "$VERSION"
fi
