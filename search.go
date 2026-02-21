package gappdirs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FindDataFileDirs returns data directories that contain filename.
func (r *resolver) FindDataFileDirs(filename string) ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedFindFileDirs(ctx, categoryData, filename)
}

// FindConfigFileDirs returns config directories that contain filename.
func (r *resolver) FindConfigFileDirs(filename string) ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedFindFileDirs(ctx, categoryConfig, filename)
}

// FindLogFileDirs returns log directories that contain filename.
func (r *resolver) FindLogFileDirs(filename string) ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedFindFileDirs(ctx, categoryLog, filename)
}

// FindCacheFileDirs returns cache directories that contain filename.
func (r *resolver) FindCacheFileDirs(filename string) ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedFindFileDirs(ctx, categoryCache, filename)
}

// DataFile returns the first matching data file path by precedence.
func (r *resolver) DataFile(filename string) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedFile(ctx, categoryData, filename)
}

// ConfigFile returns the first matching config file path by precedence.
func (r *resolver) ConfigFile(filename string) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedFile(ctx, categoryConfig, filename)
}

// LogFile returns the first matching log file path by precedence.
func (r *resolver) LogFile(filename string) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedFile(ctx, categoryLog, filename)
}

// CacheFile returns the first matching cache file path by precedence.
func (r *resolver) CacheFile(filename string) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedFile(ctx, categoryCache, filename)
}

func validateFilename(filename string) error {
	if strings.TrimSpace(filename) == "" {
		return errors.New("gappdirs: filename is required")
	}
	if filepath.IsAbs(filename) {
		return fmt.Errorf("gappdirs: filename must be relative basename, got absolute path %q", filename)
	}
	if filename == "." || filename == ".." {
		return fmt.Errorf("gappdirs: invalid filename %q", filename)
	}
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		return fmt.Errorf("gappdirs: filename must not include path separators: %q", filename)
	}
	if base := filepath.Base(filename); base != filename {
		return fmt.Errorf("gappdirs: filename must be basename, got %q", filename)
	}
	return nil
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
