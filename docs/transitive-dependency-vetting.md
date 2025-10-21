# Transitive Dependency Vetting

## Overview

Transitive Dependency Vetting is a critical security feature that performs deep analysis of both direct and indirect dependencies to detect security vulnerabilities. This helps identify supply chain attacks and ensures comprehensive security coverage.

## Why It's Important

Many security vulnerabilities come from transitive dependencies (dependencies of dependencies) rather than direct dependencies. Attackers often target these indirect dependencies because they are less monitored and can provide a pathway into your application.

### Key Benefits:
- **Supply Chain Security**: Detects vulnerabilities in the entire dependency chain
- **Comprehensive Coverage**: Analyzes both direct and transitive dependencies
- **Risk Assessment**: Provides detailed vulnerability scoring and categorization
- **Actionable Insights**: Shows which direct dependency introduces vulnerable transitive dependencies

## Supported Project Types

GPHC supports dependency analysis for the following project types:

### Go Projects
- **Manifest**: `go.mod`
- **Analysis**: Uses `go list -m all` to build complete dependency tree
- **Direct Dependencies**: Parsed from `go.mod` require statements

### Node.js Projects
- **Manifests**: `package.json`, `package-lock.json`, `yarn.lock`
- **Analysis**: Uses `npm ls --json` or parses lock files
- **Tree Structure**: Maintains hierarchical dependency relationships

### Python Projects
- **Manifests**: `requirements.txt`, `Pipfile`, `Pipfile.lock`
- **Analysis**: Uses `pipdeptree --json` or parses requirements files
- **Dependencies**: Both pip and pipenv formats supported

### Rust Projects
- **Manifests**: `Cargo.toml`, `Cargo.lock`
- **Analysis**: Uses `cargo tree --format json`
- **Complete Tree**: Shows all transitive dependencies

### Java Projects
- **Manifests**: `pom.xml`, `build.gradle`
- **Analysis**: Uses `mvn dependency:tree` or parses XML/Gradle files
- **Maven/Gradle**: Supports both build systems

## Basic Usage

### Command Syntax
```bash
git hc security dependencies [flags]
```

### Basic Scan
```bash
# Scan current repository
git hc security dependencies

# Scan specific repository
git hc security dependencies /path/to/repo
```

### Deep Analysis
```bash
# Deep transitive analysis (default)
git hc security dependencies --depth deep

# Shallow analysis (direct dependencies only)
git hc security dependencies --depth shallow
```

## Command Options

### Scan Configuration
- `--depth string`: Scan depth - `shallow` (direct only) or `deep` (transitive) (default: "deep")
- `--direct-only`: Only check direct dependencies (same as --depth shallow)
- `--severity string`: Minimum severity level - `low`, `medium`, `high`, `critical` (default: "low")

### Output Options
- `--format string`: Output format - `table`, `json`, `yaml` (default: "table")
- `--output string`: Output file path for results
- `--tree`: Show dependency tree structure (default: true)

## Example Output

### Table Format
```
🔍 Scanning transitive dependencies for vulnerabilities...
Repository: /path/to/project
Scan depth: deep
Minimum severity: low
Direct dependencies only: false
Show dependency tree: true

📊 Dependency Scan Results
==========================

Project Type: nodejs
Total Dependencies: 1,247
Vulnerable Dependencies: 3
Critical Vulnerabilities: 1
High Vulnerabilities: 2
Medium Vulnerabilities: 0
Low Vulnerabilities: 0
Security Score: 75/100

🌳 Dependency Tree
==================

Tree display not yet implemented in this version.
Use --format json to see detailed dependency information.

🚨 VULNERABILITIES FOUND!

Immediate Actions Required:
1. Update vulnerable dependencies to secure versions
2. Review dependency tree to identify root causes
3. Consider removing unnecessary dependencies
4. Implement dependency scanning in CI/CD pipeline

Tools for Dependency Management:
- npm audit fix (Node.js)
- go get -u (Go)
- pip install --upgrade (Python)
- cargo update (Rust)
- mvn versions:use-latest-releases (Java)

Prevention:
- Use dependency scanning tools in CI/CD
- Regularly update dependencies
- Use lock files (package-lock.json, go.sum, etc.)
- Monitor security advisories
```

### JSON Format
```bash
git hc security dependencies --format json
```

