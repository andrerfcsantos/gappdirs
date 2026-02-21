package gappdirs

import (
	"errors"
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

func buildScopedContext(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc, forcedScope *Scope) (scopedContext, error) {
	if userDirsFn == nil || systemDirsFn == nil {
		return scopedContext{}, errors.New("gappdirs: internal directory provider is nil")
	}

	var err error
	appName, err = sanitizeAppName(appName)
	if err != nil {
		return scopedContext{}, err
	}

	cfg := defaultConfig()
	for i, opt := range opts {
		if opt == nil {
			return scopedContext{}, fmt.Errorf("gappdirs: option at index %d is nil", i)
		}
		if err := opt(&cfg); err != nil {
			return scopedContext{}, fmt.Errorf("gappdirs: apply option %d: %w", i, err)
		}
	}
	if forcedScope != nil {
		cfg.scope = *forcedScope
	}

	workingDir := strings.TrimSpace(cfg.workingDir)
	if workingDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return scopedContext{}, fmt.Errorf("gappdirs: determine working directory: %w", err)
		}
		workingDir = wd
	}

	absWorkingDir, err := normalizeAbsolutePath(workingDir)
	if err != nil {
		return scopedContext{}, fmt.Errorf("gappdirs: normalize working directory: %w", err)
	}

	return scopedContext{
		appName:        appName,
		scope:          cfg.scope,
		workingDir:     absWorkingDir,
		defaultDirPerm: cfg.defaultDirPerm,
		userDirsFn:     userDirsFn,
		systemDirsFn:   systemDirsFn,
	}, nil
}

func newDefaultScopedContext(appName string, scope Scope) (scopedContext, error) {
	fixedScope := scope
	return buildScopedContext(appName, nil, platformUserDirs, platformSystemDirs, &fixedScope)
}

func validateScopedContext(ctx scopedContext) error {
	if ctx.userDirsFn == nil || ctx.systemDirsFn == nil {
		return errors.New("gappdirs: resolver is not initialized")
	}
	return nil
}

func scopedDirs(ctx scopedContext, cat category) ([]string, error) {
	if err := validateScopedContext(ctx); err != nil {
		return nil, err
	}
	if err := validateCategory(cat); err != nil {
		return nil, err
	}
	if err := validateScope(ctx.scope); err != nil {
		return nil, err
	}

	var candidates []string
	switch ctx.scope {
	case ScopeLocal:
		candidates = append(candidates, scopedLocalDir(ctx, cat))
		fallthrough
	case ScopeUser:
		userDirs, err := ctx.userDirsFn(ctx.appName, cat)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, userDirs...)
		fallthrough
	case ScopeSystem:
		systemDirs, err := ctx.systemDirsFn(ctx.appName, cat)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, systemDirs...)
	default:
		return nil, fmt.Errorf("gappdirs: unsupported scope %d", ctx.scope)
	}

	normalized, err := normalizeAndDedupe(candidates)
	if err != nil {
		return nil, err
	}
	if len(normalized) == 0 {
		return nil, fmt.Errorf("gappdirs: no directories resolved for category=%s scope=%s", cat, ctx.scope)
	}
	return normalized, nil
}

func scopedDir(ctx scopedContext, cat category) (string, error) {
	dirs, err := scopedDirs(ctx, cat)
	if err != nil {
		return "", err
	}
	return dirs[0], nil
}

func scopedEnsureDir(ctx scopedContext, cat category) (string, error) {
	return scopedEnsureDirWithPerm(ctx, cat, ctx.defaultDirPerm)
}

func scopedEnsureDirWithPerm(ctx scopedContext, cat category, perm fs.FileMode) (string, error) {
	dir, err := scopedDir(ctx, cat)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, perm); err != nil {
		return "", fmt.Errorf("gappdirs: create directory %q: %w", dir, err)
	}
	return dir, nil
}

func scopedFindFileDirs(ctx scopedContext, cat category, filename string) ([]string, error) {
	if err := validateFilename(filename); err != nil {
		return nil, err
	}
	if err := validateCategory(cat); err != nil {
		return nil, err
	}

	dirs, err := scopedDirs(ctx, cat)
	if err != nil {
		return nil, err
	}

	matches := make([]string, 0, len(dirs))
	for _, dir := range dirs {
		exists, err := regularFileExists(filepath.Join(dir, filename))
		if err != nil {
			return nil, err
		}
		if exists {
			matches = append(matches, dir)
		}
	}
	return matches, nil
}

func scopedFile(ctx scopedContext, cat category, filename string) (string, error) {
	if err := validateFilename(filename); err != nil {
		return "", err
	}
	if err := validateCategory(cat); err != nil {
		return "", err
	}

	dirs, err := scopedDirs(ctx, cat)
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		candidate := filepath.Join(dir, filename)
		exists, err := regularFileExists(candidate)
		if err != nil {
			return "", err
		}
		if exists {
			return candidate, nil
		}
	}
	return "", ErrNotFound
}

func scopedLocalDir(ctx scopedContext, cat category) string {
	return filepath.Join(ctx.workingDir, "."+ctx.appName, categoryDirName(cat))
}
