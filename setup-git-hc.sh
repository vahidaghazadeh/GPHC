#!/bin/bash
# Git HC Setup Script
# This script sets up GPHC as 'git hc' subcommand using git aliases

set -e

echo "Setting up GPHC as 'git hc' subcommand..."

# Find GPHC binary
GPHC_BINARY=""

# Check if gphc is in PATH
if command -v gphc >/dev/null 2>&1; then
    GPHC_BINARY="gphc"
# Check common installation paths
elif [ -f "$HOME/go/bin/gphc" ]; then
    GPHC_BINARY="$HOME/go/bin/gphc"
elif [ -f "/usr/local/bin/gphc" ]; then
    GPHC_BINARY="/usr/local/bin/gphc"
elif [ -f "/usr/bin/gphc" ]; then
    GPHC_BINARY="/usr/bin/gphc"
# Check GOPATH/bin
elif [ -n "$GOPATH" ] && [ -f "$GOPATH/bin/gphc" ]; then
    GPHC_BINARY="$GOPATH/bin/gphc"
# Check if go is available and try to find gphc
elif command -v go >/dev/null 2>&1; then
    GO_BIN_PATH="$(go env GOPATH)/bin"
    if [ -f "$GO_BIN_PATH/gphc" ]; then
        GPHC_BINARY="$GO_BIN_PATH/gphc"
    fi
fi

# If gphc is not found, show error and installation instructions
if [ -z "$GPHC_BINARY" ]; then
    echo "Error: GPHC binary not found!"
    echo ""
    echo "Please install GPHC first:"
    echo "  go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest"
    echo ""
    echo "Then run this setup script again."
    exit 1
fi

echo "Found GPHC binary at: $GPHC_BINARY"

# Create wrapper script for git hc
WRAPPER_SCRIPT="$HOME/.local/bin/git-hc-wrapper"

# Create directory if it doesn't exist
mkdir -p "$(dirname "$WRAPPER_SCRIPT")"

# Create the wrapper script
cat > "$WRAPPER_SCRIPT" << 'EOF'
#!/bin/bash
# Git HC Wrapper Script
# This script is called by git alias 'hc'

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Try to find gphc binary
GPHC_BINARY=""

# Check if gphc is in PATH
if command -v gphc >/dev/null 2>&1; then
    GPHC_BINARY="gphc"
# Check common installation paths
elif [ -f "$HOME/go/bin/gphc" ]; then
    GPHC_BINARY="$HOME/go/bin/gphc"
elif [ -f "/usr/local/bin/gphc" ]; then
    GPHC_BINARY="/usr/local/bin/gphc"
elif [ -f "/usr/bin/gphc" ]; then
    GPHC_BINARY="/usr/bin/gphc"
# Check GOPATH/bin
elif [ -n "$GOPATH" ] && [ -f "$GOPATH/bin/gphc" ]; then
    GPHC_BINARY="$GOPATH/bin/gphc"
# Check if go is available and try to find gphc
elif command -v go >/dev/null 2>&1; then
    GO_BIN_PATH="$(go env GOPATH)/bin"
    if [ -f "$GO_BIN_PATH/gphc" ]; then
        GPHC_BINARY="$GO_BIN_PATH/gphc"
    fi
fi

# If gphc is not found, show error
if [ -z "$GPHC_BINARY" ]; then
    echo "Error: GPHC binary not found!"
    echo "Please install GPHC: go install github.com/vahidaghazadeh/gphc/cmd/gphc@latest"
    exit 1
fi

# Check if we're in a git repository
if ! git rev-parse --git-dir >/dev/null 2>&1; then
    echo "Error: Not in a git repository!"
    echo "Please run 'git hc' from within a git repository."
    exit 1
fi

# Execute gphc with all passed arguments
exec "$GPHC_BINARY" "$@"
EOF

# Make wrapper script executable
chmod +x "$WRAPPER_SCRIPT"

echo "Created wrapper script at: $WRAPPER_SCRIPT"

# Add ~/.local/bin to PATH if not already there
if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
    echo ""
    echo "Adding ~/.local/bin to PATH..."
    
    # Detect shell and add to appropriate profile
    if [ -n "$ZSH_VERSION" ]; then
        PROFILE_FILE="$HOME/.zshrc"
    elif [ -n "$BASH_VERSION" ]; then
        PROFILE_FILE="$HOME/.bashrc"
    else
        PROFILE_FILE="$HOME/.profile"
    fi
    
    echo "export PATH=\"\$PATH:\$HOME/.local/bin\"" >> "$PROFILE_FILE"
    echo "Added PATH export to: $PROFILE_FILE"
    echo "Please run: source $PROFILE_FILE"
fi

# Set up git alias
echo ""
echo "Setting up git alias..."

# Set the git alias to use our wrapper script
git config --global alias.hc "!$WRAPPER_SCRIPT"

echo "Git alias 'hc' configured successfully!"

# Test the setup
echo ""
echo "Testing setup..."
if git hc version >/dev/null 2>&1; then
    echo "✅ Setup successful! You can now use 'git hc' commands."
    echo ""
    echo "Available commands:"
    echo "  git hc check          - Check repository health"
    echo "  git hc pre-commit     - Run pre-commit checks"
    echo "  git hc tui           - Launch interactive terminal UI"
    echo "  git hc serve         - Start web dashboard"
    echo "  git hc scan          - Scan multiple repositories"
    echo "  git hc update        - Update GPHC"
    echo "  git hc version       - Show version information"
    echo "  git hc --help        - Show help"
else
    echo "❌ Setup failed. Please check the error messages above."
    exit 1
fi

echo ""
echo "Setup complete! 🎉"
echo "You can now use 'git hc' instead of 'gphc' for all commands."
