# Executable & Large File Audit

## Overview

Executable & Large File Audit is a comprehensive security feature that scans repositories for executable files, large files, and suspicious file types that pose security risks or repository health issues. This helps maintain repository hygiene and prevents security vulnerabilities.

## Why It's Important

Binary and large files in Git repositories can cause several problems:

- **Security Risks**: Executable files can contain malware or malicious code
- **Repository Bloat**: Large files make repositories slow and difficult to clone
- **Storage Issues**: Binary files consume significant storage space
- **Performance Impact**: Large files slow down Git operations
- **Compliance Issues**: Some organizations prohibit binary files in source control

## Key Features

### 1. Executable File Detection
- **Windows Executables**: `.exe`, `.dll`, `.msi`, `.pkg`
- **Unix/Linux Executables**: `.so`, `.dylib`, `.bin`
- **Java Archives**: `.jar`, `.war`, `.ear`
- **Package Files**: `.deb`, `.rpm`, `.dmg`, `.iso`
- **Virtual Machine Images**: `.vmdk`, `.vdi`, `.qcow2`, `.ova`

### 2. Large File Detection
- **Size Thresholds**: Configurable size limits (default: 10MB)
- **Severity Levels**: Different severity based on file size
- **Size Reporting**: Detailed size information in MB/GB
- **Git LFS Recommendations**: Suggests using Git LFS for large files

### 3. Suspicious File Detection
- **Script Files**: `.bat`, `.cmd`, `.vbs`, `.ps1`, `.js`
- **Office Documents**: `.doc`, `.docx`, `.xls`, `.xlsx`, `.ppt`, `.pptx`
- **System Files**: `.inf`, `.reg`, `.lnk`, `.scf`
- **Malware Patterns**: Files with suspicious names or patterns

### 4. Git History Analysis
- **Historical Scan**: Checks entire Git history for binary files
- **Commit Tracking**: Identifies when binary files were added
- **Cleanup Recommendations**: Suggests history cleanup methods

## Supported File Types

### Critical Severity
- **Windows Executables**: `.exe`, `.dll`, `.msi`
- **Java Archives**: `.jar`, `.war`, `.ear`
- **Binary Files**: `.bin`
- **Package Files**: `.deb`, `.rpm`, `.pkg`, `.dmg`

### High Severity
- **Libraries**: `.so`, `.dylib`
- **Disk Images**: `.iso`, `.img`, `.raw`
- **Virtual Machine Images**: `.vmdk`, `.vdi`, `.qcow2`
- **Application Bundles**: `.app`, `.ova`, `.ovf`

### Medium Severity
- **Script Files**: `.bat`, `.cmd`, `.vbs`, `.ps1`
- **Office Documents**: `.doc`, `.docx`, `.xls`, `.xlsx`
- **System Files**: `.inf`, `.reg`, `.lnk`

### Low Severity
- **Configuration Files**: `.cfg`, `.ini`
- **Log Files**: `.log`
- **Temporary Files**: `.tmp`, `.temp`

## Basic Usage

### Command Syntax
```bash
git hc security binaries [flags]
```

### Basic Binary Audit
```bash
# Run complete binary audit
git hc security binaries

# Audit specific repository
git hc security binaries /path/to/repo
```

### Size Threshold Configuration
```bash
# Set custom size threshold
git hc security binaries --max-size 50mb

# Different size formats
git hc security binaries --max-size 100mb
git hc security binaries --max-size 1gb
git hc security binaries --max-size 500kb
```

### Focused Checks
```bash
# Check only executable files
git hc security binaries --check-executables --check-large=false --check-suspicious=false

# Check only large files
git hc security binaries --check-executables=false --check-large --check-suspicious=false

# Check only suspicious files
git hc security binaries --check-executables=false --check-large=false --check-suspicious
```

## Command Options

### Check Configuration
- `--check-executables`: Check for executable files (default: true)
- `--check-large`: Check for large files (default: true)
- `--check-suspicious`: Check for suspicious file types (default: true)
- `--check-history`: Check Git history for binary files (default: true)

### Size Configuration
- `--max-size string`: Maximum file size threshold (e.g., 10mb, 50mb, 100mb) (default: "10mb")

### Output Options
- `--severity string`: Minimum severity level (low, medium, high, critical) (default: "low")
- `--format string`: Output format (table, json, yaml) (default: "table")
- `--output string`: Output file path for results

## Example Output

