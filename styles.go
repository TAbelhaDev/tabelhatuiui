package tuiui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Header renders the top status bar: primary accent background, base text,
// full width.
func (t Theme) Header(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Primary).
		Foreground(t.Base).
		Bold(true).
		Width(width).
		Padding(0, 2)
}

// Footer renders the bottom help/status bar: mantle background, subtext
// text, full width.
func (t Theme) Footer(width int) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.Mantle).
		Foreground(t.Subtext0).
		Width(width).
		Padding(0, 2)
}

// Panel intentionally has no Width(): calling Width() on a style makes
// lipgloss re-wrap its content, and that wrap logic miscounts lines that
// already carry their own nested ANSI (a table's selected-row highlight, a
// highlighted card), breaking alignment. Content is pre-padded to a uniform
// width with PadLines instead, so the border ends up sized correctly on its
// own. focused switches the border to Primary (also used for headers/
// titles), so the currently-navigable panel is obvious.
func (t Theme) Panel(focused bool) lipgloss.Style {
	border := t.Surface1
	if focused {
		border = t.Primary
	}
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(border).
		Background(t.Base).
		Padding(0, 1)
}

// Title is a panel's first line, bold in the primary accent. Background is
// set explicitly for the same reason Panel's is: its own reset must not
// blank out the panel's background for that line.
func (t Theme) Title() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(t.Primary).Background(t.Base)
}

func (t Theme) Dim() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Overlay0).Background(t.Base)
}

// Modal is a centered overlay box with an accent (Primary) border — used by
// djobs for its form/confirm dialogs, generic enough to live in the lib.
func (t Theme) Modal() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Primary).
		Background(t.Base).
		Padding(1, 2)
}

// Semantic status styles follow the Catppuccin style guide's meanings:
// Success for active/ok, Warning for paused/delayed, Error for a hard
// problem, Info for a neutral "done" state. Background is set to Base so the
// colored text reads correctly wherever it's dropped in (a table cell, a
// description line, a modal) — djobs' colorizeStatusColumn needs this.
func (t Theme) Success() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Green).Background(t.Base).Bold(true)
}

func (t Theme) Warning() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Yellow).Background(t.Base).Bold(true)
}

func (t Theme) Error() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Red).Background(t.Base).Bold(true)
}

func (t Theme) Info() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Blue).Background(t.Base).Bold(true)
}

// Muted is the "no longer relevant" fallback — same as Dim but one step
// more visible (Overlay1), used for stale/secondary content.
func (t Theme) Muted() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.Overlay1).Background(t.Base)
}

// PadLines pads or truncates every line of s so its ANSI-aware visible width
// is exactly `width` — a border box only ends up sized (and positioned)
// correctly if every line it wraps is uniform. The truncating side uses an
// ANSI-aware truncate (ansi.Truncate, not go-runewidth's): lines here can
// carry nested SGR codes (table header/selected-row colors, highlighted
// cards), and a truncate that counts escape-sequence bytes as visible width
// cuts real content far too early, garbling it.
func PadLines(s string, width int) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		switch w := lipgloss.Width(line); {
		case w < width:
			lines[i] = line + strings.Repeat(" ", width-w)
		case w > width:
			lines[i] = ansi.Truncate(line, width, "…")
		}
	}
	return strings.Join(lines, "\n")
}

// WrapText word-wraps plain text (no embedded ANSI — see Panel's comment on
// why Width()-based wrapping doesn't mix with already-styled content) to
// width, so long natural-language strings break onto new lines instead of
// being cut short with "…" by PadLines.
func WrapText(s string, width int) string {
	if width <= 0 {
		return s
	}
	return lipgloss.NewStyle().Width(width).Render(s)
}

// PadToHeight pads s with blank lines until it has exactly `lines` lines (or
// truncates extra ones) — used for panels whose natural content is shorter
// than their computed box height, so a fixed-height neighbor panel's border
// still lines up with this one's bottom border.
func PadToHeight(s string, lines int) string {
	got := strings.Split(s, "\n")
	for len(got) < lines {
		got = append(got, "")
	}
	if len(got) > lines {
		got = got[:lines]
	}
	return strings.Join(got, "\n")
}
