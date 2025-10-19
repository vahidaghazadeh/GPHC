# Health Checks Guide

This guide explains the different health check categories and how GPHC evaluates repository health.

## Health Check Categories

### 1. Documentation & Project Structure (25 points)

#### Essential Files Check
- **README.md** (5 points): Project documentation and setup instructions
- **LICENSE** (5 points): License file for legal clarity
- **CONTRIBUTING.md** (5 points): Guidelines for contributors
- **CODE_OF_CONDUCT.md** (5 points): Community behavior standards
- **Setup Instructions** (5 points): Clear installation and usage guide

#### Gitignore Validation
- **Common Patterns**: Checks for standard .gitignore patterns
- **Language-specific**: Validates language-specific ignore patterns
- **Build Artifacts**: Ensures build outputs are ignored
- **IDE Files**: Checks for IDE-specific ignore patterns

### 2. Commit History Quality (30 points)

#### Conventional Commits
- **Format Validation**: Ensures commits follow conventional format
- **Type Checking**: Validates commit types (feat:, fix:, docs:, etc.)
- **Scope Validation**: Checks for proper scope usage
- **Breaking Changes**: Identifies breaking change indicators

#### Message Quality
- **Length Validation**: Ensures messages stay within 72 characters
- **Subject Line**: Validates clear, descriptive subject lines
- **Body Content**: Checks for detailed commit descriptions
- **Imperative Mood**: Ensures commands use imperative mood

#### Commit Size Analysis
- **Large Commits**: Identifies commits with excessive changes
- **God Commits**: Detects commits that change too many files
- **Atomic Commits**: Encourages small, focused commits
- **Change Distribution**: Analyzes commit size patterns

### 3. Git Cleanup & Hygiene (25 points)

#### Branch Management
- **Merged Branches**: Identifies branches that can be deleted
- **Stale Branches**: Finds branches with no activity for 60+ days
- **Branch Protection**: Checks for main branch protection rules
- **Remote Tracking**: Validates remote branch synchronization

#### Stash Management
- **Stash Analysis**: Reviews Git stash entries
- **Age Detection**: Identifies old stashes (>30 days)
- **Stash Content**: Analyzes stash contents and relevance
- **Cleanup Recommendations**: Suggests stash cleanup actions

#### Repository Hygiene
- **Bare Repository Check**: Validates repository structure
- **Orphaned Objects**: Identifies unreachable Git objects
- **Repository Size**: Monitors repository size growth
- **Garbage Collection**: Suggests Git maintenance tasks

### 4. Codebase Structure (20 points)

#### Test Coverage
- **Test Directory**: Checks for test directories
- **Test Files**: Identifies test file patterns
- **Test Ratio**: Analyzes code-to-test ratio
- **Test Quality**: Validates test file structure

#### Directory Organization
- **Source Structure**: Validates source code organization
- **Directory Depth**: Checks for excessive nesting
- **File Distribution**: Analyzes file distribution patterns
- **Naming Conventions**: Validates directory naming

#### Code Quality Indicators
- **Large Directories**: Identifies directories with too many files
- **Empty Directories**: Finds empty or unused directories
- **Documentation Files**: Checks for inline documentation
- **Configuration Files**: Validates configuration structure

## Health Score Calculation

### Scoring System
- **Pass**: Check passes, full points awarded
- **Warning**: Check has issues, partial points awarded
- **Fail**: Check fails, no points awarded

### Weighted Scoring
Each category has different weights based on importance:
- Documentation: 25% of total score
- Commit Quality: 30% of total score
- Git Hygiene: 25% of total score
- Codebase Structure: 20% of total score

### Grade Assignment
- **A+ (95-100)**: Excellent repository health
- **A (90-94)**: Very good repository health
- **A- (85-89)**: Good repository health
- **B+ (80-84)**: Above average repository health
- **B (75-79)**: Average repository health
- **B- (70-74)**: Below average repository health
- **C+ (65-69)**: Poor repository health
- **C (60-64)**: Very poor repository health
- **D (50-59)**: Failing repository health
- **F (0-49)**: Critical repository health issues

## Understanding Check Results

### Check Status
- **PASS**: Check passed successfully
- **WARN**: Check has minor issues
- **FAIL**: Check failed with significant issues

### Check Details
Each check provides:
- **ID**: Unique identifier (e.g., DOC-101)
- **Name**: Human-readable check name
- **Message**: Detailed description of the issue
- **Recommendations**: Suggested actions to improve
- **Score**: Points awarded for this check

### Example Check Result
```
FAIL [DOC-101] README.md exists
Message: README.md file is missing from repository root
Recommendations:
  - Create a README.md file in the repository root
  - Include project description and setup instructions
  - Add usage examples and contribution guidelines
Score: 0/5
```

## Customizing Health Checks

### Configuration File
Create a `gphc.yml` file to customize health checks:

```yaml
# Health check configuration
health_check:
  min_score: 70
  fail_on_warnings: false
  
  # Category weights
  weights:
    documentation: 25
    commits: 30
    hygiene: 25
    structure: 20

# Custom checks
custom_checks:
  - id: CUSTOM-900
    name: "Has SECURITY.md"
    path: "SECURITY.md"
    score: 5
    required: true
```

### Custom Rules
Define project-specific health checks:

```yaml
custom_checks:
  - id: CUSTOM-901
    name: "Has API Documentation"
    path: "docs/api.md"
    score: 3
    
  - id: CUSTOM-902
    name: "No TODO Comments"
    pattern: "TODO|FIXME|HACK"
    score: 2
    required: false
```

## Best Practices

### For High Health Scores
1. **Maintain Documentation**: Keep README.md updated
2. **Follow Conventions**: Use conventional commit format
3. **Clean Branches**: Regularly delete merged branches
4. **Organize Code**: Maintain clear directory structure
5. **Write Tests**: Include comprehensive test coverage

### For Team Projects
1. **Set Standards**: Define team coding standards
2. **Use Templates**: Create issue and PR templates
3. **Automate Checks**: Integrate GPHC into CI/CD
4. **Regular Reviews**: Schedule regular health check reviews
5. **Continuous Improvement**: Track health score trends

## Troubleshooting

### Common Issues

#### Low Documentation Score
- Missing README.md file
- Incomplete setup instructions
- Missing license file
- No contribution guidelines

#### Poor Commit Quality
- Non-conventional commit messages
- Commit messages too long
- Large commits with many changes
- Inconsistent commit formatting

#### Git Hygiene Issues
- Many merged branches not deleted
- Old stashes not cleaned up
- Stale branches with no activity
- Missing branch protection rules

#### Codebase Structure Problems
- No test directories
- Poor directory organization
- Too many files in root directory
- Missing source code structure

## Next Steps

- [Pre-commit Hooks](pre-commit-hooks.md) - Pre-commit integration guide
- [Custom Rules](custom-rules.md) - Custom rule engine configuration
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
