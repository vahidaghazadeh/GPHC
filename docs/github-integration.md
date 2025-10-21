# GitHub Integration Guide

This guide covers integrating GPHC with GitHub for advanced repository analysis.

## Overview

GitHub integration allows GPHC to access GitHub API for advanced repository analysis, including branch protection rules, required reviewers, and CI/CD configuration.

## Setup

### Authentication
```bash
# Set GitHub token
export GITHUB_TOKEN=your_github_token

# Or use GPHC_TOKEN
export GPHC_TOKEN=your_github_token
```

### Configuration
```yaml
# gphc.yml
github:
  enabled: true
  token: "${GITHUB_TOKEN}"
  base_url: "https://api.github.com"  # For GitHub Enterprise
  
  # Repository settings
  repository:
    owner: "your-username"
    name: "your-repository"
```

## Features

### Branch Protection Analysis
```bash
# Check branch protection rules
git hc check --github

# Example output:
PASS [GH-101] Branch protection enabled
  Message: Main branch has protection rules enabled
  Details: Requires pull request reviews, status checks, and up-to-date branches

PASS [GH-102] Required reviewers configured
  Message: Pull requests require at least 2 reviewers
  Details: Code review process is properly configured
```

### CI/CD Configuration Check
```bash
# Check GitHub Actions configuration
git hc check --github

# Example output:
PASS [GH-201] CI/CD configuration exists
  Message: GitHub Actions workflow found
  Details: .github/workflows/ci.yml contains CI configuration

WARN [GH-202] Missing security checks
  Message: No security scanning in CI pipeline
  Details: Consider adding security scanning to CI workflow
```

### Contributor Activity Analysis
```bash
# Analyze contributor activity
git hc check --github

# Example output:
PASS [GH-301] Active contributors
  Message: Repository has 5 active contributors
  Details: Contributors active in last 30 days: 5

WARN [GH-302] Single contributor dominance
  Message: One contributor has 80% of commits
  Details: Consider encouraging more team participation
```

## API Endpoints

### Branch Protection
```bash
# Check branch protection rules
GET /repos/{owner}/{repo}/branches/{branch}/protection

# Example response:
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["ci", "test"]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 2
  }
}
```

### Repository Settings
```bash
# Get repository information
GET /repos/{owner}/{repo}

# Example response:
{
  "name": "project-name",
  "full_name": "owner/project-name",
  "private": false,
  "has_issues": true,
  "has_projects": true,
  "has_wiki": true
}
```

### Workflow Runs
```bash
# Get workflow run information
GET /repos/{owner}/{repo}/actions/runs

# Example response:
{
  "workflow_runs": [
    {
      "id": 123456,
      "status": "completed",
      "conclusion": "success",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Integration Examples

### GitHub Actions
```yaml
name: Health Check with GitHub Integration
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

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
      - name: Health Check with GitHub Integration
        run: git hc check --github
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit
export GITHUB_TOKEN=your_token
git hc pre-commit --github
```

## Configuration Options

### GitHub Settings
```yaml
# gphc.yml
github:
  enabled: true
  token: "${GITHUB_TOKEN}"
  base_url: "https://api.github.com"
  
  # Repository information
  repository:
    owner: "your-username"
    name: "your-repository"
  
  # Check settings
  checks:
    branch_protection: true
    required_reviewers: true
    ci_cd_configuration: true
    contributor_activity: true
    security_settings: true
  
  # Thresholds
  thresholds:
    min_reviewers: 2
    min_contributors: 3
    max_contributor_dominance: 70  # percentage
```

### Custom Rules
```yaml
# gphc.yml
github:
  custom_rules:
    - id: "GH-CUSTOM-001"
      name: "Has Security Policy"
      check: "security_policy"
      score: 5
      
    - id: "GH-CUSTOM-002"
      name: "Has Issue Templates"
      check: "issue_templates"
      score: 3
      
    - id: "GH-CUSTOM-003"
      name: "Has PR Templates"
      check: "pr_templates"
      score: 3
```

## Troubleshooting

### Common Issues
- **Authentication Failed**: Check GitHub token permissions
- **Rate Limiting**: Implement rate limiting for API calls
- **Repository Not Found**: Verify repository owner and name
- **Permission Denied**: Check token scopes

### Debugging
```bash
# Test GitHub connection
git hc check --github --test-connection

# Verbose GitHub output
git hc check --github --verbose

# Check specific repository
git hc check --github --repo owner/repository
```

## Best Practices

### For Teams
1. **Token Security**: Use secure token storage
2. **Permission Management**: Use minimal required permissions
3. **Rate Limiting**: Implement appropriate rate limiting
4. **Error Handling**: Handle API errors gracefully
5. **Monitoring**: Monitor API usage and limits

### For Organizations
1. **Centralized Tokens**: Use organization-level tokens
2. **Security Policies**: Implement security policies
3. **Compliance**: Ensure compliance with GitHub terms
4. **Monitoring**: Monitor API usage across teams
5. **Documentation**: Document integration procedures

## Next Steps
- [GitLab Integration](gitlab-integration.md) - GitLab API integration
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
