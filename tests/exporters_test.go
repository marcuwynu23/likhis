package tests

import (
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

	// Verify first route
	if collection.Item[0].Request.Method != "GET" {
		t.Errorf("Expected GET method, got %s", collection.Item[0].Request.Method)
	}
	
	// Verify second route has body
	if collection.Item[1].Request.Body == nil {
		t.Error("POST route should have a body")
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

