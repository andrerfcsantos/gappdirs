package gappdirs

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestCreateFileCreatesInMostRelevantDirectory(t *testing.T) {
	wd := t.TempDir()
	localData := filepath.Join(wd, ".demo", "data")
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")

	r := mustResolver(t, ScopeLocal, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	created, path, err := r.CreateDataFile("settings.db")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected created=true for new file")
	}

	want := filepath.Join(localData, "settings.db")
	if path != want {
		t.Fatalf("created path mismatch: want %q, got %q", want, path)
	}
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
}

func TestCreateFileWritesReaderData(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	created, path, err := r.CreateDataFile("state.txt", WithContentsFromReader(strings.NewReader("payload")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected created=true for new file")
	}
	if path != filepath.Join(userData, "state.txt") {
		t.Fatalf("created path mismatch: got %q", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(content) != "payload" {
		t.Fatalf("content mismatch: want %q, got %q", "payload", string(content))
	}
}

func TestCreateFileDoesNotOverwriteByDefault(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	existing := filepath.Join(userData, "state.txt")
	if err := os.MkdirAll(userData, 0o755); err != nil {
		t.Fatalf("mkdir user data: %v", err)
	}
	if err := os.WriteFile(existing, []byte("old"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	created, path, err := r.CreateDataFile("state.txt", WithContentsFromReader(strings.NewReader("new")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Fatal("expected created=false when file already exists")
	}
	if path != existing {
		t.Fatalf("path mismatch: want %q, got %q", existing, path)
	}

	content, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(content) != "old" {
		t.Fatalf("existing file should stay untouched: got %q", string(content))
	}
}

func TestCreateFileOverwriteWithReader(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	existing := filepath.Join(userData, "state.txt")
	if err := os.MkdirAll(userData, 0o755); err != nil {
		t.Fatalf("mkdir user data: %v", err)
	}
	if err := os.WriteFile(existing, []byte("old"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	created, path, err := r.CreateDataFile(
		"state.txt",
		WithOverwriteExisting(),
		WithContentsFromReader(strings.NewReader("new")),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Fatal("expected created=false when overwriting existing file")
	}
	if path != existing {
		t.Fatalf("path mismatch: want %q, got %q", existing, path)
	}

	content, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(content) != "new" {
		t.Fatalf("overwritten content mismatch: got %q", string(content))
	}
}

func TestCreateFileOverwriteWithNilReaderTruncates(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	existing := filepath.Join(userData, "state.txt")
	if err := os.MkdirAll(userData, 0o755); err != nil {
		t.Fatalf("mkdir user data: %v", err)
	}
	if err := os.WriteFile(existing, []byte("old"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	created, path, err := r.CreateDataFile("state.txt", WithOverwriteExisting())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if created {
		t.Fatal("expected created=false when overwriting existing file")
	}
	if path != existing {
		t.Fatalf("path mismatch: want %q, got %q", existing, path)
	}

	content, err := os.ReadFile(existing)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if len(content) != 0 {
		t.Fatalf("expected file to be truncated, got %d bytes", len(content))
	}
}

func TestCreateFileOverwriteEnabledStillCreatesMissingFile(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	created, path, err := r.CreateDataFile(
		"new.txt",
		WithOverwriteExisting(),
		WithContentsFromReader(strings.NewReader("new")),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected created=true when file did not exist")
	}
	if path != filepath.Join(userData, "new.txt") {
		t.Fatalf("path mismatch: got %q", path)
	}
}

func TestCreateFileCreatesNestedParentDirectories(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	created, path, err := r.CreateDataFile("nested/dir/state.txt", WithContentsFromReader(strings.NewReader("ok")))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Fatal("expected created=true for new nested file")
	}
	if path != filepath.Join(userData, "nested", "dir", "state.txt") {
		t.Fatalf("path mismatch: got %q", path)
	}
	if _, err := os.Stat(filepath.Join(userData, "nested", "dir")); err != nil {
		t.Fatalf("expected nested directories to exist: %v", err)
	}
}

func TestCreateFileValidatesFilename(t *testing.T) {
	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	_, _, err := r.CreateDataFile("")
	if err == nil {
		t.Fatal("expected filename validation error")
	}
}

func TestCreateFilePermissionOption(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file mode permission checks are not reliable on windows")
	}

	wd := t.TempDir()
	userData := filepath.Join(wd, "user", "data")
	systemData := filepath.Join(wd, "system", "data")
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {userData}},
		map[category][]string{categoryData: {systemData}},
	)

	_, path, err := r.CreateDataFile("mode.txt", WithFilePerm(0o600))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat file: %v", err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("permission mismatch: want %o, got %o", 0o600, got)
	}
}
