<div align="center">
  <h1>Likhis - Universal API Route Mapper</h1>
  <p><strong>Cross-Platform Deployment Automation Framework for Node.js</strong></p>
  <p>
    <img src="https://img.shields.io/npm/v/@marcuwynu23/likhis.svg" alt="npm version"/>
    <img src="https://img.shields.io/npm/dm/@marcuwynu23/likhis.svg" alt="npm downloads"/>
    <img src="https://img.shields.io/github/stars/marcuwynu23/likhis.svg" alt="GitHub Stars"/>
    <img src="https://img.shields.io/github/license/marcuwynu23/likhis.svg" alt="License Badge"/>
  </p>
</div>

---

> **Automated API route discovery and export tool for modern backend frameworks**

Likhis is a high-performance, cross-platform command-line tool written in Go that automatically analyzes backend source code to discover API routes, extract HTTP methods and parameters, and generate ready-to-import collections for popular API testing tools. It eliminates manual route documentation by intelligently parsing your codebase and producing standardized exports compatible with industry-standard testing platforms.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Usage](#usage)
- [Supported Frameworks](#supported-frameworks)
- [Export Formats](#export-formats)
- [Plugin System](#plugin-system)
- [Architecture](#architecture)
- [Limitations](#limitations)
- [Contributing](#contributing)
- [Documentation](#documentation)
- [Support](#support)
- [License](#license)

## Overview

Likhis streamlines the API documentation and testing workflow by automatically extracting route definitions from your backend codebase. The tool supports multiple popular frameworks and generates standardized exports compatible with industry-standard API testing tools, significantly reducing the time and effort required for manual route documentation and collection creation.

### Key Benefits

- **Zero Configuration**: Automatically detects framework patterns and extracts routes
- **Multi-Framework Support**: Works with Express, Flask, Django, Spring Boot, Laravel, and more
- **Extensible Architecture**: YAML-based plugin system for adding new framework support
- **Multiple Export Formats**: Generate collections for Postman, Insomnia, HTTPie Desktop, or CURL
- **Environment-Aware**: Generate separate collections for development, staging, and production environments

## Features

### Core Capabilities

- **Intelligent Route Detection**: Uses pattern matching and regex to identify route definitions across different framework conventions
- **Breadth-First Traversal**: Efficiently scans project directories using BFS algorithm, prioritizing top-level routes
- **Parameter Extraction**: Automatically detects:
  - Path parameters (`:id`, `{id}`, `<id>`)
  - Query parameters (from function signatures and annotations)
  - Request body fields (for POST/PUT operations)
- **Router Mount Detection**: Understands nested router structures and correctly prefixes base paths
- **Environment-Specific Exports**: Generate separate collections for different deployment environments

### Supported Frameworks

| Framework | Language | Detection Method |
|-----------|----------|-----------------|
| Express.js | JavaScript/TypeScript | `app.get()`, `router.post()`, etc. |
| Flask | Python | `@app.route()` decorators |
| Django | Python | `urls.py` with `path()` or `re_path()` |
| Spring Boot | Java | `@GetMapping()`, `@PostMapping()`, etc. |
| Laravel | PHP | `Route::get()`, `Route::post()`, etc. |

## Prerequisites

- **Go 1.21+** (for building from source)
- **Windows, macOS, or Linux** (cross-platform support)
- Access to your backend project source code

## Installation

### Building from Source

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd likhis
   ```

2. **Build the executable**:
   
   **Windows (Batch)**:
   ```cmd
   scripts\build.bat
   ```
   
   **Windows (PowerShell)**:
   ```powershell
   scripts\build.ps1
   ```
   
   **Manual Build**:
   ```bash
   go build -o build/likhis.exe main.go
   ```

3. **Optional: Create symbolic link** (Windows):
   ```cmd
   scripts\link.bat
   ```
   or
   ```powershell
   scripts\link.ps1
   ```
   
   This creates a link at `C:\Bin\webserve\likhis.exe` for easy access.

The compiled executable will be located in the `build/` directory.

## Quick Start

```bash
# Scan current directory and generate Postman collection
likhis -p . -o postman

# Scan specific project with auto-detection
likhis -p ./my-backend -o insomnia

# Generate for specific framework
likhis -p ./express-app -o postman -F express

# Generate environment-specific exports (dev, staging, prod)
likhis -p ./my-api -o postman --full
```

## Usage

### Command-Line Interface

```bash
likhis [OPTIONS]
```

### Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--path` | `-p` | Path to project root directory | Current directory |
| `--output` | `-o` | Output format: `postman`, `insomnia`, `httpie`, `curl` | `postman` |
| `--file` | `-f` | Custom output file name | Auto-generated based on format |
| `--output-path` | `-O` | Output directory path | Current directory |
| `--framework` | `-F` | Target framework: `auto` or plugin name | `auto` |
| `--full` | | Generate exports for dev, staging, and prod | `false` |

### Examples

#### Basic Usage

```bash
# Express.js project - Postman export
likhis -p ./express-app -o postman -F express

# Flask project - Insomnia export
likhis -p ./flask-app -o insomnia -F flask

# Spring Boot project - HTTPie Desktop export
likhis -p ./spring-app -o httpie -F spring

# Laravel project - CURL script
likhis -p ./laravel-app -o curl -F laravel
```

#### Advanced Usage

```bash
# Custom output file name
likhis -p ./my-api -o postman -f custom-collection.json

# Custom output directory
likhis -p ./my-api -o postman --output-path ./exports

# Custom output directory and file name
likhis -p ./my-api -o postman -f my-api -O ./exports
# Generates: ./exports/my-api.json

# Auto-detect framework
likhis -p ./my-project -o postman

# Generate environment-specific collections
likhis -p ./my-api -o postman --full
# Generates:
# - postman-collection-dev.json
# - postman-collection-staging.json
# - postman-collection-prod.json

# Generate environment-specific collections to custom directory
likhis -p ./my-api -o postman --full --output-path ./api-exports
# Generates:
# - ./api-exports/postman-collection-dev.json
# - ./api-exports/postman-collection-staging.json
# - ./api-exports/postman-collection-prod.json
```

#### Framework-Specific Examples

**Express.js with Router Mounting**:
```bash
likhis -p ./express-app/src -o postman -F express
```

**Django URL Patterns**:
```bash
likhis -p ./django-project -o insomnia -F django
```

**Spring Boot Controllers**:
```bash
likhis -p ./spring-app/src/main/java -o postman -F spring
```

## Supported Frameworks

### Node.js/Express

Detects routes defined using:
- `app.get()`, `app.post()`, `app.put()`, `app.delete()`, `app.patch()`
- `router.get()`, `router.post()`, etc.
- Router mounting: `app.use('/api', router)`

**Example**:
```javascript
app.get('/users/:id', handler);
router.post('/products', handler);
```

### Python/Flask

Detects routes defined using:
- `@app.route()` decorators
- HTTP methods specified in `methods` parameter

**Example**:
```python
@app.route('/users/<id>', methods=['GET'])
def get_user(id):
    pass
```

### Django

Detects routes from `urls.py` files:
- `path()` function
- `re_path()` function

**Example**:
```python
path('users/<int:id>/', views.user_detail)
```

### Java/Spring Boot

Detects routes from controller classes:
- `@GetMapping()`, `@PostMapping()`, `@PutMapping()`, `@DeleteMapping()`, `@PatchMapping()`
- `@RequestMapping()` at class and method level
- `@RequestParam` for query parameters
- `@PathVariable` for path parameters

**Example**:
```java
@RestController
@RequestMapping("/users")
public class UserController {
    @GetMapping("/{id}")
    public ResponseEntity<User> getUser(@PathVariable String id) {
        // ...
    }
}
```

### PHP/Laravel

Detects routes from route files:
- `Route::get()`, `Route::post()`, `Route::put()`, `Route::delete()`, `Route::patch()`

**Example**:
```php
Route::get('/users/{id}', [UserController::class, 'show']);
```

## Export Formats

### Postman Collection (v2.1)

Generates a complete Postman Collection v2.1 JSON file ready for import.

**Features**:
- Environment variables for base URLs
- Organized request structure
- Path and query parameters
- Request body templates

**Usage**:
```bash
likhis -p ./my-api -o postman
# Output: postman-collection.json
```

### Insomnia Export

Generates a native Insomnia export compatible with Insomnia Desktop.

**Features**:
- Workspace structure
- Request groups
- Environment variables
- Cookie jar support

**Usage**:
```bash
likhis -p ./my-api -o insomnia
# Output: insomnia-export.json
```

### HTTPie Desktop Collection

Generates a collection file compatible with HTTPie Desktop.

**Usage**:
```bash
likhis -p ./my-api -o httpie
# Output: httpie-collection.json
```

**Note**: HTTPie Desktop may have limited import support. Consider using the CURL format as an alternative.

### CURL Script

Generates a shell script containing ready-to-use `curl` commands for each route.

**Usage**:
```bash
likhis -p ./my-api -o curl
# Output: curl-commands.sh
```

## Plugin System

Likhis features an extensible YAML-based plugin architecture that allows you to add support for new frameworks without modifying the source code.

### Plugin Structure

Plugins are YAML files located in the `plugins/` directory (next to the executable). Each plugin defines:

```yaml
name: framework-name
description: Framework description
extensions:
  - .ext1
  - .ext2
patterns:
  - method: "GET|POST|PUT|DELETE|PATCH"
    route_regex: "regex pattern to match routes"
    param_regex: "regex pattern to extract path parameters"
router_mount:
  use_pattern: "regex for router mounting"
  require_pattern: "regex for module imports"
  var_pattern: "regex for variable declarations"
```

### Creating a Custom Plugin

1. **Create a YAML file** in the `plugins/` directory:
   ```bash
   plugins/myframework.yml
   ```

2. **Define the plugin structure**:
   ```yaml
   name: myframework
   description: My Custom Framework
   extensions:
     - .myext
   patterns:
     - method: "GET|POST"
       route_regex: "route\\.(get|post)\\s*\\(['\"]([^'\"]+)['\"]"
       param_regex: "\\{(\\w+)\\}"
   ```

3. **Use your plugin**:
   ```bash
   likhis -p ./my-project -o postman -F myframework
   ```

### Included Plugins

The following plugins are included by default:

- `express.yml` - Node.js Express.js framework
- `flask.yml` - Python Flask framework
- `django.yml` - Python Django framework
- `spring.yml` - Java Spring Boot framework
- `laravel.yml` - PHP Laravel framework

### Plugin Location

Plugins are loaded from:
1. `{executable_directory}/plugins/` (primary location)
2. `./plugins/` (fallback if executable directory doesn't exist)

## Architecture

### Processing Pipeline

1. **File Traversal**: Performs breadth-first search (BFS) through project directories
   - Skips common dependency folders (`node_modules`, `vendor`, `.git`, `__pycache__`, etc.)
   - Filters files by extension based on framework

2. **Route Detection**: Analyzes source files using framework-specific patterns
   - Regex-based pattern matching
   - Router mounting detection (for Express.js)
   - Class-level annotation detection (for Spring Boot)

3. **Parameter Extraction**: Extracts route metadata
   - Path parameters from route patterns
   - Query parameters from function signatures
   - Request body fields (heuristic-based)

4. **Normalization**: Converts framework-specific routes to unified structure

5. **Export Generation**: Transforms normalized routes to target format

### Internal Route Structure

Routes are normalized to a unified JSON structure:

```json
{
  "path": "/users/:id",
  "method": "GET",
  "params": ["id"],
  "query": ["page", "limit"],
  "body": ["name", "email"],
  "file": "/path/to/file.js",
  "line": 42
}
```

## Limitations

While Likhis provides comprehensive route detection capabilities, please be aware of the following limitations:

- **Heuristic Detection**: Query parameter and body field detection uses heuristic algorithms and may not capture all parameters in complex scenarios
- **Complex Patterns**: Some advanced or unconventional route patterns may not be detected automatically
- **Middleware Configuration**: Authentication headers and middleware configurations are not automatically extracted from route definitions
- **Dynamic Routes**: Routes generated dynamically at runtime through code execution may not be detected during static analysis
- **Parsing Method**: Currently uses regex-based parsing; future versions may incorporate AST (Abstract Syntax Tree) parsing for improved accuracy and coverage

For the most up-to-date information on limitations and planned improvements, please refer to the [CHANGELOG.md](CHANGELOG.md) and project issues.

## Contributing

We welcome contributions from the community! Whether you're fixing bugs, adding features, or improving documentation, your help makes Likhis better for everyone.

For detailed information on how to contribute, please see our [Contributing Guide](CONTRIBUTING.md). Key areas where contributions are particularly valuable:

- Additional framework support via plugins
- Improved parameter detection algorithms
- AST-based parsing for better accuracy
- Additional export formats
- Documentation improvements

### Quick Start for Contributors

1. Fork the repository
2. Create a feature branch
3. Make your changes following our [Development Guidelines](GUIDELINES.md)
4. Test with the example applications in `exp/`:
   ```bash
   scripts\test.bat
   ```
5. Submit a pull request

For comprehensive development guidelines, coding standards, and plugin creation instructions, please refer to [GUIDELINES.md](GUIDELINES.md).

## Documentation

- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Guidelines for contributing to the project
- **[GUIDELINES.md](GUIDELINES.md)** - Development guidelines, architecture, and best practices
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes

## Support

### Sponsorship

If you find Likhis useful and would like to support its development, please consider sponsoring the project. See [FUNDING.yml](FUNDING.yml) for sponsorship options.

### Getting Help

- **Issues**: Report bugs or request features via [GitHub Issues](https://github.com/marcuwynu23/likhis/issues)
- **Discussions**: Ask questions and share ideas in [GitHub Discussions](https://github.com/marcuwynu23/likhis/discussions)

## Author

**Mark Wayne Menorca**

## License

This project is open source and available for use. See the repository for license details.

---

**Note**: This tool is designed to assist with API documentation and testing. Always verify generated routes against your actual API implementation.
