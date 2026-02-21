package gappdirs

import "io/fs"

// EnsureDataDir creates the most relevant data directory with resolver default permissions.
func (r *resolver) EnsureDataDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDir(ctx, categoryData)
}

// EnsureDataDirWithPerm creates the most relevant data directory with explicit permissions.
func (r *resolver) EnsureDataDirWithPerm(perm fs.FileMode) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDirWithPerm(ctx, categoryData, perm)
}

// EnsureConfigDir creates the most relevant config directory with resolver default permissions.
func (r *resolver) EnsureConfigDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDir(ctx, categoryConfig)
}

// EnsureConfigDirWithPerm creates the most relevant config directory with explicit permissions.
func (r *resolver) EnsureConfigDirWithPerm(perm fs.FileMode) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDirWithPerm(ctx, categoryConfig, perm)
}

// EnsureLogDir creates the most relevant log directory with resolver default permissions.
func (r *resolver) EnsureLogDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDir(ctx, categoryLog)
}

// EnsureLogDirWithPerm creates the most relevant log directory with explicit permissions.
func (r *resolver) EnsureLogDirWithPerm(perm fs.FileMode) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDirWithPerm(ctx, categoryLog, perm)
}

// EnsureCacheDir creates the most relevant cache directory with resolver default permissions.
func (r *resolver) EnsureCacheDir() (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDir(ctx, categoryCache)
}

// EnsureCacheDirWithPerm creates the most relevant cache directory with explicit permissions.
func (r *resolver) EnsureCacheDirWithPerm(perm fs.FileMode) (string, error) {
	ctx, err := r.scopedCtx()
	if err != nil {
		return "", err
	}
	return scopedEnsureDirWithPerm(ctx, categoryCache, perm)
}
