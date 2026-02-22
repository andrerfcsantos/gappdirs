package gappdirs

import (
	"io/fs"
	"testing"
)

func TestResolveEnsureConfig(t *testing.T) {
	t.Run("uses default permission with no options", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, nil)
		if got.dirPerm != 0o755 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o755, got.dirPerm)
		}
	})

	t.Run("applies valid override", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, []EnsureOption{WithEnsureDirPerm(0o700)})
		if got.dirPerm != 0o700 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o700, got.dirPerm)
		}
	})

	t.Run("ignores zero permission and falls back to default", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, []EnsureOption{WithEnsureDirPerm(0)})
		if got.dirPerm != 0o755 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o755, got.dirPerm)
		}
	})

	t.Run("ignores mode-type bits and falls back to default", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, []EnsureOption{WithEnsureDirPerm(fs.ModeDir | 0o700)})
		if got.dirPerm != 0o755 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o755, got.dirPerm)
		}
	})

	t.Run("ignores nil option", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, []EnsureOption{nil, WithEnsureDirPerm(0o700)})
		if got.dirPerm != 0o700 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o700, got.dirPerm)
		}
	})

	t.Run("last valid option wins", func(t *testing.T) {
		got := resolveEnsureConfig(0o755, []EnsureOption{
			WithEnsureDirPerm(0o700),
			WithEnsureDirPerm(0),
			WithEnsureDirPerm(fs.ModeDir | 0o711),
			WithEnsureDirPerm(0o750),
		})
		if got.dirPerm != 0o750 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o750, got.dirPerm)
		}
	})
}
