# AGENTS.md

Este arquivo é lido por agents que trabalham neste repo. Leia também o
`README.md` antes de qualquer mudança — ele concentra as decisões de design
(tema, chrome, convenções) e a API exposta.

## Requests de features

Pedidos de novos componentes/features chegam pelo dir `requests/`:

1. Leia `requests/README.md` (fluxo) e `requests/_template.md` (o que preencher).
2. Se for criar um pedido novo: copie o template para
   `requests/<AAAAMMDD>-<feature>.md` e preencha.
3. Se for implementar: veja se já existe um request para a feature, siga os
   critérios de aceite dele, e referencie-o no commit. Ao concluir, mova o
   arquivo para `requests/done/`.

Uso externo (terceiro usando o TabelaTuiUI publicado, sem repo local): o
pedido chega como GitHub issue via
`gh issue create --repo TabelaDev/tabelatuiui --template feature_request.yml`.

## Testes

- `go vet ./...`, `go test ./...` e `go build ./...` antes de qualquer commit.
