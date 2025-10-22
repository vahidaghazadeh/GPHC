# Git Policy Validation

## Overview

Git Policy Validation is a comprehensive security feature that validates Git security policies and configurations to ensure proper security practices are in place. This includes commit signature verification, push policies, sensitive file detection, and branch protection settings.

## Why It's Important

Git repositories can have various security vulnerabilities if not properly configured. This feature helps identify and remediate common security issues that could lead to:

- **Identity Spoofing**: Unsigned commits can be forged
- **Data Exposure**: Sensitive files accidentally committed
- **Unauthorized Changes**: Weak push policies allow dangerous operations
- **Branch Tampering**: Unprotected branches vulnerable to malicious changes

## Key Features

### 1. Commit Signature Verification
- **GPG Signature Analysis**: Checks commit signature rates and validity
- **Signature Statistics**: Provides detailed metrics on signed vs unsigned commits
- **Invalid Signature Detection**: Identifies corrupted or invalid signatures
- **Policy Enforcement**: Ensures minimum signature requirements

### 2. Sensitive File Detection
- **Pattern Matching**: Detects common sensitive file patterns
- **Git History Scanning**: Checks entire Git history for sensitive files
- **Gitignore Validation**: Ensures sensitive files are properly ignored
- **Real-time Detection**: Scans current working directory

### 3. Push Policy Validation
- **Force Push Protection**: Checks for dangerous push settings
- **Default Push Behavior**: Validates push.default configuration
- **Credential Storage**: Identifies unsafe credential storage methods
- **Merge Strategy**: Validates merge strategy settings

### 4. Branch Protection Analysis
- **Protected Branch Detection**: Identifies important branches
- **Protection Rule Validation**: Checks for branch protection rules
- **Access Control**: Validates branch access permissions

## Supported File Types

### Critical Severity
- **SSH Keys**: `id_rsa`, `id_dsa`, `id_ed25519`
- **Private Keys**: `*.key`
- **Secrets Files**: `secrets.json`, `credentials.json`
- **Production Configs**: `.env.production`

### High Severity
- **Environment Files**: `.env`, `.env.local`
- **Certificates**: `*.pem`, `*.p12`, `*.pfx`
- **Kubernetes Configs**: `kubeconfig`, `.kube/config`

### Medium Severity
- **Configuration Files**: `config.json`
- **Development Files**: `.env.development`

## Basic Usage

### Command Syntax
```bash
git hc security policy [flags]
```

### Basic Policy Validation
```bash
# Run complete policy validation
git hc security policy

# Validate specific repository
git hc security policy /path/to/repo
```

### Focused Checks
```bash
# Focus on commit signatures
git hc security policy --check-signing

# Focus on sensitive files
git hc security policy --check-files

# Focus on push policies
git hc security policy --check-push

# Focus on branch protection
git hc security policy --check-branches
```

## Command Options

### Check Configuration
- `--check-signing`: Check commit signature verification (default: true)
- `--check-files`: Check for sensitive files (default: true)
- `--check-push`: Check push policies (default: true)
- `--check-branches`: Check branch protection (default: true)

### Output Options
- `--severity string`: Minimum severity level (low, medium, high, critical) (default: "low")
- `--format string`: Output format (table, json, yaml) (default: "table")
- `--output string`: Output file path for results

## Example Output

### Table Format
```
🔍 Validating Git security policies...
Repository: /path/to/project
Check signing: true
Check files: true
Check push policies: true
Check branch protection: true
Minimum severity: low

📊 Git Policy Validation Results
=================================

Total Violations: 3
Signature Rate: 25.0%
Sensitive Files: 2
Push Policies: 1
Branch Protection: 0
Security Score: 65/100

🚨 POLICY VIOLATIONS FOUND!

Immediate Actions Required:
1. Review and fix policy violations
2. Enable commit signing for important commits
3. Add sensitive files to .gitignore
4. Configure branch protection rules
5. Review push policies and permissions

Security Best Practices:
- Enable GPG commit signing
- Use .gitignore for sensitive files
- Configure branch protection
- Use signed commits for releases
- Regular security policy audits
```

### JSON Format
```bash
git hc security policy --format json
```

```json
{
  "id": "GIT-POLICY",
  "name": "Git Policy Validation",
  "status": 1,
  "score": 65,
  "message": "Found 3 Git security policy violations",
  "details": [
    "Total Violations: 3",
    "Signature Rate: 25.0%",
    "Sensitive Files: 2",
    "Push Policies: 1",
    "Branch Protection: 0"
  ],
  "category": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Signature Verification Details

### Signature Status Codes
- **G**: Good signature
- **B**: Bad signature
- **U**: Good signature with unknown validity
- **X**: Good signature that expired
- **Y**: Good signature made by an expired key
- **R**: Good signature made by a revoked key
- **E**: Signature could not be checked
- **N**: No signature

### Signature Rate Thresholds
- **< 50%**: High severity violation
- **50-80%**: Medium severity violation
- **> 80%**: Acceptable signature rate

## Sensitive File Patterns

### Environment Files
```bash
.env
.env.local
.env.production
.env.development
```

### SSH Keys
```bash
id_rsa
id_dsa
id_ed25519
id_ecdsa
```

### Certificates
```bash
*.pem
*.key
*.p12
*.pfx
*.crt
*.cer
```

### Configuration Files
```bash
kubeconfig
.kube/config
config.json
secrets.json
credentials.json
```

## Git Configuration Checks

### Dangerous Settings
- **push.default = matching**: Can push to multiple branches
- **credential.helper = store**: Stores credentials in plain text
- **merge.ours = true**: Can hide merge conflicts

### Recommended Settings
- **push.default = simple**: Safer push behavior
- **credential.helper = cache**: Temporary credential storage
- **receive.denyNonFastForwards = true**: Prevents force pushes

## Branch Protection

### Protected Branches
- `main`
- `master`
- `develop`
- `production`

### Protection Rules
- Require pull request reviews
- Require status checks
- Require up-to-date branches
- Restrict pushes to specific users/teams

## CI/CD Integration

### GitHub Actions
```yaml
name: Git Policy Validation
on: [push, pull_request]

