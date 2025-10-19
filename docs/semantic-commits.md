# Semantic Commits Guide

This guide covers semantic commit verification and analysis.

## Overview

Semantic commit verification ensures that commit messages accurately reflect the actual changes made in the code, maintaining transparency and consistency in project history.

## How It Works

### Commit Message Analysis
GPHC analyzes commit messages to determine the intended change type:
- **feat**: New features
- **fix**: Bug fixes
- **docs**: Documentation changes
- **style**: Code style changes
- **refactor**: Code refactoring
- **test**: Test additions/changes
- **chore**: Maintenance tasks

### Change Detection
GPHC analyzes actual file changes to determine the real change type:
- **File additions**: New files added
- **File modifications**: Existing files changed
- **File deletions**: Files removed
- **Line changes**: Lines added/removed

### Mismatch Detection
GPHC identifies when commit messages don't match actual changes:

```
WARN [SEM-101] Commit message mismatch detected
  Message: "fix: resolve authentication issue"
  Actual Changes: Added 800 lines of new code
  Expected: Bug fix (small changes)
  Actual: Feature addition (large changes)
  Recommendation: Use "feat:" prefix for new features
```

## Configuration

### Semantic Rules
```yaml
# gphc.yml
semantic_commits:
  enabled: true
  
  # Commit message validation
  message_validation:
    max_subject_length: 72
    require_imperative_mood: true
    require_conventional_format: true
  
  # Change analysis
  change_analysis:
    max_lines_added: 100
    max_lines_deleted: 100
    max_files_changed: 10
    
  # Mismatch detection
  mismatch_detection:
    enabled: true
    threshold_lines: 50
    threshold_files: 5
```

### Custom Rules
```yaml
# gphc.yml
semantic_commits:
  custom_rules:
    - type: "feat"
      max_lines_added: 200
      max_files_changed: 20
      
    - type: "fix"
      max_lines_added: 50
      max_files_changed: 5
      
    - type: "docs"
      allowed_file_types: [".md", ".rst", ".txt"]
      max_code_changes: 10
```

## Examples

### Valid Commits
```
feat: add user authentication system
- Added login/logout functionality
- Implemented JWT token handling
- Added user session management

fix: resolve memory leak in data processing
- Fixed memory allocation issue
- Added proper cleanup in data handlers
- Reduced memory usage by 30%

docs: update API documentation
- Added endpoint descriptions
- Updated request/response examples
- Fixed typos in documentation
```

### Invalid Commits
```
fix: add new user authentication system
# ❌ This is a feature addition, not a bug fix

feat: resolve authentication issue
# ❌ This is a bug fix, not a new feature

docs: implement new database schema
# ❌ This is code changes, not documentation
```

## Integration

### Pre-commit Hook
```bash
# Add to .git/hooks/pre-commit
#!/bin/sh
gphc pre-commit --semantic-check
```

### GitHub Actions
```yaml
name: Semantic Commit Check
on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  semantic-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
      - name: Semantic Commit Check
        run: gphc check --semantic-commits
```

## Best Practices

### Commit Message Format
```
<type>(<scope>): <subject>

<body>

<footer>
```

### Examples
```
feat(auth): add OAuth2 integration

Implement OAuth2 authentication flow with Google and GitHub providers.
Adds support for social login and user profile management.

Closes #123
```

```
fix(api): resolve rate limiting issue

Fix incorrect rate limiting calculation that was causing 429 errors.
The issue was in the token bucket algorithm implementation.

Fixes #456
```

### Change Guidelines
- **feat**: New features, significant functionality
- **fix**: Bug fixes, small corrections
- **docs**: Documentation only changes
- **style**: Code style, formatting changes
- **refactor**: Code restructuring without behavior changes
- **test**: Test additions, modifications
- **chore**: Build process, dependency updates

## Troubleshooting

### Common Issues
- **Message Mismatch**: Commit message doesn't match changes
- **Format Issues**: Non-conventional commit format
- **Size Issues**: Commit too large for claimed type
- **Type Confusion**: Wrong commit type used

### Debugging
```bash
# Check semantic commit analysis
gphc check --semantic-commits --verbose

# Validate commit message format
gphc check --validate-commit-message

# Analyze specific commit
gphc check --analyze-commit HEAD
```

## Next Steps
- [Pre-commit Hooks](pre-commit-hooks.md) - Pre-commit integration guide
- [Health Checks](health-checks.md) - Understanding health check categories
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
