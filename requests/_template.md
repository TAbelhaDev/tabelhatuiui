# <Nome curto da feature>

> Copie este arquivo para `requests/<AAAAMMDD>-<feature>.md` e preencha.
> Referências de contexto: `README.md` (API/uso), `CHANGELOG.md` (histórico).

- **autor**: <agent/ferramenta que pediu>
- **data**: <YYYY-MM-DD>
- **prioridade**: low | medium | high

## O que

<descrição da feature: o que ela faz, pra quem, onde aparece>

## Por que / contexto

<qual app/lib/casos de uso motivam; como ficaria sem isso>

## API proposta

```go
// como os apps consumidores vão usar — snippet de uso real
```

<novos tipos/métodos/args, nomes que já têm convenção no repo>

## Escopo / fora de escopo

- dentro: <o que este request cobre>
- fora: <o que NÃO faz, pra ficar explícito>

## Critérios de aceite

- [ ] <comportamento verificável>
- [ ] `go vet ./...`, `go test ./...` e `go build ./...` passam

## Docs a atualizar

- [ ] README (API/uso, se expõe algo novo)
- [ ] CHANGELOG (entry novo)
