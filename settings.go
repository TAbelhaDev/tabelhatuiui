package tuiui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// SettingsModal is the keybinding editor: a centered overlay listing every
// registered action with its current keys, letting the user rebind (enter →
// "press a key"), reset an action (r) or everything (R). Overrides go
// straight to the KeyRegistry, which persists them.
//
// Built on top of a KeyRegistry, so the footer, the help modal and the key
// dispatch all pick up changes immediately. Bind a key (e.g. "," or "s") to
// Toggle(), forward Update() results, and render it last in View.
type SettingsModal struct {
	registry *KeyRegistry

	visible bool
	cursor  int
	scroll  int
	width   int
	height  int

	// capturing means the next keypress rebinds captureID.
	capturing bool
	captureID string
	// err holds the last failed Set message (e.g. a key conflict).
	err string
}

// NewSettingsModal creates a settings modal editing the given registry.
func NewSettingsModal(r *KeyRegistry) *SettingsModal {
	return &SettingsModal{registry: r}
}

// Visible reports whether the modal is currently open.
func (m *SettingsModal) Visible() bool { return m.visible }

// Toggle opens the modal if closed and closes it if open.
func (m *SettingsModal) Toggle() {
	m.visible = !m.visible
	if !m.visible {
		m.capturing = false
		m.err = ""
	}
}

// Open and Close explicitly control visibility.
func (m *SettingsModal) Open()  { m.visible = true }
func (m *SettingsModal) Close() { m.visible = false; m.capturing = false; m.err = "" }

// SetSize records the current viewport so the modal can size/center itself.
func (m *SettingsModal) SetSize(width, height int) {
	m.width, m.height = width, height
}

// Update handles keys while the modal is open and reports whether the message
// was consumed. While visible, every message is consumed — the app must not
// process its own keys.
func (m *SettingsModal) Update(msg tea.Msg) bool {
	if !m.visible {
		return false
	}
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return true
	}

	if m.capturing {
		return m.updateCapturing(km)
	}

	switch km.String() {
	case "q", "esc", "ctrl+c":
		m.Close()
	case "j", "down":
		if m.cursor < len(m.registry.order)-1 {
			m.cursor++
		}
	case "k", "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "g":
		m.cursor = 0
	case "G":
		m.cursor = len(m.registry.order) - 1
	case "pgup":
		m.cursor -= 5
		if m.cursor < 0 {
			m.cursor = 0
		}
	case "pgdown":
		m.cursor += 5
		if m.cursor > len(m.registry.order)-1 {
			m.cursor = len(m.registry.order) - 1
		}
	case "enter":
		if m.cursor < len(m.registry.order) {
			m.capturing = true
			m.captureID = m.registry.order[m.cursor]
			m.err = ""
		}
	case "r":
		if m.cursor < len(m.registry.order) {
			if err := m.registry.Reset(m.registry.order[m.cursor]); err != nil {
				m.err = "erro ao resetar: " + err.Error()
			} else {
				m.err = ""
			}
		}
	case "R":
		if _, err := m.registry.ResetAll(); err != nil {
			m.err = "erro ao resetar tudo: " + err.Error()
		} else {
			m.err = ""
		}
	}
	return true
}

// updateCapturing handles the "press a key" state. esc/ctrl+c cancel; any
// other key becomes the new binding for captureID.
func (m *SettingsModal) updateCapturing(km tea.KeyMsg) bool {
	switch km.String() {
	case "esc", "ctrl+c":
		m.capturing = false
		m.captureID = ""
		return true
	}
	keys := []string{km.String()}
	if err := m.registry.Set(m.captureID, keys...); err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
	m.capturing = false
	m.captureID = ""
	return true
}

// View renders the settings overlay centered on the screen, or "" when
// closed. Render it last in the app's View.
func (m *SettingsModal) View(theme Theme) string {
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

	modalW := min(width-4, 70)
	if modalW < 30 {
		modalW = 30
	}
	modalH := min(height-4, 28)
	if modalH < 12 {
		modalH = 12
	}
	innerW := modalW - 6
	innerH := modalH - 6

	rows := m.contentRows(theme, innerW)

	maxScroll := len(rows) - (innerH - 2) // 1 hint + 1 capture/err line
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.scroll > maxScroll {
		m.scroll = maxScroll
	}
	if m.scroll < 0 {
		m.scroll = 0
	}

	visible := rows[m.scroll:]
	if len(visible) > innerH-2 {
		visible = visible[:innerH-2]
	}
	body := strings.Join(visible, "\n")
	body = PadLines(body, innerW)

	var statusLine string
	if m.capturing {
		statusLine = theme.Warning().Render("nova tecla para '" + m.captureID + "'... (esc cancela)")
	} else if m.err != "" {
		statusLine = theme.Error().Render(m.err)
	} else {
		statusLine = theme.Dim().Render("enter editar · r reset · R reset todos · q/esc fechar")
	}

	box := theme.Modal().Render(body + "\n\n" + statusLine)
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, box)
}

// contentRows renders each action as a row with a cursor marker, the current
// keys (custom ones highlighted) and the description.
func (m *SettingsModal) contentRows(theme Theme, width int) []string {
	actions := m.registry.Actions()
	// Widest keys column across all rows.
	keyW := 0
	for _, a := range actions {
		if w := lipgloss.Width(a.Label); w > keyW {
			keyW = w
		}
	}
	keyW = min(keyW, 20)

	dim := theme.Dim()
	sel := lipgloss.NewStyle().Foreground(theme.Base).Background(theme.Primary).Bold(true)
	custom := theme.Info()
	rows := make([]string, 0, len(actions))
	for i, a := range actions {
		keysCell := a.Label + strings.Repeat(" ", keyW-lipgloss.Width(a.Label)+2)
		desc := a.Help
		if width > 0 && lipgloss.Width(keysCell)+lipgloss.Width(desc)+2 > width {
			desc = ansi.Truncate(desc, max(1, width-lipgloss.Width(keysCell)-2), "…")
		}
		marker := "  "
		if m.registry.IsCustom(a.ID) {
			marker = "● "
		}
		row := marker + keysCell + desc
		if i == m.cursor {
			row = sel.Render(" " + strings.TrimLeft(row, " ") + " ")
		} else if m.registry.IsCustom(a.ID) {
			row = custom.Render(row)
		} else {
			row = dim.Render(row)
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		rows = []string{dim.Render("nenhuma ação registrada")}
	}
	return rows
}
