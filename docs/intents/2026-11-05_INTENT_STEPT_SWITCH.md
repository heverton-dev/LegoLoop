# 🎯 INTENT: Alternador Global e Granular de Output STE-PT para Harnesses

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-05  
> **Autor:** Antigravity (Pair Programming com Usuário)  
> **Alvo:** `internal/ste/`, `internal/cli/`, `cmd/legoloop/`, `C:\Users\trcnologia\Desktop\stept-cli`  
> **Sessão:** 8d960f94-04a7-4717-9a72-2b34281e3a80  
> **Confiança da Entrevista:** 100% (concluída via método Addy Osmani / `interview-me`)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Os agentes de IA operando nos diversos harnesses do mercado (Gemini CLI / Antigravity, Claude Code, OMP, OpenCode, Mimo) geram saídas verbosas, imprevisíveis e sem garantia de comunicação prévia clara de risco ou plano de reversão antes de aplicarem mutações no disco (edições, criação de arquivos, comandos de terminal). O protocolo **ASD-STE100 / STE-PT** (Português Técnico Simplificado com Cartões de Ação) foi implementado com sucesso no núcleo do LegoLoop (`internal/ste`), mas permaneceu restrito às mensagens da CLI do produto e ausente no chat dos agentes por falta de uma chave de ativação universal e dinâmica.
* **Quem Sofre:** O desenvolvedor e engenheiro de software que utiliza múltiplos harnesses (individualmente ou em orquestrações multi-agência como ORCA-AIDD), ficando refém de alucinações, saídas inchadas de tokens e ações com efeitos destrutivos não autorizados.
* **O Porquê:** Padronizar um contrato visual de alta previsibilidade, enxuto e seguro (Cartão ASD / STE-PT) antes de qualquer efeito colateral, reduzindo o consumo de tokens de saída, permitindo reversibilidade auditável e garantindo controle total por harness (ex.: estrito e pausado no Claude Code, híbrido no Antigravity, desligado no Mimo).
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Apenas regras estáticas manuais residem nos harnesses (ex.: `caveman-global.md`). O protocolo STE-PT não atinge o chat dos agentes. Não existe mecanismo de CLI para ligar, desligar ou calibrar a rigidez das respostas por harness.
  * **Amanhã:** Um comando CLI ultra-rápido (`legoloop stept` e o binário standalone `stept.exe` em `Desktop\stept-cli`) que gerencia `~/.stept/config.json` e sincroniza atomicamente as regras de system prompt dos harnesses detectados com uma matriz de 2 eixos (Modo de Saída: `off | hibrido | estrito` e Política de Pausa: `sempre | destrutivo | nunca`), operando tanto em escopo global quanto granular por harness.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: **Motor de Regras e Persistência Granular:** Criar o modelo de dados e repositório de configuração em `~/.stept/config.json`, com bloco `padrao` e personalização por harness (`gemini`, `claude`, `omp`, `opencode`, `mimo`).
* [ ] Objetivo 2: **Matriz de Configuração de Dois Eixos:**
  * **Modo de Saída:** `off` (desativa/remove regra), `hibrido` (card ASD / STE-PT apenas em ações com efeito) e `estrito` (100% das mensagens em formato STE-PT).
  * **Política de Pausa:** `sempre` (obriga confirmação `[1] seguir [2] personalizar [0] cancelar`), `destrutivo` (pausa somente em ações irreversíveis) e `nunca` (renderiza o card como contrato mas executa de forma autônoma).
* [ ] Objetivo 3: **Sincronizador Universal de Harnesses:** Injetores/extratores atômicos para os caminhos de regras oficiais de cada harness suportado:
  * Gemini CLI / Antigravity (`~/.gemini/config/rules/stept-rules.md`)
  * Claude Code (`~/.claude/rules/stept-rules.md` ou `.claude/CLAUDE.md`)
  * OMP (`~/.omp/rules/stept-rules.md` ou equivalente)
  * OpenCode (`~/.config/opencode/rules/stept-rules.md`)
  * Mimo (`~/.mimocode/rules/stept-rules.md`)
