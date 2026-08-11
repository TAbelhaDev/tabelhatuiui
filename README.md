<div align="center">

# TabelaTuiUI

**Shared theme and chrome for ianptkcs's Bubble Tea TUIs.**

**English** · [Português](README.pt-BR.md)

[![Go Version](https://img.shields.io/github/go-mod/go-version/TabelaDev/tabelatuiui?style=flat-square&logo=go&logoColor=white&color=00ADD8)](go.mod)
[![License: AGPL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue?style=flat-square)](LICENSE)
[![Built with Bubble Tea](https://img.shields.io/badge/built%20with-Bubble%20Tea-ff69b4?style=flat-square)](https://github.com/charmbracelet/bubbletea)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/ianptkcs)

</div>

---

## What it is

A small library holding whatever repeats across my Bubble Tea TUIs — the theme
(Catppuccin Mocha + the DankMaterialShell accent), the chrome styles
(header/footer/panels/modals), the layout helpers (ANSI-aware
pad/wrap/truncate) and the `ipc <method> [key=value...] --json` convention.

Each app keeps only what is its own (model, keys, business logic); the library
takes care of what everyone was redrawing from scratch in every new project.

## Who uses it

| Project | What it is | Origin |
|---|---|---|
| [djobs](https://github.com/ianptkcs/dankjobs) | TUI for scheduled jobs (systemd) | migrated |
| [tabelaradar](https://github.com/TabelaDev/tabelaradar) | TUI that audits git repos | migrated |
| [tabelakanban](https://github.com/TabelaDev/tabelakanban) | kanban TUI | born consuming it |

## Installation

Requires Go 1.26+.

```bash
go get github.com/ianptkcs/tabelatuiui@latest
```

## Usage

Resolve the theme in your `main` and use the styles wherever needed:

```go
package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/ianptkcs/tabelatuiui"
)

// Theme through env: reads MEUAPP_DMS_SETTINGS and MEUAPP_ACCENT (mauve by default).
var theme = tuiui.NewThemeFromEnv("MEUAPP")

var (
	keyQuit    = key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "sair"))
	keyRefresh = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "atualizar"))
	keyNavL    = key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "próx. painel"))
)

// Footer: context on the left + hints on the right, generated from the bindings.
footer := tuiui.NewFooter(keyQuit, keyRefresh, keyNavL).
	Status("3 jobs").
	Render(80, theme)

// HelpModal: the "?" every TUI has, listing the bindings in sections.
helpModal := tuiui.NewHelpModal(tuiui.HelpSection{
	Title:    "Navegação",
	Bindings: []tuiui.Binding{keyNavL},
})

func main() {
	header := theme.Header(80).Render("meu app")
	panel := theme.Panel(true).Render(tuiui.PadLines("conteúdo", 40))
	modal := theme.Modal().Render("tem certeza?")
	_ = header + footer + panel + modal
	_ = helpModal
}
```

### Theme resolution

`NewThemeFromEnv(appPrefix)` is the recommended shortcut: it reads
`<PREFIX>_DMS_SETTINGS` (default
`~/.config/DankMaterialShell/settings.json`) and `<PREFIX>_ACCENT` (default
`mauve`) and calls `ResolveTheme`. For full control,
`ResolveTheme(settingsPath, fallbackAccent)` reads the installed
DankMaterialShell's `settings.json` and returns the hex of the accent DMS is
rendering today (the same lookup DMS itself does). When DMS is not installed, or
is on another theme, it falls back to the Catppuccin accent passed in.

### Footer and HelpModal

- `NewFooter(bindings...)` assembles the status/help bar. The right-hand side is
  generated from the bindings, so the hints never diverge from what `key.Matches`
  accepts. `Status()` sets the left-hand context; `Render(width, theme)` returns
  the finished line. Hints that do not fit are dropped token by token (` · `).
- `NewHelpModal(sections...)` creates a centred overlay with the bindings in
  sections. Bind it to `"?"` (or another key), call `Update()` (it returns `true`
  while the modal is open — the app must not process its keys meanwhile),
  `SetSize()` on resize, and render it last in `View`. Sections can use
  `BindingsFn` to fetch the live bindings on every render.

### Customisable keybindings (config-file-first)

`KeyRegistry` centralises an app's keybindings and makes the config file the
**source of truth**: defaults registered in code + the user's overrides in
`<ConfigDir()>/<app>/keybindings.json`. Create it, register the actions and load:

```go
reg := tuiui.NewKeyRegistry(tuiui.ConfigPath("meuapp", "keybindings.json"))
reg.RegisterMany(
	tuiui.Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
	tuiui.Action{ID: "nav", Help: "mover", Keys: []string{"ctrl+h", "ctrl+l", "ctrl+j", "ctrl+k"}, Label: "ctrl+h/j/k/l"},
)
if err := reg.Load(); err != nil { /* corrupt file: falls back to defaults */ }
```

The format is `{"bindings": {"<id>": ["key", ...]}}`, where `<id>` is the
`Action`'s `ID`; the file only exists while there is at least one override. The
primary flow is **edit the file and reload** — expose a reload key (or call
`Reload()` on every `View()`) so an external edit takes effect without a restart:

```go
case key.Matches(msg, reg.Resolve("reload")): // e.g. "r"
	if changed, err := reg.Reload(); err != nil {
		status = theme.Error().Render("keybindings: " + err.Error())
	} else if changed {
		status = theme.Success().Render("keybindings recarregados")
	}
```

- **Dispatch**: `case key.Matches(msg, reg.Resolve("quit")):`
- **Footer**: `tuiui.NewFooter(reg.Bindings()...)` — reflects rebinds and reloads
  immediately.
- **Help modal**: pass `BindingsFn: reg.Bindings` in the section.
- **Grouped hints (the canonical navigation pattern)**: a single action with
  several keys + a `Label` (e.g. `nav` with `ctrl+h/l/j/k`, `Label:
  "ctrl+h/j/k/l"`) becomes **one line** in the footer/help instead of N entries;
  dispatch by direction stays with the app (position-based), not with the key. Use
  `Label` when the keys share one concept.
- **SettingsModal (optional convenience)**: `tuiui.NewSettingsModal(reg)` — bind
  it to `","`/`"s"` to open the action list with interactive rebinding (`enter`
  rebinds, `r`/`R` reset, conflicts in red, custom marked with `●`). With the
  file-based flow above it becomes a shortcut, not the only path.

### TOML config (config-file-first)

`Config[T]` is `KeyRegistry`'s equivalent for the app's other preferences:
defaults compiled into the code + the user's overrides in
`<ConfigDir()>/<app>/config.toml`. Each app defines its own `T`; the library has
no opinion about the schema.

```go
type config struct {
	Editor string `toml:"editor"`
	Layout struct {
		SidebarWidth int `toml:"sidebar_width"`
	} `toml:"layout"`
}

var defaults = config{Editor: "nvim"} // defaults.Layout.SidebarWidth = 22, etc.

cfg := tuiui.NewConfig(tuiui.ConfigPath("meuapp", "config.toml"), defaults)
if err := cfg.Load(); err != nil { /* corrupt file: keeps the defaults */ }

width := cfg.Get().Layout.SidebarWidth
```

The merge is per key: the TOML only overrides the fields that **appear in the
file**, and everything else keeps its value from `defaults`. Slices are replaced
whole (a `roots = [...]` in the file swaps the entire list, it does not
concatenate). A missing file is not an error — the app runs on pure defaults.

As with the keybindings, the primary flow is **edit the file and reload**, on the
same key:

```go
case key.Matches(msg, reg.Resolve("reload")): // f5
	kChanged, kErr := reg.Reload()
	cChanged, cErr := cfg.Reload()
	switch {
	case kErr != nil || cErr != nil:
		status = theme.Error().Render("config: " + errors.Join(kErr, cErr).Error())
	case kChanged || cChanged:
		status = theme.Success().Render("config recarregada")
	default:
		status = theme.Muted().Render("config sem mudanças")
	}
```

`Reload()` reports whether the **effective** config changed, and a malformed TOML
returns an error while **preserving the previous value** — a typo mid-edit does
not drop a running app to its defaults. Not every field is hot-reloadable (a
database path, a worker count): `Reload()` says something changed, and it is up to
the app to decide what to do about it.

> `T` should be a struct of values. Pointer/map/slice fields are shared with
> `defaults` until the file overrides them, so do not mutate what `Get()` returns.

**`ConfigPath(app, file)`** resolves `~/.config/<app>/<file>` while respecting
`XDG_CONFIG_HOME` — use it for both `config.toml` and `keybindings.json` instead
of assembling the path by hand.

### Layout helpers

- `PadLines(s, width)` — ANSI-aware pad/truncate of each line to an exact width
  (needed for panel borders to line up).
- `WrapText(s, width)` — line wrapping over plain text.
- `PadToHeight(s, lines)` — pads or truncates a block to an exact height.

### Semantic styles

`Success` (green), `Warning` (yellow), `Error` (red), `Info` (blue) and `Muted`
(grey) follow Catppuccin's semantics guide — for statuses each app colours its own
way but with the same meaning.

### IPC

```go
args, err := tuiui.ParseIPCArgs(os.Args[2:]) // method + filters + --json
if err != nil {
	fmt.Fprintln(os.Stderr, "uso: meuapp ipc <método> [key=value...] --json")
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
// dispatch on args.Method...
os.Exit(tuiui.WriteJSON(resultado))
```

`EnvOr`, `ExpandHome`, `HomeDir` and `ConfigDir` close out the rest of the
skeleton every app used to repeat.

## Development

```bash
go test ./...
```

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for the version history.

## Support the project

- **Global**: [ko-fi.com/ianptkcs](https://ko-fi.com/ianptkcs)
- **Brazil (Pix)**: scan the QR below or copy the code

  <img src="pix-qr.png" alt="Pix QR" width="200" />

  <details><summary>Pix code (copy)</summary>

  ```
00020126580014BR.GOV.BCB.PIX01365ad933b0-dcdc-4525-a736-0759902aeec65204000053039865802BR5925Ian Patrick da Costa Soar6009SAO PAULO62140510tQA85x6Dov63041FB6
  ```

  </details>

## License

[GNU AGPL-3.0](LICENSE) — the same license as the TUIs that use it.
