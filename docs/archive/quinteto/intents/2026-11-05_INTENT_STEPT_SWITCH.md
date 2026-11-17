# 🎯 INTENT: Alternador Universal e Granular de Governança, Economia de Tokens e Output STE-PT para Harnesses

> **Status:** APPROVED  
> **Data:** 2026-11-06 (revalidado e aprovado)  
> **Autor:** Antigravity (Pair Programming com Usuário)  
> **Alvo:** `internal/ste/`, `internal/cli/`, `cmd/legoloop/`, `C:\Users\trcnologia\Desktop\stept-cli`  
> **Sessão:** 78979240-b03f-4f85-8d01-66cd6312f633  
> **Confiança da Entrevista:** 100% (concluída e validada via método Addy Osmani / Mega Pipeline)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Os agentes de IA operando nos diversos harnesses do mercado (Gemini CLI / Antigravity, Claude Code, OMP, OpenCode, Mimo) geram saídas verbosas, raciocínios prolixos e executam mutações no disco (edições, deleções, comandos de terminal) sem fornecer um contrato visual claro de risco prévio nem plano de reversão. Além disso, consomem cota excessiva de tokens em prompts desnecessariamente longos e logs desproporcionais. O protocolo **ASD-STE100 / STE-PT** (Português Técnico Simplificado com Cartões de Ação) foi implementado com sucesso no núcleo do LegoLoop (`internal/ste`), mas permaneceu restrito às mensagens da CLI do produto e ausente no chat dos agentes por falta de uma chave de ativação universal, dinâmica e composível.
* **Quem Sofre:** O desenvolvedor e engenheiro de software que utiliza múltiplos harnesses (individualmente ou em orquestrações multi-agência como ORCA-AIDD), ficando refém de alucinações, saídas inchadas de tokens, queima desnecessária de créditos e ações com efeitos destrutivos não autorizados.
* **O Porquê:** Padronizar um pipeline determinístico de economia de tokens e controle de risco cobrindo a Tríade do Turno da IA (**Input ➔ Thinking ➔ Output**), permitindo ao desenvolvedor ligar, desligar, calibrar e mesclar modos (Caveman Ultra + Cartão STE-PT) com latência sub-5ms em qualquer harness instalado no sistema operacional.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Regras estáticas manuais residem nos harnesses (ex.: `caveman-global.md`). Não há descoberta dinâmica de ferramentas instaladas. O protocolo STE-PT não atinge o chat dos agentes. Não existe mecanismo de CLI para ligar, desligar, calibrar a rigidez das respostas ou mesclar Caveman e STE-PT de forma unificada.
  * **Amanhã:** Um comando CLI ultra-rápido (`legoloop stept` e o binário standalone nativo em Go `stept.exe` em `Desktop\stept-cli`) que gerencia `~/.stept/config.json`, detecta dinamicamente harnesses no SO via scanner e sincroniza atomicamente as regras de system prompt dos harnesses detectados com uma matriz de modos (`off | stept | caveman | hibrido`) e política de pausa (`sempre | destrutivo | nunca`).

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [x] Objetivo 1: **Motor de Regras, Persistência e Scanner Dinâmico de Harnesses:**
  * Reuso do motor de descoberta de harnesses (`internal/assets`) para varrer dinamicamente ferramentas instaladas (Gemini CLI / Antigravity, Claude Code, OMP, OpenCode, Mimo).
  * Subcomando `scan` e flag `--rescan` para atualizar o mapa de ferramentas quando novos harnesses forem instalados no SO.
  * Repositório central de configuração em `~/.stept/config.json`, com bloco `padrao` e sobreposições granulares opcionais por harness.
