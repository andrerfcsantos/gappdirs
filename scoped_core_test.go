package gappdirs

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestBuildScopedContext(t *testing.T) {
	wd := t.TempDir()
	userFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "user")}, nil }
	systemFn := func(string, category) ([]string, error) { return []string{filepath.Join(wd, "system")}, nil }

	t.Run("applies options and explicit scope", func(t *testing.T) {
		ctx := buildScopedContext(
			"A/B Demo",
			ScopeSystem,
			[]ResolverOption{WithLocalDir(wd), WithDefaultDirPerm(0o700)},
			userFn,
			systemFn,
		)
		if ctx.appName != "A_B_Demo" {
			t.Fatalf("sanitized app name mismatch: want %q, got %q", "A_B_Demo", ctx.appName)
		}
		if ctx.scope != ScopeSystem {
			t.Fatalf("scope mismatch: want %s, got %s", ScopeSystem, ctx.scope)
		}
		if ctx.defaultDirPerm != 0o700 {
			t.Fatalf("default dir perm mismatch: want %o, got %o", 0o700, ctx.defaultDirPerm)
		}
		if !reflect.DeepEqual(ctx.workingDirs, []string{wd}) {
			t.Fatalf("working dirs mismatch: want %#v, got %#v", []string{wd}, ctx.workingDirs)
		}
	})
}

func TestNewDefaultScopedContext(t *testing.T) {
	ctx := newDefaultScopedContext("Demo App", ScopeLocal)
	if ctx.appName != "Demo_App" {
		t.Fatalf("sanitized app name mismatch: want %q, got %q", "Demo_App", ctx.appName)
	}
	if ctx.scope != ScopeLocal {
		t.Fatalf("scope mismatch: want %s, got %s", ScopeLocal, ctx.scope)
	}
	if ctx.defaultDirPerm != DefaultDirPerm {
		t.Fatalf("default dir perm mismatch: want %o, got %o", DefaultDirPerm, ctx.defaultDirPerm)
	}
	if ctx.userDirsFn == nil || ctx.systemDirsFn == nil {
		t.Fatal("default providers must be set")
	}
}

func TestScopedEnsureDir(t *testing.T) {
	wd := t.TempDir()
	userCache := filepath.Join(wd, "user", "cache")
	systemCache := filepath.Join(wd, "system", "cache")

	ctx := buildScopedContext(
		"demo",
		ScopeUser,
		[]ResolverOption{WithLocalDir(wd), WithDefaultDirPerm(0o700)},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryCache {
				return nil, nil
			}
			return []string{userCache}, nil
		},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryCache {
				return nil, nil
			}
			return []string{systemCache}, nil
		},
	)

	created, err := scopedEnsureDir(ctx, categoryCache)
	if err != nil {
		t.Fatalf("scoped ensure dir: %v", err)
	}
	if created != userCache {
		t.Fatalf("created dir mismatch: want %q, got %q", userCache, created)
	}
	if info, err := os.Stat(created); err != nil || !info.IsDir() {
		t.Fatalf("expected created directory %q to exist", created)
	}

	createdWithPerm, err := scopedEnsureDir(ctx, categoryCache, WithEnsureDirPerm(0o755))
	if err != nil {
		t.Fatalf("scoped ensure dir with option perm: %v", err)
	}
	if createdWithPerm != userCache {
		t.Fatalf("created dir mismatch: want %q, got %q", userCache, createdWithPerm)
	}
}

