package gappdirs

import (
	"errors"
	"fmt"
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
		r := newResolverForScope(scope, appName)

		for _, cat := range testCategories {
			t.Run(fmt.Sprintf("dirs_%s_%s", scope, cat), func(t *testing.T) {
				topDirs := callTopDirs(scope, cat, appName)
				resolverDirs := callResolverDirs(r, cat)
				assertValueParity(t, topDirs, resolverDirs)
			})

			t.Run(fmt.Sprintf("dir_%s_%s", scope, cat), func(t *testing.T) {
				topDir := callTopDir(scope, cat, appName)
				resolverDir := callResolverDir(r, cat)
				assertValueParity(t, topDir, resolverDir)
			})
		}
	}
}

func TestTopLevelScopedFilePathParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_file_paths"
	filenames := []string{
		"settings.yaml",
		"   ",
		filepath.Join(string(filepath.Separator), "absolute", "rooted.json"),
	}

	for _, scope := range testScopes {
		r := newResolverForScope(scope, appName)

		for _, cat := range testCategories {
			for _, filename := range filenames {
				t.Run(fmt.Sprintf("file_paths_%s_%s_%q", scope, cat, filename), func(t *testing.T) {
					topPaths := callTopFilePaths(scope, cat, appName, filename)
					resolverPaths := callResolverFilePaths(r, cat, filename)
					assertValueParity(t, topPaths, resolverPaths)
				})

				t.Run(fmt.Sprintf("file_path_%s_%s_%q", scope, cat, filename), func(t *testing.T) {
					topPath := callTopFilePath(scope, cat, appName, filename)
					resolverPath := callResolverFilePath(r, cat, filename)
					assertValueParity(t, topPath, resolverPath)

					topPaths := callTopFilePaths(scope, cat, appName, filename)
					if len(topPaths) == 0 {
						if topPath != "" {
							t.Fatalf("expected empty path when no candidate dirs are found, got %q", topPath)
						}
						return
					}
					if topPath != topPaths[0] {
						t.Fatalf("expected top path %q to match first candidate %q", topPath, topPaths[0])
					}
				})
			}
		}
	}
}

func TestTopLevelScopedEnsureParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_ensure"

	for _, scope := range testScopes {
		r := newResolverForScope(scope, appName)

		for _, cat := range testCategories {
			t.Run(fmt.Sprintf("ensure_%s_%s", scope, cat), func(t *testing.T) {
				topPath, topErr := callTopEnsure(scope, cat, appName)
				resolverPath, resolverErr := callResolverEnsure(r, cat)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)
				if topErr == nil {
					assertDirExists(t, topPath)
				}
			})

			t.Run(fmt.Sprintf("ensure_with_options_%s_%s", scope, cat), func(t *testing.T) {
				topPath, topErr := callTopEnsureWithOpts(scope, cat, appName, WithEnsureDirPerm(0o755))
				resolverPath, resolverErr := callResolverEnsureWithOpts(r, cat, WithEnsureDirPerm(0o755))
				assertParity(t, topPath, topErr, resolverPath, resolverErr)
				if topErr == nil {
					assertDirExists(t, topPath)
				}
			})
		}
	}
}

func TestTopLevelScopedCreateFileParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_create"

	for _, scope := range testScopes {
		r := newResolverForScope(scope, appName)

		for _, cat := range testCategories {
			filename := fmt.Sprintf("created_%s_%s.txt", scope, cat)
			targetPath := filepath.Join(callTopDir(scope, cat, appName), filename)

			t.Run(fmt.Sprintf("create_new_%s_%s", scope, cat), func(t *testing.T) {
				_ = os.Remove(targetPath)

				topCreated, topPath, topErr := callTopCreateFile(scope, cat, appName, filename, WithContentsFromReader(strings.NewReader("seed")))
				if topErr == nil {
					_ = os.Remove(topPath)
				}

				resolverCreated, resolverPath, resolverErr := callResolverCreateFile(r, cat, filename, WithContentsFromReader(strings.NewReader("seed")))
				assertParity(t, topCreated, topErr, resolverCreated, resolverErr)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)
			})

			t.Run(fmt.Sprintf("existing_no_overwrite_%s_%s", scope, cat), func(t *testing.T) {
				seeded := seedExactFilePath(targetPath, "existing")

				topCreated, topPath, topErr := callTopCreateFile(scope, cat, appName, filename, WithContentsFromReader(strings.NewReader("new-content")))
				resolverCreated, resolverPath, resolverErr := callResolverCreateFile(r, cat, filename, WithContentsFromReader(strings.NewReader("new-content")))
				assertParity(t, topCreated, topErr, resolverCreated, resolverErr)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)

				if seeded && topErr == nil {
					content, err := os.ReadFile(targetPath)
					if err != nil {
						t.Fatalf("read seeded file: %v", err)
					}
					if string(content) != "existing" {
						t.Fatalf("expected existing content to remain untouched, got %q", string(content))
					}
				}
			})

			t.Run(fmt.Sprintf("existing_overwrite_%s_%s", scope, cat), func(t *testing.T) {
				seeded := seedExactFilePath(targetPath, "existing")

				topCreated, topPath, topErr := callTopCreateFile(
					scope,
					cat,
					appName,
					filename,
					WithOverwriteExisting(),
					WithContentsFromReader(strings.NewReader("new-content")),
				)

				if seeded {
					_ = seedExactFilePath(targetPath, "existing")
				}

				resolverCreated, resolverPath, resolverErr := callResolverCreateFile(
					r,
					cat,
					filename,
					WithOverwriteExisting(),
					WithContentsFromReader(strings.NewReader("new-content")),
				)

				assertParity(t, topCreated, topErr, resolverCreated, resolverErr)
				assertParity(t, topPath, topErr, resolverPath, resolverErr)

				if seeded && resolverErr == nil {
					content, err := os.ReadFile(targetPath)
					if err != nil {
						t.Fatalf("read overwritten file: %v", err)
					}
					if string(content) != "new-content" {
						t.Fatalf("expected overwritten content %q, got %q", "new-content", string(content))
					}
				}
			})
		}
	}
}

