package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marcuwynu23/likhis/internal/exporters"
	"github.com/marcuwynu23/likhis/internal/parser"
	"github.com/marcuwynu23/likhis/internal/plugins"
	"github.com/marcuwynu23/likhis/internal/traversal"
)

func main() {
	// Get executable path for plugin loading
	executablePath, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		executablePath, _ = os.Getwd()
		executablePath = filepath.Join(executablePath, "rrs.exe")
	}
	// Always resolve to absolute path to ensure consistent plugin directory resolution
	if absExecPath, err := filepath.Abs(executablePath); err == nil {
		executablePath = absExecPath
	}

	// Initial plugin load (from executable and current directory)
	// We'll reload after parsing flags to include project directory plugins
	pluginMap, loadedDirs, err := plugins.LoadPlugins(executablePath, ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load plugins: %v\n", err)
	}
	if loadedDirs == nil {
		loadedDirs = []string{}
	}

	// Build initial framework list from plugins
	frameworkList := []string{"auto"}
	for name := range pluginMap {
		frameworkList = append(frameworkList, name)
	}
	frameworkHelp := fmt.Sprintf("Framework to target: %s (default: 'auto')", strings.Join(frameworkList, ", "))

	var (
		projectPath = flag.String("path", ".", "Path to the project root directory")
		output      = flag.String("output", "postman", "Output format: 'postman', 'insomnia', 'httpie', or 'curl'")
		outputFile  = flag.String("file", "", "Output file path (default: based on output format)")
		outputPath  = flag.String("output-path", "", "Output directory path (default: current directory)")
		framework   = flag.String("framework", "auto", frameworkHelp)
		full        = flag.Bool("full", false, "Generate exports for dev, staging, and production environments")
	)

	// Add short flags
	flag.StringVar(projectPath, "p", ".", "Path to the project root directory (short for --path)")
	flag.StringVar(output, "o", "postman", "Output format (short for --output)")
	flag.StringVar(outputFile, "f", "", "Output file path (short for --file)")
	flag.StringVar(outputPath, "O", "", "Output directory path (short for --output-path)")
	flag.StringVar(framework, "F", "auto", "Framework to target (short for --framework)")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Universal API Mapper Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s --path ./my-app --output postman\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -p ./my-app -o insomnia -F express\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --path ./api --output curl --file requests.sh\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --path ./api --output postman --output-path ./exports\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --path ./api --output postman --full\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    (generates dev, staging, and production exports)\n")
		fmt.Fprintf(os.Stderr, "\nSupported Output Formats: postman, insomnia, httpie, curl\n")
		if len(pluginMap) > 0 {
			fmt.Fprintf(os.Stderr, "\nAvailable Framework Plugins:\n")
			for name, plugin := range pluginMap {
				fmt.Fprintf(os.Stderr, "  - %s: %s\n", name, plugin.Description)
			}
		}
	}

	flag.Parse()

	// Reload plugins with the actual project path (to pick up project-level plugins)
	// This ensures we get plugins from: project directory, executable directory, and current directory
	allPlugins, allLoadedDirs, err := plugins.LoadPlugins(executablePath, *projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not load plugins: %v\n", err)
	}
	// Replace plugin map with fully loaded plugins
	pluginMap = allPlugins
	if allLoadedDirs != nil {
		loadedDirs = allLoadedDirs
	} else {
		loadedDirs = []string{}
	}

	// Show help if no arguments provided (except flags)
	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	if *projectPath == "" {
		fmt.Fprintf(os.Stderr, "Error: project path is required\n")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(*projectPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Scanning project at: %s\n", absPath)
	fmt.Printf("Framework detection: %s\n", *framework)
	fmt.Printf("Output format: %s\n", *output)
	if len(pluginMap) > 0 {
		fmt.Printf("Loaded %d plugin(s)", len(pluginMap))
		if len(loadedDirs) > 0 {
			fmt.Printf(" from: %s", strings.Join(loadedDirs, ", "))
		}
		fmt.Println()
	}
	
	// Validate framework plugin if specified (not "auto")
	if *framework != "auto" {
		_, err := plugins.GetPlugin(pluginMap, *framework)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	// BFS file traversal
	files, err := traversal.BFSTraverse(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error traversing files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d files to analyze\n", len(files))

	// Parse routes from all files using plugins
	routeParser := parser.NewRouteParserWithPlugins(*framework, pluginMap, executablePath)

	// Build router map if plugin supports it
	if *framework == "auto" || routeParser.HasRouterMountSupport(*framework) {
		routeParser.BuildRouterMap(files, absPath)
	}

	var allRoutes []parser.Route

	for _, file := range files {
		routes, err := routeParser.ParseFile(file)
		if err != nil {
			// Skip files that can't be parsed (not an error for our use case)
			continue
		}
		allRoutes = append(allRoutes, routes...)
	}

	fmt.Printf("Extracted %d routes\n", len(allRoutes))

	if len(allRoutes) == 0 {
		fmt.Println("No routes found. Make sure you're pointing to a valid backend project.")
		os.Exit(0)
	}

	// Determine environments to generate
	environments := []string{"dev"}
	if *full {
		environments = []string{"dev", "staging", "prod"}
		fmt.Printf("Generating exports for: %s\n", strings.Join(environments, ", "))
	}

	// Generate output for each environment
	for _, env := range environments {
		var outputData []byte
		var defaultFileName string
		var baseFileName string

		// Determine base file name
		if *outputFile != "" {
			baseFileName = *outputFile
			// Remove extension if present
			ext := filepath.Ext(baseFileName)
			baseFileName = strings.TrimSuffix(baseFileName, ext)
		} else {
			switch strings.ToLower(*output) {
			case "postman":
				baseFileName = "postman-collection"
			case "insomnia":
				baseFileName = "insomnia-export"
			case "httpie":
				baseFileName = "httpie-collection"
			case "curl":
				baseFileName = "api-requests"
			}
		}

		// Generate output based on format
		switch strings.ToLower(*output) {
		case "postman":
			collection := exporters.GeneratePostmanCollection(allRoutes, absPath, env)
			outputData, err = json.MarshalIndent(collection, "", "  ")
			if *full {
				defaultFileName = fmt.Sprintf("%s-%s.json", baseFileName, env)
			} else {
				defaultFileName = baseFileName + ".json"
			}
		case "insomnia":
			insomniaData := exporters.GenerateInsomniaExport(allRoutes, absPath, env)
			outputData, err = json.MarshalIndent(insomniaData, "", "  ")
			if *full {
				defaultFileName = fmt.Sprintf("%s-%s.json", baseFileName, env)
			} else {
				defaultFileName = baseFileName + ".json"
			}
		case "httpie":
			httpieData := exporters.GenerateHTTPieExport(allRoutes, absPath, env)
			outputData, err = json.MarshalIndent(httpieData, "", "  ")
			if *full {
				defaultFileName = fmt.Sprintf("%s-%s.json", baseFileName, env)
			} else {
				defaultFileName = baseFileName + ".json"
			}
		case "curl":
			curlScript := exporters.GenerateCURLScript(allRoutes, absPath, env)
			outputData = []byte(curlScript)
			if *full {
				defaultFileName = fmt.Sprintf("%s-%s.sh", baseFileName, env)
			} else {
				defaultFileName = baseFileName + ".sh"
			}
		default:
			fmt.Fprintf(os.Stderr, "Error: unsupported output format '%s'. Use 'postman', 'insomnia', 'httpie', or 'curl'\n", *output)
			os.Exit(1)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating output: %v\n", err)
			os.Exit(1)
		}

		// Determine output directory
		outputDir := "."
		if *outputPath != "" {
			outputDir = *outputPath
			// Create directory if it doesn't exist
			err = os.MkdirAll(outputDir, 0755)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
				os.Exit(1)
			}
		}

		// Construct full file path
		fullPath := filepath.Join(outputDir, defaultFileName)

		// Write to file
		fileMode := os.FileMode(0644)
		if strings.ToLower(*output) == "curl" {
			fileMode = 0755 // Make curl script executable
		}
		err = os.WriteFile(fullPath, outputData, fileMode)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully generated %s\n", fullPath)
	}

	outputName := strings.Title(*output)
	if outputName == "Curl" {
		outputName = "CURL"
	}
	if *full {
		fmt.Printf("Generated %d environment file(s) for %s\n", len(environments), outputName)
	} else {
		fmt.Printf("You can now use this file with %s\n", outputName)
	}
}