func TestScopedFindFilesAndFile(t *testing.T) {
	wd := t.TempDir()
	localConfig := filepath.Join(wd, ".demo", "config")
	userConfig := filepath.Join(wd, "user", "config")
	systemConfig := filepath.Join(wd, "system", "config")

	ctx := buildScopedContext(
		"demo",
		ScopeLocal,
		[]ResolverOption{WithLocalDir(wd)},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryConfig {
				return nil, nil
			}
			return []string{userConfig}, nil
		},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryConfig {
				return nil, nil
			}
			return []string{systemConfig}, nil
		},
	)

	if err := os.MkdirAll(localConfig, 0o755); err != nil {
		t.Fatalf("mkdir local config: %v", err)
	}
	if err := os.MkdirAll(systemConfig, 0o755); err != nil {
		t.Fatalf("mkdir system config: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localConfig, "settings.yaml"), []byte("local"), 0o644); err != nil {
		t.Fatalf("write local file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(systemConfig, "settings.yaml"), []byte("system"), 0o644); err != nil {
		t.Fatalf("write system file: %v", err)
	}

	gotFiles, err := findExistingScopedFiles(ctx, categoryConfig, "settings.yaml")
	if err != nil {
		t.Fatalf("scoped find files: %v", err)
	}
	wantFiles := []string{
		filepath.Join(localConfig, "settings.yaml"),
		filepath.Join(systemConfig, "settings.yaml"),
	}
	if !reflect.DeepEqual(gotFiles, wantFiles) {
		t.Fatalf("matching files mismatch:\nwant: %#v\ngot:  %#v", wantFiles, gotFiles)
	}

	gotFile, err := findExistingScopedFile(ctx, categoryConfig, "settings.yaml")
	if err != nil {
		t.Fatalf("scoped file: %v", err)
	}
	wantFile := filepath.Join(localConfig, "settings.yaml")
	if gotFile != wantFile {
		t.Fatalf("most relevant file mismatch: want %q, got %q", wantFile, gotFile)
	}

	_, err = findExistingScopedFile(ctx, categoryConfig, "missing.yaml")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	_, err = findExistingScopedFiles(ctx, categoryConfig, "")
	if err == nil {
		t.Fatal("expected filename validation error")
	}
}

func TestScopedFilePathsAndFilePath(t *testing.T) {
	wd := t.TempDir()
	localConfig := filepath.Join(wd, ".demo", "config")
	userConfig := filepath.Join(wd, "user", "config")
	systemConfig := filepath.Join(wd, "system", "config")

	ctx := buildScopedContext(
		"demo",
		ScopeLocal,
		[]ResolverOption{WithLocalDir(wd)},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryConfig {
				return nil, nil
			}
			return []string{userConfig}, nil
		},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryConfig {
				return nil, nil
			}
			return []string{systemConfig}, nil
		},
	)

	t.Run("relative filename keeps precedence", func(t *testing.T) {
		gotPaths := scopedFilePaths(ctx, categoryConfig, "settings.yaml")
		wantPaths := []string{
			filepath.Join(localConfig, "settings.yaml"),
			filepath.Join(userConfig, "settings.yaml"),
			filepath.Join(systemConfig, "settings.yaml"),
		}
		if !reflect.DeepEqual(gotPaths, wantPaths) {
			t.Fatalf("computed file paths mismatch:\nwant: %#v\ngot:  %#v", wantPaths, gotPaths)
		}

		gotPath := scopedFilePath(ctx, categoryConfig, "settings.yaml")
		if gotPath != wantPaths[0] {
			t.Fatalf("highest-precedence computed file path mismatch: want %q, got %q", wantPaths[0], gotPath)
		}
	})

	t.Run("empty filename resolves to directory paths", func(t *testing.T) {
		gotPaths := scopedFilePaths(ctx, categoryConfig, "   ")
		wantPaths := []string{
			localConfig,
			userConfig,
			systemConfig,
		}
		if !reflect.DeepEqual(gotPaths, wantPaths) {
			t.Fatalf("computed file paths mismatch for empty filename:\nwant: %#v\ngot:  %#v", wantPaths, gotPaths)
		}

		gotPath := scopedFilePath(ctx, categoryConfig, "   ")
		if gotPath != localConfig {
			t.Fatalf("highest-precedence computed file path mismatch for empty filename: want %q, got %q", localConfig, gotPath)
		}
	})

	t.Run("absolute filename is anchored under category dirs", func(t *testing.T) {
		absolute := filepath.Join(string(filepath.Separator), "absolute", "rooted.yaml")
		relative := filepath.Join("absolute", "rooted.yaml")

		gotPaths := scopedFilePaths(ctx, categoryConfig, absolute)
		wantPaths := []string{
			filepath.Join(localConfig, relative),
			filepath.Join(userConfig, relative),
			filepath.Join(systemConfig, relative),
		}
		if !reflect.DeepEqual(gotPaths, wantPaths) {
			t.Fatalf("computed file paths mismatch for absolute filename:\nwant: %#v\ngot:  %#v", wantPaths, gotPaths)
		}
	})
}

func TestScopedFilePathsAndFilePathWhenNoDirs(t *testing.T) {
	ctx := buildScopedContext(
		"demo",
		ScopeUser,
		nil,
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return nil, errors.New("system provider failed") },
	)

	gotPaths := scopedFilePaths(ctx, categoryData, "settings.yaml")
	if len(gotPaths) != 0 {
		t.Fatalf("expected empty candidate paths when there are no dirs, got %#v", gotPaths)
	}

	gotPath := scopedFilePath(ctx, categoryData, "settings.yaml")
	if gotPath != "" {
		t.Fatalf("expected empty candidate path when there are no dirs, got %q", gotPath)
	}
}

func TestWorkingDirResolutionIsDeferredToLocalScope(t *testing.T) {
	preBuildWD := t.TempDir()
	runtimeWD := t.TempDir()
	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	if err := os.Chdir(preBuildWD); err != nil {
		t.Fatalf("chdir pre-build wd: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(originalWD)
	})

	ctx := buildScopedContext(
		"demo",
		ScopeLocal,
		nil,
		func(string, category) ([]string, error) { return []string{filepath.Join(preBuildWD, "user")}, nil },
		func(string, category) ([]string, error) { return []string{filepath.Join(preBuildWD, "system")}, nil },
	)

	if err := os.Chdir(runtimeWD); err != nil {
		t.Fatalf("chdir runtime wd: %v", err)
	}
	runtimeResolvedWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get runtime wd: %v", err)
	}

	dirs := scopedDirs(ctx, categoryData)
	wantLocal := filepath.Join(runtimeResolvedWD, ".demo", "data")
	if len(dirs) == 0 || dirs[0] != wantLocal {
		t.Fatalf("expected runtime local dir %q, got %#v", wantLocal, dirs)
	}
}