jobs:
  policy-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup GPHC
        run: go install github.com/opsource/gphc/cmd/gphc@latest
      - name: Validate Git Policies
        run: |
          git hc security policy --severity medium --format json --output policy-report.json
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: git-policy-results
          path: policy-report.json
```

### GitLab CI
```yaml
git_policy_validation:
  stage: security
  image: golang:1.19
  before_script:
    - go install github.com/opsource/gphc/cmd/gphc@latest
  script:
    - git hc security policy --severity high --format json --output policy-report.json
  artifacts:
    reports:
      junit: policy-report.json
  only:
    - merge_requests
    - main
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Git Policy Validation') {
            steps {
                sh 'go install github.com/opsource/gphc/cmd/gphc@latest'
                sh 'git hc security policy --severity medium --format json --output policy-report.json'
                archiveArtifacts artifacts: 'policy-report.json'
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
                reportFiles: 'policy-report.json',
                reportName: 'Git Policy Validation Report'
            ])
        }
    }
}
```

## Configuration

### gphc.yml Configuration
```yaml
# Git policy validation configuration
git_policy_validation:
  enabled: true
  severity_threshold: "medium"
  
  # Signature verification settings
  signature_verification:
    enabled: true
    minimum_rate: 80.0
    require_releases: true
  
  # Sensitive file detection
  sensitive_files:
    enabled: true
    patterns:
      - "*.env*"
      - "*.key"
      - "*.pem"
      - "kubeconfig"
      - "secrets.json"
    check_history: true
    check_gitignore: true
  
  # Push policy validation
  push_policies:
    enabled: true
    deny_force_push: true
    require_signed_commits: false
  
  # Branch protection
  branch_protection:
    enabled: true
    protected_branches:
      - "main"
      - "master"
      - "develop"
    require_protection_rules: true
```

## Best Practices

### Commit Signing
1. **Generate GPG Key**: Create a GPG key for commit signing
2. **Configure Git**: Set up Git to use GPG signing
3. **Sign Important Commits**: Sign all release commits
4. **Team Policy**: Require signatures for important branches

### Sensitive File Management
1. **Comprehensive .gitignore**: Include all sensitive file patterns
2. **Regular Audits**: Periodically check for sensitive files
3. **History Cleanup**: Remove sensitive files from Git history
4. **Environment Variables**: Use environment variables instead of files

### Push Policies
1. **Safe Defaults**: Use `push.default = simple`
2. **Force Push Protection**: Enable `receive.denyNonFastForwards`
3. **Branch Protection**: Configure protection rules for important branches
4. **Access Control**: Limit push permissions to trusted users

### Branch Protection
1. **Protect Main Branches**: Enable protection for `main`/`master`
2. **Require Reviews**: Mandate pull request reviews
3. **Status Checks**: Require CI/CD status checks
4. **Up-to-date Branches**: Require branches to be up-to-date

## Troubleshooting

### Common Issues

#### Low Signature Rate
```
Signature Rate: 25.0%
```
**Solution**: Enable GPG commit signing and configure Git to sign commits automatically.

#### Sensitive Files Found
```
Sensitive file found: .env
```
**Solution**: Add the file pattern to `.gitignore` and remove from Git history.

#### Force Push Allowed
```
Force pushes are not denied
```
**Solution**: Configure `receive.denyNonFastForwards = true` in Git config.

#### Missing Branch Protection
```
Protected branch 'main' has no protection rules
```
**Solution**: Configure branch protection rules in your Git hosting platform.

### Debug Mode
```bash
# Enable verbose output
git hc security policy --verbose

# Check specific components
git hc security policy --check-signing --check-files
```

## Integration with Other Tools

### Pre-commit Hooks
```bash
# Add to .pre-commit-config.yaml
- repo: local
  hooks:
    - id: git-policy-validation
      name: Git Policy Validation
      entry: git hc security policy --check-files
      language: system
      pass_filenames: false
```

### Git Hooks
```bash
#!/bin/bash
# .git/hooks/pre-push
git hc security policy --check-signing --check-files
if [ $? -ne 0 ]; then
    echo "Policy validation failed. Push aborted."
    exit 1
fi
```

## Next Steps

1. [Pre-commit Hooks](pre-commit-hooks.md) - Set up policy validation in pre-commit hooks
2. [CI/CD Integration](ci-cd-integration.md) - Add to your CI/CD pipeline
3. [Secret Scanning](secret-scanning.md) - Scan for exposed secrets and credentials
4. [Transitive Dependency Vetting](transitive-dependency-vetting.md) - Analyze dependency vulnerabilities
