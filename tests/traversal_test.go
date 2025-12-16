package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/likhis/internal/traversal"
)

func TestBFSTraverse(t *testing.T) {
	// Create a temporary directory structure for testing
	tmpDir, err := os.MkdirTemp("", "likhis_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directory structure
	testDirs := []string{
		"src",
		"src/api",
		"src/api/routes",
		"node_modules", // Should be skipped
		".git",         // Should be skipped
	}
	for _, dir := range testDirs {
		err := os.MkdirAll(filepath.Join(tmpDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create test dir: %v", err)
		}
	}

	// Create test files
	testFiles := []struct {
		path     string
		content  string
		shouldInclude bool
	}{
		{"src/api/routes/users.js", "// users route", true},
		{"src/api/routes/products.js", "// products route", true},
		{"src/main.py", "# main file", true},
		{"src/controller/UserController.java", "// controller", true},
		{"node_modules/express/index.js", "// should be skipped", false},
		{"README.md", "# readme", false}, // .md not in extensions
		{"package.json", "{}", false},    // .json not in extensions
	}

	for _, tf := range testFiles {
		fullPath := filepath.Join(tmpDir, tf.path)
		dir := filepath.Dir(fullPath)
		os.MkdirAll(dir, 0755)
		err := os.WriteFile(fullPath, []byte(tf.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Run traversal
	files, err := traversal.BFSTraverse(tmpDir)
	if err != nil {
		t.Fatalf("BFSTraverse failed: %v", err)
	}

	// Verify results
	expectedCount := 0
	for _, tf := range testFiles {
		if tf.shouldInclude {
			expectedCount++
		}
	}

	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Verify specific files are included
	foundUsers := false
	foundProducts := false
	for _, file := range files {
		if filepath.Base(file) == "users.js" {
			foundUsers = true
		}
		if filepath.Base(file) == "products.js" {
			foundProducts = true
		}
		// Verify node_modules files are not included
		if filepath.Base(filepath.Dir(file)) == "node_modules" {
			t.Errorf("Found file in node_modules: %s", file)
		}
	}

	if !foundUsers {
		t.Error("users.js not found in traversal results")
	}
	if !foundProducts {
		t.Error("products.js not found in traversal results")
	}
}

func TestBFSTraverseEmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_test_empty_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	files, err := traversal.BFSTraverse(tmpDir)
	if err != nil {
		t.Fatalf("BFSTraverse failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files, got %d", len(files))
	}
}

func TestBFSTraverseNonexistentDir(t *testing.T) {
	// This should handle gracefully
	files, err := traversal.BFSTraverse("/nonexistent/directory/path")
	if err != nil {
		// Error is acceptable for nonexistent directory
		return
	}
	
	// If no error, should return empty slice
	if len(files) != 0 {
		t.Errorf("Expected 0 files for nonexistent dir, got %d", len(files))
	}
}

