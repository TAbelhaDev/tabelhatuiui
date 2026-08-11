<div align="center">

# TabelaTuiUI

**Tema + chrome compartilhado dos TUIs Bubble Tea do ianptkcs.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/TabelaDev/tabelatuiui?style=flat-square&logo=go&logoColor=white&color=00ADD8)](go.mod)
[![License: AGPL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue?style=flat-square)](LICENSE)
[![Built with Bubble Tea](https://img.shields.io/badge/built%20with-Bubble%20Tea-ff69b4?style=flat-square)](https://github.com/charmbracelet/bubbletea)

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/ianptkcs)

</div>

---

## O que é

Uma bibliotinha que concentra o que se repete entre os meus TUIs Bubble
Tea — o tema (Catppuccin Mocha + accent do DankMaterialShell), os estilos de
chrome (header/footer/panels/modais), os helpers de layout (pad/wrap/truncate
ANSI-aware) e a convenção de IPC `ipc <método> [key=value...] --json`.

Cada app mantém só o que é dele (modelo, teclas, negócio); a lib cuida do que
todo mundo desenhava do zero em todo projeto novo.

## Quem usa

| Projeto | O que é | De onde veio |
|---|---|---|
| [djobs](https://github.com/ianptkcs/dankjobs) | TUI de jobs agendados (systemd) | migrado |
| [tabelaradar](https://github.com/TabelaDev/tabelaradar) | TUI de fiscalização de repos git | migrado |
| [tabelakanban](https://github.com/TabelaDev/tabelakanban) | TUI kanban | nasceu consumindo |

## Instalação

Requer Go 1.26+.

```bash
go get github.com/ianptkcs/tabelatuiui@latest
```

## Uso

Resolve o tema no seu `main` e usa os estilos onde precisar:

```go
package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/ianptkcs/tabelatuiui"
)

// Theme via env: lê MEUAPP_DMS_SETTINGS e MEUAPP_ACCENT (padrão mauve).
var theme = tuiui.NewThemeFromEnv("MEUAPP")

var (
	keyQuit    = key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "sair"))
	keyRefresh = key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "atualizar"))
	keyNavL    = key.NewBinding(key.WithKeys("ctrl+l"), key.WithHelp("ctrl+l", "próx. painel"))
)

// Footer: contexto à esquerda + hints à direita, gerados das bindings.
footer := tuiui.NewFooter(keyQuit, keyRefresh, keyNavL).
	Status("3 jobs").
	Render(80, theme)

// HelpModal: o "?" de todo TUI, listando as bindings em seções.
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

### Resolução de tema

`NewThemeFromEnv(appPrefix)` é o atalho recomendado: lê
`<PREFIX>_DMS_SETTINGS` (padrão `~/.config/DankMaterialShell/settings.json`)
e `<PREFIX>_ACCENT` (padrão `mauve`) e chama `ResolveTheme`. Para controle
total, `ResolveTheme(settingsPath, fallbackAccent)` lê o `settings.json` do
DankMaterialShell instalado e devolve o hex do accent que o DMS está
renderizando hoje (mesmo lookup que o próprio DMS faz). Se o DMS não está
instalado/for outro tema, cai no accent Catppuccin passado.

### Footer e HelpModal

- `NewFooter(bindings...)` monta a barra de status/help. O lado direito é
  gerado das bindings, então os hints nunca divergem do que `key.Matches`
  aceita. `Status()` define o contexto esquerdo; `Render(width, theme)`
  devolve a linha pronta. Hints que não cabem são dropados por token (` · `).
- `NewHelpModal(sections...)` cria um overlay centralizado com as bindings em
  seções. Bind em `"?"` (ou outra tecla), chame `Update()` (retorna `true`
  quando o modal está aberto — o app não deve processar as teclas dele
  enquanto isso), `SetSize()` no resize, e renderize por último no `View`.
  Seções podem usar `BindingsFn` pra buscar os bindings vivos a cada render.

### Keybindings customizáveis (config-file-first)

O `KeyRegistry` centraliza as keybindings de um app e faz do arquivo de config
a **fonte da verdade**: defaults registrados em código + overrides do usuário
em `<ConfigDir()>/<app>/keybindings.json`. Crie, registre as ações e carregue:

```go
reg := tuiui.NewKeyRegistry(tuiui.ConfigPath("meuapp", "keybindings.json"))
reg.RegisterMany(
	tuiui.Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
	tuiui.Action{ID: "nav", Help: "mover", Keys: []string{"ctrl+h", "ctrl+l", "ctrl+j", "ctrl+k"}, Label: "ctrl+h/j/k/l"},
)
if err := reg.Load(); err != nil { /* arquivo corrompido: usa defaults */ }
```

O formato é `{"bindings": {"<id>": ["tecla", ...]}}`, onde `<id>` é o `ID` da
`Action`; o arquivo só existe enquanto houver pelo menos um override. O fluxo
primário é **editar o arquivo e recarregar** — exponha uma tecla de reload (ou
chame `Reload()` a cada `View()`) pra edição externa valer sem reiniciar:

```go
case key.Matches(msg, reg.Resolve("reload")): // ex.: "r"
	if changed, err := reg.Reload(); err != nil {
		status = theme.Error().Render("keybindings: " + err.Error())
	} else if changed {
		status = theme.Success().Render("keybindings recarregados")
	}
```

- **Dispatch**: `case key.Matches(msg, reg.Resolve("quit")):`
- **Footer**: `tuiui.NewFooter(reg.Bindings()...)` — reflete rebinds e reloads
  na hora.
- **Help modal**: passe `BindingsFn: reg.Bindings` na seção.
- **Hints agrupados (padrão canônico de navegação)**: uma única ação com
  várias teclas + `Label` (ex.: `nav` com `ctrl+h/l/j/k`, `Label:
  "ctrl+h/j/k/l"`) vira **uma linha** no footer/help em vez de N entradas; o
  dispatch por direção fica com o app (position-based), não com a tecla. Use
  `Label` quando as teclas compartilham o mesmo conceito.
- **SettingsModal (conveniência opcional)**: `tuiui.NewSettingsModal(reg)` —
  bind em `","`/`"s"` abre a lista de ações com rebind interativo
  (`enter` rebind, `r`/`R` resetam, conflitos em vermelho, custom marca com
  `●`). Com o fluxo de arquivo acima ele vira um atalho, não o único caminho.

### Config em TOML (config-file-first)

`Config[T]` é o equivalente do `KeyRegistry` pras demais preferências do app:
defaults compilados no código + overrides do usuário em
`<ConfigDir()>/<app>/config.toml`. Cada app define seu próprio `T`; a lib não
opina sobre o schema.

```go
type config struct {
	Editor string `toml:"editor"`
	Layout struct {
		SidebarWidth int `toml:"sidebar_width"`
	} `toml:"layout"`
}

var defaults = config{Editor: "nvim"} // defaults.Layout.SidebarWidth = 22, etc.

cfg := tuiui.NewConfig(tuiui.ConfigPath("meuapp", "config.toml"), defaults)
if err := cfg.Load(); err != nil { /* arquivo corrompido: mantém os defaults */ }

width := cfg.Get().Layout.SidebarWidth
```

O merge é por chave: o TOML só sobrescreve os campos que **aparecem no
arquivo**, e todo o resto fica com o valor de `defaults`. Slices são
substituídas inteiras (um `roots = [...]` no arquivo troca a lista toda, não
concatena). Um arquivo ausente não é erro — o app roda em defaults puros.

Assim como nas keybindings, o fluxo primário é **editar o arquivo e
recarregar**, na mesma tecla:

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

`Reload()` reporta se a config **efetiva** mudou, e um TOML malformado devolve
erro **preservando o valor anterior** — um typo no meio da edição não derruba
um app rodando pros defaults. Nem todo campo é recarregável a quente (path de
banco, número de workers): `Reload()` diz que mudou, e cabe ao app decidir o
que fazer com isso.

> `T` deve ser um struct de valores. Campos ponteiro/map/slice são
> compartilhados com o `defaults` até o arquivo sobrescrevê-los, então não
> mute o que o `Get()` devolve.

**`ConfigPath(app, file)`** resolve `~/.config/<app>/<file>` respeitando
`XDG_CONFIG_HOME` — use tanto pro `config.toml` quanto pro
`keybindings.json`, em vez de montar o path na mão.

### Helpers de layout

- `PadLines(s, width)` — pad/truncate ANSI-aware de cada linha pra largura
  exata (necessário pra bordas de panel alinharem).
- `WrapText(s, width)` — quebra linha em texto puro.
- `PadToHeight(s, lines)` — pad/trunca um bloco pra altura exata.

### Estilos semânticos

`Success` (verde), `Warning` (amarelo), `Error` (vermelho), `Info` (azul) e
`Muted` (cinza) seguem o guia de semântica do Catppuccin — para status que
cada app colore do seu jeito mas com o mesmo significado.

### IPC

```go
args, err := tuiui.ParseIPCArgs(os.Args[2:]) // método + filtros + --json
if err != nil {
	fmt.Fprintln(os.Stderr, "uso: meuapp ipc <método> [key=value...] --json")
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
// dispatch no args.Method...
os.Exit(tuiui.WriteJSON(resultado))
```

`EnvOr`, `ExpandHome`, `HomeDir` e `ConfigDir` fecham o resto do esqueleto que
todo app repetia.

## Desenvolvimento

```bash
go test ./...
```

## Changelog

Veja [CHANGELOG.md](CHANGELOG.md) para o histórico de versões.

## Apoie o projeto

- **Global**: [ko-fi.com/ianptkcs](https://ko-fi.com/ianptkcs)
- **Brasil (Pix)**: escaneie o QR abaixo ou copie o código

  <img src="pix-qr.png" alt="Pix QR" width="200" />

  <details><summary>Código Pix (copiar)</summary>

  ```
00020126580014BR.GOV.BCB.PIX01365ad933b0-dcdc-4525-a736-0759902aeec65204000053039865802BR5925Ian Patrick da Costa Soar6009SAO PAULO62140510tQA85x6Dov63041FB6
  ```

  </details>

## Licença

[GNU AGPL-3.0](LICENSE) — mesma licença dos TUIs que a usam.
