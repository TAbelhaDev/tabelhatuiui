package tuiui

import (
	"encoding/json"
	"os"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

// Theme is the resolved Catppuccin Mocha palette plus the DMS accent used
// as Primary — the shared chrome every ianptkcs TUI renders with. Colors
// follow the official semantic guide:
// https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md
type Theme struct {
	Base     lipgloss.Color
	Mantle   lipgloss.Color
	Surface0 lipgloss.Color
	Surface1 lipgloss.Color
	Overlay0 lipgloss.Color
	Overlay1 lipgloss.Color
	Text     lipgloss.Color
	Subtext0 lipgloss.Color

	// Primary mirrors the installed DankMaterialShell's own configured
	// accent (falling back to a manually chosen Catppuccin accent) — see
	// resolvePrimaryHex. Consumers read the same DMS settings.json djobs and
	// tabelaradar do, so every tool's chrome matches whatever accent DMS is
	// set to.
	Primary lipgloss.Color

	Red      lipgloss.Color
	Green    lipgloss.Color
	Yellow   lipgloss.Color
	Blue     lipgloss.Color
	Pink     lipgloss.Color
	Lavender lipgloss.Color
}

// ResolveTheme builds the theme: Primary comes from the DMS settings.json at
// settingsPath (empty string means "not resolvable", falling back to the
// Catppuccin accent fallbackAccent).
func ResolveTheme(settingsPath, fallbackAccent string) Theme {
	mocha := catppuccin.Mocha
	return Theme{
		Base:     lipgloss.Color(mocha.Base().Hex),
		Mantle:   lipgloss.Color(mocha.Mantle().Hex),
		Surface0: lipgloss.Color(mocha.Surface0().Hex),
		Surface1: lipgloss.Color(mocha.Surface1().Hex),
		Overlay0: lipgloss.Color(mocha.Overlay0().Hex),
		Overlay1: lipgloss.Color(mocha.Overlay1().Hex),
		Text:     lipgloss.Color(mocha.Text().Hex),
		Subtext0: lipgloss.Color(mocha.Subtext0().Hex),

		Primary: lipgloss.Color(resolvePrimaryHex(settingsPath, fallbackAccent)),

		Red:      lipgloss.Color(mocha.Red().Hex),
		Green:    lipgloss.Color(mocha.Green().Hex),
		Yellow:   lipgloss.Color(mocha.Yellow().Hex),
		Blue:     lipgloss.Color(mocha.Blue().Hex),
		Pink:     lipgloss.Color(mocha.Pink().Hex),
		Lavender: lipgloss.Color(mocha.Lavender().Hex),
	}
}

type dmsThemeVariant struct {
	Flavor string `json:"flavor"`
	Accent string `json:"accent"`
}

type dmsSettings struct {
	CurrentThemeCategory  string `json:"currentThemeCategory"`
	CustomThemeFile       string `json:"customThemeFile"`
	RegistryThemeVariants map[string]struct {
		Dark dmsThemeVariant `json:"dark"`
	} `json:"registryThemeVariants"`
}

type catppuccinFlavorColors struct {
	Primary string `json:"primary"`
}

type catppuccinAccent struct {
	ID        string                 `json:"id"`
	Frappe    catppuccinFlavorColors `json:"frappe"`
	Latte     catppuccinFlavorColors `json:"latte"`
	Macchiato catppuccinFlavorColors `json:"macchiato"`
	Mocha     catppuccinFlavorColors `json:"mocha"`
}

func (a catppuccinAccent) primaryForFlavor(flavor string) string {
	switch flavor {
	case "frappe":
		return a.Frappe.Primary
	case "latte":
		return a.Latte.Primary
	case "macchiato":
		return a.Macchiato.Primary
	case "mocha":
		return a.Mocha.Primary
	default:
		return ""
	}
}

type dmsTheme struct {
	ID       string `json:"id"`
	Variants struct {
		Accents []catppuccinAccent `json:"accents"`
	} `json:"variants"`
}

// dmsAccentHex resolves the primary accent hex the installed
// DankMaterialShell is currently rendering, by reading its own settings.json
// + the theme.json it references — the same lookup DMS itself performs. Only
// understands DMS's Catppuccin registry theme; any other theme category,
// registry theme, or missing/malformed file returns "" so the caller falls
// back.
func dmsAccentHex(settingsPath string) string {
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return ""
	}
	var settings dmsSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return ""
	}
	if settings.CurrentThemeCategory != "registry" || settings.CustomThemeFile == "" {
		return ""
	}

	themeData, err := os.ReadFile(settings.CustomThemeFile)
	if err != nil {
		return ""
	}
	var theme dmsTheme
	if err := json.Unmarshal(themeData, &theme); err != nil {
		return ""
	}
	if theme.ID != "catppuccin" {
		return ""
	}

	variant, ok := settings.RegistryThemeVariants[theme.ID]
	if !ok {
		return ""
	}
	for _, accent := range theme.Variants.Accents {
		if accent.ID == variant.Dark.Accent {
			return accent.primaryForFlavor(variant.Dark.Flavor)
		}
	}
	return ""
}

// catppuccinAccentHex maps a Catppuccin Mocha accent id to its hex. Unknown
// ids fall back to mauve.
func catppuccinAccentHex(id string) string {
	mocha := catppuccin.Mocha
	switch id {
	case "rosewater":
		return mocha.Rosewater().Hex
	case "flamingo":
		return mocha.Flamingo().Hex
	case "pink":
		return mocha.Pink().Hex
	case "red":
		return mocha.Red().Hex
	case "maroon":
		return mocha.Maroon().Hex
	case "peach":
		return mocha.Peach().Hex
	case "yellow":
		return mocha.Yellow().Hex
	case "green":
		return mocha.Green().Hex
	case "teal":
		return mocha.Teal().Hex
	case "sky":
		return mocha.Sky().Hex
	case "sapphire":
		return mocha.Sapphire().Hex
	case "blue":
		return mocha.Blue().Hex
	case "lavender":
		return mocha.Lavender().Hex
	case "mauve":
		return mocha.Mauve().Hex
	default:
		return mocha.Mauve().Hex
	}
}

// resolvePrimaryHex is ResolveTheme's accent source: DMS's own live accent
// when available, else the manually configured (or default "mauve") fallback.
func resolvePrimaryHex(settingsPath, fallbackAccent string) string {
	if hex := dmsAccentHex(settingsPath); hex != "" {
		return hex
	}
	return catppuccinAccentHex(fallbackAccent)
}
