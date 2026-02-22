package gappdirs

import (
	"io/fs"
	"strings"
)

const defaultDirPerm fs.FileMode = 0o755

// Option configures resolver creation.
type Option func(*newConfig)

type newConfig struct {
	scope          Scope
	workingDir     string
	defaultDirPerm fs.FileMode
}

func defaultConfig() newConfig {
	return newConfig{
		scope:          ScopeUser,
		defaultDirPerm: defaultDirPerm,
	}
}

// WithScope sets resolver lookup scope.
func WithScope(scope Scope) Option {
	return func(cfg *newConfig) {
		if cfg == nil {
			return
		}
		switch scope {
		case ScopeLocal, ScopeUser, ScopeSystem:
			cfg.scope = scope
		default:
			cfg.scope = ScopeUser
		}
	}
}

// WithWorkingDir sets the resolver working directory root used for local scope.
func WithWorkingDir(dir string) Option {
	return func(cfg *newConfig) {
		if cfg == nil {
			return
		}
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return
		}
		cfg.workingDir = dir
	}
}

// WithDefaultDirPerm sets the default permission used by Ensure*Dir helpers.
func WithDefaultDirPerm(perm fs.FileMode) Option {
	return func(cfg *newConfig) {
		if cfg == nil {
			return
		}
		if perm == 0 || perm&fs.ModeType != 0 {
			cfg.defaultDirPerm = defaultDirPerm
			return
		}
		cfg.defaultDirPerm = perm
	}
}
