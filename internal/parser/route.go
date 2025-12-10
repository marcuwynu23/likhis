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

	// Apply ignore patterns from plugins
	if len(routes) > 0 {
		var ignorePatterns []string
		
		if rp.framework == "auto" {
			// Collect ignore patterns from all plugins
			for _, plugin := range rp.plugins {
				if len(plugin.Ignore) > 0 {
					ignorePatterns = append(ignorePatterns, plugin.Ignore...)
				}
			}
		} else {
			// Use ignore patterns from the specific framework plugin
			plugin, err := plugins.GetPlugin(rp.plugins, rp.framework)
			if err == nil && plugin != nil && len(plugin.Ignore) > 0 {
				ignorePatterns = plugin.Ignore
			}
		}
		
		// Apply filtering if we have any ignore patterns
		if len(ignorePatterns) > 0 {
			routes = rp.filterIgnoredRoutes(routes, ignorePatterns)
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

	// Apply ignore patterns if any plugin has them
	for _, plugin := range matchingPlugins {
		if len(plugin.Ignore) > 0 {
			allRoutes = rp.filterIgnoredRoutes(allRoutes, plugin.Ignore)
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

				// Detect query and body parameters
				query := rp.detectQueryParams(filePath, lineNum)
				body := rp.detectBodyParams(filePath, lineNum)

				// Create routes for each method
				for _, m := range methods {
					routes = append(routes, Route{
						Path:   path,
						Method: m,
						Params: params,
						Query:  query,
						Body:   body,
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

				// Detect query and body parameters
				query := rp.detectQueryParams(filePath, lineNum)
				body := rp.detectBodyParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  query,
					Body:   body,
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

			// Detect query and body parameters for Flask
			query := rp.detectFlaskQueryParams(filePath, lineNum)
			body := rp.detectFlaskBodyParams(filePath, lineNum)

			// Create routes for each method
			for _, method := range methods {
				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  query,
					Body:   body,
					File:   filePath,
					Line:   lineNum,
				})
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

				// Detect query and body parameters for Spring
				queryParams := rp.detectSpringQueryParams(filePath, lineNum)
				bodyParams := rp.detectSpringBodyParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  queryParams,
					Body:   bodyParams,
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

				// Detect query and body parameters for Spring
				queryParams := rp.detectSpringQueryParams(filePath, lineNum)
				bodyParams := rp.detectSpringBodyParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  queryParams,
					Body:   bodyParams,
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

				// Detect query and body parameters for Laravel
				query := rp.detectLaravelQueryParams(filePath, lineNum)
				body := rp.detectLaravelBodyParams(filePath, lineNum)

				routes = append(routes, Route{
					Path:   path,
					Method: method,
					Params: params,
					Query:  query,
					Body:   body,
					File:   filePath,
					Line:   lineNum,
				})
			}
		}
	}

	return routes
}

// detectQueryParams detects query parameters by scanning code around the route definition
// Supports Express.js patterns: req.query.param, req.query['param'], const { param } = req.query
func (rp *RouteParser) detectQueryParams(filePath string, lineNum int) []string {
	var queryParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return queryParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 50 // Scan 50 lines after route definition to catch handler functions
	
	// Track if we're inside a handler function
	inHandler := false
	braceCount := 0
	
	// Patterns for Express.js query parameters
	queryPatterns := []*regexp.Regexp{
		regexp.MustCompile(`req\.query\.(\w+)`),                    // req.query.param
		regexp.MustCompile(`req\.query\[['"](\w+)['"]\]`),         // req.query['param']
		regexp.MustCompile(`req\.query\[(\w+)\]`),                 // req.query[param]
		regexp.MustCompile(`const\s*\{([^}]+)\}\s*=\s*req\.query`), // const { param1, param2 } = req.query
		regexp.MustCompile(`let\s*\{([^}]+)\}\s*=\s*req\.query`),  // let { param1, param2 } = req.query
		regexp.MustCompile(`var\s*\{([^}]+)\}\s*=\s*req\.query`),  // var { param1, param2 } = req.query
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		
		// Detect handler function start: arrow function (req, res) => or function(req, res)
		if strings.Contains(line, "=>") || (strings.Contains(line, "function") && strings.Contains(line, "req")) {
			inHandler = true
			braceCount = 0
		}
		
		// Track braces to detect function end
		if inHandler {
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			if braceCount < 0 {
				braceCount = 0
			}
			// If we hit a closing brace and brace count is back to 0, we've left the handler
			if strings.Contains(line, "}") && braceCount == 0 && strings.Count(line, "}") > 0 {
				// Check if this is the end of the route handler (look for closing paren or comma)
				if strings.Contains(line, ")") || strings.Contains(line, ",") {
					inHandler = false
				}
			}
		}
		
		// Only check for query params if we're in a handler function or within first few lines
		if inHandler || currentLine <= startLine + 5 {
			// Check each pattern
			for _, pattern := range queryPatterns {
				matches := pattern.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if len(match) >= 2 {
						// Handle destructuring: const { param1, param2 } = req.query
						if strings.Contains(match[1], ",") {
							params := strings.Split(match[1], ",")
							for _, p := range params {
								p = strings.TrimSpace(p)
								// Remove default value if present: param = defaultValue
								if idx := strings.Index(p, "="); idx > 0 {
									p = strings.TrimSpace(p[:idx])
								}
								if p != "" && !contains(queryParams, p) {
									queryParams = append(queryParams, p)
								}
							}
						} else {
							param := strings.TrimSpace(match[1])
							if param != "" && !contains(queryParams, param) {
								queryParams = append(queryParams, param)
							}
						}
					}
				}
			}
		}
	}

	return queryParams
}

// detectBodyParams detects body parameters by scanning code around the route definition
// Supports Express.js patterns: req.body.param, req.body['param'], const { param } = req.body
func (rp *RouteParser) detectBodyParams(filePath string, lineNum int) []string {
	var bodyParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return bodyParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 50 // Scan 50 lines after route definition to catch handler functions (including arrow functions)
	
	// Track if we're inside a handler function
	inHandler := false
	braceCount := 0
	
	// Patterns for Express.js body parameters
	bodyPatterns := []*regexp.Regexp{
		regexp.MustCompile(`req\.body\.(\w+)`),                    // req.body.param
		regexp.MustCompile(`req\.body\[['"](\w+)['"]\]`),         // req.body['param']
		regexp.MustCompile(`req\.body\[(\w+)\]`),                 // req.body[param]
		regexp.MustCompile(`const\s*\{([^}]+)\}\s*=\s*req\.body`), // const { param1, param2 } = req.body
		regexp.MustCompile(`let\s*\{([^}]+)\}\s*=\s*req\.body`),  // let { param1, param2 } = req.body
		regexp.MustCompile(`var\s*\{([^}]+)\}\s*=\s*req\.body`),  // var { param1, param2 } = req.body
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		
		// Detect handler function start: arrow function (req, res) => or function(req, res)
		if strings.Contains(line, "=>") || (strings.Contains(line, "function") && strings.Contains(line, "req")) {
			inHandler = true
			braceCount = 0
		}
		
		// Track braces to detect function end
		if inHandler {
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			if braceCount < 0 {
				braceCount = 0
			}
			// If we hit a closing brace and brace count is back to 0, we've left the handler
			if strings.Contains(line, "}") && braceCount == 0 && strings.Count(line, "}") > 0 {
				// Check if this is the end of the route handler (look for closing paren or comma)
				if strings.Contains(line, ")") || strings.Contains(line, ",") {
					inHandler = false
				}
			}
		}
		
		// Only check for body params if we're in a handler function or within first few lines
		if inHandler || currentLine <= startLine + 5 {
			// Check each pattern
			for _, pattern := range bodyPatterns {
				matches := pattern.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if len(match) >= 2 {
						// Handle destructuring: const { param1, param2 } = req.body
						if strings.Contains(match[1], ",") {
							params := strings.Split(match[1], ",")
							for _, p := range params {
								p = strings.TrimSpace(p)
								// Remove default value if present: param = defaultValue
								if idx := strings.Index(p, "="); idx > 0 {
									p = strings.TrimSpace(p[:idx])
								}
								if p != "" && !contains(bodyParams, p) {
									bodyParams = append(bodyParams, p)
								}
							}
						} else {
							param := strings.TrimSpace(match[1])
							if param != "" && !contains(bodyParams, param) {
								bodyParams = append(bodyParams, param)
							}
						}
					}
				}
			}
		}
	}

	return bodyParams
}

