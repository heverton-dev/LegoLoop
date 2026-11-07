# PARTE III — NÍVEL MICRO: ANATOMIA DOS CINCO SUB-LOOPS

Esta seção detalha o nível atômico do Mega Pipeline. Cada um dos cinco sub-loops é dissecado em suas etapas operacionais, apresentando os contratos rigorosos de entrada, processamento, ferramental acoplado, validação mecânica e handoff imutável.

---

## 1. Sub-Loop 1: Concepção & Alinhamento de Intenção

O primeiro sub-loop tem por objetivo extrair a verdadeira necessidade do usuário, eliminando a ambiguidade natural da linguagem humana antes de qualquer especificação técnica.

```{=typst}
#ficha(
  ("Papel", "Alinhamento socrático de escopo e premissas"),
  ("Disparo", "Linguagem natural do usuário via chat ou trigger manual"),
  ("Entrada", "Briefing inicial ou descrição da demanda"),
  ("Saída Principal", "docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md"),
  ("Portão de Saída", "STOP-GATE HUMANO 01 (Aprovação da Intenção)"),
)
```

```{=typst}
#esteira(
  no-claro("1.1 ENTREVISTA", sub: "aidd-grill socrático"),
  no("1.2 ESPELHO", sub: "template ISO"),
  no("1.3 HUMAN-GATE", sub: "parada obrigatória"),
)
```

### Etapa 1.1 — Entrevista Socrática (Briefing Ativo)
* **Consome:** Linguagem natural do usuário via chat ou comando CLI.
* **Peças:** Skill [`aidd-grill`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-grill/SKILL.md) + Hook de telemetria de tokens.
* **Processamento:** O agente assume uma postura socrática empática e incisiva, formulando uma pergunta por vez. Cada intervenção traz quatro marcadores formais: `HIPÓTESE`, `CONFIANÇA [X%]`, `PERGUNTA` afuniladora e `PALPITE` sugerido. O interrogatório se encerra quando a confiança atinge $\ge 95\%$ ou o modelo prevê com certeza as próximas respostas.
* **Gate:** Confiança calculada $\ge 95\%$ ou encerramento socrático acordado.
* **Handoff:** Transfere a estrutura de hipóteses validadas em memória para a Etapa 1.2.

### Etapa 1.2 — Geração do Espelho Canônico e Persistência
* **Consome:** Dados consolidados da Etapa 1.1.
* **Peças:** Script gerador de template ISO + Gate mecânico.
* **Processamento:** Materializa o documento formal estruturado em seis seções fixas: Resultado Esperado, Quem se Beneficia, Por que Agora, Critérios de Sucesso, Restrições Inegociáveis e Fora de Escopo.
* **Gate:** Expressão DSL: `file_exists('docs/intents/YYYY-MM-DD_INTENT_<TEMA>.md')`.
* **Handoff:** Grava o arquivo físico persistido no disco.

### Etapa 1.3 — Human-Gate 01 (Validação da Intenção)
* **Consome:** Documento de intenção gerado na Etapa 1.2.
* **Peças:** Controle interativo `human-gate`.
* **Processamento:** O runtime entra em pausa física obrigatória. Exibe o espelho canônico e aguarda o comando humano (`[1] seguir`, `[2] personalizar`, `[0] cancelar`).
* **Gate:** Confirmação explícita do usuário.
* **Handoff:** Destrava a execução para o Sub-Loop 2 ou Sub-Loop 3.

---

## 2. Sub-Loop 2: Engenharia Reversa & Design Pass *(Condicional)*

Acionado quando a demanda envolve recriar telas existentes, reconstruir interfaces legadas ou extrair design systems de produtos de terceiros.

```{=typst}
#ficha(
  ("Papel", "Captura de runtime, dissecção DOM e tokens de design"),
  ("Disparo", "Condicional: presença de URL, mock, binário ou tela legada"),
  ("Entrada", "URL ou arquivos de interface existentes"),
  ("Saída Principal", "DESIGN.md e TOKENS.css"),
  ("Portão de Saída", "STOP-GATE HUMANO 02 (Validação Visual)"),
)
```

```{=typst}
#esteira(
  no("2.1 INSPEÇÃO", sub: "Playwright / Obscura"),
  no-claro("2.2 NORMALIZAÇÃO", sub: "impeccable + OKLCH"),
  no("2.3 GATE DESIGN", sub: "validação de contraste"),
)
```

