# Secret Scanning in Git History

## Overview

Secret scanning is a critical security feature that scans the entire Git history and stashes for exposed secrets, credentials, and sensitive information. This is one of the most important security vulnerabilities in most Git repositories.

## Why It's Important

Even if you remove sensitive files from your current working directory, they remain in the Git history and can be recovered by anyone who clones the repository. This includes:

- API keys and tokens
- Passwords and credentials
- Private keys and certificates
- Database connection strings
- Cloud service credentials

## Features

### Deep History Scanning
- Scans all commits in the repository history
- Analyzes Git stashes for hidden secrets
- Performs comprehensive pattern matching

### Pattern Detection
Detects common secret formats including:

- **AWS Credentials**: Access keys, secret keys
- **GitHub Tokens**: Personal access tokens, app tokens
- **GitLab Tokens**: Personal access tokens
- **Slack Tokens**: Bot tokens, webhook URLs
- **Discord Tokens**: Bot tokens
- **JWT Tokens**: JSON Web Tokens
- **Private Keys**: SSH keys, SSL certificates
- **API Keys**: Generic API key patterns
- **Passwords**: Password field patterns

### Entropy Analysis
- Analyzes random-looking strings for high entropy
- Detects base64-encoded secrets
- Identifies potential secrets that don't match known patterns

## Usage

### Basic Command
```bash
git hc security secrets
```

### Advanced Options
```bash
# Scan with custom settings
git hc security secrets --severity high --confidence 0.8

# Scan specific repository
git hc security secrets /path/to/repository

# Export results to file
git hc security secrets --format json --output secrets-report.json
```

### Command Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--history` | Scan entire Git history | `true` |
| `--stashes` | Scan Git stashes | `true` |
| `--entropy` | Perform entropy analysis | `true` |
| `--severity` | Minimum severity level (low, medium, high) | `medium` |
| `--confidence` | Minimum confidence threshold (0.0-1.0) | `0.7` |
| `--format` | Output format (table, json, yaml) | `table` |
| `--output` | Output file path | stdout |

## Example Output

### No Secrets Found
```
🔍 Scanning for secrets in Git history...
Repository: /path/to/repo
Scanning history: true
Scanning stashes: true
Entropy analysis: true
Minimum severity: medium
Minimum confidence: 0.70

✅ No secrets found in Git history!
```

### Secrets Found
```
🔍 Scanning for secrets in Git history...
Repository: /path/to/repo
Scanning history: true
Scanning stashes: true
Entropy analysis: true
Minimum severity: medium
Minimum confidence: 0.70

🚨 Secrets found in Git history!

• Total secrets found: 3
• High severity secrets: 2
• 1. AWS Access Key (high) in config/aws.yml:15
• 2. GitHub Token (high) in .env:8
• 3. High Entropy String (medium) in scripts/deploy.sh:42

🚨 CRITICAL: Secrets found in Git history!

Immediate Actions Required:
1. Rotate/revoke all exposed credentials immediately
2. Rewrite Git history to remove secrets
3. Notify team members about the exposure

Tools for History Rewriting:
- git filter-repo: https://github.com/newren/git-filter-repo
- BFG Repo-Cleaner: https://rtyley.github.io/bfg-repo-cleaner/

Commands:
# Using git filter-repo
git filter-repo --replace-text <(echo 'SECRET_VALUE==>REDACTED')

# Using BFG
java -jar bfg.jar --replace-text replacements.txt

After rewriting history:
git push --force-with-lease origin main
```

## Secret Types and Patterns

### AWS Credentials
- **Access Key ID**: `AKIA[0-9A-Z]{16}`
- **Secret Access Key**: `[A-Za-z0-9/+=]{40}`

### GitHub Tokens
- **Personal Access Token**: `ghp_[A-Za-z0-9]{36}`
- **App Token**: `ghs_[A-Za-z0-9]{36}`

### GitLab Tokens
- **Personal Access Token**: `glpat-[A-Za-z0-9_-]{20}`

### Slack Tokens
- **Bot Token**: `xox[baprs]-[A-Za-z0-9-]+`

### Discord Tokens
- **Bot Token**: `[MN][A-Za-z\d]{23}\.[\w-]{6}\.[\w-]{27}`

### JWT Tokens
- **Pattern**: `eyJ[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*\.[A-Za-z0-9_-]*`

### Private Keys
- **Pattern**: `-----BEGIN [A-Z ]+ PRIVATE KEY-----`

### Generic Patterns
- **API Keys**: `(?i)(api[_-]?key|apikey)[\s]*[:=][\s]*['"]?([A-Za-z0-9_-]{20,})['"]?`
- **Passwords**: `(?i)(password|passwd|pwd)[\s]*[:=][\s]*['"]?([A-Za-z0-9@#$%^&+=]{8,})['"]?`

