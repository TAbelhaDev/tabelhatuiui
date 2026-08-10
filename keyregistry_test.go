package tuiui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKeyRegistryReload(t *testing.T) {
	r := testRegistry(t)
	if err := r.Load(); err != nil {
		t.Fatal(err)
	}
	if err := r.Set("refresh", "f"); err != nil {
		t.Fatal(err)
	}

	// External edit to the file: Reload must pick it up and report a change.
	data := []byte(`{"bindings": {"refresh": ["x"], "quit": ["ctrl+q"]}}`)
	if err := os.WriteFile(r.path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	changed, err := r.Reload()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("Reload() should report a change after an external edit")
	}
	if !keyMatch(r.Resolve("refresh"), "x") {
		t.Fatalf("after Reload, refresh = %v, want x", r.Resolve("refresh").Keys())
	}
	if !keyMatch(r.Resolve("quit"), "ctrl+q") {
		t.Fatalf("after Reload, quit = %v, want ctrl+q", r.Resolve("quit").Keys())
	}

	// Unchanged file → no change reported.
	if changed, err = r.Reload(); err != nil || changed {
		t.Fatalf("Reload() on unchanged file = changed:%v, err:%v; want false,nil", changed, err)
	}
}

func TestKeyRegistryReloadDeletedFile(t *testing.T) {
	r := testRegistry(t)
	if err := r.Set("refresh", "f"); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(r.path); err != nil {
		t.Fatal(err)
	}
	changed, err := r.Reload()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("Reload() should report a change when the file was deleted")
	}
	if r.CustomCount() != 0 {
		t.Fatalf("CustomCount = %d, want 0 after deleting the file", r.CustomCount())
	}
	if keyMatch(r.Resolve("refresh"), "f") {
		t.Fatal("refresh should be back to default r after the file was deleted")
	}
}

func TestFooterGroupedLabelFromRegistry(t *testing.T) {
	reg := NewKeyRegistry(filepath.Join(t.TempDir(), "keybindings.json"))
	reg.RegisterMany(
		Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
		Action{ID: "nav", Help: "mover", Keys: []string{"ctrl+h", "ctrl+l", "ctrl+j", "ctrl+k"}, Label: "ctrl+h/j/k/l"},
	)
	got := NewFooter(reg.Bindings()...).Line(80)
	if !strings.Contains(got, "ctrl+h/j/k/l mover") {
		t.Errorf("footer = %q, want grouped label %q as one hint", got, "ctrl+h/j/k/l mover")
	}
	for _, k := range []string{"ctrl+h", "ctrl+l", "ctrl+j", "ctrl+k"} {
		if strings.Contains(got, k+" mover") {
			t.Errorf("footer = %q, individual key %q should not render as its own hint", got, k)
		}
	}
}