* [x] Objetivo 2: **Governança da Tríade de Economia de Tokens (Input ➔ Thinking ➔ Output):**
  * **Input (Entrada):** Regra contratual explícita que instrui a IA a isolar comandos imperativos e substantivos de ação, eliminando ruído cortês e saudações; hook local determinístico (sem LLM) onde suportado nativamente pelo harness (ex.: Claude Code).
  * **Thinking (Raciocínio):** Injeção compulsória do bloco **Caveman Ultra** no raciocínio interno (frases telegráficas, sem artigos, proibição de carregar logs inteiros no contexto).
  * **Output (Resposta):** Matriz modular de modos:
    * `off`: Desativa a governança e remove o arquivo de regra do harness.
    * `stept`: Cartão de Ação estruturado ASD-STE100 em 100% das saídas de ação.
    * `caveman`: Respostas curtas, secas e telegráficas.
    * `hibrido` *(Padrão de Ouro)*: Caveman no Thinking + Resposta Concisa no Chat + Cartão STE-PT obrigatório antes de qualquer ação com efeito.
  * **Política de Pausa:** `sempre` (obriga confirmação `[1] seguir [2] personalizar [0] cancelar`), `destrutivo` (pausa somente em ações irreversíveis) e `nunca` (renderiza o card como contrato de intenção mas executa autonomamente).
* [x] Objetivo 3: **Sincronizador Universal e Não-Invasivo de Harnesses:**
  * Injeção atômica de `stept-rules.md` nos diretórios oficiais de regras de cada harness suportado:
    * Gemini CLI / Antigravity (`~/.gemini/config/rules/stept-rules.md`)
    * Claude Code (`~/.claude/rules/stept-rules.md`)
    * OMP (`~/.omp/rules/stept-rules.md`)
    * OpenCode (`~/.config/opencode/rules/stept-rules.md`)
    * Mimo (`~/.mimocode/rules/stept-rules.md`)
  * **Política Não-Invasiva:** Não criar pastas fantasmas para harnesses ausentes (salvo `--create-dirs`). No modo `off`, remover estritamente o arquivo `stept-rules.md`, preservando 100% das outras regras do usuário (ex.: `caveman-global.md`).
* [x] Objetivo 4: **Subcomando Nativo no LegoLoop (`legoloop stept`):**
  * Integrado ao binário Go existente, com suporte a `status`, `set`, `on`, `off`, `scan` com flags `--harness`, `--modo`, `--pausa`, `--all`, `--create-dirs`.
