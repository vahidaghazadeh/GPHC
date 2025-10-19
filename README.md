# Git Project Health Checker (GPHC)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-orange.svg)](https://github.com/opsource/gphc)

**GPHC** (pronounced "githlth") is a Command-Line Interface (CLI) tool written in Go that audits local Git repositories against established Open Source best practices. It evaluates documentation quality, commit history standards, and repository hygiene, providing a comprehensive Health Score with actionable feedback.

## 🌟 Features

### 📚 Documentation & Project Structure
- **Essential Files Check**: Validates presence of README.md, LICENSE, CONTRIBUTING.md, and CODE_OF_CONDUCT.md
- **Setup Instructions**: Ensures clear installation and usage instructions
- **Gitignore Validation**: Checks for proper .gitignore configuration with common patterns

### 📝 Commit History Quality
- **Conventional Commits**: Validates adherence to conventional commit format (feat:, fix:, etc.)
- **Message Length**: Ensures commit messages stay within 72-character limit
- **Commit Size Analysis**: Identifies oversized commits that might indicate "God Commits"

### 🧹 Git Cleanup & Hygiene
- **Local Branch Cleanup**: Identifies merged branches that should be deleted
- **Stale Branch Detection**: Finds branches with no activity for 60+ days
- **Branch Protection**: Checks for main branch protection (requires GitHub API)
- **Stash Management**: Analyzes Git stash entries and warns about old stashes (>30 days)

## 🚀 Installation

### Prerequisites
- Go 1.19 or higher
- Git repository

### Method 1: Install Globally (Recommended)
```bash
# Install from GitHub
go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# Or install from local source
git clone https://github.com/vahidaghazadeh/gphc.git
cd gphc
go install ./cmd/gphc
```

**Note:** Make sure your `$GOPATH/bin` is in your `$PATH`. Add this to your shell profile:
```bash
# For zsh (macOS/Linux)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc

# For bash (Linux)
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.bashrc
source ~/.bashrc
```

### Method 2: Build from Source
```bash
git clone https://github.com/vahidaghazadeh/gphc.git
cd gphc
go build -o gphc cmd/gphc/main.go

# Use with full path
./gphc check
```

### Verify Installation
```bash
gphc version
# Should output: GPHC (Git Project Health Checker) v1.0.0
```

## 📖 Usage

### Basic Usage
```bash
# Check current directory (must be a git repository)
gphc check

# Check specific repository
gphc check /path/to/repository

# Run pre-commit checks on staged files
gphc pre-commit

# Update GPHC to latest version
gphc update

# Show version information
gphc version

# Show help
gphc --help
```

### Pre-Commit Hook Mode

GPHC includes a fast pre-commit mode designed for integration with pre-commit frameworks and Husky:

```bash
# Run pre-commit checks on staged files
gphc pre-commit
gphc badge
gphc github
gphc gitlab
gphc authors
gphc codebase
```

**Export Formats:**
```bash
# JSON output
gphc check --format json

# YAML output  
gphc check --format yaml

# Markdown output (perfect for README)
gphc check --format markdown --output health-report.md

# HTML output
gphc check --format html --output health-report.html
```

**Health Badge:**
```bash
# Generate badge URL and markdown
gphc badge

# Output example:
# 🔗 Badge URL:
# https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github
# 
# 📝 Markdown Badge:
# ![Health Score](https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github)
```

**Pre-Commit Hook Mode:**
```bash
gphc pre-commit
```

**Features:**
- ✅ **Fast execution:** Only checks staged files
- ✅ **File formatting:** Validates Go code formatting
- ✅ **Commit message:** Checks conventional commit format
- ✅ **Large files:** Prevents files >1MB from being committed
- ✅ **Sensitive files:** Blocks .env, keys, credentials
- ✅ **Exit codes:** Returns non-zero for failed checks

**Integration Examples:**

**Pre-commit framework:**
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: gphc-pre-commit
        name: GPHC Pre-commit Check
        entry: gphc pre-commit
        language: system
        pass_filenames: false
