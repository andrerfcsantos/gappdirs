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

	t.Run("applies options and forced scope", func(t *testing.T) {
		forcedScope := ScopeSystem
		ctx := buildScopedContext(
			"A/B Demo",
			[]Option{WithScope(ScopeLocal), WithWorkingDir(wd), WithDefaultDirPerm(0o700)},
			userFn,
			systemFn,
			&forcedScope,
		)
		if ctx.appName != "A_B_Demo" {
			t.Fatalf("sanitized app name mismatch: want %q, got %q", "A_B_Demo", ctx.appName)
		}
		if ctx.scope != ScopeSystem {
			t.Fatalf("forced scope mismatch: want %s, got %s", ScopeSystem, ctx.scope)
		}
		if ctx.defaultDirPerm != 0o700 {
			t.Fatalf("default dir perm mismatch: want %o, got %o", 0o700, ctx.defaultDirPerm)
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
	if ctx.defaultDirPerm != defaultDirPerm {
		t.Fatalf("default dir perm mismatch: want %o, got %o", defaultDirPerm, ctx.defaultDirPerm)
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
		[]Option{WithScope(ScopeUser), WithWorkingDir(wd), WithDefaultDirPerm(0o700)},
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
		nil,
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
		[]Option{WithScope(ScopeLocal), WithWorkingDir(wd)},
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
		nil,
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
		[]Option{WithScope(ScopeLocal)},
		func(string, category) ([]string, error) { return []string{filepath.Join(preBuildWD, "user")}, nil },
		func(string, category) ([]string, error) { return []string{filepath.Join(preBuildWD, "system")}, nil },
		nil,
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

func TestUserScopeDoesNotRequireWorkingDirResolution(t *testing.T) {
	base := t.TempDir()
	ctx := buildScopedContext(
		"demo",
		[]Option{WithScope(ScopeUser), WithWorkingDir("relative-working-dir")},
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "user")}, nil },
		func(string, category) ([]string, error) { return []string{filepath.Join(base, "system")}, nil },
		nil,
	)

	dirs := scopedDirs(ctx, categoryData)
	want := []string{filepath.Join(base, "user"), filepath.Join(base, "system")}
	if !reflect.DeepEqual(dirs, want) {
		t.Fatalf("expected user/system dirs without local resolution:\nwant: %#v\ngot:  %#v", want, dirs)
	}
}

func TestInvalidScopeFallsBackToUserSemantics(t *testing.T) {
	base := t.TempDir()
	ctx := scopedContext{
		appName:      "demo",
		scope:        Scope(-1),
		workingDir:   "",
		userDirsFn:   func(string, category) ([]string, error) { return []string{filepath.Join(base, "user")}, nil },
		systemDirsFn: func(string, category) ([]string, error) { return []string{filepath.Join(base, "system")}, nil },
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
		[]Option{WithScope(ScopeUser), WithWorkingDir(base)},
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return []string{systemData}, nil },
		nil,
	)
	gotUser := scopedDirs(ctxUser, categoryData)
	wantUser := []string{systemData}
	if !reflect.DeepEqual(gotUser, wantUser) {
		t.Fatalf("user-scope graceful fallback mismatch:\nwant: %#v\ngot:  %#v", wantUser, gotUser)
	}

	ctxLocal := buildScopedContext(
		"demo",
		[]Option{WithScope(ScopeLocal), WithWorkingDir(base)},
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return []string{systemData}, nil },
		nil,
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
		[]Option{WithScope(ScopeUser)},
		func(string, category) ([]string, error) { return nil, errors.New("user provider failed") },
		func(string, category) ([]string, error) { return nil, errors.New("system provider failed") },
		nil,
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
