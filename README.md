<div align="center">
  <h1>Likhis</h1>
  <p><strong>Cross-Platform API Route Discovery and Export Tool</strong></p>
  <p>
    <img src="https://img.shields.io/github/v/release/marcuwynu23/likhis?include_prereleases&style=flat-square" alt="Release"/>
    <img src="https://img.shields.io/github/go-mod/go-version/marcuwynu23/likhis?style=flat-square" alt="Go Version"/>
    <img src="https://img.shields.io/github/stars/marcuwynu23/likhis?style=flat-square" alt="GitHub Stars"/>
    <img src="https://img.shields.io/github/forks/marcuwynu23/likhis?style=flat-square" alt="GitHub Forks"/>
    <img src="https://img.shields.io/badge/license-Apache%202.0-blue?style=flat-square" alt="License"/>
    <img src="https://img.shields.io/github/issues/marcuwynu23/likhis?style=flat-square" alt="GitHub Issues"/>
    <img src="https://img.shields.io/github/actions/workflow/status/marcuwynu23/likhis/test.yml?style=flat-square" alt="CI"/>
  </p>
  <p>
    <strong>Automatically discover API routes from backend source code and export ready-to-import collections.</strong>
  </p>
  ➡️ <strong><a href="USER-GUIDE.md">Read the full user guide →</a></strong>
</div>

### Pronunciation & Origin

- **Pronunciation:** `/lik-hees/` (**Lik** as in *lick* + **his** with a long *ee* sound like *heez*)
- **Etymology:** Named after the Tagalog word **Likha** (to create/build) combined with a phonetic twist on **Lihis** (to deviate/pivot).

## Table of Contents

