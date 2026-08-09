package tuiui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Footer builds the bottom status/help bar: a context string on the left
// (a message, a counter, "perfil: dev") and the keybinding hints on the
// right, generated from a slice of bubbles/key bindings so the hints can
// never drift out of sync with what Update() actually matches.
//
// The layout follows what every ianptkcs TUI hand-rolled before this type
// existed: left and right sides separated by a run of spaces, with
// token-aware truncation (hints are dropped at " · " boundaries from the
// right, then the left side is truncated) when the line overflows.
type Footer struct {
	status   string
	bindings []Binding
}

// NewFooter starts a footer whose right side is generated from the given
// keybindings. Disabled bindings are skipped.
func NewFooter(bindings ...Binding) *Footer {
	return &Footer{bindings: bindings}
}

// Status sets the left-side context text. It can carry ANSI (callers pass
// already-styled message strings); the width math is ANSI-aware.
func (f *Footer) Status(s string) *Footer {
	f.status = s
	return f
}

// Line returns the footer's unstyled content, sized to fit a Footer(width)
// render — that is, trimmed to width-4 columns (the Footer style's horizontal
// padding + border consume the other 4). The line never wraps: whatever
// doesn't fit gets dropped or truncated.
func (f *Footer) Line(width int) string {
	avail := width - 4
	if avail <= 0 {
		return ""
	}
	parts := f.hints()
	hints := strings.Join(parts, " · ")
	leftW := lipgloss.Width(f.status)
	hintsW := lipgloss.Width(hints)

	switch {
	case f.status == "":
		return fitHints(parts, avail)
	case leftW+hintsW <= avail:
		return f.status + strings.Repeat(" ", avail-leftW-hintsW) + hints
	case leftW+3 > avail:
		return ansi.Truncate(f.status, avail, "…")
	default:
		return f.status + " " + fitHints(parts, avail-leftW-1)
	}
}

// Render renders the footer fully styled with the given theme. Drop the
// result straight into the View output.
func (f *Footer) Render(width int, theme Theme) string {
	return theme.Footer(width).Render(f.Line(width))
}

// hints returns "key desc" for every enabled binding, in order.
func (f *Footer) hints() []string {
	parts := make([]string, 0, len(f.bindings))
	for _, b := range f.bindings {
		if !b.Enabled() {
			continue
		}
		h := b.Help()
		parts = append(parts, h.Key+" "+h.Desc)
	}
	return parts
}

// fitHints drops whole " · "-separated hints from the right until the joined
// string fits width, then ANSI-truncates the first hint if even it doesn't
// fit (or width <= 0, in which case everything is dropped).
func fitHints(parts []string, width int) string {
	if width <= 0 {
		return ""
	}
	joined := strings.Join(parts, " · ")
	if lipgloss.Width(joined) <= width {
		return joined
	}
	for len(parts) > 1 {
		parts = parts[:len(parts)-1]
		if w := lipgloss.Width(strings.Join(parts, " · ")); w <= width {
			return strings.Join(parts, " · ")
		}
	}
	return ansi.Truncate(parts[0], width, "…")
}
