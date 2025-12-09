package exporters

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/parser"
)

// InsomniaExport represents an Insomnia export structure
type InsomniaExport struct {
	Type         string           `json:"_type"`
	ExportFormat int              `json:"__export_format"`
	ExportDate   string           `json:"__export_date"`
	ExportSource string           `json:"__export_source"`
	Resources    []InsomniaResource `json:"resources"`
}

// InsomniaResource represents a resource in Insomnia
type InsomniaResource struct {
	ID       string                 `json:"_id"`
	ParentID interface{}            `json:"parentId,omitempty"`
	Modified int64                  `json:"modified"`
	Created  int64                  `json:"created"`
	Name     string                 `json:"name"`
	Type     string                 `json:"_type"`
	Description string              `json:"description,omitempty"`
	Scope    string                 `json:"scope,omitempty"`
	URL      string                 `json:"url,omitempty"`
	Method   string                 `json:"method,omitempty"`
	Body     interface{}            `json:"body,omitempty"`
	Headers  []interface{}          `json:"headers,omitempty"`
	Parameters []interface{}        `json:"parameters,omitempty"`
	Authentication interface{}      `json:"authentication,omitempty"`
	MetaSortKey int64               `json:"metaSortKey,omitempty"`
	IsPrivate bool                  `json:"isPrivate,omitempty"`
	SettingStoreCookies bool        `json:"settingStoreCookies,omitempty"`
	SettingSendCookies bool         `json:"settingSendCookies,omitempty"`
	SettingDisableRenderRequestBody bool `json:"settingDisableRenderRequestBody,omitempty"`
	SettingEncodeUrl bool           `json:"settingEncodeUrl,omitempty"`
	SettingRebuildPath bool         `json:"settingRebuildPath,omitempty"`
	SettingFollowRedirects string   `json:"settingFollowRedirects,omitempty"`
	Environment interface{}         `json:"environment,omitempty"`
	EnvironmentPropertyOrder interface{} `json:"environmentPropertyOrder,omitempty"`
	Data     interface{}            `json:"data,omitempty"`
	DataPropertyOrder interface{}   `json:"dataPropertyOrder,omitempty"`
	Color    interface{}            `json:"color,omitempty"`
	Cookies  []interface{}          `json:"cookies,omitempty"`
	FileName string                 `json:"fileName,omitempty"`
	Contents string                 `json:"contents,omitempty"`
	ContentType string              `json:"contentType,omitempty"`
}

// InsomniaRequestBody represents request body
type InsomniaRequestBody struct {
	MimeType string                 `json:"mimeType"`
	Text     string                 `json:"text,omitempty"`
	Params   []InsomniaFormParam    `json:"params,omitempty"`
}