```

**Husky (Node.js):**
```json
// package.json
{
  "husky": {
    "hooks": {
      "pre-commit": "gphc pre-commit"
    }
  }
}
```

### Troubleshooting

**Command not found error:**
```bash
# If you get "command not found: gphc"
# Make sure GOPATH/bin is in your PATH
echo $PATH | grep -q "$(go env GOPATH)/bin" || echo "PATH issue detected"

# Add to PATH if missing
export PATH=$PATH:$(go env GOPATH)/bin
```

**Permission denied:**
```bash
# If you get permission errors, make sure the binary is executable
chmod +x $(go env GOPATH)/bin/gphc
```

### Example Output
```
🔍 Analyzing repository: /path/to/repo

✅ Repository Health Check (GPHC v1.0.0)

🌟 Overall Health Score: 85/100 (B+)

---------------------------------------------------
[A] Documentation & Project Structure (Passed: 3/4)
---------------------------------------------------
   ✅ DOC-101: README.md found (Score: +10)
   ✅ DOC-102: LICENSE file found (Score: +10)
   ❌ DOC-103: CONTRIBUTING.md is missing (Deduct: -10)
   ✅ IG-201: .gitignore is present and valid (Score: +10)

---------------------------------------------------
[B] Commit History Quality (Passed: 8/10)
---------------------------------------------------
   ⚠️ CHQ-301: 2 of 10 recent commits violate Conventional Commit standard (Deduct: -5)
      - Non-Standard Commit: "Initial work"
   ✅ CHQ-302: Commit message length is compliant (Avg. 55 chars) (Score: +10)
   ✅ CHQ-303: Average commit size is moderate (Avg. 120 lines) (Score: +10)

---------------------------------------------------
[C] Git Cleanup & Hygiene (Needs Attention)
---------------------------------------------------
   ❌ CLEAN-401: 3 local branches are merged but not deleted (Deduct: -10)
   ⚠️ CLEAN-402: Branch 'experiment-beta' is stale (last activity: 95 days ago) (Deduct: -5)
   ⚠️ STASH-501: Found 2 stash entries (1 old) (Deduct: -5)
      ⚠️ Old stash@{0}: WIP on feature-branch (45 days ago) [feature-branch]
      ✅ Recent stash@{1}: Quick fix attempt (2 days ago) [main]

💡 Next Steps:
   1. Create CONTRIBUTING.md
   2. Delete 3 stale local branches
   3. Review and clean up old stash entries
```

## 📊 Export Formats & Badges

GPHC supports multiple output formats for integration with CI/CD pipelines, documentation, and reporting systems.

### Export Formats

**JSON Output:**
```bash
gphc check --format json
```
Perfect for CI/CD integration and automated processing.

**YAML Output:**
```bash
gphc check --format yaml
```
Human-readable format for configuration files.

**Markdown Output:**
```bash
gphc check --format markdown --output health-report.md
```
Ready-to-use markdown for README files or documentation.

**HTML Output:**
```bash
gphc check --format html --output health-report.html
```
Beautiful HTML reports for web dashboards.

### Health Badges

Generate shields.io-style badges for your repository:

```bash
gphc badge
```

**Output:**
```
📊 Health Score: 85/100 (B+)

🔗 Badge URL:
https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github

📝 Markdown Badge:
![Health Score](https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github)
```

**Badge Colors:**
- 🟢 **90-100:** `brightgreen` (A+)
- 🟢 **80-89:** `green` (A, B+)
- 🟡 **70-79:** `yellowgreen` (B, B-)
- 🟡 **60-69:** `yellow` (C+, C)
- 🟠 **50-59:** `orange` (C-, D+)
- 🔴 **0-49:** `red` (D, F)

**Add to README:**
```markdown
# My Project

![Health Score](https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github)

A well-maintained project with excellent health metrics.
```

## 🔗 GitHub Integration

GPHC provides deep integration with GitHub repositories to check advanced features and configurations.

### Setup

Set your GitHub Personal Access Token:

```bash
# Option 1: Using GPHC_TOKEN (recommended)
export GPHC_TOKEN=your_github_token

# Option 2: Using GITHUB_TOKEN (also supported)
export GITHUB_TOKEN=your_github_token
```

**Required Token Permissions:**
- `repo` (Full control of private repositories)
- `read:org` (Read org and team membership)

### GitHub Integration Check

```bash
# Check GitHub integration features
gphc github

