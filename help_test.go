package tuiui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestHelpModalToggleAndClose(t *testing.T) {
	m := NewHelpModal(HelpSection{Title: "Nav", Bindings: []Binding{mkKeys("h", "esq")}})
	if m.Visible() {
		t.Fatal("new modal should be hidden")
	}
	m.Toggle()
	if !m.Visible() {
		t.Fatal("Toggle() should open the modal")
	}
	m.Toggle()
	if m.Visible() {
		t.Fatal("second Toggle() should close the modal")
	}
	m.Open()
	m.Close()
	if m.Visible() {
		t.Fatal("Close() should hide the modal")
	}
}

func TestHelpModalUpdate(t *testing.T) {
	m := NewHelpModal(HelpSection{Title: "Nav", Bindings: []Binding{mkKeys("h", "esq")}})

	// Closed modal consumes nothing.
	if m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}) {
		t.Fatal("closed modal should not consume messages")
	}

	m.Open()
	if !m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}) {
		t.Fatal("open modal should consume messages")
	}
	if m.Visible() {
		t.Fatal("q should close the modal")
	}

	m.Open()
	m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if m.Visible() {
		t.Fatal("esc should close the modal")
	}
}

func TestHelpModalView(t *testing.T) {
	m := NewHelpModal(
		HelpSection{Title: "Navegação", Bindings: []Binding{
			mkKeys("ctrl+h", "painel esq"),
			mkKeys("ctrl+l", "painel dir"),
		}},
		HelpSection{Title: "Ações", Bindings: []Binding{mkKeys("q", "sair")}},
	)
	m.SetSize(80, 24)
	m.Open()

	theme := ResolveTheme("", "mauve")
	out := m.View(theme)
	if out == "" {
		t.Fatal("open modal should render content")
	}
	for _, want := range []string{"Navegação", "Ações", "painel esq", "sair"} {
		if !strings.Contains(out, want) {
			t.Errorf("View() missing %q in output", want)
		}
	}
	if !strings.Contains(out, "fechar") {
		t.Errorf("View() missing close hint")
	}

	m.Close()
	if out := m.View(theme); out != "" {
		t.Errorf("closed modal should render empty, got %q", out)
	}
}

func TestHelpModalSkipsDisabled(t *testing.T) {
	hidden := mkKeys("x", "escondido")
	hidden.SetEnabled(false)
	m := NewHelpModal(HelpSection{Title: "T", Bindings: []Binding{hidden, mkKeys("q", "sair")}})
	m.SetSize(60, 20)
	m.Open()
	out := m.View(ResolveTheme("", "mauve"))
	if strings.Contains(out, "escondido") {
		t.Errorf("View() rendered disabled binding")
	}
}
