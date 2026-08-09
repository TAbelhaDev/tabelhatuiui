# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.2.0] - 2026-08-09

### Added

- `NewThemeFromEnv(appPrefix)` — resolve o tema lendo `<PREFIX>_DMS_SETTINGS`
  e `<PREFIX>_ACCENT` do ambiente, substituindo o boilerplate de `theme.go`
  que todo app repetia.
- `Footer` — builder da barra de status/help: contexto à esquerda + hints de
  teclado à direita, gerados a partir de uma slice de `key.Binding` (nunca
  mais drift entre o que o `Update` aceita e o que o footer mostra). O
  overflow é truncado por token (` · `) em vez de caractere.
- `HelpModal` — overlay centralizado listando todas as keybindings em seções
  (o "?" de todo TUI), com scroll e fechamento por `q`/`esc`. Só precisa do
  tema; a lib cuida do posicionamento.
- `Binding` — alias para `key.Binding` (bubbles/key), exposto pelo pacote
  pra apps não precisarem importar a bubbles diretamente pra declarar atalhos.

### Changed

- `tabelatuiui` agora depende de `bubbles`/`bubbletea` (para os tipos
  `key.Binding` usados por `Footer`/`HelpModal`).

## [0.1.1] - 2026-08-09

### Changed

- `Header` agora usa o background do accent (Primary) com texto Base, em vez
  de mantle + texto accent — o cabeçalho fica visível contra terminal escuro.
- `Modal` passou a usar borda do accent (Primary) no lugar de lavender.
- `Header`/`Footer` ficam flush contra as bordas do terminal, com padding
  horizontal de 2 colunas (removido o antigo padding interno que os apps
  precisavam compensar).