func TestTopLevelScopedFindAndFileParity(t *testing.T) {
	setupTopLevelTestEnv(t)
	appName := "demo_top_level_search"

	for _, scope := range testScopes {
		r := newResolverForScope(scope, appName)

		for _, cat := range testCategories {
			filename := fmt.Sprintf("file_%s_%s.txt", scope, cat)
			dirs := callResolverDirs(r, cat)
			seededDir, seeded := seedFileInFirstWritableDir(dirs, filename)

			t.Run(fmt.Sprintf("find_%s_%s", scope, cat), func(t *testing.T) {
				topFiles, topErr := callTopFindFiles(scope, cat, appName, filename)
				resolverFiles, resolverErr := callResolverFindFiles(r, cat, filename)
				assertParity(t, topFiles, topErr, resolverFiles, resolverErr)

				if seeded {
					wantFirst := filepath.Join(seededDir, filename)
					if len(topFiles) == 0 || topFiles[0] != wantFirst {
						t.Fatalf("expected first matching file %q, got %#v", wantFirst, topFiles)
					}
				}
			})

			t.Run(fmt.Sprintf("file_%s_%s", scope, cat), func(t *testing.T) {
				topFile, topErr := callTopFindFile(scope, cat, appName, filename)
				resolverFile, resolverErr := callResolverFindFile(r, cat, filename)
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
		t.Run(fmt.Sprintf("defaulted_app_%s", scope), func(t *testing.T) {
			r := newResolverForScope(scope, "")
			topDir := callTopDir(scope, categoryData, "")
			resolverDir := callResolverDir(r, categoryData)
			assertValueParity(t, topDir, resolverDir)
			if !strings.Contains(topDir, "unnamed_app") {
				t.Fatalf("expected defaulted app name in path, got %q", topDir)
			}
		})
	}

	appName := "demo_top_level_errors"
	for _, scope := range testScopes {
		r := newResolverForScope(scope, appName)

		t.Run(fmt.Sprintf("invalid_filename_%s", scope), func(t *testing.T) {
			_, topErr := callTopFindFiles(scope, categoryData, appName, "")
			_, resolverErr := callResolverFindFiles(r, categoryData, "")
			assertErrorParity(t, topErr, resolverErr)
		})

		t.Run(fmt.Sprintf("missing_file_%s", scope), func(t *testing.T) {
			filename := fmt.Sprintf("missing_%s.txt", scope)
			_, topErr := callTopFindFile(scope, categoryData, appName, filename)
			_, resolverErr := callResolverFindFile(r, categoryData, filename)
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

func newResolverForScope(scope Scope, appName string) *Resolver {
	switch scope {
	case ScopeLocal:
		return NewLocalResolver(appName)
	case ScopeUser:
		return NewUserResolver(appName)
	case ScopeSystem:
		return NewSystemResolver(appName)
	default:
		return nil
	}
}

func assertParity[T any](t *testing.T, topValue T, topErr error, resolverValue T, resolverErr error) {
	t.Helper()
	assertErrorParity(t, topErr, resolverErr)
	if topErr == nil && !reflect.DeepEqual(topValue, resolverValue) {
		t.Fatalf("result mismatch:\nresolver: %#v\ntop:      %#v", resolverValue, topValue)
	}
}

func assertValueParity[T any](t *testing.T, topValue T, resolverValue T) {
	t.Helper()
	if !reflect.DeepEqual(topValue, resolverValue) {
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

func seedExactFilePath(path string, content string) bool {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return false
	}
	return true
}

func callTopDirs(scope Scope, cat category, appName string) []string {
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
	return nil
}

func callTopDir(scope Scope, cat category, appName string) string {
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
	return ""
}

func callTopFilePaths(scope Scope, cat category, appName string, filename string) []string {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return LocalDataFilePaths(appName, filename)
		case categoryConfig:
			return LocalConfigFilePaths(appName, filename)
		case categoryLog:
			return LocalLogFilePaths(appName, filename)
		case categoryCache:
			return LocalCacheFilePaths(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return UserDataFilePaths(appName, filename)
		case categoryConfig:
			return UserConfigFilePaths(appName, filename)
		case categoryLog:
			return UserLogFilePaths(appName, filename)
		case categoryCache:
			return UserCacheFilePaths(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return SystemDataFilePaths(appName, filename)
		case categoryConfig:
			return SystemConfigFilePaths(appName, filename)
		case categoryLog:
			return SystemLogFilePaths(appName, filename)
		case categoryCache:
			return SystemCacheFilePaths(appName, filename)
		}
	}
	return nil
}

func callTopFilePath(scope Scope, cat category, appName string, filename string) string {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return LocalDataFilePath(appName, filename)
		case categoryConfig:
			return LocalConfigFilePath(appName, filename)
		case categoryLog:
			return LocalLogFilePath(appName, filename)
		case categoryCache:
			return LocalCacheFilePath(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return UserDataFilePath(appName, filename)
		case categoryConfig:
			return UserConfigFilePath(appName, filename)
		case categoryLog:
			return UserLogFilePath(appName, filename)
		case categoryCache:
			return UserCacheFilePath(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return SystemDataFilePath(appName, filename)
		case categoryConfig:
			return SystemConfigFilePath(appName, filename)
		case categoryLog:
			return SystemLogFilePath(appName, filename)
		case categoryCache:
			return SystemCacheFilePath(appName, filename)
		}
	}
	return ""
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

func callTopEnsureWithOpts(scope Scope, cat category, appName string, opts ...EnsureOption) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return EnsureLocalDataDir(appName, opts...)
		case categoryConfig:
			return EnsureLocalConfigDir(appName, opts...)
		case categoryLog:
			return EnsureLocalLogDir(appName, opts...)
		case categoryCache:
			return EnsureLocalCacheDir(appName, opts...)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return EnsureUserDataDir(appName, opts...)
		case categoryConfig:
			return EnsureUserConfigDir(appName, opts...)
		case categoryLog:
			return EnsureUserLogDir(appName, opts...)
		case categoryCache:
			return EnsureUserCacheDir(appName, opts...)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return EnsureSystemDataDir(appName, opts...)
		case categoryConfig:
			return EnsureSystemConfigDir(appName, opts...)
		case categoryLog:
			return EnsureSystemLogDir(appName, opts...)
		case categoryCache:
			return EnsureSystemCacheDir(appName, opts...)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopFindFiles(scope Scope, cat category, appName string, filename string) ([]string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return FindLocalDataFiles(appName, filename)
		case categoryConfig:
			return FindLocalConfigFiles(appName, filename)
		case categoryLog:
			return FindLocalLogFiles(appName, filename)
		case categoryCache:
			return FindLocalCacheFiles(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return FindUserDataFiles(appName, filename)
		case categoryConfig:
			return FindUserConfigFiles(appName, filename)
		case categoryLog:
			return FindUserLogFiles(appName, filename)
		case categoryCache:
			return FindUserCacheFiles(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return FindSystemDataFiles(appName, filename)
		case categoryConfig:
			return FindSystemConfigFiles(appName, filename)
		case categoryLog:
			return FindSystemLogFiles(appName, filename)
		case categoryCache:
			return FindSystemCacheFiles(appName, filename)
		}
	}
	return nil, fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopFindFile(scope Scope, cat category, appName string, filename string) (string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return FindLocalDataFile(appName, filename)
		case categoryConfig:
			return FindLocalConfigFile(appName, filename)
		case categoryLog:
			return FindLocalLogFile(appName, filename)
		case categoryCache:
			return FindLocalCacheFile(appName, filename)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return FindUserDataFile(appName, filename)
		case categoryConfig:
			return FindUserConfigFile(appName, filename)
		case categoryLog:
			return FindUserLogFile(appName, filename)
		case categoryCache:
			return FindUserCacheFile(appName, filename)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return FindSystemDataFile(appName, filename)
		case categoryConfig:
			return FindSystemConfigFile(appName, filename)
		case categoryLog:
			return FindSystemLogFile(appName, filename)
		case categoryCache:
			return FindSystemCacheFile(appName, filename)
		}
	}
	return "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callTopCreateFile(scope Scope, cat category, appName string, filename string, opts ...CreateFileOption) (bool, string, error) {
	switch scope {
	case ScopeLocal:
		switch cat {
		case categoryData:
			return CreateLocalDataFile(appName, filename, opts...)
		case categoryConfig:
			return CreateLocalConfigFile(appName, filename, opts...)
		case categoryLog:
			return CreateLocalLogFile(appName, filename, opts...)
		case categoryCache:
			return CreateLocalCacheFile(appName, filename, opts...)
		}
	case ScopeUser:
		switch cat {
		case categoryData:
			return CreateUserDataFile(appName, filename, opts...)
		case categoryConfig:
			return CreateUserConfigFile(appName, filename, opts...)
		case categoryLog:
			return CreateUserLogFile(appName, filename, opts...)
		case categoryCache:
			return CreateUserCacheFile(appName, filename, opts...)
		}
	case ScopeSystem:
		switch cat {
		case categoryData:
			return CreateSystemDataFile(appName, filename, opts...)
		case categoryConfig:
			return CreateSystemConfigFile(appName, filename, opts...)
		case categoryLog:
			return CreateSystemLogFile(appName, filename, opts...)
		case categoryCache:
			return CreateSystemCacheFile(appName, filename, opts...)
		}
	}
	return false, "", fmt.Errorf("unsupported scope/category: %s/%s", scope, cat)
}

func callResolverDirs(r *Resolver, cat category) []string {
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
		return nil
	}
}

func callResolverDir(r *Resolver, cat category) string {
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
		return ""
	}
}

func callResolverFilePaths(r *Resolver, cat category, filename string) []string {
	switch cat {
	case categoryData:
		return r.DataFilePaths(filename)
	case categoryConfig:
		return r.ConfigFilePaths(filename)
	case categoryLog:
		return r.LogFilePaths(filename)
	case categoryCache:
		return r.CacheFilePaths(filename)
	default:
		return nil
	}
}

func callResolverFilePath(r *Resolver, cat category, filename string) string {
	switch cat {
	case categoryData:
		return r.DataFilePath(filename)
	case categoryConfig:
		return r.ConfigFilePath(filename)
	case categoryLog:
		return r.LogFilePath(filename)
	case categoryCache:
		return r.CacheFilePath(filename)
	default:
		return ""
	}
}

func callResolverEnsure(r *Resolver, cat category) (string, error) {
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

func callResolverEnsureWithOpts(r *Resolver, cat category, opts ...EnsureOption) (string, error) {
	switch cat {
	case categoryData:
		return r.EnsureDataDir(opts...)
	case categoryConfig:
		return r.EnsureConfigDir(opts...)
	case categoryLog:
		return r.EnsureLogDir(opts...)
	case categoryCache:
		return r.EnsureCacheDir(opts...)
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverFindFiles(r *Resolver, cat category, filename string) ([]string, error) {
	switch cat {
	case categoryData:
		return r.FindDataFiles(filename)
	case categoryConfig:
		return r.FindConfigFiles(filename)
	case categoryLog:
		return r.FindLogFiles(filename)
	case categoryCache:
		return r.FindCacheFiles(filename)
	default:
		return nil, fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverFindFile(r *Resolver, cat category, filename string) (string, error) {
	switch cat {
	case categoryData:
		return r.FindDataFile(filename)
	case categoryConfig:
		return r.FindConfigFile(filename)
	case categoryLog:
		return r.FindLogFile(filename)
	case categoryCache:
		return r.FindCacheFile(filename)
	default:
		return "", fmt.Errorf("unsupported category %s", cat)
	}
}

func callResolverCreateFile(r *Resolver, cat category, filename string, opts ...CreateFileOption) (bool, string, error) {
	switch cat {
	case categoryData:
		return r.CreateDataFile(filename, opts...)
	case categoryConfig:
		return r.CreateConfigFile(filename, opts...)
	case categoryLog:
		return r.CreateLogFile(filename, opts...)
	case categoryCache:
		return r.CreateCacheFile(filename, opts...)
	default:
		return false, "", fmt.Errorf("unsupported category %s", cat)
	}
}