# Check specific repository
gphc github /path/to/repository
```

**Features Checked:**
- ✅ **Branch Protection:** Required reviewers, status checks, code owner reviews
- ✅ **GitHub Actions:** Workflow configuration and status
- ✅ **Repository Settings:** Issues, projects, wiki enabled
- ✅ **Contributors:** Multi-contributor analysis
- ✅ **Repository Info:** Stars, forks, activity metrics

**Example Output:**
```
🔍 Checking GitHub integration: /path/to/repo
✅ GitHub token found

📊 GitHub Integration Check Results:
Status: PASS
Score: 85
Message: Excellent GitHub integration and configuration

Details:
  Repository: owner/repository
  ✅ Issues are enabled
  ✅ Projects are enabled
  ✅ Wiki is enabled
  ✅ Branch protection is enabled
  ✅ Required 2 reviewer(s)
  ✅ Code owner reviews required
  ✅ Required status checks: ci, test, build
  ✅ Found 3 workflow(s)
  ✅ 2 active workflow(s)
  📊 Found 5 contributor(s)
  ✅ Multiple contributors
```

### Integration with Health Check

GitHub integration is automatically included in the main health check:

```bash
gphc check
```

The GitHub integration checker (`GH-601`) will appear in the results when:
- Repository is hosted on GitHub
- GitHub token is available
- Repository is accessible via GitHub API

## 🔗 GitLab Integration

GPHC provides comprehensive integration with GitLab repositories to check advanced features and configurations.

### Setup

Set your GitLab Personal Access Token:

```bash
# Option 1: Using GPHC_TOKEN (recommended)
export GPHC_TOKEN=your_gitlab_token

# Option 2: Using GITLAB_TOKEN (also supported)
export GITLAB_TOKEN=your_gitlab_token

# For custom GitLab instances
export GITLAB_URL=https://your-gitlab-instance.com
```

**Required Token Permissions:**
- `api` (Full API access)
- `read_repository` (Read repository data)
- `read_user` (Read user information)

### GitLab Integration Check

```bash
# Check GitLab integration features
gphc gitlab

# Check specific repository
gphc gitlab /path/to/repository
```

**Features Checked:**
- ✅ **Branch Protection:** Push/merge access restrictions, code owner approval
- ✅ **GitLab CI/CD:** Pipeline configuration and status
- ✅ **Project Settings:** Issues, merge requests, wiki, snippets enabled
- ✅ **Contributors:** Multi-contributor analysis
- ✅ **Merge Requests:** Open MRs and development activity
- ✅ **Project Info:** Stars, forks, activity metrics

**Example Output:**
```
🔍 Checking GitLab integration: /path/to/repo
✅ GitLab token found

📊 GitLab Integration Check Results:
Status: PASS
Score: 85
Message: Excellent GitLab integration and configuration

Details:
  Project: owner/repository
  ✅ Issues are enabled
  ✅ Merge requests are enabled
  ✅ Wiki is enabled
  ✅ Snippets are enabled
  ✅ Branch protection is enabled
  ✅ Push access is restricted
  ✅ Merge access is restricted
  ✅ Code owner approval required
  ✅ Found 3 pipeline(s)
  ✅ 2 successful pipeline(s)
  📊 Found 5 contributor(s)
  ✅ Multiple contributors
  📊 Found 2 open merge request(s)
  ✅ Active development with open MRs
```

### Integration with Health Check

GitLab integration is automatically included in the main health check:

```bash
gphc check
```

The GitLab integration checker (`GL-602`) will appear in the results when:
- Repository is hosted on GitLab
- GitLab token is available
- Repository is accessible via GitLab API

## 👥 Commit Author Insights

GPHC analyzes commit history to identify contributor patterns and bus factor risks, helping teams understand project dependencies and team participation.

### Author Analysis

```bash
# Analyze commit authors and bus factor risk
gphc authors

