package gappdirs

import "io/fs"

func topLevelScopedContext(appName string, scope Scope) (scopedContext, error) {
	return newDefaultScopedContext(appName, scope)
}

func topLevelScopedDirs(appName string, scope Scope, cat category) ([]string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return nil, err
	}
	return scopedDirs(ctx, cat)
}

func topLevelScopedDir(appName string, scope Scope, cat category) (string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return "", err
	}
	return scopedDir(ctx, cat)
}

func topLevelScopedEnsureDir(appName string, scope Scope, cat category) (string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return "", err
	}
	return scopedEnsureDir(ctx, cat)
}

func topLevelScopedEnsureDirWithPerm(appName string, scope Scope, cat category, perm fs.FileMode) (string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return "", err
	}
	return scopedEnsureDirWithPerm(ctx, cat, perm)
}

func topLevelScopedFindFileDirs(appName string, scope Scope, cat category, filename string) ([]string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return nil, err
	}
	return scopedFindFileDirs(ctx, cat, filename)
}

func topLevelScopedFile(appName string, scope Scope, cat category, filename string) (string, error) {
	ctx, err := topLevelScopedContext(appName, scope)
	if err != nil {
		return "", err
	}
	return scopedFile(ctx, cat, filename)
}

func LocalDataDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeLocal, categoryData)
}

func LocalConfigDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeLocal, categoryConfig)
}

func LocalLogDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeLocal, categoryLog)
}

func LocalCacheDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeLocal, categoryCache)
}

func UserDataDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeUser, categoryData)
}

func UserConfigDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeUser, categoryConfig)
}

func UserLogDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeUser, categoryLog)
}

func UserCacheDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeUser, categoryCache)
}

func SystemDataDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeSystem, categoryData)
}

func SystemConfigDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeSystem, categoryConfig)
}

func SystemLogDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeSystem, categoryLog)
}

func SystemCacheDirs(appName string) ([]string, error) {
	return topLevelScopedDirs(appName, ScopeSystem, categoryCache)
}

func LocalDataDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeLocal, categoryData)
}

func LocalConfigDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeLocal, categoryConfig)
}

func LocalLogDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeLocal, categoryLog)
}

func LocalCacheDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeLocal, categoryCache)
}

func UserDataDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeUser, categoryData)
}

func UserConfigDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeUser, categoryConfig)
}

func UserLogDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeUser, categoryLog)
}

func UserCacheDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeUser, categoryCache)
}

func SystemDataDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeSystem, categoryData)
}

func SystemConfigDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeSystem, categoryConfig)
}

func SystemLogDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeSystem, categoryLog)
}

func SystemCacheDir(appName string) (string, error) {
	return topLevelScopedDir(appName, ScopeSystem, categoryCache)
}

func EnsureLocalDataDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryData)
}

func EnsureLocalConfigDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryConfig)
}

func EnsureLocalLogDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryLog)
}

func EnsureLocalCacheDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeLocal, categoryCache)
}

func EnsureUserDataDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryData)
}

func EnsureUserConfigDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryConfig)
}

func EnsureUserLogDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryLog)
}

func EnsureUserCacheDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeUser, categoryCache)
}

func EnsureSystemDataDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryData)
}

func EnsureSystemConfigDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryConfig)
}

func EnsureSystemLogDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryLog)
}

func EnsureSystemCacheDir(appName string) (string, error) {
	return topLevelScopedEnsureDir(appName, ScopeSystem, categoryCache)
}

func EnsureLocalDataDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeLocal, categoryData, perm)
}

func EnsureLocalConfigDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeLocal, categoryConfig, perm)
}

func EnsureLocalLogDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeLocal, categoryLog, perm)
}

func EnsureLocalCacheDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeLocal, categoryCache, perm)
}

func EnsureUserDataDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeUser, categoryData, perm)
}

func EnsureUserConfigDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeUser, categoryConfig, perm)
}

func EnsureUserLogDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeUser, categoryLog, perm)
}

func EnsureUserCacheDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeUser, categoryCache, perm)
}

func EnsureSystemDataDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeSystem, categoryData, perm)
}

func EnsureSystemConfigDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeSystem, categoryConfig, perm)
}

func EnsureSystemLogDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeSystem, categoryLog, perm)
}

func EnsureSystemCacheDirWithPerm(appName string, perm fs.FileMode) (string, error) {
	return topLevelScopedEnsureDirWithPerm(appName, ScopeSystem, categoryCache, perm)
}

func FindLocalDataFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeLocal, categoryData, filename)
}

func FindLocalConfigFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeLocal, categoryConfig, filename)
}

func FindLocalLogFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeLocal, categoryLog, filename)
}

func FindLocalCacheFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeLocal, categoryCache, filename)
}

func FindUserDataFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeUser, categoryData, filename)
}

func FindUserConfigFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeUser, categoryConfig, filename)
}

func FindUserLogFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeUser, categoryLog, filename)
}

func FindUserCacheFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeUser, categoryCache, filename)
}

func FindSystemDataFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeSystem, categoryData, filename)
}

func FindSystemConfigFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeSystem, categoryConfig, filename)
}

func FindSystemLogFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeSystem, categoryLog, filename)
}

func FindSystemCacheFileDirs(appName string, filename string) ([]string, error) {
	return topLevelScopedFindFileDirs(appName, ScopeSystem, categoryCache, filename)
}

func LocalDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeLocal, categoryData, filename)
}

func LocalConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeLocal, categoryConfig, filename)
}

func LocalLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeLocal, categoryLog, filename)
}

func LocalCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeLocal, categoryCache, filename)
}

func UserDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeUser, categoryData, filename)
}

func UserConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeUser, categoryConfig, filename)
}

func UserLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeUser, categoryLog, filename)
}

func UserCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeUser, categoryCache, filename)
}

func SystemDataFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeSystem, categoryData, filename)
}

func SystemConfigFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeSystem, categoryConfig, filename)
}

func SystemLogFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeSystem, categoryLog, filename)
}

func SystemCacheFile(appName string, filename string) (string, error) {
	return topLevelScopedFile(appName, ScopeSystem, categoryCache, filename)
}
