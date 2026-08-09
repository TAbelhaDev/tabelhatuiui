package tuiui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
)

func mkKeys(keys, desc string) key.Binding {
	return key.NewBinding(key.WithKeys(keys), key.WithHelp(keys, desc))
}

func TestNewThemeFromEnv(t *testing.T) {
	t.Setenv("FOO_DMS_SETTINGS", "/tmp/settings.json")
	t.Setenv("FOO_ACCENT", "blue")
	th := NewThemeFromEnv("FOO")
	if got, want := string(th.Primary), catppuccinAccentHex("blue"); got != want {
		t.Fatalf("NewThemeFromEnv accent = %v, want %v", got, want)
	}

	t.Setenv("FOO_DMS_SETTINGS", "/nonexistent/settings.json")
	t.Setenv("FOO_ACCENT", "")
	th = NewThemeFromEnv("FOO")
	if got, want := string(th.Primary), catppuccinAccentHex("mauve"); got != want {
		t.Fatalf("NewThemeFromEnv default accent = %v, want %v", got, want)
	}
}

func TestFooterLine(t *testing.T) {
	keys := []key.Binding{
		mkKeys("q", "sair"),
		mkKeys("r", "refresh"),
		mkKeys("n", "novo"),
	}

	t.Run("fits", func(t *testing.T) {
		f := NewFooter(keys...).Status("3 jobs")
		got := f.Line(60)
		// Left status + spaces + right hints, total content width = 56.
		if !strings.HasPrefix(got, "3 jobs") {
			t.Errorf("Line() = %q, want prefix %q", got, "3 jobs")
		}
		if !strings.HasSuffix(got, "q sair · r refresh · n novo") {
			t.Errorf("Line() = %q, want suffix %q", got, "q sair · r refresh · n novo")
		}
	})

	t.Run("hints dropped on overflow", func(t *testing.T) {
		f := NewFooter(keys...).Status("")
		got := f.Line(30)
		// "q sair · r refresh" (18 cols) fits in 26; the last hint doesn't.
		if got != "q sair · r refresh" {
			t.Errorf("Line() = %q, want %q", got, "q sair · r refresh")
		}
	})

	t.Run("status truncated, hints dropped", func(t *testing.T) {
		f := NewFooter(keys...).Status("status bem longo que não cabe aqui")
		got := f.Line(24)
		// Status alone exceeds 20 cols: truncate it, drop all hints.
		if !strings.HasPrefix(got, "status bem longo") {
			t.Errorf("Line() = %q, want truncated status prefix", got)
		}
		if strings.Contains(got, "refresh") {
			t.Errorf("Line() = %q, hints should have been dropped", got)
		}
	})

	t.Run("disabled binding skipped", func(t *testing.T) {
		disabled := mkKeys("x", "extra")
		disabled.SetEnabled(false)
		f := NewFooter(keys[0], disabled)
		got := f.Line(60)
		if strings.Contains(got, "extra") {
			t.Errorf("Line() = %q, disabled binding rendered", got)
		}
		if !strings.Contains(got, "sair") {
			t.Errorf("Line() = %q, enabled binding missing", got)
		}
	})

	t.Run("narrow width", func(t *testing.T) {
		f := NewFooter(keys...).Status("long status")
		if got := f.Line(3); got != "" {
			t.Errorf("Line(3) = %q, want empty", got)
		}
	})
}

func TestFooterRender(t *testing.T) {
	theme := ResolveTheme("", "mauve")
	f := NewFooter(mkKeys("q", "sair")).Status("ok")
	out := f.Render(40, theme)
	// The styled footer is exactly 40 columns wide.
	if w := lipgloss.Width(out); w != 40 {
		t.Errorf("Render() width = %d, want 40", w)
	}
}
