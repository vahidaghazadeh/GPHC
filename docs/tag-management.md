# Tag Management

Tag Management module provides comprehensive Git tag and release validation, semantic versioning enforcement, and automated changelog generation.

## Overview

The Tag Management feature helps maintain healthy release practices by:
- Validating semantic versioning compliance
- Monitoring tag freshness and frequency
- Detecting unreleased commits
- Enforcing annotated tag policies
- Suggesting next semantic versions
- Generating changelogs automatically

## Basic Usage

### Check Tag Health
```bash
# Basic tag health check
git hc tags

# Check specific repository
git hc tags /path/to/repo
```

### Suggest Next Version
```bash
# Get suggested next semantic version
git hc tags --suggest
```

### Generate Changelog
```bash
# Generate changelog to file
git hc tags --changelog CHANGELOG.md

# Generate changelog to stdout
git hc tags --changelog ""
```

### Policy Enforcement
```bash
# Fail if tag policies are violated (for CI/CD)
git hc tags --enforce-tags
```

### Combined Operations
```bash
# Full tag analysis with suggestions and changelog
git hc tags --suggest --changelog CHANGELOG.md --enforce-tags
```

## Features

### 1. Semantic Version Validation
- Validates tags follow Semantic Versioning (`vX.Y.Z`)
- Supports pre-release tags (`-beta`, `-rc`, `-alpha`)
- Identifies invalid or inconsistent tag formats

**Example Output:**
```
Semantic Versioning: OK
Invalid tags (non-semver): v1.0, release-2023
```

### 2. Tag Freshness Monitoring
- Tracks days since last tag
- Configurable threshold (default: 45 days)
- Warns about stale releases

**Example Output:**
```
Last tag: v1.4.2 (60 days ago)
Last tag older than 45 days
```

### 3. Unreleased Commits Detection
- Counts commits since latest tag
- Configurable threshold (default: 3 commits)
- Helps maintain regular release cycles

**Example Output:**
```
Unreleased commits since last tag: 5
Too many unreleased commits (>3)
```

### 4. Annotated vs Lightweight Tags
- Enforces annotated tag requirements
- Calculates annotation percentage
- Recommends best practices

**Example Output:**
```
Annotated tags: 100%
Some release tags are lightweight; annotate release tags
```

### 5. Auto-Suggest Next Tag
- Analyzes commit messages since last tag
- Detects conventional commit patterns
- Suggests appropriate version bump

**Commit Pattern Detection:**
- `feat:` → Minor version bump
- `fix:` → Patch version bump
- `feat!:` or `BREAKING CHANGE` → Major version bump

**Example Output:**
```
Auto-suggested next tag: v1.5.0
```

### 6. Changelog Generation
- Groups commits by conventional commit types
- Generates structured changelog
- Supports multiple output formats

**Generated Sections:**
- Features
- Fixes
- Docs
- Refactors
- Others

**Example Changelog:**
```markdown
# Changelog

Changes since v1.4.2

## Features
- feat: add new authentication system
- feat: implement user dashboard

## Fixes
- fix: resolve memory leak in cache
- fix: correct validation logic

## Docs
- docs: update API documentation
```

## Configuration

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--suggest` | Suggest next semantic version | false |
| `--changelog` | Generate changelog to file | "" |
| `--enforce-tags` | Fail if policies violated | false |

### Health Check Integration

Tag Management is automatically included in the main health check:

```bash
# Tag checker runs as part of comprehensive health check
git hc check
```

**Health Check Output:**
```
[TAGS-901] Tag & Release Health: PASS (85/100)
- Semantic Versioning: OK
- Last tag: v1.4.2 (15 days ago)
- Unreleased commits since last tag: 2
- Annotated tags: 100%
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Tag Policy Check
on: [push, pull_request]

jobs:
  tag-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check Tag Policies
        run: git hc tags --enforce-tags
