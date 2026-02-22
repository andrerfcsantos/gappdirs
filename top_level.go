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

func topLevelScopedFindFiles(appName string, scope Scope, cat category, filename string) ([]string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return findExistingScopedFiles(ctx, cat, filename)
}

func topLevelScopedFindFile(appName string, scope Scope, cat category, filename string) (string, error) {
	ctx := topLevelScopedContext(appName, scope)
	return findExistingScopedFile(ctx, cat, filename)
}

func LocalDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryData)
}

func LocalConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryConfig)
}

func LocalLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryLog)
}

func LocalCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeLocal, categoryCache)
}

func UserDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryData)
}

func UserConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryConfig)
}

func UserLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryLog)
}

func UserCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeUser, categoryCache)
}

func SystemDataDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryData)
}

func SystemConfigDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryConfig)
}

func SystemLogDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryLog)
}

func SystemCacheDirs(appName string) []string {
	return topLevelScopedDirs(appName, ScopeSystem, categoryCache)
}

func LocalDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryData)
}

func LocalConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryConfig)
}

func LocalLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryLog)
}

func LocalCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeLocal, categoryCache)
}

func UserDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryData)
}

func UserConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryConfig)
}

func UserLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryLog)
}

func UserCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeUser, categoryCache)
}

func SystemDataDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryData)
}

func SystemConfigDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryConfig)
}

func SystemLogDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryLog)
}

func SystemCacheDir(appName string) string {
	return topLevelScopedDir(appName, ScopeSystem, categoryCache)
}

func EnsureLocalDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryData, opts...)
}

func EnsureLocalConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryConfig, opts...)
}

func EnsureLocalLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryLog, opts...)
}

func EnsureLocalCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryCache, opts...)
}

func EnsureUserDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryData, opts...)
}

func EnsureUserConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryConfig, opts...)
}

func EnsureUserLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryLog, opts...)
}

func EnsureUserCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryCache, opts...)
}

func EnsureSystemDataDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryData, opts...)
}

func EnsureSystemConfigDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryConfig, opts...)
}

func EnsureSystemLogDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryLog, opts...)
}

func EnsureSystemCacheDir(appName string, opts ...EnsureOption) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryCache, opts...)
}

func FindLocalDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryData, filename)
}

func FindLocalConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryConfig, filename)
}

func FindLocalLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryLog, filename)
}

func FindLocalCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeLocal, categoryCache, filename)
}

func FindUserDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryData, filename)
}

func FindUserConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryConfig, filename)
}

func FindUserLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryLog, filename)
}

func FindUserCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeUser, categoryCache, filename)
}

func FindSystemDataFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryData, filename)
}

func FindSystemConfigFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryConfig, filename)
}

func FindSystemLogFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryLog, filename)
}

func FindSystemCacheFiles(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFiles(appName, ScopeSystem, categoryCache, filename)
}

func FindLocalDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryData, filename)
}

func FindLocalConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryConfig, filename)
}

func FindLocalLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryLog, filename)
}

func FindLocalCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeLocal, categoryCache, filename)
}

func FindUserDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryData, filename)
}

func FindUserConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryConfig, filename)
}

func FindUserLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryLog, filename)
}

func FindUserCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeUser, categoryCache, filename)
}

func FindSystemDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryData, filename)
}

func FindSystemConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryConfig, filename)
}

func FindSystemLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryLog, filename)
}

func FindSystemCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFindFile(appName, ScopeSystem, categoryCache, filename)
}
