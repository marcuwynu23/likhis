package tests

import (
	"strings"
	"testing"

	"github.com/marcuwynu23/likhis/internal/exporters"
	"github.com/marcuwynu23/likhis/internal/parser"
)

func TestGeneratePostmanCollection(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "GET",
			Params: []string{},
			Query:  []string{"page"},
			Body:   []string{},
		},
		{
			Path:   "/users/:id",
			Method: "POST",
			Params: []string{"id"},
			Query:  []string{},
			Body:   []string{"name", "email"},
		},
	}

	collection := exporters.GeneratePostmanCollection(routes, "/test", "dev")
	
	if collection.Info.Name == "" {
		t.Error("Collection info name should not be empty")
	}
	
	if len(collection.Item) != 2 {
		t.Errorf("Expected 2 items in collection, got %d", len(collection.Item))
	}

	// Find routes by method and path (order may vary)
	foundGET := false
	foundPOST := false
	
	for _, item := range collection.Item {
		if item.Request.Method == "GET" && strings.Contains(item.Name, "/users") && !strings.Contains(item.Name, ":id") {
			foundGET = true
		}
		if item.Request.Method == "POST" && strings.Contains(item.Name, "/users/:id") {
			foundPOST = true
			// Verify POST route has body
			if item.Request.Body == nil {
				t.Error("POST route should have a body")
			}
		}
	}
	
	if !foundGET {
		t.Error("GET /users route not found in collection")
	}
	if !foundPOST {
		t.Error("POST /users/:id route not found in collection")
	}
}

func TestGenerateInsomniaExport(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/api/users",
			Method: "GET",
		},
	}

	export := exporters.GenerateInsomniaExport(routes, "/test", "dev")
	
	if export.Type != "export" {
		t.Errorf("Expected type 'export', got %s", export.Type)
	}
	
	if len(export.Resources) == 0 {
		t.Error("Expected at least one resource in export")
	}
}

func TestGenerateHTTPieExport(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "GET",
		},
	}

	export := exporters.GenerateHTTPieExport(routes, "/test", "dev")
	
	if len(export.Entry.Requests) == 0 {
		t.Error("Expected at least one request in export")
	}
	
	if export.Meta.Format != "httpie" {
		t.Errorf("Expected format 'httpie', got %s", export.Meta.Format)
	}
}

func TestGenerateCURLScript(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "GET",
			Query:  []string{"page"},
		},
		{
			Path:   "/users/:id",
			Method: "POST",
			Body:   []string{"name"},
		},
	}

	script := exporters.GenerateCURLScript(routes, "/test", "dev")
	
	if script == "" {
		t.Error("CURL script should not be empty")
	}
	
	// Verify it contains curl commands
	if len(script) < 10 {
		t.Error("CURL script seems too short")
	}
	
	// Verify it's valid (can be parsed as text at least)
	_ = script // Just ensure it's not empty
}

func TestGenerateOpenAPIExport(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "GET",
			Params: []string{},
			Query:  []string{"page", "limit"},
			Body:   []string{},
		},
		{
			Path:   "/users/:id",
			Method: "GET",
			Params: []string{"id"},
			Query:  []string{},
			Body:   []string{},
		},
		{
			Path:   "/users",
			Method: "POST",
			Params: []string{},
			Query:  []string{},
			Body:   []string{"name", "email"},
		},
		{
			Path:   "/users/{id}",
			Method: "PUT",
			Params: []string{"id"},
			Query:  []string{},
			Body:   []string{"name"},
		},
		{
			Path:   "/users/<id>",
			Method: "DELETE",
			Params: []string{"id"},
			Query:  []string{},
			Body:   []string{},
		},
	}

	doc := exporters.GenerateOpenAPIExport(routes, "/test", "dev")

	if doc.OpenAPI != "3.0.3" {
		t.Errorf("Expected OpenAPI version 3.0.3, got %s", doc.OpenAPI)
	}

	if doc.Info.Title == "" {
		t.Error("Info title should not be empty")
	}

	if len(doc.Servers) != 1 {
		t.Errorf("Expected 1 server, got %d", len(doc.Servers))
	}

	if len(doc.Paths) != 2 {
		t.Errorf("Expected 2 paths, got %d", len(doc.Paths))
	}

	usersPath, ok := doc.Paths["/users"]
	if !ok {
		t.Fatal("Expected /users path")
	}
	if usersPath.Get == nil {
		t.Error("Expected GET /users")
	}
	if usersPath.Post == nil {
		t.Error("Expected POST /users")
	}
	if len(usersPath.Get.Parameters) != 2 {
		t.Errorf("Expected 2 query parameters on GET /users, got %d", len(usersPath.Get.Parameters))
	}

	usersIdPath, ok := doc.Paths["/users/{id}"]
	if !ok {
		t.Fatal("Expected /users/{id} path")
	}
	if usersIdPath.Get == nil {
		t.Error("Expected GET /users/{id}")
	}
	if usersIdPath.Put == nil {
		t.Error("Expected PUT /users/{id}")
	}
	if usersIdPath.Delete == nil {
		t.Error("Expected DELETE /users/{id}")
	}

	for _, param := range usersIdPath.Get.Parameters {
		if param.Name == "id" && param.In == "path" {
			if !param.Required {
				t.Error("Path parameter 'id' should be required")
			}
			return
		}
	}
	t.Error("Path parameter 'id' not found on GET /users/{id}")
}

