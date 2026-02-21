package gappdirs

import (
	"fmt"
	"io/fs"
	"strings"
)

const defaultDirPerm fs.FileMode = 0o755

// Option configures resolver creation.
type Option func(*newConfig) error

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
	return func(cfg *newConfig) error {
		if cfg == nil {
			return fmt.Errorf("gappdirs: nil option config")
		}
		if err := validateScope(scope); err != nil {
			return err
		}
		cfg.scope = scope
		return nil
	}
}

// WithWorkingDir sets the resolver working directory root used for local scope.
func WithWorkingDir(dir string) Option {
	return func(cfg *newConfig) error {
		if cfg == nil {
			return fmt.Errorf("gappdirs: nil option config")
		}
		dir = strings.TrimSpace(dir)
		if dir == "" {
			return fmt.Errorf("gappdirs: working directory is required")
		}
		cfg.workingDir = dir
		return nil
	}
}

// WithDefaultDirPerm sets the default permission used by Ensure*Dir helpers.
func WithDefaultDirPerm(perm fs.FileMode) Option {
	return func(cfg *newConfig) error {
		if cfg == nil {
			return fmt.Errorf("gappdirs: nil option config")
		}
		if perm == 0 {
			return fmt.Errorf("gappdirs: default directory permission must not be zero")
		}
		if perm&fs.ModeType != 0 {
			return fmt.Errorf("gappdirs: default directory permission must only contain permission bits")
		}
		cfg.defaultDirPerm = perm
		return nil
	}
}
