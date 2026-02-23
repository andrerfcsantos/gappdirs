package gappdirs

func topLevelScopedContext(appName string, scope Scope) scopedContext {
	return newDefaultScopedContext(appName, scope)
}

func topLevelScopedDirs(appName string, scope Scope, cat category) []string {
	ctx := topLevelScopedContext(appName, scope)
	return scopedDirs(ctx, cat)
}

func topLevelScopedDir(appName string, scope Scope, cat category) string {
	ctx := topLevelScopedContext(appName, scope)
	return scopedDir(ctx, cat)
}

func topLevelScopedEnsureDir(appName string, scope Scope, cat category, opts ...EnsureOption) (string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return scopedEnsureDir(ctx, cat, opts...)
}

func topLevelScopedCreateFile(appName string, scope Scope, cat category, filename string, opts ...CreateFileOption) (bool, string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return scopedCreateFile(ctx, cat, filename, opts...)
}

func topLevelScopedFindFiles(appName string, scope Scope, cat category, filename string) ([]string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return findExistingScopedFiles(ctx, cat, filename)
}

func topLevelScopedFindFile(appName string, scope Scope, cat category, filename string) (string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return findExistingScopedFile(ctx, cat, filename)
}

// LocalDataDirs returns all candidate data directories paths, in precedence order for appName selected by local, user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureLocalDataDir to ensure the highest-precedence data directory exists and get its path.
func LocalDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryData)
}

// LocalConfigDirs returns all candidate config directories paths, in precedence order for appName selected by local, user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureLocalConfigDir to ensure the highest-precedence config directory exists and get its path.
func LocalConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryConfig)
}

// LocalLogDirs returns all candidate log directories paths, in precedence order for appName selected by local, user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureLocalLogDir to ensure the highest-precedence log directory exists and get its path.
func LocalLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryLog)
}

// LocalCacheDirs returns all candidate cache directories paths, in precedence order for appName selected by local, user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureLocalCacheDir to ensure the highest-precedence cache directory exists and get its path.
func LocalCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryCache)
}

// UserDataDirs returns all candidate data directories paths, in precedence order for appName selected by user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureUserDataDir to ensure the highest-precedence data directory exists and get its path.
func UserDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryData)
}

// UserConfigDirs returns all candidate config directories paths, in precedence order for appName selected by user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureUserConfigDir to ensure the highest-precedence config directory exists and get its path.
func UserConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryConfig)
}

// UserLogDirs returns all candidate log directories paths, in precedence order for appName selected by user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureUserLogDir to ensure the highest-precedence log directory exists and get its path.
func UserLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryLog)
}

// UserCacheDirs returns all candidate cache directories paths, in precedence order for appName selected by user, then system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureUserCacheDir to ensure the highest-precedence cache directory exists and get its path.
func UserCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryCache)
}

// SystemDataDirs returns all candidate data directories paths, in precedence order for appName selected by system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureSystemDataDir to ensure the highest-precedence data directory exists and get its path.
func SystemDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryData)
}

// SystemConfigDirs returns all candidate config directories paths, in precedence order for appName selected by system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureSystemConfigDir to ensure the highest-precedence config directory exists and get its path.
func SystemConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryConfig)
}

// SystemLogDirs returns all candidate log directories paths, in precedence order for appName selected by system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureSystemLogDir to ensure the highest-precedence log directory exists and get its path.
func SystemLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryLog)
}

// SystemCacheDirs returns all candidate cache directories paths, in precedence order for appName selected by system precedence.
//
// The paths returned are not guaranteed to exist on the filesystem.
// Use EnsureSystemCacheDir to ensure the highest-precedence cache directory exists and get its path.
func SystemCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryCache)
}

// LocalDataDir returns the single highest-precedence data directory path for appName selected by local, user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureLocalDataDir to ensure the highest-precedence data directory exists and get its path.
func LocalDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryData)
}

// LocalConfigDir returns the single highest-precedence config directory path for appName selected by local, user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureLocalConfigDir to ensure the highest-precedence config directory exists and get its path.
func LocalConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryConfig)
}

// LocalLogDir returns the single highest-precedence log directory path for appName selected by local, user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureLocalLogDir to ensure the highest-precedence log directory exists and get its path.
func LocalLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryLog)
}

