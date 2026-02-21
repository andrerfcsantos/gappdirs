//go:build darwin

package gappdirs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDarwinUserAndSystemDirs(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("home directory error: %v", err)
	}

	for _, tc := range []struct {
		category   category
		wantUser   []string
		wantSystem []string
	}{
		{
			category:   categoryData,
			wantUser:   []string{filepath.Join(homeDir, "Library", "Application Support", "demo", "data")},
			wantSystem: []string{filepath.Join("/Library", "Application Support", "demo", "data")},
		},
		{
			category:   categoryConfig,
			wantUser:   []string{filepath.Join(homeDir, "Library", "Application Support", "demo", "config")},
			wantSystem: []string{filepath.Join("/Library", "Application Support", "demo", "config")},
		},
		{
			category:   categoryLog,
			wantUser:   []string{filepath.Join(homeDir, "Library", "Logs", "demo")},
			wantSystem: []string{filepath.Join("/Library", "Logs", "demo")},
		},
		{
			category:   categoryCache,
			wantUser:   []string{filepath.Join(homeDir, "Library", "Caches", "demo")},
			wantSystem: []string{filepath.Join("/Library", "Caches", "demo")},
		},
	} {
		gotUser, err := platformUserDirs("demo", tc.category)
		if err != nil {
			t.Fatalf("user dirs error for %s: %v", tc.category, err)
		}
		if !reflect.DeepEqual(gotUser, tc.wantUser) {
			t.Fatalf("user dirs mismatch for %s:\nwant: %#v\ngot:  %#v", tc.category, tc.wantUser, gotUser)
		}

		gotSystem, err := platformSystemDirs("demo", tc.category)
		if err != nil {
			t.Fatalf("system dirs error for %s: %v", tc.category, err)
		}
		if !reflect.DeepEqual(gotSystem, tc.wantSystem) {
			t.Fatalf("system dirs mismatch for %s:\nwant: %#v\ngot:  %#v", tc.category, tc.wantSystem, gotSystem)
		}
	}
}
