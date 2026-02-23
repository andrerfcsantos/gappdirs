# gappdirs

[![Go Reference](https://pkg.go.dev/badge/github.com/andrerfcsantos/gappdirs.svg)](https://pkg.go.dev/github.com/andrerfcsantos/gappdirs)
[![Go Report Card](https://goreportcard.com/badge/github.com/andrerfcsantos/gappdirs)](https://goreportcard.com/report/github.com/andrerfcsantos/gappdirs)

`gappdirs` allows you to retrieve platform-specifc directories and files for application data, configuration, logs, and cache.
This module handles platform-specific directory conventions and provides a consistent API for applications to manage their files across different operating systems.

Currently supports Linux, macOS and Windows.
On Linux, the directory and file lookup honors the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir/0.8/) v0.8 behavior.

Application files can be classified into four categories based on their purpose:
  - **Data** - application data files that should be preserved across sessions.
  - **Config** - configuration files that are typically looked up at the startup of the application.
  - **Log** - files that record the state of application's execution and can be used for debugging, monitoring, or auditing purposes.
  - **Cache** - temporary files used by the application.

Directory and file lookup of the files and directories is organized by scope:
  - **System** - for system files that are relevant for all users.
    A system scope lookup will resolve to system directories only.
  - **User** - for files that are relevant to the current user.
    A user scope lookup resolves to user directories first, then system directories.
	This is the default scope most applications might want to use.
  - **Local** - for files that are in specific folders like the current directory, or a project-specific directories.
    Local scope resolves to local directories first, then user directories, then system directories.

## Prerequisites

- Go `1.24` or newer

## Install

```bash
go get github.com/andrerfcsantos/gappdirs
```

## Usage

### Global functions

Use the global functions of the module to retrieve directories and files.
This is the easiest way to interact with the library.

```go
package main

import (
	"fmt"
	appdirs "github.com/andrerfcsantos/gappdirs"
)

func main() {
	// User data directories
	fmt.Println(appdirs.UserDataDirs("myapp")) // returns the user data directories for "myapp" in precedence order
	fmt.Println(appdirs.UserDataDir("myapp")) // returns the user data directory with most precedence for "myapp"

	// System log directories
	fmt.Println(appdirs.SystemLogDirs("myapp")) // returns the system log directories for "myapp" in precedence order
	fmt.Println(appdirs.SystemLogDir("myapp")) // returns the system log directory with most precedence for "myapp"

	// Create a config file in the directory with the most precedence for the user scope
	created, configPath, err := appdirs.CreateUserConfigFile("myapp", "settings.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config file path:", configPath, "created:", created)

	// Find all existing config files names matching "settings.yaml" in the config dirs for user scope
	configFiles, err := appdirs.FindUserConfigFiles("myapp", "settings.yaml")
	if err != nil {
		fmt.Println("error finding config files:", err)
	}
	fmt.Println("config files:", configFiles)
}
```

### Resolver for multiple scoped operations

If you do operations on a scope for a single app name often, you can create a Resolver for that app name and scope, and re-use it for several operations in files and directories:


```go
package main

import (
	"fmt"
	appdirs "github.com/andrerfcsantos/gappdirs"
)

func main() {
	// The same example as above, but using a Resolver
	// Create a resolver for the user scope for the app "myapp"
	resolver := appdirs.NewUserResolver("myapp")

	// User data directories
	fmt.Println(resolver.DataDirs()) // returns the user data directories for "myapp" in precedence order
	fmt.Println(resolver.DataDir()) // returns the user data directory with most precedence for "myapp"

	// User log directories
	fmt.Println(resolver.LogDirs()) // returns the user log directories for "myapp" in precedence order
	fmt.Println(resolver.LogDir()) // returns the user log directory with most precedence for "myapp"

	// Create a config file in the directory with the most precedence for the user scope
	created, configPath, err := resolver.CreateConfigFile("settings.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("config file path:", configPath, "created:", created)

	// Find all existing config files names matching "settings.yaml" in the config dirs for user scope
	configFiles, err := resolver.FindConfigFiles("settings.yaml")
	if err != nil {
		fmt.Println("error finding config files:", err)
	}
	fmt.Println("config files:", configFiles)
}
```

### Local Scope

Local scope allows you to configure lookup on directories that might not be user or system-related.
For instance, before looking up user's and system's folder, you might want to check the user's local directory,
or a config folder inside the local directory:

```go
r := appdirs.NewLocalResolver("myapp",
	appdirs.WithLocalDirs(".", "./config"),
)
```

## Directory Resolution Matrix

Each cell shows lookup precedence in order: `Local -> User -> System`.

| Operating System | Data | Config | Log | Cache |
| --- | --- | --- | --- | --- |
| Linux | 1. Local: `<cwd>/.<appName>/data`<br>2. User: `$XDG_DATA_HOME/<appName>` (default: `~/.local/share/<appName>`)<br>3. System: `/var/lib/<appName>`, then each entry in `$XDG_DATA_DIRS/<appName>` (default: `/usr/local/share/<appName>`, `/usr/share/<appName>`) | 1. Local: `<cwd>/.<appName>/config`<br>2. User: `$XDG_CONFIG_HOME/<appName>` (default: `~/.config/<appName>`)<br>3. System: each entry in `$XDG_CONFIG_DIRS/<appName>` (default: `/etc/xdg/<appName>`) | 1. Local: `<cwd>/.<appName>/log`<br>2. User: `$XDG_STATE_HOME/<appName>/log` (default: `~/.local/state/<appName>/log`)<br>3. System: `/var/log/<appName>` | 1. Local: `<cwd>/.<appName>/cache`<br>2. User: `$XDG_CACHE_HOME/<appName>` (default: `~/.cache/<appName>`)<br>3. System: `/var/cache/<appName>` |
| macOS | 1. Local: `<cwd>/.<appName>/data`<br>2. User: `~/Library/Application Support/<appName>/data`<br>3. System: `/Library/Application Support/<appName>/data` | 1. Local: `<cwd>/.<appName>/config`<br>2. User: `~/Library/Application Support/<appName>/config`<br>3. System: `/Library/Application Support/<appName>/config` | 1. Local: `<cwd>/.<appName>/log`<br>2. User: `~/Library/Logs/<appName>`<br>3. System: `/Library/Logs/<appName>` | 1. Local: `<cwd>/.<appName>/cache`<br>2. User: `~/Library/Caches/<appName>`<br>3. System: `/Library/Caches/<appName>` |
| Windows | 1. Local: `<cwd>\\.<appName>\\data` (normalized form: `<cwd>/.<appName>/data`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Data`, `%APPDATA%\\<appName>\\Data`<br>3. System: `%ProgramData%\\<appName>\\Data` | 1. Local: `<cwd>\\.<appName>\\config` (normalized form: `<cwd>/.<appName>/config`)<br>2. User: `%APPDATA%\\<appName>\\Config`, `%LOCALAPPDATA%\\<appName>\\Config`<br>3. System: `%ProgramData%\\<appName>\\Config` | 1. Local: `<cwd>\\.<appName>\\log` (normalized form: `<cwd>/.<appName>/log`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Logs`<br>3. System: `%ProgramData%\\<appName>\\Logs` | 1. Local: `<cwd>\\.<appName>\\cache` (normalized form: `<cwd>/.<appName>/cache`)<br>2. User: `%LOCALAPPDATA%\\<appName>\\Cache`<br>3. System: `%ProgramData%\\<appName>\\Cache` |
