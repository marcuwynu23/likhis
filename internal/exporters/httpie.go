package exporters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/parser"
)

// HTTPieCollection represents an HTTPie Desktop collection structure
type HTTPieCollection struct {
	Meta  HTTPieMeta  `json:"meta"`
	Entry HTTPieEntry `json:"entry"`
}

// HTTPieMeta represents metadata for HTTPie collection
type HTTPieMeta struct {
	Format      string `json:"format"`
	Version     string `json:"version"`
	ContentType string `json:"contentType"`
	Schema      string `json:"schema"`
	Docs        string `json:"docs"`
	Source      string `json:"source"`
}

// HTTPieEntry represents the main entry/collection
type HTTPieEntry struct {
	Name  string         `json:"name"`
	Icon  HTTPieIcon     `json:"icon"`
	Auth  HTTPieAuth     `json:"auth"`
	Requests []HTTPieRequest `json:"requests"`
}

// HTTPieIcon represents icon information
type HTTPieIcon struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// HTTPieAuth represents authentication
type HTTPieAuth struct {
	Type string `json:"type"`
}

// HTTPieRequest represents a request in HTTPie Desktop
type HTTPieRequest struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Method      string            `json:"method"`
	Headers     []interface{}     `json:"headers"`
	QueryParams []interface{}     `json:"queryParams"`
	PathParams  []interface{}     `json:"pathParams"`
	Auth        HTTPieAuth        `json:"auth"`
	Body        HTTPieRequestBody `json:"body"`
}

// HTTPieRequestBody represents request body
type HTTPieRequestBody struct {
	Type    string          `json:"type"`
	File    HTTPieFile      `json:"file"`
	Text    HTTPieText      `json:"text"`
	Form    HTTPieForm      `json:"form"`
	GraphQL HTTPieGraphQL   `json:"graphql"`
}

// HTTPieFile represents file body
type HTTPieFile struct {
	Name string `json:"name"`
}

// HTTPieText represents text body
type HTTPieText struct {
	Value  string `json:"value"`
	Format string `json:"format"`
}

// HTTPieForm represents form body
type HTTPieForm struct {
	IsMultipart bool              `json:"isMultipart"`
	Fields      []HTTPieFormField `json:"fields"`
}

// HTTPieFormField represents a form field
type HTTPieFormField struct {
	Name    string `json:"name"`
	Value   string `json:"value"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

// HTTPieGraphQL represents GraphQL body
type HTTPieGraphQL struct {
	Query     string `json:"query"`
	Variables string `json:"variables"`
}

// HTTPieBody represents request body
type HTTPieBody struct {
	MimeType string                 `json:"mimeType"`
	Text     string                 `json:"text,omitempty"`
	Params   []HTTPieFormParam      `json:"params,omitempty"`
}

// HTTPieFormParam represents form parameter
type HTTPieFormParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GenerateHTTPieExport generates an HTTPie Desktop collection from routes
func GenerateHTTPieExport(routes []parser.Route, projectPath string, env string) HTTPieCollection {
	envName := getEnvironmentName(env)
	collectionName := fmt.Sprintf("%s API (%s)", filepath.Base(projectPath), envName)
	
	collection := HTTPieCollection{
		Meta: HTTPieMeta{
			Format:      "httpie",
			Version:     "1.0.0",
			ContentType: "collection",
			Schema:      "https://schema.httpie.io/1.0.0.json",
			Docs:        "https://httpie.io/r/help/export-from-httpie",
			Source:      "API Mapper Tool",
		},
		Entry: HTTPieEntry{
			Name:  collectionName,
			Icon: HTTPieIcon{
				Name:  "default",
				Color: "gray",
			},
			Auth: HTTPieAuth{
				Type: "none",
			},
			Requests: []HTTPieRequest{},
		},
	}

	// Create requests for all routes
	for _, route := range routes {
		request := createHTTPieRequest(route, route.Method, env)
		collection.Entry.Requests = append(collection.Entry.Requests, request)
	}

	return collection
}

// createHTTPieRequest creates an HTTPie request from a route
func createHTTPieRequest(route parser.Route, method string, env string) HTTPieRequest {
	// Build URL with path parameters
	baseURL := getBaseURL(env)
	url := baseURL + route.Path

	// Replace path parameters with placeholders
	for _, param := range route.Params {
		url = strings.ReplaceAll(url, ":"+param, "{{"+param+"}}")
		url = strings.ReplaceAll(url, "{"+param+"}", "{{"+param+"}}")
		url = strings.ReplaceAll(url, "<"+param+">", "{{"+param+"}}")
	}

	// Build path params array
	pathParams := []interface{}{}
	for _, param := range route.Params {
		pathParams = append(pathParams, map[string]interface{}{
			"name":    param,
			"value":   "",
			"enabled": true,
		})
	}

	// Build query params array
	queryParams := []interface{}{}
	for _, qp := range route.Query {
		queryParams = append(queryParams, map[string]interface{}{
			"name":    qp,
			"value":   "",
			"enabled": true,
		})
	}

	// Build body
	body := HTTPieRequestBody{
		Type: "none",
		File: HTTPieFile{
			Name: "",
		},
		Text: HTTPieText{
			Value:  "",
			Format: "application/json",
		},
		Form: HTTPieForm{
			IsMultipart: false,
			Fields:      []HTTPieFormField{},
		},
		GraphQL: HTTPieGraphQL{
			Query:     "",
			Variables: "",
		},
	}

	// Set body type and content for POST, PUT, PATCH
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if len(route.Body) > 0 {
			// Form data
			body.Type = "form"
			var formFields []HTTPieFormField
			for _, field := range route.Body {
				formFields = append(formFields, HTTPieFormField{
					Name:    field,
					Value:   "",
					Type:    "text",
					Enabled: true,
				})
			}
			body.Form.Fields = formFields
			body.Form.IsMultipart = false
		} else {
			// JSON body
			body.Type = "text"
			body.Text.Value = "{}"
			body.Text.Format = "application/json"
		}
	}

	request := HTTPieRequest{
		Name:        fmt.Sprintf("%s %s", method, route.Path),
		URL:         url,
		Method:      method,
		Headers:     []interface{}{},
		QueryParams: queryParams,
		PathParams:  pathParams,
		Auth: HTTPieAuth{
			Type: "inherited",
		},
		Body: body,
	}

	return request
}

