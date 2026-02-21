package gappdirs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var testScopes = []Scope{ScopeLocal, ScopeUser, ScopeSystem}

var testCategories = []category{
	categoryData,
	categoryConfig,
	categoryLog,
	categoryCache,
}

func TestTopLevelScopedDirsAndDirParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_parity"

	for _, scope := range testScopes {
		r, err := newResolverForScope(scope, appName)
		if err != nil {
			t.Fatalf("new resolver for scope %s: %v", scope, err)
		}

		for _, cat := range testCategories {
			t.Run(fmt.Sprintf("dirs_%s_%s", scope, cat), func(t *testing.T) {
				topDirs, topErr := callTopDirs(scope, cat, appName)
				resolverDirs, resolverErr := callResolverDirs(r, cat)
				assertParity(t, topDirs, topErr, resolverDirs, resolverErr)
			})

			t.Run(fmt.Sprintf("dir_%s_%s", scope, cat), func(t *testing.T) {
				topDir, topErr := callTopDir(scope, cat, appName)
				resolverDir, resolverErr := callResolverDir(r, cat)
				assertParity(t, topDir, topErr, resolverDir, resolverErr)
			})
		}
	}
}

func TestTopLevelScopedEnsureParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_ensure"

	for _, scope := range testScopes {
		r, err := newResolverForScope(scope, appName)
		if err != nil {
			t.Fatalf("new resolver for scope %s: %v", scope, err)
		}

		for _, cat := range testCategories {
			t.Run(fmt.Sprintf("ensure_%s_%s", scope, cat), func(t *testing.T) {
				topPath, topErr := callTopEnsure(scope, cat, appName)
				resolverPath, resolverErr := callResolverEnsure(r, cat)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)
				if topErr == nil {
					assertDirExists(t, topPath)
				}
			})

			t.Run(fmt.Sprintf("ensure_with_perm_%s_%s", scope, cat), func(t *testing.T) {
				topPath, topErr := callTopEnsureWithPerm(scope, cat, appName, 0o755)
				resolverPath, resolverErr := callResolverEnsureWithPerm(r, cat, 0o755)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)
				if topErr == nil {
					assertDirExists(t, topPath)
				}
			})
		}
	}
}

func TestTopLevelScopedFindAndFileParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_search"

	for _, scope := range testScopes {
		r, err := newResolverForScope(scope, appName)
		if err != nil {
			t.Fatalf("new resolver for scope %s: %v", scope, err)
		}

		for _, cat := range testCategories {
			filename := fmt.Sprintf("file_%s_%s.txt", scope, cat)
			dirs, err := callResolverDirs(r, cat)
			if err != nil {
				t.Fatalf("resolve dirs for %s/%s: %v", scope, cat, err)
			}
			seededDir, seeded := seedFileInFirstWritableDir(dirs, filename)

			t.Run(fmt.Sprintf("find_%s_%s", scope, cat), func(t *testing.T) {
				topDirs, topErr := callTopFindDirs(scope, cat, appName, filename)
				resolverDirs, resolverErr := callResolverFindDirs(r, cat, filename)
				assertParity(t, topDirs, topErr, resolverDirs, resolverErr)

				if seeded {
					if len(topDirs) == 0 || topDirs[0] != seededDir {
						t.Fatalf("expected first matching dir %q, got %#v", seededDir, topDirs)
					}
				}
			})

			t.Run(fmt.Sprintf("file_%s_%s", scope, cat), func(t *testing.T) {
				topFile, topErr := callTopFile(scope, cat, appName, filename)
				resolverFile, resolverErr := callResolverFile(r, cat, filename)
				assertParity(t, topFile, topErr, resolverFile, resolverErr)

				if seeded {
					want := filepath.Join(seededDir, filename)
					if topFile != want {
						t.Fatalf("expected first matching file %q, got %q", want, topFile)
					}
				} else {
					if !errors.Is(topErr, ErrNotFound) {
						t.Fatalf("expected ErrNotFound when no seeded file exists, got %v", topErr)
					}
				}
			})
		}
	}
}

