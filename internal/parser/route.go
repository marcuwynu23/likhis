package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/marcuwynu23/likhis/internal/plugins"
)

// Route represents a unified API route structure
type Route struct {
	Path      string   `json:"path"`
	Method    string   `json:"method"`
	Params    []string `json:"params"`
	Query     []string `json:"query"`
	Body      []string `json:"body"`
	File      string   `json:"file"`
	Line      int      `json:"line"`
}

// RouteParser handles parsing routes from different frameworks
type RouteParser struct {
	framework      string
	routerBasePath map[string]string // Maps router file paths to their base paths
	plugins        map[string]*plugins.Plugin
	executablePath string
}

// NewRouteParser creates a new route parser (backward compatibility)
func NewRouteParser(framework string) *RouteParser {
	return &RouteParser{
		framework:      framework,
		routerBasePath: make(map[string]string),
		plugins:        make(map[string]*plugins.Plugin),
	}
}

// NewRouteParserWithPlugins creates a new route parser with plugin support
func NewRouteParserWithPlugins(framework string, pluginMap map[string]*plugins.Plugin, executablePath string) *RouteParser {
	return &RouteParser{
		framework:      framework,
		routerBasePath: make(map[string]string),
		plugins:        pluginMap,
		executablePath: executablePath,
	}
}

// HasRouterMountSupport checks if the framework plugin supports router mounting
func (rp *RouteParser) HasRouterMountSupport(framework string) bool {
	if framework == "auto" {
		// Check all plugins
		for _, plugin := range rp.plugins {
			if plugin.RouterMount.UsePattern != "" {
				return true
			}
		}
		return false
	}
	
	plugin, err := plugins.GetPlugin(rp.plugins, framework)
	if err != nil {
		return false
	}
	return plugin != nil && plugin.RouterMount.UsePattern != ""
}

// BuildRouterMap scans files to build a map of router files to their base paths
func (rp *RouteParser) BuildRouterMap(files []string, projectRoot string) {
	// Find plugin with router mount support
	var plugin *plugins.Plugin
	if rp.framework == "auto" {
		// Find first plugin with router mount support
		for _, p := range rp.plugins {
			if p.RouterMount.UsePattern != "" {
				plugin = p
				break
			}
		}
	} else {
		var err error
		plugin, err = plugins.GetPlugin(rp.plugins, rp.framework)
		if err != nil {
			// Plugin not found, will fall back to hardcoded patterns
			plugin = nil
		}
	}

	// Fallback to hardcoded Express patterns if no plugin
	if plugin == nil || plugin.RouterMount.UsePattern == "" {
		if rp.framework != "auto" && rp.framework != "express" {
			return
		}
		// Use hardcoded Express patterns
		usePattern := regexp.MustCompile(`app\.use\s*\(\s*['"]([^'"]+)['"]\s*,\s*(\w+)`)
		requirePattern := regexp.MustCompile(`require\s*\(['"]([^'"]+)['"]\)`)
		rp.buildRouterMapWithPatterns(files, usePattern, requirePattern, regexp.MustCompile(`(?:const|let|var)\s+(\w+)\s*=.*require`))
		return
	}

	// Use plugin patterns
	usePattern := regexp.MustCompile(plugin.RouterMount.UsePattern)
	requirePattern := regexp.MustCompile(plugin.RouterMount.RequirePattern)
	varPattern := regexp.MustCompile(plugin.RouterMount.VarPattern)
	rp.buildRouterMapWithPatterns(files, usePattern, requirePattern, varPattern)
}

