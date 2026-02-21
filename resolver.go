package gappdirs

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"unicode"
)

// Scope determines which directory layers are searched.
type Scope int

const (
	// ScopeLocal searches current working directory, then user directories, then system directories.
	ScopeLocal Scope = iota
	// ScopeUser searches user directories, then system directories.
	ScopeUser
	// ScopeSystem searches only system directories.
	ScopeSystem
)

type dirLookupFunc func(appName string, cat category) ([]string, error)

// Resolver is the public interface for resolving application directories.
type Resolver interface {
	DataDirs() ([]string, error)
	ConfigDirs() ([]string, error)
	LogDirs() ([]string, error)
	CacheDirs() ([]string, error)
	DataDir() (string, error)
	ConfigDir() (string, error)
	LogDir() (string, error)
	CacheDir() (string, error)
	EnsureDataDir() (string, error)
	EnsureDataDirWithPerm(perm fs.FileMode) (string, error)
	EnsureConfigDir() (string, error)
	EnsureConfigDirWithPerm(perm fs.FileMode) (string, error)
	EnsureLogDir() (string, error)
	EnsureLogDirWithPerm(perm fs.FileMode) (string, error)
	EnsureCacheDir() (string, error)
	EnsureCacheDirWithPerm(perm fs.FileMode) (string, error)
	FindDataFileDirs(filename string) ([]string, error)
	FindConfigFileDirs(filename string) ([]string, error)
	FindLogFileDirs(filename string) ([]string, error)
	FindCacheFileDirs(filename string) ([]string, error)
	DataFile(filename string) (string, error)
	ConfigFile(filename string) (string, error)
	LogFile(filename string) (string, error)
	CacheFile(filename string) (string, error)
}

// resolver resolves application directories for the current OS.
type resolver struct {
	ctx scopedContext
}

var ErrNotFound = errors.New("gappdirs: no matching path found")

// ensures that resolver implements the Resolver interface at compile time.
var _ Resolver = (*resolver)(nil)

// NewResolver creates a resolver with functional options.
func NewResolver(appName string, opts ...Option) (Resolver, error) {
	return newResolver(appName, opts, platformUserDirs, platformSystemDirs)
}

// NewUserResolver creates a resolver that always uses user scope.
func NewUserResolver(appName string, opts ...Option) (Resolver, error) {
	fixedScope := ScopeUser
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

// NewSystemResolver creates a resolver that always uses system scope.
func NewSystemResolver(appName string, opts ...Option) (Resolver, error) {
	fixedScope := ScopeSystem
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

// NewLocalResolver creates a resolver that always uses local scope.
func NewLocalResolver(appName string, opts ...Option) (Resolver, error) {
	fixedScope := ScopeLocal
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

func newResolver(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc) (*resolver, error) {
	return newResolverWithScope(appName, opts, userDirsFn, systemDirsFn, nil)
}

func newResolverWithScope(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc, forcedScope *Scope) (*resolver, error) {
	ctx, err := buildScopedContext(appName, opts, userDirsFn, systemDirsFn, forcedScope)
	if err != nil {
		return nil, err
	}
	return &resolver{ctx: ctx}, nil
}

// DataDirs returns all data directories in precedence order for the resolver scope.
func (r *resolver) DataDirs() ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedDirs(ctx, categoryData)
}

// ConfigDirs returns all config directories in precedence order for the resolver scope.
func (r *resolver) ConfigDirs() ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedDirs(ctx, categoryConfig)
}

// LogDirs returns all log directories in precedence order for the resolver scope.
func (r *resolver) LogDirs() ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedDirs(ctx, categoryLog)
}

// CacheDirs returns all cache directories in precedence order for the resolver scope.
func (r *resolver) CacheDirs() ([]string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return nil, err
	}
	return scopedDirs(ctx, categoryCache)
}

// DataDir returns the first data directory by precedence.
func (r *resolver) DataDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedDir(ctx, categoryData)
}

// ConfigDir returns the first config directory by precedence.
func (r *resolver) ConfigDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedDir(ctx, categoryConfig)
}

// LogDir returns the first log directory by precedence.
func (r *resolver) LogDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedDir(ctx, categoryLog)
}

// CacheDir returns the first cache directory by precedence.
func (r *resolver) CacheDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedDir(ctx, categoryCache)
}

func (r *resolver) scopedCtx() (scopedContext, error) {
	if err := r.validateReady(); err != nil {
		return scopedContext{}, err
	}
	return r.ctx, nil
}

func sanitizeAppName(appName string) (string, error) {
	appName = strings.TrimSpace(strings.ToLower(appName))
	if appName == "" {
		return "", errors.New("gappdirs: app name is required")
	}
	if appName == "." {
		return "_", nil
	}
	if appName == ".." {
		return "__", nil
	}

	var out strings.Builder
	out.Grow(len(appName))
	for _, r := range appName {
		if unicode.IsSpace(r) || isInvalidAppNameRune(r) {
			out.WriteByte('_')
			continue
		}
		out.WriteRune(r)
	}

	sanitized := collapseUnderscores(out.String())
	sanitized = strings.Trim(sanitized, "_")
	if sanitized == "" {
		return "", errors.New("gappdirs: app name is required")
	}
	return sanitized, nil
}

func isInvalidAppNameRune(r rune) bool {
	switch r {
	case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
		return true
	default:
		return false
	}
}

func collapseUnderscores(value string) string {
	var out strings.Builder
	out.Grow(len(value))
	prevUnderscore := false

	for _, r := range value {
		if r == '_' {
			if prevUnderscore {
				continue
			}
			prevUnderscore = true
		} else {
			prevUnderscore = false
		}
		out.WriteRune(r)
	}

	return out.String()
}

func validateScope(scope Scope) error {
	switch scope {
	case ScopeLocal, ScopeUser, ScopeSystem:
		return nil
	default:
		return fmt.Errorf("gappdirs: unsupported scope %d", scope)
	}
}

func (r *resolver) validateReady() error {
	if r == nil {
		return errors.New("gappdirs: nil resolver")
	}
	if r.ctx.userDirsFn == nil || r.ctx.systemDirsFn == nil {
		return errors.New("gappdirs: resolver is not initialized")
	}
	return nil
}
