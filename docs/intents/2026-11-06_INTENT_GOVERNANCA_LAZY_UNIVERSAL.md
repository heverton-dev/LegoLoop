# 🎯 INTENT: Governança Lazy Universal de Ferramentas e Conhecimentos no LegoLoop

> **Status:** APPROVED  
> **Data:** 2026-11-06  
> **Autor:** Antigravity & Usuário  
> **Alvo:** `internal/mcp`, `internal/assets`, `cmd/legoloop`, `docs/routes/`  
> **Sessão:** 088c6d3f-4908-4fd5-b40e-896ee8fcf556  
> **Confiança da Entrevista:** 100% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Sessões em múltiplos harnesses de IA (MIMO Code, Claude Code, OpenCode) iniciam com inchaço massivo de 100.000 a 120.000 tokens de contexto logo no primeiro turno. A causa-raiz é a injeção cega de 207 ferramentas de engenharia reversa do MCP REA e de 596 skills redundantes formatadas no prompt do sistema.
* **Quem Sofre:** Todos os agentes e desenvolvedores operando no ecossistema local do usuário, com esgotamento prematuro de limites diários de tokens, alta latência em requisições simples e degradação de raciocínio cognitivo pelo ruído de contexto.
* **O Porquê:** Reduzir o contexto de largada de qualquer sessão de ~120k tokens para menos de 4k tokens, mantendo 100% das ferramentas e skills disponíveis sob demanda via chamadas dinâmicas.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O harness conecta diretamente em todos os servidores MCP e lê centenas de arquivos de skills do disco, enviando schemas pesados e cardápios gigantes no prompt de entrada.
  * **Amanhã:** O harness conecta exclusivamente ao Meta-MCP Router local via stdio e acessa o Skill Engine via busca BM25 com metadados de Tier e Natureza. Apenas o que for explicitamente requisitado pela conversa ou pelo nó do pipeline entra no contexto.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Objetivo 1: Implementar o subcomando `legoloop mcp-router` em Go puro via transporte `stdio`, expondo estritamente duas ferramentas no prompt: `mcp_search` e `call_mcp_tool`.
* [x] Objetivo 2: Implementar o motor de despacho `legoloop skill-router`, expondo `skill_search` com retorno enriquecido (ID, Nome, Score, Natureza, Tier, Categoria, Resumo) e `call_skill` sob demanda.
* [x] Objetivo 3: Criar comando de sincronização de harnesses (`legoloop sync-harnesses`) para distribuir as configurações canônicas para `~/.claude.json`, `.mimocode/mimocode.jsonc` e `.cursor/mcp.json`.
* [x] Objetivo 4: Integrar com ferramentas consagradas do ambiente (ex: AST e grafo semântico via `codebase-memory-mcp` existente, busca rápida com `ripgrep`), evitando reescrever do zero analisadores de código.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto (100% Go stdlib).
* ⛔ Não criaremos servidor HTTP daemon com portas de rede abertas na Fase 1 (arquitetura restrita a stdio nativo).
* ⛔ Não criaremos analisador AST próprio (usar o `codebase-memory-mcp` já integrado).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD Dado/Quando/Então. |
| `BLOCO-02` | `unit_task` | Código de produção Go (`internal/mcp`) | Implementação do Meta-MCP Router stdio e do Lazy Skill Dispatcher via TDD. |
| `BLOCO-03` | `script` | `go test -count=1 ./...` | Regressão mecânica e integridade em 100% dos pacotes Go. |
| `BLOCO-04` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória e teste de linha de comando. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (invocado pelos harnesses via stdio ou pelo usuário via CLI).
* **Nível de Risco & Blast Radius:** Baixo (operações em memória via stdio e alterações apenas em arquivos de configuração locais de harnesses com backup prévio).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente local, stdio puro, sem envio de telemetria externa e sem chaves de API trafegando fora do cofre seguro.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):** Uma sessão iniciada em qualquer harness (MIMO, Claude Code, Cursor) consome comprovadamente < 4.000 tokens de entrada no primeiro turno, mantendo a capacidade de chamar ferramentas do REA ou skills especializadas sob demanda.
* **Premissas Assumidas:** O transporte stdio é universalmente suportado pelo protocolo MCP em Windows; ferramentas existentes como `codebase-memory-mcp` permanecem funcionais em suas localizações canônicas.
* **Riscos e Perguntas em Aberto:** Manter a resposta do router < 50ms para não introduzir latência perceptível no arranque do harness.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e funclen limpo | `go test -count=1 ./...` e `go run ./tools/funclen internal` |
| 2 | Meta-MCP Router responde ao handshake MCP via stdio e lista apenas 2 tools. | Execução de teste unitário e E2E stdio | `go test -v ./internal/mcp` |
| 3 | `skill_search` retorna metadados completos de Tier, Natureza e Categoria. | Teste de contrato JSON | Validação de schema JSON na saída |
| 4 | Recompilação fresca bem sucedida e instalada. | Binário atualizado | `go install ./cmd/legoloop` e `legoloop --version` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md`](../specs/2026-11-06_SPEC_GOVERNANCA_LAZY_UNIVERSAL.md)
* **Gate de Transição:** Este Intent foi aprovado com 100% de confiança pelo usuário em sessão. Autorizado o avanço para a elaboração da SPEC técnica.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*Aprovado explicitamente pelo usuário na sessão 088c6d3f-4908-4fd5-b40e-896ee8fcf556.*
