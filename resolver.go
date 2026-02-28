package gappdirs

import (
	"errors"
	"strings"
	"unicode"
)

// Scope determines which directory layers are searched for each lookup.
//
// It allows callers to control the level of visibility and persistence for their
// files, choosing between project-local, per-user, or machine-wide system directories.
type Scope int

const (
	// ScopeLocal includes the current working directory as the highest priority,
	// falling back to user directories, and then system directories.
	ScopeLocal Scope = iota
	// ScopeUser includes user-specific directories as the highest priority,
	// falling back to system directories.
	ScopeUser
	// ScopeSystem includes only machine-wide system directories, ignoring user
	// and local paths.
	ScopeSystem
)

type dirLookupFunc func(appName string, cat category) ([]string, error)

// Resolver resolves application directories for the current platform.
// It can be used for several operations of lookup and creation of application files and directories for the same scope and app name.
// If you just need a single operation, you might want to use the global functions of the module instead.
type Resolver struct {
	ctx scopedContext
}

// ErrNotFound indicates that a requested file was not found in any of the searched directories.
//
// It is returned by all Find*File methods to allow callers to handle missing
// files—like falling back to defaults or creating new files—without performing
// string matching on errors.
var ErrNotFound = errors.New("gappdirs: no matching path found")

// NewUserResolver creates a Resolver configured to search user directories first,
// followed by system directories.
//
// It is the standard choice for most applications that store per-user configuration,
// data, and caches, while optionally reading machine-wide system defaults as a fallback.
func NewUserResolver(appName string, opts ...ResolverOption) *Resolver {
	return newScopedResolver(appName, ScopeUser, opts, platformUserDirs, platformSystemDirs)
}

// NewSystemResolver creates a Resolver configured to search only machine-wide system directories.
//
// It is useful for system services, daemons, or tools that manage shared assets
// and must intentionally ignore any user-specific or local overrides.
func NewSystemResolver(appName string, opts ...ResolverOption) *Resolver {
	return newScopedResolver(appName, ScopeSystem, opts, platformUserDirs, platformSystemDirs)
}

// NewLocalResolver creates a Resolver configured to search local directories first,
// followed by user and then system directories.
//
// It is typically used for development tools or applications that support project-local
// configuration overrides while keeping standard fallback behavior.
// Use the WithLocalDir or WithLocalDirs options to set the local directories.
func NewLocalResolver(appName string, opts ...ResolverOption) *Resolver {
	return newScopedResolver(appName, ScopeLocal, opts, platformUserDirs, platformSystemDirs)
}

func newScopedResolver(appName string, scope Scope, opts []ResolverOption, userDirsFn, systemDirsFn dirLookupFunc) *Resolver {
	ctx := buildScopedContext(appName, scope, opts, userDirsFn, systemDirsFn)
	return &Resolver{ctx: ctx}
}

// DataDirs returns all candidate data directories paths, in precedence order for the resolver's scope and app name.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureDataDir to ensure the highest-precedence data directory exists and get its path.
func (r *Resolver) DataDirs() []string {
	return scopedDirs(r.ctx, categoryData)
}

// ConfigDirs returns all candidate config directories paths, in precedence order for the resolver's scope and app name.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureConfigDir to ensure the highest-precedence config directory exists and get its path.
func (r *Resolver) ConfigDirs() []string {
	return scopedDirs(r.ctx, categoryConfig)
}

// LogDirs returns all candidate log directories paths, in precedence order for the resolver's scope and app name.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureLogDir to ensure the highest-precedence log directory exists and get its path.
func (r *Resolver) LogDirs() []string {
	return scopedDirs(r.ctx, categoryLog)
}

// CacheDirs returns all candidate cache directories paths, in precedence order for the resolver's scope and app name.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureCacheDir to ensure the highest-precedence cache directory exists and get its path.
func (r *Resolver) CacheDirs() []string {
	return scopedDirs(r.ctx, categoryCache)
}

// DataDir returns the single highest-precedence data directory path for the resolver's scope and app name.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureDataDir to ensure the highest-precedence data directory exists and get its path.
func (r *Resolver) DataDir() string {
	return scopedDir(r.ctx, categoryData)
}

// ConfigDir returns the single highest-precedence config directory path for the resolver's scope and app name.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureConfigDir to ensure the highest-precedence config directory exists and get its path.
func (r *Resolver) ConfigDir() string {
	return scopedDir(r.ctx, categoryConfig)
}

// LogDir returns the single highest-precedence log directory path for the resolver's scope and app name.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureLogDir to ensure the highest-precedence log directory exists and get its path.
func (r *Resolver) LogDir() string {
	return scopedDir(r.ctx, categoryLog)
}

// CacheDir returns the single highest-precedence cache directory path for the resolver's scope and app name.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureCacheDir to ensure the highest-precedence cache directory exists and get its path.
func (r *Resolver) CacheDir() string {
	return scopedDir(r.ctx, categoryCache)
}

// DataFilePaths returns all candidate data file paths for filename, in precedence order for the resolver's scope and app name.
//
// The paths returned are computed candidates and are not guaranteed to exist on the filesystem.
func (r *Resolver) DataFilePaths(filename string) []string {
	return scopedFilePaths(r.ctx, categoryData, filename)
}

// ConfigFilePaths returns all candidate config file paths for filename, in precedence order for the resolver's scope and app name.
//
// The paths returned are computed candidates and are not guaranteed to exist on the filesystem.
func (r *Resolver) ConfigFilePaths(filename string) []string {
	return scopedFilePaths(r.ctx, categoryConfig, filename)
}

// LogFilePaths returns all candidate log file paths for filename, in precedence order for the resolver's scope and app name.
//
// The paths returned are computed candidates and are not guaranteed to exist on the filesystem.
func (r *Resolver) LogFilePaths(filename string) []string {
	return scopedFilePaths(r.ctx, categoryLog, filename)
}

// CacheFilePaths returns all candidate cache file paths for filename, in precedence order for the resolver's scope and app name.
//
// The paths returned are computed candidates and are not guaranteed to exist on the filesystem.
func (r *Resolver) CacheFilePaths(filename string) []string {
	return scopedFilePaths(r.ctx, categoryCache, filename)
}

// DataFilePath returns the single highest-precedence data file path for filename for the resolver's scope and app name.
//
// The path returned is a computed candidate and is not guaranteed to exist on the filesystem.
func (r *Resolver) DataFilePath(filename string) string {
	return scopedFilePath(r.ctx, categoryData, filename)
}

// ConfigFilePath returns the single highest-precedence config file path for filename for the resolver's scope and app name.
//
// The path returned is a computed candidate and is not guaranteed to exist on the filesystem.
func (r *Resolver) ConfigFilePath(filename string) string {
	return scopedFilePath(r.ctx, categoryConfig, filename)
}

// LogFilePath returns the single highest-precedence log file path for filename for the resolver's scope and app name.
//
// The path returned is a computed candidate and is not guaranteed to exist on the filesystem.
func (r *Resolver) LogFilePath(filename string) string {
	return scopedFilePath(r.ctx, categoryLog, filename)
}

// CacheFilePath returns the single highest-precedence cache file path for filename for the resolver's scope and app name.
//
// The path returned is a computed candidate and is not guaranteed to exist on the filesystem.
func (r *Resolver) CacheFilePath(filename string) string {
	return scopedFilePath(r.ctx, categoryCache, filename)
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
