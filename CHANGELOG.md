# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.4.0] - 2026-08-10

### Added

- `Config[T]` — settings de app em TOML (`<ConfigDir()>/<app>/config.toml`),
  seguindo o mesmo modelo config-file-first do `KeyRegistry`:
  `NewConfig(path, defaults)` + `Load()`/`Reload()`/`Get()`. O merge é por
  chave — o arquivo só sobrescreve os campos que aparecem nele, o resto fica
  no default — e cada app define o próprio schema (`T`). Arquivo ausente não é
  erro; arquivo malformado devolve erro **preservando o valor anterior**, pra
  um typo no meio da edição não derrubar um app rodando pros defaults.

- `ConfigPath(app, file)` — resolve `~/.config/<app>/<file>` respeitando
  `XDG_CONFIG_HOME`, substituindo o `filepath.Join(ConfigDir(), app, file)`
  que era copiado em cada app.

- `KeyRegistry.Reload()` — re-lê o arquivo de keybindings do disco (fonte da
  verdade) e reporta se os overrides mudaram; `Load()` agora delega pra ele.
  Edições externas no `keybindings.json` passam a valer chamando `Reload()`
  numa tecla de recarga (ou a cada `View()`), sem reiniciar.

### Changed

- Fluxo de keybindings documentado como **config-file-first** no README:
  editar `keybindings.json` + recarregar é o caminho primário, e `Label`
  multi-tecla (ex.: `nav` com `ctrl+h/j/k/l`) é o padrão canônico de hints
  agrupados — uma linha no footer/help em vez de uma entrada por tecla.
  `SettingsModal` fica como conveniência opcional, não o único caminho.

- `KeyRegistry` — fonte única de verdade das keybindings: defaults
  registrados em código + overrides por ação persistidos em JSON
  (`bindings.json`). `Resolve(id)` devolve o binding efetivo (override vence),
  `Bindings()` alimenta Footer/HelpModal, `Set` valida conflito de tecla,
  `Reset`/`ResetAll` restauram os defaults (e apagam o arquivo quando nada
  está customizado).
- `SettingsModal` — overlay de edição de teclas: navega pela lista de ações
  (j/k, g/G, pgup/pgdn), `enter` entra no modo "pressione a nova tecla",
  `r` reseta uma ação, `R` reseta tudo, mostra conflito em vermelho e marca
  com `●` as ações customizadas.
- `HelpSection.BindingsFn` — provider dinâmico de bindings pro HelpModal,
  pra ele refletir rebinds feitos no SettingsModal na hora.

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
