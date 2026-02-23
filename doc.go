// Package gappdirs allows you to retrieve platform-specifc directories and files for application data, configuration, logs, and cache.
// This module handles platform-specific directory conventions and provides a consistent API for applications to manage their files across different operating systems.
// On Linux, the directory and file lookup honors the XDG Base Directory Specification v0.8 behavior.
//
// Application files can be classified into four categories based on their purpose:
//   - Data - application data files that should be preserved across sessions.
//   - Config - configuration files that are typically looked up at the startup of the application.
//   - Log - files that record the state of application's execution and can be used for debugging, monitoring, or auditing purposes.
//   - Cache - temporary files used by the application.
//
// Directory and file lookup of the files and directories is organized by scope:
//   - System - for system files that are relevant for all users.
//     A system scope lookup will resolve to system directories only.
//   - User - for files that are relevant to the current user.
//     A user scope lookup resolves to user directories first, then system directories.
//   - Local - for files that are in specific folders like the current directory, or a project-specific directories.
//     Local scope resolves to local directories first, then user directories, then system directories.
//
// Use the global functions of the module to retrieve directories and files:
//
//	package main
//
//	import (
//		"fmt"
//		appdirs "github.com/andrerfcsantos/gappdirs"
//	)
//
//	func main() {
//		// User data directories
//		fmt.Println(appdirs.UserDataDirs("myapp")) // returns the user data directories for "myapp" in precedence order
//		fmt.Println(appdirs.UserDataDir("myapp")) // returns the user data directory with most precedence for "myapp"
//
//		// System log directories
//		fmt.Println(appdirs.SystemLogDirs("myapp")) // returns the system log directories for "myapp" in precedence order
//		fmt.Println(appdirs.SystemLogDir("myapp")) // returns the system log directory with most precedence for "myapp"
//
//		// Create a config file in the directory with the most precedence for the user scope
//		created, configPath, err := appdirs.CreateUserConfigFile("myapp", "settings.yaml")
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println("config file path:", configPath, "created:", created)
//
//		// Find all existing config files names matching "settings.yaml" in the config dirs for user scope
//		configFiles, err := appdirs.FindUserConfigFiles("myapp", "settings.yaml")
//		if err != nil {
//			fmt.Println("error finding config files:", err)
//		}
//		fmt.Println("config files:", configFiles)
//	}
//
// If you do operations on a scope for a single app name often, you can create a Resolver for that app name and scope, and use it for operations in files and directories:
//
//	package main
//
//	import (
//		"fmt"
//		appdirs "github.com/andrerfcsantos/gappdirs"
//	)
//
//	func main() {
//		// The same example as above, but using a Resolver
//
//		// Create a resolver for the user scope for the app "myapp"
//		resolver := appdirs.NewUserResolver("myapp")
//
//		// User data directories
//		fmt.Println(resolver.DataDirs()) // returns the user data directories for "myapp" in precedence order
//		fmt.Println(resolver.DataDir()) // returns the user data directory with most precedence for "myapp"
//
//		// User log directories
//		fmt.Println(resolver.LogDirs()) // returns the user log directories for "myapp" in precedence order
//		fmt.Println(resolver.LogDir()) // returns the user log directory with most precedence for "myapp"
//
//		// Create a config file in the directory with the most precedence for the user scope
//		created, configPath, err := resolver.CreateConfigFile("settings.yaml")
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println("config file path:", configPath, "created:", created)
//
//		// Find all existing config files names matching "settings.yaml" in the config dirs for user scope
//		configFiles, err := resolver.FindConfigFiles("settings.yaml")
//		if err != nil {
//			fmt.Println("error finding config files:", err)
//		}
//		fmt.Println("config files:", configFiles)
//	}
package gappdirs
