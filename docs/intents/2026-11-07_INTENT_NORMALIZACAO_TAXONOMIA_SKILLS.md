# 🎯 INTENT: Normalização da Taxonomia das Skills Proprietárias do LegoLoop

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Autor:** Antigravity  
> **Iniciativa:** I03 — Governança e Taxonomia Canônica  
> **Alvo:** `skills/embed.go`, `internal/cli/instalar_skills.go`, `skills/`, `docs/manuals/`, `AGENTS.md`  
> **Sessão:** 38311f3b-f20c-4f72-b15e-6826ffe786cd  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Três skills operacionais que nasceram para o LegoLoop (`aidd-piece-upgrader`, `aidd-legoloop-mega-pipeline` e `aidd-loop-map`) ainda mantêm o prefixo obsoleto `aidd-*`. Isso gera estranheza e confusão de descoberta: o usuário busca pelas capacidades do LegoLoop (como elevar peças com `/upgrade` ou orquestrar pipelines) e não as encontra associadas à marca `legoloop-*`, assumindo que são resquícios legados de outro ecossistema.
* **Quem Sofre:** Usuários da CLI, agentes operando em múltiplos harnesses (Antigravity, Claude, Cursor, Windsurf) e desenvolvedores de loops.
* **O Porquê:** O LegoLoop já padronizou as suas 7 funções canônicas embutidas (`legoloop-fabricar`, `legoloop-matriz`, `legoloop-listar`, `legoloop-refinar`, `legoloop-remover`, `legoloop-compor`, `legoloop-legopress`). A consolidação das 3 skills proprietárias restantes unifica a taxonomia sob `legoloop-*`, fecha o cofre canônico em 10 funções nativas e permite distribuição automatizada via `legoloop instalar-skills`.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** As skills moram no cofre com nomes `aidd-*`, apenas 7 skills estão embutidas no binário `legoloop.exe` (`skills/embed.go`), e a elevação de peças (`piece-upgrader`) fica isolada e invisível na descoberta de comandos.
  * **Amanhã:** Todas as 10 skills proprietárias possuem taxonomia uniforme `legoloop-*`, estão embutidas no binário via `skills/embed.go`, mantêm aliases retrocompatíveis para os comandos no chat (`/upgrade-piece`, `/mega-pipeline`, etc.), possuem manuais operacionais em `docs/manuals/` e são distribuídas em 1 comando com `legoloop instalar-skills`.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Normalizar as 3 skills proprietárias para taxonomia `legoloop-*`:
  * `aidd-piece-upgrader` ➔ `legoloop-piece-upgrader`
  * `aidd-legoloop-mega-pipeline` ➔ `legoloop-mega-pipeline`
  * `aidd-loop-map` ➔ `legoloop-loop-map`
* [x] Embutir as 3 skills no repositório em `skills/` e expandir `skills/embed.go` de 7 para 10 funções nativas.
* [x] Criar os manuais operacionais canônicos em `docs/manuals/`:
  * `MANUAL_FUNCAO_PIECE_UPGRADER_LEGOLOOP.md`
  * `MANUAL_FUNCAO_MEGA_PIPELINE_LEGOLOOP.md`
  * `MANUAL_FUNCAO_LOOP_MAP_LEGOLOOP.md`
* [x] Atualizar a suíte de testes de instalação de skills (`internal/cli/instalar_skills_test.go`) e templates.
* [x] Garantir retrocompatibilidade total de comandos no chat via aliases frontmatter.
* [x] Distribuir e sincronizar as skills no cofre e harnesses via `legoloop instalar-skills`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não renomearemos peças de terceiros/workhorses de execução de tarefas dentro das DAGs (ex.: `spec`, `tdd`, `rea`, `grill` do cofre continuam válidas como peças de execução).
* ⛔ Não quebraremos contratos públicos de CLI existentes (`legoloop upgrade`, `legoloop matriz`, etc.).
* ⛔ Não alteraremos a lógica interna de cálculo de reputação ou diagnósticos do `legoloop upgrade`.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/legoloop-spec` | Elaborar especificação técnica detalhada com contratos e critérios de aceite BDD. |
| `BLOCO-02` | `unit_task` | Código de produção Go & Markdown | Criação dos SKILL.md, manuais e atualização de `skills/embed.go`. |
| `BLOCO-03` | `script` | `go test ./internal/cli/... ./docs/templates/...` | Execução de testes de regressão dos embeds e conformidade. |
| `BLOCO-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória do binário LegoLoop. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (solicitação do usuário no chat).
* **Nível de Risco & Blast Radius:** Baixo (arquivos de documentação de skills e lista embutida no Go; zero impacto na estabilidade da CLI).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente local, stdlib Go pura, sem requisições de rede.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. `legoloop instalar-skills --json` reporta 10 skills instaladas e vinculadas aos harnesses.
  2. `go test -count=1 ./internal/cli/...` e `go test -count=1 ./docs/templates/...` verdes.
  3. Comandos `/legoloop-upgrade` e `/upgrade-piece` apontam para `legoloop-piece-upgrader`.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | `skills.FuncoesNativas` contém exatamente as 10 funções canônicas. | Teste unitário em Go | `go test -count=1 ./internal/cli/...` |
| 2 | As 3 skills existem em `skills/legoloop-*/SKILL.md` com frontmatter válido. | Leitura do FS embutido | `go test -count=1 ./internal/cli -run TestInstalarSkills` |
| 3 | Binário recompilado e atualizado no PATH local. | Build do binário sem erros | `go build -o bin/legoloop.exe ./cmd/legoloop` |
| 4 | Distribuição nos harnesses executada com sucesso. | Execução da CLI | `go run ./cmd/legoloop instalar-skills` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md`](../specs/2026-11-07_SPEC_NORMALIZACAO_TAXONOMIA_SKILLS.md)
* **Gate de Transição:** Este Intent foi formalizado com 100% de confiança após alinhamento explícito com o usuário. Procede-se para a emissão da SPEC formal.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*Intenção aprovada pelo usuário via comando explícito no chat.*
