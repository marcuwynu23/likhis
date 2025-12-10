package plugins

import (
	"fmt"
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

// LoadPlugins loads all plugin YAML files from multiple possible plugin directories
// It checks in order: project directory, executable directory, current directory
func LoadPlugins(executablePath string, projectPath string) (map[string]*Plugin, []string, error) {
	plugins := make(map[string]*Plugin)
	var loadedDirs []string
	
	// List of plugin directories to check (in priority order)
	var pluginDirs []string
	
	// 1. Project directory plugins (highest priority for custom plugins)
	if projectPath != "" {
		absProjectPath, err := filepath.Abs(projectPath)
		if err == nil {
			projectPluginsDir := filepath.Join(absProjectPath, "plugins")
			if _, err := os.Stat(projectPluginsDir); err == nil {
				pluginDirs = append(pluginDirs, projectPluginsDir)
			}
		}
	}
	
	// 2. Executable directory plugins (and parent directories)
	// First, try to resolve symlinks to get the actual executable path
	actualExecPath := executablePath
	if resolved, err := filepath.EvalSymlinks(executablePath); err == nil {
		actualExecPath = resolved
	}
	
	// Check executable directory and walk up to find plugins directory
	// This handles cases where executable is in a subdirectory (e.g., build/)
	execDir := filepath.Dir(actualExecPath)
	currentCheckDir := execDir
	maxDepth := 5 // Increased depth to handle deeper directory structures
	depth := 0
	
	for depth < maxDepth {
		execPluginsDir := filepath.Join(currentCheckDir, "plugins")
		if _, err := os.Stat(execPluginsDir); err == nil {
			// Only add if not already in list
			alreadyAdded := false
			for _, dir := range pluginDirs {
				// Use absolute paths for comparison
				absDir, _ := filepath.Abs(dir)
				absExecPluginsDir, _ := filepath.Abs(execPluginsDir)
				if absDir == absExecPluginsDir {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				// Use absolute path
				if absPath, err := filepath.Abs(execPluginsDir); err == nil {
					pluginDirs = append(pluginDirs, absPath)
				} else {
					pluginDirs = append(pluginDirs, execPluginsDir)
				}
			}
			break // Found plugins directory, stop searching
		}
		
		// Move up one directory
		parentDir := filepath.Dir(currentCheckDir)
		if parentDir == currentCheckDir {
			// Reached root, stop
			break
		}
		currentCheckDir = parentDir
		depth++
	}
	
	// 3. Current directory plugins (fallback)
	currentDir, err := os.Getwd()
	if err == nil {
		currentPluginsDir := filepath.Join(currentDir, "plugins")
		if _, err := os.Stat(currentPluginsDir); err == nil {
			// Only add if not already in list (avoid duplicates)
			alreadyAdded := false
			for _, dir := range pluginDirs {
				if dir == currentPluginsDir {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				pluginDirs = append(pluginDirs, currentPluginsDir)
			}
		}
	}
	
	// Load plugins from each directory (later directories can override earlier ones)
	for _, pluginsDir := range pluginDirs {
		loaded, err := loadPluginsFromDir(pluginsDir)
		if err != nil {
			// Log error but continue with other directories
			continue
		}
		
		// Merge plugins (later directories override earlier ones)
		for key, plugin := range loaded {
			plugins[key] = plugin
		}
		
		if len(loaded) > 0 {
			loadedDirs = append(loadedDirs, pluginsDir)
		}
	}
	
	return plugins, loadedDirs, nil
}

// loadPluginsFromDir loads all plugin YAML files from a specific directory
func loadPluginsFromDir(pluginsDir string) (map[string]*Plugin, error) {
	plugins := make(map[string]*Plugin)
	
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
			// Log error but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: Could not read plugin file %s: %v\n", pluginPath, err)
			continue
		}
		
		var plugin Plugin
		if err := yaml.Unmarshal(data, &plugin); err != nil {
			// Log error but continue with other files
			fmt.Fprintf(os.Stderr, "Warning: Could not parse plugin file %s: %v\n", pluginPath, err)
			continue
		}
		
		// Validate plugin has required fields
		if len(plugin.Extensions) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: Plugin file %s has no extensions defined, skipping\n", pluginPath)
			continue
		}
		if len(plugin.Patterns) == 0 {
			fmt.Fprintf(os.Stderr, "Warning: Plugin file %s has no patterns defined, skipping\n", pluginPath)
			continue
		}
		
		// Use filename (without extension) as key - this allows express-v2.yaml to be keyed as "express-v2"
		// The filename takes priority over the plugin name field to support custom plugin variants
		key := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		
		// Only use plugin.Name as key if it's explicitly different and filename is generic
		// This maintains backward compatibility while allowing custom plugin names
		if plugin.Name != "" && key == "plugin" {
			key = strings.ToLower(plugin.Name)
		}
		
		plugins[key] = &plugin
	}
	
	return plugins, nil
}

// GetPlugin returns a plugin by name with exact matching only
// Plugin names must match exactly (case-insensitive) to the filename without extension
func GetPlugin(plugins map[string]*Plugin, name string) (*Plugin, error) {
	if name == "" {
		return nil, fmt.Errorf("plugin name cannot be empty")
	}
	
	normalizedName := strings.ToLower(name)
	
	// Only try exact match - plugin names must match exactly
	if plugin, ok := plugins[normalizedName]; ok {
		return plugin, nil
	}
	
	// Build list of available plugins for error message
	availablePlugins := make([]string, 0, len(plugins))
	for key := range plugins {
		availablePlugins = append(availablePlugins, key)
	}
	
	return nil, fmt.Errorf("plugin '%s' not found. Available plugins: %s", name, strings.Join(availablePlugins, ", "))
}


