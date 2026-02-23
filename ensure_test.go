package gappdirs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnsureCacheDirWithOptionPerm(t *testing.T) {
	wd := t.TempDir()
	userCache := filepath.Join(wd, "user", "cache")
	systemCache := filepath.Join(wd, "system", "cache")

	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryCache: {userCache}},
		map[category][]string{categoryCache: {systemCache}},
	)

	created, err := r.EnsureCacheDir(WithEnsureDirPerm(0o755))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created != userCache {
		t.Fatalf("created dir mismatch: want %q, got %q", userCache, created)
	}

	info, err := os.Stat(userCache)
	if err != nil {
		t.Fatalf("expected directory to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected %q to be a directory", userCache)
	}
}

func TestEnsureCacheDirUsesDefaultPermSetting(t *testing.T) {
	wd := t.TempDir()
	userCache := filepath.Join(wd, "user", "cache")
	systemCache := filepath.Join(wd, "system", "cache")

	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryCache: {userCache}},
		map[category][]string{categoryCache: {systemCache}},
		WithDefaultDirPerm(0o700),
	)

	if r.ctx.defaultDirPerm != 0o700 {
		t.Fatalf("Resolver default permission mismatch: want %o, got %o", 0o700, r.ctx.defaultDirPerm)
	}

	created, err := r.EnsureCacheDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created != userCache {
		t.Fatalf("created dir mismatch: want %q, got %q", userCache, created)
	}
}
