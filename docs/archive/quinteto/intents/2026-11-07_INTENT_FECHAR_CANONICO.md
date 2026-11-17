# 🎯 INTENT: Motor Canônico Universal de Fechamento e Canonização (legoloop fechar)

> **Status:** APPROVED  
> **Data:** 2026-11-07  
> **Autor:** Ecossistema LegoLoop  
> **Iniciativa:** Canonização Universal do Ciclo de Vida do LegoLoop  
> **Alvo:** `internal/fechar`, `internal/cli`, `skills/legoloop-fechar`, `loops/fechar`, `AGENTS.md`  
> **Sessão:** 9408084d-aa3e-40a3-a924-418ca1f67855  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** A conclusão de implementações no LegoLoop (como pipelines de vídeo, novas ferramentas e processos) dependia da memória humana ou de disciplina de prompt para que manuais, livros didáticos LegoPress, skills universais nos harnesses, recompilação de binários e reindexação de grafo fossem disparados. O motor `internal/fabricar` gerava apenas a Tríade e a skill básica, sem executar a cadeia final em 8 passos, gerando risco de "entregas pela metade" e skills invisíveis nos chats.
* **Quem Sofre:** Todos os desenvolvedores, agentes e operadores do ecossistema que utilizam pipelines gerados e precisam que eles sejam universais, documentados e compilados imediatamente.
* **O Porquê:** Garantir o princípio inegociável do LegoLoop: "Tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de um LLM ou à memória humana".
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O encerramento de um pipeline requer intervenção manual do operador ou orquestração manual do agente para emitir livro, compilar PDF, instalar skills em harnesses, dar rebuild e comitar.
  * **Amanhã:** Todo pipeline de criação, fabricação ou evolução compõe deterministicamente o comando `legoloop fechar <slug>` (ou o sub-loop `fechar`) como seu nó final obrigatório, executando de forma atômica e mecânica as 8 etapas de canonização.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Criar o motor em Go `internal/fechar` com orquestração mecânica das 8 etapas de canonização universal.
* [ ] Objetivo 2: Expor o subcomando CLI `legoloop fechar <slug> [--yes] [--json]` e a skill nativa universal `skills/legoloop-fechar/SKILL.md` (distribuída a todos os harnesses).
* [ ] Objetivo 3: Criar o grafo DAG acíclico `loops/fechar/pipeline.dag.json` e o manifesto `.loop.json` para permitir composição como sub-loop final de qualquer processo.
* [ ] Objetivo 4: Acoplar o gancho final de fechamento no motor `internal/fabricar`, no `internal/videorea` e no `loops/mega-pipeline.loop.json`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto (Go stdlib pura).
* ⛔ Não faremos commits automáticos sem a aprovação explícita no Stop-Gate humano (`--yes` ou confirmação no terminal).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD e contratos. |
| `BLOCO-02` | `unit_task` | `internal/fechar` | Motor de orquestração mecânica das 8 etapas em Go via TDD. |
| `BLOCO-03` | `unit_task` | `internal/cli/fechar.go` | Subcomando CLI `legoloop fechar` com flags e telemetria STE-PT. |
| `BLOCO-04` | `subloop` | `loops/fechar` | Grafo DAG acíclico executável e composível em outros loops. |
| `BLOCO-05` | `script` | `go test -count=1 ./...` | Regressão mecânica global, funclen, go vet e gofmt. |
| `BLOCO-06` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca e instalação no PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `loop_end: "*"` e `manual` (`legoloop fechar <slug>`).
* **Nível de Risco & Blast Radius:** Baixo (workspace local com escrita atômica e rollback em caso de falha).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente Local (Go stdlib pura, zero dependências externas, teto zero tokens de LLM para as etapas determinísticas).
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Execução de `legoloop fechar <slug>` valida e gera o Manual, compila o Livro em PDF via LegoPress, vincula a Skill nos harnesses detectados, valida o DAG, recompila o binário e verifica todos os portões de qualidade com exit code 0.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  - O slug do processo existe ou possui pelo menos a Tríade básica emitida.
  - As ferramentas `pandoc` e `typst` estão instaladas no SO para compilação editorial de livros.
* **Riscos e Perguntas em Aberto:** Nenhum. A arquitetura de cada etapa já foi provada isoladamente nas demandas anteriores.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green). | Teste novo em `internal/fechar` falhando antes do código | `go test -v ./internal/fechar` |
| 3 | Execução idempotente do fechamento para pipelines existentes | Executar contra `video-motion` e verificar saída | `legoloop fechar video-motion --dry-run` |
| 4 | Recompilação fresca do binário garantida como etapa intrínseca | Verificação de data/hora do binário em `bin/legoloop.exe` | `legoloop version` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_FECHAR_CANONICO.md`](../specs/2026-11-07_SPEC_FECHAR_CANONICO.md)
* **Gate de Transição:** Este Intent foi aprovado com 100% de confiança pelo operador no chat.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende à sua expectativa?*  
*Aprovado pelo operador para avanço à SPEC (`TEMPLATE_SPEC.md`).*
