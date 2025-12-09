package plugins

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Plugin represents a framework plugin configuration
type Plugin struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Extensions  []string `yaml:"extensions"`
	Patterns    []Pattern `yaml:"patterns"`
	RouterMount RouterMount `yaml:"router_mount,omitempty"`
}

// Pattern represents a route pattern to match
type Pattern struct {
	Method      string   `yaml:"method"`
	RouteRegex  string   `yaml:"route_regex"`
	ParamRegex  string   `yaml:"param_regex,omitempty"`
	QueryRegex  string   `yaml:"query_regex,omitempty"`
	BodyRegex   string   `yaml:"body_regex,omitempty"`
}

// RouterMount represents router mounting patterns (for Express, etc.)
type RouterMount struct {
	UsePattern    string `yaml:"use_pattern,omitempty"`
	RequirePattern string `yaml:"require_pattern,omitempty"`
	VarPattern    string `yaml:"var_pattern,omitempty"`
}

// LoadPlugins loads all plugin YAML files from the plugins directory
func LoadPlugins(executablePath string) (map[string]*Plugin, error) {
	plugins := make(map[string]*Plugin)
	
	// Get directory where executable is located
	execDir := filepath.Dir(executablePath)
	pluginsDir := filepath.Join(execDir, "plugins")
	
	// Check if plugins directory exists
	if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
		// Try current directory as fallback
		pluginsDir = "plugins"
		if _, err := os.Stat(pluginsDir); os.IsNotExist(err) {
			return plugins, nil // No plugins directory, return empty map
		}
	}
	
	// Read all YAML files in plugins directory
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return plugins, err
	}
	
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".yml") && 
		   !strings.HasSuffix(strings.ToLower(entry.Name()), ".yaml") {
			continue
		}
		
		pluginPath := filepath.Join(pluginsDir, entry.Name())
		data, err := os.ReadFile(pluginPath)
		if err != nil {
			continue
		}
		
		var plugin Plugin
		if err := yaml.Unmarshal(data, &plugin); err != nil {
			continue
		}
		
		// Use filename (without extension) as key, or plugin name
		key := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		if plugin.Name != "" {
			key = strings.ToLower(plugin.Name)
		}
		plugins[key] = &plugin
	}
	
	return plugins, nil
}

// GetPlugin returns a plugin by name
func GetPlugin(plugins map[string]*Plugin, name string) *Plugin {
	// Try exact match
	if plugin, ok := plugins[strings.ToLower(name)]; ok {
		return plugin
	}
	
	// Try partial match
	for key, plugin := range plugins {
		if strings.Contains(strings.ToLower(key), strings.ToLower(name)) {
			return plugin
		}
	}
	
	return nil
}

