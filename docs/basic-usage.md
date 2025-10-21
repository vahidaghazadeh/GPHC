# Basic Usage Guide

This guide covers the fundamental usage of GPHC for checking Git repository health using the `git hc` command.

## Installation

### Prerequisites
- Go 1.19 or higher
- Git repository

### Install GPHC
```bash
# Install from GitHub
go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# Or install from local source
git clone https://github.com/vahidaghazadeh/gphc.git
cd gphc
go install ./cmd/gphc

# Setup as git subcommand
./setup-git-hc.sh
```

**Note:** Make sure your `$GOPATH/bin` is in your `$PATH`. Add this to your shell profile:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## Basic Commands

### Health Check
```bash
# Check current directory (must be a git repository)
git hc check

# Check specific repository
git hc check /path/to/repository

# Check with custom output format
git hc check --format json
git hc check --format yaml
git hc check --format markdown
git hc check --format html

# Save output to file
git hc check --output health-report.json
```

### Pre-commit Hooks
```bash
# Run pre-commit checks on staged files
git hc pre-commit

# This command will:
# - Check staged files for formatting issues
# - Validate commit message format
# - Detect large files (>1MB)
# - Check for sensitive files
# - Return appropriate exit codes for CI/CD
```

### Interactive Terminal UI
```bash
# Launch interactive terminal interface
git hc tui

# Start TUI with specific repository
git hc tui /path/to/repository
```

### Web Dashboard
```bash
# Start web dashboard server
git hc serve

# Start with custom port
git hc serve --port 3000

# Start with custom host and port
git hc serve --host 0.0.0.0 --port 8080
```

### Multi-Repository Scan
```bash
# Scan current directory
git hc scan

# Scan specific directory
git hc scan ~/projects

# Recursive scan
git hc scan ~/projects --recursive

# Scan with minimum score threshold
git hc scan ~/projects --min-score 80
```

### Utility Commands
```bash
# Update GPHC to latest version
git hc update

# Show version information
git hc version

# Show help
git hc --help
git hc check --help
git hc serve --help
```

## Understanding Health Scores

### Score Ranges
- **90-100**: Excellent (A+, A, A-)
- **80-89**: Good (B+, B, B-)
- **70-79**: Fair (C+, C, C-)
- **60-69**: Poor (D+, D, D-)
- **0-59**: Failing (F)

### Health Categories
1. **Documentation & Project Structure** (25 points)
2. **Commit History Quality** (30 points)
3. **Git Cleanup & Hygiene** (25 points)
4. **Codebase Structure** (20 points)

## Example Output

```bash
$ git hc check

Git Project Health Checker
==========================

Repository: /path/to/project
Last Updated: 2024-01-15 10:30:00

Overall Health Score: 85/100 (B+)
Status: PASS

Documentation & Project Structure: 90/100 (A-)
Commit History Quality: 85/100 (B+)
Git Cleanup & Hygiene: 80/100 (B-)
Codebase Structure: 75/100 (C+)

Summary:
  Total Checks: 12
  Passed: 8
  Failed: 2
  Warnings: 2
```

## Configuration

Create a `git-hc.yml` file in your repository root:

```yaml
# Basic configuration
health_check:
  min_score: 70
  fail_on_warnings: false

# Custom rules
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5

# Server configuration
server:
  port: 8080
  host: "localhost"
  auth:
    enabled: false
```

## Troubleshooting

### Common Issues

#### Command Not Found
```bash
# If git hc command is not found
git config --global alias.hc

# If alias is not set, run setup script
./setup-git-hc.sh

# Check if wrapper script exists
ls -la ~/.local/bin/git-hc-wrapper
```

#### Repository Not Found
```bash
# Make sure you're in a Git repository
git status

# Or specify the repository path
git hc check /path/to/git/repository
```

#### Permission Issues
```bash
# Make sure you have read access to the repository
ls -la /path/to/repository

# Check Git permissions
git log --oneline -5
```

## Next Steps

- [Health Checks](health-checks.md) - Understanding health check categories
- [Pre-commit Hooks](pre-commit-hooks.md) - Pre-commit integration guide
- [Configuration](configuration.md) - Advanced configuration options
