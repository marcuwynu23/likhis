package exporters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/parser"
)

// GenerateCURLScript generates a shell script with curl commands from routes
func GenerateCURLScript(routes []parser.Route, projectPath string, env string) string {
	var script strings.Builder
	envName := getEnvironmentName(env)
	baseURL := getBaseURL(env)

	script.WriteString("#!/bin/bash\n")
	script.WriteString("# Auto-generated CURL commands for ")
	script.WriteString(filepath.Base(projectPath))
	script.WriteString(fmt.Sprintf(" API - %s Environment\n", envName))
	script.WriteString("# Base URL for ")
	script.WriteString(envName)
	script.WriteString("\n")
	script.WriteString(fmt.Sprintf("BASE_URL=\"%s\"\n\n", baseURL))

	// Group routes by method for better organization
	routesByMethod := make(map[string][]parser.Route)
	for _, route := range routes {
		routesByMethod[route.Method] = append(routesByMethod[route.Method], route)
	}

	// Generate curl commands grouped by method
	for method, methodRoutes := range routesByMethod {
		script.WriteString(fmt.Sprintf("# %s Requests\n", method))
		script.WriteString("# " + strings.Repeat("=", 50) + "\n\n")

		for _, route := range methodRoutes {
			script.WriteString(fmt.Sprintf("# %s %s\n", method, route.Path))
			script.WriteString(createCURLCommand(route, method, "${BASE_URL}"))
			script.WriteString("\n\n")
		}
	}

	return script.String()
}

// createCURLCommand creates a curl command string from a route
func createCURLCommand(route parser.Route, method string, baseURL string) string {
	var cmd strings.Builder

	// Build URL
	url := baseURL + route.Path

	// Replace path parameters with example values
	for _, param := range route.Params {
		exampleValue := "1" // Default example value
		if strings.Contains(strings.ToLower(param), "id") {
			exampleValue = "1"
		} else if strings.Contains(strings.ToLower(param), "name") {
			exampleValue = "example"
		}
		url = strings.ReplaceAll(url, ":"+param, exampleValue)
		url = strings.ReplaceAll(url, "{"+param+"}", exampleValue)
		url = strings.ReplaceAll(url, "<"+param+">", exampleValue)
	}

	// Start curl command
	cmd.WriteString("curl -X ")
	cmd.WriteString(method)

	// Add headers
	cmd.WriteString(" \\\n  -H \"Content-Type: application/json\"")

	// Add query parameters
	if len(route.Query) > 0 {
		queryParts := []string{}
		for _, qp := range route.Query {
			queryParts = append(queryParts, fmt.Sprintf("%s=VALUE", qp))
		}
		if len(queryParts) > 0 {
			url += "?" + strings.Join(queryParts, "&")
		}
	}

	// Add body for POST, PUT, PATCH
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if len(route.Body) > 0 {
			// Build JSON body from body fields
			bodyFields := []string{}
			for _, field := range route.Body {
				bodyFields = append(bodyFields, fmt.Sprintf("\"%s\": \"VALUE\"", field))
			}
			body := "{" + strings.Join(bodyFields, ", ") + "}"
			cmd.WriteString(" \\\n  -d '")
			cmd.WriteString(body)
			cmd.WriteString("'")
		} else {
			cmd.WriteString(" \\\n  -d '{}'")
		}
	}

	// Add URL
	cmd.WriteString(" \\\n  \"")
	cmd.WriteString(url)
	cmd.WriteString("\"")

	return cmd.String()
}

// GenerateCURLMarkdown generates a markdown file with curl command examples
func GenerateCURLMarkdown(routes []parser.Route, projectPath string) string {
	var md strings.Builder

	md.WriteString("# API CURL Commands\n\n")
	md.WriteString(fmt.Sprintf("Auto-generated CURL commands for **%s** API\n\n", filepath.Base(projectPath)))
	md.WriteString("## Base URL\n\n")
	md.WriteString("```bash\n")
	md.WriteString("BASE_URL=\"http://localhost:3000\"\n")
	md.WriteString("```\n\n")

	// Group routes by method
	routesByMethod := make(map[string][]parser.Route)
	for _, route := range routes {
		routesByMethod[route.Method] = append(routesByMethod[route.Method], route)
	}

	// Generate curl commands grouped by method
	for method, methodRoutes := range routesByMethod {
		md.WriteString(fmt.Sprintf("## %s Requests\n\n", method))

		for _, route := range methodRoutes {
			md.WriteString(fmt.Sprintf("### %s %s\n\n", method, route.Path))
			md.WriteString("```bash\n")
			md.WriteString(createCURLCommand(route, method, "${BASE_URL}"))
			md.WriteString("\n```\n\n")
		}
	}

	return md.String()
}