### Table Format
```
🔍 Auditing executable and large files...
Repository: /path/to/project
Max size threshold: 50mb (50.0 MB)
Check history: true
Check executables: true
Check large files: true
Check suspicious files: true
Minimum severity: low

❌ Binary audit found issues: Found 3 binary/large file issues
📊 Binary File Audit Results
============================

Executable Files: 1
Large Files: 2
Suspicious Files: 0
Total Size: 125.4 MB
File Count: 3
Security Score: 65/100

🚨 BINARY FILE ISSUES FOUND!

Immediate Actions Required:
1. Review and remove unnecessary binary files
2. Use Git LFS for large files
3. Add binary file patterns to .gitignore
4. Remove suspicious files from repository
5. Clean up Git history if needed

Best Practices:
- Use Git LFS for files > 100MB
- Avoid committing executable files
- Use .gitignore for binary patterns
- Regular binary file audits
- Use package managers for dependencies
```

### JSON Format
```bash
git hc security binaries --format json
```

```json
{
  "id": "BINARY-AUDIT",
  "name": "Executable & Large File Audit",
  "status": 1,
  "score": 65,
  "message": "Found 3 binary/large file issues",
  "details": [
    "Executable Files: 1",
    "Large Files: 2",
    "Suspicious Files: 0",
    "Total Size: 125.4 MB",
    "File Count: 3"
  ],
  "category": 0,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

## Size Thresholds

### Default Thresholds
- **10MB**: Default threshold for large file detection
- **50MB**: Recommended threshold for most projects
- **100MB**: Critical threshold requiring Git LFS

### Severity Levels by Size
- **< 10MB**: Low severity
- **10-20MB**: Medium severity
- **20-50MB**: High severity
- **> 50MB**: Critical severity

## Executable File Types

### Windows Executables
```bash
.exe    # Windows executable
.dll    # Dynamic link library
.msi    # Windows installer
.pkg    # macOS package
```

### Unix/Linux Executables
```bash
.so     # Shared object library
.dylib  # macOS dynamic library
.bin    # Binary executable
```

### Java Archives
```bash
.jar    # Java archive
.war    # Web application archive
.ear    # Enterprise application archive
```

### Package Files
```bash
.deb    # Debian package
.rpm    # RPM package
.dmg    # macOS disk image
.iso    # ISO disk image
```

### Virtual Machine Images
```bash
.vmdk   # VMware disk image
.vdi    # VirtualBox disk image
.qcow2  # QEMU disk image
.ova    # Open Virtual Appliance
.ovf    # Open Virtualization Format
```

## Suspicious File Patterns

### Script Files
```bash
.bat    # Windows batch file
.cmd    # Windows command file
.vbs    # VBScript file
.ps1    # PowerShell script
.js     # JavaScript file
```

### Office Documents
```bash
.doc    # Microsoft Word document
.docx   # Microsoft Word document (new format)
.xls    # Microsoft Excel spreadsheet
.xlsx   # Microsoft Excel spreadsheet (new format)
.ppt    # Microsoft PowerPoint presentation
.pptx   # Microsoft PowerPoint presentation (new format)
```

### System Files
```bash
.inf    # Windows information file
.reg    # Windows registry file
.lnk    # Windows shortcut
.scf    # Windows shell command file
```

### Malware Patterns
Files containing these patterns in their names:
- `malware`, `virus`, `trojan`, `backdoor`
- `keylogger`, `rootkit`, `spyware`, `adware`
- `ransomware`, `botnet`, `exploit`, `payload`
- `inject`, `bypass`, `crack`, `keygen`
- `patch`, `hack`, `cracked`, `pirated`

## CI/CD Integration

### GitHub Actions
```yaml
name: Binary File Audit
on: [push, pull_request]

jobs:
  binary-audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup GPHC
        run: go install github.com/opsource/gphc/cmd/gphc@latest
      - name: Audit Binary Files
        run: |
          git hc security binaries --max-size 50mb --format json --output binary-audit.json
      - name: Upload Results
        uses: actions/upload-artifact@v3
        with:
          name: binary-audit-results
          path: binary-audit.json
