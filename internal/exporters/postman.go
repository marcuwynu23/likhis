package exporters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/parser"
)

// PostmanCollection represents a Postman Collection v2.1 structure
type PostmanCollection struct {
	Info struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Schema      string `json:"schema"`
	} `json:"info"`
	Item []PostmanItem `json:"item"`
}

// PostmanItem represents a request item in Postman
type PostmanItem struct {
	Name     string         `json:"name"`
	Request  PostmanRequest `json:"request"`
	Response []interface{}  `json:"response"`
}

// PostmanRequest represents a Postman request
type PostmanRequest struct {
	Method string                 `json:"method"`
	Header []interface{}          `json:"header"`
	Body   *PostmanRequestBody    `json:"body,omitempty"`
	URL    PostmanURL             `json:"url"`
}

// PostmanRequestBody represents request body
type PostmanRequestBody struct {
	Mode     string            `json:"mode"`
	Raw      string            `json:"raw,omitempty"`
	Formdata []PostmanFormField `json:"formdata,omitempty"`
}

// PostmanFormField represents form data field
type PostmanFormField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// PostmanURL represents a Postman URL
type PostmanURL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
	Query    []PostmanQueryParam `json:"query,omitempty"`
	Variable []PostmanVariable   `json:"variable,omitempty"`
}

// PostmanQueryParam represents a query parameter
type PostmanQueryParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Description string `json:"description,omitempty"`
}

// PostmanVariable represents a path variable
type PostmanVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GeneratePostmanCollection generates a Postman Collection v2.1 from routes
func GeneratePostmanCollection(routes []parser.Route, projectPath string, env string) PostmanCollection {
	envName := getEnvironmentName(env)
	collection := PostmanCollection{}
	collection.Info.Name = fmt.Sprintf("%s API (%s)", filepath.Base(projectPath), envName)
	collection.Info.Description = fmt.Sprintf("Auto-generated API collection from %s - %s environment", projectPath, envName)
	collection.Info.Schema = "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"

	// Group routes by method for better organization
	routesByMethod := make(map[string][]parser.Route)
	for _, route := range routes {
		routesByMethod[route.Method] = append(routesByMethod[route.Method], route)
	}

	// Create items grouped by HTTP method
	for method, methodRoutes := range routesByMethod {
		for _, route := range methodRoutes {
			item := createPostmanItem(route, method, env)
			collection.Item = append(collection.Item, item)
		}
	}

	return collection
}

// createPostmanItem creates a Postman item from a route
func createPostmanItem(route parser.Route, method string, env string) PostmanItem {
	item := PostmanItem{
		Name:     fmt.Sprintf("%s %s", method, route.Path),
		Request:  PostmanRequest{Method: method},
		Response: []interface{}{},
	}

	// Parse path and extract variables
	pathParts := strings.Split(strings.Trim(route.Path, "/"), "/")
	var urlPath []string
	var variables []PostmanVariable

	for _, part := range pathParts {
		if part == "" {
			continue
		}
		// Check if it's a parameter (Express :id, Flask <id>, Spring {id}, Laravel {id})
		if strings.HasPrefix(part, ":") {
			paramName := part[1:]
			urlPath = append(urlPath, fmt.Sprintf(":%s", paramName))
			variables = append(variables, PostmanVariable{
				Key:   paramName,
				Value: "",
			})
		} else if strings.HasPrefix(part, "<") && strings.HasSuffix(part, ">") {
			// Flask/Django style: <id> or <int:id>
			paramName := strings.Trim(part, "<>")
			if idx := strings.Index(paramName, ":"); idx != -1 {
				paramName = paramName[idx+1:]
			}
			urlPath = append(urlPath, fmt.Sprintf(":%s", paramName))
			variables = append(variables, PostmanVariable{
				Key:   paramName,
				Value: "",
			})
		} else if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			// Spring/Laravel style: {id}
			paramName := strings.Trim(part, "{}")
			urlPath = append(urlPath, fmt.Sprintf(":%s", paramName))
			variables = append(variables, PostmanVariable{
				Key:   paramName,
				Value: "",
			})
		} else {
			urlPath = append(urlPath, part)
		}
	}

	// Build URL with environment-specific base URL
	baseURL := getBaseURL(env)
	// For Postman, use variable if it's a placeholder, otherwise use the actual URL
	if strings.HasPrefix(baseURL, "{{") {
		item.Request.URL = PostmanURL{
			Raw:      baseURL + route.Path,
			Protocol: "https",
			Host:     []string{baseURL},
			Path:     urlPath,
			Variable: variables,
		}
	} else {
		// Parse the base URL to extract protocol and host
		protocol := "https"
		host := baseURL
		if strings.HasPrefix(baseURL, "http://") {
			protocol = "http"
			host = strings.TrimPrefix(baseURL, "http://")
		} else if strings.HasPrefix(baseURL, "https://") {
			host = strings.TrimPrefix(baseURL, "https://")
		}
		item.Request.URL = PostmanURL{
			Raw:      baseURL + route.Path,
			Protocol: protocol,
			Host:     []string{host},
			Path:     urlPath,
			Variable: variables,
		}
	}

	// Add query parameters
	if len(route.Query) > 0 {
		for _, qp := range route.Query {
			item.Request.URL.Query = append(item.Request.URL.Query, PostmanQueryParam{
				Key:   qp,
				Value: "",
			})
		}
	}

	// Add body for POST, PUT, PATCH
	if method == "POST" || method == "PUT" || method == "PATCH" {
		if len(route.Body) > 0 {
			// Create form data
			var formData []PostmanFormField
			for _, field := range route.Body {
				formData = append(formData, PostmanFormField{
					Key:   field,
					Value: "",
					Type:  "text",
				})
			}
			item.Request.Body = &PostmanRequestBody{
				Mode:     "formdata",
				Formdata: formData,
			}
		} else {
			// Default JSON body
			item.Request.Body = &PostmanRequestBody{
				Mode: "raw",
				Raw:  "{}",
			}
		}
	}

	// Add headers
	item.Request.Header = []interface{}{}

	return item
}