* [ ] Objetivo 4: **Subcomando Nativo no LegoLoop (`legoloop stept`):** Integrado ao binário Go existente, com suporte a `set`, `status`, `on`, `off` com flags `--harness`, `--modo`, `--pausa`, `--all`.
* [ ] Objetivo 5: **Projeto Standalone Comercial no Desktop:** Criação da pasta `C:\Users\trcnologia\Desktop\stept-cli`, com código Go independente, testes unitários, compilação do executável `stept.exe`, `README.md`, `MANUAL.md` e documentação de arquitetura comercial.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além da standard library em Go (zero dependências de terceiros).
* ⛔ Não faremos chamadas a APIs pagas de LLM ou redes externas; a sincronização é 100% local no disco.
* ⛔ Não tocaremos no cofre real `~/.stept` durante a execução de testes automatizados (testes sempre herméticos em `t.TempDir()`).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD (`2026-11-05_SPEC_STEPT_SWITCH.md`). |
| `BLOCO-02` | `skill` | `skills/aidd-legoloop-mega-pipeline` | Decomposição formal em Plano de Tickets (`PLAN`) e Grafo Executável (`ROUTE.json`). |
| `BLOCO-03` | `unit_task` | `internal/ste/switch/` | Implementação do motor de configuração, persistência e injeção de regras nos harnesses via TDD. |
| `BLOCO-04` | `unit_task` | `internal/cli/stept.go` | Implementação do subcomando CLI `legoloop stept` e tabela de status no terminal. |
| `BLOCO-05` | `unit_task` | `C:\Users\trcnologia\Desktop\stept-cli` | Materialização do projeto standalone no Desktop, documentação completa e compilação do binário `stept.exe`. |
| `BLOCO-06` | `script` | `go test ./...` | Suíte de regressão mecânica, testes de staleness e linters (`funclen`, `gofmt`). |
| `BLOCO-07` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória do LegoLoop no PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (invocado via CLI `legoloop stept` ou `stept`).
* **Nível de Risco & Blast Radius:** Baixo (opera exclusivamente gravando arquivos de regras nos diretórios conhecidos de configuração de harnesses do usuário e no diretório de dados `~/.stept/`).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente Local. Zero dados enviados à rede. Modelos Flash/Low-effort para código e testes; raciocínio focado apenas em contratos e invariantes.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. Comando `legoloop stept status` exibindo a tabela limpa com o estado de cada harness no terminal.
  2. Arquivos de regras `stept-rules.md` criados e removidos fisicamente conforme comandos `on` e `off`.
  3. Binário standalone `Desktop\stept-cli\bin\stept.exe` executando de forma independente e produzindo o mesmo comportamento.
  4. Suíte `go test -count=1 ./...` passando com zero falhas e cobertura completa de mutação.
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**
  - Os harnesses reconhecem regras em formato markdown colocadas em seus respectivos subdiretórios de regras de usuário.
  - A linguagem Go stdlib atende com latência sub-5ms toda a operação de I/O de regras no disco.
* **Riscos e Perguntas em Aberto:**
  - Caso um harness não esteja instalado na máquina, o comando deve registrar "não instalado / ausente" no status e não criar pastas fantasmas desnecessárias, a menos que forçado com `--create-dirs`.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura de testes mantida com TDD estrito (Red-Green). | Teste unitário de cada método do motor falhando antes da implementação | Histórico de execução de testes |
| 3 | Matriz de 2 eixos granular por harness salva e recuperada fielmente. | Teste de persistência JSON com herança padrão vs override | Asserts em `internal/ste/switch_test.go` |
| 4 | Sincronização atômica de regras nos 5 harnesses (Gemini, Claude, OMP, OpenCode, Mimo). | Presença e formato do template de prompt injetado no disco | Verificação de escrita e remoção de arquivo |
| 5 | Projeto autônomo funcional no Desktop com manual e executável compilado. | Execução de `stept.exe status` na pasta do Desktop | `C:\Users\trcnologia\Desktop\stept-cli\bin\stept.exe status` retornando 0 |
| 6 | Nenhuma função ultrapassa 35 linhas de código (Lei do LegoLoop). | Linter `funclen` executado na pasta `internal/` | `go run ./tools/funclen internal` com 0 violações |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md`](../specs/2026-11-05_SPEC_STEPT_SWITCH.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat, com confiança da entrevista em 100%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 100% de certeza factual atende integralmente à sua expectativa para o alternador do protocolo ASD / STE-PT?*  
*Somente após a sua aprovação explícita, avançaremos para o Loop 2 (Geração da SPEC v2 com critérios BDD).*