// detectFlaskQueryParams detects Flask query parameters
// Supports: request.args.get('param'), request.args['param']
func (rp *RouteParser) detectFlaskQueryParams(filePath string, lineNum int) []string {
	var queryParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return queryParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`request\.args\.get\(['"](\w+)['"]`),
		regexp.MustCompile(`request\.args\[['"](\w+)['"]\]`),
		regexp.MustCompile(`request\.args\.get\((\w+)`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(queryParams, param) {
						queryParams = append(queryParams, param)
					}
				}
			}
		}
	}

	return queryParams
}

// detectFlaskBodyParams detects Flask body parameters
// Supports: request.form.get('param'), request.json.get('param'), request.get_json()
func (rp *RouteParser) detectFlaskBodyParams(filePath string, lineNum int) []string {
	var bodyParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return bodyParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`request\.form\.get\(['"](\w+)['"]`),
		regexp.MustCompile(`request\.json\.get\(['"](\w+)['"]`),
		regexp.MustCompile(`request\.form\[['"](\w+)['"]\]`),
		regexp.MustCompile(`request\.json\[['"](\w+)['"]\]`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(bodyParams, param) {
						bodyParams = append(bodyParams, param)
					}
				}
			}
		}
	}

	return bodyParams
}

// detectDjangoQueryParams detects Django query parameters
// Supports: request.GET.get('param'), request.GET['param']
func (rp *RouteParser) detectDjangoQueryParams(filePath string, lineNum int) []string {
	var queryParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return queryParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`request\.GET\.get\(['"](\w+)['"]`),
		regexp.MustCompile(`request\.GET\[['"](\w+)['"]\]`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(queryParams, param) {
						queryParams = append(queryParams, param)
					}
				}
			}
		}
	}

	return queryParams
}

