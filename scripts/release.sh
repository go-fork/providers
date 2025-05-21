#!/bin/bash
# filepath: /Users/cluster/dev/go/github.com/go-fork/providers/scripts/release.sh

set -e

# Display help message
function show_help {
  echo "Usage: ./release.sh [options]"
  echo ""
  echo "Options:"
  echo "  -m, --module MODULE   Specify the module to release (e.g., cache, log)"
  echo "  -v, --version VERSION Specify the version to release (e.g., v0.1.0)"
  echo "  -a, --all VERSION     Release all modules with the same version"
  echo "  -r, --repo VERSION    Release the entire repository with given version"
  echo "  -h, --help            Display this help message"
  echo ""
  echo "Examples:"
  echo "  ./release.sh --module cache --version v0.1.0  # Release cache module v0.1.0"
  echo "  ./release.sh --all v0.1.0                     # Release all modules with v0.1.0"
  echo "  ./release.sh --repo v0.1.0                    # Release the entire repository as v0.1.0"
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

# Function to release a single module
function release_module {
  local module=$1
  local version=$2
  
  echo "Preparing release for module: $module, version: $version"
  
  # Verify module exists
  if [[ ! -d "$module" || ! -f "$module/go.mod" ]]; then
    echo "Error: Module directory $module does not exist or doesn't contain go.mod" >&2
    exit 1
  fi
  
  # Check for uncommitted changes
  if [[ -n "$(git status --porcelain "$module")" ]]; then
    echo "Error: Module $module has uncommitted changes" >&2
    echo "Please commit or stash changes before releasing" >&2
    exit 1
  fi
  
  # Update CHANGELOG.md if it exists
  if [[ -f "$module/CHANGELOG.md" ]]; then
    echo "Updating $module/CHANGELOG.md"
    
    # Get current date
    CURRENT_DATE=$(date +"%Y-%m-%d")
    
    # Create changelog entry
    CHANGELOG_ENTRY="## $version - $CURRENT_DATE\n\n* See GitHub release notes\n\n"
    
    # Add entry to the top of the changelog (after the header)
    sed -i '' -e "4i\\
$CHANGELOG_ENTRY
" "$module/CHANGELOG.md"
    
    # Commit changelog update
    git add "$module/CHANGELOG.md"
    git commit -m "docs($module): update CHANGELOG for $version"
  fi
  
  # Create tag for module
  echo "Creating tag: $module/$version"
  git tag -a "$module/$version" -m "Release $module $version"
  
  echo "Module $module $version prepared for release"
  echo "To push this release, run: git push origin $module/$version"
}

# Function to release the entire repository
function release_repo {
  local version=$1
  
  echo "Preparing release for entire repository, version: $version"
  
  # Check for uncommitted changes
  if [[ -n "$(git status --porcelain)" ]]; then
    echo "Error: Repository has uncommitted changes" >&2
    echo "Please commit or stash changes before releasing" >&2
    exit 1
  fi
  
  # Update root CHANGELOG.md if it exists
  if [[ -f "CHANGELOG.md" ]]; then
    echo "Updating CHANGELOG.md"
    
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
  
  # Create tag for repository
  echo "Creating tag: $version"
  git tag -a "$version" -m "Release $version"
  
  echo "Repository $version prepared for release"
  echo "To push this release, run: git push origin $version"
}

# Execute the requested action
if [[ "$RELEASE_ALL" == true ]]; then
  echo "Releasing all modules with version $VERSION"
  for module in $MODULES; do
    release_module "$module" "$VERSION"
  done
  echo "All modules prepared for release with version $VERSION"
  echo "To push all releases, run: git push origin --tags"
elif [[ "$RELEASE_REPO" == true ]]; then
  release_repo "$VERSION"
else
  release_module "$MODULE" "$VERSION"
fi
