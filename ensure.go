package gappdirs

// EnsureDataDir creates the most relevant data directory using resolver defaults or per-call options.
func (r *resolver) EnsureDataDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryData, opts...)
}

// EnsureConfigDir creates the most relevant config directory using resolver defaults or per-call options.
func (r *resolver) EnsureConfigDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryConfig, opts...)
}

// EnsureLogDir creates the most relevant log directory using resolver defaults or per-call options.
func (r *resolver) EnsureLogDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryLog, opts...)
}

// EnsureCacheDir creates the most relevant cache directory using resolver defaults or per-call options.
func (r *resolver) EnsureCacheDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryCache, opts...)
}
