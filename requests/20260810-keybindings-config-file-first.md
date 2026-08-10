# Keybindings config-file-first

> Rebind editando o arquivo de config e recarregando, em vez de depender do
> modal interativo — com hints agrupados por conceito (padrão herdr).

- **autor**: opencode (djobs)
- **data**: 2026-08-10
- **prioridade**: medium

## O que

Trocar o fluxo de "rebind via modal" por **config-file-first**:

1. O `KeyRegistry` já persiste em `<ConfigDir()>/<app>/keybindings.json`.
   Tornar esse arquivo o fluxo primário: documentar o formato, e expor uma
   tecla de reload (ex.: `r` dentro de um modal de ajuda/config, ou re-ler o
   arquivo a cada `View()`).
2. Suportar **agrupamento de hints**: uma ação única com várias teclas e
   `Label` (ex.: `nav` com `ctrl+h/j/k/l`, `Label: "ctrl+h/j/k/l"`) já existe,
   mas o dispatch por direção fica por conta do app (position-based). Vale
   documentar esse padrão como o jeito canônico de navegação.
3. Manter o `SettingsModal` como conveniência opcional (abre o arquivo/lista
   as ações), não como o único caminho.

## Por que / contexto

O `SettingsModal` (ativada por `,`) é o único jeito de customizar keybindings
nos TUIs que usam o `KeyRegistry`, mas:

- Editar à mão o arquivo final (`{"bindings": {"<id>": ["tecla", ...]}}`) não é
  um fluxo suportado — só existe via modal, e não há tecla de "recarregar
  config".
- O herdr usa outro padrão: cada app tem um arquivo de config com as
  keys/settings, você edita o arquivo e recarrega a config. Sem modal, sem
  "rebind key" interativo.
- Com muitas ações (ex.: djobs tem ~20), o modal vira uma lista longa e o
  footer/help herdam hints por ação em vez de hints agrupados por conceito
  (ex.: `ctrl+h/j/k/l move focus` em vez de 4 entradas separadas).

## API proposta

```go
// Fluxo: usuário edita ~/.config/<app>/keybindings.json e recarrega.
reg := tuiui.NewKeyRegistry(filepath.Join(tuiui.ConfigDir(), "meuapp", "keybindings.json"))
reg.RegisterMany(
    tuiui.Action{ID: "quit", Help: "sair", Keys: []string{"q", "ctrl+c"}},
    tuiui.Action{ID: "nav", Help: "mover", Keys: []string{"ctrl+h", "ctrl+l", "ctrl+j", "ctrl+k"}, Label: "ctrl+h/j/k/l"},
)
reg.Load() // re-ler o arquivo a cada View()/tecla de reload
```

- Sem assinatura nova obrigatória: `Load()` já existe; a questão é *quando*
  chamar + documentar o formato + agrupar hints.
- Se preciso, um método de "reload" explícito no `KeyRegistry` pra apps
  deixarem o arquivo como fonte da verdade.

## Escopo / fora de escopo

- dentro: config-file como fluxo primário (documentação + reload), hints
  agrupados como padrão canônico de navegação
- fora: `fsnotify`/auto-reload contínuo (dependência nova, re-scan frequente);
  mudança de formato do arquivo (mantém `{"bindings": ...}`)

## Critérios de aceite

- [ ] Editar `keybindings.json` e recarregar reflete no footer/help/dispatch
- [ ] `nav` multi-tecla com `Label` agrupado aparece como uma linha no footer
- [ ] `go vet ./...`, `go test ./...` e `go build ./...` passam

## Docs a atualizar

- [ ] README (seção Keybindings customizáveis)
- [ ] CHANGELOG (entry novo)