// buildRouterMapWithPatterns builds router map using provided patterns
func (rp *RouteParser) buildRouterMapWithPatterns(files []string, usePattern, requirePattern, varPattern *regexp.Regexp) {

	for _, filePath := range files {
		ext := filepath.Ext(filePath)
		if ext != ".js" && ext != ".ts" {
			continue
		}

		file, err := os.Open(filePath)
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(file)
		var routerVars map[string]string // Maps variable name to file path

		// First pass: find require statements for routers
		for scanner.Scan() {
			line := scanner.Text()
			// Pattern: const usersRouter = require('./routes/users')
			requireMatch := requirePattern.FindStringSubmatch(line)
			if requireMatch != nil && len(requireMatch) >= 2 {
				// Extract variable name using varPattern
				varNameMatch := varPattern.FindStringSubmatch(line)
				if varNameMatch != nil && len(varNameMatch) >= 2 {
					if routerVars == nil {
						routerVars = make(map[string]string)
					}
					routerVarName := varNameMatch[1]
					requirePath := requireMatch[1]

					// Resolve relative path to absolute
					routerFile := rp.resolveRouterPath(filePath, requirePath)
					routerVars[routerVarName] = routerFile
				}
			}
		}
		file.Close()

		// Second pass: find app.use() calls
		file, err = os.Open(filePath)
		if err != nil {
			continue
		}
		scanner = bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()

			useMatch := usePattern.FindStringSubmatch(line)
			if useMatch != nil && len(useMatch) >= 3 {
				basePath := useMatch[1]
				routerVarName := useMatch[2]

				// Look up the router file
				if routerFile, ok := routerVars[routerVarName]; ok {
					// Normalize paths for comparison
					normalizedRouterFile := filepath.ToSlash(routerFile)
					rp.routerBasePath[normalizedRouterFile] = basePath
				}
			}
		}
		file.Close()
	}
}

// resolveRouterPath resolves a relative require path to an absolute file path
func (rp *RouteParser) resolveRouterPath(baseFile, requirePath string) string {
	baseDir := filepath.Dir(baseFile)

	// Remove leading ./ if present
	requirePath = strings.TrimPrefix(requirePath, "./")

	// Resolve to absolute path
	absPath := filepath.Join(baseDir, requirePath)

	// Try with .js extension if no extension and file doesn't exist
	if !strings.HasSuffix(absPath, ".js") && !strings.HasSuffix(absPath, ".ts") {
		// Check if file exists without extension first
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			// Try with .js extension
			if _, err := os.Stat(absPath + ".js"); err == nil {
				absPath += ".js"
			}
		}
	}

	// Convert to absolute and normalize
	absPath, _ = filepath.Abs(absPath)
	return filepath.ToSlash(absPath)
}

// ParseFile parses routes from a single file
func (rp *RouteParser) ParseFile(filePath string) ([]Route, error) {
	ext := filepath.Ext(filePath)
	var routes []Route

	// For specific frameworks, prefer hardcoded parsers (they handle edge cases better)
	// For "auto" mode, try plugins first
	useHardcoded := false
	switch ext {
	case ".java":
		if rp.framework == "spring" {
			useHardcoded = true
			routes = rp.parseSpring(filePath)
		}
	case ".js", ".ts":
		if rp.framework == "express" {
			useHardcoded = true
			routes = rp.parseExpress(filePath)
		}
	case ".php":
		if rp.framework == "laravel" {
			useHardcoded = true
			routes = rp.parseLaravel(filePath)
		}
	}

	// Try to use plugins if not using hardcoded parser
	if !useHardcoded && len(rp.plugins) > 0 {
		routes = rp.parseWithPlugins(filePath, ext)
	}

	// Fallback to hardcoded parsers if no plugins or no routes found
	if len(routes) == 0 && !useHardcoded {
		switch ext {
		case ".js", ".ts":
			if rp.framework == "auto" || rp.framework == "express" {
				routes = rp.parseExpress(filePath)
			}
		case ".py":
			if rp.framework == "auto" || rp.framework == "flask" || rp.framework == "django" {
				routes = append(routes, rp.parseFlask(filePath)...)
				routes = append(routes, rp.parseDjango(filePath)...)
			}
		case ".java":
			if rp.framework == "auto" || rp.framework == "spring" {
				routes = rp.parseSpring(filePath)
			}
		case ".php":
			if rp.framework == "auto" || rp.framework == "laravel" {
				routes = rp.parseLaravel(filePath)
			}
		}
	}

	// Apply base path prefix for router files
	normalizedPath := filepath.ToSlash(filePath)
	// Try exact match first
	basePath, ok := rp.routerBasePath[normalizedPath]
	if !ok {
		// Try with different path separators and case variations
		for routerFile, bp := range rp.routerBasePath {
			// Compare normalized paths
			if strings.EqualFold(filepath.ToSlash(routerFile), normalizedPath) ||
				strings.EqualFold(routerFile, normalizedPath) {
				basePath = bp
				ok = true
				break
			}
		}
	}

	if ok {
		for i := range routes {
			// Prepend base path, ensuring proper path joining
			if routes[i].Path == "/" {
				routes[i].Path = basePath
			} else {
				// Ensure base path ends with / or route starts with /
				if !strings.HasSuffix(basePath, "/") && !strings.HasPrefix(routes[i].Path, "/") {
					routes[i].Path = basePath + "/" + routes[i].Path
				} else {
					routes[i].Path = basePath + routes[i].Path
				}
			}
		}
	}

	return routes, nil
}