// detectDjangoBodyParams detects Django body parameters
// Supports: request.POST.get('param'), request.POST['param']
func (rp *RouteParser) detectDjangoBodyParams(filePath string, lineNum int) []string {
	var bodyParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return bodyParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`request\.POST\.get\(['"](\w+)['"]`),
		regexp.MustCompile(`request\.POST\[['"](\w+)['"]\]`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(bodyParams, param) {
						bodyParams = append(bodyParams, param)
					}
				}
			}
		}
	}

	return bodyParams
}

// detectLaravelQueryParams detects Laravel query parameters
// Supports: $request->query('param'), $request->input('param')
func (rp *RouteParser) detectLaravelQueryParams(filePath string, lineNum int) []string {
	var queryParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return queryParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\$request->query\(['"](\w+)['"]`),
		regexp.MustCompile(`\$request->input\(['"](\w+)['"]`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(queryParams, param) {
						queryParams = append(queryParams, param)
					}
				}
			}
		}
	}

	return queryParams
}

// detectLaravelBodyParams detects Laravel body parameters
// Supports: $request->input('param'), $request->json('param')
func (rp *RouteParser) detectLaravelBodyParams(filePath string, lineNum int) []string {
	var bodyParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return bodyParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 20

	patterns := []*regexp.Regexp{
		regexp.MustCompile(`\$request->input\(['"](\w+)['"]`),
		regexp.MustCompile(`\$request->json\(['"](\w+)['"]`),
	}

	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		for _, pattern := range patterns {
			matches := pattern.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				if len(match) >= 2 {
					param := strings.TrimSpace(match[1])
					if param != "" && !contains(bodyParams, param) {
						bodyParams = append(bodyParams, param)
					}
				}
			}
		}
	}

	return bodyParams
}

