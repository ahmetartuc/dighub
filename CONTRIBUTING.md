# Contributing to Dighub

Thank you for your interest in contributing to Dighub! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Sample output** or error messages

### Suggesting Features

Feature suggestions are welcome! Please:

- **Check existing issues** for similar suggestions
- **Describe the feature** in detail
- **Explain the use case** and benefits
- **Consider implementation** complexity

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards
3. **Add tests** if applicable
4. **Update documentation** if needed
5. **Ensure tests pass** and code is formatted
6. **Submit a pull request** with a clear description

### Coding Standards

- Follow Go conventions and best practices
- Use `gofmt` to format your code
- Write clear, descriptive commit messages
- Add comments for complex logic
- Keep functions focused and small

### Adding New Dorks

To add new dork patterns:

1. Open `internal/dorks/dorks.go`
2. Add your dork to the appropriate priority level:

```go
{
    Pattern:     "filename:.env NEW_PATTERN",
    Priority:    PriorityHigh,  // or PriorityMedium, PriorityLow
    Category:    "CategoryName",
    Description: "Description of what this detects",
}
```

3. Test the dork pattern manually
4. Submit a PR with:
   - Description of the pattern
   - Why it's important
   - Example of what it detects

### Testing

```bash
# Run tests
go test ./...

# Build the project
go build

# Test your changes
./dighub -org test-org -token xxx
```

### Commit Message Format

Use clear, descriptive commit messages:

```
type: brief description

Longer explanation if needed

Fixes #issue-number
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/dighub.git
cd dighub

# Install dependencies
go mod download

# Build
go build

# Run
./dighub -h
```

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Questions?

Feel free to open an issue for any questions or concerns.

Thank you for contributing to Dighub! 🚀