### Etapa 2.1 — Inspeção de Runtime & Captura DOM
* **Consome:** URL de referência externa, binário compilado ou mock de tela.
* **Peças:** MCP `playwright`, MCP `chrome-devtools` ou motor [`obscura`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/obscura/SKILL.md) em Rust.
* **Processamento:** Inicia o navegador em ambiente isolado, captura a árvore DOM semântica, calcula os estilos reais de CSS e realiza capturas de tela multi-viewport.
* **Gate:** `exit_code == 0 && file_exists('scratch/raw_dom.json')`.
* **Handoff:** Grava o payload JSON da estrutura da tela em `scratch/raw_dom.json`.

### Etapa 2.2 — Normalização e Extração de Tokens
* **Consome:** `scratch/raw_dom.json` da Etapa 2.1.
* **Peças:** Skill [`impeccable`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/impeccable/SKILL.md) + Script de paleta OKLCH.
* **Processamento:** Erradica 100% de emojis amadores, extrai tokens semânticos de espaçamento, tipografia e cores, e valida contrastes em conformidade com as normas WCAG.
* **Gate:** Zero emojis detectados e contraste de cores aprovado pelo script de validação.
* **Handoff:** Grava os artefatos `DESIGN.md` e `TOKENS.css`.

### Etapa 2.3 — Gate de Design & Human-Gate 02
* **Consome:** `DESIGN.md` e `TOKENS.css`.
* **Peças:** Controle `human-gate` + Script linter de design.
* **Processamento:** Apresentação visual da paleta e layout propostos para aprovação do desenvolvedor.
* **Gate:** Aval humano formal registrado no Studio ou terminal.
* **Handoff:** Libera os tokens aprovados para a especificação do Sub-Loop 3.

---

## 3. Sub-Loop 3: Engenharia & Tríade de Planejamento

É o coração arquitetural do LegoLoop. Ele converte a intenção humana e o design visual em um contrato de engenharia matematicamente executável.

```{=typst}
#ficha(
  ("Papel", "Especificação BDD, decomposição hierárquica e rota DAG"),
  ("Disparo", "Conclusão e aprovação dos Sub-Loops 1 e 2"),
  ("Entrada", "INTENT.md e opcionalmente DESIGN.md"),
  ("Saída Principal", "A Tríade: SPEC.md, PLAN.md e ROUTE.json"),
  ("Portão de Saída", "STOP-GATE HUMANO 03 (Auditoria da Tríade)"),
)
```

```{=typst}
#esteira(
  no-claro("3.1 SPEC BDD", sub: "aidd-spec + boundaries"),
  no-claro("3.2 TICKETS", sub: "fatiamento vertical"),
  no("3.3 ROTA JSON", sub: "compilador DAG"),
  no("3.4 STOP-GATE", sub: "bloqueio pré-código"),
)
```

### Etapa 3.1 — Especificação BDD (SPEC)
* **Consome:** `docs/intents/INTENT.md` (+ `DESIGN.md` se presente).
* **Peças:** Skill [`aidd-spec`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-spec/SKILL.md).
* **Processamento:** Redige as premissas técnicas explícitas, define as fronteiras em 3 níveis (*Sempre Fazer*, *Perguntar Antes*, *Nunca Fazer*), demarca as costuras de teste (*Test Seams*) e formaliza os critérios de aceite binários `AC-E{N}.{M}` no formato rigoroso Dado/Quando/Então.
* **Gate:** `file_exists('docs/specs/YYYY-MM-DD_SPEC_<TEMA>.md') && has_binary_criteria(...)`.
* **Handoff:** Grava o documento de especificação técnica no disco.

### Etapa 3.2 — Decomposição Hierárquica em Tickets (PLAN)
* **Consome:** `docs/specs/SPEC.md` da Etapa 3.1 *(Dependência estrita)*.
* **Peças:** Skill [`aidd-tickets`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-tickets/SKILL.md).
* **Processamento:** Fatiamento vertical da especificação em uma hierarquia de tickets: Épicos (>20k tokens de contexto), Micro-tickets (4k a 20k tokens para TDD pontual) e Nano-tickets (<4k tokens para correções de lint).
* **Gate:** `file_exists('docs/plans/YYYY-MM-DD_PLAN_<TEMA>.md')`.
* **Handoff:** Grava o plano unificado com tabela estruturada.

