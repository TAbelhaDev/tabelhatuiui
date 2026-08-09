package tuiui

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpSection is one titled group of keybindings in the help modal — for
// example "Navegação" with the panel-switching keys, "Ações" with the
// per-app commands.
type HelpSection struct {
	Title    string
	Bindings []Binding
}

// Binding is the keybinding type the help modal and footer consume. It's the
// bubbles/key Binding, aliased so apps only need to import tabelatuiui (the
// bubbles/key types stay available if an app wants to do its own matching).
type Binding = key.Binding

// HelpModal is a centered overlay that lists every keybinding in the app,
// grouped by section — the "?" view every vim/tmux-ish TUI has. Build it
// once with the app's sections, bind "?" to Toggle(), forward Update()
// results, and render it last in View (on top of the app).
//
// Scrolls with j/k (or arrows), closes with q/esc. Sizes itself from the
// latest SetSize call and centers itself over the whole screen.
type HelpModal struct {
	visible  bool
	scroll   int
	width    int
	height   int
	sections []HelpSection
}

// NewHelpModal creates a help modal from the given sections. Bindings that
// are disabled (or have no keys) are skipped at render time.
func NewHelpModal(sections ...HelpSection) *HelpModal {
	return &HelpModal{sections: sections}
}

// Visible reports whether the modal is currently open.
func (m *HelpModal) Visible() bool { return m.visible }

// Toggle opens the modal if it's closed and closes it if it's open.
func (m *HelpModal) Toggle() { m.visible = !m.visible }

// Open and Close explicitly control visibility.
func (m *HelpModal) Open()  { m.visible = true }
func (m *HelpModal) Close() { m.visible = false }

// SetSize records the current viewport so the modal can size/center itself.
// Call it from the app's WindowSizeMsg handler.
func (m *HelpModal) SetSize(width, height int) {
	m.width, m.height = width, height
}

// Update handles keys while the modal is open and reports whether the
// message was consumed. When visible, every message is consumed — the app
// must not process its own keys while the modal is up.
func (m *HelpModal) Update(msg tea.Msg) bool {
	if !m.visible {
		return false
	}
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return true
	}
	switch km.String() {
	case "q", "esc", "ctrl+c":
		m.visible = false
	case "k", "up":
		m.scroll--
	case "j", "down":
		m.scroll++
	case "pgup":
		m.scroll -= 5
	case "pgdown":
		m.scroll += 5
	}
	return true
}

// View renders the modal overlay centered on the screen, or "" when closed.
// Render it last in the app's View so it sits on top of everything else.
func (m *HelpModal) View(theme Theme) string {
	if !m.visible {
		return ""
	}
	width, height := m.width, m.height
	if width <= 0 {
		width = 100
	}
	if height <= 0 {
		height = 30
	}

	modalW := min(width-4, 84)
	if modalW < 20 {
		modalW = 20
	}
	modalH := min(height-4, 34)
	if modalH < 10 {
		modalH = 10
	}
	innerW := modalW - 6
	innerH := modalH - 6

	lines := m.contentLines(theme)

	// Footer hint line (kept inside the box, on top of the content).
	hint := "j/k rolar · q/esc fechar"
	if len(lines) <= innerH {
		hint = "q/esc fechar"
	}
	maxScroll := len(lines) - innerH + 1 // +1 for the hint line
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
	if m.scroll < 0 {
		m.scroll = 0
	}

	visible := lines[m.scroll:]
	avail := innerH - 1
	if len(visible) > avail {
		visible = visible[:avail]
	}
	body := strings.Join(visible, "\n")
	body = PadLines(body, innerW)

	box := theme.Modal().Render(body + "\n" + theme.Dim().Render(hint))
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

// contentLines renders every section into plain rows (already styled), with
// the section title as a dim header line.
func (m *HelpModal) contentLines(theme Theme) []string {
	var out []string
	dim := theme.Dim()
	text := lipgloss.NewStyle().Foreground(theme.Text).Background(theme.Base)

	for _, sec := range m.sections {
		var keys, descs []string
		keyW := 0
		for _, b := range sec.Bindings {
			if !b.Enabled() {
				continue
			}
			h := b.Help()
			keys = append(keys, h.Key)
			descs = append(descs, h.Desc)
			if w := lipgloss.Width(h.Key); w > keyW {
				keyW = w
			}
		}
		if len(keys) == 0 {
			continue
		}
		out = append(out, dim.Render(sec.Title))
		for i := range keys {
			row := keys[i] + strings.Repeat(" ", keyW-lipgloss.Width(keys[i])+2) + descs[i]
			out = append(out, text.Render(row))
		}
		out = append(out, "")
	}
	return out
}
