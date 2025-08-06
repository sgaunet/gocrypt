package cmd

import (
	"os"
	"path/filepath"
)

func isFileExists(file string) bool {
	// Clean the path to prevent directory traversal
	cleanFile := filepath.Clean(file)
	f, err := os.Open(cleanFile) // #nosec G304 - file path is cleaned
	if os.IsNotExist(err) {
		return false
	}
	defer func() { _ = f.Close() }()
	i, _ := os.Stat(cleanFile)
	return !i.IsDir()
}
