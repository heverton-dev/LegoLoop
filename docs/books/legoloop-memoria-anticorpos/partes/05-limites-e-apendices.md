# Parte V — Limites e Apêndices

## 1. Limites declarados

- O lock por arquivo cobre o mesmo host. Duas máquinas não compartilham lock.
- O gate anti-reincidência depende da assinatura do erro. Erro com texto variável gera assinaturas diferentes.
- A telemetria de `token_budget` só aparece para nós com orçamento declarado.
- O zeramento da sessão cobre a sessão dona do Fechamento. Sessões alheias seguem vivas.

## 2. Apêndice A — Divergências anotadas na integração

| Tema                    | Registro                                                 |
| :---------------------- | :------------------------------------------------------- |
| Contagem de tickets     | O brief cita 37 tickets; o plano paralelo tem 11 tickets e a ROUTE tem 15 nós |
| Nome de erro            | Termo único `ErrLockOcupado` (Lei 15); divergência `ErrLockAtivo` resolvida em 2026-11-17 |
| Pasta de Processos      | Os schemas citam `processos/<slug>`; o acervo real vive em `processos/<slug>` |

## 3. Apêndice B — Portões de qualidade

- `go test -count=1 ./...`
- `go vet ./...`
- `gofmt -l .` com saída vazia
- `go run ./tools/funclen internal` sem estouro novo
- `go run ./tools/taxlinter internal cmd loops skills docs/templates` sem violação nova
- `git diff --exit-code -- go.mod go.sum`

## 4. Apêndice C — Artefatos da demanda

| Artefato                | Caminho                                                             |
| :---------------------- | :------------------------------------------------------------------ |
| Manual operacional      | `docs/manuals/2026-11-16_MANUAL_FUNCAO_MEMORIA_ANTICORPOS_LEGOLOOP.md`     |
| Skill canônica          | `skills/legoloop-memoria-anticorpos/SKILL.md`                         |
| Livro didático          | `docs/books/memoria-anticorpos/`                                      |
| Relatório de integração | `docs/reports/2026-11-16_RELATORIO_INTEGRACAO_MEMORIA_ANTICORPOS.md`      |
