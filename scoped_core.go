package gappdirs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type scopedContext struct {
	appName        string
	scope          Scope
	workingDirs    []string
	defaultDirPerm fs.FileMode
	userDirsFn     dirLookupFunc
	systemDirsFn   dirLookupFunc
}

func buildScopedContext(appName string, scope Scope, opts []ResolverOption, userDirsFn, systemDirsFn dirLookupFunc) scopedContext {
	appName = sanitizeAppName(appName)

	cfg := defaultConfig()
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}

	switch scope {
	case ScopeLocal, ScopeUser, ScopeSystem:
	default:
		scope = ScopeUser
	}

	return scopedContext{
		appName:        appName,
		scope:          scope,
		workingDirs:    append([]string(nil), cfg.workingDirs...),
		defaultDirPerm: cfg.defaultDirPerm,
		userDirsFn:     userDirsFn,
		systemDirsFn:   systemDirsFn,
	}
}

func newDefaultScopedContext(appName string, scope Scope) scopedContext {
	return buildScopedContext(appName, scope, nil, platformUserDirs, platformSystemDirs)
}

func scopedDirs(ctx scopedContext, cat category) []string {
	var candidates []string
	switch ctx.scope {
	case ScopeLocal:
		localDirs := scopedLocalDirs(ctx, cat)
		candidates = append(candidates, localDirs...)
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

func scopedFilePaths(ctx scopedContext, cat category, filename string) []string {
	dirs := scopedDirs(ctx, cat)
	normalizedFilename := normalizeFilePathFilename(filename)
	paths := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		paths = append(paths, filepath.Join(dir, normalizedFilename))
	}
	return paths
}

func scopedFilePath(ctx scopedContext, cat category, filename string) string {
	paths := scopedFilePaths(ctx, cat, filename)
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func scopedEnsureDir(ctx scopedContext, cat category, opts ...EnsureOption) (string, error) {
	dir := scopedDir(ctx, cat)
	cfg := resolveEnsureConfig(ctx.defaultDirPerm, opts)
	if err := os.MkdirAll(dir, cfg.dirPerm); err != nil {
		return "", fmt.Errorf("gappdirs: create directory %q: %w", dir, err)
	}
	return dir, nil
}

func scopedCreateFile(ctx scopedContext, cat category, filename string, opts ...CreateFileOption) (bool, string, error) {
	dir, err := scopedEnsureDir(ctx, cat)
	if err != nil {
		return false, "", err
	}

	normalizedFileName, err := normalizeFilename(dir, filename)
	if err != nil {
		return false, "", err
	}
	pathToFile := filepath.Join(dir, normalizedFileName)

	parentDir := filepath.Dir(pathToFile)
	if err := os.MkdirAll(parentDir, ctx.defaultDirPerm); err != nil {
		return false, pathToFile, fmt.Errorf("gappdirs: create parent directories for %q: %w", pathToFile, err)
	}

	cfg := resolveCreateFileOptions(opts)

	exists, err := regularFileExists(pathToFile)
	if err != nil {
		return false, pathToFile, err
	}
	if exists && !cfg.overwriteExisting {
		return false, pathToFile, nil
	}

	flags := os.O_CREATE | os.O_WRONLY
	if cfg.overwriteExisting {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(pathToFile, flags, cfg.filePerm)
	if err != nil {
		if !cfg.overwriteExisting && errors.Is(err, os.ErrExist) {
			return false, pathToFile, nil
		}
		return false, pathToFile, fmt.Errorf("gappdirs: open file %q: %w", pathToFile, err)
	}

	created := !exists
	if cfg.reader != nil {
		if _, err := io.Copy(file, cfg.reader); err != nil {
			_ = file.Close()
			return created, pathToFile, fmt.Errorf("gappdirs: write file %q: %w", pathToFile, err)
		}
	}
	if err := file.Close(); err != nil {
		return created, pathToFile, fmt.Errorf("gappdirs: close file %q: %w", pathToFile, err)
	}

	return created, pathToFile, nil
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

func scopedLocalDirs(ctx scopedContext, cat category) []string {
	configuredWorkingDirs := ctx.workingDirs
	if len(configuredWorkingDirs) == 0 {
		configuredWorkingDirs = []string{""}
	}

	localDirs := make([]string, 0, len(configuredWorkingDirs))
	for _, configuredWorkingDir := range configuredWorkingDirs {
		workingDir := resolveLocalWorkingDir(configuredWorkingDir)
		localDirs = append(localDirs, filepath.Join(workingDir, "."+ctx.appName, categoryDirName(cat)))
	}

	return localDirs
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

func normalizeFilePathFilename(filename string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		return ""
	}
	if !filepath.IsAbs(filename) {
		return filename
	}

	volume := filepath.VolumeName(filename)
	filename = strings.TrimPrefix(filename, volume)
	filename = strings.TrimLeft(filename, `/\`)

	return filename
}

func scopedScopeForPath(ctx scopedContext, inputPath string) (Scope, bool) {
	normalizedPath, ok := normalizeScopeLookupPath(inputPath)
	if !ok {
		return ScopeLocal, false
	}

	for _, scopeLayer := range scopeLookupLayers(ctx.scope) {
		dirs := scopedLayerDirs(ctx, scopeLayer)
		for _, dir := range dirs {
			if isPathWithinBase(normalizedPath, dir) {
				return scopeLayer, true
			}
		}
	}

	return ScopeLocal, false
}

func scopeLookupLayers(scope Scope) []Scope {
	switch scope {
	case ScopeLocal:
		return []Scope{ScopeLocal, ScopeUser, ScopeSystem}
	case ScopeUser:
		return []Scope{ScopeUser, ScopeSystem}
	case ScopeSystem:
		return []Scope{ScopeSystem}
	default:
		return []Scope{ScopeUser, ScopeSystem}
	}
}

func scopedLayerDirs(ctx scopedContext, scopeLayer Scope) []string {
	categories := []category{categoryData, categoryConfig, categoryLog, categoryCache}
	candidates := make([]string, 0, len(categories))

	for _, cat := range categories {
		switch scopeLayer {
		case ScopeLocal:
			candidates = append(candidates, scopedLocalDirs(ctx, cat)...)
		case ScopeUser:
			userDirs, _ := ctx.userDirsFn(ctx.appName, cat)
			candidates = append(candidates, userDirs...)
		case ScopeSystem:
			systemDirs, _ := ctx.systemDirsFn(ctx.appName, cat)
			candidates = append(candidates, systemDirs...)
		}
	}

	return normalizePaths(candidates)
}

func normalizeScopeLookupPath(path string) (string, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", false
	}

	normalizedPath, err := tryNormalizeAbsolutePath(path)
	if err != nil {
		return filepath.Clean(path), true
	}
	return normalizedPath, true
}

func isPathWithinBase(path string, base string) bool {
	base = strings.TrimSpace(base)
	if base == "" {
		return false
	}
	if path == base {
		return true
	}

	rel, err := filepath.Rel(base, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	if rel == ".." {
		return false
	}

	parentPrefix := ".." + string(filepath.Separator)
	return !strings.HasPrefix(rel, parentPrefix)
}
