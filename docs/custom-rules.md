# Custom Rules Guide

This guide covers creating custom health checks using the rule engine.

## Overview

The Custom Rules Engine allows you to define project-specific health checks that go beyond the standard GPHC checks.

## Basic Configuration

### Creating Custom Rules
Create a `gphc.yml` file in your repository root:

```yaml
# Custom rules configuration
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5
    required: true
    
  - id: CUSTOM-901
    name: "Has API Documentation"
    path: "docs/api.md"
    score: 3
    required: false
    
  - id: CUSTOM-902
    name: "No TODO Comments"
    pattern: "TODO|FIXME|HACK"
    score: 2
    required: false
```

## Rule Types

### File Existence Rules
```yaml
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5
    required: true
    description: "Security policy file must exist"
```

### Content Pattern Rules
```yaml
custom_checks:
  - id: CUSTOM-901
    name: "No TODO Comments"
    pattern: "TODO|FIXME|HACK"
    score: 2
    required: false
    description: "Code should not contain TODO comments"
```

### Directory Structure Rules
```yaml
custom_checks:
  - id: CUSTOM-902
    name: "Has Tests Directory"
    path: "tests/"
    type: "directory"
    score: 3
    required: true
```

### File Size Rules
```yaml
custom_checks:
  - id: CUSTOM-903
    name: "No Large Files"
    pattern: ".*"
    max_size: 1048576  # 1MB
    score: 2
    required: false
```

## Advanced Configuration

### Complex Rules
```yaml
custom_checks:
  - id: CUSTOM-904
    name: "Has Proper License"
    path: "LICENSE"
    score: 5
    required: true
    validation:
      min_lines: 10
      must_contain: ["MIT", "Apache", "GPL"]
      
  - id: CUSTOM-905
    name: "No Hardcoded Secrets"
    pattern: "(password|secret|key)\\s*=\\s*['\"][^'\"]+['\"]"
    score: 3
    required: false
    exclude_paths: ["*.example", "*.template"]
```

### Conditional Rules
```yaml
custom_checks:
  - id: CUSTOM-906
    name: "Has Dockerfile"
    path: "Dockerfile"
    score: 3
    required: false
    condition:
      if_file_exists: "package.json"
      then_required: true
```

## Rule Execution

### Running Custom Rules
```bash
# Run all checks including custom rules
git hc check

# Run only custom rules
git hc check --custom-only

# Validate custom rules configuration
git hc check --validate-config
```

### Rule Results
```
Custom Rules Results
===================

PASS [CUSTOM-900] Has SECURITY.md
  Message: SECURITY.md file exists
  Score: 5/5

FAIL [CUSTOM-901] Has API Documentation
  Message: docs/api.md file is missing
  Recommendations:
    - Create API documentation file
    - Include endpoint descriptions
    - Add request/response examples
  Score: 0/3

WARN [CUSTOM-902] No TODO Comments
  Message: Found 3 TODO comments in code
  Recommendations:
    - Remove or resolve TODO comments
    - Create issues for remaining TODOs
  Score: 1/2
```

## Best Practices

### Rule Design
1. **Clear Names**: Use descriptive rule names
2. **Appropriate Scores**: Assign realistic scores
3. **Helpful Messages**: Provide actionable recommendations
4. **Team Consensus**: Get team agreement on rules
5. **Regular Review**: Review and update rules regularly

### Rule Categories
```yaml
# Documentation rules
documentation:
  - id: DOC-CUSTOM-001
    name: "Has CHANGELOG.md"
    path: "CHANGELOG.md"
    score: 3

# Security rules
security:
  - id: SEC-CUSTOM-001
    name: "No Hardcoded Secrets"
    pattern: "password\\s*=\\s*['\"][^'\"]+['\"]"
    score: 5

# Code quality rules
quality:
  - id: QUAL-CUSTOM-001
    name: "No Console Logs"
    pattern: "console\\.log"
    score: 2
```

## Integration

### With CI/CD
```yaml
# GitHub Actions
- name: Custom Rules Check
  run: git hc check --custom-only --min-score 80
```

### With Pre-commit
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: gphc-custom-rules
        name: GPHC Custom Rules
        entry: git hc check --custom-only
        language: system
```

## Troubleshooting

### Common Issues
- **Rule Not Found**: Check rule ID and configuration
- **Pattern Issues**: Validate regex patterns
- **Performance**: Optimize complex rules
- **False Positives**: Adjust rule sensitivity

### Debugging
```bash
# Verbose output
git hc check --verbose

# Debug custom rules
git hc check --debug-custom-rules

# Validate configuration
git hc check --validate-config
```

## Next Steps
- [Basic Usage](basic-usage.md) - Getting started with GPHC
- [Health Checks](health-checks.md) - Understanding health check categories
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