// parseWithPlugins parses routes using plugin patterns
func (rp *RouteParser) parseWithPlugins(filePath string, ext string) []Route {
	var allRoutes []Route

	// Find matching plugins for this file extension
	var matchingPlugins []*plugins.Plugin
	if rp.framework == "auto" {
		// Try all plugins that match the extension
		for _, plugin := range rp.plugins {
			for _, pluginExt := range plugin.Extensions {
				if ext == pluginExt {
					matchingPlugins = append(matchingPlugins, plugin)
					break
				}
			}
		}
	} else {
		// Use specific framework plugin
		plugin, err := plugins.GetPlugin(rp.plugins, rp.framework)
		if err == nil && plugin != nil {
			for _, pluginExt := range plugin.Extensions {
				if ext == pluginExt {
					matchingPlugins = append(matchingPlugins, plugin)
					break
				}
			}
		}
	}

	// Parse with each matching plugin
	for _, plugin := range matchingPlugins {
		for _, pattern := range plugin.Patterns {
			routes := rp.parseWithPattern(filePath, pattern)
			allRoutes = append(allRoutes, routes...)
		}
	}

	return allRoutes
}

// parseWithPattern parses routes using a specific pattern
func (rp *RouteParser) parseWithPattern(filePath string, pattern plugins.Pattern) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	routeRegex := regexp.MustCompile(pattern.RouteRegex)
	var paramRegex *regexp.Regexp
	if pattern.ParamRegex != "" {
		paramRegex = regexp.MustCompile(pattern.ParamRegex)
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		matches := routeRegex.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			// Extract method (first capture group)
			method := strings.ToUpper(match[1])
			if len(match) >= 3 {
				// Path is in second capture group
				path := match[2]

				// Extract path parameters
				var params []string
				if paramRegex != nil {
					paramMatches := paramRegex.FindAllStringSubmatch(path, -1)
					for _, pm := range paramMatches {
						if len(pm) >= 2 {
							params = append(params, pm[1])
						}
					}
				}

				// Handle method patterns like "GET|POST"
				methods := []string{method}
				if strings.Contains(pattern.Method, "|") {
					// Check if method matches any in pattern
					methodParts := strings.Split(pattern.Method, "|")
					found := false
					for _, mp := range methodParts {
						if strings.EqualFold(method, mp) {
							found = true
							break
						}
					}
					if !found {
						continue
					}
				}

				// Check for multiple methods in query regex (Flask style)
				if pattern.QueryRegex != "" && strings.Contains(line, "methods") {
					methodRegex := regexp.MustCompile(pattern.QueryRegex)
					methodMatch := methodRegex.FindStringSubmatch(line)
					if methodMatch != nil && len(methodMatch) >= 2 {
						methodsStr := methodMatch[1]
						methods = []string{}
						for _, m := range strings.Split(methodsStr, ",") {
							m = strings.Trim(strings.TrimSpace(m), "'\"")
							if m != "" {
								methods = append(methods, strings.ToUpper(m))
							}
						}
					}
				}

				// Create routes for each method
				for _, m := range methods {
					routes = append(routes, Route{
						Path:   path,
						Method: m,
						Params: params,
						Query:  rp.detectQueryParams(filePath, lineNum),
						File:   filePath,
						Line:   lineNum,
					})
				}
			}
		}
	}

	return routes
}

// parseExpress parses Express.js routes
func (rp *RouteParser) parseExpress(filePath string) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Patterns for Express routes
	methodPattern := regexp.MustCompile(`(?:app|router|express)\.(get|post|put|delete|patch|all)\s*\(\s*['"]([^'"]+)['"]`)
	paramPattern := regexp.MustCompile(`:(\w+)`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		matches := methodPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				method := strings.ToUpper(match[1])
				path := match[2]

				// Extract path parameters
				paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
				var params []string
				for _, pm := range paramMatches {
					if len(pm) >= 2 {
						params = append(params, pm[1])
					}
				}

				// Try to detect query parameters (heuristic: look for req.query usage)
				query := rp.detectQueryParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  query,
					File:   filePath,
					Line:   lineNum,
				})
			}
		}
	}

	return routes
}

