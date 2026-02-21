//go:build windows

package gappdirs

import (
	"fmt"
	"path/filepath"
)

func platformUserDirs(appName string, cat category) ([]string, error) {
	switch cat {
	case categoryData:
		local, err := knownFolderPath(knownFolderLocalAppData)
		if err != nil {
			return nil, err
		}
		roaming, err := knownFolderPath(knownFolderRoamingAppData)
		if err != nil {
			return nil, err
		}
		return []string{
			filepath.Join(local, appName, "Data"),
			filepath.Join(roaming, appName, "Data"),
		}, nil
	case categoryConfig:
		roaming, err := knownFolderPath(knownFolderRoamingAppData)
		if err != nil {
			return nil, err
		}
		local, err := knownFolderPath(knownFolderLocalAppData)
		if err != nil {
			return nil, err
		}
		return []string{
			filepath.Join(roaming, appName, "Config"),
			filepath.Join(local, appName, "Config"),
		}, nil
	case categoryLog:
		local, err := knownFolderPath(knownFolderLocalAppData)
		if err != nil {
			return nil, err
		}
		return []string{filepath.Join(local, appName, "Logs")}, nil
	case categoryCache:
		local, err := knownFolderPath(knownFolderLocalAppData)
		if err != nil {
			return nil, err
		}
		return []string{filepath.Join(local, appName, "Cache")}, nil
	default:
		return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
	}
}

func platformSystemDirs(appName string, cat category) ([]string, error) {
	programData, err := knownFolderPath(knownFolderProgramData)
	if err != nil {
		return nil, err
	}

	switch cat {
	case categoryData:
		return []string{filepath.Join(programData, appName, "Data")}, nil
	case categoryConfig:
		return []string{filepath.Join(programData, appName, "Config")}, nil
	case categoryLog:
		return []string{filepath.Join(programData, appName, "Logs")}, nil
	case categoryCache:
		return []string{filepath.Join(programData, appName, "Cache")}, nil
	default:
		return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
	}
}
