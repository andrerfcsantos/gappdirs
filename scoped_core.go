package gappdirs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type scopedContext struct {
	appName        string
	scope          Scope
	workingDir     string
	defaultDirPerm fs.FileMode
	userDirsFn     dirLookupFunc
	systemDirsFn   dirLookupFunc
}

func buildScopedContext(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc, forcedScope *Scope) scopedContext {
	appName = sanitizeAppName(appName)

	cfg := defaultConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}
	if forcedScope != nil {
		switch *forcedScope {
		case ScopeLocal, ScopeUser, ScopeSystem:
			cfg.scope = *forcedScope
		default:
			cfg.scope = ScopeUser
		}
	}

	return scopedContext{
		appName:        appName,
		scope:          cfg.scope,
		workingDir:     strings.TrimSpace(cfg.workingDir),
		defaultDirPerm: cfg.defaultDirPerm,
		userDirsFn:     userDirsFn,
		systemDirsFn:   systemDirsFn,
	}
}

func newDefaultScopedContext(appName string, scope Scope) scopedContext {
	fixedScope := scope
	return buildScopedContext(appName, nil, platformUserDirs, platformSystemDirs, &fixedScope)
}

func scopedDirs(ctx scopedContext, cat category) []string {
	var candidates []string
	switch ctx.scope {
	case ScopeLocal:
		localDir := scopedLocalDir(ctx, cat)
		candidates = append(candidates, localDir)
		fallthrough
	case ScopeUser:
		userDirs, _ := ctx.userDirsFn(ctx.appName, cat)
		candidates = append(candidates, userDirs...)
		fallthrough
	case ScopeSystem:
		systemDirs, _ := ctx.systemDirsFn(ctx.appName, cat)
		candidates = append(candidates, systemDirs...)
	default:
		userDirs, _ := ctx.userDirsFn(ctx.appName, cat)
		systemDirs, _ := ctx.systemDirsFn(ctx.appName, cat)
		candidates = append(candidates, userDirs...)
		candidates = append(candidates, systemDirs...)
	}

	normalized := normalizePaths(candidates)

	return normalized
}

func scopedDir(ctx scopedContext, cat category) string {
	dirs := scopedDirs(ctx, cat)
	if len(dirs) == 0 {
		return ""
	}
	return dirs[0]
}

func scopedEnsureDir(ctx scopedContext, cat category, opts ...EnsureOption) (string, error) {
	dir := scopedDir(ctx, cat)
	cfg := resolveEnsureConfig(ctx.defaultDirPerm, opts)
	if err := os.MkdirAll(dir, cfg.dirPerm); err != nil {
		return "", fmt.Errorf("gappdirs: create directory %q: %w", dir, err)
	}
	return dir, nil
}

func findExistingScopedFiles(ctx scopedContext, cat category, filename string) ([]string, error) {
	dirs := scopedDirs(ctx, cat)

	matches := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		var normalizedFileName string
		var err error
		if normalizedFileName, err = normalizeFilename(dir, filename); err != nil {
			return nil, err
		}
		pathToFile := filepath.Join(dir, normalizedFileName)
		exists, err := regularFileExists(pathToFile)
		if err != nil {
			return nil, err
		}
		if exists {
			matches = append(matches, pathToFile)
		}
	}
	return matches, nil
}

func findExistingScopedFile(ctx scopedContext, cat category, filename string) (string, error) {

	dirs := scopedDirs(ctx, cat)
	for _, dir := range dirs {
		var normalizedFileName string
		var err error
		if normalizedFileName, err = normalizeFilename(dir, filename); err != nil {
			return "", err
		}
		pathToFile := filepath.Join(dir, normalizedFileName)
		exists, err := regularFileExists(pathToFile)
		if err != nil {
			return "", err
		}

		if exists {
			return pathToFile, nil
		}
	}
	return "", ErrNotFound
}

func scopedLocalDir(ctx scopedContext, cat category) string {
	workingDir := resolveLocalWorkingDir(ctx.workingDir)
	return filepath.Join(workingDir, "."+ctx.appName, categoryDirName(cat))
}

func resolveLocalWorkingDir(configuredWorkingDir string) string {
	configuredWorkingDir = strings.TrimSpace(configuredWorkingDir)

	if configuredWorkingDir == "" || configuredWorkingDir == "." {
		wd, err := os.Getwd()
		if err != nil {
			configuredWorkingDir = "."
		} else {
			configuredWorkingDir = wd
		}
	}

	return configuredWorkingDir
}
