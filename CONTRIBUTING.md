# Contributing to Likhis

Thank you for your interest in contributing to Likhis! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Creating Plugins](#creating-plugins)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow. Please be respectful, inclusive, and professional in all interactions.

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please open an issue with:
- A clear, descriptive title
- Steps to reproduce the issue
- Expected vs. actual behavior
- Your environment (OS, Go version, etc.)
- Any relevant error messages or logs

### Suggesting Enhancements

Feature requests are welcome! Please include:
- A clear description of the feature
- Use cases and examples
- Potential implementation approach (if you have ideas)

### Contributing Code

- Fix bugs
- Implement new features
- Improve documentation
- Add support for new frameworks via plugins
- Optimize performance
- Improve test coverage

## Development Setup

1. **Fork and clone the repository**:
   ```bash
   git clone https://github.com/your-username/likhis.git
   cd likhis
   ```

2. **Ensure you have Go 1.21+ installed**:
   ```bash
   go version
   ```

3. **Build the project**:
   ```bash
   # Windows
   scripts\build.bat
   
   # Unix/Linux/macOS
   ./scripts/build.sh
   ```

4. **Run tests**:
   ```bash
   # Windows
   scripts\test.bat
   
   # Unix/Linux/macOS
   ./scripts/test.sh
   ```

## Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**:
   - Write clean, readable code
   - Follow the coding standards
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   # Run the test suite
   scripts\test.bat  # or ./scripts/test.sh
   
   # Test manually with example projects
   ./build/likhis -p exp/express -o postman -F express
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add support for new framework"
   ```
   
   Use conventional commit messages:
   - `feat:` for new features
   - `fix:` for bug fixes
   - `docs:` for documentation
   - `refactor:` for code refactoring
   - `test:` for tests
   - `chore:` for maintenance

5. **Push and create a Pull Request**:
   ```bash
   git push origin feature/your-feature-name
   ```

## Coding Standards

### Go Style Guide

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format code
- Keep functions focused and small
- Add comments for exported functions and types
- Use meaningful variable and function names

### Code Organization

- Keep related code together
- Separate concerns (parsing, exporting, traversal)
- Use interfaces where appropriate
- Avoid deep nesting

### Error Handling

- Always handle errors explicitly
- Return errors from functions that can fail
- Provide meaningful error messages
- Use `fmt.Errorf` with `%w` for error wrapping when appropriate

## Testing

### Running Tests

```bash
# Run all tests
scripts\test.bat  # Windows
./scripts/test.sh  # Unix/Linux/macOS

# Test specific framework
./build/likhis -p exp/express -o postman -F express
```

### Writing Tests

- Test edge cases and error conditions
- Test with the example projects in `exp/`
- Verify output formats are correct
- Test across different frameworks

### Test Coverage

Aim for good test coverage, especially for:
- Route parsing logic
- Export format generation
- Plugin system
- Edge cases and error handling

## Submitting Changes

### Pull Request Process

1. **Update your branch**:
   ```bash
   git checkout main
   git pull upstream main
   git checkout your-branch
   git rebase main
   ```

2. **Ensure all tests pass**:
   ```bash
   scripts\test.bat  # or ./scripts/test.sh
   ```

3. **Create a descriptive PR**:
   - Clear title and description
   - Reference related issues
   - Include screenshots/examples if applicable
   - List any breaking changes

4. **Respond to feedback**:
   - Address review comments
   - Make requested changes
   - Keep discussions constructive

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated (if needed)
- [ ] Commit messages follow conventions
- [ ] No merge conflicts
- [ ] Changes tested with example projects

## Creating Plugins

### Plugin Structure

See [GUIDELINES.md](GUIDELINES.md) for detailed plugin creation guidelines.

### Plugin Testing

1. Create your plugin YAML file in `plugins/`
2. Test with a sample project:
   ```bash
   ./build/likhis -p /path/to/project -o postman -F your-framework
   ```
3. Verify routes are detected correctly
4. Test with different output formats

### Submitting Plugins

When submitting a new framework plugin:
- Include the plugin YAML file
- Add example project in `exp/` (if possible)
- Update documentation
- Test with multiple route patterns

## Questions?

- Open an issue for questions or discussions
- Check existing issues and PRs first
- Be patient and respectful

Thank you for contributing to Likhis! 🎉

