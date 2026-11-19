# 03 — Tickets de Remediação e Melhoria

> **Alvo:** cmd/legoloop, internal/cli e governança do verbo Matriz / Matriciar
> **Data:** 2026-11-18
> **Auditor:** Agente Auditor LegoLoop
> **Metodologia:** BDD / EARS / Zero Stubs

---

## Ticket TICK-MAT-01: Registro Canônico do Comando `matriciar` na CLI com Aliases Retrocompatíveis
- **Alvo:** `internal/cli/cli.go:110` e `internal/cli/matriz.go:22, 95-98`.
- **Objetivo:** Adicionar `"matriciar": cmdMatriz` na tabela principal de comandos de `internal/cli/cli.go`, preservando `"matriz"` e `"matrix"` como aliases retrocompatíveis funcionais. Atualizar o texto de uso canônico para priorizar `matriciar`.
- **Critério de Aceite (EARS):**
  - Dado que o usuário execute `legoloop matriciar ver --json`;
  - Quando a CLI despachar o comando através de `cli.Run`;
  - Então o resultado deve ser exatamente idêntico ao comando `legoloop matriz ver --json`, com exit code 0;
  - E dado que scripts legados executem `legoloop matriz <subcomando>`, a execução deve continuar funcionando sem regressão.
- **Gate Mecânico:** `go test -count=1 -run "Matriz" ./internal/cli/...` com exit code 0.

---

## Ticket TICK-MAT-02: Atualização do Dicionário Canônico ASD-STE100 e Linter Taxonômico
- **Alvo:** `docs/vocab/2026-11-10_ASD_STE100_CANONICAL_DICTIONARY.json` e `internal/taxonomy/`.
- **Objetivo:** Registrar o verbo canônico `matriciar` no domínio de Governança e Peças, explicitando que `Matriz` é a entidade/substantivo e `matriciar` é o verbo de ação correspondente.
- **Critério de Aceite (EARS):**
  - Dado o dicionário canônico ASD-STE100;
  - Quando o validador de taxonomia for executado (`go run ./tools/taxlinter`);
  - Então o termo `matriciar` deve ser reconhecido como termo canônico aprovado para ações de catalogação e consulta de reputação na matriz.
- **Gate Mecânico:** `go run ./tools/taxlinter internal cmd loops skills docs/templates` com exit code 0.

---

## Ticket TICK-MAT-03: Sincronização dos 7 Verbos na Governança Viva
- **Alvo:** `AGENTS.md:124-132`, `README.md` e diagramas de arquitetura de fluxo de verbos.
- **Objetivo:** Atualizar a documentação do Ciclo dos 7 Verbos para apresentar formalmente: Fabricar, Matriciar, Listar, Refinar, Remover, Compor e Fechar, destacando a preservação do substantivo `Matriz` para a tabela de aptidão.
- **Critério de Aceite (EARS):**
  - Dado o arquivo `AGENTS.md`;
  - Quando inspecionada a seção dos 7 Verbos de Processo;
  - Então todos os 7 verbos devem constar harmonizados no infinitivo, com indicação do comando `legoloop matriciar ver`.
- **Gate Mecânico:** Validação de links e conformidade de formatação sem tags HTML em tabelas.

---

## Ticket TICK-MAT-04: Disponibilização da Skill Canônica e Slash Command
- **Alvo:** `skills/legoloop-matriciar/SKILL.md` e distribuição para os harnesses.
- **Objetivo:** Criar a definição canônica da skill `legoloop-matriciar` espelhando a funcionalidade de `legoloop-matriz` e mantendo esta última como referência legada para compatibilidade de agentes.
- **Critério de Aceite (EARS):**
  - Dado o comando `legoloop instalar-skills`;
  - Quando executado no ambiente;
  - Então a skill `legoloop-matriciar` deve ser vinculada aos diretórios de skills de Claude, Antigravity, Cursor e Windsurf.
- **Gate Mecânico:** `legoloop instalar-skills --dry-run --json` validado com exit code 0.