func TestTopLevelErrorParity(t *testing.T) {
	setupTopLevelTestEnv(t)

	for _, scope := range testScopes {
		t.Run(fmt.Sprintf("invalid_app_%s", scope), func(t *testing.T) {
			_, topErr := callTopDir(scope, categoryData, "")
			_, resolverErr := newResolverForScope(scope, "")
			assertErrorParity(t, topErr, resolverErr)
		})
	}

	appName := "demo_top_level_errors"
	for _, scope := range testScopes {
		r, err := newResolverForScope(scope, appName)
		if err != nil {
			t.Fatalf("new resolver for scope %s: %v", scope, err)
		}

		t.Run(fmt.Sprintf("invalid_filename_%s", scope), func(t *testing.T) {
			_, topErr := callTopFindDirs(scope, categoryData, appName, "a/b")
			_, resolverErr := callResolverFindDirs(r, categoryData, "a/b")
			assertErrorParity(t, topErr, resolverErr)
		})

		t.Run(fmt.Sprintf("missing_file_%s", scope), func(t *testing.T) {
			filename := fmt.Sprintf("missing_%s.txt", scope)
			_, topErr := callTopFile(scope, categoryData, appName, filename)
			_, resolverErr := callResolverFile(r, categoryData, filename)
			assertErrorParity(t, topErr, resolverErr)
			if !errors.Is(topErr, ErrNotFound) {
				t.Fatalf("expected ErrNotFound, got %v", topErr)
			}
		})
	}
}

func setupTopLevelTestEnv(t *testing.T) {
	t.Helper()

	wd := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(wd); err != nil {
		t.Fatalf("chdir %q: %v", wd, err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	homeDir := filepath.Join(wd, "home")
	if err := os.MkdirAll(homeDir, 0o755); err != nil {
		t.Fatalf("mkdir home dir %q: %v", homeDir, err)
	}
	t.Setenv("HOME", homeDir)

	switch runtime.GOOS {
	case "linux":
		t.Setenv("XDG_DATA_HOME", filepath.Join(homeDir, "xdg-data-home"))
		t.Setenv("XDG_CONFIG_HOME", filepath.Join(homeDir, "xdg-config-home"))
		t.Setenv("XDG_STATE_HOME", filepath.Join(homeDir, "xdg-state-home"))
		t.Setenv("XDG_CACHE_HOME", filepath.Join(homeDir, "xdg-cache-home"))
		xdgDataSystem := []string{
			filepath.Join(wd, "xdg-system-data-a"),
			filepath.Join(wd, "xdg-system-data-b"),
		}
		t.Setenv("XDG_DATA_DIRS", strings.Join(xdgDataSystem, string(os.PathListSeparator)))
		t.Setenv("XDG_CONFIG_DIRS", filepath.Join(wd, "xdg-system-config"))
	case "windows":
		local := filepath.Join(wd, "LocalAppData")
		roaming := filepath.Join(wd, "RoamingAppData")
		programData := filepath.Join(wd, "ProgramData")
		t.Setenv("LOCALAPPDATA", local)
		t.Setenv("APPDATA", roaming)
		t.Setenv("ProgramData", programData)
		t.Setenv("PROGRAMDATA", programData)
	}
}

func newResolverForScope(scope Scope, appName string) (Resolver, error) {
	switch scope {
	case ScopeLocal:
		return NewLocalResolver(appName)
	case ScopeUser:
		return NewUserResolver(appName)
	case ScopeSystem:
		return NewSystemResolver(appName)
	default:
		return nil, fmt.Errorf("unsupported scope %s", scope)
	}
}

func assertParity[T any](t *testing.T, topValue T, topErr error, resolverValue T, resolverErr error) {
	t.Helper()
	assertErrorParity(t, topErr, resolverErr)
	if topErr == nil && !reflect.DeepEqual(topValue, resolverValue) {
		t.Fatalf("result mismatch:\nresolver: %#v\ntop:      %#v", resolverValue, topValue)
	}
}

func assertErrorParity(t *testing.T, topErr, resolverErr error) {
	t.Helper()

	switch {
	case topErr == nil && resolverErr == nil:
		return
	case topErr == nil || resolverErr == nil:
		t.Fatalf("error mismatch:\nresolver: %v\ntop:      %v", resolverErr, topErr)
	case topErr.Error() != resolverErr.Error():
		t.Fatalf("error mismatch:\nresolver: %v\ntop:      %v", resolverErr, topErr)
	}
}

func assertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("expected directory %q to exist: %v", path, err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", path)
	}
}

func seedFileInFirstWritableDir(dirs []string, filename string) (string, bool) {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			continue
		}

		candidate := filepath.Join(dir, filename)
		if err := os.WriteFile(candidate, []byte("seed"), 0o644); err != nil {
			continue
		}
		return dir, true
	}

	return "", false
}

