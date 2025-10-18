# Contributing to Git Project Health Checker (GPHC)

Thank you for your interest in contributing to GPHC! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites
- Go 1.19 or higher
- Git
- Basic understanding of Go development

### Setting Up Development Environment

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/gphc.git
   cd gphc
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build the Project**
   ```bash
   go build -o gphc cmd/gphc/main.go
   ```

4. **Test Your Changes**
   ```bash
   ./gphc check
   ```

## 📋 How to Contribute

### Reporting Issues
- Use the [GitHub Issues](https://github.com/opsource/gphc/issues) page
- Include detailed description, steps to reproduce, and expected behavior
- Use appropriate labels (bug, enhancement, documentation, etc.)

### Suggesting Features
- Open a [GitHub Discussion](https://github.com/opsource/gphc/discussions) for feature requests
- Describe the use case and expected behavior
- Consider implementation complexity and maintainability

### Code Contributions

#### 1. Choose an Issue
- Look for issues labeled `good first issue` for beginners
- Comment on the issue to indicate you're working on it
- Ask questions if anything is unclear

#### 2. Create a Branch
```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number-description
```

#### 3. Make Changes
- Follow Go coding standards and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure all existing tests pass

#### 4. Commit Changes
```bash
git add .
git commit -m "feat: add new checker for X"
# or
git commit -m "fix: resolve issue with Y"
```

#### 5. Push and Create PR
```bash
git push origin feature/your-feature-name
```
Then create a Pull Request on GitHub.

## 🏗️ Project Structure

```
gphc/
├── cmd/gphc/           # CLI entry point
├── internal/
│   ├── checkers/       # Health check implementations
│   ├── git/           # Git repository analysis
│   ├── scorer/        # Scoring engine
│   └── reporter/      # Output formatting
├── pkg/
│   ├── types/         # Core data structures
│   └── config/        # Configuration management
├── tests/             # Test files
├── docs/              # Additional documentation
└── examples/          # Usage examples
```

## 🔧 Development Guidelines

### Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions and complex logic
- Keep functions focused and small
- Use interfaces for testability

### Testing
- Write unit tests for new functionality
- Aim for good test coverage
- Use table-driven tests where appropriate
- Test error conditions and edge cases

### Documentation
- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update configuration documentation
- Include examples for new features

## 🎯 Areas for Contribution

### New Checkers
Implement new health checks by:
1. Creating a new checker in `internal/checkers/`
2. Implementing the `Checker` interface
3. Adding appropriate scoring logic
4. Updating the main checker list

### Enhanced Git Analysis
- Improve branch analysis algorithms
- Add support for more Git features
- Optimize performance for large repositories

### Output Formats
- Add JSON/XML output support
- Create HTML reports
- Add integration with CI/CD systems

### Configuration
- Add more configuration options
- Support for custom rules
- Environment-specific settings

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment Information**
   - Operating System
   - Go version
   - GPHC version
   - Git version

2. **Steps to Reproduce**
   - Clear, numbered steps
   - Sample repository (if applicable)
   - Expected vs actual behavior

3. **Additional Context**
   - Error messages
   - Screenshots (if applicable)
   - Related issues

## 💡 Feature Requests

When suggesting features:

1. **Describe the Problem**
   - What problem does this solve?
   - Who would benefit from this feature?

2. **Propose a Solution**
   - How should this work?
   - Any implementation ideas?

3. **Consider Alternatives**
   - Are there existing solutions?
   - Could this be implemented differently?

## 🔍 Code Review Process

### For Contributors
- Respond to review feedback promptly
- Make requested changes or explain why they're not needed
- Test changes thoroughly
- Keep PRs focused and reasonably sized

### For Reviewers
- Be constructive and respectful
- Focus on code quality and correctness
- Consider maintainability and performance
- Approve when ready, request changes when needed

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Testing](https://golang.org/doc/tutorial/add-a-test)
- [Git Documentation](https://git-scm.com/doc)

## 🤝 Community Guidelines

- Be respectful and inclusive
- Help others learn and grow
- Share knowledge and best practices
- Follow the [Code of Conduct](CODE_OF_CONDUCT.md)

## 📞 Getting Help

- 💬 [GitHub Discussions](https://github.com/opsource/gphc/discussions)
- 🐛 [GitHub Issues](https://github.com/opsource/gphc/issues)
- 📧 Email: contributors@gphc.dev

## 🎉 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to GPHC! 🚀
