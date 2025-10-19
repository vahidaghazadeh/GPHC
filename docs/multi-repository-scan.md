# Multi-Repository Scan Guide

This guide covers scanning multiple repositories simultaneously for comprehensive health analysis.

## Overview

The Multi-Repository Scan feature allows you to analyze multiple Git repositories at once, perfect for organizations with many projects or developers managing multiple repositories.

## Basic Usage

### Simple Scan
```bash
# Scan current directory
gphc scan

# Scan specific directory
gphc scan ~/projects

# Scan with recursive option
gphc scan ~/projects --recursive
```

### Example Output
```
Multi-Repository Health Scan Results
====================================

Scanning: /Users/dev/projects/project-a
  project-a: 92/100 (A-)

Scanning: /Users/dev/projects/project-b
  project-b: 78/100 (C+)

Scanning: /Users/dev/projects/project-c
  project-c: 85/100 (B+)

Summary:
  Total Repositories: 3
  Average Health: 85.0/100
  Highest Score: project-a (92/100)
  Lowest Score: project-b (78/100)
```

## Advanced Options

### Recursive Scanning
```bash
# Find all Git repositories in directory tree
gphc scan ~/projects --recursive

# Scan with minimum score threshold
gphc scan ~/projects --recursive --min-score 80

# Scan with parallel processing
gphc scan ~/projects --recursive --parallel 8
```

### Filtering Options
```bash
# Exclude specific directories
gphc scan ~/projects --exclude "node_modules" --exclude ".git"

# Include only specific patterns
gphc scan ~/projects --include "*.go" --include "*.py"

# Combine filters
gphc scan ~/projects --recursive --exclude "test" --min-score 70
```

### Output Options
```bash
# Generate detailed report
gphc scan ~/projects --detailed

# Save output to file
gphc scan ~/projects --output scan-results.json

# Export in different formats
gphc scan ~/projects --format json
gphc scan ~/projects --format yaml
gphc scan ~/projects --format markdown
```

## Configuration

### Scan Settings
```yaml
# gphc.yml
scan:
  recursive: true
  min_score: 70
  parallel_jobs: 4
  detailed_report: false
  
  # Exclude patterns
  exclude_patterns:
    - "node_modules"
    - ".git"
    - "test"
    - "tmp"
  
  # Include patterns
  include_patterns:
    - "*.go"
    - "*.py"
    - "*.js"
  
  # Output settings
  output:
    format: "terminal"  # terminal, json, yaml, markdown
    file: ""  # empty for stdout
    detailed: false
```

### Command Line Flags
```bash
# Basic options
gphc scan [path] [flags]

Flags:
  -r, --recursive          Recursively scan subdirectories
  -m, --min-score int      Minimum health score threshold
  -e, --exclude strings    Exclude directories matching patterns
  -i, --include strings    Include only files matching patterns
  -p, --parallel int       Number of parallel jobs (default 4)
  -d, --detailed          Generate detailed report
  -o, --output string     Output file path (default: stdout)
```

## Use Cases

### Organization-wide Analysis
```bash
# Scan all company repositories
gphc scan /opt/repositories --recursive --min-score 80

# Generate organization health report
gphc scan /opt/repositories --recursive --detailed --output org-health.json
```

### Team Project Monitoring
```bash
# Monitor team projects
gphc scan ~/team/projects --recursive --min-score 75

# Track team progress
gphc scan ~/team/projects --recursive --trends
```

### Personal Repository Management
```bash
# Check all personal projects
gphc scan ~/dev --recursive

# Find projects needing attention
gphc scan ~/dev --recursive --min-score 60
```

### CI/CD Integration
```bash
# Scan in CI pipeline
gphc scan /workspace --recursive --min-score 85 --format json

# Fail pipeline if average score is too low
gphc scan /workspace --recursive --min-score 80 || exit 1
```

## Performance Optimization

### Parallel Processing
```bash
# Use more parallel jobs for faster scanning
gphc scan ~/projects --parallel 8

# Adjust based on system resources
gphc scan ~/projects --parallel 16  # For powerful systems
gphc scan ~/projects --parallel 2   # For limited resources
```

### Caching
```bash
# Enable caching for repeated scans
gphc scan ~/projects --cache

# Cache with TTL
gphc scan ~/projects --cache --cache-ttl 300s
```

