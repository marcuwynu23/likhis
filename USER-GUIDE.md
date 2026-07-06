<div align="center">

# User Guide

  <p><strong>Likhis — Cross-Platform API Route Discovery and Export Tool</strong></p>
  <p>Comprehensive reference for installation, configuration, and workflows.</p>
</div>

## Table of Contents

- [Installation](#installation)
- [Quick Start](#quick-start)
- [Command Reference](#command-reference)
- [Configuration](#configuration)
- [Plugin System](#plugin-system)
- [Concepts](#concepts)
- [Export Formats](#export-formats)
- [CI/CD Integration](#cicd-integration)
- [Workflows](#workflows)
- [Troubleshooting](#troubleshooting)
- [FAQ](#faq)

---

## Installation

### Prerequisites

| Requirement | Details                                    |
| ----------- | ------------------------------------------ |
| **Go**      | 1.21+ (only required to build from source) |
| **OS**      | Windows, macOS, or Linux                   |
| **Disk**    | ~10 MB for the binary                      |

### Option 1: Build from Source

```bash
git clone https://github.com/marcuwynu23/likhis.git
cd likhis
make build
```

The binary is placed in `build/likhis` (or `build/likhis.exe` on Windows).

### Option 2: Download Prebuilt Binary

Download the latest release for your platform from the [releases page](https://github.com/marcuwynu23/likhis/releases).

| Platform | Architecture          | File                       |
| -------- | --------------------- | -------------------------- |
| Windows  | amd64                 | `likhis-windows-amd64.exe` |
| Linux    | amd64                 | `likhis-linux-amd64`       |
| macOS    | amd64                 | `likhis-darwin-amd64`      |
| macOS    | arm64 (Apple Silicon) | `likhis-darwin-arm64`      |

### Option 3: Install via `go install`

```bash
go install github.com/marcuwynu23/likhis@latest
```

### Verify Installation

```bash
likhis --help
```

### Linking to PATH

**Using Make:**

```bash
make link
```

This creates a symbolic link in `C:\Bin\tools` (Windows) or `/usr/local/bin` (macOS/Linux).

---

## Quick Start

### 1. Scan a Project

```bash
likhis -p ./my-backend -o postman
```

This scans `./my-backend` for API routes and generates a Postman v2.1 collection.

### 2. Specify a Framework

```bash
likhis -p ./my-backend -o insomnia -F express
```

Skips auto-detection and uses the Express.js plugin directly.

### 3. Generate Environment Exports

```bash
likhis -p ./my-backend -o postman --full
```

Generates three files: dev, staging, and production.

### 4. Export for Different Tools

```bash
likhis -p ./my-backend -o curl      # cURL script (executable)
likhis -p ./my-backend -o httpie    # HTTPie Desktop collection
likhis -p ./my-backend -o openapi   # OpenAPI 3.0 YAML spec
```

---

## Command Reference

### Global Flags

| Flag            | Short | Default        | Description                                                       |
| --------------- | ----- | -------------- | ----------------------------------------------------------------- |
| `--path`        | `-p`  | `.`            | Path to the project root directory                                |
| `--output`      | `-o`  | `postman`      | Output format: `postman`, `insomnia`, `httpie`, `curl`, `openapi` |
| `--file`        | `-f`  | Auto-generated | Custom output file name (without extension)                       |
| `--output-path` | `-O`  | `.`            | Output directory path                                             |
| `--framework`   | `-F`  | `auto`         | Framework to target (`auto` or plugin name)                       |
| `--full`        |       | `false`        | Generate exports for dev, staging, and production environments    |

### `--path` / `-p`

Specify the root directory to scan. Likhis traverses this directory using breadth-first search, skipping common dependency folders (`node_modules`, `vendor`, `.git`, `__pycache__`, etc.).

```bash
likhis -p ./src
likhis -p /absolute/path/to/project
```

### `--output` / `-o`

Select the output format. See [Export Formats](#export-formats) for details on each.

```bash
likhis -p . -o postman    # Postman Collection v2.1 (JSON)
likhis -p . -o insomnia   # Insomnia Export (JSON)
likhis -p . -o httpie     # HTTPie Desktop (JSON)
likhis -p . -o curl       # cURL script (Bash)
likhis -p . -o openapi    # OpenAPI 3.0 (YAML)
```

### `--file` / `-f`

Customize the output file name. The extension is added automatically based on format.

```bash
likhis -p . -o postman -f my-collection   # → my-collection.json
likhis -p . -o curl -f api-scripts         # → api-scripts.sh
```

### `--output-path` / `-O`

Specify where to write the output file. Creates the directory if it doesn't exist.

```bash
likhis -p . -o postman -O ./api-exports
```

Combine with `--file`:

```bash
likhis -p . -o postman -f my-api -O ./exports
# → ./exports/my-api.json
```

### `--framework` / `-F`

Select the framework plugin to use. Use `auto` for automatic detection.

```bash
likhis -p . -o postman -F express    # Force Express.js
likhes -p . -o postman -F auto       # Auto-detect (default)
```

Available built-in frameworks: `auto`, `express`, `flask`, `django`, `spring`, `laravel`.

Custom plugin names match their YAML filename (without extension).

### `--full`

Generate exports for all three environments simultaneously.

```bash
likhis -p . -o postman --full
# Generates:
#   postman-collection-dev.json
#   postman-collection-staging.json
#   postman-collection-prod.json
```

Each environment uses a different base URL:
| Environment | Base URL |
|---|---|
| dev | `http://localhost:3000` |
| staging | `https://staging-api.example.com` |
| prod | `https://api.example.com` |

---

## Configuration

Likhis requires no configuration file — all options are passed via CLI flags.

### Configuration Precedence

1. CLI flags (highest priority)
2. Plugin YAML definitions
3. Hardcoded defaults

### Environment Export Behavior

The `--full` flag controls multi-environment export. When enabled, the output file name pattern becomes:

```
{base-name}-{environment}.{ext}
```

Example: `postman-collection-dev.json`, `postman-collection-staging.json`, `postman-collection-prod.json`.

---

## Plugin System

### Overview

Plugins are YAML files that define how Likhis detects routes for a specific framework. Each plugin specifies:

- File extensions to scan
- Regex patterns for route matching
- Regex patterns for parameter extraction
- Router mounting patterns (optional)
- Route ignore patterns (optional)

### Built-in Plugins

| Plugin File   | Framework          | Extensions   |
| ------------- | ------------------ | ------------ |
| `express.yml` | Node.js Express.js | `.js`, `.ts` |
| `flask.yml`   | Python Flask       | `.py`        |
| `django.yml`  | Python Django      | `.py`        |
| `spring.yml`  | Java Spring Boot   | `.java`      |
| `laravel.yml` | PHP Laravel        | `.php`       |

### Plugin YAML Schema

```yaml
# Required
name: string # Plugin identifier
description: string # Human-readable description
extensions:
  - string # File extensions to scan (e.g., ".js")
patterns:
  - method: string # HTTP method or pipe-separated list: "GET|POST|PUT|DELETE|PATCH"
    route_regex: string # Regex with two capture groups: method and path
    param_regex: string # Regex to extract path parameters from route path
    query_regex: string # Optional regex to extract methods from Flask-style methods=[] syntax

# Optional
router_mount:
  use_pattern: string # Regex to detect app.use('/path', router) calls
  require_pattern: string # Regex to detect require/import statements
  var_pattern: string # Regex to capture the router variable name

ignore:
  - string # Regex patterns for routes to exclude
```

### Pattern Reference

#### `route_regex`

Must capture two groups:

1. HTTP method (lowercase): `get`, `post`, `put`, `delete`, `patch`, `all`
2. Route path (string)

**Express.js example:**

```regex
\w+\.(get|post|put|delete|patch|all)\s*\(\s*['"]([^'"]+)['"]
```

Matches: `app.get("/users", handler)`, `router.post("/products", handler)`

#### `param_regex`

Extracts path parameter names from the route path.

| Framework    | Pattern     | Example       |
| ------------ | ----------- | ------------- |
| Express      | `:(\w+)`    | `:id` → `id`  |
| Spring       | `\{(\w+)\}` | `{id}` → `id` |
| Flask/Django | `<(\w+)>`   | `<id>` → `id` |
| Laravel      | `\{(\w+)\}` | `{id}` → `id` |

#### `router_mount`

For frameworks like Express that support nested routers, these patterns detect:

- **`use_pattern`**: Matches `app.use('/api', router)` calls, capturing base path and router variable
- **`require_pattern`**: Matches `require('./routes/users')` calls, capturing the module path
- **`var_pattern`**: Matches `const usersRouter = require(...)`, capturing the variable name

### Creating a Custom Plugin

**Example: Adding support for a `fastify` framework:**

```yaml
# plugins/fastify.yml
name: fastify
description: Node.js Fastify framework
extensions:
  - .js
  - .ts
patterns:
  - method: "GET|POST|PUT|DELETE|PATCH"
    route_regex: "fastify\\.(get|post|put|delete|patch)\\s*\\(\\s*['\"]([^'\"]+)['\"]"
    param_regex: ":(\\w+)"
```

Then use it:

```bash
likhis -p ./fastify-app -o postman -F fastify
```

### Plugin Load Order

Plugins are discovered in this priority order (later directories override earlier ones):

1. `{project_directory}/plugins/`
2. `{executable_directory}/plugins/`
3. `{current_working_directory}/plugins/`

This allows project-specific custom plugins to coexist with globally installed ones.

### Ignore Patterns

The `ignore` field in a plugin YAML specifies regex patterns for routes to exclude from output. This is useful for filtering out:

- Health check endpoints: `^/health$`
- WebSocket endpoints: `^/io(/.*)?$`
- Internal/admin routes

**Example from `express.yml`:**

```yaml
ignore:
  - "^/io(/.*)?$"
  - "^/health$"
```

Routes whose path matches any ignore pattern are excluded from all exports.

---

## Concepts

### How Route Detection Works

Likhis performs static analysis on source files. It does not execute code or make network requests.

1. **BFS Traversal**: The scanner walks the directory tree breadth-first, collecting files whose extensions match the target framework's plugin.

2. **Regex Matching**: Each file is scanned line by line. Lines matching the plugin's `route_regex` are extracted as route candidates.

3. **Router Mount Resolution**: For frameworks with router support (Express.js), Likhis first builds a map of router files to their base paths by scanning for `app.use('/prefix', routerVar)` calls and resolving `require()` chains.

4. **Parameter Detection**: After identifying a route line, Likhis scans the following lines (within the handler function) for patterns like `req.query.param`, `req.body.param`, `@RequestParam`, etc.

5. **Route Filtering**: Routes matching any plugin-level ignore patterns are removed from the result set.

### Framework Auto-Detection

When `-F auto` (default), Likhis checks each file's extension against all loaded plugins and applies every matching plugin's patterns to the file. This allows mixed-framework projects to be scanned in a single pass.

### Internal Route Structure

All routes are normalized to this structure before export:

```go
type Route struct {
    Path   string   // Route path, e.g., "/users/:id"
    Method string   // HTTP method, e.g., "GET"
    Params []string // Path parameters, e.g., ["id"]
    Query  []string // Query parameters, e.g., ["page", "limit"]
    Body   []string // Request body fields, e.g., ["name", "email"]
    File   string   // Source file path
    Line   int      // Line number in source file
}
```

### Supported Framework Patterns

#### Express.js

Detects:

- `app.get()`, `app.post()`, `app.put()`, `app.delete()`, `app.patch()`, `app.all()`
- `router.get()`, `router.post()`, etc.
- Router mounting: `app.use('/api', router)`
- Path parameters: `:id`, `:userId`
- Query params: `req.query.page`, `req.query['page']`, destructured queries

#### Flask

Detects:

- `@app.route()` decorators
- `@blueprint.route()` decorators
- HTTP methods from `methods=['GET', 'POST']` parameter
- Path parameters: `<id>`, `<int:user_id>`
- Query params: `request.args.get('page')`

#### Django

Detects:

- `path()` function calls in `urls.py`
- `re_path()` function calls in `urls.py`
- Path parameters: `<int:id>`, `<slug:slug>`
- Defaults to GET method (Django views determine the actual method)

#### Spring Boot

Detects:

- `@GetMapping`, `@PostMapping`, `@PutMapping`, `@DeleteMapping`, `@PatchMapping`
- `@RequestMapping` at class and method level
- `@RequestParam` for query parameters
- `@PathVariable` for path parameters
- `@RequestBody` for body parameters
- Path parameters: `{id}`, `{userId}`

#### Laravel

Detects:

- `Route::get()`, `Route::post()`, `Route::put()`, `Route::delete()`, `Route::patch()`, `Route::any()`
- Path parameters: `{id}`, `{slug}`
- Query params: `$request->query('page')`, `$request->input('email')`

---

## Export Formats

### Postman Collection v2.1

The default export format. Generates a JSON file compatible with Postman's import feature.

**Features:**

- Environment variable for base URL (`{{base_url}}`)
- Path variables for route parameters
- Query parameter placeholders
- Request body templates for POST/PUT

```bash
likhis -p ./my-api -o postman
# Output: postman-collection.json
```

### Insomnia Export

Generates a native Insomnia workspace export file.

**Features:**

- Request group organization
- Environment variable for base URL
- Cookie jar support

```bash
likhis -p ./my-api -o insomnia
# Output: insomnia-export.json
```

### HTTPie Desktop Collection

Generates a collection compatible with HTTPie Desktop.

```bash
likhis -p ./my-api -o httpie
# Output: httpie-collection.json
```

### cURL Script

Generates an executable Bash script with `curl` commands for each route.

**Features:**

- `BASE_URL` environment variable (defaults to `http://localhost:3000`)
- Executable file permissions set automatically
- Content-Type headers for POST/PUT requests

```bash
likhis -p ./my-api -o curl
# Output: api-requests.sh (chmod +x)
```

### OpenAPI 3.0 Specification

Generates an OpenAPI 3.0 YAML specification file.

**Features:**

- Paths organized by route
- Parameter definitions
- Request body schemas
- Environment variable for server URL

```bash
likhis -p ./my-api -o openapi
# Output: openapi-spec.yml
```

---

## CI/CD Integration

### GitHub Actions — Full Workflow

```yaml
name: API Collection Generation

on:
  push:
    branches: [main, develop]
  schedule:
    - cron: "0 0 * * 0" # Weekly regeneration

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Download Likhis
        run: |
          curl -L -o likhis https://github.com/marcuwynu23/likhis/releases/latest/download/likhis-linux-amd64
          chmod +x likhis

      - name: Generate All Formats
        run: |
          mkdir -p api-collections
          ./likhis -p ./backend -o postman -O api-collections
          ./likhis -p ./backend -o openapi -O api-collections
          ./likhis -p ./backend -o curl -O api-collections

      - name: Upload Artifacts
        uses: actions/upload-artifact@v4
        with:
          name: api-collections
          path: api-collections/
```

### GitLab CI

```yaml
stages:
  - test
  - api-docs

generate-api-docs:
  stage: api-docs
  image: alpine:latest
  before_script:
    - apk add curl
    - curl -L -o likhis https://github.com/marcuwynu23/likhis/releases/latest/download/likhis-linux-amd64
    - chmod +x likhis
  script:
    - ./likhis -p ./backend -o postman -O public/
    - ./likhis -p ./backend -o openapi -O public/
  artifacts:
    paths:
      - public/
    expire_in: 30 days
```

---

## Workflows

### Monorepo with Multiple Backend Services

```bash
# Generate collections for each service
likhis -p ./services/users -o postman -O ./api-collections/users
likhis -p ./services/orders -o postman -O ./api-collections/orders
likhis -p ./services/payments -o postman -O ./api-collections/payments
```

### Generating Collections Before a Release

```bash
#!/bin/bash
# pre-release.sh
VERSION=$1

likhis -p ./backend -o postman -f "api-collection-v$VERSION" -O ./release-assets
likhis -p ./backend -o openapi -f "openapi-spec-v$VERSION" -O ./release-assets
```

### Using with Docker

```dockerfile
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o /likhis main.go

FROM alpine:latest
COPY --from=builder /likhis /usr/local/bin/likhis
COPY ./backend /backend
RUN likhis -p /backend -o postman -O /exports
```

---

## Troubleshooting

### "No routes found"

| Cause                  | Fix                                                                             |
| ---------------------- | ------------------------------------------------------------------------------- |
| Wrong project path     | Verify `-p` points to the directory containing source files, not a subdirectory |
| Framework not detected | Use `-F <framework>` to force a specific plugin                                 |
| Unsupported framework  | Create a custom plugin YAML or check if the framework is supported              |
| Files filtered out     | Check if `ignore` patterns in the plugin are too aggressive                     |

### "Plugin not found"

| Cause                                | Fix                                                                                  |
| ------------------------------------ | ------------------------------------------------------------------------------------ |
| Typo in plugin name                  | Run `likhis --help` to see the list of available plugins                             |
| Plugin file not in plugins directory | Place your `.yml` file in `./plugins/` or use the project-level `plugins/` directory |
| Invalid YAML                         | Validate your YAML with `yamllint` or a YAML validator                               |

### "Routes missing parameters"

| Cause                                | Fix                                                                      |
| ------------------------------------ | ------------------------------------------------------------------------ |
| Complex handler logic                | Parameters extracted via heuristic scanning; some patterns may be missed |
| Parameters defined in separate files | The scanner only looks within the same file as the route definition      |
| Non-standard patterns                | Custom regex in a plugin may be needed for unusual parameter patterns    |

### "Export file is empty"

| Cause                    | Fix                                                           |
| ------------------------ | ------------------------------------------------------------- |
| No routes found          | Check the console output for the "Extracted X routes" message |
| All routes filtered      | Review plugin `ignore` patterns                               |
| Output path not writable | Ensure `--output-path` exists or is creatable                 |

### "Binary won't run on my system"

| Cause                      | Fix                                                     |
| -------------------------- | ------------------------------------------------------- |
| Wrong architecture         | Download the correct binary for your OS/architecture    |
| Missing execute permission | Run `chmod +x likhis` (macOS/Linux)                     |
| Windows Defender block     | Add an exclusion or use `make build` to compile locally |

---

## FAQ

**Q: Does Likhis execute my code?**

A: No. Likhis performs static analysis using regex pattern matching. It only reads files, never runs them.

**Q: Can it scan a running server?**

A: No. Likhis scans source files on disk, not running services.

**Q: Does it support GraphQL?**

A: Not directly. Likhis focuses on REST API routes. GraphQL endpoints (typically a single `POST /graphql` route) can be added via a custom plugin.

**Q: How accurate is parameter detection?**

A: Path parameters are highly accurate (extracted from route patterns). Query and body parameters use heuristic scanning of handler function bodies and may miss some fields in complex scenarios.

**Q: Can I customize the base URLs for environments?**

A: The default base URLs are hardcoded (localhost:3000 for dev, staging-api.example.com for staging, api.example.com for prod). Future versions may support user configuration.

**Q: Does Likhis support OpenAPI import?**

A: Currently Likhis only exports to OpenAPI 3.0 YAML. Import/roundtrip is not yet supported.

**Q: Can I add my own framework?**

A: Yes. Create a YAML plugin file with the appropriate regex patterns. See [Creating a Custom Plugin](#creating-a-custom-plugin).

**Q: How do I exclude certain routes?**

A: Add `ignore` patterns to the plugin YAML file. See [Ignore Patterns](#ignore-patterns).

**Q: Does Likhis work with TypeScript?**

A: Yes. The Express.js plugin scans both `.js` and `.ts` file extensions.

**Q: Can I scan multiple projects at once?**

A: Not in a single command, but you can script it: `for dir in ./services/*/; do likhis -p "$dir" -o postman; done`

**Q: What happens if a project uses multiple frameworks?**

A: Use `-F auto` to let Likhis try all registered plugins. Each file is checked against all matching plugins.

**Q: Is there a configuration file?**

A: No. All options are CLI flags. This keeps the tool simple and scriptable.

**Q: How do I update Likhis?**

A: Rebuild from source (`git pull && make build`) or download the latest release binary.

**Q: Does it scan minified files?**

A: It scans all files matching the plugin's extensions. You can exclude directories with large minified files by placing them in folders Likhis skips (like `node_modules`).

**Q: Can I contribute a new framework plugin?**

A: Absolutely. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. Submit a PR with the plugin YAML and an example project in `exp/`.

---

_For additional help, please [open an issue](https://github.com/marcuwynu23/likhis/issues) or start a [discussion](https://github.com/marcuwynu23/likhis/discussions)._
