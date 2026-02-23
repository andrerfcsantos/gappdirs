package gappdirs

import (
	"io/fs"
	"strings"
)

// DefaultDirPerm is the default directory permission used for directory creation operations.
const DefaultDirPerm fs.FileMode = 0o755

// ResolverOption sets an option when creating a resolver.
// It allows adjusting the local directories and the default directory permissions of a resolver.
type ResolverOption func(*newResolverConfig)

type newResolverConfig struct {
	workingDirs    []string
	defaultDirPerm fs.FileMode
}

func defaultConfig() newResolverConfig {
	return newResolverConfig{
		defaultDirPerm: DefaultDirPerm,
	}
}

func (cfg *newResolverConfig) appendWorkingDirs(dirs ...string) {
	if cfg == nil {
		return
	}

	seen := make(map[string]struct{}, len(cfg.workingDirs)+len(dirs))
	for _, existing := range cfg.workingDirs {
		seen[existing] = struct{}{}
	}

	for _, dir := range dirs {
		dir = strings.TrimSpace(dir)
		if dir == "" {
			continue
		}
		if _, exists := seen[dir]; exists {
			continue
		}
		cfg.workingDirs = append(cfg.workingDirs, dir)
		seen[dir] = struct{}{}
	}
}

// WithLocalDir sets a local working directory root used for local scope.
//
// This option can be used multiple times to set multiple local directories.
// Directories given first have highest precedence.
// Duplicate directories are ignored and only the first occurrence is used.
func WithLocalDir(dir string) ResolverOption {
	return func(cfg *newResolverConfig) {
		cfg.appendWorkingDirs(dir)
	}
}

// WithLocalDirs sets one or multiple local working directory roots used for local scope.
//
// Directories given first have highest precedence.
// Duplicate directories are ignored and only the first occurrence is used.
func WithLocalDirs(dirs ...string) ResolverOption {
	return func(cfg *newResolverConfig) {
		cfg.appendWorkingDirs(dirs...)
	}
}

// WithDefaultDirPerm sets the default permission used when creating a directory.
//
// Invalid permissions are ignored and the default is used instead.
func WithDefaultDirPerm(perm fs.FileMode) ResolverOption {
	return func(cfg *newResolverConfig) {
		if cfg == nil {
			return
		}
		if perm == 0 || perm&fs.ModeType != 0 {
			cfg.defaultDirPerm = DefaultDirPerm
			return
		}
		cfg.defaultDirPerm = perm
	}
}
