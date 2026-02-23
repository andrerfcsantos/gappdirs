package gappdirs

import "io/fs"

// EnsureOption configures options for directory operations.
// It allows adjusting directory creation details for a single directory operation.
type EnsureOption func(*ensureConfig)

type ensureConfig struct {
	dirPerm fs.FileMode
}

func defaultEnsureConfig(defaultPerm fs.FileMode) ensureConfig {
	return ensureConfig{dirPerm: defaultPerm}
}

// WithEnsureDirPerm sets the directory permission when creating a directory.
//
// Invalid permissions are ignored and the default is used instead.
func WithEnsureDirPerm(perm fs.FileMode) EnsureOption {
	return func(cfg *ensureConfig) {
		if cfg == nil {
			return
		}
		if !isValidEnsureDirPerm(perm) {
			return
		}
		cfg.dirPerm = perm
	}
}

func resolveEnsureConfig(defaultPerm fs.FileMode, opts []EnsureOption) ensureConfig {
	cfg := defaultEnsureConfig(defaultPerm)
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&cfg)
	}
	if !isValidEnsureDirPerm(cfg.dirPerm) {
		cfg.dirPerm = defaultPerm
	}
	return cfg
}

func isValidEnsureDirPerm(perm fs.FileMode) bool {
	return perm != 0 && perm&fs.ModeType == 0
}
