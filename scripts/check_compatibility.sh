#!/bin/bash
# filepath: /Users/cluster/dev/go/github.com/go-fork/providers/scripts/check_compatibility.sh

set -e

# This script checks compatibility between modules before release
# It analyzes module dependencies and reports potential issues

# Ensure we're at the repo root
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$REPO_ROOT"

# Get all modules
MODULES=$(find . -name "go.mod" -not -path "*/vendor/*" -not -path "*/.git/*" | xargs dirname | sed 's/^\.\///' | sort)

# Function to extract module dependencies
function get_module_dependencies() {
  local module_path=$1
  if [[ -f "$module_path/go.mod" ]]; then
    # Extract lines that contain github.com/go-fork/providers
    grep -E "github.com/go-fork/providers/.+" "$module_path/go.mod" | grep -v "replace" | awk '{print $1 " " $2}' || true
  fi
}

# Function to check if a module depends on another module
function depends_on() {
  local module=$1
  local dependency=$2
  
  if grep -q "github.com/go-fork/providers/$dependency" "$module/go.mod"; then
    return 0
  else
    return 1
  fi
}

echo "Checking module compatibility..."
echo

# Check for compatibility issues
ISSUES_FOUND=0
REPORT=""

for module in $MODULES; do
  MODULE_DEPS=$(get_module_dependencies "$module")
  if [[ -n "$MODULE_DEPS" ]]; then
    REPORT+="Module $module depends on:\n"
    
    while IFS= read -r dep; do
      if [[ -n "$dep" ]]; then
        dep_module=$(echo "$dep" | awk '{print $1}' | sed 's|github.com/go-fork/providers/||')
        dep_version=$(echo "$dep" | awk '{print $2}')
        
        REPORT+="  - $dep_module ($dep_version)\n"
        
        # Check if dependency has a different version in the workspace
        if [[ -f "$dep_module/go.mod" ]]; then
          current_version=$(grep -E "^module" "$dep_module/go.mod" | awk '{print $2}' | xargs basename)
          if [[ "$current_version" != "$dep_version" ]]; then
            REPORT+="    ⚠️ Warning: $module depends on $dep_module $dep_version, but current workspace has $current_version\n"
            ISSUES_FOUND=1
          fi
        else
          REPORT+="    ⚠️ Warning: Dependency $dep_module not found in workspace\n"
          ISSUES_FOUND=1
        fi
      fi
    done <<< "$MODULE_DEPS"
    
    REPORT+="\n"
  fi
done

echo -e "$REPORT"

if [[ $ISSUES_FOUND -eq 1 ]]; then
  echo "⚠️ Compatibility issues found between modules."
  echo "Please review the warnings above before releasing."
  echo "You may need to update module dependencies or make sure they are compatible."
else
  echo "✅ No compatibility issues found between modules."
fi

# Check for circular dependencies
echo
echo "Checking for circular dependencies..."
CIRCULAR_DEPS_FOUND=0

for module1 in $MODULES; do
  for module2 in $MODULES; do
    if [[ "$module1" != "$module2" ]]; then
      if depends_on "$module1" "$module2" && depends_on "$module2" "$module1"; then
        echo "⚠️ Circular dependency detected: $module1 <-> $module2"
        CIRCULAR_DEPS_FOUND=1
      fi
    fi
  done
done

if [[ $CIRCULAR_DEPS_FOUND -eq 0 ]]; then
  echo "✅ No circular dependencies found."
fi
