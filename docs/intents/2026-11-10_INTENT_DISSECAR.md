# 🎯 INTENT: Oitavo Verbo Canônico `dissecar` — Inspeção Anatômica Profunda, Engenharia Reversa Estática e Visualização Motion Graphics

> **Status:** APPROVED  
> **Data:** 2026-11-10  
> **Autor:** Antigravity AI + Engenheiro Líder (Sessão Socrática Factual Addy Osmani)  
> **Iniciativa:** I08 — Verbo e Processo de Dissecação Universal do LegoLoop  
> **Alvo:** `internal/dissecar`, `cmd/legoloop`, `internal/cli`, `internal/studio`, `skills/legoloop-dissecar`, `loops/legoloop-dissecar`  
> **Sessão:** intv_1794354050170134000  
> **Confiança da Entrevista:** 100% (Consolidado em Sessão Socrática Factual)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação

### 1.1 O Problema Real
O ecossistema LegoLoop foi construído sobre a fundição de peças autônomas governadas por 7 verbos principais: `fabricar` (criação), `matriz` (avaliação de aptidão), `listar` (inventário e vínculos), `refinar` (substituição atômica de passos), `remover` (desvinculação segura), `compor` (hierarquia de sub-loops) e `fechar` (canonização determinística em 9 etapas).

