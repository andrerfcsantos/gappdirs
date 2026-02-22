package gappdirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindDataFiles returns matching data file paths in precedence order.
func (r *resolver) FindDataFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryData, filename)
}

// FindConfigFiles returns matching config file paths in precedence order.
func (r *resolver) FindConfigFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryConfig, filename)
}

// FindLogFiles returns matching log file paths in precedence order.
func (r *resolver) FindLogFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryLog, filename)
}

// FindCacheFiles returns matching cache file paths in precedence order.
func (r *resolver) FindCacheFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryCache, filename)
}

// FindDataFile returns the first matching data file path by precedence.
func (r *resolver) FindDataFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryData, filename)
}

// FindConfigFile returns the first matching config file path by precedence.
func (r *resolver) FindConfigFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryConfig, filename)
}

// FindLogFile returns the first matching log file path by precedence.
func (r *resolver) FindLogFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryLog, filename)
}

// FindCacheFile returns the first matching cache file path by precedence.
func (r *resolver) FindCacheFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryCache, filename)
}

func normalizeFilename(base, filename string) (string, error) {
	if strings.TrimSpace(filename) == "" {
		return filename, errors.New("gappdirs: filename is required")
	}
	if filepath.IsAbs(filename) {
		rel, err := filepath.Rel(base, filename)
		if err != nil {
			return filename, fmt.Errorf("gappdirs: filename is an absolute path and it is not possible to get a relative path from %q to %q: %w", base, filename, err)
		}

		filename = rel
	}
	return filename, nil
}

func regularFileExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("gappdirs: stat %q: %w", path, err)
}
