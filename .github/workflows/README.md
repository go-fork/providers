# GitHub Actions Workflows

This directory contains all GitHub Actions workflows for the Go-Fork Providers repository.

## Core Workflows

### [CI](ci.yml)
- Runs on: push to main, pull requests to main
- Purpose: Runs tests and linting
- Features:
  - Tests across multiple Go versions (1.21, 1.22, 1.23)
  - Tests across multiple operating systems (Linux, macOS)
  - Reports code coverage to Codecov

### [Release](release.yml)
- Runs on: new tags (v*)
- Purpose: Creates GitHub releases
- Features:
  - Uses GoReleaser for release creation
  - Automatically generates release notes

### [Update README](update-readme.yml)
- Runs on: push to main affecting code files
- Purpose: Updates README files with dependency graphs
- Features:
  - Generates SVG dependency graphs for each module
  - Updates README files with latest graphs

## Security & Dependencies

### [CodeQL Analysis](codeql-analysis.yml)
- Runs on: push to main, pull requests to main, weekly schedule
- Purpose: Identifies security vulnerabilities
- Features:
  - Static code analysis for security issues
  - Reports findings as GitHub Security alerts

### [Dependency Review](dependency-review.yml)
- Runs on: pull requests to main
- Purpose: Reviews dependency changes
- Features:
  - Detects vulnerable dependencies
  - Checks license compliance

### [Auto-merge Dependabot PRs](auto-merge-dependabot.yml)
- Runs on: pull requests from Dependabot
- Purpose: Automates dependency updates
- Features:
  - Automatically approves and merges safe dependency updates
  - Only merges minor and patch updates

### [Go Module Validation](go-mod-validation.yml)
- Runs on: changes to go.mod files
- Purpose: Ensures go.mod files are properly maintained
- Features:
  - Verifies go.mod files are tidy
  - Checks for consistent dependency versions across modules

## Quality & Testing

### [Code Quality Analysis](code-quality.yml)
- Runs on: push to main, pull requests affecting code
- Purpose: Comprehensive code quality checks
- Features:
  - Advanced linting with multiple tools
  - Complexity analysis
  - Code coverage reporting

### [Documentation Check](doc-check.yml)
- Runs on: changes to Go files or doc.go files
- Purpose: Enforces documentation standards
- Features:
  - Verifies each package has doc.go
  - Checks spelling in documentation
  - Ensures godoc quality

### [Backwards Compatibility Check](compatibility-check.yml)
- Runs on: pull requests affecting code
- Purpose: Prevents breaking changes
- Features:
  - Uses apidiff to detect API changes
  - Warns about potentially breaking changes

### [Cross-Version Compatibility](cross-version-check.yml)
- Runs on: code changes, weekly schedule
- Purpose: Ensures compatibility across Go versions
- Features:
  - Tests on Go 1.20 through 1.23 and tip
  - Tests on multiple operating systems

### [Integration Tests](integration-tests.yml)
- Runs on: code changes, weekly schedule
- Purpose: Tests integration between packages
- Features:
  - Sets up test services (Redis, MongoDB)
  - Tests cross-package functionality

## Examples & Documentation

### [Generate Documentation](docs.yml)
- Runs on: changes to Go files or doc.go files
- Purpose: Builds comprehensive documentation
- Features:
  - Generates markdown docs from godoc
  - Creates documentation site structure

### [Validate Examples](examples.yml)
- Runs on: changes to examples
- Purpose: Ensures example quality
- Features:
  - Checks example structure and required files
  - Verifies examples compile

### [Test Examples](examples-test.yml)
- Runs on: changes to examples or code
- Purpose: Tests example functionality
- Features:
  - Compiles all examples
  - Runs examples where possible
  - Generates example documentation

## Metrics & Reporting

### [Benchmark](benchmark.yml)
- Runs on: code changes, weekly schedule
- Purpose: Tracks performance
- Features:
  - Runs Go benchmarks
  - Records and compares results over time
  - Alerts on significant regressions

### [Code Metrics Report](metrics-report.yml)
- Runs on: push to main, monthly schedule
- Purpose: Provides code quality insights
- Features:
  - Generates detailed code metrics
  - Creates visualizations
  - Submits reports as PRs

## Best Practices for Workflow Development

When modifying or adding workflows:

1. **Use Consistent Naming**: Follow the existing naming pattern for workflow files.

2. **Provide Detailed Comments**: Explain what each section of the workflow does.

3. **Use Matrix Testing**: When applicable, test across multiple environments.

4. **Share Common Steps**: Use composite actions for steps used in multiple workflows.

5. **Set Resource Limits**: Be mindful of resource usage, especially in frequently run workflows.

6. **Add Scheduled Runs**: Consider adding scheduled runs for critical checks.

7. **Optimize Caching**: Use appropriate caching to speed up workflows.

8. **Test Locally First**: Test complex workflows using [act](https://github.com/nektos/act) before pushing.

9. **Limit Triggers**: Only run workflows when needed by using specific path filters.

10. **Monitor Usage**: Keep an eye on workflow usage to stay within GitHub Actions limits.