# Analyze specific repository
gphc authors /path/to/repository
```

**Features Analyzed:**
- ✅ **Contributor Count:** Total number of unique contributors
- ✅ **Commit Distribution:** Percentage of commits per author
- ✅ **Single Author Dominance:** Detection of >70% contribution by one person
- ✅ **Bus Factor Risk:** Assessment of project dependency on individuals
- ✅ **Email Consistency:** Validation of author email addresses
- ✅ **Team Participation:** Analysis of contributor engagement

**Example Output:**
```
👥 Analyzing commit authors: /path/to/repo

📊 Commit Author Insights:
Status: WARNING
Score: 60
Message: Single author dominance detected

Details:
  👥 Contributors: 4
  📊 Total commits: 31
  1. vahidaghazadeh (24 commits, 77.4%)
  2. john.doe (4 commits, 12.9%)
  3. jane.smith (2 commits, 6.5%)
  4. mike.wilson (1 commits, 3.2%)
  ⚠️ Single Author Dominance Detected (>70%)
  Top author: vahidaghazadeh (77.4%)
  Consider encouraging more team participation
  ✅ Email addresses are consistent

💡 Bus Factor Analysis:
  ⚠️ MODERATE RISK: Low contributor count
  📊 Contributors: 4
  🚨 Bus Factor: 4 (Acceptable)
  💡 Recommendation: Maintain current team size
```

**Bus Factor Risk Levels:**
- **🚨 Critical (1 contributor):** Single person project - immediate action needed
- **⚠️ High (2 contributors):** Low contributor count - expand team
- **✅ Acceptable (3-5 contributors):** Small team - maintain current size
- **✅ Excellent (6+ contributors):** Well-distributed team - excellent distribution

### Integration with Health Check

Author insights are automatically included in the main health check:

```bash
gphc check
```

The author insights checker (`CAI-701`) will appear in the results and provides:
- Contributor distribution analysis
- Single author dominance warnings
- Bus factor risk assessment
- Team participation recommendations

## 🔍 Codebase Smell Check

GPHC performs lightweight codebase structure analysis to detect common organizational issues and maintainability problems without requiring AST parsing.

### Structure Analysis

```bash
# Analyze codebase structure and detect code smells
gphc codebase

# Analyze specific repository
gphc codebase /path/to/repository
```

**Features Analyzed:**
- ✅ **Test Directory Detection:** Missing test directories and files
- ✅ **Directory Size Analysis:** Oversized directories (>1000 files)
- ✅ **Code-to-Test Ratio:** Test coverage ratio analysis
- ✅ **Empty Directory Detection:** Unused empty directories
- ✅ **Directory Depth Analysis:** Deep nesting detection
- ✅ **Root File Count:** Too many files in root directory
- ✅ **Standard Directory Structure:** Missing src/, lib/, app/ directories
- ✅ **Documentation Files:** Presence of documentation files

**Example Output:**
```
🔍 Analyzing codebase structure: /path/to/repo

📊 Codebase Structure Analysis:
Status: WARNING
Score: 65
Message: Codebase structure warnings detected

Details:
  ⚠️ No test directory found
  Consider adding tests/ or test/ directory
  ✅ No oversized directories
  ⚠️ Low test coverage ratio (5.2%)
  Consider adding more test files
  ℹ️ 2 empty directory(ies) found
  Consider removing empty directories
  ✅ Reasonable directory depth (4)
  ✅ Reasonable number of root files (8)
  ℹ️ No standard source directories (src/, lib/, app/)
  Consider organizing code into standard directories
  ✅ Found 3 documentation file(s)

📊 Codebase Statistics:
  📁 Total directories: 12
  📄 Total files: 45
  🧪 Test files: 2
  📝 Documentation files: 3
  📏 Max directory depth: 4

💡 Structure Recommendations:
  ⚠️ Codebase structure needs improvement
  📋 Consider the following actions:
    • Add test directories and test files
    • Organize code into logical subdirectories
    • Split oversized directories (>1000 files)
    • Add documentation files
    • Remove empty directories
