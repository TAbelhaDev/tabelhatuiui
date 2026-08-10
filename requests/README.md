# Requests de features

Inbox de pedidos de features/componentes pro TabelaTuiUI, feitos por outros
agents. Cada pedido é um arquivo markdown em `requests/`.

## Como pedir

1. Leia o template: `requests/_template.md`.
2. Copie o template para `requests/<AAAAMMDD>-<feature>.md` (data do dia +
   slug curto; o formato evita colisão entre agents trabalhando em paralelo).
3. Preencha e deixe o pedido lá — quem estiver implementando vai passar por
   aí, e a ordem de implementação segue a prioridade declarada em cada arquivo.

## Lifecycle

| Estado    | Onde vive                                             |
| --------- | ----------------------------------------------------- |
| ativo     | `requests/<AAAAMMDD>-<feature>.md`                    |
| concluído | `requests/done/<AAAAMMDD>-<feature>.md` (movido)      |
| recusado  | `requests/done/<AAAAMMDD>-<feature>.md` com `wontfix` |

- Ativo vira concluído **movendo o arquivo** pra `requests/done/`.
- O commit que implementa referencia o arquivo do request no corpo da mensagem.
- Recusado (`wontfix`) também é movido pra `requests/done/`, e fica o motivo
  no corpo.
- Antes de criar um request, confira se já existe um parecido em `requests/`
  (evite duplicar) — o duplicado deve se anexar ao existente.

## Onde isso aponta

- `requests/_template.md` — o que preencher em cada arquivo.
- Ao implementar: seguir o padrão do repo (ver `README.md` — API/uso, e
  `CHANGELOG.md` — bump + entry).

## Abrir request pelo GitHub (uso externo)

Usuários do TabelaTuiUI publicado **fora deste repo** (um terceiro, sem o
clonado localmente) abrem o request como **issue de feature no GitHub**:

```bash
gh issue create --repo TabelaDev/tabelatuiui --template feature_request.yml
```

O template `feature_request.yml` funciona no formulário do GitHub via
"New Issue → Feature request". Issues de feature criadas pelo GitHub entram
no mesmo fluxo: vão para `requests/done/` quando implementadas e o commit de
implementação referencia a issue no corpo da mensagem.
