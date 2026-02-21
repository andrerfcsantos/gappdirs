//go:build darwin

package gappdirs

import (
	"fmt"
	"os"
	"path/filepath"
)

func platformUserDirs(appName string, cat category) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("gappdirs: resolve home directory: %w", err)
	}

	libraryDir := filepath.Join(homeDir, "Library")
	switch cat {
	case categoryData:
		return []string{filepath.Join(libraryDir, "Application Support", appName, "data")}, nil
	case categoryConfig:
		return []string{filepath.Join(libraryDir, "Application Support", appName, "config")}, nil
	case categoryLog:
		return []string{filepath.Join(libraryDir, "Logs", appName)}, nil
	case categoryCache:
		return []string{filepath.Join(libraryDir, "Caches", appName)}, nil
	default:
		return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
	}
}

func platformSystemDirs(appName string, cat category) ([]string, error) {
	switch cat {
	case categoryData:
		return []string{filepath.Join("/Library", "Application Support", appName, "data")}, nil
	case categoryConfig:
		return []string{filepath.Join("/Library", "Application Support", appName, "config")}, nil
	case categoryLog:
		return []string{filepath.Join("/Library", "Logs", appName)}, nil
	case categoryCache:
		return []string{filepath.Join("/Library", "Caches", appName)}, nil
	default:
		return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
	}
}
