# Git Project Health Checker (GPHC)

[![Go Version](https://img.shields.io/badge/go-1.24-blue.svg)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/vahidaghazadeh/gphc)](https://github.com/vahidaghazadeh/gphc/releases)

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

# Setup as git subcommand
./setup-git-hc.sh
```

### Updating
```bash
# Preferred: update the binary used by git hc
git hc update

# If your installed version does not have update yet, install into the
# directory used by the git hc wrapper.
GOBIN="$HOME/.local/bin" go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# Verify the active version
git hc version
```

If `git hc update` fails with an old updater error such as
`Could not find GPHC source directory`, install the new binary directly into the
wrapper directory:

```bash
mkdir -p ~/.local/bin
GOBIN="$HOME/.local/bin" go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
hash -r
git hc version
```

If the latest changes are pushed to `main` but have not been released/tagged
yet, install from `main`:

```bash
GOBIN="$HOME/.local/bin" go install github.com/vahidaghazadeh/gphc/cmd/gphc@main
hash -r
git hc version
```

`@latest` installs the latest published tag/release. Local changes are not
available on other machines until they are committed, tagged, and pushed.

### Basic Usage
```bash
# Check current directory (must be a git repository)
git hc check

# Check specific repository
git hc check /path/to/repository

# Run pre-commit checks on staged files
git hc pre-commit

# Launch interactive terminal UI
git hc tui

# Start web dashboard server
git hc serve

# Scan multiple repositories
git hc scan ~/projects --recursive

# Analyze and manage Git tags
git hc tags --suggest --changelog CHANGELOG.md

# Scan for secrets in Git history
git hc security secrets --history

# Scan transitive dependencies for vulnerabilities
git hc security dependencies --depth deep

# Validate Git security policies
git hc security policy --check-signing

# Audit executable and large files
git hc security binaries --max-size 50mb

# Update GPHC to latest version
git hc update

# Show version information
git hc version

# Show help
git hc --help
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
- **Tag Management**: Git tag validation, semantic versioning, and release management
- **Secret Scanning**: Deep scan of Git history for exposed secrets and credentials
- **Transitive Dependency Vetting**: Comprehensive analysis of direct and indirect dependencies for security vulnerabilities
- **Git Policy Validation**: Validate Git security policies including commit signatures, push policies, and sensitive file detection
- **Binary File Audit**: Scan for executable files, large files, and suspicious file types that pose security risks

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
- [🏷️ Tag Management](docs/tag-management.md) - Git tag validation and release management
- [🔒 Secret Scanning](docs/secret-scanning.md) - Git history secret detection and remediation
- [🛡️ Transitive Dependency Vetting](docs/transitive-dependency-vetting.md) - Deep dependency vulnerability analysis
- [⚙️ Git Policy Validation](docs/git-policy-validation.md) - Git security policy validation and compliance
- [🔍 Binary File Audit](docs/binary-file-audit.md) - Executable and large file security audit

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

### Phase 3: Team & Automation
- [x] Multi-repository analysis
- [x] Integration with popular Git hosting platforms
- [x] Interactive Terminal UI (TUI)
- [x] Web Dashboard
- [ ] Historical trend analysis
- [ ] Team collaboration metrics

---

**Made with ❤️ for the Open Source community**