func callTopDirs(scope Scope, cat category, appName string) ([]string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return LocalDataDirs(appName)
		case categoryConfig:
			return LocalConfigDirs(appName)
		case categoryLog:
			return LocalLogDirs(appName)
		case categoryCache:
			return LocalCacheDirs(appName)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return UserDataDirs(appName)
		case categoryConfig:
			return UserConfigDirs(appName)
		case categoryLog:
			return UserLogDirs(appName)
		case categoryCache:
			return UserCacheDirs(appName)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return SystemDataDirs(appName)
		case categoryConfig:
			return SystemConfigDirs(appName)
		case categoryLog:
			return SystemLogDirs(appName)
		case categoryCache:
			return SystemCacheDirs(appName)
		}
	}
	return nil, fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopDir(scope Scope, cat category, appName string) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return LocalDataDir(appName)
		case categoryConfig:
			return LocalConfigDir(appName)
		case categoryLog:
			return LocalLogDir(appName)
		case categoryCache:
			return LocalCacheDir(appName)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return UserDataDir(appName)
		case categoryConfig:
			return UserConfigDir(appName)
		case categoryLog:
			return UserLogDir(appName)
		case categoryCache:
			return UserCacheDir(appName)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return SystemDataDir(appName)
		case categoryConfig:
			return SystemConfigDir(appName)
		case categoryLog:
			return SystemLogDir(appName)
		case categoryCache:
			return SystemCacheDir(appName)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopEnsure(scope Scope, cat category, appName string) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return EnsureLocalDataDir(appName)
		case categoryConfig:
			return EnsureLocalConfigDir(appName)
		case categoryLog:
			return EnsureLocalLogDir(appName)
		case categoryCache:
			return EnsureLocalCacheDir(appName)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return EnsureUserDataDir(appName)
		case categoryConfig:
			return EnsureUserConfigDir(appName)
		case categoryLog:
			return EnsureUserLogDir(appName)
		case categoryCache:
			return EnsureUserCacheDir(appName)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return EnsureSystemDataDir(appName)
		case categoryConfig:
			return EnsureSystemConfigDir(appName)
		case categoryLog:
			return EnsureSystemLogDir(appName)
		case categoryCache:
			return EnsureSystemCacheDir(appName)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopEnsureWithPerm(scope Scope, cat category, appName string, perm fs.FileMode) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return EnsureLocalDataDirWithPerm(appName, perm)
		case categoryConfig:
			return EnsureLocalConfigDirWithPerm(appName, perm)
		case categoryLog:
			return EnsureLocalLogDirWithPerm(appName, perm)
		case categoryCache:
			return EnsureLocalCacheDirWithPerm(appName, perm)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return EnsureUserDataDirWithPerm(appName, perm)
		case categoryConfig:
			return EnsureUserConfigDirWithPerm(appName, perm)
		case categoryLog:
			return EnsureUserLogDirWithPerm(appName, perm)
		case categoryCache:
			return EnsureUserCacheDirWithPerm(appName, perm)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return EnsureSystemDataDirWithPerm(appName, perm)
		case categoryConfig:
			return EnsureSystemConfigDirWithPerm(appName, perm)
		case categoryLog:
			return EnsureSystemLogDirWithPerm(appName, perm)
		case categoryCache:
			return EnsureSystemCacheDirWithPerm(appName, perm)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopFindDirs(scope Scope, cat category, appName string, filename string) ([]string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return FindLocalDataFileDirs(appName, filename)
		case categoryConfig:
			return FindLocalConfigFileDirs(appName, filename)
		case categoryLog:
			return FindLocalLogFileDirs(appName, filename)
		case categoryCache:
			return FindLocalCacheFileDirs(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return FindUserDataFileDirs(appName, filename)
		case categoryConfig:
			return FindUserConfigFileDirs(appName, filename)
		case categoryLog:
			return FindUserLogFileDirs(appName, filename)
		case categoryCache:
			return FindUserCacheFileDirs(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return FindSystemDataFileDirs(appName, filename)
		case categoryConfig:
			return FindSystemConfigFileDirs(appName, filename)
		case categoryLog:
			return FindSystemLogFileDirs(appName, filename)
		case categoryCache:
			return FindSystemCacheFileDirs(appName, filename)
		}
	}
	return nil, fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopFile(scope Scope, cat category, appName string, filename string) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return LocalDataFile(appName, filename)
		case categoryConfig:
			return LocalConfigFile(appName, filename)
		case categoryLog:
			return LocalLogFile(appName, filename)
		case categoryCache:
			return LocalCacheFile(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return UserDataFile(appName, filename)
		case categoryConfig:
			return UserConfigFile(appName, filename)
		case categoryLog:
			return UserLogFile(appName, filename)
		case categoryCache:
			return UserCacheFile(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return SystemDataFile(appName, filename)
		case categoryConfig:
			return SystemConfigFile(appName, filename)
		case categoryLog:
			return SystemLogFile(appName, filename)
		case categoryCache:
			return SystemCacheFile(appName, filename)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callResolverDirs(r Resolver, cat category) ([]string, error) {
	switch cat {
	case categoryData:
		return r.DataDirs()
	case categoryConfig:
		return r.ConfigDirs()
	case categoryLog:
		return r.LogDirs()
	case categoryCache:
		return r.CacheDirs()
	default:
		return nil, fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverDir(r Resolver, cat category) (string, error) {
	switch cat {
	case categoryData:
		return r.DataDir()
	case categoryConfig:
		return r.ConfigDir()
	case categoryLog:
		return r.LogDir()
	case categoryCache:
		return r.CacheDir()
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverEnsure(r Resolver, cat category) (string, error) {
	switch cat {
	case categoryData:
		return r.EnsureDataDir()
	case categoryConfig:
		return r.EnsureConfigDir()
	case categoryLog:
		return r.EnsureLogDir()
	case categoryCache:
		return r.EnsureCacheDir()
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverEnsureWithPerm(r Resolver, cat category, perm fs.FileMode) (string, error) {
	switch cat {
	case categoryData:
		return r.EnsureDataDirWithPerm(perm)
	case categoryConfig:
		return r.EnsureConfigDirWithPerm(perm)
	case categoryLog:
		return r.EnsureLogDirWithPerm(perm)
	case categoryCache:
		return r.EnsureCacheDirWithPerm(perm)
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverFindDirs(r Resolver, cat category, filename string) ([]string, error) {
	switch cat {
	case categoryData:
		return r.FindDataFileDirs(filename)
	case categoryConfig:
		return r.FindConfigFileDirs(filename)
	case categoryLog:
		return r.FindLogFileDirs(filename)
	case categoryCache:
		return r.FindCacheFileDirs(filename)
	default:
		return nil, fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverFile(r Resolver, cat category, filename string) (string, error) {
	switch cat {
	case categoryData:
		return r.DataFile(filename)
	case categoryConfig:
		return r.ConfigFile(filename)
	case categoryLog:
		return r.LogFile(filename)
	case categoryCache:
		return r.CacheFile(filename)
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}
