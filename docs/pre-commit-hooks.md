# Pre-commit Hooks Guide

This guide covers integrating GPHC with pre-commit hooks for automated health checks.

## Overview

Pre-commit hooks allow you to run GPHC checks automatically before commits are made, ensuring code quality and repository health standards are maintained.

## Basic Usage

### Running Pre-commit Checks
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

### Exit Codes
- **0**: All checks passed
- **1**: One or more checks failed
- **2**: Error occurred during check

## Integration Examples

### Pre-commit Framework
Add to `.pre-commit-config.yaml`:
```yaml
repos:
  - repo: local
    hooks:
      - id: git-hc-pre-commit
        name: Git HC Pre-commit Check
        entry: git hc pre-commit
        language: system
        stages: [pre-commit]
        pass_filenames: false
```

### Husky (Node.js)
Add to `.husky/pre-commit`:
```bash
#!/bin/sh
git hc pre-commit
```

### Git Hooks
Create `.git/hooks/pre-commit`:
```bash
#!/bin/sh
git hc pre-commit
```

## Configuration

### Pre-commit Settings
```yaml
# git-hc.yml
pre_commit:
  enabled: true
  check_staged_files: true
  validate_commit_message: true
  check_large_files: true
  check_sensitive_files: true
  
  # File size limits
  max_file_size: 1048576  # 1MB
  
  # Sensitive file patterns
  sensitive_patterns:
    - ".env"
    - "*.key"
    - "*.pem"
    - "secrets.json"
```

## Troubleshooting

### Common Issues
- **Command not found**: Ensure Git HC is set up properly
- **Permission denied**: Check file permissions
- **Slow execution**: Optimize file patterns

## Next Steps
- [Basic Usage](basic-usage.md) - Getting started with GPHC
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
