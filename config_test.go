package tuiui

import (
	"os"
	"path/filepath"
	"testing"
)

type layoutSection struct {
	SidebarWidth int `toml:"sidebar_width"`
	CardHeight   int `toml:"card_height"`
}

type testConfig struct {
	Editor string        `toml:"editor"`
	Roots  []string      `toml:"roots"`
	Layout layoutSection `toml:"layout"`
}

func testDefaults() testConfig {
	return testConfig{
		Editor: "nvim",
		Roots:  []string{"~/codigo/pessoal", "~/codigo/tabeladev"},
		Layout: layoutSection{SidebarWidth: 22, CardHeight: 4},
	}
}

// testConfigAt builds a Config pointing at a file inside a temp dir. The file
// itself is written per-test (or left missing on purpose).
func testConfigAt(t *testing.T) (*Config[testConfig], string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.toml")
	return NewConfig(path, testDefaults()), path
}

func writeConfig(t *testing.T, path, body string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestConfigMissingFileUsesDefaults(t *testing.T) {
	c, _ := testConfigAt(t)
	if err := c.Load(); err != nil {
		t.Fatalf("Load() on missing file = %v, want nil", err)
	}
	if got := c.Get(); got.Editor != "nvim" || got.Layout.SidebarWidth != 22 {
		t.Fatalf("Get() = %+v, want defaults", got)
	}
}

func TestConfigPartialFileMergesOntoDefaults(t *testing.T) {
	c, path := testConfigAt(t)
	// Only one key of one section is set; everything else must survive.
	writeConfig(t, path, "[layout]\nsidebar_width = 40\n")
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	got := c.Get()
	if got.Layout.SidebarWidth != 40 {
		t.Fatalf("Layout.SidebarWidth = %d, want 40", got.Layout.SidebarWidth)
	}
	if got.Layout.CardHeight != 4 {
		t.Fatalf("Layout.CardHeight = %d, want the default 4", got.Layout.CardHeight)
	}
	if got.Editor != "nvim" {
		t.Fatalf("Editor = %q, want the default %q", got.Editor, "nvim")
	}
	if len(got.Roots) != 2 {
		t.Fatalf("Roots = %v, want the 2 defaults", got.Roots)
	}
}

func TestConfigSliceReplacesDefaultEntirely(t *testing.T) {
	c, path := testConfigAt(t)
	writeConfig(t, path, `roots = ["/only/this/one"]`+"\n")
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	got := c.Get().Roots
	if len(got) != 1 || got[0] != "/only/this/one" {
		t.Fatalf("Roots = %v, want exactly [/only/this/one] (no leftover defaults)", got)
	}
}

func TestConfigMalformedFileKeepsPreviousValue(t *testing.T) {
	c, path := testConfigAt(t)
	writeConfig(t, path, "[layout]\nsidebar_width = 40\n")
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	// A typo mid-edit must not drop a running app back to defaults.
	writeConfig(t, path, "[layout\nsidebar_width = ")
	changed, err := c.Reload()
	if err == nil {
		t.Fatal("Reload() on malformed TOML = nil error, want a parse error")
	}
	if changed {
		t.Fatal("Reload() on malformed TOML should not report a change")
	}
	if got := c.Get().Layout.SidebarWidth; got != 40 {
		t.Fatalf("after a failed Reload, SidebarWidth = %d, want the previous 40", got)
	}
}

func TestConfigReloadReportsChange(t *testing.T) {
	c, path := testConfigAt(t)
	writeConfig(t, path, "editor = \"hx\"\n")
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	// External edit → picked up and reported.
	writeConfig(t, path, "editor = \"vim\"\n")
	changed, err := c.Reload()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Fatal("Reload() should report a change after an external edit")
	}
	if got := c.Get().Editor; got != "vim" {
		t.Fatalf("Editor = %q, want vim", got)
	}

	// Unchanged file → no change reported.
	if changed, err = c.Reload(); err != nil || changed {
		t.Fatalf("Reload() on unchanged file = changed:%v, err:%v; want false,nil", changed, err)
	}
}

func TestConfigReloadDeletedFileFallsBackToDefaults(t *testing.T) {
	c, path := testConfigAt(t)
	writeConfig(t, path, "editor = \"hx\"\n")
	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	changed, err := c.Reload()
	if err != nil {
		t.Fatalf("Reload() after deleting the file = %v, want nil", err)
	}
	if !changed {
		t.Fatal("Reload() should report a change back to defaults")
	}
	if got := c.Get().Editor; got != "nvim" {
		t.Fatalf("Editor = %q, want the default nvim", got)
	}
}

func TestConfigPath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/tmp/xdg-test")
	if got, want := ConfigPath("tabelaradar", "config.toml"), "/tmp/xdg-test/tabelaradar/config.toml"; got != want {
		t.Fatalf("ConfigPath() = %q, want %q", got, want)
	}
}
