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

	t.Run("validates providers", func(t *testing.T) {
		_, err := buildScopedContext("demo", nil, nil, systemFn, nil)
		if err == nil {
			t.Fatal("expected error for nil provider")
		}
	})

	t.Run("applies options and forced scope", func(t *testing.T) {
		forcedScope := ScopeSystem
		ctx, err := buildScopedContext(
			"A/B Demo",
			[]Option{WithScope(ScopeLocal), WithWorkingDir(wd), WithDefaultDirPerm(0o700)},
			userFn,
			systemFn,
			&forcedScope,
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.appName != "a_b_demo" {
			t.Fatalf("sanitized app name mismatch: want %q, got %q", "a_b_demo", ctx.appName)
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
	ctx, err := newDefaultScopedContext("Demo App", ScopeLocal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ctx.appName != "demo_app" {
		t.Fatalf("sanitized app name mismatch: want %q, got %q", "demo_app", ctx.appName)
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

func TestScopedDirsAndDir(t *testing.T) {
	wd := t.TempDir()
	dupBase := filepath.Join(wd, "shared")
	systemData := filepath.Join(wd, "system", "data")

	ctx, err := buildScopedContext(
		"demo",
		[]Option{WithScope(ScopeLocal), WithWorkingDir(wd)},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryData {
				return nil, nil
			}
			return []string{dupBase, filepath.Join(dupBase, ".")}, nil
		},
		func(_ string, cat category) ([]string, error) {
			if cat != categoryData {
				return nil, nil
			}
			return []string{dupBase, systemData}, nil
		},
		nil,
	)
	if err != nil {
		t.Fatalf("build scoped context: %v", err)
	}

	dirs, err := scopedDirs(ctx, categoryData)
	if err != nil {
		t.Fatalf("scoped dirs: %v", err)
	}

	wantDirs := []string{
		filepath.Join(wd, ".demo", "data"),
		dupBase,
		systemData,
	}
	if !reflect.DeepEqual(dirs, wantDirs) {
		t.Fatalf("dirs mismatch:\nwant: %#v\ngot:  %#v", wantDirs, dirs)
	}

	dir, err := scopedDir(ctx, categoryData)
	if err != nil {
		t.Fatalf("scoped dir: %v", err)
	}
	if dir != wantDirs[0] {
		t.Fatalf("dir mismatch: want %q, got %q", wantDirs[0], dir)
	}
}

func TestScopedEnsureDir(t *testing.T) {
	wd := t.TempDir()
	userCache := filepath.Join(wd, "user", "cache")
	systemCache := filepath.Join(wd, "system", "cache")

	ctx, err := buildScopedContext(
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
	if err != nil {
		t.Fatalf("build scoped context: %v", err)
	}

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

	createdWithPerm, err := scopedEnsureDirWithPerm(ctx, categoryCache, 0o755)
	if err != nil {
		t.Fatalf("scoped ensure dir with perm: %v", err)
	}
	if createdWithPerm != userCache {
		t.Fatalf("created dir mismatch: want %q, got %q", userCache, createdWithPerm)
	}
}

func TestScopedFindFileDirsAndFile(t *testing.T) {
	wd := t.TempDir()
	localConfig := filepath.Join(wd, ".demo", "config")
	userConfig := filepath.Join(wd, "user", "config")
	systemConfig := filepath.Join(wd, "system", "config")

	ctx, err := buildScopedContext(
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
	if err != nil {
		t.Fatalf("build scoped context: %v", err)
	}

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

	gotDirs, err := scopedFindFileDirs(ctx, categoryConfig, "settings.yaml")
	if err != nil {
		t.Fatalf("scoped find dirs: %v", err)
	}
	wantDirs := []string{localConfig, systemConfig}
	if !reflect.DeepEqual(gotDirs, wantDirs) {
		t.Fatalf("matching dirs mismatch:\nwant: %#v\ngot:  %#v", wantDirs, gotDirs)
	}

	gotFile, err := scopedFile(ctx, categoryConfig, "settings.yaml")
	if err != nil {
		t.Fatalf("scoped file: %v", err)
	}
	wantFile := filepath.Join(localConfig, "settings.yaml")
	if gotFile != wantFile {
		t.Fatalf("most relevant file mismatch: want %q, got %q", wantFile, gotFile)
	}

	_, err = scopedFile(ctx, categoryConfig, "missing.yaml")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	_, err = scopedFindFileDirs(ctx, categoryConfig, "a/b")
	if err == nil {
		t.Fatal("expected filename validation error")
	}
}

func TestScopedFunctionsValidateContext(t *testing.T) {
	_, err := scopedDirs(scopedContext{}, categoryData)
	if err == nil {
		t.Fatal("expected error for uninitialized scoped context")
	}

	ctx := scopedContext{
		appName:      "demo",
		scope:        Scope(-1),
		workingDir:   t.TempDir(),
		userDirsFn:   func(string, category) ([]string, error) { return nil, nil },
		systemDirsFn: func(string, category) ([]string, error) { return nil, nil },
	}
	_, err = scopedDirs(ctx, categoryData)
	if err == nil {
		t.Fatal("expected unsupported scope error")
	}
}