// parseFlask parses Flask routes
func (rp *RouteParser) parseFlask(filePath string) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Pattern for Flask routes: @app.route('/path', methods=['GET', 'POST'])
	routePattern := regexp.MustCompile(`@(?:app|blueprint)\.route\s*\(\s*['"]([^'"]+)['"]`)
	methodPattern := regexp.MustCompile(`methods\s*=\s*\[([^\]]+)\]`)
	paramPattern := regexp.MustCompile(`<(\w+)(?::[^>]+)?>`)

	var currentRoute *Route

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for route decorator
		routeMatch := routePattern.FindStringSubmatch(line)
		if routeMatch != nil && len(routeMatch) >= 2 {
			path := routeMatch[1]

			// Extract path parameters
			paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
			var params []string
			for _, pm := range paramMatches {
				if len(pm) >= 2 {
					params = append(params, pm[1])
				}
			}

			// Check for methods
			methodMatch := methodPattern.FindStringSubmatch(line)
			methods := []string{"GET"} // default
			if methodMatch != nil && len(methodMatch) >= 2 {
				methodsStr := methodMatch[1]
				methods = []string{}
				for _, m := range strings.Split(methodsStr, ",") {
					m = strings.Trim(strings.TrimSpace(m), "'\"")
					if m != "" {
						methods = append(methods, strings.ToUpper(m))
					}
				}
			}

			// Create routes for each method
			for _, method := range methods {
				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  rp.detectQueryParams(filePath, lineNum),
					File:   filePath,
					Line:   lineNum,
				})
			}
			currentRoute = nil
		} else if currentRoute != nil {
			// Look for request.args or request.form usage
			if strings.Contains(line, "request.args") || strings.Contains(line, "request.form") {
				// Heuristic: try to extract parameter names
				argPattern := regexp.MustCompile(`(?:args|form)\[['"](\w+)['"]\]`)
				argMatches := argPattern.FindAllStringSubmatch(line, -1)
				for _, am := range argMatches {
					if len(am) >= 2 {
						currentRoute.Query = append(currentRoute.Query, am[1])
					}
				}
			}
		}
	}

	return routes
}

// parseDjango parses Django URL patterns
func (rp *RouteParser) parseDjango(filePath string) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	// Only parse urls.py files
	if !strings.HasSuffix(filepath.Base(filePath), "urls.py") {
		return routes
	}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Pattern for Django: path('users/<int:id>/', views.user_detail)
	pathPattern := regexp.MustCompile(`(?:path|re_path)\s*\(\s*['"]([^'"]+)['"]`)
	paramPattern := regexp.MustCompile(`<(\w+)(?::[^>]+)?>`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		pathMatch := pathPattern.FindStringSubmatch(line)
		if pathMatch != nil && len(pathMatch) >= 2 {
			path := pathMatch[1]

			// Extract path parameters
			paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
			var params []string
			for _, pm := range paramMatches {
				if len(pm) >= 2 {
					params = append(params, pm[1])
				}
			}

			// Django doesn't specify method in urls.py, so we'll default to GET
			// In practice, views handle methods
			routes = append(routes, Route{
				Path:   path,
				Method: "GET", // Default, could be enhanced to check views
				Params: params,
				Query:  []string{},
				File:   filePath,
				Line:   lineNum,
			})
		}
	}

	return routes
}

