package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/marcuwynu23/likhis/internal/plugins"
)

func TestLoadPlugins(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_plugins_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create plugins directory
	pluginsDir := filepath.Join(tmpDir, "plugins")
	err = os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins dir: %v", err)
	}

	// Create a test plugin file
	pluginFile := filepath.Join(pluginsDir, "test.yml")
	pluginContent := `name: test
description: Test plugin
extensions:
  - .test
patterns:
  - method: GET
    route_regex: 'test\.(get|post)\([''"]([^''"]+)[''"]'
    param_regex: ':(\w+)'
`
	err = os.WriteFile(pluginFile, []byte(pluginContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create plugin file: %v", err)
	}

	// Load plugins
	pluginMap, loadedDirs, err := plugins.LoadPlugins("/fake/executable", tmpDir)
	if err != nil {
		t.Fatalf("LoadPlugins failed: %v", err)
	}

	// Verify plugin was loaded
	if len(pluginMap) == 0 {
		t.Error("No plugins loaded")
	}

	// Check if our test plugin is loaded
	testPlugin, ok := pluginMap["test"]
	if !ok {
		t.Error("Test plugin not found in loaded plugins")
	} else {
		if testPlugin.Name != "test" {
			t.Errorf("Expected plugin name 'test', got '%s'", testPlugin.Name)
		}
		if len(testPlugin.Extensions) == 0 {
			t.Error("Plugin has no extensions")
		}
		if len(testPlugin.Patterns) == 0 {
			t.Error("Plugin has no patterns")
		}
	}

	// Verify loaded directories
	if len(loadedDirs) == 0 {
		t.Error("No plugin directories reported as loaded")
	}
}

func TestGetPlugin(t *testing.T) {
	// Create a plugin map
	pluginMap := make(map[string]*plugins.Plugin)
	pluginMap["express"] = &plugins.Plugin{
		Name:        "express",
		Description: "Express.js plugin",
		Extensions:  []string{".js"},
		Patterns:    []plugins.Pattern{},
	}

	// Test getting existing plugin
	plugin, err := plugins.GetPlugin(pluginMap, "express")
	if err != nil {
		t.Fatalf("GetPlugin failed: %v", err)
	}
	if plugin == nil {
		t.Error("GetPlugin returned nil for existing plugin")
	}
	if plugin.Name != "express" {
		t.Errorf("Expected plugin name 'express', got '%s'", plugin.Name)
	}

	// Test getting non-existent plugin
	_, err = plugins.GetPlugin(pluginMap, "nonexistent")
	if err == nil {
		t.Error("Expected error for nonexistent plugin")
	}
}

func TestGetPluginEmptyName(t *testing.T) {
	pluginMap := make(map[string]*plugins.Plugin)
	_, err := plugins.GetPlugin(pluginMap, "")
	if err == nil {
		t.Error("Expected error for empty plugin name")
	}
}

func TestLoadPluginsInvalidYAML(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_plugins_invalid_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pluginsDir := filepath.Join(tmpDir, "plugins")
	err = os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins dir: %v", err)
	}

	// Create invalid YAML file
	invalidFile := filepath.Join(pluginsDir, "invalid.yml")
	err = os.WriteFile(invalidFile, []byte("invalid: yaml: content: [unclosed"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid plugin file: %v", err)
	}

	// Load plugins - should handle invalid YAML gracefully
	pluginMap, _, err := plugins.LoadPlugins("/fake/executable", tmpDir)
	// Should not fail completely, just skip invalid file
	if err != nil {
		// Error is acceptable, but should not crash
		t.Logf("LoadPlugins returned error (acceptable): %v", err)
	}

	// Invalid plugin should not be in the map
	if _, ok := pluginMap["invalid"]; ok {
		t.Error("Invalid plugin should not be loaded")
	}
}

func TestLoadPluginsMissingExtensions(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "likhis_plugins_no_ext_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	pluginsDir := filepath.Join(tmpDir, "plugins")
	err = os.MkdirAll(pluginsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create plugins dir: %v", err)
	}

	// Create plugin without extensions
	pluginFile := filepath.Join(pluginsDir, "noext.yml")
	pluginContent := `name: noext
description: Plugin without extensions
patterns:
  - method: GET
    route_regex: 'test'
`
	err = os.WriteFile(pluginFile, []byte(pluginContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create plugin file: %v", err)
	}

	// Load plugins - should skip plugin without extensions
	pluginMap, _, _ := plugins.LoadPlugins("/fake/executable", tmpDir)
	
	// Plugin without extensions should not be loaded
	if _, ok := pluginMap["noext"]; ok {
		t.Error("Plugin without extensions should not be loaded")
	}
}

