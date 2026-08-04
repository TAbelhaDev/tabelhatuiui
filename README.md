<div align="center">

# tabelatuiui

**Tema + chrome compartilhado dos TUIs Bubble Tea do ianptkcs.**

[![Go Version](https://img.shields.io/github/go-mod/go-version/ianptkcs/tabelatuiui?style=flat-square&logo=go&logoColor=white&color=00ADD8)](go.mod)
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
| [djobs](https://github.com/ianptkcs/djobs) | TUI de jobs agendados (systemd) | migrado |
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
	"path/filepath"

	"github.com/ianptkcs/tabelatuiui"
)

var theme = tuiui.ResolveTheme(
	tuiui.EnvOr("MEUAPP_DMS_SETTINGS", filepath.Join(tuiui.HomeDir(), ".config", "DankMaterialShell", "settings.json")),
	tuiui.EnvOr("MEUAPP_ACCENT", "mauve"),
)

func main() {
	header := theme.Header(80).Render("meu app")
	footer := theme.Footer(80).Render("navegar · q sair")
	panel := theme.Panel(true).Render(tuiui.PadLines("conteúdo", 40))
	modal := theme.Modal().Render("tem certeza?")
	_ = header + footer + panel + modal
}
```

### Resolução de tema

`ResolveTheme(settingsPath, fallbackAccent)` lê o `settings.json` do
DankMaterialShell instalado e devolve o hex do accent que o DMS está
renderizando hoje (mesmo lookup que o próprio DMS faz). Se o DMS não está
instalado/for outro tema, cai no accent Catppuccin passado (padrão `mauve`).
O prefixo de env é escolha de cada app — a lib é agnóstica.

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

## Licença

[GNU AGPL-3.0](LICENSE) — mesma licença dos TUIs que a usam.
