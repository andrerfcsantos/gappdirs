//go:build windows

package gappdirs

import (
	"errors"
	"path/filepath"
	"reflect"
	"testing"

	"golang.org/x/sys/windows"
)

func TestWindowsUserAndSystemDirsOrdering(t *testing.T) {
	oldKnownFolderPath := knownFolderPath
	knownFolderPath = func(folder knownFolder) (string, error) {
		switch folder {
		case knownFolderLocalAppData:
			return filepath.Join(`C:\`, "Users", "alice", "AppData", "Local"), nil
		case knownFolderRoamingAppData:
			return filepath.Join(`C:\`, "Users", "alice", "AppData", "Roaming"), nil
		case knownFolderProgramData:
			return filepath.Join(`C:\`, "ProgramData"), nil
		default:
			return "", errors.New("unknown folder")
		}
	}
	defer func() { knownFolderPath = oldKnownFolderPath }()

	gotData, err := platformUserDirs("demo", categoryData)
	if err != nil {
		t.Fatalf("data dirs error: %v", err)
	}
	wantData := []string{
		filepath.Join(`C:\`, "Users", "alice", "AppData", "Local", "demo", "Data"),
		filepath.Join(`C:\`, "Users", "alice", "AppData", "Roaming", "demo", "Data"),
	}
	if !reflect.DeepEqual(gotData, wantData) {
		t.Fatalf("data dirs mismatch:\nwant: %#v\ngot:  %#v", wantData, gotData)
	}

	gotConfig, err := platformUserDirs("demo", categoryConfig)
	if err != nil {
		t.Fatalf("config dirs error: %v", err)
	}
	wantConfig := []string{
		filepath.Join(`C:\`, "Users", "alice", "AppData", "Roaming", "demo", "Config"),
		filepath.Join(`C:\`, "Users", "alice", "AppData", "Local", "demo", "Config"),
	}
	if !reflect.DeepEqual(gotConfig, wantConfig) {
		t.Fatalf("config dirs mismatch:\nwant: %#v\ngot:  %#v", wantConfig, gotConfig)
	}

	gotSystem, err := platformSystemDirs("demo", categoryLog)
	if err != nil {
		t.Fatalf("system dirs error: %v", err)
	}
	wantSystem := []string{filepath.Join(`C:\`, "ProgramData", "demo", "Logs")}
	if !reflect.DeepEqual(gotSystem, wantSystem) {
		t.Fatalf("system dirs mismatch:\nwant: %#v\ngot:  %#v", wantSystem, gotSystem)
	}
}

func TestWindowsKnownFolderFallbackToEnv(t *testing.T) {
	oldLookup := knownFolderPathLookup
	knownFolderPathLookup = func(folderID *windows.KNOWNFOLDERID, flags uint32) (string, error) {
		return "", errors.New("forced failure")
	}
	defer func() { knownFolderPathLookup = oldLookup }()

	t.Setenv("LOCALAPPDATA", filepath.Join(`C:\`, "Users", "alice", "AppData", "Local"))

	got, err := lookupKnownFolderPath(knownFolderLocalAppData)
	if err != nil {
		t.Fatalf("lookup error: %v", err)
	}
	want := filepath.Join(`C:\`, "Users", "alice", "AppData", "Local")
	if got != want {
		t.Fatalf("lookup mismatch: want %q, got %q", want, got)
	}
}