## Remediation Steps

### 1. Immediate Actions
1. **Rotate Credentials**: Immediately revoke and regenerate all exposed credentials
2. **Notify Team**: Inform all team members about the exposure
3. **Assess Impact**: Determine what systems were accessible with exposed credentials

### 2. History Rewriting
Use one of these tools to remove secrets from Git history:

#### git filter-repo
```bash
# Install git filter-repo
pip install git-filter-repo

# Create replacement file
echo "AKIA1234567890ABCDEF==>REDACTED" > replacements.txt
echo "ghp_1234567890abcdef1234567890abcdef12345678==>REDACTED" >> replacements.txt

# Rewrite history
git filter-repo --replace-text replacements.txt
```

#### BFG Repo-Cleaner
```bash
# Download BFG
wget https://repo1.maven.org/maven2/com/madgag/bfg/1.14.0/bfg-1.14.0.jar

# Create replacement file
echo "AKIA1234567890ABCDEF==>REDACTED" > replacements.txt

# Rewrite history
java -jar bfg-1.14.0.jar --replace-text replacements.txt
```

### 3. Force Push Changes
```bash
# Force push the cleaned history
git push --force-with-lease origin main

# Notify all team members to re-clone the repository
```

## Prevention

### 1. Use Environment Variables
```bash
# Instead of hardcoding secrets
export AWS_ACCESS_KEY_ID="your-key"
export AWS_SECRET_ACCESS_KEY="your-secret"
```

### 2. Use Secret Management Tools
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **Azure Key Vault**
- **Google Secret Manager**

### 3. Use .gitignore
```gitignore
# Environment files
.env
.env.local
.env.production

# Configuration files with secrets
config/secrets.yml
config/production.yml

# Key files
*.pem
*.key
*.p12
```

### 4. Pre-commit Hooks
```bash
# Install git-secrets
git secrets --install
git secrets --register-aws
git secrets --add 'password.*=.*'
```

### 5. Regular Scanning
```bash
# Add to CI/CD pipeline
git hc security secrets --severity high --confidence 0.8
```

## Integration with CI/CD

### GitHub Actions
```yaml
name: Secret Scan
on: [push, pull_request]

jobs:
  secret-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0  # Fetch full history
      
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
      
      - name: Scan for secrets
        run: git hc security secrets --severity high --confidence 0.8
```

### GitLab CI
```yaml
secret_scan:
  stage: security
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - git hc security secrets --severity high --confidence 0.8
  rules:
    - if: $CI_PIPELINE_SOURCE == "push"
```

## Configuration

### git-hc.yml
```yaml
security:
  secret_scanning:
    enabled: true
    severity_threshold: "medium"
    confidence_threshold: 0.7
    scan_history: true
    scan_stashes: true
    entropy_analysis: true
    
    # Custom patterns
    custom_patterns:
      - name: "Custom API Key"
        pattern: "myapi_[A-Za-z0-9]{32}"
        severity: "high"
        confidence: 0.9
    
    # Exclude patterns
    exclude_patterns:
      - "test/"
      - "*.test.js"
      - "docs/"
```

## Best Practices

### 1. Regular Scanning
- Run secret scans before every release
- Include in CI/CD pipelines
- Schedule weekly scans for active repositories

### 2. Team Education
- Train developers on secret management
- Establish clear policies for credential handling
- Regular security awareness sessions

### 3. Monitoring
- Set up alerts for secret detection
- Monitor for new secret patterns
- Track remediation progress

### 4. Documentation
- Document all secret management procedures
- Maintain incident response plans
- Keep remediation steps updated

## Troubleshooting

### High False Positive Rate
```bash
# Increase confidence threshold
git hc security secrets --confidence 0.9

# Disable entropy analysis
git hc security secrets --entropy=false

# Focus on high severity only
git hc security secrets --severity high
```

### Performance Issues
```bash
# Scan only recent history
git hc security secrets --history-depth 100

# Skip stashes
git hc security secrets --stashes=false
```

### Memory Issues
```bash
# Process in smaller chunks
git hc security secrets --batch-size 1000
```

## Related Features

- [Pre-commit Hooks](pre-commit-hooks.md) - Prevent secrets from being committed
- [Health Checks](health-checks.md) - Overall repository health assessment
- [CI/CD Integration](ci-cd-integration.md) - Automated security scanning

## Next Steps

1. [Pre-commit Hooks](pre-commit-hooks.md) - Set up pre-commit secret detection
2. [CI/CD Integration](ci-cd-integration.md) - Add to your CI/CD pipeline
3. [Health Checks](health-checks.md) - Run comprehensive health assessment
