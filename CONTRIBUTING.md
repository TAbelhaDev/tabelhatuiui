# Contribuindo com o tabelatuiui

Obrigado pelo interesse. Antes de abrir um PR de código, dá uma olhada no
`README.md` pra entender a proposta (tema + chrome compartilhado dos TUIs
Bubble Tea — não um framework de UI).

## Reportando bugs / sugerindo features

Pedidos de feature (inclusive de agents) seguem o fluxo `requests/`: veja
`requests/README.md` — copie `requests/_template.md` para
`requests/<AAAAMMDD>-<feature>.md` e preencha. Bugs continuam sendo issues de
bug no [GitHub](../../issues/new/choose).

## Enviando um PR

1. Fork o repositório.
2. Crie uma branch a partir de `main`.
3. Rode `go vet ./...`, `go test ./...` e `go build ./...` localmente antes
   de abrir o PR.
4. Abra o PR usando o template — descreva o quê e o porquê da mudança.

## Licença

Ao contribuir, você concorda que sua contribuição será licenciada sob a
[AGPL-3.0](LICENSE), a mesma licença do projeto.

## Código de conduta

Seja respeitoso. Críticas técnicas são bem-vindas; ataques pessoais não.