// parseSpring parses Spring Boot routes
func (rp *RouteParser) parseSpring(filePath string) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	var classBasePath string = ""
	
	// Pattern for class-level @RequestMapping
	classMappingPattern := regexp.MustCompile(`@RequestMapping\s*\([^)]*['"]([^'"]+)['"]`)
	// Pattern for method-level annotations with path
	methodWithPathPattern := regexp.MustCompile(`@(Get|Post|Put|Delete|Patch|RequestMapping)\w*\s*\([^)]*['"]([^'"]+)['"]`)
	// Pattern for method-level annotations without path (may or may not have empty parentheses)
	methodWithoutPathPattern := regexp.MustCompile(`@(Get|Post|Put|Delete|Patch)Mapping\s*(?:\([^)]*\))?`)
	paramPattern := regexp.MustCompile(`\{(\w+)\}`)

	// First pass: find class-level @RequestMapping
	for scanner.Scan() {
		line := scanner.Text()
		classMatch := classMappingPattern.FindStringSubmatch(line)
		if classMatch != nil && len(classMatch) >= 2 {
			classBasePath = classMatch[1]
			break
		}
	}
	file.Close()

	// Second pass: find method-level annotations
	file, err = os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()
	scanner = bufio.NewScanner(file)
	lineNum = 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for method annotation with path
		matches := methodWithPathPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				method := strings.ToUpper(strings.TrimSuffix(match[1], "Mapping"))
				if method == "REQUEST" {
					method = "GET" // Default for @RequestMapping
				}
				path := match[2]

				// Combine with class base path if exists
				if classBasePath != "" {
					if strings.HasPrefix(path, "/") {
						path = classBasePath + path
					} else {
						path = classBasePath + "/" + path
					}
				}

				// Extract path parameters
				paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
				var params []string
				for _, pm := range paramMatches {
					if len(pm) >= 2 {
						params = append(params, pm[1])
					}
				}

				// Detect query parameters for Spring
				queryParams := rp.detectSpringQueryParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  queryParams,
					File:   filePath,
					Line:   lineNum,
				})
			}
		}

		// Check for method annotation without path (uses class base path)
		noPathMatches := methodWithoutPathPattern.FindAllStringSubmatch(line, -1)
		for _, match := range noPathMatches {
			if len(match) >= 2 && classBasePath != "" {
				method := strings.ToUpper(strings.TrimSuffix(match[1], "Mapping"))
				
				// Use class base path
				path := classBasePath

				// Extract path parameters from base path
				paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
				var params []string
				for _, pm := range paramMatches {
					if len(pm) >= 2 {
						params = append(params, pm[1])
					}
				}

				// Detect query parameters for Spring
				queryParams := rp.detectSpringQueryParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  queryParams,
					File:   filePath,
					Line:   lineNum,
				})
			}
		}
	}

	return routes
}

// detectSpringQueryParams detects Spring @RequestParam annotations
func (rp *RouteParser) detectSpringQueryParams(filePath string, lineNum int) []string {
	var queryParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return queryParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	// Pattern: @RequestParam (optional attributes) type paramName
	// Example: @RequestParam(required = false) Integer page
	requestParamPattern := regexp.MustCompile(`@RequestParam\s*(?:\([^)]*\))?\s+\w+\s+(\w+)`)

	// Scan from the annotation line to the end of the method (find closing brace)
	for scanner.Scan() {
		currentLine++
		if currentLine < lineNum {
			continue
		}
		
		line := scanner.Text()
		
		// Look for @RequestParam annotations
		matches := requestParamPattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				paramName := match[1] // The parameter name
				queryParams = append(queryParams, paramName)
			}
		}
		
		// Stop at method closing brace (heuristic: line with just "}" or "},")
		if strings.TrimSpace(line) == "}" || strings.TrimSpace(line) == "}," {
			break
		}
	}

	return queryParams
}

// parseLaravel parses Laravel routes
func (rp *RouteParser) parseLaravel(filePath string) []Route {
	var routes []Route
	file, err := os.Open(filePath)
	if err != nil {
		return routes
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	// Pattern for Laravel: Route::get('/users/{id}', ...)
	routePattern := regexp.MustCompile(`Route::(get|post|put|delete|patch|any)\s*\(\s*['"]([^'"]+)['"]`)
	paramPattern := regexp.MustCompile(`\{(\w+)\}`)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		matches := routePattern.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				method := strings.ToUpper(match[1])
				path := match[2]

				// Extract path parameters
				paramMatches := paramPattern.FindAllStringSubmatch(path, -1)
				var params []string
				for _, pm := range paramMatches {
					if len(pm) >= 2 {
						params = append(params, pm[1])
					}
				}

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  []string{},
					File:   filePath,
					Line:   lineNum,
				})
			}
		}
	}

	return routes
}

// detectQueryParams attempts to detect query parameters (heuristic)
func (rp *RouteParser) detectQueryParams(filePath string, lineNum int) []string {
	// This is a simplified heuristic - in a real implementation,
	// you'd analyze more context around the route definition
	return []string{}
}