func TestScopedLocalDirsHonorConfiguredPriority(t *testing.T) {
	base := t.TempDir()
	localA := filepath.Join(base, "project-a")
	localB := filepath.Join(base, "project-b")
	userData := filepath.Join(base, "user", "data")
	systemData := filepath.Join(base, "system", "data")

	ctx := buildScopedContext(
		"demo",
		ScopeLocal,
		[]ResolverOption{WithLocalDirs(localA, localB)},
		func(string, category) ([]string, error) { return []string{userData}, nil },
		func(string, category) ([]string, error) { return []string{systemData}, nil },
	)

	got := scopedDirs(ctx, categoryData)
	want := []string{
		filepath.Join(localA, ".demo", "data"),
		filepath.Join(localB, ".demo", "data"),
		userData,
		systemData,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("local scope order mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestScopedEnsureDirUsesFirstLocalDirWhenMultipleConfigured(t *testing.T) {
	base := t.TempDir()
	localA := filepath.Join(base, "project-a")
	localB := filepath.Join(base, "project-b")

	ctx := buildScopedContext(
		"demo",
		ScopeLocal,
		[]ResolverOption{WithLocalDirs(localA, localB)},
		func(string, category) ([]string, error) { return nil, nil },
		func(string, category) ([]string, error) { return nil, nil },
	)

	created, err := scopedEnsureDir(ctx, categoryCache)
	if err != nil {
		t.Fatalf("scoped ensure dir: %v", err)
	}

	want := filepath.Join(localA, ".demo", "cache")
	if created != want {
		t.Fatalf("created dir mismatch: want %q, got %q", want, created)
	}

	if info, err := os.Stat(created); err != nil || !info.IsDir() {
		t.Fatalf("expected created directory %q to exist", created)
	}
}

func TestUserScopeDoesNotRequireWorkingDirResolution(t *testing.T) {
	base := t.TempDir()
	ctx := buildScopedContext(
		"demo",
		ScopeUser,
		[]ResolverOption{WithLocalDir("relative-working-dir")},
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "user")}, nil },
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "system")}, nil },
	)

	dirs := scopedDirs(ctx, categoryData)
	want := []string{filepath.Join(base, "user"), filepath.Join(base, "system")}
	if !reflect.DeepEqual(dirs, want) {
		t.Fatalf("expected user/system dirs without local resolution:\nwant: %#v\ngot:  %#v", want, dirs)
	}
}

func TestBuildScopedContextInvalidScopeFallsBackToUser(t *testing.T) {
	base := t.TempDir()
	ctx := buildScopedContext(
		"demo",
		Scope(-1),
		nil,
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "user")}, nil },
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "system")}, nil },
	)

	if ctx.scope != ScopeUser {
		t.Fatalf("scope fallback mismatch: want %s, got %s", ScopeUser, ctx.scope)
	}

	dirs := scopedDirs(ctx, categoryData)
	want := []string{filepath.Join(base, "user"), filepath.Join(base, "system")}
	if !reflect.DeepEqual(dirs, want) {
		t.Fatalf("fallback dirs mismatch:\nwant: %#v\ngot:  %#v", want, dirs)
	}
}

func TestScopedDirsIgnoresProviderErrorsWhenOtherSourcesExist(t *testing.T) {
	base := t.TempDir()
	systemData := filepath.Join(base, "system", "data")

	ctxUser := buildScopedContext(
		"demo",
		ScopeUser,
		[]ResolverOption{WithLocalDir(base)},
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return []string{systemData}, nil },
	)
	gotUser := scopedDirs(ctxUser, categoryData)
	wantUser := []string{systemData}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Fatalf("user-scope graceful fallback mismatch:\nwant: %#v\ngot:  %#v", wantUser, gotUser)
	}

	ctxLocal := buildScopedContext(
		"demo",
		ScopeLocal,
		[]ResolverOption{WithLocalDir(base)},
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return []string{systemData}, nil },
	)
	gotLocal := scopedDirs(ctxLocal, categoryData)
	wantLocal := []string{filepath.Join(base, ".demo", "data"), systemData}
	if !reflect.DeepEqual(gotLocal, wantLocal) {
		t.Fatalf("local-scope graceful fallback mismatch:\nwant: %#v\ngot:  %#v", wantLocal, gotLocal)
	}
}

func TestScopedDirsReturnsEmptyWhenAllProvidersFail(t *testing.T) {
	ctx := buildScopedContext(
		"demo",
		ScopeUser,
		nil,
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return nil, errors.New("system provider failed") },
	)

	gotDirs := scopedDirs(ctx, categoryData)
	if len(gotDirs) != 0 {
		t.Fatalf("expected empty dirs when both providers fail, got %#v", gotDirs)
	}

	gotDir := scopedDir(ctx, categoryData)
	if gotDir != "" {
		t.Fatalf("expected empty dir when both providers fail, got %q", gotDir)
	}
}
