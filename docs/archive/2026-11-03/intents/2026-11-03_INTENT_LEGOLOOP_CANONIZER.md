# 🎯 INTENT: Fábrica Universal de Processos Canônicos (legoloop-canonizer)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-03  
> **Autor:** Antigravity AI & Everton  
> **Alvo:** `internal/canonizer/`, `internal/cli/`, `internal/studio/`, `docs/templates/`  
> **Sessão:** 1da38997-8926-43dc-a9c7-4c885fc12303  

---

## 1. Contexto e Motivação
* **O Problema Real:** A criação de novos processos, esteiras, loops e subetapas no ecossistema assistido por IA ainda ocorria de forma artesanal e fragmentada. Cada autor ou tarefa impunha convenções arbitrárias, gerando sobrecarga cognitiva, risco de alucinação de modelos, falta de testes padronizados e incompatibilidade entre os meios de execução (CLI, Chat e Canvas Visual).
* **O Porquê:** Unificar a criação de qualquer processo sob uma máquina de complexidade $O(1)$ constante e fractal. Não importa se a escala é $1x$ (um nano-ticket de 5 linhas) ou $1.000.000x$ (um mega-pipeline enterprise), a fábrica matricial de processos opera sob as mesmas 5 camadas determinísticas invariáveis.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Criação manual ou semi-estruturada de artefatos de planejamento; dispersão de onde salvar e como invocar em cada harness; ausência de garantia de prefixo uniforme.
  * **Amanhã:** A Fábrica de Processos (`legoloop-canonizer`) recebe uma intenção via CLI, Chat ou Studio, executa as 5 camadas canônicas (Clarificação ➔ Curadoria SO ➔ Fusão ➔ Tríade ➔ Template), crava o prefixo imutável `legoloop-<kebab-case>` e materializa o pacote interoperável pronto para rodar em qualquer um dos 3 meios.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] **1. Prefixo Fixo e Imutável:** Impor obrigatoriamente o prefixo `legoloop-` em formato `kebab-case` para todas as peças e processos gerados pela fábrica (ex: `legoloop-processo-do-joaozinho`).
* [ ] **2. A Esteira das 5 Camadas Canônicas:**
  * **Camada 1 (Clarificação):** Entrevista socrática estilo Addy Osmani adaptada em PT-BR até $\ge 95\%$ de confiança.
  * **Camada 2 (Curadoria Multi-Drive & Agnóstica a SO):** Scanner universal e agnóstico (Windows, Linux e macOS) que varre não apenas o disco principal do sistema (`C:\` ou `/`), mas **todos os discos locais, volumes montados e HDs externos / pen-drives USB espetados na máquina** (`D:\`, `E:\`, `F:\`, `/Volumes/*`, `/media/*`, `/mnt/*`), catalogando qualquer peça de IA disponível em qualquer dispositivo de armazenamento conectado.
  * **Camada 3 (Fusão Arquitetural):** Combinação semântica dos pontos fortes das peças curadas.
  * **Camada 4 (Tríade de Governança):** Geração de `SPEC.md`, `PLAN.md` e `ROUTE.json` com isolamento estéril, allowlist de ferramentas e gates booleanos.
  * **Camada 5 (Template SSOT):** Geração do gabarito contratual imutável de saída.
* [ ] **3. Interoperabilidade Tripartida Nativa:**
  * **Para a CLI:** Rota executável JSON compatível com `legoloop run legoloop-<nome>`.
  * **Para o Chat / Harness:** Manifesto `SKILL.md` instalado em diretório de skills para invocação por prompt ou slash command.
  * **Para o LegoLoop Studio:** Grafo `.loop.json` renderizável no Canvas e listado na barra lateral de Catálogo.
* [ ] **4. Subcomando CLI Go:** Implementar `legoloop canonize <nome>` com modos interativo e autônomo (`--yes`).
* [ ] **5. Interface no Studio:** Modal de Canonização com seletor de modelo cognitivo e plotagem automática no Canvas.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não alterar os contratos das peças existentes no cofre que já estejam no Padrão Ouro.
* ⛔ Não permitir remoção ou substituição do prefixo `legoloop-`.
* ⛔ Não misturar a *Fábrica de Processos* (Design-time) com a *Execução da Aplicação* (Run-time).

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `CANON-E01.1` | `unit_task` | Go Core (`internal/canonizer/`) | Motor de validação de nomenclatura e sanitização `legoloop-<kebab-case>`. |
| `CANON-E01.2` | `unit_task` | Go Core (`internal/canonizer/`) | Orquestrador das 5 camadas (Entrevista ➔ Curadoria ➔ Fusão ➔ Tríade ➔ Template). |
| `CANON-E02.1` | `unit_task` | CLI Go (`internal/cli/canonize.go`) | Subcomando `legoloop canonize` com flags `--interactive` e `--yes`. |
| `CANON-E03.1` | `skill` | `skills/legoloop-canonizer/SKILL.md` | Skill canônica instalada no SO para uso em Chat/NL com stop-gates. |
| `CANON-E04.1` | `ui_task` | Studio Web (`internal/studio/`) | Modal de Entrevista Socrática com seletor de modelos cognitivos e botão de plotagem no Canvas. |
| `CANON-E05.1` | `script` | `go test ./...` | Suíte completa de testes unitários e regressão TDD. |
| `CANON-E05.2` | `stop_gate` | `go build ./cmd/legoloop` | Recompilação fresca obrigatória e atualização no PATH. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** Tripartido (`cli: legoloop canonize`, `chat: /legoloop-canonizer`, `studio: botão Novo Processo`).
* **Nível de Risco & Blast Radius:** Baixo a Médio (gravação em diretórios canônicos dedicados e `docs/`).
* **Política de Privacidade & Orçamento Cognitivo:** Suporte a modelos locais e nuvem; effort clamping por camada (Clarificação: High, Curadoria: Low/Local, Tríade: Medium/High).
* **Contrato de Evidência Factual (Definition of Done):**
  1. O comando `legoloop canonize "teste do joaozinho" --yes` gera no disco o pacote tripartido `legoloop-teste-do-joaozinho` com `SKILL.md`, `.loop.json` e `ROUTE.json`.
  2. O processo aparece na paleta do Studio Canvas imediatamente após a materialização.
  3. Suíte de testes `go test ./internal/canonizer/...` passa com 100% verde e zero stubs.

---

## 5. Critérios de Sucesso de Alto Nível
1. Todas as invariantes da esteira passam com zero stubs e zero mocks vazios.
2. Prefixo `legoloop-` garantido por teste unitário em todos os fluxos de criação.
3. Binário recompilado e atualizado no PATH com sucesso (`go/bin/legoloop.exe`).

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-03_SPEC_LEGOLOOP_CANONIZER.md`](file:///C:/Users/trcnologia/Desktop/LegoLoop/docs/specs/2026-11-03_SPEC_LEGOLOOP_CANONIZER.md)
* **Gate de Transição:** Este Intent aguarda validação humana explícita (*Human-in-the-Loop Gate*) antes do início da redação da Spec técnica.