// LocalCacheDir returns the single highest-precedence cache directory path for appName selected by local, user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureLocalCacheDir to ensure the highest-precedence cache directory exists and get its path.
func LocalCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryCache)
}

// UserDataDir returns the single highest-precedence data directory path for appName selected by user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureUserDataDir to ensure the highest-precedence data directory exists and get its path.
func UserDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryData)
}

// UserConfigDir returns the single highest-precedence config directory path for appName selected by user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureUserConfigDir to ensure the highest-precedence config directory exists and get its path.
func UserConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryConfig)
}

// UserLogDir returns the single highest-precedence log directory path for appName selected by user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureUserLogDir to ensure the highest-precedence log directory exists and get its path.
func UserLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryLog)
}

// UserCacheDir returns the single highest-precedence cache directory path for appName selected by user, then system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureUserCacheDir to ensure the highest-precedence cache directory exists and get its path.
func UserCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryCache)
}

// SystemDataDir returns the single highest-precedence data directory path for appName selected by system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureSystemDataDir to ensure the highest-precedence data directory exists and get its path.
func SystemDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryData)
}

// SystemConfigDir returns the single highest-precedence config directory path for appName selected by system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureSystemConfigDir to ensure the highest-precedence config directory exists and get its path.
func SystemConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryConfig)
}

// SystemLogDir returns the single highest-precedence log directory path for appName selected by system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureSystemLogDir to ensure the highest-precedence log directory exists and get its path.
func SystemLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryLog)
}

// SystemCacheDir returns the single highest-precedence cache directory path for appName selected by system precedence.
//
// The path returned is not guaranteed to exist on the filesystem.
// Use EnsureSystemCacheDir to ensure the highest-precedence cache directory exists and get its path.
func SystemCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryCache)
}

// EnsureLocalDataDir ensures the highest-precedence data directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureLocalDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryData, opts...)
}

// EnsureLocalConfigDir ensures the existence of the highest-precedence config directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureLocalConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryConfig, opts...)
}

// EnsureLocalLogDir ensures the existence of the highest-precedence log directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureLocalLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryLog, opts...)
}

// EnsureLocalCacheDir ensures the existence of the highest-precedence cache directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureLocalCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryCache, opts...)
}

// EnsureUserDataDir ensures the existence of the highest-precedence data directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureUserDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryData, opts...)
}

// EnsureUserConfigDir ensures the existence of the highest-precedence config directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureUserConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryConfig, opts...)
}

// EnsureUserLogDir ensures the existence of the highest-precedence log directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureUserLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryLog, opts...)
}

// EnsureUserCacheDir ensures the existence of the highest-precedence cache directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureUserCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryCache, opts...)
}

// EnsureSystemDataDir ensures the existence of the highest-precedence data directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureSystemDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryData, opts...)
}

// EnsureSystemConfigDir ensures the existence of the highest-precedence config directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureSystemConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryConfig, opts...)
}

// EnsureSystemLogDir ensures the existence of the highest-precedence log directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureSystemLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryLog, opts...)
}

// EnsureSystemCacheDir ensures the existence of the highest-precedence cache directory for appName and returns its path.
//
// It returns the absolute path to the directory, and any error encountered.
func EnsureSystemCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryCache, opts...)
}

// CreateLocalDataFile creates a new file in the highest-precedence data directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateLocalDataFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeLocal, categoryData, filename, opts...)
}

// CreateLocalConfigFile creates a new file in the highest-precedence config directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateLocalConfigFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeLocal, categoryConfig, filename, opts...)
}

// CreateLocalLogFile creates a new file in the highest-precedence log directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateLocalLogFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeLocal, categoryLog, filename, opts...)
}

// CreateLocalCacheFile creates a new file in the highest-precedence cache directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateLocalCacheFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeLocal, categoryCache, filename, opts...)
}

// CreateUserDataFile creates a new file in the highest-precedence data directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateUserDataFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeUser, categoryData, filename, opts...)
}

// CreateUserConfigFile creates a new file in the highest-precedence config directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateUserConfigFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeUser, categoryConfig, filename, opts...)
}

