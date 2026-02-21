//go:build windows

package gappdirs

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

type knownFolder int

const (
	knownFolderLocalAppData knownFolder = iota
	knownFolderRoamingAppData
	knownFolderProgramData
)

var knownFolderPath = lookupKnownFolderPath

type knownFolderPathLookupFunc func(folderID *windows.KNOWNFOLDERID, flags uint32) (string, error)

var knownFolderPathLookup knownFolderPathLookupFunc = windows.KnownFolderPath

func lookupKnownFolderPath(folder knownFolder) (string, error) {
	folderID, envKeys, err := folderSpec(folder)
	if err != nil {
		return "", err
	}

	if knownFolderPathLookup != nil {
		if path, err := knownFolderPathLookup(folderID, windows.KF_FLAG_DEFAULT); err == nil {
			if path = strings.TrimSpace(path); path != "" {
				return path, nil
			}
		}
	}

	for _, key := range envKeys {
		if value := strings.TrimSpace(os.Getenv(key)); value != "" {
			return value, nil
		}
	}

	return "", fmt.Errorf("gappdirs: could not resolve windows known folder %s", folder)
}

func folderSpec(folder knownFolder) (*windows.KNOWNFOLDERID, []string, error) {
	switch folder {
	case knownFolderLocalAppData:
		return windows.FOLDERID_LocalAppData, []string{"LOCALAPPDATA"}, nil
	case knownFolderRoamingAppData:
		return windows.FOLDERID_RoamingAppData, []string{"APPDATA"}, nil
	case knownFolderProgramData:
		return windows.FOLDERID_ProgramData, []string{"ProgramData", "PROGRAMDATA"}, nil
	default:
		return nil, nil, fmt.Errorf("gappdirs: unsupported windows known folder %d", folder)
	}
}

func (folder knownFolder) String() string {
	switch folder {
	case knownFolderLocalAppData:
		return "LocalAppData"
	case knownFolderRoamingAppData:
		return "RoamingAppData"
	case knownFolderProgramData:
		return "ProgramData"
	default:
		return fmt.Sprintf("knownFolder(%d)", int(folder))
	}
}