```

**Structure Quality Levels:**
- **🚨 Poor (<70):** Significant structural issues - immediate action needed
- **⚠️ Fair (70-89):** Minor improvements needed - good foundation
- **✅ Excellent (90+):** Well-organized codebase - maintain current structure

### Integration with Health Check

Codebase structure analysis is automatically included in the main health check:

```bash
gphc check
```

The codebase smell checker (`CBS-801`) will appear in the results and provides:
- Directory structure analysis
- Test coverage assessment
- Organization recommendations
- Maintainability insights

## ⚙️ Configuration

Create a `gphc.yml` file in your repository root to customize settings:

```yaml
# Commit analysis settings
max_commits_to_analyze: 50

# Branch analysis settings
stale_branch_threshold_days: 60

# Commit message settings
max_commit_message_length: 72

# Commit size settings
max_commit_size_lines: 500

# Scoring weights (1-10)
weights:
  documentation: 3
  commits: 4
  hygiene: 2
```

## 🏗️ Architecture

GPHC follows a modular architecture with discrete checkers reporting to a central scoring engine:

```
Input: gphc check [path]
    ↓
Data Collector: go-git repository analysis
    ↓
Checker Modules: Independent Go structs implementing Checker interface
    ↓
Scoring Engine: Aggregates results with weighted scoring
    ↓
Reporter: Colorful terminal output with structured results
```

### Core Components

- **`internal/git/`**: Repository data collection and analysis
- **`internal/checkers/`**: Individual health check implementations
- **`internal/scorer/`**: Scoring engine and health report generation
- **`internal/reporter/`**: Terminal output formatting
- **`pkg/types/`**: Core data structures and interfaces
- **`pkg/config/`**: Configuration management

## 🔧 Development

### Setting Up Development Environment
```bash
# Clone the repository
git clone https://github.com/vahidaghazadeh/gphc.git
cd gphc

# Install dependencies
go mod download

# Build the project
go build -o gphc cmd/gphc/main.go

# Test locally
./gphc check

# Install for development
go install ./cmd/gphc
```

### Project Structure
```
gphc/
├── cmd/gphc/           # CLI entry point
├── internal/
│   ├── checkers/       # Health check implementations
│   ├── git/           # Git repository analysis
│   ├── scorer/        # Scoring engine
│   └── reporter/      # Output formatting
├── pkg/
│   ├── types/         # Core data structures
│   └── config/        # Configuration management
├── gphc.yml          # Default configuration
├── .gitignore        # Git ignore patterns
└── README.md
```

### Building and Testing
```bash
# Build binary
go build -o gphc cmd/gphc/main.go

# Run tests (when available)
go test ./...

# Run with race detection
go run -race cmd/gphc/main.go check

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o gphc-linux cmd/gphc/main.go
GOOS=windows GOARCH=amd64 go build -o gphc.exe cmd/gphc/main.go
```

### Adding New Checkers

1. Implement the `Checker` interface in `internal/checkers/`
2. Add your checker to the main checker list in `cmd/gphc/main.go`
3. Define appropriate scoring weights in the configuration

### Example Checker Implementation
```go
type MyChecker struct {
    BaseChecker
}

func NewMyChecker() *MyChecker {
    return &MyChecker{
        BaseChecker: NewBaseChecker("My Checker", "MY", types.CategoryDocs, 5),
    }
}

func (mc *MyChecker) Check(data *types.RepositoryData) *types.CheckResult {
    // Implementation here
}
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📋 Roadmap

### Phase 1: Core Features ✅
- [x] Basic CLI framework
- [x] Git repository analysis
- [x] Documentation checks
- [x] Commit quality analysis
- [x] Branch hygiene checks
- [x] Scoring engine
- [x] Colorful terminal output

### Phase 2: Enhanced Features
- [ ] GitHub API integration for branch protection
- [ ] Pre-commit hook validation
- [ ] Custom rule definitions
- [ ] JSON/XML output formats
- [ ] CI/CD integration

### Phase 3: Advanced Features
- [ ] Multi-repository analysis
- [ ] Historical trend analysis
- [ ] Team collaboration metrics
- [ ] Integration with popular Git hosting platforms

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [go-git](https://github.com/go-git/go-git) for Git repository access
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [LipGloss](https://github.com/charmbracelet/lipgloss) for terminal styling
- [Viper](https://github.com/spf13/viper) for configuration management

---

**Made with ❤️ for the Open Source community**