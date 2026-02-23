package gappdirs

// CreateDataFile creates a file with the given name in the highest-precedence data directory for the Resolver scope.
//
// If the file already exists, it will not be modified, unless the WithOverwriteExisting option is used.
// If the file does not exist, it will be created with the specified permissions and contents (if provided).
//
// It returns a boolean indicating whether a new file was created and didn't exist before, the absolute path to the file, and any error encountered.
func (r *Resolver) CreateDataFile(filename string, opts ...CreateFileOption) (created bool, path string, err error) {
	return scopedCreateFile(r.ctx, categoryData, filename, opts...)
}

// CreateConfigFile creates a file with the given name in the highest-precedence config directory for the Resolver scope.
//
// If the file already exists, it will not be modified, unless the WithOverwriteExisting option is used.
// If the file does not exist, it will be created with the specified permissions and contents (if provided).
//
// It returns a boolean indicating whether a new file was created and didn't exist before, the absolute path to the file, and any error encountered.
func (r *Resolver) CreateConfigFile(filename string, opts ...CreateFileOption) (created bool, path string, err error) {
	return scopedCreateFile(r.ctx, categoryConfig, filename, opts...)
}

// CreateLogFile creates a file with the given name in the highest-precedence log directory for the Resolver scope.
//
// If the file already exists, it will not be modified, unless the WithOverwriteExisting option is used.
// If the file does not exist, it will be created with the specified permissions and contents (if provided).
//
// It returns a boolean indicating whether a new file was created and didn't exist before, the absolute path to the file, and any error encountered.
func (r *Resolver) CreateLogFile(filename string, opts ...CreateFileOption) (created bool, path string, err error) {
	return scopedCreateFile(r.ctx, categoryLog, filename, opts...)
}

// CreateCacheFile creates a file with the given name in the highest-precedence cache directory for the Resolver scope.
//
// If the file already exists, it will not be modified, unless the WithOverwriteExisting option is used.
// If the file does not exist, it will be created with the specified permissions and contents (if provided).
//
// It returns a boolean indicating whether a new file was created and didn't exist before, the absolute path to the file, and any error encountered.
func (r *Resolver) CreateCacheFile(filename string, opts ...CreateFileOption) (created bool, path string, err error) {
	return scopedCreateFile(r.ctx, categoryCache, filename, opts...)
}
