package gappdirs

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFindFileDirsAndFile(t *testing.T) {
	wd := t.TempDir()
	localConfig := filepath.Join(wd, ".demo", "config")
	userConfig := filepath.Join(wd, "user", "config")
	systemConfig := filepath.Join(wd, "system", "config")

	r := mustResolver(t, ScopeLocal, wd,
		map[category][]string{categoryConfig: {userConfig}},
		map[category][]string{categoryConfig: {systemConfig}},
	)

	if err := os.MkdirAll(localConfig, 0o755); err != nil {
		t.Fatalf("mkdir local: %v", err)
	}
	if err := os.WriteFile(filepath.Join(localConfig, "settings.yaml"), []byte("local"), 0o644); err != nil {
		t.Fatalf("write local file: %v", err)
	}
	if err := os.MkdirAll(systemConfig, 0o755); err != nil {
		t.Fatalf("mkdir system: %v", err)
	}
	if err := os.WriteFile(filepath.Join(systemConfig, "settings.yaml"), []byte("system"), 0o644); err != nil {
		t.Fatalf("write system file: %v", err)
	}

	gotDirs, err := r.FindConfigFileDirs("settings.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantDirs := []string{localConfig, systemConfig}
	if !reflect.DeepEqual(gotDirs, wantDirs) {
		t.Fatalf("matching dirs mismatch:\nwant: %#v\ngot:  %#v", wantDirs, gotDirs)
	}

	gotFile, err := r.ConfigFile("settings.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantFile := filepath.Join(localConfig, "settings.yaml")
	if gotFile != wantFile {
		t.Fatalf("most relevant file mismatch: want %q, got %q", wantFile, gotFile)
	}
}

func TestFileNotFound(t *testing.T) {
	wd := t.TempDir()
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {filepath.Join(wd, "user", "data")}},
		map[category][]string{categoryData: {filepath.Join(wd, "system", "data")}},
	)

	_, err := r.DataFile("missing.db")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFilenameValidation(t *testing.T) {
	wd := t.TempDir()
	r := mustResolver(t, ScopeUser, wd,
		map[category][]string{categoryData: {filepath.Join(wd, "user", "data")}},
		map[category][]string{categoryData: {filepath.Join(wd, "system", "data")}},
	)

	invalidNames := []string{"", "a/b", "a\\b", "/tmp/x"}
	for _, filename := range invalidNames {
		_, err := r.FindDataFileDirs(filename)
		if err == nil {
			t.Fatalf("expected error for filename %q", filename)
		}
	}
}
