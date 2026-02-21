package gappdirs

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type category int

const (
	categoryData category = iota
	categoryConfig
	categoryLog
	categoryCache
)

func (s Scope) String() string {
	switch s {
	case ScopeLocal:
		return "local"
	case ScopeUser:
		return "user"
	case ScopeSystem:
		return "system"
	default:
		return fmt.Sprintf("scope(%d)", int(s))
	}
}

func (c category) String() string {
	switch c {
	case categoryData:
		return "data"
	case categoryConfig:
		return "config"
	case categoryLog:
		return "log"
	case categoryCache:
		return "cache"
	default:
		return fmt.Sprintf("category(%d)", int(c))
	}
}

func categoryDirName(c category) string {
	switch c {
	case categoryData:
		return "data"
	case categoryConfig:
		return "config"
	case categoryLog:
		return "log"
	case categoryCache:
		return "cache"
	default:
		return ""
	}
}

func validateCategory(c category) error {
	switch c {
	case categoryData, categoryConfig, categoryLog, categoryCache:
		return nil
	default:
		return fmt.Errorf("gappdirs: unsupported category %d", c)
	}
}

func normalizeAndDedupe(paths []string) ([]string, error) {
	out := make([]string, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))

	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}

		absPath, err := normalizeAbsolutePath(path)
		if err != nil {
			return nil, err
		}

		key := absPath
		if runtime.GOOS == "windows" {
			key = strings.ToLower(absPath)
		}
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		out = append(out, absPath)
	}

	return out, nil
}

func normalizeAbsolutePath(path string) (string, error) {
	clean := filepath.Clean(path)
	absPath, err := filepath.Abs(clean)
	if err != nil {
		return "", fmt.Errorf("gappdirs: normalize path %q: %w", path, err)
	}
	return absPath, nil
}

func splitAbsolutePathList(pathList string) []string {
	rawItems := filepath.SplitList(pathList)
	out := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		if item == "" || !filepath.IsAbs(item) {
			continue
		}
		out = append(out, filepath.Clean(item))
	}
	return out
}
