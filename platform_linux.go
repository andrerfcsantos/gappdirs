//go:build linux

package gappdirs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	linuxDefaultDataHome   = ".local/share"
	linuxDefaultConfigHome = ".config"
	linuxDefaultStateHome  = ".local/state"
	linuxDefaultCacheHome  = ".cache"

	linuxDefaultDataDirs   = "/usr/local/share:/usr/share"
	linuxDefaultConfigDirs = "/etc/xdg"
)

func platformUserDirs(appName string, cat category) ([]string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "~"
	}

	switch cat {
	case categoryData:
		base := linuxEnvDirOrDefault("XDG_DATA_HOME", filepath.Join(homeDir, linuxDefaultDataHome))
		return []string{filepath.Join(base, appName)}, err
	case categoryConfig:
		base := linuxEnvDirOrDefault("XDG_CONFIG_HOME", filepath.Join(homeDir, linuxDefaultConfigHome))
		return []string{filepath.Join(base, appName)}, err
	case categoryLog:
		base := linuxEnvDirOrDefault("XDG_STATE_HOME", filepath.Join(homeDir, linuxDefaultStateHome))
		return []string{filepath.Join(base, appName, "log")}, err
	case categoryCache:
		base := linuxEnvDirOrDefault("XDG_CACHE_HOME", filepath.Join(homeDir, linuxDefaultCacheHome))
		return []string{filepath.Join(base, appName)}, err
	}

	return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
}

func platformSystemDirs(appName string, cat category) ([]string, error) {
	switch cat {
	case categoryData:
		dirs := []string{filepath.Join("/var/lib", appName)}
		for _, base := range linuxEnvDirListOrDefault("XDG_DATA_DIRS", linuxDefaultDataDirs) {
			dirs = append(dirs, filepath.Join(base, appName))
		}
		return dirs, nil
	case categoryConfig:
		bases := linuxEnvDirListOrDefault("XDG_CONFIG_DIRS", linuxDefaultConfigDirs)
		dirs := make([]string, 0, len(bases))
		for _, base := range bases {
			dirs = append(dirs, filepath.Join(base, appName))
		}
		return dirs, nil
	case categoryLog:
		return []string{filepath.Join("/var/log", appName)}, nil
	case categoryCache:
		return []string{filepath.Join("/var/cache", appName)}, nil
	default:
		return nil, fmt.Errorf("gappdirs: unsupported category %d", cat)
	}
}

func linuxEnvDirOrDefault(envKey, fallback string) string {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value != "" && filepath.IsAbs(value) {
		return filepath.Clean(value)
	}
	return filepath.Clean(fallback)
}

func linuxEnvDirListOrDefault(envKey, fallback string) []string {
	value := strings.TrimSpace(os.Getenv(envKey))
	if value == "" {
		value = fallback
	}
	dirs := splitAbsolutePathList(value)
	if len(dirs) > 0 {
		return dirs
	}
	return splitAbsolutePathList(fallback)
}

func splitAbsolutePathList(pathList string) []string {
	rawItems := filepath.SplitList(pathList)
	out := make([]string, 0, len(rawItems))
	for _, item := range rawItems {
		item = strings.TrimSpace(item)
		if item == "" || !filepath.IsAbs(item) {
			continue
		}
		out = append(out, filepath.Clean(item))
	}
	return out
}