### Filtering for Performance
```bash
# Exclude large directories
gphc scan ~/projects --exclude "node_modules" --exclude "vendor"

# Include only relevant repositories
gphc scan ~/projects --include "*.go" --include "*.py"
```

## Output Formats

### Terminal Output
```
Multi-Repository Health Scan Results
====================================

Scanning: /path/to/project-a
  project-a: 92/100 (A-)

Scanning: /path/to/project-b
  project-b: 78/100 (C+)

Summary:
  Total Repositories: 2
  Average Health: 85.0/100
  Highest Score: project-a (92/100)
  Lowest Score: project-b (78/100)
```

### JSON Output
```json
{
  "scan_results": [
    {
      "name": "project-a",
      "path": "/path/to/project-a",
      "score": 92,
      "grade": "A-"
    },
    {
      "name": "project-b",
      "path": "/path/to/project-b",
      "score": 78,
      "grade": "C+"
    }
  ],
  "summary": {
    "total_repositories": 2,
    "average_health": 85.0,
    "highest_score": {
      "name": "project-a",
      "score": 92
    },
    "lowest_score": {
      "name": "project-b",
      "score": 78
    }
  }
}
```

### Markdown Output
```markdown
# Multi-Repository Health Scan Results

## Summary
- **Total Repositories**: 2
- **Average Health**: 85.0/100
- **Highest Score**: project-a (92/100)
- **Lowest Score**: project-b (78/100)

## Results

| Repository | Score | Grade | Status |
|------------|-------|-------|--------|
| project-a  | 92    | A-    | PASS   |
| project-b  | 78    | C+    | WARN   |
```

## Integration Examples

### GitHub Actions
```yaml
name: Multi-Repository Health Check
on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly

jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
      - name: Scan Repositories
        run: gphc scan . --recursive --min-score 80 --format json
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: health-scan-results
          path: health-scan-results.json
```

### GitLab CI
```yaml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - gphc scan . --recursive --min-score 80 --format json
  artifacts:
    reports:
      junit: health-scan-results.json
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Health Check') {
            steps {
                sh 'go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest'
                sh 'gphc scan . --recursive --min-score 80 --format json'
            }
        }
    }
    post {
        always {
            archiveArtifacts artifacts: 'health-scan-results.json'
        }
    }
}
```

## Troubleshooting

### Common Issues

#### No Repositories Found
```bash
# Check if path contains Git repositories
find ~/projects -name ".git" -type d

# Use recursive flag
gphc scan ~/projects --recursive
```

#### Permission Denied
```bash
# Check repository permissions
ls -la ~/projects

# Use sudo if necessary
sudo gphc scan /opt/repositories --recursive
```

#### Performance Issues
```bash
# Reduce parallel jobs
gphc scan ~/projects --parallel 2

# Exclude large directories
gphc scan ~/projects --exclude "node_modules" --exclude "vendor"
```

### Error Handling
```bash
# Continue on errors
gphc scan ~/projects --continue-on-error

# Verbose output for debugging
gphc scan ~/projects --verbose

# Check specific repositories
gphc scan ~/projects --check-specific project-a project-b
```

## Best Practices

### For Organizations
1. **Regular Scanning**: Schedule weekly scans of all repositories
2. **Quality Thresholds**: Set minimum health score requirements
3. **Trend Monitoring**: Track health improvements over time
4. **Team Accountability**: Share scan results with teams
5. **Continuous Improvement**: Use results to improve processes

### For Teams
1. **Team Standards**: Establish team-wide health standards
2. **Regular Reviews**: Review scan results in team meetings
3. **Improvement Plans**: Create action plans for low-scoring projects
4. **Knowledge Sharing**: Share best practices across teams
5. **Tool Integration**: Integrate scanning into team workflows

### For Individuals
1. **Personal Projects**: Regularly scan personal repositories
2. **Learning Tool**: Use scans to learn best practices
3. **Portfolio Maintenance**: Keep portfolio repositories healthy
4. **Skill Development**: Improve Git and project management skills
5. **Community Contribution**: Contribute to open source projects

## Next Steps

- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
- [Historical Tracking](historical-tracking.md) - Health trend analysis
- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
- [Terminal UI](terminal-ui.md) - Interactive terminal interface
