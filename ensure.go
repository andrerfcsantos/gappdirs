package gappdirs

// EnsureDataDir ensures the existence of the highest-precedence data directory
// for the Resolver scope by creating it if it doesn't exist.
//
// It returns the absolute path to the directory, and any error encountered.
func (r *Resolver) EnsureDataDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryData, opts...)
}

// EnsureConfigDir ensures the existence of the highest-precedence config directory
// for the Resolver scope by creating it if it doesn't exist.
//
// It returns the absolute path to the directory, and any error encountered.
func (r *Resolver) EnsureConfigDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryConfig, opts...)
}

// EnsureLogDir ensures the existence of the highest-precedence log directory
// for the Resolver scope by creating it if it doesn't exist.
//
// It returns the absolute path to the directory, and any error encountered.
func (r *Resolver) EnsureLogDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryLog, opts...)
}

// EnsureCacheDir ensures the existence of the highest-precedence cache directory
// for the Resolver scope by creating it if it doesn't exist.
//
// It returns the absolute path to the directory, and any error encountered.
func (r *Resolver) EnsureCacheDir(opts ...EnsureOption) (string, error) {
	return scopedEnsureDir(r.ctx, categoryCache, opts...)
}
