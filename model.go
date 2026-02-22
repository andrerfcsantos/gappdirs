package gappdirs

import (
	"fmt"
	"path/filepath"
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

func normalizePaths(paths []string) []string {
	out := make([]string, 0, len(paths))

	for _, path := range paths {
		absPath, _ := tryNormalizeAbsolutePath(path)
		out = append(out, absPath)
	}

	return out
}

func tryNormalizeAbsolutePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return path, fmt.Errorf("gappdirs: normalize path %q: %w", path, err)
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
