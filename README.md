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
	r := appdirs.NewResolver("myapp", appdirs.WithScope(appdirs.ScopeUser))

	dataDirs := r.DataDirs()
	fmt.Println("data dirs:", dataDirs)

	configDirs := r.ConfigDirs()
	fmt.Println("config dirs:", configDirs)

	logDirs := r.LogDirs()
	fmt.Println("log dirs:", logDirs)

	cacheDirs := r.CacheDirs()
	fmt.Println("cache dirs:", cacheDirs)

	cacheDir, err := r.EnsureCacheDir(appdirs.WithEnsureDirPerm(0o700))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ensured cache dir:", cacheDir)

	configFile, err := r.FindConfigFile("settings.yaml")
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

- `NewResolver(appName, opts...) Resolver` uses normal option-driven scope behavior.
- `NewUserResolver(appName, opts...) Resolver` always uses `ScopeUser`.
- `NewSystemResolver(appName, opts...) Resolver` always uses `ScopeSystem`.
- `NewLocalResolver(appName, opts...) Resolver` always uses `ScopeLocal`.

For fixed-scope constructors, scope is enforced by the constructor even if `WithScope(...)` is provided in options.

## Option Behavior

Options never return errors and invalid values fall back to defaults.

- `WithScope(scope)` uses `ScopeUser` when `scope` is invalid.
- `WithWorkingDir(dir)` trims spaces and ignores the option when `dir` is empty.
- `WithDefaultDirPerm(perm)` falls back to `0o755` when `perm` is zero or contains non-permission bits.
- `nil` options are ignored.

For local scope, working directory resolution is deferred to local-scope operations. This keeps constructor paths non-failing.

Ensure call options:

- `WithEnsureDirPerm(perm)` sets the permission for a single `Ensure*Dir` call.
- Invalid values (`0` or values containing mode-type bits) are ignored and fallback to resolver default permission.
- `nil` ensure options are ignored.

## Top-Level Scoped Helpers

You can also use top-level scoped helpers without creating a resolver instance.
Each helper takes only `appName` and uses the default resolver behavior.

- Directory lists: `LocalDataDirs(appName)`, `UserConfigDirs(appName)`, `SystemCacheDirs(appName)`
- Most relevant directory: `LocalDataDir(appName)`, `UserConfigDir(appName)`, `SystemCacheDir(appName)`
- Ensure directory: `EnsureLocalDataDir(appName)`, `EnsureUserConfigDir(appName)`, `EnsureSystemCacheDir(appName)`
- Ensure directory with options: `EnsureLocalDataDir(appName, WithEnsureDirPerm(0o700))`
- Find matching files: `FindLocalDataFiles(appName, filename)`, `FindUserConfigFiles(appName, filename)`, `FindSystemCacheFiles(appName, filename)`
- Most relevant file: `FindLocalDataFile(appName, filename)`, `FindUserConfigFile(appName, filename)`, `FindSystemCacheFile(appName, filename)`

`Dir`/`Dirs` helpers do not return errors. `Ensure*`, `Find*Files`, and `Find*File` helpers still return errors.

Naming pattern:
- Non-`Ensure`/`Find`: `Scope + Category + ...` (for example, `LocalDataDir`)
- `Ensure` methods: `Ensure + Scope + Category + ...` (for example, `EnsureLocalDataDir`)
- `Find` methods: `Find + Scope + Category + ...` (for example, `FindLocalDataFiles`)

## App Name Sanitization

App names are sanitized during constructor creation:

- Trim leading/trailing whitespace.
- Convert spaces/whitespace to `_`.
- Convert `/`, `\\`, `:`, `*`, `?`, `"`, `<`, `>`, and `|` to `_`.
- Keep all other characters unchanged (including letter casing and repeated underscores).

If the sanitized name is empty, it defaults to `unnamed_app`.

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
