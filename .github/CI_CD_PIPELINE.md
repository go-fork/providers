# CI/CD Pipeline

This document provides an overview of the CI/CD pipeline configuration for the Go-Fork Providers repository.

## Workflows

The repository uses GitHub Actions to automate various tasks:

### Core Workflows

- **CI (`ci.yml`)**: Runs tests and linting on pull requests and pushes to main branch.
- **Release (`release.yml`)**: Creates a new release when a tag is pushed.
- **Update README (`update-readme.yml`)**: Automatically updates README files with dependency graphs.

### Security & Dependencies

- **CodeQL Analysis (`codeql-analysis.yml`)**: Performs security analysis on the codebase.
- **Dependency Review (`dependency-review.yml`)**: Reviews dependencies in pull requests for vulnerabilities.
- **Auto-merge Dependabot PRs (`auto-merge-dependabot.yml`)**: Automatically approves and merges minor/patch dependency updates.
- **Go Module Validation (`go-mod-validation.yml`)**: Validates go.mod files and checks for consistent dependency versions.

### Quality & Testing

- **Code Quality Analysis (`code-quality.yml`)**: Performs comprehensive linting and code quality checks.
- **Documentation Check (`doc-check.yml`)**: Ensures all packages have proper documentation.
- **Backwards Compatibility Check (`compatibility-check.yml`)**: Verifies API compatibility between releases.
- **Cross-Version Compatibility (`cross-version-check.yml`)**: Tests code against multiple Go versions.
- **Integration Tests (`integration-tests.yml`)**: Runs tests that verify interactions between different packages.

### Examples & Documentation

- **Generate Documentation (`docs.yml`)**: Automatically generates and updates documentation from godoc comments.
- **Validate Examples (`examples.yml`)**: Ensures examples compile and follow project standards.
- **Test Examples (`examples-test.yml`)**: Compiles and tests example code to ensure it works correctly.

### Metrics & Reporting

- **Benchmark (`benchmark.yml`)**: Runs performance benchmarks weekly and on code changes.
- **Code Metrics Report (`metrics-report.yml`)**: Generates detailed code metrics and reports.

## Schedule

Some workflows run on schedule:

- CodeQL Analysis: Every Monday at 3 AM UTC
- Benchmarks: Every Sunday at midnight UTC
- Cross-Version Compatibility: Every Monday at 4 AM UTC
- Integration Tests: Every Monday at 2 AM UTC
- Code Metrics Report: First day of each month at midnight UTC

## Dependencies

The repository uses Dependabot to keep dependencies up-to-date:

- Go modules in root and subdirectories (weekly)
- GitHub Actions (weekly)

## Pull Request Validation

When creating a pull request, the following checks will run:

1. CI tests and linting
2. CodeQL security analysis
3. Dependency review
4. Documentation quality check
5. Example validation (if applicable)
6. Backwards compatibility check
7. Go module validation

## Release Process

To create a new release:

1. Create and push a new tag (e.g., `git tag v1.2.3 && git push origin v1.2.3`)
2. The release workflow will automatically:
   - Create a GitHub release
   - Generate release notes based on commits
   - Publish artifacts

## Local Development

Developers can run the same checks locally before pushing:

```bash
# Run tests
go test -v -race ./...

# Run linting
golangci-lint run

# Validate docs
go install github.com/client9/misspell/cmd/misspell@latest
misspell -error **/README.md **/doc.go

# Compile examples
find examples -name "*.go" | while read -r file; do (cd "$(dirname "$file")" && go build); done

# Check for breaking changes (requires installed tools)
go install golang.org/x/exp/cmd/apidiff@latest
apidiff ./path/to/old ./path/to/new

# Run benchmarks
go test -bench=. -benchmem ./...
```