// InsomniaFormParam represents form parameter
type InsomniaFormParam struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// InsomniaHeader represents a header
type InsomniaHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// InsomniaParameter represents a URL parameter
type InsomniaParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// GenerateInsomniaExport generates an Insomnia export from routes
func GenerateInsomniaExport(routes []parser.Route, projectPath string, env string) InsomniaExport {
	envName := getEnvironmentName(env)
	baseURL := getBaseURL(env)
	
	// Generate timestamps
	baseTime := int64(1765300000000) // Base timestamp
	
	export := InsomniaExport{
		Type:         "export",
		ExportFormat: 4,
		ExportDate:   "2025-12-10T00:00:00.000Z",
		ExportSource: "likhis:api-mapper",
		Resources:    []InsomniaResource{},
	}

	collectionName := fmt.Sprintf("%s API (%s)", filepath.Base(projectPath), envName)
	
	// Create workspace
	workspaceID := "wrk_backend_api_001"
	workspace := InsomniaResource{
		ID:          workspaceID,
		Type:        "workspace",
		ParentID:    nil,
		Modified:    baseTime,
		Created:     baseTime,
		Name:        collectionName,
		Description: "",
		Scope:       "collection",
	}
	export.Resources = append(export.Resources, workspace)

	// Create request group
	groupID := "fld_backend_api_group"
	group := InsomniaResource{
		ID:                      groupID,
		Type:                    "request_group",
		ParentID:                workspaceID,
		Modified:                baseTime + 100,
		Created:                 baseTime + 100,
		Name:                    fmt.Sprintf("%s API", filepath.Base(projectPath)),
		Description:             "Auto-generated API collection",
		Environment:             map[string]interface{}{},
		EnvironmentPropertyOrder: nil,
		MetaSortKey:             -(baseTime + 100),
	}
	export.Resources = append(export.Resources, group)

	// Create requests
	for i, route := range routes {
		requestID := fmt.Sprintf("req_%s_%d", sanitizeID(route.Path), i+1)
		request := createInsomniaRequest(route, requestID, groupID, env, baseTime+int64(200+i*100))
		export.Resources = append(export.Resources, request)
	}

	// Create environment
	envID := "env_backend_base"
	environment := InsomniaResource{
		ID:                envID,
		Type:              "environment",
		ParentID:          workspaceID,
		Modified:          baseTime + int64(400+len(routes)*100),
		Created:           baseTime + int64(400+len(routes)*100),
		Name:              "Base Environment",
		Data:              map[string]interface{}{"base_url": baseURL},
		DataPropertyOrder: nil,
		Color:             nil,
		IsPrivate:         false,
		MetaSortKey:       baseTime + int64(400+len(routes)*100),
	}
	export.Resources = append(export.Resources, environment)

	// Create cookie jar
	jarID := "jar_backend"
	cookieJar := InsomniaResource{
		ID:       jarID,
		Type:     "cookie_jar",
		ParentID: workspaceID,
		Modified: baseTime + int64(500+len(routes)*100),
		Created:  baseTime + int64(500+len(routes)*100),
		Name:     "Default Jar",
		Cookies:  []interface{}{},
	}
	export.Resources = append(export.Resources, cookieJar)

	// Create API spec
	specID := "spc_backend_api_spec"
	apiSpec := InsomniaResource{
		ID:        specID,
		Type:      "api_spec",
		ParentID:  workspaceID,
		Modified:  baseTime + int64(600+len(routes)*100),
		Created:   baseTime + int64(600+len(routes)*100),
		FileName:  collectionName,
		Contents:  "",
		ContentType: "yaml",
	}
	export.Resources = append(export.Resources, apiSpec)

	return export
}

// sanitizeID creates a sanitized ID from a path
func sanitizeID(path string) string {
	// Remove special characters and convert to lowercase
	sanitized := strings.ToLower(path)
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, ":", "")
	sanitized = strings.ReplaceAll(sanitized, "{", "")
	sanitized = strings.ReplaceAll(sanitized, "}", "")
	sanitized = strings.ReplaceAll(sanitized, "<", "")
	sanitized = strings.ReplaceAll(sanitized, ">", "")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	return sanitized
}

// createInsomniaRequest creates an Insomnia request from a route
func createInsomniaRequest(route parser.Route, requestID, parentID string, env string, timestamp int64) InsomniaResource {
	// Use {{base_url}} placeholder for URL
	url := "{{base_url}}" + route.Path

	// Build body object
	body := map[string]interface{}{}
	if route.Method == "POST" || route.Method == "PUT" || route.Method == "PATCH" {
		if len(route.Body) > 0 {
			// Form data
			body["mimeType"] = "application/x-www-form-urlencoded"
			params := []map[string]interface{}{}
			for _, field := range route.Body {
				params = append(params, map[string]interface{}{
					"name":  field,
					"value": "",
				})
			}
			body["params"] = params
		} else {
			// JSON body
			body["mimeType"] = "application/json"
			body["text"] = "{}"
		}
	}

	request := InsomniaResource{
		ID:                          requestID,
		Type:                        "request",
		ParentID:                    parentID,
		Modified:                    timestamp,
		Created:                     timestamp,
		URL:                         url,
		Name:                        fmt.Sprintf("%s %s", route.Method, route.Path),
		Description:                 "",
		Method:                      route.Method,
		Body:                        body,
		Parameters:                  []interface{}{},
		Headers:                     []interface{}{},
		Authentication:              map[string]interface{}{},
		MetaSortKey:                 -timestamp,
		IsPrivate:                   false,
		SettingStoreCookies:         true,
		SettingSendCookies:          true,
		SettingDisableRenderRequestBody: false,
		SettingEncodeUrl:            true,
		SettingRebuildPath:          true,
		SettingFollowRedirects:      "global",
	}

	// Add query parameters
	for _, qp := range route.Query {
		request.Parameters = append(request.Parameters, map[string]interface{}{
			"name":  qp,
			"value": "",
		})
	}

	return request
}

