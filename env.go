package tuiui

import (
	"os"
	"path/filepath"
	"strings"
)

// HomeDir is the current user's home directory, the base most config/data
// paths are anchored to. Falls back to $HOME when UserHomeDir fails.
func HomeDir() string {
	if dir, err := os.UserHomeDir(); err == nil {
		return dir
	}
	return os.Getenv("HOME")
}

// ConfigDir is the user config base (os.UserConfigDir()), falling back to
// ~/.config — the parent of every project's own config directory.
func ConfigDir() string {
	if dir, err := os.UserConfigDir(); err == nil {
		return dir
	}
	return filepath.Join(HomeDir(), ".config")
}

// ConfigPath is ~/.config/<app>/<file>, the canonical location of an app's
// own config files (config.toml, keybindings.json). It honors XDG_CONFIG_HOME
// through ConfigDir, so callers never need to resolve the base themselves.
func ConfigPath(app, file string) string {
	return filepath.Join(ConfigDir(), app, file)
}

// EnvOr returns the value of env var `key`, or `fallback` when it's unset or
// empty.
func EnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// ExpandHome expands a leading "~" or "~/" into the user's home directory,
// leaving any other path untouched.
func ExpandHome(path string) string {
	if path == "~" {
		return HomeDir()
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(HomeDir(), path[2:])
	}
	return path
}