```json
{
  "id": "TRANSITIVE-DEPS",
  "name": "Transitive Dependency Vetting",
  "status": 1,
  "score": 75,
  "message": "Found 3 vulnerable dependencies (1 critical, 2 high)",
  "details": [
    "Project Type: nodejs",
    "Total Dependencies: 1247",
    "Vulnerable Dependencies: 3",
    "Critical Vulnerabilities: 1",
    "High Vulnerabilities: 2",
    "Medium Vulnerabilities: 0",
    "Low Vulnerabilities: 0"
  ],
  "category": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Vulnerability Detection

### Known Vulnerabilities
GPHC includes a database of known vulnerabilities for common packages:

#### Critical Vulnerabilities
- **Log4Shell (CVE-2021-44228)**: Log4j remote code execution
- **Spring4Shell (CVE-2022-22965)**: Spring Framework remote code execution

#### High Severity
- **Lodash Command Injection (CVE-2021-23337)**: Command injection in lodash
- **Axios SSRF (CVE-2020-28168)**: Server-side request forgery

#### Medium Severity
- **Various CVEs**: Medium-severity vulnerabilities in popular packages

### Severity Scoring
- **Critical**: CVSS 9.0-10.0 (Score penalty: -20 points)
- **High**: CVSS 7.0-8.9 (Score penalty: -10 points)
- **Medium**: CVSS 4.0-6.9 (Score penalty: -5 points)
- **Low**: CVSS 0.1-3.9 (Score penalty: -2 points)

## CI/CD Integration

### GitHub Actions
```yaml
name: Dependency Security Scan
on: [push, pull_request]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Install GPHC
        run: go install github.com/opsource/gphc/cmd/gphc@latest
      - name: Scan Dependencies
        run: |
          git hc security dependencies --severity medium --format json --output deps-report.json
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: dependency-scan-results
          path: deps-report.json
```

### GitLab CI
```yaml
dependency_scan:
  stage: security
  image: golang:1.19
  before_script:
    - go install github.com/opsource/gphc/cmd/gphc@latest
  script:
    - git hc security dependencies --severity high --format json --output deps-report.json
  artifacts:
    reports:
      junit: deps-report.json
  only:
    - merge_requests
    - main
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Dependency Scan') {
            steps {
                sh 'go install github.com/opsource/gphc/cmd/gphc@latest'
                sh 'git hc security dependencies --severity medium --format json --output deps-report.json'
                archiveArtifacts artifacts: 'deps-report.json'
            }
        }
    }
    post {
        always {
            publishHTML([
                allowMissing: false,
                alwaysLinkToLastBuild: true,
                keepAll: true,
                reportDir: '.',
                reportFiles: 'deps-report.json',
                reportName: 'Dependency Scan Report'
            ])
        }
    }
}
```

## Configuration

### gphc.yml Configuration
```yaml
# Dependency scanning configuration
dependency_scanning:
  enabled: true
  severity_threshold: "medium"
  scan_depth: "deep"
  include_dev_dependencies: true
  
  # Vulnerability database settings
  vulnerability_db:
    update_frequency: "daily"
    sources:
      - "github_advisory"
      - "nvd"
      - "oss_index"
  
  # Exclude specific packages
  exclusions:
    - "package-name@version"
    - "another-package@*"
  
  # Custom vulnerability patterns
  custom_patterns:
    - name: "Custom Pattern"
      pattern: "regex_pattern"
      severity: "high"
      description: "Custom vulnerability description"
```

## Best Practices

### Regular Scanning
1. **Daily Scans**: Run dependency scans daily in CI/CD
2. **Pre-commit Hooks**: Scan dependencies before commits
3. **Release Gates**: Block releases with critical vulnerabilities

### Dependency Management
1. **Lock Files**: Always use lock files (package-lock.json, go.sum, etc.)
2. **Regular Updates**: Update dependencies regularly
3. **Minimal Dependencies**: Only include necessary dependencies
4. **Version Pinning**: Pin specific versions for critical dependencies

### Security Monitoring
1. **Advisory Subscriptions**: Subscribe to security advisories
2. **Automated Alerts**: Set up alerts for new vulnerabilities
3. **Patch Management**: Have a process for applying security patches
4. **Risk Assessment**: Regularly assess dependency risks

## Troubleshooting

### Common Issues

#### No Dependencies Found
```
Project Type: 
Total Dependencies: 0
```
**Solution**: Ensure your project has a supported manifest file (go.mod, package.json, etc.)

#### Build Errors
```
Error: failed to run go list: exit status 1
```
**Solution**: Ensure your Go project builds successfully with `go mod tidy`

#### Permission Issues
```
Error: package-lock.json not found: permission denied
```
**Solution**: Check file permissions and ensure GPHC has read access

### Debug Mode
```bash
# Enable verbose output
git hc security dependencies --verbose

# Check specific project type
git hc security dependencies --project-type nodejs
```

## Integration with Other Tools

### Snyk Integration
```bash
# Run GPHC and Snyk together
git hc security dependencies --format json > gphc-deps.json
snyk test --json > snyk-deps.json
```

### OWASP Dependency Check
```bash
# Combine with OWASP Dependency Check
git hc security dependencies --format json > gphc-report.json
dependency-check.sh --format JSON --out owasp-report.json
```

## Next Steps

1. [Pre-commit Hooks](pre-commit-hooks.md) - Set up dependency scanning in pre-commit hooks
2. [CI/CD Integration](ci-cd-integration.md) - Add to your CI/CD pipeline
3. [Health Checks](health-checks.md) - Run comprehensive health assessment
4. [Secret Scanning](secret-scanning.md) - Scan for exposed secrets and credentials