### Etapa 3.3 — Compilação da Rota Executável (ROUTE)
* **Consome:** `docs/plans/PLAN.md` e `docs/specs/SPEC.md`.
* **Peças:** Skill [`aidd-orchestrate-pipeline-creator`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-orchestrate-pipeline-creator/SKILL.md) + Script validador JSON Schema.
* **Processamento:** Converte a tabela legível do plano em um mapa formal JSON com declaração de IDs únicos, dependências, modelo roteado (`code-pro`, `code-fast`, `code-nano`), esforço (`--effort`) e comandos de verificação mecânica.
* **Gate:** Validação contra o schema oficial (`exit_code == 0`).
* **Handoff:** Grava `docs/routes/YYYY-MM-DD_ROUTE_<TEMA>.json`.

### Etapa 3.4 — Human-Gate 03 (Auditoria da Tríade)
* **Consome:** Tríade completa (`SPEC.md` + `PLAN.md` + `ROUTE.json`).
* **Peças:** Controle `human-gate`.
* **Processamento:** Parada obrigatória antes de tocar em qualquer arquivo da aplicação. O desenvolvedor valida se o plano reflete as decisões arquiteturais corretas.
* **Gate:** Aprovação explícita do usuário.
* **Handoff:** Destrava a inicialização da construção no Sub-Loop 4.

---

## 4. Sub-Loop 4: Construção TDD por Épico *(Orchestrix)*

É o sub-loop onde o código real de software é escrito. Ele executa cada micro-ticket isoladamente dentro de um ciclo rigoroso de Test-Driven Development.

```{=typst}
#ficha(
  ("Papel", "Construção de código de produção orientado a testes unitários"),
  ("Disparo", "Aprovação formal da Tríade de Planejamento"),
  ("Entrada", "ROUTE.json e âncoras da SPEC.md"),
  ("Saída Principal", "Código de produção testado e TDD_REPORT.md"),
  ("Portão de Saída", "STOP-GATE POR ÉPICO (Revisão Humana de Código)"),
)
```

```{=typst}
#esteira(
  no("4.1 RED", sub: "teste falha primeiro"),
  no-claro("4.2 GREEN", sub: "implementação real"),
  no("4.3 REFACTOR", sub: "linters e funclen"),
  no("4.4 STOP-GATE", sub: "auditoria do épico"),
)
```

### Etapa 4.1 — Ciclo RED (Teste Falhando Primeiro)
* **Consome:** Seção âncora do critério de aceite na SPEC (`spec_anchor`) e arquivo de teste alvo.
* **Peças:** Skill [`aidd-tdd`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-tdd/SKILL.md) + Script executor de testes locais (`go test`).
* **Processamento:** O agente implementa o teste unitário estrito reproduzindo fielmente o cenário do critério BDD.
* **Gate:** O teste **deve falhar comprovadamente** com erro de asserção funcional (`exit_code != 0`). Uma falha de compilação invalida o teste.
* **Handoff:** Grava o código do teste e o log de falha em `scratch/test_red.log`.

### Etapa 4.2 — Ciclo GREEN (Implementação Mínima de Produção)
* **Consome:** Teste falhando gerado na Etapa 4.1 *(Dependência estrita)*.
* **Peças:** Skill [`implementacao`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/implementacao/SKILL.md) + Script de compilação.
* **Processamento:** Escreve o código de produção mínimo e limpo para fazer o teste unitário passar. Proibição absoluta de stubs, mocks falsos, `TODO` ou `pass`.
* **Gate:** A suíte de teste daquele micro-ticket deve retornar código de sucesso (`exit_code == 0`).
* **Handoff:** Salva o código de produção funcional no workspace do ticket.

### Etapa 4.3 — Ciclo REFACTOR & Quality Gates
* **Consome:** Código funcional aprovado na Etapa 4.2.
* **Peças:** Scripts locais de qualidade: `go vet ./...`, `funclen` (nenhuma função > 35 linhas), `gofmt -l .`.
* **Processamento:** Limpeza de redundâncias, padronização idiomática de código e redução da complexidade ciclomática.
* **Gate:** Zero violações de linter e saída vazia no verificador de formatação.
* **Handoff:** Consolida os commits locais do micro-ticket no branch do Épico.