```

### GitLab CI
```yaml
tag_check:
  stage: test
  script:
    - git hc tags --enforce-tags
  rules:
    - if: $CI_PIPELINE_SOURCE == "push"
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Tag Policy Check') {
            steps {
                sh 'git hc tags --enforce-tags'
            }
        }
    }
}
```

## Best Practices

### 1. Semantic Versioning
- Use `vX.Y.Z` format for all release tags
- Follow semantic versioning rules:
  - Major (X): Breaking changes
  - Minor (Y): New features (backward compatible)
  - Patch (Z): Bug fixes (backward compatible)

### 2. Annotated Tags
- Always use annotated tags for releases
- Include meaningful tag messages
- Use `git tag -a v1.0.0 -m "Release v1.0.0"`

### 3. Regular Releases
- Maintain regular release cycles
- Don't accumulate too many unreleased commits
- Use automated release workflows

### 4. Conventional Commits
- Follow conventional commit format
- Use `feat:`, `fix:`, `docs:` prefixes
- Include breaking change indicators

## Troubleshooting

### Common Issues

**No tags found:**
```
Status: WARNING
Message: No tags found in repository
```
- Create your first release tag: `git tag v0.1.0`

**Invalid semantic version:**
```
Invalid tags (non-semver): v1.0, release-2023
```
- Rename tags to follow `vX.Y.Z` format
- Use `git tag -d old-tag` and `git tag v1.0.0`

**Too many unreleased commits:**
```
Too many unreleased commits (>3)
```
- Create a new release tag
- Or adjust the threshold in configuration

**Lightweight tags detected:**
```
Some release tags are lightweight; annotate release tags
```
- Recreate tags as annotated: `git tag -a v1.0.0 -m "Release v1.0.0"`

### Debug Mode

For detailed debugging, check git commands manually:

```bash
# List all tags
git tag

# Check tag types
git for-each-ref --format='%(refname:short) %(objecttype)' refs/tags

# Count unreleased commits
git rev-list $(git describe --tags --abbrev=0)..HEAD --count

# Get latest tag date
git for-each-ref --sort=-creatordate --format='%(refname:short)|%(creatordate:iso8601)' refs/tags | head -n1
```

## Advanced Usage

### Custom Thresholds

While not yet configurable via CLI, thresholds can be modified in the source code:

```go
// In internal/checkers/tag_checker.go
maxDaysSinceLastTag:  45,  // Adjust threshold
maxUnreleasedCommits: 3,   // Adjust threshold
requireAnnotatedTags: true, // Enable/disable requirement
```

### Integration with Release Tools

Tag Management can be integrated with release automation tools:

```bash
# Pre-release check
git hc tags --enforce-tags || exit 1

# Generate changelog for release notes
git hc tags --changelog RELEASE_NOTES.md

# Get suggested version for automation
NEXT_VERSION=$(git hc tags --suggest | grep "Auto-suggested" | cut -d' ' -f4)
```

## Examples

### Complete Release Workflow
```bash
# 1. Check current tag health
git hc tags

# 2. Get suggested next version
git hc tags --suggest

# 3. Generate changelog
git hc tags --changelog CHANGELOG.md

# 4. Create annotated tag
git tag -a v1.5.0 -m "Release v1.5.0"

# 5. Verify tag health
git hc tags --enforce-tags
```

### CI/CD Pipeline Integration
```bash
# Fail pipeline if tag policies violated
if ! git hc tags --enforce-tags; then
    echo "Tag policy violations detected"
    exit 1
fi

# Generate release notes
git hc tags --changelog RELEASE_NOTES.md

# Get next version for automated releases
NEXT_VERSION=$(git hc tags --suggest | grep "Auto-suggested" | cut -d' ' -f4)
echo "Next version: $NEXT_VERSION"
```

## Related Documentation

- [Basic Usage](basic-usage.md)
- [Health Checks](health-checks.md)
- [CI/CD Integration](ci-cd-integration.md)
- [Conventional Commits](semantic-commits.md)
