# Git Project Health Checker (GPHC)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-orange.svg)](https://github.com/opsource/gphc)

**GPHC** (pronounced "githlth") is a Command-Line Interface (CLI) tool written in Go that audits local Git repositories against established Open Source best practices. It evaluates documentation quality, commit history standards, and repository hygiene, providing a comprehensive Health Score with actionable feedback.

## Quick Start

### Installation
```bash
# Install from GitHub
go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# Or install from local source
git clone https://github.com/vahidaghazadeh/gphc.git
cd gphc
go install ./cmd/gphc

# Setup as git subcommand (optional)
./setup-git-hc.sh
```

### Basic Usage
```bash
# Using gphc directly
gphc check                    # Check current directory (must be a git repository)
gphc check /path/to/repository # Check specific repository
gphc pre-commit               # Run pre-commit checks on staged files
gphc tui                      # Launch interactive terminal UI
gphc serve                    # Start web dashboard server
gphc scan ~/projects --recursive # Scan multiple repositories
gphc update                   # Update GPHC to latest version
gphc version                  # Show version information
gphc --help                   # Show help

# Using git hc (after running setup-git-hc.sh)
git hc check                  # Check current directory
git hc pre-commit             # Run pre-commit checks
git hc tui                    # Launch interactive terminal UI
git hc serve                  # Start web dashboard server
git hc scan ~/projects --recursive # Scan multiple repositories
git hc update                 # Update GPHC to latest version
git hc version                # Show version information
git hc --help                 # Show help
```

## Features Overview

### Core Features
- **Documentation & Project Structure**: Essential files validation, setup instructions, gitignore checks
- **Commit History Quality**: Conventional commits, message length validation, commit size analysis
- **Git Cleanup & Hygiene**: Branch cleanup, stale branch detection, stash management
- **Pre-commit Hook Mode**: Fast execution for staged files with exit codes for CI/CD

### Advanced Features
- **Historical Health Tracking**: Track project health over time with trend analysis
- **Multi-Repository Scan**: Analyze multiple repositories simultaneously
- **CI/CD Integration**: Quality gates and pipeline integration
- **Custom Rules Engine**: Define project-specific health checks
- **Slack/Webhook Notifications**: Team notifications and real-time updates
- **Semantic Commit Verification**: Verify commit messages match actual changes
- **Interactive Terminal UI (TUI)**: Beautiful terminal interface for health monitoring
- **Web Dashboard**: Local web server for team collaboration

## Documentation

Detailed documentation for each feature is available in the `docs/` directory:

- [📋 Basic Usage](docs/basic-usage.md) - Getting started with GPHC
- [🔧 Git HC Integration](docs/git-hc-integration.md) - Using GPHC as git subcommand
- [📊 Health Checks](docs/health-checks.md) - Understanding health check categories
- [🔧 Pre-commit Hooks](docs/pre-commit-hooks.md) - Pre-commit integration guide
- [📈 Historical Tracking](docs/historical-tracking.md) - Health trend analysis
- [🔍 Multi-Repository Scan](docs/multi-repository-scan.md) - Batch repository analysis
- [🚀 CI/CD Integration](docs/ci-cd-integration.md) - Pipeline integration guide
- [⚙️ Custom Rules](docs/custom-rules.md) - Custom rule engine configuration
- [📢 Notifications](docs/notifications.md) - Slack and webhook setup
- [✅ Semantic Commits](docs/semantic-commits.md) - Commit verification guide
- [🖥️ Terminal UI](docs/terminal-ui.md) - Interactive terminal interface
- [🌐 Web Dashboard](docs/web-dashboard.md) - Web server and team collaboration
- [📤 Export Formats](docs/export-formats.md) - Report export options
- [🔗 GitHub Integration](docs/github-integration.md) - GitHub API integration
- [🔗 GitLab Integration](docs/gitlab-integration.md) - GitLab API integration
- [👥 Author Insights](docs/author-insights.md) - Contributor analysis
- [🏗️ Codebase Analysis](docs/codebase-analysis.md) - Structure and smell detection

## Example Output

```bash
$ gphc check

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

Create a `gphc.yml` file in your repository root to customize behavior:

```yaml
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

# Dashboard settings
dashboard:
  title: "My Project Dashboard"
  theme: "dark"
  refresh_interval: "30s"
```

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- 📖 [Documentation](docs/) - Comprehensive guides for all features
- 🐛 [Issues](https://github.com/vahidaghazadeh/gphc/issues) - Report bugs or request features
- 💬 [Discussions](https://github.com/vahidaghazadeh/gphc/discussions) - Community discussions
- 📧 [Contact](mailto:support@gphc.dev) - Direct support

## Roadmap

### Phase 1: Core Features ✅
- [x] Basic health checks
- [x] Documentation validation
- [x] Commit history analysis
- [x] Git hygiene checks
- [x] Pre-commit hook mode

### Phase 2: Advanced Features ✅
- [x] Export formats (JSON, YAML, Markdown, HTML)
- [x] Health badges
- [x] GitHub/GitLab integration
- [x] Commit author insights
- [x] Codebase smell detection

### Phase 3: Team, Trends & Automation ✅
- [x] Multi-repository analysis
- [x] Historical trend analysis
- [x] Team collaboration metrics
- [x] Integration with popular Git hosting platforms
- [x] Interactive Terminal UI (TUI)
- [x] Web Dashboard

---

**Made with ❤️ for the Open Source community**