- [What Is Likhis?](#what-is-likhis)
- [Use Cases](#use-cases)
- [Benefits](#benefits)
- [Advantages Over Other Tools](#advantages-over-other-tools)
- [User Guide](USER-GUIDE.md)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [CLI Reference](#cli-reference)
- [Plugin System](#plugin-system)
- [Example Output](#example-output)
- [CI/CD Integration](#cicd-integration)
- [Development](#development)
- [Architecture](#architecture)
- [Limitations](#limitations)
- [Contributing](CONTRIBUTING.md)
- [License](#license)

## What Is Likhis?

**Likhis** is a high-performance, cross-platform CLI tool written in Go that automatically analyzes backend source code to discover API routes, extract HTTP methods and parameters, and generate ready-to-import collections for popular API testing tools.

### What It Does

- **Discovers** — Scans backend projects using breadth-first traversal, automatically detecting route definitions across multiple frameworks
- **Extracts** — Pulls HTTP methods, path parameters, query parameters, and request body fields from source code
- **Exports** — Generates collections for Postman, Insomnia, HTTPie Desktop, cURL, and OpenAPI formats
- **Detects** — Router mounting structures (e.g., Express.js `app.use()`) and correctly prefixes base paths
- **Filters** — Route filtering via plugin-level ignore patterns to exclude health checks, WebSocket endpoints, etc.
- **Environment-Aware** — Generates separate exports for development, staging, and production environments
- **Extends** — YAML-based plugin architecture lets you add framework support without modifying source code

### Why Use It?

| Problem | How Likhis Solves It |
|---|---|
| **Manual API documentation is tedious and error-prone** | Likhis scans your codebase in seconds and extracts every route automatically |
| **Keeping Postman/Insomnia collections in sync** | Regenerate collections any time routes change with a single command |
| **Onboarding new team members** | Give them a complete API collection generated directly from your source |
| **Multiple frameworks in one codebase** | Auto-detects framework patterns; plugin system handles any framework |
| **Environment-specific API testing** | `--full` flag generates dev, staging, and production collections at once |

### The Philosophy

1. **Minimal setup, maximum value.** Point Likhis at your project and go — zero configuration required for most frameworks.
2. **Your code stays your source of truth.** Collections are derived from code, not maintained separately.
3. **Extensible by design.** Adding a new framework means writing a YAML file, not Go code.

## Use Cases

| Scenario | How Likhis Helps |
|---|---|
| **You just joined a team with a large Express.js API** | Generate a complete Postman collection from the project to explore endpoints |
| **You need to share API endpoints with frontend developers** | Export an Insomnia workspace they can import immediately |
| **You maintain a Django REST API with dozens of endpoints** | Keep auto-generated collections in CI to catch route changes |
| **You're evaluating different API testing tools** | Likhis supports Postman, Insomnia, HTTPie, cURL, and OpenAPI from the same scan |
| **You have a monorepo with multiple backend services** | Run Likhis on each service directory for per-service collections |

## Benefits

- **Zero configuration** — automatically detects framework patterns and extracts routes
- **Multi-framework** — Express.js, Flask, Django, Spring Boot, Laravel, and more via plugins
- **Multiple export formats** — Postman v2.1, Insomnia, HTTPie Desktop, cURL scripts, OpenAPI specs
- **Extensible plugin system** — add new framework support with a YAML file, no Go code required
- **Environment-aware exports** — generate separate dev, staging, and production collections
- **Router mounting detection** — correctly resolves Express.js-style nested routers with base paths
- **Parameter extraction** — detects path params, query params, and request body fields per route
- **Route filtering** — plugin-level ignore patterns to exclude internal or irrelevant endpoints
- **Cross-platform** — works on Windows, macOS, and Linux with prebuilt binaries
- **Fast** — breadth-first traversal with BFS algorithm; skips dependency folders automatically

## Advantages Over Other Tools

| Aspect | Likhis | Postman CLI | Insomnia CLI | API Extractor | Manual |
|---|---|---|---|---|---|
| **Setup time** | ~10 seconds | Minutes | Minutes | Hours | Ongoing effort |
| **Framework auto-detection** | Yes | No | No | Limited | N/A |
| **Plugin system** | YAML-based | No | No | No | N/A |
| **Export formats** | Postman, Insomnia, HTTPie, cURL, OpenAPI | Postman only | Insomnia only | Postman only | Any |
| **Router mount detection** | Yes | No | No | No | N/A |
| **Environment exports** | Dev, staging, prod | Limited | Limited | Limited | Full control |
| **OpenAPI support** | Yes (export) | Import only | Import only | Yes | Yes |
| **Runtime** | Go (single binary) | Node.js | Electron | Java | Any |
| **License** | Apache 2.0 | Proprietary | MIT | Proprietary | N/A |
| **Offline capable** | Yes | Yes | Yes | Yes | Yes |
| **Parameter extraction** | Path, query, body | No | No | Limited | Manual |

## Installation

### Building from Source

1. **Clone the repository**:

```bash
git clone https://github.com/marcuwynu23/likhis.git
cd likhis
```

2. **Build the executable**:

**Using Make (cross-platform)**:

```bash
make build
```

**Manual Build**:

```bash
go build -o build/likhis main.go
```

The compiled executable will be located in the `build/` directory.

### Verify Installation

```bash
./build/likhis --help
```

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

## CLI Reference

```bash
likhis [OPTIONS]
```

| Flag | Short | Default | Description |
|---|---|---|---|
| `--path` | `-p` | `.` | Path to project root directory |
| `--output` | `-o` | `postman` | Output format: `postman`, `insomnia`, `httpie`, `curl`, `openapi` |
| `--file` | `-f` | Auto-generated | Custom output file name (without extension) |
| `--output-path` | `-O` | `.` | Output directory path |
| `--framework` | `-F` | `auto` | Target framework plugin or `auto` for detection |
| `--full` | | `false` | Generate exports for dev, staging, and production |

### Basic Examples

**Express.js project — Postman export:**

```bash
likhis -p ./express-app -o postman -F express
```

**Flask project — Insomnia export:**

```bash
likhis -p ./flask-app -o insomnia -F flask
```

**Spring Boot — HTTPie Desktop export:**

```bash
likhis -p ./spring-app -o httpie -F spring
```

**Laravel project — cURL script:**

```bash
likhis -p ./laravel-app -o curl -F laravel
```

**Generate OpenAPI spec:**

```bash
likhis -p ./my-api -o openapi
```

### Advanced Examples

**Custom output file and directory:**

```bash
likhis -p ./my-api -o postman -f my-api -O ./exports
# Generates: ./exports/my-api.json
```

**Auto-detect framework:**

```bash
likhis -p ./my-project -o postman
```

**Environment-specific collections:**

```bash
likhis -p ./my-api -o postman --full --output-path ./api-exports
# Generates:
# - ./api-exports/postman-collection-dev.json
# - ./api-exports/postman-collection-staging.json
# - ./api-exports/postman-collection-prod.json
```

## Plugin System

Likhis features an extensible YAML-based plugin architecture that lets you add support for new frameworks without modifying the Go source code.

### Included Plugins

| Plugin | Framework | Files |
|---|---|---|
| `express` | Node.js Express.js | `.js`, `.ts` |
| `flask` | Python Flask | `.py` |
| `django` | Python Django | `.py` |
| `spring` | Java Spring Boot | `.java` |
| `laravel` | PHP Laravel | `.php` |

### Plugin Structure

Plugins are YAML files in the `plugins/` directory:

```yaml
name: express
description: Node.js Express.js framework
extensions:
  - .js
  - .ts
patterns:
  - method: "GET|POST|PUT|DELETE|PATCH|ALL"
    route_regex: "\\w+\\.(get|post|put|delete|patch|all)\\s*\\(\\s*['\"]([^'\"]+)['\"]"
    param_regex: ":(\\w+)"
router_mount:
  use_pattern: "\\w+\\.use\\s*\\(\\s*['\"]([^'\"]+)['\"]\\s*,\\s*(\\w+)"
  require_pattern: "require\\s*\\(['\"]([^'\"]+)['\"]\\)"
  var_pattern: "(?:const|let|var)\\s+(\\w+)\\s*=.*require"
```

### Creating a Custom Plugin

1. Create a YAML file in the `plugins/` directory:

```bash
plugins/myframework.yml
```

2. Define the plugin structure with regex patterns for route detection.

3. Use your plugin:

```bash
likhis -p ./my-project -o postman -F myframework
```

### Plugin Load Order

Plugins are loaded from these directories (later overrides earlier):

1. `{project}/plugins/` — project-specific custom plugins
2. `{executable_directory}/plugins/` — bundled plugins
3. `./plugins/` — fallback directory

## Example Output

### Postman Collection (v2.1)

```json
{
  "info": {
    "name": "Likhis Export (Development)",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "GET /users/:id",
      "request": {
        "method": "GET",
        "url": {
          "raw": "{{base_url}}/users/:id",
          "host": ["{{base_url}}"],
          "path": ["users", ":id"],
          "variable": [
            { "key": "id", "value": "" }
          ]
        }
      }
    }
  ]
}
```

### Insomnia Export

```json
{
  "_type": "export",
  "__export_format": 4,
  "__export_source": "likhis",
  "resources": [
    {
      "_type": "request",
      "method": "GET",
      "url": "{{ base_url }}/users/{{ id }}"
    }
  ]
}
```

### cURL Script

```bash
#!/bin/bash
# API Requests - Generated by Likhis
BASE_URL="${BASE_URL:-http://localhost:3000}"
curl -X GET "$BASE_URL/users/:id"
curl -X POST "$BASE_URL/products" -H "Content-Type: application/json"
```

## CI/CD Integration

### GitHub Actions

```yaml
name: Generate API Collection

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Download Likhis
        run: |
          curl -L -o likhis https://github.com/marcuwynu23/likhis/releases/latest/download/likhis-linux-amd64
          chmod +x likhis
      - name: Generate Postman Collection
        run: ./likhis -p ./backend -o postman -O ./api-collections
      - name: Upload Artifact
        uses: actions/upload-artifact@v4
        with:
          name: api-collections
          path: ./api-collections/
```

### GitLab CI

```yaml
generate-api-collection:
  stage: test
  script:
    - curl -L -o likhis https://github.com/marcuwynu23/likhis/releases/latest/download/likhis-linux-amd64
    - chmod +x likhis
    - ./likhis -p ./backend -o postman -O ./api-collections
  artifacts:
    paths:
      - ./api-collections/
```

## Development

### Prerequisites

| Tool | Version | Purpose |
|---|---|---|
| Go | 1.21+ | Compile and run |
| Make | Any | Build automation |

### Quick Start for Contributors

```bash
make build          # Build the binary
make test           # Run unit tests
make test-integration # Run integration tests against example projects
make lint           # Run golangci-lint
make coverage       # Run tests with coverage report
```

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed development guidelines.

## Architecture

### Processing Pipeline

1. **File Traversal** — BFS through project directories, skipping `node_modules`, `vendor`, `.git`, etc.
2. **Route Detection** — Framework-specific regex patterns (hardcoded or plugin-sourced) identify route definitions
3. **Router Mounting** — For Express.js-style frameworks, resolves nested router structures and prepends base paths
4. **Parameter Extraction** — Scans handler functions for `req.query`, `req.body`, `@RequestParam`, etc.
5. **Route Filtering** — Applies plugin-level ignore patterns to exclude matched routes
6. **Export Generation** — Transforms normalized routes into target format (Postman, Insomnia, etc.)

### Internal Route Structure

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

### Package Layout

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
└── build/                   # Compiled binaries
```

## Limitations

- **Heuristic detection** — query parameter and body field detection is heuristic; complex patterns may be missed
- **Static analysis only** — routes generated dynamically at runtime are not detected
- **No AST parsing** — currently uses regex; future versions may incorporate AST-based parsing for accuracy
- **Middleware not extracted** — authentication headers and middleware config are not captured

For the most up-to-date information, see [CHANGELOG.md](CHANGELOG.md) and [GitHub Issues](https://github.com/marcuwynu23/likhis/issues).

## License

This project is licensed under the Apache License 2.0 — see [LICENSE](LICENSE) for details.

All contributors must adhere to our [Code of Conduct](CODE_OF_CONDUCT.md).

---

**Note:** This tool is designed to assist with API documentation and testing. Always verify generated routes against your actual API implementation.
