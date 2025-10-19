# Git Project Health Checker (GPHC)

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.19-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-1.0.0-orange.svg)](https://github.com/opsource/gphc)

**GPHC** (pronounced "githlth") is a Command-Line Interface (CLI) tool written in Go that audits local Git repositories against established Open Source best practices. It evaluates documentation quality, commit history standards, and repository hygiene, providing a comprehensive Health Score with actionable feedback.

## Features

### Documentation & Project Structure
- **Essential Files Check**: Validates presence of README.md, LICENSE, CONTRIBUTING.md, and CODE_OF_CONDUCT.md
- **Setup Instructions**: Ensures clear installation and usage instructions
- **Gitignore Validation**: Checks for proper .gitignore configuration with common patterns

### Commit History Quality
- **Conventional Commits**: Validates adherence to conventional commit format (feat:, fix:, etc.)
- **Message Length**: Ensures commit messages stay within 72-character limit
- **Commit Size Analysis**: Identifies oversized commits that might indicate "God Commits"

### Git Cleanup & Hygiene
- **Local Branch Cleanup**: Identifies merged branches that should be deleted
- **Stale Branch Detection**: Finds branches with no activity for 60+ days
- **Branch Protection**: Checks for main branch protection (requires GitHub API)
- **Stash Management**: Analyzes Git stash entries and warns about old stashes (>30 days)

### Historical Health Tracking
- **Health History**: Automatically saves health scores to `.gphc-history.json`
- **Trend Analysis**: Shows project health improvement over time
- **CI Integration**: Track quality metrics in continuous integration
- **Team Insights**: Monitor team progress and code quality trends

### Multi-Repository Scan
- **Batch Analysis**: Scan multiple repositories simultaneously
- **Recursive Scanning**: Find and analyze all Git repositories in directories
- **Aggregate Reporting**: Compare health scores across projects
- **Organization Overview**: Get company-wide code quality insights

### CI/CD Integration
- **Automated Health Checks**: Integrate with CI/CD pipelines for quality gates
- **Quality Thresholds**: Set minimum health score requirements
- **Pipeline Integration**: Fail builds when quality standards aren't met
- **Continuous Monitoring**: Track project health in every build

### Custom Rules (Rule Engine)
- **Custom Checks**: Define project-specific health checks in `gphc.yml`
- **File Existence Rules**: Check for required files (SECURITY.md, CONTRIBUTING.md, etc.)
- **Regex Pattern Matching**: Search for patterns in files using regular expressions
- **Flexible Scoring**: Assign custom scores to custom rules
- **Project-Specific Policies**: Enforce organization-specific requirements

### Slack / Webhook Notifications
- **Team Notifications**: Send health reports directly to Slack or Discord channels
- **Real-time Updates**: Get immediate notifications about repository health status
- **Customizable Messages**: Format notifications with project-specific information
- **Webhook Integration**: Support for any webhook-compatible service
- **Team Collaboration**: Keep entire team informed about code quality trends

### Semantic Commit Verification
- **Commit Message Analysis**: Verify commit messages match actual code changes
- **Change Type Detection**: Analyze file modifications to determine change type
- **Mismatch Detection**: Identify commits with misleading commit messages
- **Semantic Consistency**: Ensure commit history accurately reflects changes
- **Quality Assurance**: Prevent confusing commit history and maintain transparency

### Interactive Terminal UI (TUI)
- **Graphical Terminal Interface**: Beautiful terminal-based user interface like htop
- **Colorful Score Display**: Interactive and colorful score visualization
- **Real-time Filtering**: Filter results by category, status, or score range
- **Rule Explanations**: View detailed explanations for each health check rule
- **Trend Navigation**: Browse through historical score trends interactively
- **Enhanced Developer Experience**: Improved user experience for developers and technical teams

### Web Dashboard
- **Local Web Server**: Launch local HTML dashboard with `gphc serve`
- **Multi-Project Monitoring**: View health scores for multiple repositories simultaneously
- **Trend Analysis**: Historical health tracking and trend visualization
- **Export Capabilities**: Generate comprehensive reports in various formats
- **Team Collaboration**: Shared dashboard accessible to entire team
- **Real-time Updates**: Live health monitoring with automatic refresh

## Installation

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

# Launch interactive terminal UI
gphc tui

# Start web dashboard server
gphc serve

# Scan multiple repositories
gphc scan ~/projects --recursive

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
gphc trend
gphc scan
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
# Markdown Badge:
# ![Health Score](https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github)
```

**Pre-Commit Hook Mode:**
```bash
gphc pre-commit
```

**Features:**
- **Fast execution:** Only checks staged files
- **File formatting:** Validates Go code formatting
- **Commit message:** Checks conventional commit format
- **Large files:** Prevents files >1MB from being committed
- **Sensitive files:** Blocks .env, keys, credentials
- **Exit codes:** Returns non-zero for failed checks

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
Analyzing repository: /path/to/repo

Repository Health Check (GPHC v1.0.0)

Overall Health Score: 85/100 (B+)

---------------------------------------------------
[A] Documentation & Project Structure (Passed: 3/4)
---------------------------------------------------
   PASS DOC-101: README.md found (Score: +10)
   PASS DOC-102: LICENSE file found (Score: +10)
   FAIL DOC-103: CONTRIBUTING.md is missing (Deduct: -10)
   PASS IG-201: .gitignore is present and valid (Score: +10)

---------------------------------------------------
[B] Commit History Quality (Passed: 8/10)
---------------------------------------------------
   WARN CHQ-301: 2 of 10 recent commits violate Conventional Commit standard (Deduct: -5)
      - Non-Standard Commit: "Initial work"
   PASS CHQ-302: Commit message length is compliant (Avg. 55 chars) (Score: +10)
   PASS CHQ-303: Average commit size is moderate (Avg. 120 lines) (Score: +10)

---------------------------------------------------
[C] Git Cleanup & Hygiene (Needs Attention)
---------------------------------------------------
   FAIL CLEAN-401: 3 local branches are merged but not deleted (Deduct: -10)
   WARN CLEAN-402: Branch 'experiment-beta' is stale (last activity: 95 days ago) (Deduct: -5)
   WARN STASH-501: Found 2 stash entries (1 old) (Deduct: -5)
      WARN Old stash@{0}: WIP on feature-branch (45 days ago) [feature-branch]
      PASS Recent stash@{1}: Quick fix attempt (2 days ago) [main]

Next Steps:
   1. Create CONTRIBUTING.md
   2. Delete 3 stale local branches
   3. Review and clean up old stash entries
```

## Export Formats & Badges

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
Health Score: 85/100 (B+)

Badge URL:
https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github

Markdown Badge:
![Health Score](https://img.shields.io/badge/Health_Score-85%2F100-green?style=for-the-badge&logo=github)
```

**Badge Colors:**
- **90-100:** `brightgreen` (A+)
- **80-89:** `green` (A, B+)
- **70-79:** `yellowgreen` (B, B-)
- **60-69:** `yellow` (C+, C)
- **50-59:** `orange` (C-, D+)
- **0-49:** `red` (D, F)

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
- **Branch Protection:** Required reviewers, status checks, code owner reviews
- **GitHub Actions:** Workflow configuration and status
- **Repository Settings:** Issues, projects, wiki enabled
- **Contributors:** Multi-contributor analysis
- **Repository Info:** Stars, forks, activity metrics

**Example Output:**
```
Checking GitHub integration: /path/to/repo
GitHub token found

GitHub Integration Check Results:
Status: PASS
Score: 85
Message: Excellent GitHub integration and configuration

Details:
  Repository: owner/repository
  PASS Issues are enabled
  PASS Projects are enabled
  PASS Wiki is enabled
  PASS Branch protection is enabled
  PASS Required 2 reviewer(s)
  PASS Code owner reviews required
  PASS Required status checks: ci, test, build
  PASS Found 3 workflow(s)
  PASS 2 active workflow(s)
  Found 5 contributor(s)
  PASS Multiple contributors
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
- **Branch Protection:** Push/merge access restrictions, code owner approval
- **GitLab CI/CD:** Pipeline configuration and status
- **Project Settings:** Issues, merge requests, wiki, snippets enabled
- **Contributors:** Multi-contributor analysis
- **Merge Requests:** Open MRs and development activity
- **Project Info:** Stars, forks, activity metrics

**Example Output:**
```
Checking GitLab integration: /path/to/repo
GitLab token found

GitLab Integration Check Results:
Status: PASS
Score: 85
Message: Excellent GitLab integration and configuration

Details:
  Project: owner/repository
  PASS Issues are enabled
  PASS Merge requests are enabled
  PASS Wiki is enabled
  PASS Snippets are enabled
  PASS Branch protection is enabled
  PASS Push access is restricted
  PASS Merge access is restricted
  PASS Code owner approval required
  PASS Found 3 pipeline(s)
  PASS 2 successful pipeline(s)
  Found 5 contributor(s)
  PASS Multiple contributors
  Found 2 open merge request(s)
  PASS Active development with open MRs
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

## Commit Author Insights

GPHC analyzes commit history to identify contributor patterns and bus factor risks, helping teams understand project dependencies and team participation.

### Author Analysis

```bash
# Analyze commit authors and bus factor risk
gphc authors

# Analyze specific repository
gphc authors /path/to/repository
```

**Features Analyzed:**
- **Contributor Count:** Total number of unique contributors
- **Commit Distribution:** Percentage of commits per author
- **Single Author Dominance:** Detection of >70% contribution by one person
- **Bus Factor Risk:** Assessment of project dependency on individuals
- **Email Consistency:** Validation of author email addresses
- **Team Participation:** Analysis of contributor engagement

**Example Output:**
```
Analyzing commit authors: /path/to/repo

Commit Author Insights:
Status: WARNING
Score: 60
Message: Single author dominance detected

Details:
  Contributors: 4
  Total commits: 31
  1. vahidaghazadeh (24 commits, 77.4%)
  2. john.doe (4 commits, 12.9%)
  3. jane.smith (2 commits, 6.5%)
  4. mike.wilson (1 commits, 3.2%)
  WARN Single Author Dominance Detected (>70%)
  Top author: vahidaghazadeh (77.4%)
  Consider encouraging more team participation
  PASS Email addresses are consistent

Bus Factor Analysis:
  WARN MODERATE RISK: Low contributor count
  Contributors: 4
  Bus Factor: 4 (Acceptable)
  Recommendation: Maintain current team size
```

**Bus Factor Risk Levels:**
- **Critical (1 contributor):** Single person project - immediate action needed
- **High (2 contributors):** Low contributor count - expand team
- **Acceptable (3-5 contributors):** Small team - maintain current size
- **Excellent (6+ contributors):** Well-distributed team - excellent distribution

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

## Codebase Smell Check

GPHC performs lightweight codebase structure analysis to detect common organizational issues and maintainability problems without requiring AST parsing.

### Structure Analysis

```bash
# Analyze codebase structure and detect code smells
gphc codebase

# Analyze specific repository
gphc codebase /path/to/repository
```

**Features Analyzed:**
- **Test Directory Detection:** Missing test directories and files
- **Directory Size Analysis:** Oversized directories (>1000 files)
- **Code-to-Test Ratio:** Test coverage ratio analysis
- **Empty Directory Detection:** Unused empty directories
- **Directory Depth Analysis:** Deep nesting detection
- **Root File Count:** Too many files in root directory
- **Standard Directory Structure:** Missing src/, lib/, app/ directories
- **Documentation Files:** Presence of documentation files

**Example Output:**
```
Analyzing codebase structure: /path/to/repo

Codebase Structure Analysis:
Status: WARNING
Score: 65
Message: Codebase structure warnings detected

Details:
  WARN No test directory found
  Consider adding tests/ or test/ directory
  PASS No oversized directories
  WARN Low test coverage ratio (5.2%)
  Consider adding more test files
  INFO 2 empty directory(ies) found
  Consider removing empty directories
  PASS Reasonable directory depth (4)
  PASS Reasonable number of root files (8)
  INFO No standard source directories (src/, lib/, app/)
  Consider organizing code into standard directories
  PASS Found 3 documentation file(s)

Codebase Statistics:
  Total directories: 12
  Total files: 45
  Test files: 2
  Documentation files: 3
  Max directory depth: 4

Structure Recommendations:
  WARN Codebase structure needs improvement
  Consider the following actions:
    • Add test directories and test files
    • Organize code into logical subdirectories
    • Split oversized directories (>1000 files)
    • Add documentation files
    • Remove empty directories
```

**Structure Quality Levels:**
- **Poor (<70):** Significant structural issues - immediate action needed
- **Fair (70-89):** Minor improvements needed - good foundation
- **Excellent (90+):** Well-organized codebase - maintain current structure

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

## Historical Health Tracking

GPHC automatically tracks your project's health score over time, providing valuable insights into code quality trends and team progress.

### Health History Storage

Every time you run `gphc check`, the results are automatically saved to `.gphc-history.json` in your repository root:

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "score": 85,
  "grade": "B+",
  "checks_passed": 12,
  "checks_failed": 2,
  "checks_warning": 3,
  "repository": "owner/repository"
}
```

### Trend Analysis

View your project's health improvement over time:

```bash
# Show health trends
gphc trend

# Show trends for specific time period
gphc trend --days 30
gphc trend --weeks 12
gphc trend --months 6
```

**Example Output:**
```
Health Trend Analysis (Last 30 days)

Score Progression:
  Jan 01: 72/100 (C+)
  Jan 08: 75/100 (B-)
  Jan 15: 78/100 (B-)
  Jan 22: 82/100 (B+)
  Jan 29: 85/100 (B+)

Improvement: +13 points (+18.1%)
Trend: Improving
Average Score: 78.4/100

Key Improvements:
  PASS Added comprehensive test coverage (+8 points)
  PASS Improved commit message quality (+3 points)
  PASS Cleaned up stale branches (+2 points)

Recommendations:
  • Continue current improvement trajectory
  • Focus on documentation completeness
  • Consider adding more contributors
```

### CI/CD Integration

Perfect for continuous integration pipelines:

```yaml
# GitHub Actions example
- name: Run GPHC Health Check
  run: |
    gphc check --format json --output health-report.json
    gphc trend --days 7 --format json --output trend-report.json

- name: Upload Health Reports
  uses: actions/upload-artifact@v3
  with:
    name: health-reports
    path: |
      health-report.json
      trend-report.json
      .gphc-history.json
```

### Team Benefits

- **Progress Tracking**: See how your team's efforts improve code quality
- **Goal Setting**: Set health score targets and track progress
- **Quality Metrics**: Monitor technical debt reduction over time
- **Team Motivation**: Visualize improvements and celebrate progress
- **CI Insights**: Integrate health trends into your deployment pipeline

### Integration with Health Check

Historical tracking is automatically included in the main health check:

```bash
gphc check
```

The trend data will be referenced in recommendations and next steps when available.

## Multi-Repository Scan

GPHC can analyze multiple repositories simultaneously, making it perfect for organizations with many projects.

### Basic Usage

Scan multiple repositories in a directory:

```bash
# Scan all repositories in a directory
gphc scan ~/projects

# Recursive scan (find all Git repos in subdirectories)
gphc scan ~/projects --recursive

# Scan specific repositories
gphc scan ~/project-a ~/project-b ~/project-c

# Scan with custom output format
gphc scan ~/projects --format json
gphc scan ~/projects --format yaml
```

### Example Output

```bash
$ gphc scan ~/projects --recursive

Multi-Repository Health Scan Results
====================================

project-a: 92/100 (A-)
project-b: 78/100 (C+)
project-c: 85/100 (B+)
project-d: 67/100 (D+)
project-e: 91/100 (A-)

Summary:
  Total Repositories: 5
  Average Health: 82.6/100
  Highest Score: project-a (92/100)
  Lowest Score: project-d (67/100)
  Health Distribution:
    A Grade (90-100): 2 repositories
    B Grade (80-89): 1 repository
    C Grade (70-79): 1 repository
    D Grade (60-69): 1 repository

Recommendations:
  • Focus on project-d for immediate improvement
  • Share best practices from project-a and project-e
  • Consider standardizing documentation across all projects
```

### Advanced Options

```bash
# Filter by minimum score
gphc scan ~/projects --min-score 80

# Exclude specific directories
gphc scan ~/projects --exclude "*/node_modules" --exclude "*/vendor"

# Include only specific file types
gphc scan ~/projects --include "*.go" --include "*.js"

# Parallel processing (faster for many repos)
gphc scan ~/projects --parallel 4

# Generate detailed report
gphc scan ~/projects --detailed --output scan-report.json
```

### Organization Benefits

- **Portfolio Overview**: Get a bird's-eye view of all your projects
- **Quality Benchmarking**: Compare projects against each other
- **Resource Allocation**: Identify which projects need attention
- **Best Practice Sharing**: Find your best-performing projects
- **Compliance Monitoring**: Ensure all projects meet quality standards
- **Team Productivity**: Track improvements across the entire organization

### CI/CD Integration

Perfect for automated organization-wide health monitoring:

```yaml
# GitHub Actions example
- name: Multi-Repository Health Scan
  run: |
    gphc scan ./repos --recursive --format json --output org-health.json
    
- name: Upload Organization Health Report
  uses: actions/upload-artifact@v3
  with:
    name: organization-health-report
    path: org-health.json
```

### Use Cases

- **Software Companies**: Monitor all client projects
- **Open Source Organizations**: Track health of multiple repositories
- **Enterprise Teams**: Ensure consistency across departments
- **Consulting Firms**: Maintain quality across client projects
- **Educational Institutions**: Monitor student project portfolios

## CI/CD Integration

GPHC integrates seamlessly with CI/CD pipelines to ensure code quality standards are maintained throughout the development process.

### Quality Gates

Set minimum health score requirements to prevent low-quality code from being merged:

```bash
# Fail if health score is below 85
gphc check --min-score 85

# Fail if any critical checks fail
gphc check --fail-on-critical

# Fail if health score drops below threshold
gphc check --min-score 80 --fail-on-decline
```

### GitHub Actions Integration

**Basic Health Check:**
```yaml
name: Health Check
on: [push, pull_request]

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
        
      - name: Run Health Check
        run: gphc check --min-score 85 --format json
```

**Advanced Pipeline with Quality Gates:**
```yaml
name: Quality Pipeline
on: [push, pull_request]

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
        
      - name: Run Health Check
        id: health
        run: |
          gphc check --min-score 85 --format json --output health-report.json
          echo "score=$(jq -r '.summary.score' health-report.json)" >> $GITHUB_OUTPUT
          
      - name: Upload Health Report
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: health-report
          path: health-report.json
          
      - name: Comment PR with Health Score
        if: github.event_name == 'pull_request'
        uses: actions/github-script@v6
        with:
          script: |
            const score = '${{ steps.health.outputs.score }}';
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `Health Score: ${score}/100`
            });
```

### GitLab CI Integration

```yaml
# .gitlab-ci.yml
stages:
  - health-check

health-check:
  stage: health-check
  image: golang:1.23
  before_script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
  script:
    - gphc check --min-score 85 --format json --output health-report.json
  artifacts:
    reports:
      junit: health-report.json
    paths:
      - health-report.json
  only:
    - merge_requests
    - main
```

### Jenkins Integration

```groovy
pipeline {
    agent any
    
    stages {
        stage('Health Check') {
            steps {
                sh 'go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest'
                sh 'gphc check --min-score 85 --format json --output health-report.json'
            }
            post {
                always {
                    archiveArtifacts artifacts: 'health-report.json'
                    publishHTML([
                        allowMissing: false,
                        alwaysLinkToLastBuild: true,
                        keepAll: true,
                        reportDir: '.',
                        reportFiles: 'health-report.json',
                        reportName: 'Health Report'
                    ])
                }
            }
        }
    }
}
```

### Pre-commit Integration

**pre-commit configuration:**
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
        always_run: true
```

**Husky integration:**
```json
{
  "husky": {
    "hooks": {
      "pre-commit": "gphc pre-commit"
    }
  }
}
```

### Quality Thresholds

Configure different thresholds for different environments:

```bash
# Development branch - lower threshold
gphc check --min-score 70

# Staging branch - medium threshold  
gphc check --min-score 80

# Production branch - high threshold
gphc check --min-score 90
```

### Exit Codes

GPHC provides specific exit codes for CI/CD integration:

- **0**: All checks passed, health score above threshold
- **1**: Health score below minimum threshold
- **2**: Critical checks failed
- **3**: Configuration or system error

### Benefits for Teams

- **Quality Assurance**: Prevent low-quality code from reaching production
- **Consistent Standards**: Ensure all team members follow quality guidelines
- **Early Detection**: Catch quality issues before they become problems
- **Automated Enforcement**: No manual review needed for basic quality checks
- **Team Accountability**: Everyone contributes to maintaining code quality
- **Continuous Improvement**: Track quality trends over time

### Advanced CI/CD Features

```bash
# Generate badges for README
gphc badge --min-score 85

# Multi-repository health monitoring
gphc scan ./repos --min-score 80 --format json

# Historical trend analysis in CI
gphc trend --days 7 --format json
```

## Custom Rules (Rule Engine)

GPHC allows you to define custom health checks tailored to your project's specific requirements and organizational policies.

### Basic Custom Rules

Define custom checks in your `gphc.yml` configuration file:

```yaml
# Custom health checks
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5
    description: "Project must have a SECURITY.md file for vulnerability reporting"
    
  - id: CUSTOM-901
    name: "Has CONTRIBUTING.md"
    path: "CONTRIBUTING.md"
    score: 3
    description: "Project should have contribution guidelines"
    
  - id: CUSTOM-902
    name: "Has CODE_OF_CONDUCT.md"
    path: "CODE_OF_CONDUCT.md"
    score: 2
    description: "Project should have a code of conduct"
```

### Regex-Based Rules

Search for patterns in files using regular expressions:

```yaml
custom_checks:
  - id: CUSTOM-903
    name: "No TODO Comments"
    pattern: "TODO|FIXME|HACK"
    file_pattern: "*.go"
    score: -2
    description: "Code should not contain TODO comments"
    
  - id: CUSTOM-904
    name: "Has License Header"
    pattern: "Copyright.*\\d{4}"
    file_pattern: "*.go"
    score: 3
    description: "Source files should have copyright headers"
    
  - id: CUSTOM-905
    name: "No Hardcoded Secrets"
    pattern: "(password|secret|key)\\s*=\\s*['\"][^'\"]+['\"]"
    file_pattern: "*.{go,js,py}"
    score: -5
    description: "No hardcoded secrets should be present in code"
```

### Advanced Custom Rules

More sophisticated rule configurations:

```yaml
custom_checks:
  - id: CUSTOM-906
    name: "API Documentation Complete"
    type: "file_content"
    path: "docs/api.md"
    pattern: "GET|POST|PUT|DELETE"
    min_matches: 4
    score: 5
    description: "API documentation should cover all HTTP methods"
    
  - id: CUSTOM-907
    name: "Test Coverage Threshold"
    type: "file_content"
    path: "coverage.txt"
    pattern: "coverage: (\\d+\\.\\d+)%"
    min_value: 80.0
    score: 5
    description: "Test coverage should be at least 80%"
    
  - id: CUSTOM-908
    name: "No Deprecated Dependencies"
    type: "file_content"
    path: "go.mod"
    pattern: "deprecated|obsolete"
    score: -3
    description: "No deprecated dependencies should be used"
```

### Rule Types

GPHC supports several types of custom rules:

#### File Existence Rules
```yaml
- id: CUSTOM-910
  name: "Required File Check"
  type: "file_exists"
  path: "required-file.txt"
  score: 5
```

#### Content Pattern Rules
```yaml
- id: CUSTOM-911
  name: "Pattern Search"
  type: "file_content"
  pattern: "your-regex-pattern"
  file_pattern: "*.{go,js,py}"
  score: 3
```

#### Directory Structure Rules
```yaml
- id: CUSTOM-912
  name: "Directory Structure"
  type: "directory_exists"
  path: "src/main"
  score: 2
```

#### File Size Rules
```yaml
- id: CUSTOM-913
  name: "File Size Check"
  type: "file_size"
  path: "large-file.txt"
  max_size: "10MB"
  score: -2
```

### Configuration Examples

**Security-Focused Project:**
```yaml
custom_checks:
  - id: SEC-001
    name: "Security Policy"
    path: "SECURITY.md"
    score: 10
    
  - id: SEC-002
    name: "No Hardcoded Secrets"
    pattern: "(api_key|secret|password)\\s*=\\s*['\"][^'\"]+['\"]"
    file_pattern: "*.{go,js,py}"
    score: -10
    
  - id: SEC-003
    name: "Dependency Security Scan"
    path: "security-scan.txt"
    pattern: "vulnerabilities found: 0"
    score: 5
```

**Open Source Project:**
```yaml
custom_checks:
  - id: OSS-001
    name: "License File"
    path: "LICENSE"
    score: 8
    
  - id: OSS-002
    name: "Contributing Guidelines"
    path: "CONTRIBUTING.md"
    score: 5
    
  - id: OSS-003
    name: "Code of Conduct"
    path: "CODE_OF_CONDUCT.md"
    score: 3
    
  - id: OSS-004
    name: "Issue Templates"
    path: ".github/ISSUE_TEMPLATE"
    score: 2
```

**Enterprise Project:**
```yaml
custom_checks:
  - id: ENT-001
    name: "Architecture Documentation"
    path: "docs/architecture.md"
    score: 8
    
  - id: ENT-002
    name: "Deployment Guide"
    path: "docs/deployment.md"
    score: 5
    
  - id: ENT-003
    name: "Monitoring Configuration"
    path: "monitoring.yml"
    score: 3
    
  - id: ENT-004
    name: "No Console Logs"
    pattern: "console\\.log|print\\(|fmt\\.Print"
    file_pattern: "*.{go,js,py}"
    score: -2
```

### Rule Execution

Custom rules are executed alongside built-in checks:

```bash
# Run health check with custom rules
gphc check

# Validate custom rules configuration
gphc validate-config

# Test specific custom rule
gphc test-rule CUSTOM-900
```

### Example Output

```bash
$ gphc check

Custom Rules Results:
====================

CUSTOM-900: Has SECURITY.md - PASS (5 points)
CUSTOM-901: Has CONTRIBUTING.md - PASS (3 points)
CUSTOM-902: Has CODE_OF_CONDUCT.md - FAIL (0 points)
CUSTOM-903: No TODO Comments - FAIL (-2 points)

Custom Rules Summary:
  Total Rules: 4
  Passed: 2
  Failed: 2
  Score: 6/10
```

### Benefits

- **Project-Specific Requirements**: Enforce policies unique to your project
- **Organizational Standards**: Maintain consistency across multiple repositories
- **Compliance**: Ensure adherence to industry standards and regulations
- **Quality Gates**: Prevent specific issues from reaching production
- **Team Guidelines**: Enforce coding standards and best practices
- **Flexible Scoring**: Customize the impact of each rule on overall health score

### Best Practices

1. **Start Simple**: Begin with basic file existence checks
2. **Use Descriptive Names**: Make rule names clear and actionable
3. **Appropriate Scoring**: Balance rule scores with built-in checks
4. **Document Rules**: Always include descriptions for team understanding
5. **Test Rules**: Validate rules before deploying to production
6. **Regular Review**: Update rules as project requirements evolve

## Slack / Webhook Notifications

GPHC can automatically send health reports to your team's communication channels, keeping everyone informed about repository health status in real-time.

### Basic Setup

Configure webhook notifications in your `gphc.yml`:

```yaml
# Notification settings
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    channel: "#dev-team"
    enabled: true
    
  discord:
    webhook_url: "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
    enabled: false
    
  custom_webhook:
    url: "https://your-webhook-endpoint.com/notify"
    enabled: false
```

### Slack Integration

**Basic Slack Notification:**
```bash
# Run health check and send to Slack
gphc check --notify slack

# Send specific report to Slack
gphc check --format json --notify slack --webhook-url "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
```

**Example Slack Message:**
```
GPHC Health Report: project-name
Score: 83/100 (B+)
Status: PASS

Issues Found:
• Missing CONTRIBUTING.md
• 2 stale branches found
• Low test coverage (65%)

Recommendations:
• Add contribution guidelines
• Clean up old branches
• Increase test coverage

Full Report: https://github.com/owner/repo/actions/runs/123456
```

### Discord Integration

**Discord Webhook Setup:**
```bash
# Send to Discord
gphc check --notify discord --webhook-url "https://discord.com/api/webhooks/YOUR/DISCORD/WEBHOOK"
```

**Example Discord Message:**
```
**GPHC Health Report: project-name**

**Score:** 83/100 (B+)
**Status:** PASS

**Issues Found:**
• Missing CONTRIBUTING.md
• 2 stale branches found
• Low test coverage (65%)

**Recommendations:**
• Add contribution guidelines
• Clean up old branches
• Increase test coverage

[View Full Report](https://github.com/owner/repo/actions/runs/123456)
```

### Custom Webhook Integration

**Generic Webhook Support:**
```bash
# Send to any webhook endpoint
gphc check --notify webhook --webhook-url "https://your-service.com/webhook"
```

**Webhook Payload Format:**
```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "repository": "owner/repository",
  "score": 83,
  "grade": "B+",
  "status": "PASS",
  "checks_passed": 12,
  "checks_failed": 2,
  "checks_warning": 3,
  "issues": [
    {
      "id": "DOC-101",
      "name": "Missing CONTRIBUTING.md",
      "status": "FAIL",
      "score": 0
    }
  ],
  "recommendations": [
    "Add contribution guidelines",
    "Clean up old branches",
    "Increase test coverage"
  ],
  "trend": {
    "previous_score": 78,
    "change": "+5",
    "direction": "improving"
  }
}
```

### Advanced Notification Features

**Conditional Notifications:**
```yaml
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    conditions:
      - score_below: 80
      - critical_failures: true
      - score_dropped: 10
    enabled: true
```

**Custom Message Templates:**
```yaml
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"
    template: |
      🚀 *GPHC Report: {{.Repository}}*
      Score: {{.Score}}/100 ({{.Grade}})
      Status: {{.Status}}
      
      {{if .Issues}}
      Issues Found:
      {{range .Issues}}
      • {{.Name}}
      {{end}}
      {{end}}
      
      {{if .Recommendations}}
      Recommendations:
      {{range .Recommendations}}
      • {{.}}
      {{end}}
      {{end}}
```

### CI/CD Integration

**GitHub Actions with Slack Notifications:**
```yaml
name: Health Check with Notifications
on: [push, pull_request]

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
        
      - name: Run Health Check with Slack Notification
        env:
          SLACK_WEBHOOK_URL: ${{ secrets.SLACK_WEBHOOK_URL }}
        run: |
          gphc check --min-score 85 --notify slack --webhook-url $SLACK_WEBHOOK_URL
```

**GitLab CI with Discord Notifications:**
```yaml
# .gitlab-ci.yml
health-check:
  stage: health-check
  image: golang:1.23
  before_script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
  script:
    - gphc check --min-score 85 --notify discord --webhook-url $DISCORD_WEBHOOK_URL
  only:
    - merge_requests
    - main
```

### Notification Triggers

**Automatic Triggers:**
```bash
# Notify on health score below threshold
gphc check --notify slack --min-score 80

# Notify on critical failures
gphc check --notify slack --fail-on-critical

# Notify on score decline
gphc check --notify slack --min-score 85 --fail-on-decline
```

**Manual Triggers:**
```bash
# Send current health status
gphc notify slack

# Send trend analysis
gphc trend --days 7 --notify slack

# Send multi-repository scan results
gphc scan ~/projects --notify slack
```

### Team Benefits

- **Real-time Awareness**: Team members get immediate updates about repository health
- **Proactive Communication**: Issues are communicated before they become problems
- **Team Accountability**: Everyone stays informed about code quality standards
- **Progress Tracking**: Celebrate improvements and track quality trends together
- **Centralized Updates**: All health information in one communication channel
- **Integration**: Seamlessly fits into existing team workflows

### Configuration Examples

**Development Team Setup:**
```yaml
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/TEAM/DEV/WEBHOOK"
    channel: "#development"
    conditions:
      - score_below: 75
      - critical_failures: true
    enabled: true
```

**DevOps Team Setup:**
```yaml
notifications:
  slack:
    webhook_url: "https://hooks.slack.com/services/TEAM/DEVOPS/WEBHOOK"
    channel: "#devops-alerts"
    conditions:
      - score_below: 90
      - score_dropped: 5
    enabled: true
    
  discord:
    webhook_url: "https://discord.com/api/webhooks/TEAM/DEVOPS/WEBHOOK"
    channel: "devops-alerts"
    conditions:
      - critical_failures: true
    enabled: true
```

**Enterprise Setup:**
```yaml
notifications:
  custom_webhook:
    url: "https://enterprise-monitoring.com/gphc-webhook"
    headers:
      Authorization: "Bearer $GPHC_TOKEN"
      Content-Type: "application/json"
    conditions:
      - score_below: 85
      - critical_failures: true
      - score_dropped: 10
    enabled: true
```

## Semantic Commit Verification

GPHC analyzes commit messages against actual code changes to ensure semantic consistency and prevent misleading commit history.

### How It Works

GPHC compares commit message prefixes with the actual changes made in the commit:

```bash
# Run semantic commit verification
gphc check --semantic-verification

# Analyze specific commits
gphc analyze-commits --semantic-check

# Verify commit message consistency
gphc verify-semantic --commits 10
```

### Detection Examples

**Mismatch Detection:**

```bash
Commit: fix: resolve authentication issue
Changes: +800 lines added, -50 lines deleted
Files: src/auth.go, src/user.go, src/session.go
Analysis: MAJOR MISMATCH - 'fix' commit added 800 lines
Warning: This appears to be a feature addition, not a bug fix
```

**Correct Semantic Usage:**

```bash
Commit: feat: add user authentication system
Changes: +800 lines added, -0 lines deleted
Files: src/auth.go, src/user.go, src/session.go
Analysis: SEMANTIC MATCH - 'feat' commit with significant additions
Status: PASS
```

### Semantic Commit Types

GPHC recognizes standard semantic commit prefixes:

#### Primary Types
- **feat**: New features or functionality
- **fix**: Bug fixes and corrections
- **docs**: Documentation changes only
- **style**: Code formatting, no logic changes
- **refactor**: Code restructuring without behavior changes
- **test**: Adding or modifying tests
- **chore**: Maintenance tasks, dependencies, build changes

#### Extended Types
- **perf**: Performance improvements
- **ci**: Continuous integration changes
- **build**: Build system changes
- **revert**: Reverting previous commits

### Analysis Rules

**Feature Commits (feat:):**
```yaml
semantic_rules:
  feat:
    min_lines_added: 10
    max_lines_deleted: 5
    allowed_file_types: ["*.go", "*.js", "*.py", "*.ts"]
    forbidden_patterns: ["test", "spec", "mock"]
    description: "New functionality should add significant code"
```

**Fix Commits (fix:):**
```yaml
semantic_rules:
  fix:
    max_lines_added: 100
    max_lines_deleted: 50
    allowed_file_types: ["*.go", "*.js", "*.py", "*.ts"]
    required_patterns: ["bug", "error", "issue"]
    description: "Bug fixes should be minimal changes"
```

**Documentation Commits (docs:):**
```yaml
semantic_rules:
  docs:
    max_lines_added: 200
    max_lines_deleted: 100
    allowed_file_types: ["*.md", "*.rst", "*.txt", "*.yml", "*.yaml"]
    forbidden_patterns: ["src/", "lib/", "app/", "*.go", "*.js"]
    description: "Documentation changes should not modify code"
```

**Style Commits (style:):**
```yaml
semantic_rules:
  style:
    max_lines_added: 50
    max_lines_deleted: 50
    allowed_file_types: ["*.go", "*.js", "*.py", "*.ts", "*.css"]
    description: "Style changes should be formatting only"
```

### Configuration

Configure semantic verification in your `gphc.yml`:

```yaml
# Semantic commit verification settings
semantic_verification:
  enabled: true
  strict_mode: false
  check_last_commits: 20
  
  # Custom rules for commit types
  rules:
    feat:
      min_lines_added: 10
      max_lines_added: 1000
      max_lines_deleted: 20
      allowed_file_types: ["*.go", "*.js", "*.py", "*.ts"]
      
    fix:
      max_lines_added: 100
      max_lines_deleted: 50
      allowed_file_types: ["*.go", "*.js", "*.py", "*.ts"]
      
    docs:
      max_lines_added: 200
      max_lines_deleted: 100
      allowed_file_types: ["*.md", "*.rst", "*.txt", "*.yml"]
      forbidden_patterns: ["src/", "lib/", "*.go", "*.js"]
      
    style:
      max_lines_added: 50
      max_lines_deleted: 50
      allowed_file_types: ["*.go", "*.js", "*.py", "*.ts", "*.css"]
      
    test:
      max_lines_added: 300
      max_lines_deleted: 100
      allowed_file_types: ["*test*", "*spec*", "*mock*", "*.go", "*.js"]
      
    refactor:
      max_lines_added: 200
      max_lines_deleted: 200
      allowed_file_types: ["*.go", "*.js", "*.py", "*.ts"]
      
    chore:
      max_lines_added: 50
      max_lines_deleted: 50
      allowed_file_types: ["*.yml", "*.yaml", "*.json", "*.lock", "*.toml"]
```

### Example Output

```bash
$ gphc check --semantic-verification

Semantic Commit Verification Results:
=====================================

Commit: abc1234 - feat: add user authentication
Changes: +450 lines, -12 lines
Files: src/auth.go, src/user.go, src/session.go
Status: PASS - Semantic match ✓

Commit: def5678 - fix: resolve login bug
Changes: +800 lines, -50 lines  
Files: src/auth.go, src/user.go, src/session.go, src/middleware.go
Status: FAIL - Major mismatch ⚠️
Warning: 'fix' commit added 800 lines - appears to be a feature

Commit: ghi9012 - docs: update API documentation
Changes: +25 lines, -5 lines
Files: docs/api.md, README.md
Status: PASS - Semantic match ✓

Commit: jkl3456 - docs: update user interface
Changes: +200 lines, -50 lines
Files: src/components/UserForm.js, src/styles/user.css
Status: FAIL - Major mismatch ⚠️
Warning: 'docs' commit modified code files

Summary:
  Total Commits Analyzed: 4
  Semantic Matches: 2
  Mismatches Detected: 2
  Semantic Score: 50/100
```

### Advanced Analysis

**File Type Analysis:**
```bash
# Analyze by file type patterns
gphc verify-semantic --file-patterns "*.go,*.js" --commits 10

# Check specific commit types
gphc verify-semantic --types "feat,fix" --commits 20

# Strict mode verification
gphc verify-semantic --strict --commits 15
```

**Custom Rule Testing:**
```bash
# Test custom semantic rules
gphc test-semantic-rules --config custom-semantic.yml

# Validate commit message format
gphc validate-commit-format --commits 5
```

### Integration with CI/CD

**GitHub Actions Integration:**
```yaml
name: Semantic Commit Verification
on: [push, pull_request]

jobs:
  semantic-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 20
          
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
        
      - name: Verify Semantic Commits
        run: gphc verify-semantic --commits 10 --fail-on-mismatch
```

**Pre-commit Hook:**
```bash
#!/bin/bash
# .git/hooks/commit-msg

# Run semantic verification on commit message
gphc verify-commit-message --message "$1" --staged-files
if [ $? -ne 0 ]; then
    echo "Semantic commit verification failed"
    exit 1
fi
```

### Benefits

- **Commit History Accuracy**: Ensure commit messages accurately describe changes
- **Team Consistency**: Maintain consistent semantic commit practices across team
- **Code Review Quality**: Help reviewers understand changes from commit messages
- **Release Notes**: Generate accurate release notes from semantic commit history
- **Debugging**: Make it easier to find specific changes in git history
- **Automation**: Enable automated versioning and changelog generation

### Best Practices

1. **Use Descriptive Messages**: Commit messages should clearly describe the change
2. **Follow Semantic Format**: Use standard prefixes (feat:, fix:, docs:, etc.)
3. **Keep Commits Focused**: Each commit should address one logical change
4. **Review Before Committing**: Check that message matches actual changes
5. **Use Conventional Commits**: Follow the conventional commits specification
6. **Regular Verification**: Run semantic verification regularly in CI/CD

### Troubleshooting

**Common Issues:**

```bash
# Commit message doesn't match changes
Warning: 'fix' commit added 500 lines - consider using 'feat' instead

# Documentation commit modified code
Error: 'docs' commit modified .go files - use 'feat' or 'fix' instead

# Style commit changed logic
Warning: 'style' commit modified function logic - use 'refactor' instead
```

**Configuration Issues:**
```bash
# Invalid semantic rules
Error: semantic_rules.feat.max_lines_added must be positive

# Missing file type patterns
Warning: No file type patterns defined for 'docs' commits
```

## Interactive Terminal UI (TUI)

GPHC provides a beautiful, interactive terminal user interface that makes health checking and monitoring an engaging experience for developers and technical teams.

### Getting Started

Launch the interactive terminal interface:

```bash
# Start the TUI
gphc tui

# Start TUI with specific repository
gphc tui --path /path/to/repository

# Start TUI with auto-refresh
gphc tui --refresh 30s

# Start TUI in full-screen mode
gphc tui --fullscreen
```

### Interface Overview

The TUI provides multiple views and interactive features:

```
┌─────────────────────────────────────────────────────────────┐
│ GPHC - Git Project Health Checker                          │
│ Repository: /path/to/project                                │
│ Last Updated: 2024-01-15 10:30:00                          │
├─────────────────────────────────────────────────────────────┤
│ Overall Health Score: 85/100 (B+)                          │
│ Status: PASS                                                │
│                                                             │
│ ┌─ Health Overview ──────────────────────────────────────┐ │
│ │ Documentation & Project Structure: 90/100 (A-)         │ │
│ │ Commit History Quality: 85/100 (B+)                    │ │
│ │ Git Cleanup & Hygiene: 80/100 (B-)                    │ │
│ │ Codebase Structure: 75/100 (C+)                        │ │
│ └─────────────────────────────────────────────────────────┘ │
│                                                             │
│ ┌─ Quick Actions ────────────────────────────────────────┐ │
│ │ [F1] Help  [F2] Filter  [F3] Trends  [F4] Settings    │ │
│ │ [F5] Refresh  [F6] Export  [F7] Notify  [F8] Quit     │ │
│ └─────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Main Views

#### 1. Health Overview Dashboard
- **Overall Score Display**: Large, colorful health score with grade
- **Category Breakdown**: Visual representation of each health category
- **Status Indicators**: Pass/Fail/Warning status for each check
- **Quick Stats**: Number of checks passed, failed, and warnings

#### 2. Detailed Check Results
```
┌─ Check Results ────────────────────────────────────────────┐
│ DOC-101: README.md exists                                 │
│ Status: PASS  Score: 5/5  Category: Documentation        │
│ Details: Project has a comprehensive README.md file      │
│                                                           │
│ DOC-102: LICENSE file exists                              │
│ Status: FAIL  Score: 0/5  Category: Documentation        │
│ Details: Project is missing a LICENSE file               │
│ Recommendation: Add a LICENSE file for legal clarity     │
│                                                           │
│ COMM-201: Conventional commit format                      │
│ Status: WARN  Score: 3/5  Category: Commit Quality       │
│ Details: 15% of commits don't follow conventional format │
│ Recommendation: Use feat:, fix:, docs: prefixes          │
└───────────────────────────────────────────────────────────┘
```

#### 3. Interactive Filtering
- **Category Filter**: Filter by Documentation, Commits, Hygiene, Structure
- **Status Filter**: Show only PASS, FAIL, or WARNING checks
- **Score Range**: Filter by minimum/maximum score thresholds
- **Search**: Find specific checks by name or ID

#### 4. Trend Analysis View
```
┌─ Health Trends ────────────────────────────────────────────┐
│ Score History (Last 30 days)                              │
│                                                           │
│ 100 ┤                                                    │
│  90 ┤     ●                                              │
│  80 ┤   ●   ●                                            │
│  70 ┤ ●       ●                                          │
│  60 ┤           ●                                        │
│  50 ┤                                                     │
│     └─────────────────────────────────────────────────────│
│      Jan 01  Jan 08  Jan 15  Jan 22  Jan 29             │
│                                                           │
│ Current: 85/100  Previous: 78/100  Change: +7 (+9.0%)    │
│ Trend: Improving  Average: 81.2/100                      │
└───────────────────────────────────────────────────────────┘
```

### Interactive Features

#### Keyboard Shortcuts
```
Navigation:
  ↑/↓/←/→    Move cursor
  Tab        Switch between panels
  Enter      Select/expand item
  Esc        Go back/close dialog

Actions:
  F1         Help and shortcuts
  F2         Open filter menu
  F3         View trend analysis
  F4         Settings and configuration
  F5         Refresh data
  F6         Export results
  F7         Send notifications
  F8         Quit application

Search and Filter:
  /          Search checks
  Ctrl+F     Advanced filtering
  Ctrl+R     Reset filters
  Space      Toggle selection
```

#### Real-time Updates
```bash
# Auto-refresh every 30 seconds
gphc tui --refresh 30s

# Auto-refresh every 2 minutes
gphc tui --refresh 2m

# Watch mode - refresh on file changes
gphc tui --watch
```

#### Rule Explanations
When you select a check, view detailed explanations:

```
┌─ Rule Details: DOC-101 ────────────────────────────────────┐
│ Name: README.md exists                                     │
│ ID: DOC-101                                               │
│ Category: Documentation & Project Structure               │
│ Score: 5 points                                           │
│                                                           │
│ Description:                                              │
│ Checks if the project has a README.md file in the root   │
│ directory. This file is essential for project             │
│ documentation and helps new contributors understand        │
│ the project.                                              │
│                                                           │
│ Requirements:                                             │
│ • README.md file must exist in repository root           │
│ • File should not be empty                               │
│ • Should contain project description                      │
│                                                           │
│ Benefits:                                                 │
│ • Improves project discoverability                       │
│ • Helps new contributors get started                     │
│ • Provides essential project information                  │
│                                                           │
│ [Press Enter to go back]                                 │
└───────────────────────────────────────────────────────────┘
```

### Advanced Features

#### Multi-Repository View
```bash
# Compare multiple repositories
gphc tui --multi-repo ~/projects/*

# Scan and display all repositories
gphc tui --scan-recursive ~/projects
```

#### Custom Themes
```yaml
# gphc.yml
tui:
  theme: "dark"  # dark, light, auto
  colors:
    score_excellent: "#00ff00"  # Green
    score_good: "#ffff00"       # Yellow
    score_poor: "#ff0000"      # Red
    status_pass: "#00ff00"
    status_fail: "#ff0000"
    status_warn: "#ffaa00"
  
  layout:
    show_trends: true
    show_recommendations: true
    auto_refresh: "30s"
    fullscreen: false
```

#### Export from TUI
- **Current View Export**: Export currently displayed data
- **Filtered Export**: Export only filtered results
- **Trend Export**: Export historical trend data
- **Formats**: JSON, YAML, Markdown, HTML

### Configuration

Configure TUI behavior in your `gphc.yml`:

```yaml
# TUI settings
tui:
  enabled: true
  theme: "dark"
  auto_refresh: "30s"
  fullscreen: false
  
  # Display options
  show_trends: true
  show_recommendations: true
  show_details: true
  compact_mode: false
  
  # Color scheme
  colors:
    score_excellent: "#00ff00"
    score_good: "#ffff00"
    score_poor: "#ff0000"
    status_pass: "#00ff00"
    status_fail: "#ff0000"
    status_warn: "#ffaa00"
    background: "#000000"
    foreground: "#ffffff"
  
  # Layout preferences
  layout:
    main_panel_height: 60
    sidebar_width: 30
    show_help: true
    show_shortcuts: true
```

### Benefits

- **Enhanced User Experience**: Beautiful, intuitive interface for health monitoring
- **Real-time Monitoring**: Live updates and auto-refresh capabilities
- **Interactive Exploration**: Deep dive into specific checks and trends
- **Efficient Navigation**: Quick access to different views and features
- **Visual Feedback**: Color-coded scores and status indicators
- **Team Collaboration**: Easy sharing of health status with team members
- **Accessibility**: Keyboard-driven interface suitable for all users

### Use Cases

- **Daily Health Monitoring**: Quick daily check of repository health
- **Team Meetings**: Visual presentation of project health status
- **Code Review Preparation**: Understanding project health before reviews
- **Onboarding**: New team members exploring project structure
- **CI/CD Integration**: Visual monitoring of automated health checks
- **Client Presentations**: Professional health status displays

### Troubleshooting

**Common Issues:**

```bash
# Terminal not supported
Error: Terminal does not support TUI mode
Solution: Use a modern terminal emulator (iTerm2, Windows Terminal, etc.)

# Color display issues
Warning: Terminal colors not properly configured
Solution: Set TERM environment variable or use --no-color flag

# Performance issues
Warning: TUI performance degraded with large repositories
Solution: Use --compact-mode or reduce --refresh interval
```

**Performance Optimization:**
```bash
# Use compact mode for better performance
gphc tui --compact-mode

# Reduce refresh frequency
gphc tui --refresh 2m

# Disable auto-refresh for static analysis
gphc tui --no-refresh
```

## Configuration

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

## Architecture

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

## Development

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

## Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

**How to contribute:**
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature-name`)
3. Make your changes and test them
4. Commit your changes (`git commit -m 'Add your feature'`)
5. Push to your branch (`git push origin feature/your-feature-name`)
6. Open a Pull Request with a clear description

## Roadmap

### Phase 1: Core Features
- [x] Basic CLI framework
- [x] Git repository analysis
- [x] Documentation checks
- [x] Commit quality analysis
- [x] Branch hygiene checks
- [x] Scoring engine
- [x] Colorful terminal output

### Phase 2: Enhanced Features
- [x] GitHub API integration for branch protection
- [x] Pre-commit hook validation
- [x] Custom rule definitions
- [x] JSON/XML output formats
- [x] CI/CD integration

#### CI/CD Integration

**What it does:**
Allows CIs to automatically check project health.
If the score is below the allowed threshold, the pipeline fails.

**Sample GitHub Action:**
```yaml
- name: Health Check
  run: gphc check --min-score 85 --format json
```

**Why it's important:**
So that projects pass the minimum required quality before being merged.

#### Custom Rules (Rule Engine)

**What it does:**
Users can add rules to their `gphc.yml`, for example:

```yaml
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5
```

Or even regex-based rules for searching in files.

**Why it's important:**
For projects that have specific policies or requirements (e.g., SECURITY.md file or TODO detection).

#### Semantic Commit Verification

**What it does:**
Checks whether commit messages match the actual changes made.
For example:

```
commit with title "fix:" but added 800 lines of code → warning
commit with "docs:" but code was changed → warning
```

**Why it's important:**
To prevent misleading commits and maintain transparency in history.

### Phase 3: Team, Trends & Automation
- [x] Multi-repository analysis
- [x] Historical trend analysis
- [x] Team collaboration metrics
- [x] Integration with popular Git hosting platforms
- [x] Interactive Terminal UI (TUI)
- [x] Web Dashboard

#### Interactive Terminal UI (TUI)

**What it does:**
Creates a graphical interface in the terminal (like htop):

```bash
gphc tui
```

Colorful and interactive score display with ability to filter, view rule explanations, and browse score trends.

**Why it's important:**
Much better user experience for developers and technical teams.

#### Web Dashboard

**What it does:**
With the command:

```bash
gphc serve
```

Shows health results of projects on a local HTML dashboard (like localhost:8080).

**Features:**
- Browse multiple projects
- View trends
- Export reports

**Why it's important:**
Very useful for teams and internal monitoring.

#### Multi-Repository Scan

**What it does:**
Scans multiple local repositories simultaneously:

```bash
gphc scan ~/projects --recursive
```

**Output:**
```
project-a: 92
project-b: 78
project-c: 85
Average Health: 85.0
```

**Why it's important:**
Very useful for companies or organizations that have multiple repositories.

#### Slack / Webhook Notifications

**What it does:**
After execution, the health report is sent to Slack or Discord:

```
GPHC Report: 83/100
Missing CONTRIBUTING.md
2 stale branches found
```

**Why it's important:**
Teams can immediately be informed about the health status of repositories in their work channel.

#### Historical Health Tracking

**What it does:**
Every time GPHC runs, it saves the result to `.gphc-history.json`.
With the command:

```bash
gphc trend
```

It shows how the project score has changed over time (e.g., improvement from 72 → 85 in one month).

**Why it's important:**
Useful for teams and CI to see if quality has improved over time.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [go-git](https://github.com/go-git/go-git) for Git repository access
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [LipGloss](https://github.com/charmbracelet/lipgloss) for terminal styling
- [Viper](https://github.com/spf13/viper) for configuration management

---

**Made with love for the Open Source community**