* [x] Objetivo 5: **Projeto Standalone Comercial em Go (`Desktop\stept-cli`):**
  * Criação da pasta `C:\Users\trcnologia\Desktop\stept-cli`, com código Go independente espelhando o núcleo do LegoLoop, testes unitários e compilação do executável `stept.exe`.
  * Comando `stept install` para configuração automática de PATH no sistema operacional:
    * Windows: `%LOCALAPPDATA%\stept\bin` via API de usuário.
    * Linux / macOS: `~/.local/bin` com injeção segura no shell profile.
  * Documentação comercial completa: `README.md` (visão comercial, instalação universal de 1 linha) e `MANUAL.md` (referência operacional de comandos e flags).

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não usaremos "wrappers" CLI invasivos que sequestram ou substituem binários originais dos harnesses no PATH.
* ⛔ Não faremos chamadas a APIs pagas de LLM para resumir prompts (pré-processamento é 100% local e determinístico).
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além da standard library em Go (zero dependências de terceiros em ambos os projetos).
* ⛔ Não tocaremos no cofre real `~/.stept` durante a execução de testes automatizados (testes sempre herméticos em `t.TempDir()`).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/aidd-spec` | Elaborar especificação formal com critérios BDD (`2026-11-05_SPEC_STEPT_SWITCH.md`). |
| `BLOCO-02` | `skill` | `skills/aidd-legoloop-mega-pipeline` | Decomposição formal em Plano de Tickets (`PLAN`) e Grafo Executável (`ROUTE.json`). |
| `BLOCO-03` | `unit_task` | `internal/ste/switch/` | Implementação do motor de configuração, persistência, scanner e injeção de regras via TDD. |
| `BLOCO-04` | `unit_task` | `internal/cli/stept.go` | Implementação do subcomando CLI `legoloop stept` e tabela de status no terminal. |
| `BLOCO-05` | `unit_task` | `C:\Users\trcnologia\Desktop\stept-cli` | Materialização do projeto standalone em Go no Desktop, documentação completa e compilação de `stept.exe`. |
| `BLOCO-06` | `unit_task` | Testes de Mutação (Regra 8) | Validação por mutação deliberada nos pontos críticos de herança, prompt e remoção atômica. |
| `BLOCO-07` | `script` | `go test ./...` | Suíte de regressão mecânica, staleness e linters (`funclen`, `gofmt`). |
| `BLOCO-08` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória do LegoLoop no PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (invocado via CLI `legoloop stept` ou `stept`).
* **Nível de Risco & Blast Radius:** Baixo (opera exclusivamente gravando arquivos de regras nos diretórios conhecidos de configuração de harnesses do usuário e no diretório de dados `~/.stept/`).
* **Política de Privacidade & Orçamento Cognitivo:** Exclusivamente Local. Zero dados enviados à rede. Modelos Flash/Low-effort para código e testes; raciocínio focado apenas em contratos e invariantes.
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**
  1. Comando `legoloop stept status` exibindo a tabela limpa com o estado e origem de cada harness no terminal.
  2. Arquivo de regras `stept-rules.md` criado e removido fisicamente conforme comandos `on` e `off`.
  3. Comando `scan` identificando novas ferramentas instaladas no SO sem reiniciar a máquina.
  4. Binário standalone `Desktop\stept-cli\bin\stept.exe` executando de forma independente e produzindo o mesmo comportamento.
  5. Suíte `go test -count=1 ./...` passando com zero falhas, `funclen` limpo (≤ 35 linhas) e 100% das mutações detectadas.
* **Premissas Assumidas:**
  - Os harnesses reconhecem regras em formato markdown colocadas em seus respectivos subdiretórios de regras de usuário.
  - A linguagem Go stdlib atende com latência sub-5ms toda a operação de I/O de regras no disco.
  - A instalação de PATH no escopo do usuário não requer privilégios de administrador em nenhum dos sistemas operacionais suportados.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test -count=1 ./...` com exit code 0 |
| 2 | Cobertura de testes mantida com TDD estrito (Red-Green). | Teste unitário de cada método do motor falhando antes da implementação | Histórico de execução de testes |
| 3 | Matriz de 2 eixos granular por harness salva e recuperada fielmente com herança padrão vs override. | Teste de persistência JSON com herança padrão vs override | Asserts em `internal/ste/switch_test.go` |
| 4 | Sincronização atômica e não-destrutiva de regras nos harnesses detectados. | Presença e formato do template de prompt injetado no disco; remoção limpa no `off` | Verificação de escrita e deleção física de `stept-rules.md` |
| 5 | Scanner dinâmico de harnesses funcional com comando `scan`. | Teste de detecção com diretórios mockados em `t.TempDir()` | Verificação de lista de harnesses atualizada |
| 6 | Projeto autônomo em Go funcional no Desktop com instalador de PATH, manual e executável compilado. | Execução de `stept.exe status` na pasta do Desktop | `C:\Users\trcnologia\Desktop\stept-cli\bin\stept.exe status` retornando 0 |
| 7 | Nenhuma função ultrapassa 35 linhas de código (Lei do LegoLoop). | Linter `funclen` executado na pasta `internal/` | `go run ./tools/funclen internal` com 0 violações |
| 8 | Testes de mutação deliberada comprovam a eficácia dos testes nos 4 pontos críticos. | Falha imediata de teste sob mutação | Relatório de mutações |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-05_SPEC_STEPT_SWITCH.md`](../specs/2026-11-05_SPEC_STEPT_SWITCH.md)
* **Gate de Transição:** Este Intent foi aprovado pelo usuário no chat com confiança de 100%, autorizando formalmente a transição para o Loop 2 & 3 (atualização da SPEC, PLAN e ROUTE).

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA — STATUS:** `APROVADO` ✅  
*Intenção consolidada e validada formalmente pelo usuário em 2026-11-06.*
