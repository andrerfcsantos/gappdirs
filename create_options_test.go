package gappdirs

import (
	"io/fs"
	"strings"
	"testing"
)

func TestResolveCreateFileOptions(t *testing.T) {
	t.Run("uses defaults with no options", func(t *testing.T) {
		got := resolveCreateFileOptions(nil)
		if got.filePerm != DefaultCreateFilePerm {
			t.Fatalf("permission mismatch: want %o, got %o", DefaultCreateFilePerm, got.filePerm)
		}
		if got.overwriteExisting {
			t.Fatal("overwriteExisting should default to false")
		}
		if got.reader != nil {
			t.Fatal("reader should default to nil")
		}
	})

	t.Run("applies valid permission override", func(t *testing.T) {
		got := resolveCreateFileOptions([]CreateFileOption{WithFilePerm(0o600)})
		if got.filePerm != 0o600 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o600, got.filePerm)
		}
	})

	t.Run("ignores invalid permissions and falls back to default", func(t *testing.T) {
		got := resolveCreateFileOptions([]CreateFileOption{WithFilePerm(0)})
		if got.filePerm != DefaultCreateFilePerm {
			t.Fatalf("permission mismatch: want %o, got %o", DefaultCreateFilePerm, got.filePerm)
		}

		got = resolveCreateFileOptions([]CreateFileOption{WithFilePerm(fs.ModeDir | 0o700)})
		if got.filePerm != DefaultCreateFilePerm {
			t.Fatalf("permission mismatch: want %o, got %o", DefaultCreateFilePerm, got.filePerm)
		}
	})

	t.Run("ignores nil option", func(t *testing.T) {
		got := resolveCreateFileOptions([]CreateFileOption{nil, WithFilePerm(0o640)})
		if got.filePerm != 0o640 {
			t.Fatalf("permission mismatch: want %o, got %o", 0o640, got.filePerm)
		}
	})

	t.Run("applies overwrite option", func(t *testing.T) {
		got := resolveCreateFileOptions([]CreateFileOption{WithOverwriteExisting()})
		if !got.overwriteExisting {
			t.Fatal("overwriteExisting should be true")
		}
	})

	t.Run("stores reader option", func(t *testing.T) {
		reader := strings.NewReader("data")
		got := resolveCreateFileOptions([]CreateFileOption{WithContentsFromReader(reader)})
		if got.reader != reader {
			t.Fatal("reader mismatch")
		}
	})
}