Entretanto, o ecossistema sofre de uma lacuna crítica de **visibilidade e inspeção prévia**:
1. **Opacidade de Peças Externas (Caixa-Preta):** Quando um desenvolvedor ou agente ingere um script (`.ps1`, `.sh`, `.py`), um servidor MCP, uma skill de terceiro ou um repositório clonado, não há um mecanismo determinístico para radiografar o alvo antes de sua adoção. Peças externas entram com stubs ocultos (`TODO`, `pass`), dependências de runtime não declaradas, variáveis de ambiente desconhecidas, chamadas de sistema inseguras (`exec`, `rm -rf`, `curl`) e vícios de linguagem que violam a taxonomia ASD-STE100 (Lei #15).
2. **Ausência de Engenharia Reversa Operacional:** O LegoLoop sabe compor e orquestrar nós declarados em grafos DAG (`pipeline.dag.json`), mas não possui um motor nativo capaz de ler código procedural arbitrário e sintetizar automaticamente o grafo de nós correspondente. O desenvolvedor é forçado a desenhar os nós manualmente no Studio ou escrever o JSON do zero.
3. **Defasagem Visual na Inspeção:** A avaliação de uma peça hoje é puramente textual ou tabular. Não existe representação gráfica dinâmica do fluxo de dados e da arquitetura da peça. Embora o projeto possua o acervo de referência `docs/videorea/ai-motion-graphics-10-repos/` (Three.js, PixelFX, Lemo Opuscar, React Rapier), esse arsenal permanece desconectado do motor principal do LegoLoop.

### 1.2 Quem Sofre e Impacto
* **Engenheiros de Plataforma e Arquitetos:** Perdem horas auditando manualmente linhas de scripts e pacotes para mapear entradas, saídas e dependências antes de promover peças ao Cofre (`~/.aidd/canonical/`).
* **Operadores do Studio Web:** Não conseguem visualizar a anatomia estrutural de uma peça selecionada na barra lateral, vendo apenas cards estáticos sem percepção tridimensional da topologia.
* **Agentes de IA nos Harnesses:** Sofrem com estouro de contexto ao tentar ler arquivos extensos (> 500 linhas) sem antes consultar uma ficha anatômica compacta e estruturada (violando o princípio Graph-First da Lei #14).

### 1.3 O Porquê: A Solução Definitiva
A criação do 8º verbo canônico, **`dissecar`**, preenche de forma definitiva essa lacuna arquitetural:
* **No CLI Go:** Executa análise estática de alta velocidade (< 50ms) sem executar código de terceiro, extraindo árvore de nós, dependências, variáveis de ambiente, métricas e conformidade de governança.
* **Na Engenharia Reversa:** Desconstroi a lógica sequencial do alvo e projeta imediatamente um Grafo DAG executável (`pipeline.dag.json`).
* **No Arsenal Visual:** Gera a Tríade Visual (`motion_graph.json`, `anatomia.html`, `anatomia.svg`) utilizando os estilos catalogados em `docs/videorea/ai-motion-graphics-10-repos` com o Design System estrito do LegoLoop (fundo `#131312`, laranja `#d9480f`, zero emojis e vetores proprietários).
* **No Studio Web:** Disponibiliza o botão contextual de bisturi em SVG, modal de dissecação direta, visualização 3D Canvas interativa e 4 botões de ação imediata para o ciclo de vida.

### 1.4 Estado Atual vs. Estado Desejado
* **Estado Atual (Hoje):**
  * Códigos e scripts externos são tratados como caixas-pretas opacas.
  * Ingestão de peças exige auditoria manual e propensa a erros humanos.
  * Criação de DAGs a partir de código legado requer transcrição manual.
  * O acervo de motion graphics em `docs/videorea/` está inativo no produto principal.
* **Estado Desejado (Amanhã):**
  * O comando `legoloop dissecar <alvo>` disseca qualquer peça, script ou repositório em milissegundos.
  * Relatórios padronizados são persistidos em `docs/anatomias/<slug>/` com Ficha Anatômica JSON e manual técnico Markdown.
  * A engenharia reversa forense sintetiza um `pipeline.dag.json` válido pronto para execução via `legoloop run`.
  * A Tríade Visual gera visualizações interativas em Canvas Three.js/PixelFX autônomas e sem acoplamento a binários pesados de vídeo no CLI.
  * Quatro ações de ciclo de vida (`Instalar no Cofre`, `Executar DAG`, `Elevar para Padrão Ouro`, `Exportar Pacote`) fecham o fluxo no Studio.

---

## 2. Escopo Inegociável

### 2.1 O que ESTÁ no Escopo (Objetivos Primários):
* [ ] **Objetivo 1 (Engine `internal/dissecar`):** Criar motor estático puro em Go stdlib com Path Jail, AST parser para Go e scripts, detector de chamadas ao sistema, extrator de variáveis de ambiente, cálculo de métricas (LOC, complexidade ciclomática, tokens) e linter de governança contra stubs e termos proibidos.
* [ ] **Objetivo 2 (Engenharia Reversa de DAG):** Implementar sintetizador de grafos que desconstroi o alvo e gera um arquivo canônico `pipeline.dag.json` válido conforme o schema v2.2.0, aprovado pelo validador mecânico `legoloop validate`.
* [ ] **Objetivo 3 (Tríade Visual de Motion Graphics):** Gerar três artefatos visuais desacoplados por dissecação baseados nos 10 repositórios de `docs/videorea/`:
  1. `motion_graph.json`: Cena de nós e conexões consumível pelo Studio Web.
  2. `anatomia.html`: Bundle autônomo offline com Three.js/Canvas embutido para inspeção sem servidor.
  3. `anatomia.svg`: Vetor animado via CSS keyframes para inclusão em Markdown e livros LegoPress.
* [ ] **Objetivo 4 (Design System Rígido da Marca):** Aplicar nos gráficos exclusivamente as cores canônicas do LegoLoop (`#131312`, `#d9480f`, `#2b8a3e`, `#c92a2a`), tipografia Anton / Cascadia Code, ícone vetorial proprietário do bloco LegoLoop e **zero emojis** em todo o código, UI e saídas.
* [ ] **Objetivo 5 (Subcomando CLI Go):** Criar `legoloop dissecar <alvo>` com flags `--profundidade <rasa|completa|forense>`, `--formato <json|md|html|svg|dag|todos>`, `--motion`, `--saida <dir>` e `--forcar`, com renderização de tabela ANSI/Lipgloss de alto contraste e modo headless JSON puro.
* [ ] **Objetivo 6 (Studio Web & Rotas REST):** Implementar rotas `POST /api/dissecar` e `GET /api/dissecar/cache/:hash`, botão de bisturi em SVG proprietário nos cards de peças, modal de dissecação direta, aba dedicada "Anatomia" com Canvas 3D e 4 botões de ação operacional.
* [ ] **Objetivo 7 (Governança Canônica):** Criar skill universal `skills/legoloop-dissecar/SKILL.md` (comunicação em ASD-STE100 puro sem marcações artificiais) e o loop de esteira `loops/legoloop-dissecar/pipeline.dag.json`.

### 2.2 O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos execução dinâmica de código inspecionado (proibido rodar código de terceiro no SO fora de sandbox seguro).
* ⛔ Não incluiremos dependências pesadas de vídeo (`ffmpeg`, `puppeteer`, `chromium`) no binário Go do CLI (gravação de vídeo MP4 opcional roda exclusivamente no navegador do usuário via WebCodecs).
* ⛔ Não adicionaremos dependências externas de terceiros no `go.mod` (Go stdlib puro para garantir compilação hermética).
* ⛔ Não utilizaremos emojis em nenhuma linha de código, elemento de UI ou relatório gerado.
* ⛔ Não alteraremos o código-fonte original do alvo dissecado sem confirmação humana explícita.
* ⛔ Não refatoraremos peças legadas fora do escopo estrito da dissecação.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Bloco / Subloop | Papel na Esteira de Dissecação |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `unit_task` | `internal/dissecar/resolver` | Resolução do caminho do alvo, isolamento via Path Jail e detecção de tipo de peça. |
| `BLOCO-02` | `unit_task` | `internal/dissecar/extrator` | Parsing de AST, contagem de linhas/tokens, mapeamento de imports, chamadas de sistema e variáveis de ambiente. |
| `BLOCO-03` | `unit_task` | `internal/dissecar/governanca` | Varredura de stubs (`TODO`, `pass`) e auditoria taxonômica contra `internal/taxonomy`. |
| `BLOCO-04` | `unit_task` | `internal/dissecar/sintese` | Materialização de Grafo DAG (`pipeline.dag.json`) e Tríade Visual (`motion_graph.json`, `anatomia.html`, `anatomia.svg`). |
| `BLOCO-05` | `unit_task` | `internal/cli/dissecar` | Subcomando CLI com parser de flags e renderizador de tabela ANSI de alto contraste. |
| `BLOCO-06` | `ui_task` | `internal/studio` | Handlers HTTP REST, aba de Anatomia com Canvas 3D e ações de ciclo de vida. |
| `BLOCO-07` | `skill` | `skills/legoloop-dissecar` | Instrução universal de agente em ASD-STE100 puro para todos os harnesses. |
| `BLOCO-08` | `loop` | `loops/legoloop-dissecar` | Grafo de processo DAG acíclico validado mecanicamente com `legoloop validate`. |
| `BLOCO-09` | `script` | `go test -count=1 ./...` | Suíte de testes com 100% de sucesso e zero regressões. |
| `BLOCO-10` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca e instalação no PATH do binário oficial. |

---

## 4. Governança Operacional e Restrições Inegociáveis

* **Gatilho Operacional de Disparo (Trigger):** `manual` (acionável por desenvolvedores no CLI, operadores no Studio Web ou agentes autônomos via chat).
* **Nível de Risco & Blast Radius:** Baixo (opera em modo somente-leitura; gravações restritas à pasta isolada `docs/anatomias/<slug>/`).
* **Matriz de Risco FMEA de Negócio:**

| Risco / Modo de Falha | Severidade (S) | Ocorrência (O) | Detecção (D) | RPN | Barreira Mecânica de Mitigação |
| :--- | :---: | :---: | :---: | :---: | :--- |
| Alvo malicioso com tentativa de traversal de diretório (`../../etc/passwd`) | 9 | 3 | 1 | **27** | Prisão mecânica de *Path Jail* via `filepath.Clean` e barreira de contenção no diretório raiz. |
| Execução inadvertida de script perigoso durante dissecação | 10 | 2 | 1 | **20** | O motor é 100% estático baseado em AST lexer (zero chamadas de execução `exec.Command`). |
| Travamento do CLI por geração pesada de assets gráficos | 7 | 2 | 1 | **14** | Desacoplamento estrito: o Go só emite JSON, HTML e SVG textuais. O render 3D roda no navegador. |
| Inconsistência nos nomes de campos da Ficha Anatômica | 6 | 3 | 1 | **18** | Validação automática com JSON Schema formal (`schemas/anatomia.v1.json`). |
| Poluição visual por emojis em interfaces técnicas | 5 | 2 | 1 | **10** | Linter e testes automatizados que barram emojis e exigem SVGs proprietários. |

* **Definição Mecânica de Concluído (Definition of Done):**
  1. Suíte de testes unitários em `internal/dissecar`, `internal/cli` e `internal/studio` passando com 100% de aprovação e testes de mutação certificados.
  2. Nenhuma função de produção com mais de 35 linhas (`tools/funclen`).
  3. Formatação estrita sem diffs em `gofmt -l .`.
  4. Validação formal do grafo DAG gerado via `legoloop validate`.
  5. Recompilação fresca obrigatória do binário em `bin/legoloop.exe` e `go install ./cmd/legoloop`.
  6. Emissão de relatório de fechamento em `docs/reports/` com balanço auditável de tempo e telemetria.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Resolução estática segura de qualquer peça (Skill, DAG, Script, Hook, MCP, Código Go). | Testes unitários do pacote `internal/dissecar` | `go test -count=1 ./internal/dissecar/...` |
| 2 | Subcomando CLI funcional com todas as flags e saída ANSI formatada. | Testes de integração do subcomando | `go test -count=1 ./internal/cli -run TestDissecar` |
| 3 | Síntese de Grafo DAG correspondente por engenharia reversa no nível forense. | Validação de integridade do DAG gerado | `go run ./cmd/legoloop validate docs/anatomias/teste/pipeline.dag.json` |
| 4 | Tríade de Motion Graphics gerada com Design System LegoLoop e sem emojis. | Inspeção dos arquivos JSON, HTML e SVG | Verificação de cores `#131312` e ausência de emojis |
| 5 | Rotas REST e interface no Studio para disparo e inspeção 3D interativa. | Testes das rotas HTTP no Studio | `go test -count=1 ./internal/studio -run TestDissecarAPI` |
| 6 | Governança canônica fechada (Skill, Loop, Manual e Recompilação). | Checklist de conformidade de templates | `go test -count=1 ./docs/templates/...` |

---

## 6. Handoff de Governança

* **Próximo Documento Mandatório:** [`docs/specs/2026-11-10_SPEC_DISSECAR.md`](../specs/2026-11-10_SPEC_DISSECAR.md)
* **Gate de Transição:** Este Intent foi aprovado em sessão socrática factual com 100% de confiança, autorizando o início da especificação técnica arquitetural detalhada.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*Intenção aprovada formalmente. Prosseguindo para a SPEC técnica arquitetural.*