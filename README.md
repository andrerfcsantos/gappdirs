# gappdirs

[![Go Reference](https://pkg.go.dev/badge/github.com/andrerfcsantos/gappdirs.svg)](https://pkg.go.dev/github.com/andrerfcsantos/gappdirs)

`gappdirs` is a Go library for fetching platform-specific application directories, such as:

- **Data** - for user-generated content, databases, and other persistent files that are not configuration.
- **Config** - for configuration files, settings, and other user-editable resources.
- **Log** - for log files, crash dumps, and other diagnostic output.
- **Cache** - for temporary files, caches, and other non-essential data that can be safely deleted.

It supports windows, macOS, and Linux, following the conventions of each platform.
On Linux, the resolver honors the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir/0.8/) v0.8 behavior.

## Scopes

The directories returned are an ordered list of paths from most to least specific, according to the scope of the resolver:

- **Local** - directories specific to the current working directory, typically hidden subdirectories named `.<appName>`.
- **User** - directories specific to the current user, typically located in the user's home directory or application data folders.
- **System** - directories shared across all users, typically located in system

The resolver will return the directories for its scope, and all less specific scopes.
For example, a resolver with `ScopeUser` will return both user and system directories, in this order, but not local directories.

## Prerequisites

- Go `1.24` or newer

## Install

```bash
go get github.com/andrerfcsantos/gappdirs
```

## Usage Example

```go
package main

import (
	"errors"
	"fmt"
	"log"

	appdirs "github.com/andrerfcsantos/gappdirs"
)

func main() {
	r, err := appdirs.NewResolver("myapp", appdirs.WithScope(appdirs.ScopeUser))
	if err != nil {
		log.Fatal(err)
	}

	dataDirs, err := r.DataDirs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("data dirs:", dataDirs)

	configDirs, err := r.ConfigDirs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config dirs:", configDirs)

	logDirs, err := r.LogDirs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("log dirs:", logDirs)

	cacheDirs, err := r.CacheDirs()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("cache dirs:", cacheDirs)

	cacheDir, err := r.EnsureCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ensured cache dir:", cacheDir)

	configFile, err := r.ConfigFile("settings.yaml")
	if errors.Is(err, appdirs.ErrNotFound) {
		fmt.Println("settings.yaml not found in config search path")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("settings file:", configFile)
	}
}
```

## Constructors

- `NewResolver(appName, opts...)` returns a `Resolver` and uses normal option-driven scope behavior.
- `NewUserResolver(appName, opts...)` returns a `Resolver` and always uses `ScopeUser`.
- `NewSystemResolver(appName, opts...)` returns a `Resolver` and always uses `ScopeSystem`.
- `NewLocalResolver(appName, opts...)` returns a `Resolver` and always uses `ScopeLocal`.

For fixed-scope constructors, scope is enforced by the constructor even if `WithScope(...)` is provided in options.

## Top-Level Scoped Helpers

You can also use top-level scoped helpers without creating a resolver instance.
Each helper takes only `appName` and uses the default resolver behavior.

- Directory lists: `LocalDataDirs(appName)`, `UserConfigDirs(appName)`, `SystemCacheDirs(appName)`
- Most relevant directory: `LocalDataDir(appName)`, `UserConfigDir(appName)`, `SystemCacheDir(appName)`
- Ensure directory: `EnsureLocalDataDir(appName)`, `EnsureUserConfigDir(appName)`, `EnsureSystemCacheDir(appName)`
- Ensure directory with permissions: `EnsureLocalDataDirWithPerm(appName, perm)`
- Find matching file directories: `FindLocalDataFileDirs(appName, filename)`, `FindUserConfigFileDirs(appName, filename)`, `FindSystemCacheFileDirs(appName, filename)`
- Most relevant file: `LocalDataFile(appName, filename)`, `UserConfigFile(appName, filename)`, `SystemCacheFile(appName, filename)`

Naming pattern:
- Non-`Ensure`/`Find`: `Scope + Category + ...` (for example, `LocalDataDir`)
- `Ensure` methods: `Ensure + Scope + Category + ...` (for example, `EnsureLocalDataDir`)
- `Find` methods: `Find + Scope + Category + ...` (for example, `FindLocalDataFileDirs`)

## App Name Sanitization

App names are sanitized during constructor creation:

- Trim leading/trailing whitespace.
- Convert to lowercase.
- Convert spaces/whitespace to `_`.
- Convert `/`, `\\`, `:`, `*`, `?`, `"`, `<`, `>`, and `|` to `_`.
- Collapse repeated underscores and trim leading/trailing underscores.
- Special cases: `"."` becomes `"_"`, and `".."` becomes `"__"`.

If the sanitized name is empty, constructor creation returns `gappdirs: app name is required`.

## Directory Resolution Matrix

Each cell shows lookup precedence in order: `Local -> User -> System`.

| Operating System | Data | Config | Log | Cache |
| --- | --- | --- | --- | --- |
| Linux | 1. Local: `<cwd>/.<appName>/data`<br>2. User: `$XDG_DATA_HOME/<appName>` (default: `~/.local/share/<appName>`)<br>3. System: `/var/lib/<appName>`, then each entry in `$XDG_DATA_DIRS/<appName>` (default: `/usr/local/share/<appName>`, `/usr/share/<appName>`) | 1. Local: `<cwd>/.<appName>/config`<br>2. User: `$XDG_CONFIG_HOME/<appName>` (default: `~/.config/<appName>`)<br>3. System: each entry in `$XDG_CONFIG_DIRS/<appName>` (default: `/etc/xdg/<appName>`) | 1. Local: `<cwd>/.<appName>/log`<br>2. User: `$XDG_STATE_HOME/<appName>/log` (default: `~/.local/state/<appName>/log`)<br>3. System: `/var/log/<appName>` | 1. Local: `<cwd>/.<appName>/cache`<br>2. User: `$XDG_CACHE_HOME/<appName>` (default: `~/.cache/<appName>`)<br>3. System: `/var/cache/<appName>` |
| macOS | 1. Local: `<cwd>/.<appName>/data`<br>2. User: `~/Library/Application Support/<appName>/data`<br>3. System: `/Library/Application Support/<appName>/data` | 1. Local: `<cwd>/.<appName>/config`<br>2. User: `~/Library/Application Support/<appName>/config`<br>3. System: `/Library/Application Support/<appName>/config` | 1. Local: `<cwd>/.<appName>/log`<br>2. User: `~/Library/Logs/<appName>`<br>3. System: `/Library/Logs/<appName>` | 1. Local: `<cwd>/.<appName>/cache`<br>2. User: `~/Library/Caches/<appName>`<br>3. System: `/Library/Caches/<appName>` |
| Windows | 1. Local: `<cwd>\\.<appName>\\data` (normalized form: `<cwd>/.<appName>/data`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Data`, `%APPDATA%\\<appName>\\Data`<br>3. System: `%ProgramData%\\<appName>\\Data` | 1. Local: `<cwd>\\.<appName>\\config` (normalized form: `<cwd>/.<appName>/config`)<br>2. User: `%APPDATA%\\<appName>\\Config`, `%LOCALAPPDATA%\\<appName>\\Config`<br>3. System: `%ProgramData%\\<appName>\\Config` | 1. Local: `<cwd>\\.<appName>\\log` (normalized form: `<cwd>/.<appName>/log`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Logs`<br>3. System: `%ProgramData%\\<appName>\\Logs` | 1. Local: `<cwd>\\.<appName>\\cache` (normalized form: `<cwd>/.<appName>/cache`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Cache`<br>3. System: `%ProgramData%\\<appName>\\Cache` |

## Releases

Releases are automated with GoReleaser through GitHub Actions.

- Push a semver tag (for example, `v0.1.0` or `v0.2.0-rc.1`) that points to a commit on `main`.
- The workflow in `.github/workflows/release.yml` runs tests and publishes the release.
- Assets include a source archive and `checksums.txt`.

Local verification:

```bash
goreleaser check
goreleaser release --snapshot --clean --skip=publish
```