// detectSpringBodyParams detects Spring body parameters
// Supports: @RequestBody, method parameters
func (rp *RouteParser) detectSpringBodyParams(filePath string, lineNum int) []string {
	var bodyParams []string
	file, err := os.Open(filePath)
	if err != nil {
		return bodyParams
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentLine := 0
	startLine := lineNum
	endLine := lineNum + 30

	// Look for @RequestBody annotation and method parameters
	requestBodyPattern := regexp.MustCompile(`@RequestBody`)
	methodParamPattern := regexp.MustCompile(`@RequestBody\s+(?:\w+\s+)?(\w+)\s*[,)]`)

	inMethod := false
	for scanner.Scan() {
		currentLine++
		if currentLine < startLine {
			continue
		}
		if currentLine > endLine {
			break
		}

		line := scanner.Text()
		
		// Check for method signature start
		if strings.Contains(line, "public") && strings.Contains(line, "(") {
			inMethod = true
		}
		
		// Check for @RequestBody with parameter name
		if matches := methodParamPattern.FindStringSubmatch(line); len(matches) >= 2 {
			param := strings.TrimSpace(matches[1])
			if param != "" && !contains(bodyParams, param) {
				bodyParams = append(bodyParams, param)
			}
		}
		
		// Check for @RequestBody annotation (may have parameter on next line)
		if requestBodyPattern.MatchString(line) && inMethod {
			// Try to find parameter on same or next line
			if nextLine := scanner.Scan(); nextLine {
				currentLine++
				nextLineText := scanner.Text()
				paramMatch := regexp.MustCompile(`(\w+)\s*[,)]`).FindStringSubmatch(nextLineText)
				if len(paramMatch) >= 2 {
					param := strings.TrimSpace(paramMatch[1])
					if param != "" && !contains(bodyParams, param) {
						bodyParams = append(bodyParams, param)
					}
				}
			}
		}
		
		// Check for method end
		if strings.Contains(line, "}") && inMethod {
			inMethod = false
		}
	}

	return bodyParams
}

// filterIgnoredRoutes filters out routes that match any of the ignore patterns
// Patterns are checked against route path, query parameters, and body parameters
func (rp *RouteParser) filterIgnoredRoutes(routes []Route, ignorePatterns []string) []Route {
	if len(ignorePatterns) == 0 {
		return routes
	}

	// Compile all ignore patterns
	compiledPatterns := make([]*regexp.Regexp, 0, len(ignorePatterns))
	for _, pattern := range ignorePatterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			// Log error but continue with other patterns
			continue
		}
		compiledPatterns = append(compiledPatterns, compiled)
	}

	if len(compiledPatterns) == 0 {
		return routes
	}

	// Filter routes
	filtered := make([]Route, 0, len(routes))
	for _, route := range routes {
		shouldIgnore := false
		
		// Check route path
		for _, pattern := range compiledPatterns {
			if pattern.MatchString(route.Path) {
				shouldIgnore = true
				break
			}
		}
		
		// Check query parameters if path didn't match
		if !shouldIgnore {
			for _, pattern := range compiledPatterns {
				for _, queryParam := range route.Query {
					if pattern.MatchString(queryParam) {
						shouldIgnore = true
						break
					}
				}
				if shouldIgnore {
					break
				}
			}
		}
		
		// Check body parameters if path and query didn't match
		if !shouldIgnore {
			for _, pattern := range compiledPatterns {
				for _, bodyParam := range route.Body {
					if pattern.MatchString(bodyParam) {
						shouldIgnore = true
						break
					}
				}
				if shouldIgnore {
					break
				}
			}
		}
		
		if !shouldIgnore {
			filtered = append(filtered, route)
		}
	}

	return filtered
}

// contains checks if a string slice contains a value
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

