//go:build !linux && !darwin && !windows

package gappdirs

import (
	"fmt"
	"runtime"
)

func platformUserDirs(appName string, cat category) ([]string, error) {
	return nil, fmt.Errorf("gappdirs: unsupported operating system %q", runtime.GOOS)
}

func platformSystemDirs(appName string, cat category) ([]string, error) {
	return nil, fmt.Errorf("gappdirs: unsupported operating system %q", runtime.GOOS)
}
