package tuiui

import (
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
)

// Config holds an app's settings, merged from compiled-in defaults and an
// optional TOML file. It mirrors KeyRegistry: the file on disk is the source
// of truth, and Reload() picks up external edits without restarting.
//
// T must be a struct of plain values. Pointer, map and slice fields are shared
// with the defaults value until the file overrides them, so callers must not
// mutate what Get() returns.
type Config[T any] struct {
	path     string
	defaults T
	current  T
}

// NewConfig builds a Config for path, using defaults for every key the file
// leaves out. It does not touch the filesystem; call Load first.
func NewConfig[T any](path string, defaults T) *Config[T] {
	return &Config[T]{path: path, defaults: defaults, current: defaults}
}

// Path is the config file this Config reads, useful for error messages.
func (c *Config[T]) Path() string { return c.path }

// Load reads the config file once. A missing file is not an error — the app
// runs on pure defaults. See Reload for the config-file-first flow.
func (c *Config[T]) Load() error {
	_, err := c.Reload()
	return err
}

// Reload re-reads the config file, making the on-disk config.toml the source
// of truth again — call it from the app's reload key so an external edit is
// picked up without restarting. It reports whether the effective config
// changed (a deleted file counts as a change back to pure defaults).
//
// A malformed file returns an error and leaves the previous config in place,
// so a typo mid-edit never drops a running app back to defaults.
func (c *Config[T]) Reload() (bool, error) {
	prev := c.current

	// Struct copy: keys absent from the file keep their default value, which
	// is what makes the decode below a merge rather than a replacement.
	next := c.defaults

	data, err := os.ReadFile(c.path)
	if err != nil {
		if os.IsNotExist(err) {
			c.current = next
			return !reflect.DeepEqual(prev, next), nil
		}
		return false, err
	}

	// The decoder only assigns fields whose keys appear in the file, so
	// decoding over `next` layers the overrides onto the defaults.
	if _, err := toml.Decode(string(data), &next); err != nil {
		return false, err
	}

	c.current = next
	return !reflect.DeepEqual(prev, next), nil
}

// Get returns the current effective config: defaults with the file's
// overrides applied.
func (c *Config[T]) Get() T { return c.current }
