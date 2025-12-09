# Likhis Development Guidelines

This document provides comprehensive guidelines for developing, extending, and maintaining the Likhis project.

## Table of Contents

- [Project Overview](#project-overview)
- [Architecture](#architecture)
- [Code Structure](#code-structure)
- [Plugin System](#plugin-system)
- [Adding New Frameworks](#adding-new-frameworks)
- [Export Formats](#export-formats)
- [Testing Guidelines](#testing-guidelines)
- [Performance Considerations](#performance-considerations)
- [Documentation Standards](#documentation-standards)
- [Release Process](#release-process)

## Project Overview

Likhis is a Go-based CLI tool that:
- Scans backend source code using BFS traversal
- Detects API routes using regex pattern matching
- Extracts route metadata (methods, parameters, queries)
- Generates exports for various API testing tools
- Supports extensible framework plugins via YAML

## Architecture

### Core Components

1. **Traversal** (`internal/traversal/`)
   - Breadth-First Search file traversal
   - Skips dependency directories
   - Filters files by extension

2. **Parser** (`internal/parser/`)
   - Framework-specific route parsing
   - Plugin-based pattern matching
   - Router mounting detection
   - Parameter extraction

3. **Exporters** (`internal/exporters/`)
   - Format-specific export generation
   - Environment variable handling
   - Request structure creation

4. **Plugins** (`internal/plugins/`)
   - YAML plugin loading
   - Pattern matching engine
   - Framework detection

### Data Flow

```
Source Code → Traversal → Parser → Route Objects → Exporters → Output Files
```

## Code Structure

### Directory Layout

```
likhis/
├── main.go                 # CLI entry point
├── internal/
│   ├── traversal/         # File traversal logic
│   ├── parser/            # Route parsing logic
│   ├── exporters/         # Export format generators
│   └── plugins/           # Plugin system
├── plugins/               # YAML plugin definitions
├── exp/                   # Example projects for testing
├── scripts/               # Build and test scripts
└── build/                 # Compiled executables
```

### Naming Conventions

- **Packages**: lowercase, single word when possible
- **Files**: lowercase with underscores for multiple words
- **Exported functions**: PascalCase
- **Unexported functions**: camelCase
- **Constants**: PascalCase or UPPER_SNAKE_CASE

### Code Organization Principles

1. **Separation of Concerns**: Each package has a single responsibility
2. **Dependency Injection**: Pass dependencies explicitly
3. **Interface Segregation**: Small, focused interfaces
4. **Error Handling**: Explicit error returns, no panics in library code

## Plugin System

### Plugin File Structure

Plugins are YAML files located in the `plugins/` directory:

```yaml
name: framework-name
description: Framework description
extensions:
  - .ext1
  - .ext2
patterns:
  - method: "GET|POST|PUT|DELETE|PATCH"
    route_regex: "pattern to match routes"
    param_regex: "pattern to extract path parameters"
    query_regex: "optional pattern for query parameters"
router_mount:
  use_pattern: "pattern for router mounting"
  require_pattern: "pattern for module imports"
  var_pattern: "pattern for variable declarations"
```

### Pattern Matching

#### Route Regex

The `route_regex` should capture:
1. HTTP method (first capture group)
2. Route path (second capture group)

Example for Express:
```regex
(?:app|router|express)\.(get|post|put|delete|patch|all)\s*\(\s*['"]([^'"]+)['"]
```

#### Parameter Regex

The `param_regex` extracts path parameters from the route path:
- Express: `:(\w+)` matches `:id`
- Spring: `\{(\w+)\}` matches `{id}`
- Flask: `<(\w+)>` matches `<id>`

### Router Mounting

For frameworks that support router mounting (like Express):

```yaml
router_mount:
  use_pattern: "app\\.use\\s*\\(\\s*['\"]([^'\"]+)['\"]\\s*,\\s*(\\w+)"
  require_pattern: "require\\s*\\(['\"]([^'\"]+)['\"]\\)"
  var_pattern: "(?:const|let|var)\\s+(\\w+)\\s*=.*require"
```

## Adding New Frameworks

### Step 1: Create Plugin File

Create `plugins/framework-name.yml`:

```yaml
name: framework-name
description: Framework description
extensions:
  - .ext
patterns:
  - method: "GET|POST"
    route_regex: "your pattern"
    param_regex: "your param pattern"
```

### Step 2: Test Pattern

Test your regex patterns with sample code:
- Use online regex testers
- Test with various route patterns
- Handle edge cases

### Step 3: Add Example Project

Create an example in `exp/framework-name/`:
- Include common route patterns
- Show different HTTP methods
- Include path parameters
- Add query parameters if applicable

### Step 4: Test Integration

```bash
./build/likhis -p exp/framework-name -o postman -F framework-name
```

### Step 5: Update Documentation

- Add framework to README.md
- Document any special features
- Add usage examples

## Export Formats

### Adding New Export Format

1. **Create exporter function** in `internal/exporters/`:
   ```go
   func GenerateNewFormat(routes []Route, projectPath string, env string) NewFormatData {
       // Implementation
   }
   ```

2. **Define data structures** matching the target format

3. **Handle environment variables** using helpers:
   ```go
   baseURL := getBaseURL(env)
   envName := getEnvironmentName(env)
   ```

4. **Add to main.go**:
   ```go
   case "newformat":
       data := exporters.GenerateNewFormat(allRoutes, absPath, env)
       // Serialize and write
   ```

5. **Update CLI help** and README

### Export Format Requirements

- Support environment variables for base URLs
- Include all route metadata
- Handle path and query parameters
- Support `--full` flag (dev, staging, prod)

## Testing Guidelines

### Unit Testing

- Test individual functions
- Mock external dependencies
- Test error conditions
- Aim for >80% coverage

### Integration Testing

- Use example projects in `exp/`
- Test across all frameworks
- Verify export formats
- Test edge cases

### Test Data

- Keep test data in `exp/` directory
- Use realistic examples
- Cover various patterns
- Include error cases

### Running Tests

```bash
# All tests
scripts/test.bat  # Windows
./scripts/test.sh  # Unix/Linux/macOS

# Specific framework
./build/likhis -p exp/express -o postman -F express
```

## Performance Considerations

### File Traversal

- Use BFS for predictable order
- Skip large dependency directories early
- Filter by extension before processing

### Pattern Matching

- Compile regex patterns once
- Use efficient regex patterns
- Avoid backtracking in patterns

### Memory Usage

- Process files one at a time when possible
- Avoid loading entire files into memory
- Use streaming for large files

### Optimization Tips

1. **Early Returns**: Exit early when conditions aren't met
2. **Caching**: Cache compiled regex patterns
3. **Buffering**: Use buffered I/O for file operations
4. **Parallelism**: Consider parallel processing for large projects

## Documentation Standards

### Code Comments

- Document exported functions and types
- Explain complex logic
- Include examples for public APIs
- Keep comments up-to-date

### README Updates

- Update when adding features
- Include usage examples
- Document new options
- Keep installation instructions current

### API Documentation

- Use clear function names
- Document parameters and return values
- Include error conditions
- Provide usage examples

## Release Process

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes

### Release Checklist

1. **Update CHANGELOG.md**
   - Add new features
   - List bug fixes
   - Note breaking changes

2. **Update Version**
   - Update version in code (if applicable)
   - Tag release in git

3. **Build Binaries**
   ```bash
   scripts/build.bat  # Windows
   ./scripts/build.sh  # Unix/Linux/macOS
   ```

4. **Run Full Test Suite**
   ```bash
   scripts/test.bat  # or ./scripts/test.sh
   ```

5. **Create Release Notes**
   - Summarize changes
   - Highlight new features
   - List breaking changes

6. **Create Git Tag**
   ```bash
   git tag -a v1.0.0 -m "Release version 1.0.0"
   git push origin v1.0.0
   ```

## Best Practices

### Code Quality

- Write readable, maintainable code
- Follow Go idioms and conventions
- Use `gofmt` and `golint`
- Keep functions small and focused

### Error Handling

- Always handle errors explicitly
- Provide meaningful error messages
- Use error wrapping for context
- Log errors appropriately

### Testing

- Write tests before fixing bugs
- Test edge cases
- Keep tests simple and focused
- Maintain good test coverage

### Git Workflow

- Use meaningful commit messages
- Follow conventional commits
- Keep commits focused
- Write clear PR descriptions

## Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Semantic Versioning](https://semver.org/)
- [Keep a Changelog](https://keepachangelog.com/)

---

For questions or clarifications, please open an issue or start a discussion.

