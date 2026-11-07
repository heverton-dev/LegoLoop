# 🗺️ PLAN: Normalização da Taxonomia das Skills Proprietárias do LegoLoop

> **Status:** Proposta para Revisão  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** Antigravity  
> **Iniciativa:** I03 — Governança e Taxonomia Canônica  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_NORMALIZACAO_TAXONOMIA_SKILLS.md](../intents/2026-11-07_INTENT_NORMALIZACAO_TAXONOMIA_SKILLS.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md](../specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_NORMALIZACAO_TAXONOMIA_SKILLS.json](../routes/2026-11-07_ROUTE_NORMALIZACAO_TAXONOMIA_SKILLS.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (AGENTS.md, regra 4: agentes de harness proibidos; subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Materialização de Skills e Manuais | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Injeção no Core Go e Testes TDD | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Recompilação e Sincronização Universal | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Materialização de Skills e Manuais

#### Ticket `TICK-E01.1`: Criar SKILL.md Canônicos em `skills/`
* **Spec Ref:** `AC-E01.1`
* **Alvo:** `skills/legoloop-piece-upgrader/SKILL.md`, `skills/legoloop-mega-pipeline/SKILL.md`, `skills/legoloop-loop-map/SKILL.md`
* **Descrição:** Materializar os arquivos `SKILL.md` nas novas pastas canônicas com frontmatter atualizado (`name: legoloop-*`), descrições claras e comandos contendo os aliases legados e novos.
* **Gate Mecânico:** `TestEmbeddedTemplatesMatchFiles` e validação de existência dos arquivos no disco.

#### Ticket `TICK-E01.2`: Redigir Manuais Operacionais em `docs/manuals/`
* **Spec Ref:** `AC-E01.2`
* **Alvo:** `docs/manuals/MANUAL_FUNCAO_PIECE_UPGRADER_LEGOLOOP.md`, `docs/manuals/MANUAL_FUNCAO_MEGA_PIPELINE_LEGOLOOP.md`, `docs/manuals/MANUAL_FUNCAO_LOOP_MAP_LEGOLOOP.md`
* **Descrição:** Escrever os 3 manuais operacionais ISO detalhando gatilhos, flags da CLI, fluxo determinístico e exemplos práticos.
* **Gate Mecânico:** Verificação de ausência de placeholders `{{...}}`.

---

### Épico 02 — Injeção no Core Go e Testes TDD

#### Ticket `TICK-E02.1`: Atualizar Testes de Instalação (Red-Green)
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`
* **Alvo:** `internal/cli/instalar_skills_test.go`
* **Descrição:** Atualizar asserções de contagem esperada de skills nativas de 7 para 10 no teste unitário, fazendo o teste falhar (Red).
* **Gate Mecânico:** Falha esperada comprovando mutação.

#### Ticket `TICK-E02.2`: Expandir `skills.FuncoesNativas` em Go (Green)
* **Spec Ref:** `AC-E02.1`, `AC-E02.2`
* **Alvo:** `skills/embed.go`
* **Descrição:** Adicionar as 3 novas funções no slice `FuncoesNativas`, fazendo a suíte passar (Green).
* **Gate Mecânico:** `go test -count=1 ./internal/cli -run TestInstalarSkills` com exit code 0.

---

### Épico 03 — Recompilação e Sincronização Universal

#### Ticket `TICK-E03.1`: Recompilar Binário e Regressão Global
* **Spec Ref:** `AC-E03.1`
* **Alvo:** `cmd/legoloop/`, `bin/legoloop.exe`
* **Descrição:** Recompilar o binário `bin/legoloop.exe`, rodar `go vet ./...` e `gofmt -l .`.
* **Gate Mecânico:** `go build -o bin/legoloop.exe ./cmd/legoloop`.

#### Ticket `TICK-E03.2`: Executar `legoloop instalar-skills` Real
* **Spec Ref:** `AC-E03.2`
* **Alvo:** Cofre `~/.aidd/canonical/skills/` e junctions nos harnesses do usuário
* **Descrição:** Disparar o comando de instalação para distribuir fisicamente as 10 skills nos harnesses detectados.
* **Gate Mecânico:** `bin/legoloop.exe instalar-skills --json` reportando status ok e 10 skills.

---

## 4. Handoff de Governança & Próximos Passos

* **Espelho Formal:** [`docs/routes/2026-11-07_ROUTE_NORMALIZACAO_TAXONOMIA_SKILLS.json`](../routes/2026-11-07_ROUTE_NORMALIZACAO_TAXONOMIA_SKILLS.json)
* **Execução:** Após validação deste PLAN e da ROUTE no Stop-Gate 03, inicia-se a execução estéril dos tickets.

---

🛑 **STOP-GATE 03 DE AUDITORIA HUMANA:**  
*O plano estruturado com 3 épicos e 6 tickets reflete fielmente o escopo desejado?*  
*Aguardando aprovação do usuário para dar início à execução dos tickets.*
