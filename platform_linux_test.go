//go:build linux

package gappdirs

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLinuxUserDirsDefaultsAndOverrides(t *testing.T) {
	homeDir := filepath.Join(t.TempDir(), "home")
	t.Setenv("HOME", homeDir)
	t.Setenv("XDG_DATA_HOME", "")
	t.Setenv("XDG_CONFIG_HOME", "relative/path")
	t.Setenv("XDG_CACHE_HOME", "/tmp/cache-home")
	t.Setenv("XDG_STATE_HOME", "")

	gotData, err := platformUserDirs("demo", categoryData)
	if err != nil {
		t.Fatalf("data dirs error: %v", err)
	}
	wantData := []string{filepath.Join(homeDir, ".local", "share", "demo")}
	if !reflect.DeepEqual(gotData, wantData) {
		t.Fatalf("data dirs mismatch:\nwant: %#v\ngot:  %#v", wantData, gotData)
	}

	gotConfig, err := platformUserDirs("demo", categoryConfig)
	if err != nil {
		t.Fatalf("config dirs error: %v", err)
	}
	wantConfig := []string{filepath.Join(homeDir, ".config", "demo")}
	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Fatalf("config dirs mismatch:\nwant: %#v\ngot:  %#v", wantConfig, gotConfig)
	}

	gotCache, err := platformUserDirs("demo", categoryCache)
	if err != nil {
		t.Fatalf("cache dirs error: %v", err)
	}
	wantCache := []string{filepath.Join("/tmp/cache-home", "demo")}
	if !reflect.DeepEqual(gotCache, wantCache) {
		t.Fatalf("cache dirs mismatch:\nwant: %#v\ngot:  %#v", wantCache, gotCache)
	}

	gotLog, err := platformUserDirs("demo", categoryLog)
	if err != nil {
		t.Fatalf("log dirs error: %v", err)
	}
	wantLog := []string{filepath.Join(homeDir, ".local", "state", "demo", "log")}
	if !reflect.DeepEqual(gotLog, wantLog) {
		t.Fatalf("log dirs mismatch:\nwant: %#v\ngot:  %#v", wantLog, gotLog)
	}
}

func TestLinuxSystemDataAndConfigDirs(t *testing.T) {
	t.Setenv("XDG_DATA_DIRS", "/opt/share:/usr/share")
	t.Setenv("XDG_CONFIG_DIRS", "")

	gotData, err := platformSystemDirs("demo", categoryData)
	if err != nil {
		t.Fatalf("data dirs error: %v", err)
	}
	wantData := []string{
		filepath.Join("/var/lib", "demo"),
		filepath.Join("/opt/share", "demo"),
		filepath.Join("/usr/share", "demo"),
	}
	if !reflect.DeepEqual(gotData, wantData) {
		t.Fatalf("data dirs mismatch:\nwant: %#v\ngot:  %#v", wantData, gotData)
	}

	gotConfig, err := platformSystemDirs("demo", categoryConfig)
	if err != nil {
		t.Fatalf("config dirs error: %v", err)
	}
	wantConfig := []string{filepath.Join("/etc/xdg", "demo")}
	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Fatalf("config dirs mismatch:\nwant: %#v\ngot:  %#v", wantConfig, gotConfig)
	}
}

func TestLinuxUserDirsRequireHome(t *testing.T) {
	t.Setenv("HOME", "")
	if _, err := os.UserHomeDir(); err == nil {
		t.Skip("os.UserHomeDir resolved home from system user database")
	}
	if _, err := platformUserDirs("demo", categoryData); err == nil {
		t.Fatal("expected error when home directory cannot be resolved")
	}
}
