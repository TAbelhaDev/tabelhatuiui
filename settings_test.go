package tuiui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func testRegistry(t *testing.T) *KeyRegistry {
	t.Helper()
	r := NewKeyRegistry(filepath.Join(t.TempDir(), "keybindings.json"))
	r.RegisterMany(
		Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
		Action{ID: "refresh", Help: "atualizar", Keys: []string{"r"}},
		Action{ID: "nav", Help: "mover", Keys: []string{"j", "k"}, Label: "j/k"},
	)
	return r
}

func TestKeyRegistryResolve(t *testing.T) {
	r := testRegistry(t)
	if b := r.Resolve("quit"); !keyMatch(b, "q") {
		t.Fatalf("Resolve(quit) keys = %v, want q", b.Keys())
	}
	if got := r.Resolve("quit").Help().Key; got != "q/ctrl+c" {
		t.Fatalf("Resolve(quit) label = %q, want q/ctrl+c", got)
	}
	// Custom label wins.
	if got := r.Resolve("nav").Help().Key; got != "j/k" {
		t.Fatalf("Resolve(nav) label = %q, want j/k", got)
	}
	// Unknown id → disabled, no keys.
	if r.Resolve("nope").Enabled() {
		t.Fatal("Resolve(unknown) should be disabled")
	}
}

func TestKeyRegistrySetResetPersist(t *testing.T) {
	r := testRegistry(t)
	if err := r.Set("refresh", "f"); err != nil {
		t.Fatal(err)
	}
	if !keyMatch(r.Resolve("refresh"), "f") {
		t.Fatalf("Resolve(refresh) = %v, want f", r.Resolve("refresh").Keys())
	}
	if !r.IsCustom("refresh") {
		t.Fatal("IsCustom(refresh) should be true")
	}

	// Reload from disk: override must survive.
	r2 := NewKeyRegistry(r.path)
	r2.RegisterMany(
		Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
		Action{ID: "refresh", Help: "atualizar", Keys: []string{"r"}},
	)
	if err := r2.Load(); err != nil {
		t.Fatal(err)
	}
	if !keyMatch(r2.Resolve("refresh"), "f") {
		t.Fatalf("reloaded Resolve(refresh) = %v, want f", r2.Resolve("refresh").Keys())
	}

	if err := r.Reset("refresh"); err != nil {
		t.Fatal(err)
	}
	if !keyMatch(r.Resolve("refresh"), "r") {
		t.Fatalf("after Reset, Resolve(refresh) = %v, want r", r.Resolve("refresh").Keys())
	}
	if _, err := os.Stat(r.path); !os.IsNotExist(err) {
		t.Fatalf("config file should be deleted after reset, got %v", err)
	}
}

func TestKeyRegistryConflict(t *testing.T) {
	r := testRegistry(t)
	err := r.Set("refresh", "q") // q is used by quit
	if err == nil {
		t.Fatal("Set should fail on a conflicting key")
	}
	if r.IsCustom("refresh") {
		t.Fatal("failed Set must not create an override")
	}
	// Multiple keys with one conflict also fail.
	if err := r.Set("refresh", "f", "j"); err == nil {
		t.Fatal("Set with a conflicting key among several should fail")
	}
}

func TestKeyRegistryResetAll(t *testing.T) {
	r := testRegistry(t)
	_ = r.Set("refresh", "f")
	_ = r.Set("nav", "n")
	if n, _ := r.ResetAll(); n != 2 {
		t.Fatalf("ResetAll cleared %d, want 2", n)
	}
	if r.CustomCount() != 0 {
		t.Fatal("CustomCount should be 0 after ResetAll")
	}
	if keyMatch(r.Resolve("nav"), "n") {
		t.Fatal("nav should be back to defaults")
	}
}

func TestKeyRegistryLoadMissingFile(t *testing.T) {
	r := NewKeyRegistry(filepath.Join(t.TempDir(), "nope.json"))
	if err := r.Load(); err != nil {
		t.Fatalf("Load on missing file should be nil, got %v", err)
	}
}

func TestSettingsModalCapture(t *testing.T) {
	r := testRegistry(t)
	m := NewSettingsModal(r)
	m.SetSize(80, 24)
	m.Open()

	if !m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}) {
		t.Fatal("open modal should consume messages")
	}
	// enter on action 1 (refresh)
	if !m.Update(tea.KeyMsg{Type: tea.KeyEnter}) {
		t.Fatal("enter should be consumed")
	}
	if !m.capturing {
		t.Fatal("should be capturing after enter")
	}
	// press f → rebinds refresh
	if !m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("f")}) {
		t.Fatal("capture key should be consumed")
	}
	if m.capturing {
		t.Fatal("capture should end after a key")
	}
	if !keyMatch(r.Resolve("refresh"), "f") {
		t.Fatalf("refresh should be rebound to f, got %v", r.Resolve("refresh").Keys())
	}

	// esc cancels capture without rebinding
	// (navigate to quit, enter, then esc)
	for i := 0; i < 5; i++ {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("k")})
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if !m.capturing {
		t.Fatal("should be capturing")
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.capturing {
		t.Fatal("esc should cancel capture")
	}
	if !keyMatch(r.Resolve("quit"), "q") {
		t.Fatal("quit should be unchanged after esc")
	}
}

func TestSettingsModalResetKey(t *testing.T) {
	r := testRegistry(t)
	_ = r.Set("refresh", "f")
	m := NewSettingsModal(r)
	m.SetSize(80, 24)
	m.Open()
	// cursor starts at 0 (quit); move to refresh (index 1)
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")})
	if !m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}) {
		t.Fatal("r should be consumed")
	}
	if keyMatch(r.Resolve("refresh"), "f") {
		t.Fatal("refresh should be reset to default r")
	}
}

func TestSettingsModalView(t *testing.T) {
	r := testRegistry(t)
	m := NewSettingsModal(r)
	m.SetSize(80, 24)
	if out := m.View(ResolveTheme("", "mauve")); out != "" {
		t.Fatalf("closed modal should render empty, got %q", out)
	}
	m.Open()
	out := m.View(ResolveTheme("", "mauve"))
	if out == "" {
		t.Fatal("open modal should render")
	}
	for _, want := range []string{"sair", "atualizar", "fechar"} {
		if !strings.Contains(out, want) {
			t.Errorf("View() missing %q", want)
		}
	}
}

func keyMatch(b Binding, key string) bool {
	for _, k := range b.Keys() {
		if k == key {
			return true
		}
	}
	return false
}