### Etapa 4.4 — Human-Gate por Épico
* **Consome:** Conjunto de micro-tickets finalizados do Épico ativo.
* **Peças:** Controle `human-gate`.
* **Processamento:** Pausa obrigatória ao fechar um Épico. Apresenta o sumário de testes verdes e commits para auditoria humana antes de avançar para o próximo Épico.
* **Gate:** Aprovação explícita do desenvolvedor.
* **Handoff:** Encaminha para o próximo Épico ou, ao término de todos, para o Sub-Loop 5.

---

## 5. Sub-Loop 5: Homologação E2E, Recompilação & Fechamento

O último sub-loop garante que o sistema construído opera em perfeita harmonia no mundo real, atualiza os binários executáveis e encerra a rastreabilidade da sessão.

```{=typst}
#ficha(
  ("Papel", "Homologação ponta a ponta, compilação de binários e fechamento"),
  ("Disparo", "Conclusão e aprovação de 100% dos Épicos de construção"),
  ("Entrada", "Código de produção completo no repositório"),
  ("Saída Principal", "Binário executável no PATH e REPORT.md"),
  ("Portão de Saída", "STOP-GATE FINAL (Homologação e Entrega)"),
)
```

```{=typst}
#esteira(
  no("5.1 E2E SMOKE", sub: "Playwright headless"),
  no("5.2 RECOMPILAÇÃO", sub: "go build legoloop.exe"),
  no-claro("5.3 RELATÓRIO", sub: "sessão e métricas"),
  no("5.4 STOP-GATE", sub: "entrega homologada"),
)
```

### Etapa 5.1 — Homologação E2E Multi-Viewport & Smoke Test
* **Consome:** Aplicação viva rodando na porta local do ambiente de desenvolvimento.
* **Peças:** Skill [`aidd-homologue`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-homologue/SKILL.md) + MCP `playwright`.
* **Processamento:** Navega autonomamente por todos os fluxos críticos de usuário em resoluções mobile, tablet e desktop, auditando erros de console HTTP e regressões visuais.
* **Gate:** Zero erros de console HTTP e suíte E2E 100% verde.
* **Handoff:** Grava o relatório de execução e telas de evidência em `scratch/e2e/`.

### Etapa 5.2 — Recompilação Obrigatória de Binários
* **Consome:** Código-fonte final aprovado no repositório.
* **Peças:** Script local de compilação: `go build -o bin/legoloop.exe ./cmd/legoloop`.
* **Processamento:** Compilação nativa de release do executável para garantir que as alterações no motor Go estejam ativas no PATH do sistema operacional.
* **Gate:** `exit_code == 0 && file_exists('bin/legoloop.exe')`.
* **Handoff:** Atualiza o binário oficial no diretório `bin/`.

### Etapa 5.3 — Consolidação de Relatório & Sincronização de Sessão
* **Consome:** Logs de execução, telemetria de tokens de todos os blocos e status dos gates.
* **Peças:** Skill [`aidd-sessao`](file:///C:/Users/trcnologia/.gemini/antigravity-cli/skills/aidd-sessao/SKILL.md) + Script CLI `faz-session`.
* **Processamento:** Consolida métricas de tokens de entrada e saída, custo financeiro aproximado, modelos utilizados e balanço de conformidade contra a SPEC inicial.
* **Gate:** `file_exists('docs/reports/YYYY-MM-DD_REPORT_<TEMA>.md')`.
* **Handoff:** Grava o relatório final consolidado em `docs/reports/` e salva o histórico da sessão.

### Etapa 5.4 — Human-Gate Final (Entrega do Produto)
* **Consome:** Relatório final consolidado e binário compilado.
* **Peças:** Controle `human-gate`.
* **Processamento:** Apresentação formal de encerramento da entrega para o cliente ou líder técnico.
* **Gate:** Aceite final assinado pelo usuário humano.
* **Handoff:** Conclusão bem-sucedida do Mega Pipeline.

## Rastreabilidade

* Esquema canônico de nós: `internal/dag/types.go`.
* Implementação do runner e compilação de gates: `internal/gate/gate.go`.
* Manual mestre dos sub-loops: `docs/manuals/2026-11-01_MANUAL_MEGA_PIPELINE_LEGOLOOP.md`.
