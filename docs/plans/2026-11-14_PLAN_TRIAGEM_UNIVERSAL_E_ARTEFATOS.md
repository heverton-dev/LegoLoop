# PLAN — Peça de Triagem Universal e Artefatos do Fabricar

- **Data:** 2026-11-14
- **SPEC:** `docs/specs/2026-11-14_SPEC_TRIAGEM_UNIVERSAL_E_ARTEFATOS.md`

## Tickets

- **T1 — Motor `internal/triagem` (Híbrido Determinístico + Probabilístico):**
  - Criar schema `triagem.v1`, parser/validador e motor de decisão (determinístico primeiro, LLM só na ambiguidade).
  - Testes TDD unitários e de mutação (herméticos, sem tocar `~/.aidd`).

- **T2 — Atualizar `internal/fabricar` e `internal/processo` para 13 Artefatos:**
  - Remover `template` de `plannedFiles` (`internal/fabricar/camada4.go`).
  - Adicionar projeção e emissão automática de `triagem.json`, `manual`, `livro` e `fichas_tecnicas` (por nano-processo/nó).
  - Atualizar testes de `internal/fabricar` e `internal/processo`.

- **T3 — Pontos de Entrada (CLI, Chat/Skill e Studio):**
  - Expor `legoloop triagem` em `cmd/legoloop`.
  - Expor `POST /api/triagem` e entrada no Studio (`internal/studio`).
  - Atualizar template de `SKILL.md` para instruir triagem na entrada.

- **T4 — Saneamento e Migração de `loops/`:**
  - Implementar `legoloop migrar-loops [--dry-run] [--apply]` para auditar os 10 diretórios em `loops/`, remover legados inúteis/quebrados e gerar `triagem.json` + artefatos faltantes.
  - Executar migração e validar com `legoloop validate`.

- **T5 — Portões de Qualidade e Fechamento:**
  - Rodar `go test -count=1 ./...`, `go vet ./...`, `gofmt -l .`, `go run ./tools/funclen internal`, `go build -o bin/legoloop.exe ./cmd/legoloop`.