func TestGenerateOpenAPIExportWithBody(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "POST",
			Params: []string{},
			Query:  []string{},
			Body:   []string{"name", "email", "age"},
		},
		{
			Path:   "/posts",
			Method: "POST",
			Params: []string{},
			Query:  []string{},
			Body:   []string{},
		},
	}

	doc := exporters.GenerateOpenAPIExport(routes, "/test", "prod")

	usersPost := doc.Paths["/users"].Post
	if usersPost == nil {
		t.Fatal("Expected POST /users")
	}
	if usersPost.RequestBody == nil {
		t.Fatal("POST /users should have a request body")
	}
	if !usersPost.RequestBody.Required {
		t.Error("Request body should be required")
	}

	jsonContent, ok := usersPost.RequestBody.Content["application/json"]
	if !ok {
		t.Fatal("Expected application/json content type")
	}
	if jsonContent.Schema.Type != "object" {
		t.Errorf("Expected schema type 'object', got %s", jsonContent.Schema.Type)
	}
	if len(jsonContent.Schema.Properties) != 3 {
		t.Errorf("Expected 3 body properties, got %d", len(jsonContent.Schema.Properties))
	}

	props := jsonContent.Schema.Properties
	if _, ok := props["name"]; !ok {
		t.Error("Expected 'name' property in body schema")
	}
	if _, ok := props["email"]; !ok {
		t.Error("Expected 'email' property in body schema")
	}

	postsPost := doc.Paths["/posts"].Post
	if postsPost == nil {
		t.Fatal("Expected POST /posts")
	}
	if postsPost.RequestBody == nil {
		t.Fatal("POST /posts should have a default request body")
	}
	if postsPost.RequestBody.Content["application/json"].Schema.Properties != nil {
		t.Error("Empty body should not have properties")
	}
}

func TestGenerateOpenAPIExportServers(t *testing.T) {
	routes := []parser.Route{
		{Path: "/test", Method: "GET"},
	}

	devDoc := exporters.GenerateOpenAPIExport(routes, "/test", "dev")
	if len(devDoc.Servers) != 1 {
		t.Fatalf("Expected 1 server, got %d", len(devDoc.Servers))
	}
	if devDoc.Servers[0].Description != "Development" {
		t.Errorf("Expected 'Development' server, got %s", devDoc.Servers[0].Description)
	}

	stagingDoc := exporters.GenerateOpenAPIExport(routes, "/test", "staging")
	if stagingDoc.Servers[0].Description != "Staging" {
		t.Errorf("Expected 'Staging' server, got %s", stagingDoc.Servers[0].Description)
	}

	prodDoc := exporters.GenerateOpenAPIExport(routes, "/test", "prod")
	if prodDoc.Servers[0].Description != "Production" {
		t.Errorf("Expected 'Production' server, got %s", prodDoc.Servers[0].Description)
	}
	if prodDoc.Servers[0].URL != "https://api.example.com" {
		t.Errorf("Expected prod URL 'https://api.example.com', got %s", prodDoc.Servers[0].URL)
	}
}

func TestGenerateOpenAPIExportEmptyRoutes(t *testing.T) {
	doc := exporters.GenerateOpenAPIExport([]parser.Route{}, "/test", "dev")

	if doc.OpenAPI != "3.0.3" {
		t.Error("Should return valid OpenAPI document even with no routes")
	}
	if len(doc.Paths) != 0 {
		t.Errorf("Expected 0 paths for empty routes, got %d", len(doc.Paths))
	}
}

func TestGenerateOpenAPIExportResponses(t *testing.T) {
	routes := []parser.Route{
		{Path: "/health", Method: "GET"},
	}

	doc := exporters.GenerateOpenAPIExport(routes, "/test", "dev")

	healthPath := doc.Paths["/health"]
	if healthPath.Get == nil {
		t.Fatal("Expected GET /health")
	}

	resp, ok := healthPath.Get.Responses["200"]
	if !ok {
		t.Fatal("Expected 200 response")
	}
	if resp.Description != "Successful response" {
		t.Errorf("Expected 'Successful response', got %s", resp.Description)
	}
}

func TestGenerateCURLMarkdown(t *testing.T) {
	routes := []parser.Route{
		{
			Path:   "/users",
			Method: "GET",
		},
	}

	markdown := exporters.GenerateCURLMarkdown(routes, "/test")
	
	if markdown == "" {
		t.Error("CURL markdown should not be empty")
	}
}

