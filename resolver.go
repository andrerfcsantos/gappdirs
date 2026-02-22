package gappdirs

import (
	"errors"
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
	DataDirs() []string
	ConfigDirs() []string
	LogDirs() []string
	CacheDirs() []string
	DataDir() string
	ConfigDir() string
	LogDir() string
	CacheDir() string
	EnsureDataDir(opts ...EnsureOption) (string, error)
	EnsureConfigDir(opts ...EnsureOption) (string, error)
	EnsureLogDir(opts ...EnsureOption) (string, error)
	EnsureCacheDir(opts ...EnsureOption) (string, error)
	FindDataFiles(filename string) ([]string, error)
	FindConfigFiles(filename string) ([]string, error)
	FindLogFiles(filename string) ([]string, error)
	FindCacheFiles(filename string) ([]string, error)
	FindDataFile(filename string) (string, error)
	FindConfigFile(filename string) (string, error)
	FindLogFile(filename string) (string, error)
	FindCacheFile(filename string) (string, error)
}

// resolver resolves application directories for the current OS.
type resolver struct {
	ctx scopedContext
}

var ErrNotFound = errors.New("gappdirs: no matching path found")

// ensures that resolver implements the Resolver interface at compile time.
var _ Resolver = (*resolver)(nil)

// NewResolver creates a resolver with functional options.
func NewResolver(appName string, opts ...Option) Resolver {
	return newResolver(appName, opts, platformUserDirs, platformSystemDirs)
}

// NewUserResolver creates a resolver that always uses user scope.
func NewUserResolver(appName string, opts ...Option) Resolver {
	fixedScope := ScopeUser
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

// NewSystemResolver creates a resolver that always uses system scope.
func NewSystemResolver(appName string, opts ...Option) Resolver {
	fixedScope := ScopeSystem
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

// NewLocalResolver creates a resolver that always uses local scope.
func NewLocalResolver(appName string, opts ...Option) Resolver {
	fixedScope := ScopeLocal
	return newResolverWithScope(appName, opts, platformUserDirs, platformSystemDirs, &fixedScope)
}

func newResolver(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc) *resolver {
	return newResolverWithScope(appName, opts, userDirsFn, systemDirsFn, nil)
}

func newResolverWithScope(appName string, opts []Option, userDirsFn, systemDirsFn dirLookupFunc, forcedScope *Scope) *resolver {
	ctx := buildScopedContext(appName, opts, userDirsFn, systemDirsFn, forcedScope)
	return &resolver{ctx: ctx}
}

// DataDirs returns all data directories in precedence order for the resolver scope.
func (r *resolver) DataDirs() []string {
	return scopedDirs(r.ctx, categoryData)
}

// ConfigDirs returns all config directories in precedence order for the resolver scope.
func (r *resolver) ConfigDirs() []string {
	return scopedDirs(r.ctx, categoryConfig)
}

// LogDirs returns all log directories in precedence order for the resolver scope.
func (r *resolver) LogDirs() []string {
	return scopedDirs(r.ctx, categoryLog)
}

// CacheDirs returns all cache directories in precedence order for the resolver scope.
func (r *resolver) CacheDirs() []string {
	return scopedDirs(r.ctx, categoryCache)
}

// DataDir returns the first data directory by precedence.
func (r *resolver) DataDir() string {
	return scopedDir(r.ctx, categoryData)
}

// ConfigDir returns the first config directory by precedence.
func (r *resolver) ConfigDir() string {
	return scopedDir(r.ctx, categoryConfig)
}

// LogDir returns the first log directory by precedence.
func (r *resolver) LogDir() string {
	return scopedDir(r.ctx, categoryLog)
}

// CacheDir returns the first cache directory by precedence.
func (r *resolver) CacheDir() string {
	return scopedDir(r.ctx, categoryCache)
}

func (r *resolver) scopedCtx() (scopedContext, error) {
	return r.ctx, nil
}

func sanitizeAppName(appName string) string {
	appName = strings.TrimSpace(appName)
	if appName == "" {
		return "unnamed_app"
	}

	runes := []rune(appName)
	for i, r := range runes {
		if unicode.IsSpace(r) || isInvalidAppNameRune(r) {
			runes[i] = '_'
		}
	}

	sanitized := string(runes)
	if sanitized == "" {
		return "unnamed_app"
	}
	return sanitized
}

func isInvalidAppNameRune(r rune) bool {
	switch r {
	case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
		return true
	default:
		return false
	}
}

func (r *resolver) validateReady() error {
	if r == nil {
		return errors.New("gappdirs: nil resolver")
	}
	return nil
}
