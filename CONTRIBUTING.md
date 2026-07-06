# Contributing to Likhis

Thank you for your interest in contributing to Likhis! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Prerequisites](#prerequisites)
- [Project Structure](#project-structure)
- [Makefile Reference](#makefile-reference)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Commit Conventions](#commit-conventions)
- [PR Process](#pr-process)
- [Release Process](#release-process)
- [Questions](#questions)

## Code of Conduct

This project adheres to a [Code of Conduct](CODE_OF_CONDUCT.md) that all contributors are expected to follow. Please be respectful, inclusive, and professional in all interactions.

## Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.21+ | Compile and run the project |
| Make | Any | Build automation and common tasks |
| golangci-lint | Latest | Code linting (optional, for `make lint`) |

## Project Structure

```
likhis/
├── main.go                  # CLI entry point, flag parsing, orchestration
├── internal/
│   ├── traversal/           # BFS file traversal, dependency skipping
│   ├── parser/              # Route parsing (Express, Flask, Django, Spring, Laravel + plugin engine)
│   ├── exporters/           # Output generators (Postman, Insomnia, HTTPie, cURL, OpenAPI)
│   └── plugins/             # YAML plugin loader and pattern matcher
├── plugins/                 # Built-in YAML plugin definitions
├── exp/                     # Example projects for integration testing
├── build/                   # Compiled binaries (gitignored)
├── docs/                    # Documentation site assets (Cloudflare Pages)
├── tests/                   # Unit test files
└── .github/workflows/       # CI/CD pipeline definitions
```

## Makefile Reference

| Command | Description |
|---|---|
| `make all` | Clean and build (default) |
| `make build` | Build the binary into `build/` |
| `make clean` | Remove build artifacts and coverage files |
| `make test` | Run unit tests |
| `make test-integration` | Build + test against example projects in `exp/` |
| `make coverage` | Run tests with coverage report |
| `make lint` | Run golangci-lint |
| `make run` | Build and run on the current directory |
| `make release` | Cross-compile for Windows, Linux, and macOS (amd64 + arm64) |
| `make link` | Create symbolic link to PATH for local development |
| `make help` | Show all available targets |

### Examples

```bash
# Quick build and test cycle
make build && make test

# Full integration test suite
make test-integration

# Build release binaries for all platforms
make release
```

## Development Workflow

1. **Fork the repository** and clone your fork:

```bash
git clone https://github.com/your-username/likhis.git
cd likhis
```

2. **Create a feature branch:**

```bash
git checkout -b feat/your-feature-name
```

3. **Make your changes:**
   - Write clean, readable code following the coding standards
   - Add tests for new functionality
   - Update documentation as needed

4. **Test your changes:**

```bash
make test                  # Unit tests
make test-integration      # Integration tests against example projects
make lint                  # Code linting
```

5. **Commit your changes** using [conventional commits](#commit-conventions):

```bash
git commit -m "feat(parser): add support for Fastify framework"
```

6. **Push and create a Pull Request:**

```bash
git push origin feat/your-feature-name
```

## Coding Standards

### Go Style

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` to format all Go code
- Keep functions focused and small (prefer under 50 lines)
- Use meaningful variable and function names
- Document all exported functions and types

### Code Organization

- **Separation of concerns**: Each package has a single responsibility
- **Dependency injection**: Pass dependencies explicitly, avoid global state
- **Interface segregation**: Keep interfaces small and focused
- **Error handling**: Always handle errors explicitly; use `fmt.Errorf` with `%w` for wrapping

### File Naming

- Packages: lowercase, single word when possible
- Files: lowercase with underscores for multiple words
- Exported functions: PascalCase
- Unexported functions: camelCase
- Constants: PascalCase or UPPER_SNAKE_CASE

## Testing

### Running Tests

```bash
make test                # Run all unit tests
make test-integration    # Run integration tests against examples in exp/
make coverage            # Run tests with coverage report
```

### Writing Tests

- Test edge cases and error conditions
- Test with example projects in `exp/`
- Verify output formats are correct and valid JSON/YAML
- Aim for >80% coverage on core packages (parser, exporters)

### Test Patterns

**Unit test example:**

```go
func TestParseExpressRoute(t *testing.T) {
    parser := NewRouteParser("express")
    routes, err := parser.ParseFile("testdata/express/routes.js")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if len(routes) == 0 {
        t.Error("expected at least one route")
    }
}
```

### Integration Test Example

```bash
./build/likhis -p exp/express -o postman -F express -O out/express
# Verify: out/express/postman-collection.json exists and is valid JSON
```

## Commit Conventions

Likhis uses [Conventional Commits](https://www.conventionalcommits.org/) for all commit messages. This allows automatic changelog generation and semantic versioning.

### Format

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

| Type | Usage |
|---|---|
| `feat` | A new feature |
| `fix` | A bug fix |
| `docs` | Documentation changes |
| `refactor` | Code refactoring (no functional change) |
| `test` | Adding or updating tests |
| `chore` | Maintenance, build, CI, dependencies |
| `perf` | Performance improvements |
| `style` | Code style changes (formatting, etc.) |
| `ci` | CI/CD configuration changes |

### Scope

The scope should be the package or area affected:

- `parser` — route parsing logic
- `exporters` — output format generators
- `plugins` — plugin system
- `traversal` — file traversal
- `cli` — command-line interface
- `docs` — documentation

### Examples

```
feat(parser): add support for Fastify framework
fix(exporters): handle empty route list in Postman export
docs(readme): add CI/CD integration section
refactor(plugins): extract pattern compilation to separate function
test(parser): add edge cases for Spring Boot path parameters
chore(build): update Go version to 1.22
```

### Breaking Changes

Add `BREAKING CHANGE:` in the footer or append `!` after the type/scope:

```
feat!(api): change output format from v1 to v2

BREAKING CHANGE: Output file structure has changed.
```

## PR Process

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] Tests pass locally (`make test`)
- [ ] Integration tests pass (`make test-integration`)
- [ ] Lint passes (`make lint`)
- [ ] Documentation updated (README, user guide, or plugin docs)
- [ ] Commit messages follow conventional commits
- [ ] No merge conflicts with main branch
- [ ] Changes tested with example projects in `exp/`

### Pull Request Checklist

Copy and paste this into your PR description:

```markdown
## Description

[Describe your changes and which issue they fix]

Fixes #(issue)

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation
- [ ] Refactoring
- [ ] Tests

## Checklist

- [ ] Code follows project style
- [ ] Tests pass
- [ ] Integration tests pass
- [ ] Lint passes
- [ ] Documentation updated
- [ ] Conventional commits used
```

### Review Process

1. A maintainer will review your PR within a few days
2. Address feedback by pushing additional commits
3. Once approved, a maintainer will merge your PR

### What Gets Merged

- Features that follow the project's architecture and design principles
- Bug fixes with clear reproduction steps
- Documentation improvements
- Plugin additions with example projects

### What Doesn't Get Merged

- Changes that break existing functionality without clear migration path
- Unrelated or cosmetic changes mixed with functional changes
- Code that doesn't include tests

## Release Process

1. **Update CHANGELOG.md** — ensure all changes since last release are documented

2. **Run full test suite:**

```bash
make test && make test-integration
```

3. **Build release binaries:**

```bash
make release
```

4. **Create a git tag:**

```bash
git tag -a v1.2.0 -m "Release version 1.2.0"
git push origin v1.2.0
```

5. **Create a GitHub Release:**
   - Upload the binaries from `release/`
   - Add release notes summarizing changes
   - Highlight new features and breaking changes

Version numbers follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

## Questions?

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/marcuwynu23/likhis/issues)
- **Discussions**: Ask questions and share ideas in [GitHub Discussions](https://github.com/marcuwynu23/likhis/discussions)
- **Existing issues**: Check open/closed issues before creating new ones

Thank you for contributing to Likhis!
