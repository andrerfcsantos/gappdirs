package gappdirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindDataFiles returns all existing matches for filename in data directories by precedence order.
//
// It returns thr paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func (r *Resolver) FindDataFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryData, filename)
}

// FindConfigFiles returns all existing matches for filename in config directories by precedence order.
//
// It returns thr paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func (r *Resolver) FindConfigFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryConfig, filename)
}

// FindLogFiles returns all existing matches for filename in log directories by precedence order.
//
// It returns thr paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func (r *Resolver) FindLogFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryLog, filename)
}

// FindCacheFiles returns all existing matches for filename in cache directories by precedence order.
//
// It returns thr paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func (r *Resolver) FindCacheFiles(filename string) ([]string, error) {
	return findExistingScopedFiles(r.ctx, categoryCache, filename)
}

// FindDataFile returns the first existing match for filename in data directories by precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, an error of type ErrNotFound is returned.
func (r *Resolver) FindDataFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryData, filename)
}

// FindConfigFile returns the first existing match for filename in config directories by precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, an error of type ErrNotFound is returned.
func (r *Resolver) FindConfigFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryConfig, filename)
}

// FindLogFile returns the first existing match for filename in log directories by precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, an error of type ErrNotFound is returned.
func (r *Resolver) FindLogFile(filename string) (string, error) {
	return findExistingScopedFile(r.ctx, categoryLog, filename)
}

// FindCacheFile returns the first existing match for filename in cache directories by precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, an error of type ErrNotFound is returned.
func (r *Resolver) FindCacheFile(filename string) (string, error) {
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
