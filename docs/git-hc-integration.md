# Git HC Integration Guide

This guide covers using GPHC as a Git subcommand (`git hc`) for seamless integration with Git workflows.

## Overview

Git HC integration allows you to use GPHC as a native Git subcommand, making it feel like a built-in Git feature. This provides a more integrated experience for developers who are already familiar with Git commands.

## Setup

### Automatic Setup
```bash
# Run the setup script
./setup-git-hc.sh
```

### Manual Setup
```bash
# 1. Create wrapper script directory
mkdir -p ~/.local/bin

# 2. Create wrapper script
cat > ~/.local/bin/git-hc-wrapper << 'EOF'
#!/bin/bash
# Find and execute gphc binary
if command -v gphc >/dev/null 2>&1; then
    exec gphc "$@"
else
    echo "Error: GPHC not found. Please install it first."
    exit 1
fi
EOF

# 3. Make wrapper executable
chmod +x ~/.local/bin/git-hc-wrapper

# 4. Add to PATH (if not already there)
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc
source ~/.bashrc

# 5. Set up git alias
git config --global alias.hc "!$HOME/.local/bin/git-hc-wrapper"
```

## Usage

### Basic Commands
```bash
# Health check
git hc check

# Pre-commit checks
git hc pre-commit

# Interactive terminal UI
git hc tui

# Web dashboard
git hc serve

# Multi-repository scan
git hc scan ~/projects --recursive

# Update GPHC
git hc update

# Version information
git hc version

# Help
git hc --help
```

### Advanced Usage
```bash
# Check with specific options
git hc check --min-score 80 --format json

# Scan with filters
git hc scan ~/projects --recursive --min-score 70 --exclude "node_modules"

# Serve with custom port
git hc serve --port 3000

# Pre-commit with verbose output
git hc pre-commit --verbose
```

## Benefits

### Seamless Integration
- **Native Feel**: Feels like a built-in Git command
- **Consistent Interface**: Follows Git command patterns
- **Context Awareness**: Automatically detects Git repository
- **Error Handling**: Proper Git-style error messages

### Workflow Integration
- **Pre-commit Hooks**: Easy integration with Git hooks
- **CI/CD Pipelines**: Consistent command interface
- **Team Collaboration**: Standardized commands across team
- **Documentation**: Clear command structure

## Configuration

### Git Configuration
```bash
# View current alias
git config --global alias.hc

# Remove alias
git config --global --unset alias.hc

# Update alias
git config --global alias.hc "!$HOME/.local/bin/git-hc-wrapper"
```

### Wrapper Script Customization
```bash
# Edit wrapper script
nano ~/.local/bin/git-hc-wrapper

# Add custom logic
#!/bin/bash
# Custom wrapper with additional features

# Check for specific conditions
if [ "$1" = "check" ] && [ ! -f "README.md" ]; then
    echo "Warning: No README.md found in repository"
fi

# Execute gphc
exec gphc "$@"
```

## Troubleshooting

### Common Issues

#### Command Not Found
```bash
# Check if alias is set
git config --global alias.hc

# Check if wrapper script exists
ls -la ~/.local/bin/git-hc-wrapper

# Check PATH
echo $PATH | grep -q "$HOME/.local/bin" && echo "PATH OK" || echo "PATH missing"
```

#### Permission Denied
```bash
# Fix wrapper script permissions
chmod +x ~/.local/bin/git-hc-wrapper

# Check file ownership
ls -la ~/.local/bin/git-hc-wrapper
```

#### GPHC Not Found
```bash
# Check if gphc is installed
which gphc

# Install if missing
go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Debugging
```bash
# Test wrapper script directly
~/.local/bin/git-hc-wrapper version

# Test git alias
git hc version

# Verbose output
git hc check --verbose
```

## Integration Examples

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit
git hc pre-commit
```

### GitHub Actions
```yaml
name: Health Check
on: [push, pull_request]
jobs:
  health-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
      - name: Install GPHC
        run: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
      - name: Setup Git HC
        run: ./setup-git-hc.sh
      - name: Health Check
        run: git hc check --min-score 80
```

### GitLab CI
```yaml
health_check:
  stage: test
  script:
    - go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest
    - ./setup-git-hc.sh
    - git hc check --min-score 80
```

## Best Practices

### For Teams
1. **Standardized Setup**: Use setup script for consistent installation
2. **Documentation**: Document Git HC usage in team guidelines
3. **Training**: Train team members on Git HC commands
4. **Integration**: Integrate with existing Git workflows
5. **Monitoring**: Monitor Git HC usage across team

### For Organizations
1. **Centralized Configuration**: Use organization-wide Git configuration
2. **Security Policies**: Implement security policies for Git aliases
3. **Compliance**: Ensure compliance with organizational standards
4. **Training Programs**: Implement training programs for Git HC
5. **Support**: Provide support for Git HC integration

## Migration

### From GPHC to Git HC
```bash
# 1. Install GPHC (if not already installed)
go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest

# 2. Run setup script
./setup-git-hc.sh

# 3. Test Git HC
git hc version

# 4. Update scripts and documentation
# Replace 'gphc' with 'git hc' in scripts
```

## Next Steps
- [Basic Usage](basic-usage.md) - Getting started with GPHC
- [Pre-commit Hooks](pre-commit-hooks.md) - Pre-commit integration guide
- [CI/CD Integration](ci-cd-integration.md) - Pipeline integration guide
