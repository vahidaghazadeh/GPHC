# Codebase Analysis Guide

This guide covers lightweight codebase structure analysis and smell detection.

## Overview

Codebase analysis provides quick insights into project structure, test coverage, and maintainability without requiring complex AST analysis.

## Basic Usage

### Analyzing Codebase Structure
```bash
# Analyze codebase structure
git hc codebase

# Example output:
Codebase Structure Analysis
===========================

Repository: /path/to/project
Analysis Date: 2024-01-15 10:30:00

Structure Overview:
┌─────────────────┬─────────┬─────────┬─────────┐
│ Directory       │ Files   │ Type    │ Status  │
├─────────────────┼─────────┼─────────┼─────────┤
│ src/            │ 45      │ Source  │ ✅ Good │
│ tests/          │ 23      │ Tests   │ ✅ Good │
│ docs/           │ 8       │ Docs    │ ✅ Good │
│ config/         │ 5       │ Config  │ ✅ Good │
└─────────────────┴─────────┴─────────┴─────────┘

Test Coverage Analysis:
✅ Test directory exists: tests/
✅ Test files found: 23
✅ Code-to-test ratio: 1:0.51 (Good)
✅ Test file patterns: test_*.py, *_test.py

Structure Health Score: 85/100 (B+)
```

### Detailed Analysis
```bash
# Get detailed codebase analysis
git hc codebase --detailed

# Example output:
Detailed Codebase Analysis
=========================

File Distribution:
- Total files: 81
- Source files: 45
- Test files: 23
- Documentation: 8
- Configuration: 5

Directory Analysis:
✅ src/ directory exists (45 files)
✅ tests/ directory exists (23 files)
✅ docs/ directory exists (8 files)
⚠️ config/ directory has only 5 files (consider consolidation)

Test Analysis:
✅ Test directory structure is good
✅ Test file naming follows conventions
✅ Code-to-test ratio is healthy (1:0.51)
⚠️ Some source files lack corresponding tests

Structure Issues:
⚠️ Large src/ directory (45 files) - consider subdirectories
⚠️ Some empty directories detected
⚠️ Deep nesting detected (max depth: 6)

Recommendations:
1. Consider breaking src/ into subdirectories
2. Add tests for untested source files
3. Remove empty directories
4. Reduce directory nesting depth
```

## Configuration

### Codebase Analysis Settings
```yaml
# gphc.yml
codebase_analysis:
  enabled: true
  
  # Analysis settings
  max_files_per_directory: 1000
  max_directory_depth: 10
  min_test_ratio: 0.3
  
  # File patterns
  source_patterns:
    - "*.py"
    - "*.js"
    - "*.go"
    - "*.java"
  
  test_patterns:
    - "test_*.py"
    - "*_test.py"
    - "*.test.js"
    - "*_test.go"
  
  # Directory patterns
  source_directories:
    - "src/"
    - "lib/"
    - "app/"
    - "source/"
  
  test_directories:
    - "tests/"
    - "test/"
    - "spec/"
    - "tests/"
```

### Custom Rules
```yaml
# gphc.yml
codebase_analysis:
  custom_rules:
    - id: "CODE-CUSTOM-001"
      name: "Has Source Directory"
      check: "source_directory"
      score: 5
      
    - id: "CODE-CUSTOM-002"
      name: "Has Test Directory"
      check: "test_directory"
      score: 5
      
    - id: "CODE-CUSTOM-003"
      name: "Good Test Ratio"
      check: "test_ratio"
      min_ratio: 0.3
      score: 3
```

## Analysis Categories

### Directory Structure
- **Source Directories**: Check for proper source organization
- **Test Directories**: Validate test structure
- **Documentation**: Check documentation organization
- **Configuration**: Validate config file structure

### File Distribution
- **File Counts**: Analyze file distribution
- **File Types**: Check file type patterns
- **File Sizes**: Identify large files
- **File Naming**: Validate naming conventions

### Test Coverage
- **Test Presence**: Check for test files
- **Test Structure**: Validate test organization
- **Test Ratio**: Analyze code-to-test ratio
- **Test Quality**: Check test file quality

### Structure Health
- **Directory Depth**: Check nesting levels
- **File Organization**: Validate file organization
- **Empty Directories**: Identify empty directories
- **Structure Patterns**: Check for common patterns

## Use Cases

### Project Health Monitoring
```bash
# Monitor project structure
git hc codebase --period 30

# Check for structure issues
git hc codebase --check-issues

# Analyze structure trends
git hc codebase --trends
```

### Team Onboarding
```bash
# Generate structure overview
git hc codebase --format markdown --output structure-overview.md

# Check for onboarding issues
git hc codebase --check-onboarding

# Generate team guide
git hc codebase --generate-guide
```

### Code Review Analysis
```bash
# Analyze structure for reviews
git hc codebase --review-analysis

# Check for review issues
git hc codebase --check-review-issues

# Generate review checklist
git hc codebase --review-checklist
```

## Integration Examples

### CI/CD Integration
```yaml
# GitHub Actions
- name: Codebase Analysis
  run: git hc codebase --format json --output codebase.json

- name: Check Structure Health
  run: git hc codebase --min-score 80
```

### Team Reporting
```bash
# Generate weekly structure report
git hc codebase --period 7 --format markdown --output structure-report.md

# Generate monthly analysis
git hc codebase --period 30 --format json --output monthly-analysis.json
```

## Best Practices

### For Teams
1. **Regular Analysis**: Analyze structure regularly
2. **Structure Standards**: Establish structure standards
3. **Test Coverage**: Maintain good test coverage
4. **Documentation**: Keep documentation organized
5. **Code Organization**: Maintain good code organization

### For Organizations
1. **Structure Standards**: Establish organization-wide standards
2. **Training Programs**: Implement structure training
3. **Code Reviews**: Include structure in code reviews
4. **Mentoring**: Establish structure mentoring
5. **Continuous Improvement**: Continuously improve structure

## Troubleshooting

### Common Issues
- **No Source Files**: Check file patterns
- **Missing Tests**: Verify test directory
- **Structure Issues**: Check directory organization
- **Analysis Errors**: Validate configuration

### Debugging
```bash
# Test codebase analysis
git hc codebase --test

# Verbose output
git hc codebase --verbose

# Check specific patterns
git hc codebase --check-patterns
```

## Next Steps
- [Health Checks](health-checks.md) - Understanding health check categories
- [Multi-Repository Scan](multi-repository-scan.md) - Batch repository analysis
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
