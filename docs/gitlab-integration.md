# GitLab Integration Guide

This guide covers integrating GPHC with GitLab for advanced repository analysis.

## Overview

GitLab integration allows GPHC to access GitLab API for advanced repository analysis, including merge request settings, CI/CD configuration, and project settings.

## Setup

### Authentication
```bash
# Set GitLab token
export GITLAB_TOKEN=your_gitlab_token

# Or use GPHC_TOKEN
export GPHC_TOKEN=your_gitlab_token

# Set GitLab URL (for self-hosted instances)
export GITLAB_URL=https://gitlab.example.com
```

### Configuration
```yaml
# gphc.yml
gitlab:
  enabled: true
  token: "${GITLAB_TOKEN}"
  base_url: "https://gitlab.com"  # or your GitLab instance URL
  
  # Project settings
  project:
    id: "your-project-id"
    path: "your-group/your-project"
```

## Features

### Merge Request Settings
```bash
# Check merge request settings
git hc check --gitlab

# Example output:
PASS [GL-101] Merge request settings configured
  Message: Merge requests require approval
  Details: At least 2 approvals required for merge

PASS [GL-102] Pipeline required
  Message: Merge requests require successful pipeline
  Details: CI/CD pipeline must pass before merge
```

### CI/CD Configuration Check
```bash
# Check GitLab CI configuration
git hc check --gitlab

# Example output:
PASS [GL-201] CI/CD configuration exists
  Message: GitLab CI configuration found
  Details: .gitlab-ci.yml contains CI configuration

WARN [GL-202] Missing security scanning
  Message: No security scanning in CI pipeline
  Details: Consider adding security scanning to CI workflow
```

### Project Settings Analysis
```bash
# Analyze project settings
git hc check --gitlab

# Example output:
PASS [GL-301] Project visibility configured
  Message: Project visibility is set to internal
  Details: Appropriate visibility for team project

WARN [GL-302] Missing project description
  Message: Project description is empty
  Details: Add project description for better discoverability
```

## API Endpoints

### Project Information
```bash
# Get project information
GET /projects/{project_id}

# Example response:
{
  "id": 123,
  "name": "project-name",
  "path": "group/project-name",
  "description": "Project description",
  "visibility": "internal",
  "merge_requests_enabled": true,
  "issues_enabled": true,
  "wiki_enabled": true
}
```

### Merge Request Settings
```bash
# Get merge request settings
GET /projects/{project_id}/merge_requests

# Example response:
{
  "merge_requests": [
    {
      "id": 456,
      "title": "Feature implementation",
      "state": "opened",
      "approvals_required": 2,
      "approvals_left": 1
    }
  ]
}
```

### Pipeline Information
```bash
# Get pipeline information
GET /projects/{project_id}/pipelines

# Example response:
{
  "pipelines": [
    {
      "id": 789,
      "status": "success",
      "ref": "main",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ]
}
```

## Integration Examples

### GitLab CI
```yaml
# .gitlab-ci.yml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - git hc check --gitlab
  variables:
    GITLAB_TOKEN: $GITLAB_TOKEN
```

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit
export GITLAB_TOKEN=your_token
export GITLAB_URL=https://gitlab.com
git hc pre-commit --gitlab
```

## Configuration Options

### GitLab Settings
```yaml
# gphc.yml
gitlab:
  enabled: true
  token: "${GITLAB_TOKEN}"
  base_url: "https://gitlab.com"
  
  # Project information
  project:
    id: "123"
    path: "group/project-name"
  
  # Check settings
  checks:
    merge_request_settings: true
    ci_cd_configuration: true
    project_settings: true
    contributor_activity: true
    security_settings: true
  
  # Thresholds
  thresholds:
    min_approvals: 2
    min_contributors: 3
    max_contributor_dominance: 70  # percentage
```

### Custom Rules
```yaml
# gphc.yml
gitlab:
  custom_rules:
    - id: "GL-CUSTOM-001"
      name: "Has Security Policy"
      check: "security_policy"
      score: 5
      
    - id: "GL-CUSTOM-002"
      name: "Has Issue Templates"
      check: "issue_templates"
      score: 3
      
    - id: "GL-CUSTOM-003"
      name: "Has MR Templates"
      check: "mr_templates"
      score: 3
```

## Troubleshooting

### Common Issues
- **Authentication Failed**: Check GitLab token permissions
- **Rate Limiting**: Implement rate limiting for API calls
- **Project Not Found**: Verify project ID or path
- **Permission Denied**: Check token scopes

### Debugging
```bash
# Test GitLab connection
git hc check --gitlab --test-connection

# Verbose GitLab output
git hc check --gitlab --verbose

# Check specific project
git hc check --gitlab --project 123
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
3. **Compliance**: Ensure compliance with GitLab terms
4. **Monitoring**: Monitor API usage across teams
5. **Documentation**: Document integration procedures

## Next Steps
- [GitHub Integration](github-integration.md) - GitHub API integration
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
- [Web Dashboard](web-dashboard.md) - Web server and team collaboration
