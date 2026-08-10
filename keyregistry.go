package tuiui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/charmbracelet/bubbles/key"
)

// Action is one keybinding in a KeyRegistry: a stable ID (used as the key in
// the persisted config file and in Resolve calls), the help description shown
// in the footer/help modal, the default keys, and an optional display label
// (defaults to the keys joined with "/"). Label can be a compact alias — e.g.
// Keys {"j","k","up","down"} with Label "j/k".
type Action struct {
	ID    string
	Help  string
	Keys  []string
	Label string
}

// KeyRegistry is the single source of truth for an app's keybindings:
// defaults registered in code, optional per-action overrides persisted to a
// JSON file, and live resolution so dispatch/footer/help/settings always see
// the same effective binding. Create one, Register every action, Load(), then
// Resolve in Update and feed Bindings() to Footer/HelpModal.
//
// The config file is the primary customization flow (config-file-first): the
// user edits keybindings.json and the app calls Reload() on a reload key (or
// per View) to pick the edit up without restarting. The persisted format is a
// flat map of action ID to the custom keys; the file is written only when at
// least one override exists (all-resets delete it), and is removed when the
// last override is cleared.
type KeyRegistry struct {
	path   string
	order  []string
	byID   map[string]*Action
	custom map[string][]string // id -> custom keys (from the config file)
}

// NewKeyRegistry creates an empty registry that persists overrides to path.
func NewKeyRegistry(path string) *KeyRegistry {
	return &KeyRegistry{
		path:   path,
		byID:   map[string]*Action{},
		custom: map[string][]string{},
	}
}

// Register adds (or replaces) an action with its default keys and help.
// Registration order is kept and used by Bindings()/Actions().
func (r *KeyRegistry) Register(a Action) *KeyRegistry {
	if _, ok := r.byID[a.ID]; !ok {
		r.order = append(r.order, a.ID)
	}
	r.byID[a.ID] = &a
	return r
}

// RegisterMany is a convenience for registering a slice of actions in order.
func (r *KeyRegistry) RegisterMany(actions ...Action) *KeyRegistry {
	for _, a := range actions {
		r.Register(a)
	}
	return r
}

// Resolve returns the effective binding for an action — custom keys when an
// override exists, otherwise the defaults. Unknown IDs resolve to a disabled
// empty binding.
func (r *KeyRegistry) Resolve(id string) key.Binding {
	a, ok := r.byID[id]
	if !ok {
		return key.NewBinding(key.WithKeys())
	}
	keys := a.Keys
	if c, ok := r.custom[id]; ok {
		keys = c
	}
	label := a.Label
	if label == "" {
		label = strings.Join(keys, "/")
	}
	return key.NewBinding(key.WithKeys(keys...), key.WithHelp(label, a.Help))
}

// Bindings returns the effective binding for every registered action, in
// registration order — feed this to Footer/HelpModal.
func (r *KeyRegistry) Bindings() []key.Binding {
	out := make([]key.Binding, 0, len(r.order))
	for _, id := range r.order {
		out = append(out, r.Resolve(id))
	}
	return out
}

// IsCustom reports whether an action currently has a custom override.
func (r *KeyRegistry) IsCustom(id string) bool {
	_, ok := r.custom[id]
	return ok
}

// CustomCount is the number of actions with a custom override.
func (r *KeyRegistry) CustomCount() int {
	return len(r.custom)
}

// Set overrides the keys of an action and immediately persists the change.
// It fails (without saving) when a key is already bound to another action.
func (r *KeyRegistry) Set(id string, keys ...string) error {
	if _, ok := r.byID[id]; !ok {
		return fmt.Errorf("ação desconhecida %q", id)
	}
	var clean []string
	for _, k := range keys {
		k = strings.TrimSpace(k)
		if k == "" {
			continue
		}
		clean = append(clean, k)
	}
	if len(clean) == 0 {
		return fmt.Errorf("precisa de pelo menos uma tecla")
	}
	// Conflict check against every other action's current binding.
	used := map[string]string{} // key -> owner action id
	for _, oid := range r.order {
		if oid == id {
			continue
		}
		for _, k := range r.Resolve(oid).Keys() {
			used[k] = oid
		}
	}
	for _, k := range clean {
		if owner, ok := used[k]; ok {
			return fmt.Errorf("a tecla %q já é usada por %q", k, owner)
		}
	}
	r.custom[id] = clean
	return r.Save()
}

// Reset restores an action to its default keys and persists the change.
func (r *KeyRegistry) Reset(id string) error {
	delete(r.custom, id)
	return r.Save()
}

// ResetAll restores every action to its defaults (deleting the config file)
// and returns the number of overrides that were cleared.
func (r *KeyRegistry) ResetAll() (int, error) {
	n := len(r.custom)
	r.custom = map[string][]string{}
	return n, r.Save()
}

// Actions returns every registered action with its CURRENT (effective) keys,
// in registration order, for rendering the settings list.
func (r *KeyRegistry) Actions() []Action {
	out := make([]Action, 0, len(r.order))
	for _, id := range r.order {
		a := *r.byID[id]
		if c, ok := r.custom[id]; ok {
			a.Keys = c
		}
		if a.Label == "" {
			a.Label = strings.Join(a.Keys, "/")
		}
		out = append(out, a)
	}
	return out
}

type registryFile struct {
	Bindings map[string][]string `json:"bindings"`
}

// Load reads custom overrides from the config file once. A missing file is
// not an error (no overrides yet). See Reload for the config-file-first flow.
func (r *KeyRegistry) Load() error {
	_, err := r.Reload()
	return err
}

// Reload re-reads the config file, making the on-disk keybindings.json the
// source of truth again — call it from the app's reload key (or once per
// View) so an external edit to the file is picked up without restarting. It
// reports whether the overrides changed on disk since the last load/save (a
// deleted file counts as a change back to pure defaults).
func (r *KeyRegistry) Reload() (bool, error) {
	prev := r.custom
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			r.custom = map[string][]string{}
			return len(prev) > 0, nil
		}
		return false, err
	}
	var f registryFile
	if err := json.Unmarshal(data, &f); err != nil {
		return false, err
	}
	if f.Bindings == nil {
		f.Bindings = map[string][]string{}
	}
	changed := !reflect.DeepEqual(prev, f.Bindings)
	r.custom = f.Bindings
	return changed, nil
}

// Save persists the current overrides. With nothing overridden it deletes
// the config file so the app reads as pure defaults.
func (r *KeyRegistry) Save() error {
	if len(r.custom) == 0 {
		if err := os.Remove(r.path); err != nil && !os.IsNotExist(err) {
			return err
		}
		return nil
	}
	data, err := json.MarshalIndent(registryFile{Bindings: r.custom}, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(r.path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(r.path, data, 0o600)
}