```

### GitLab CI
```yaml
binary_file_audit:
  stage: security
  image: golang:1.19
  before_script:
    - go install github.com/opsource/gphc/cmd/gphc@latest
  script:
    - git hc security binaries --max-size 100mb --format json --output binary-audit.json
  artifacts:
    reports:
      junit: binary-audit.json
  only:
    - merge_requests
    - main
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    stages {
        stage('Binary File Audit') {
            steps {
                sh 'go install github.com/opsource/gphc/cmd/gphc@latest'
                sh 'git hc security binaries --max-size 50mb --format json --output binary-audit.json'
                archiveArtifacts artifacts: 'binary-audit.json'
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
                reportFiles: 'binary-audit.json',
                reportName: 'Binary File Audit Report'
            ])
        }
    }
}
```

## Configuration

### gphc.yml Configuration
```yaml
# Binary file audit configuration
binary_file_audit:
  enabled: true
  severity_threshold: "medium"
  
  # Size thresholds
  size_thresholds:
    large_file: "10mb"
    critical_file: "50mb"
    git_lfs_threshold: "100mb"
  
  # File type checks
  checks:
    executables: true
    large_files: true
    suspicious_files: true
    git_history: true
  
  # Executable file patterns
  executable_patterns:
    - "*.exe"
    - "*.dll"
    - "*.so"
    - "*.dylib"
    - "*.bin"
    - "*.jar"
    - "*.war"
    - "*.ear"
    - "*.app"
    - "*.deb"
    - "*.rpm"
    - "*.msi"
    - "*.pkg"
    - "*.dmg"
    - "*.iso"
    - "*.img"
    - "*.raw"
    - "*.vmdk"
    - "*.vdi"
    - "*.qcow2"
    - "*.ova"
    - "*.ovf"
  
  # Suspicious file patterns
  suspicious_patterns:
    - "*.scr"
    - "*.bat"
    - "*.cmd"
    - "*.com"
    - "*.pif"
    - "*.vbs"
    - "*.js"
    - "*.jse"
    - "*.wsf"
    - "*.wsh"
    - "*.ps1"
    - "*.psm1"
    - "*.psd1"
    - "*.ps1xml"
    - "*.psc1"
    - "*.msh"
    - "*.msh1"
    - "*.msh2"
    - "*.mshxml"
    - "*.msh1xml"
    - "*.msh2xml"
    - "*.scf"
    - "*.lnk"
    - "*.inf"
    - "*.reg"
    - "*.doc"
    - "*.docx"
    - "*.xls"
    - "*.xlsx"
    - "*.ppt"
    - "*.pptx"
```

## Best Practices

### Repository Hygiene
1. **Use Git LFS**: For files larger than 100MB
2. **Avoid Binary Files**: Use package managers instead
3. **Comprehensive .gitignore**: Include binary file patterns
4. **Regular Audits**: Run binary audits regularly
5. **Clean History**: Remove binary files from Git history

### Git LFS Setup
```bash
# Install Git LFS
git lfs install

# Track large files
git lfs track "*.zip"
git lfs track "*.tar.gz"
git lfs track "*.iso"
git lfs track "*.dmg"

# Add .gitattributes
git add .gitattributes
git commit -m "Add Git LFS tracking"
```

### .gitignore Patterns
```gitignore
# Executable files
*.exe
*.dll
*.so
*.dylib
*.bin
*.jar
*.war
*.ear
*.app
*.deb
*.rpm
*.msi
*.pkg
*.dmg
*.iso
*.img
*.raw
*.vmdk
*.vdi
*.qcow2
*.ova
*.ovf

# Large files
*.zip
*.tar.gz
*.rar
*.7z
*.bz2
*.gz

# Suspicious files
*.scr
*.bat
*.cmd
*.com
*.pif
*.vbs
*.js
*.jse
*.wsf
*.wsh
*.ps1
*.psm1
*.psd1
*.ps1xml
*.psc1
*.msh
*.msh1
*.msh2
*.mshxml
*.msh1xml
*.msh2xml
*.scf
*.lnk
*.inf
*.reg
*.doc
*.docx
*.xls
*.xlsx
*.ppt
*.pptx
```

## Troubleshooting

### Common Issues

#### Large Files Found
```
Large Files: 2
Total Size: 125.4 MB
```
**Solution**: Use Git LFS or remove unnecessary large files.

#### Executable Files Found
```
Executable Files: 1
```
**Solution**: Remove executable files and use package managers.

#### Suspicious Files Found
```
Suspicious Files: 1
```
**Solution**: Review and remove suspicious files, update .gitignore.

#### Files in Git History
```
File found in Git history: malicious.exe
```
**Solution**: Use `git filter-repo` or BFG to clean history.

### Debug Mode
```bash
# Enable verbose output
git hc security binaries --verbose

# Check specific file types
git hc security binaries --check-executables --check-large=false
```

## Integration with Other Tools

### Pre-commit Hooks
```bash
# Add to .pre-commit-config.yaml
- repo: local
  hooks:
    - id: binary-file-audit
      name: Binary File Audit
      entry: git hc security binaries --max-size 10mb
      language: system
      pass_filenames: false
```

### Git Hooks
```bash
#!/bin/bash
# .git/hooks/pre-push
git hc security binaries --max-size 50mb
if [ $? -ne 0 ]; then
    echo "Binary file audit failed. Push aborted."
    exit 1
fi
```

## Next Steps

1. [Git Policy Validation](git-policy-validation.md) - Validate Git security policies
2. [Secret Scanning](secret-scanning.md) - Scan for exposed secrets and credentials
3. [Transitive Dependency Vetting](transitive-dependency-vetting.md) - Analyze dependency vulnerabilities
4. [Pre-commit Hooks](pre-commit-hooks.md) - Set up binary file checks in pre-commit hooks
5. [CI/CD Integration](ci-cd-integration.md) - Add to your CI/CD pipeline