// CreateUserLogFile creates a new file in the highest-precedence log directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateUserLogFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeUser, categoryLog, filename, opts...)
}

// CreateUserCacheFile creates a new file in the highest-precedence cache directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateUserCacheFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeUser, categoryCache, filename, opts...)
}

// CreateSystemDataFile creates a new file in the highest-precedence data directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateSystemDataFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeSystem, categoryData, filename, opts...)
}

// CreateSystemConfigFile creates a new file in the highest-precedence config directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateSystemConfigFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeSystem, categoryConfig, filename, opts...)
}

// CreateSystemLogFile creates a new file in the highest-precedence log directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateSystemLogFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeSystem, categoryLog, filename, opts...)
}

// CreateSystemCacheFile creates a new file in the highest-precedence cache directory for appName.
//
// It returns the absolute path to the file, and any error encountered.
func CreateSystemCacheFile(appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	return topLevelScopedCreateFile(appName, ScopeSystem, categoryCache, filename, opts...)
}

// FindLocalDataFiles returns all existing matches for filename in data directories for appName selected by local, user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindLocalDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryData, filename)
}

// FindLocalConfigFiles returns all existing matches for filename in config directories for appName in local, user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindLocalConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryConfig, filename)
}

// FindLocalLogFiles returns all existing matches for filename in log directories for appName in local, user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindLocalLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryLog, filename)
}

// FindLocalCacheFiles returns all existing matches for filename in cache directories for appName in local, user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindLocalCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryCache, filename)
}

// FindUserDataFiles returns all existing matches for filename in data directories for appName in user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindUserDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryData, filename)
}

// FindUserConfigFiles returns all existing matches for filename in config directories for appName selected by user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindUserConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryConfig, filename)
}

// FindUserLogFiles returns all existing matches for filename in log directories for appName selected by user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindUserLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryLog, filename)
}

// FindUserCacheFiles returns all existing matches for filename in cache directories for appName selected by user, then system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindUserCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryCache, filename)
}

// FindSystemDataFiles returns all existing matches for filename in data directories for appName selected by system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindSystemDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryData, filename)
}

// FindSystemConfigFiles returns all existing matches for filename in config directories for appName selected by system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindSystemConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryConfig, filename)
}

// FindSystemLogFiles returns all existing matches for filename in log directories for appName selected by system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindSystemLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryLog, filename)
}

// FindSystemCacheFiles returns all existing matches for filename in cache directories for appName selected by system precedence order.
//
// It returns the paths to the existing files and any error encountered during the lookup.
// If no matching files are found, it returns an empty slice and nil error.
func FindSystemCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryCache, filename)
}

// FindLocalDataFile returns the first existing match for filename in data directories for appName selected by local, user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindLocalDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryData, filename)
}

// FindLocalConfigFile returns the first existing match for filename in config directories for appName selected by local, user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindLocalConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryConfig, filename)
}

// FindLocalLogFile returns the first existing match for filename in log directories for appName selected by local, user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindLocalLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryLog, filename)
}

// FindLocalCacheFile returns the first existing match for filename in cache directories for appName selected by local, user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindLocalCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryCache, filename)
}

// FindUserDataFile returns the first existing match for filename in data directories for appName selected by user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindUserDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryData, filename)
}

// FindUserConfigFile returns the first existing match for filename in config directories for appName selected by user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindUserConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryConfig, filename)
}

// FindUserLogFile returns the first existing match for filename in log directories for appName selected by user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindUserLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryLog, filename)
}

// FindUserCacheFile returns the first existing match for filename in cache directories for appName selected by user, then system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindUserCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryCache, filename)
}

// FindSystemDataFile returns the first existing match for filename in data directories for appName selected by system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindSystemDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryData, filename)
}

// FindSystemConfigFile returns the first existing match for filename in config directories for appName selected by system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindSystemConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryConfig, filename)
}

// FindSystemLogFile returns the first existing match for filename in log directories for appName selected by system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindSystemLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryLog, filename)
}

// FindSystemCacheFile returns the first existing match for filename in cache directories for appName selected by system precedence.
//
// It returns the path to the existing file and any error encountered during the lookup.
// If no matching file is found, it returns ErrNotFound.
func FindSystemCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryCache, filename)
}
