# Relatório de Análise: All Code Is Slop (Except Mine) — Lauren Tan

**Data:** 2026-11-14
**Fonte:** [YouTube — AI Engineer](https://www.youtube.com/watch?v=O_zcc6Xl12g)
**Palestrante:** Lauren Tan (`poteto`) — SpaceXSI (Grok Bot), ex-Cursor
**Duração:** 17:31

---

## 1. Resumo Executivo

Lauren Tan apresenta como o time do Grok Bot na SpaceXSI atingiu **11.000 pull requests em um mês** com **taxa de revert de ~0,5%** usando agentes de código em escala. A tese central é que a qualidade da saída dos agentes depende diretamente da qualidade determinística do ambiente (código, arquitetura, lint, tipos e verificação), reduzindo a dependência de instruções em prosa (`AGENTS.md`, `rules`, `skills`) apenas ao que exige julgamento.

---

## 2. Achados Extraídos do Vídeo

### 2.1 Conceitos e Processos


| Categoria    | Achado                                         | Descrição~                                                                                                                                                                             |
| :--- | :--- | :--- |
| **Conceito** | **Trust Pyramid (Pirâmide Invertida)**         | Prioriza investimentos determinísticos na base (arquitetura, APIs, lint, tipos, testes), depois verificação (runtime e formal) e por último prosa (`rules`, `skills`, `AGENTS.md`).    |
| **Conceito** | **Codebase as a Garden**                       | O repositório é o maior input dos agentes. Workarounds e débito técnico viram "ervas daninhas" copiadas automaticamente pelos próximos agentes.                                        |
| **Processo** | **Banimento de APIs com Nuance**               | APIs propensas a uso incorreto por agentes (ex.: `useEffect` no React) são banidas no código de produto e encapsuladas no framework interno.                                           |
| **Processo** | **Contexto por Estrutura de Diretórios**       | Separação física explícita (ex.: `electron-renderer` vs `electron-main`) com bloqueio mecânico de imports cruzados via analisador de grafo de dependências.                            |
| **Processo** | **Substituição de Comentários por Lint Rules** | Banimento de ~99% dos comentários de código. Agentes usam comentários como *escape hatch* para justificar gambiarras. Correções viram regras de lint com mensagens de erro acionáveis. |
| **Processo** | **Feature Map**                                | Mapa estruturado das funcionalidades da aplicação sob a ótica do usuário (equivalente a um sitemap de features) para eliminar ambiguidade de vocabulário com o agente.                 |
| **Processo** | **Verificação como Artefato de Revisão**       | Revisão humana migra de leitura de código para inspeção de artefatos de verificação (vídeos, screenshots, diagramas de estado). Cerca de **70% dos tokens** são gastos em verificação. |
| **Processo** | **Fuzzing Contínuo Pós-Release**               | A cada release, o orquestrador (`Grok Bot`) dispara agentes para testar o app como usuário real, abrir PRs de correção e fazer auto-merge.                                             |


### 2.2 Ferramentas, Peças e Repositórios Citados


| Item                               | Tipo                   | Uso no Vídeo                                                                                                        | Referência                                                                                       |
| :--- | :--- | :--- | :--- |
| **pstack (Potato Stack)**          | Repositório / Plugin   | Coleção de skills e peças da Lauren Tan para Cursor, incluindo a skill `statechart`.                                | [github.com/cursor/plugins/tree/main/pstack](https://github.com/cursor/plugins/tree/main/pstack) |
| **Dune**                           | Framework interno      | Framework cliente para Electron ("Next.js para Electron") com constraints e lint rules embutidas.                   | Interno SpaceXSI (Grok Bot)                                                                      |
| **Grok Bot**                       | Orquestrador           | Outer loop que dispara projetos e agentes em nuvem (Cursor Cloud Agents) e orquestra PRs.                           | Interno SpaceXSI                                                                                 |
| **Cursor / Cloud Agents**          | Harness / Agentes      | Execução paralela de agentes de código e fuzzing.                                                                   | Cursor                                                                                           |
| **Chrome DevTools Protocol (CDP)** | Protocolo / Ferramenta | Verificação em runtime: permite ao agente rodar, inspecionar e depurar a UI.                                        | Padrão Chromium                                                                                  |
| **TLA+**                           | Verificação Formal     | Model checking de fluxos críticos (ex.: notificações entre múltiplos dispositivos).                                 | Ferramenta aberta                                                                                |
| **Lean**                           | Prova Formal           | Prova matemática de algoritmos críticos (ex.: biblioteca de virtualização sem flickering).                          | Ferramenta aberta                                                                                |
| **Bombadil**                       | Property-Based Testing | Testes baseados em propriedades para interfaces de usuário (UI).                                                    | Biblioteca                                                                                       |
| **Custom CLIs em Skills**          | Peça / Padrão          | Encapsular drivers e automações em CLIs próprias dentro das skills para economizar tokens e garantir testabilidade. | Padrão `pstack`                                                                                  |


---

## 3. Análise de Integração ao LegoLoop

### 3.1 DEVEM ser integrados (Alta Prioridade)

1. **Substituição de Prosa por Gates Determinísticos (Trust Pyramid)**
  - **Por quê:** Alinha-se 100% às **5 Prisões Mecânicas** do LegoLoop (`Path Jail`, `Schema Validation`, `Exit-Code Gate`, `Sterile Context`, `Test Mutation`) e reduz consumo de tokens em `AGENTS.md` e `SKILL.md`.
  - **Como:** Sempre que um agente errar repetidamente em um ticket, converter a correção em regra no `tools/taxlinter`, `tools/funclen` ou novo linter Go em `tools/`, em vez de adicionar mais texto em `AGENTS.md`.
2. **Contexto por Diretórios + Bloqueio de Imports Cruzados**
  - **Por quê:** O LegoLoop já possui `codebase-memory-mcp` e pacotes `internal/` bem definidos, mas pode reforçar fronteiras entre pacotes (ex.: `internal/studio` vs `internal/ste` vs `internal/provider`) de forma mecânica.
  - **Como:** Criar regra de verificação estática em Go (`tools/archlinter` ou teste arquitetural) validando imports permitidos entre pacotes `internal/`.
3. **Feature Map para o Studio**
  - **Por quê:** Reduz ambiguidade quando agentes editam a UI do Studio (`internal/studio`), evitando que inventem nomes ou quebrem fluxos existentes.
  - **Como:** Gerar um `docs/FEATURE_MAP_STUDIO.md` (ou JSON) listando rotas, painéis e ações do Studio sob a ótica do usuário, referenciado pelas skills de UI.
4. **Custom CLIs em Skills (em vez de instruções longas)**
  - **Por quê:** Já é parcialmente praticado no LegoLoop (`legoloop fechar`, `legoloop fabricar`, `legoloop stept`). Reforçar como regra obrigatória para novas skills reduz tokens e torna cada skill testável.
  - **Como:** Documentar no protocolo de autoria de skills que toda automação recorrente deve virar subcomando Go em `cmd/legoloop` ou script testável, nunca passos manuais longos em Markdown.

### 3.2 PODEM ser integrados (Média Prioridade / Sob Demanda)

1. **Skill** `statechart` **(do** `pstack`**)**
  - **Por quê:** Útil para revisar máquinas de estado complexas (ex.: ciclo de vida de Processos, Stop-Gates e execução de DAGs) visualmente.
  - **Como:** Inspecionar [github.com/cursor/plugins/tree/main/pstack](https://github.com/cursor/plugins/tree/main/pstack) e adaptar o conceito para gerar diagramas Mermaid `stateDiagram-v2` (respeitando a Lei #9: sem HTML e com aspas).
2. **Property-Based Testing para UI (**`Bombadil` **ou equivalente JS/Playwright)**
  - **Por quê:** O Studio possui interface web; testes baseados em propriedades podem detectar regressões de estado na UI sem escrever dezenas de casos manuais.
  - **Como:** Avaliar adoção pontual na suíte E2E do Studio quando houver refatorações grandes de interface.
3. **Verificação Formal (**`TLA+`**) para Concorrência de DAGs**
  - **Por quê:** O motor de execução paralela de DAGs (`legoloop-swarm` / worktrees) tem estados concorrentes críticos (deadlocks, join barriers).
  - **Como:** Opcionalmente modelar em `TLA+` as invariantes do join barrier e do disjuntor Andon Cord caso surjam bugs de concorrência difíceis de reproduzir em testes Go.

### 3.3 NÃO PODEM ser integrados diretamente

1. **Framework** `Dune` **e Orquestrador** `Grok Bot`
  - **Por quê:** São ferramentas internas e proprietárias da SpaceXSI, acopladas a Electron e à infraestrutura interna deles. O LegoLoop é um binário Go (`stdlib` apenas, sem dependências externas em `go.mod`).
  - **Alternativa:** Adotar apenas os *princípios* (constraints no código, lint rules e outer loop automatizado via `legoloop` + `pipeline.dag.json`).
2. **Auto-Merge Cego de PRs sem Stop-Gate Humano**
  - **Por quê:** Viola diretamente a **Lei #8** e os **Stop-Gates (**`SG-*`**)** do LegoLoop, que exigem aprovação humana explícita para ações irreversíveis, commits e pushes.
  - **Alternativa:** Manter o fluxo de verificação automatizada (gates mecânicos) gerando evidências para aprovação rápida humana no Stop-Gate.
3. **Banimento Total de Comentários em Go**
  - **Por quê:** A diretriz do repositório exige docstrings concisas em funções e classes novas e comentários explicando o *porquê* não óbvio.
  - **Alternativa:** Banir apenas comentários que justificam *workarounds* ou gambiarras, exigindo correção na causa raiz.
