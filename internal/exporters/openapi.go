package exporters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/parser"
)

type OpenAPIDocument struct {
	OpenAPI    string                `json:"openapi" yaml:"openapi"`
	Info       OpenAPIInfo           `json:"info" yaml:"info"`
	Servers    []OpenAPIServer       `json:"servers" yaml:"servers"`
	Paths      OpenAPIPaths          `json:"paths" yaml:"paths"`
	Components OpenAPIComponents     `json:"components,omitempty" yaml:"components,omitempty"`
}

type OpenAPIInfo struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`
	Version     string `json:"version" yaml:"version"`
}

type OpenAPIServer struct {
	URL         string `json:"url" yaml:"url"`
	Description string `json:"description" yaml:"description"`
}

type OpenAPIPaths map[string]OpenAPIPathItem

type OpenAPIPathItem struct {
	Get    *OpenAPIOperation `json:"get,omitempty" yaml:"get,omitempty"`
	Post   *OpenAPIOperation `json:"post,omitempty" yaml:"post,omitempty"`
	Put    *OpenAPIOperation `json:"put,omitempty" yaml:"put,omitempty"`
	Delete *OpenAPIOperation `json:"delete,omitempty" yaml:"delete,omitempty"`
	Patch  *OpenAPIOperation `json:"patch,omitempty" yaml:"patch,omitempty"`
}

type OpenAPIOperation struct {
	Summary    string                 `json:"summary" yaml:"summary"`
	Parameters []OpenAPIParameter     `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *OpenAPIRequestBody   `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses  OpenAPIResponses       `json:"responses" yaml:"responses"`
}

type OpenAPIParameter struct {
	Name        string      `json:"name" yaml:"name"`
	In          string      `json:"in" yaml:"in"`
	Required    bool        `json:"required" yaml:"required"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Schema      OpenAPISchema `json:"schema" yaml:"schema"`
}

type OpenAPIRequestBody struct {
	Required bool                    `json:"required" yaml:"required"`
	Content  map[string]OpenAPIMedia `json:"content" yaml:"content"`
}

type OpenAPIMedia struct {
	Schema OpenAPISchema `json:"schema" yaml:"schema"`
}

type OpenAPIResponses map[string]OpenAPIResponse

type OpenAPIResponse struct {
	Description string `json:"description" yaml:"description"`
}

type OpenAPISchema struct {
	Type       string                  `json:"type,omitempty" yaml:"type,omitempty"`
	Properties map[string]OpenAPISchema `json:"properties,omitempty" yaml:"properties,omitempty"`
	Items      *OpenAPISchema           `json:"items,omitempty" yaml:"items,omitempty"`
}

type OpenAPIComponents struct {
	Schemas map[string]OpenAPISchema `json:"schemas,omitempty" yaml:"schemas,omitempty"`
}

func GenerateOpenAPIExport(routes []parser.Route, projectPath string, env string) OpenAPIDocument {
	envName := getEnvironmentName(env)
	baseURL := getBaseURL(env)

	doc := OpenAPIDocument{
		OpenAPI: "3.0.3",
		Info: OpenAPIInfo{
			Title:       fmt.Sprintf("%s API", filepath.Base(projectPath)),
			Description: fmt.Sprintf("Auto-generated OpenAPI specification from %s - %s environment", projectPath, envName),
			Version:     "1.0.0",
		},
		Servers: []OpenAPIServer{
			{
				URL:         baseURL,
				Description: envName,
			},
		},
		Paths: make(OpenAPIPaths),
	}

	for _, route := range routes {
		normalizedPath := normalizePathForOpenAPI(route.Path)
		pathItem := doc.Paths[normalizedPath]

		operation := buildOpenAPIOperation(route)

		switch strings.ToUpper(route.Method) {
		case "GET":
			pathItem.Get = &operation
		case "POST":
			pathItem.Post = &operation
		case "PUT":
			pathItem.Put = &operation
		case "DELETE":
			pathItem.Delete = &operation
		case "PATCH":
			pathItem.Patch = &operation
		default:
			if pathItem.Get == nil {
				pathItem.Get = &operation
			}
		}

		doc.Paths[normalizedPath] = pathItem
	}

	return doc
}

func normalizePathForOpenAPI(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	var normalized []string

	for _, part := range parts {
		if part == "" {
			continue
		}
		if strings.HasPrefix(part, ":") {
			normalized = append(normalized, "{"+part[1:]+"}")
		} else if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			paramName := strings.Trim(part, "<>")
			if idx := strings.Index(paramName, ":"); idx != -1 {
				paramName = paramName[idx+1:]
			}
			normalized = append(normalized, "{"+paramName+"}")
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			normalized = append(normalized, part)
		} else {
			normalized = append(normalized, part)
		}
	}

	return "/" + strings.Join(normalized, "/")
}

func buildOpenAPIOperation(route parser.Route) OpenAPIOperation {
	operation := OpenAPIOperation{
		Summary: fmt.Sprintf("%s %s", route.Method, route.Path),
		Responses: OpenAPIResponses{
			"200": OpenAPIResponse{
				Description: "Successful response",
			},
		},
	}

	var parameters []OpenAPIParameter

	pathParams := extractPathParams(route.Path)
	for _, param := range pathParams {
		parameters = append(parameters, OpenAPIParameter{
			Name:        param,
			In:          "path",
			Required:    true,
			Description: fmt.Sprintf("Path parameter: %s", param),
			Schema: OpenAPISchema{
				Type: "string",
			},
		})
	}

	for _, qp := range route.Query {
		parameters = append(parameters, OpenAPIParameter{
			Name:        qp,
			In:          "query",
			Required:    false,
			Description: fmt.Sprintf("Query parameter: %s", qp),
			Schema: OpenAPISchema{
				Type: "string",
			},
		})
	}

	if len(parameters) > 0 {
		operation.Parameters = parameters
	}

	method := strings.ToUpper(route.Method)
	if method == "POST" || method == "PUT" || method == "PATCH" {
		operation.RequestBody = buildOpenAPIRequestBody(route)
	}

	return operation
}

func extractPathParams(path string) []string {
	var params []string
	parts := strings.Split(strings.Trim(path, "/"), "/")
	seen := make(map[string]bool)

	for _, part := range parts {
		var paramName string
		if strings.HasPrefix(part, ":") {
			paramName = part[1:]
		} else if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			paramName = strings.Trim(part, "<>")
			if idx := strings.Index(paramName, ":"); idx != -1 {
				paramName = paramName[idx+1:]
			}
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			paramName = strings.Trim(part, "{}")
		}
		if paramName != "" && !seen[paramName] {
			params = append(params, paramName)
			seen[paramName] = true
		}
	}
	return params
}

func buildOpenAPIRequestBody(route parser.Route) *OpenAPIRequestBody {
	if len(route.Body) == 0 {
		return &OpenAPIRequestBody{
			Required: true,
			Content: map[string]OpenAPIMedia{
				"application/json": {
					Schema: OpenAPISchema{
						Type: "object",
					},
				},
			},
		}
	}

	properties := make(map[string]OpenAPISchema)
	for _, field := range route.Body {
		properties[field] = OpenAPISchema{
			Type: "string",
		}
	}

	return &OpenAPIRequestBody{
		Required: true,
		Content: map[string]OpenAPIMedia{
			"application/json": {
				Schema: OpenAPISchema{
					Type:       "object",
					Properties: properties,
				},
			},
		},
	}
}
