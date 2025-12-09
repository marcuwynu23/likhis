package traversal

import (
	"os"
	"path/filepath"
)

// BFSTraverse performs a breadth-first traversal of the directory tree
// Returns a slice of file paths in BFS order
func BFSTraverse(rootPath string) ([]string, error) {
	var files []string
	queue := []string{rootPath}

	// Common directories to skip
	skipDirs := map[string]bool{
		"node_modules": true,
		".git":         true,
		"vendor":       true,
		"__pycache__":  true,
		".venv":        true,
		"venv":         true,
		"target":       true,
		"build":        true,
		"dist":         true,
		".idea":        true,
		".vscode":      true,
	}

	// Common file extensions to process
	extensions := map[string]bool{
		".js":   true,
		".ts":   true,
		".py":   true,
		".java": true,
		".php":  true,
		".go":   true,
	}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		entries, err := os.ReadDir(current)
		if err != nil {
			continue // Skip directories we can't read
		}

		for _, entry := range entries {
			fullPath := filepath.Join(current, entry.Name())

			if entry.IsDir() {
				// Skip common dependency/build directories
				if skipDirs[entry.Name()] {
					continue
				}
				queue = append(queue, fullPath)
			} else {
				// Only include files with relevant extensions
				ext := filepath.Ext(entry.Name())
				if extensions[ext] {
					files = append(files, fullPath)
				}
			}
		}
	}

	return files, nil
}

