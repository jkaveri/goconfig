# Contributing to goconfig

Thank you for your interest in contributing to goconfig! This document provides guidelines and steps for contributing to the project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and constructive in your interactions with other contributors.

## How to Contribute

### Reporting Issues

Before creating an issue, please:
1. Check if the issue has already been reported
2. Provide a clear and descriptive title
3. Include as much relevant information as possible
4. Add steps to reproduce the issue
5. Include your Go version and operating system
6. Add any relevant code snippets

### Pull Requests

1. Fork the repository
2. Create a new branch for your feature/fix (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation if needed
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/jkaveri/goconfig.git
   cd goconfig
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   go test ./...
   ```

### Code Style

- Follow standard Go formatting (`go fmt`)
- Follow Go best practices and idioms
- Write clear and concise code
- Add comments for complex logic
- Keep functions focused and small
- Use meaningful variable and function names

### Testing

- Write unit tests for new functionality
- Ensure existing tests pass
- Test edge cases and error conditions
- Use table-driven tests where appropriate

### Documentation

- Update README.md if needed
- Add comments for exported functions and types
- Include examples in documentation
- Keep documentation up to date with code changes

## Questions?

If you have any questions, feel free to:
1. Open an issue
2. Start a discussion
3. Contact the maintainers

Thank you for contributing to goconfig!
