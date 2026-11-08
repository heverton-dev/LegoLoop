# 🎯 INTENT: Mega-Pipeline REDC (Engenharia Reversa Profunda de Canais e Clonagem Metodológica/Estrutural)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-07  
> **Autor:** LegoLoop  
> **Iniciativa:** I04 — Mega-Pipelines Audiovisuais e Automação de Canais  
> **Alvo:** `internal/redc`, `internal/cli`, `loops/legoloop-canal-redc`, `skills/legoloop-canal-redc`  
> **Sessão:** 90e39a5f-4aee-4306-8527-98d77bb3d0d2  
> **Confiança da Entrevista:** 98% (mínimo de 95% para sair de DRAFT)  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criadores de conteúdo, agências e estrategistas de mídia social precisam analisar canais virais concorrentes ou referências de mercado de forma manual, lenta e intuitiva. Perdem centenas de horas tentando decifrar o ritmo de edição, a retenção dos primeiros 3 segundos (ganchos), o design system e as metodologias ensinadas nos vídeos, sem conseguir transformar esse aprendizado em processos replicáveis, componíveis e parametrizáveis no LegoLoop com sua própria marca.
* **Quem Sofre:** Operadores do LegoLoop, criadores de conteúdo e estrategistas de canais sociais automatizados que necessitam escalar produção de vídeo mantendo o rigor técnico e a eficácia algorítmica comprovada.
* **O Porquê:** Automatizar a Engenharia Reversa Profunda de Canais (REDC) e estruturar a clonagem metodológica e estética de canais sociais cria um ativo exponencial: transforma qualquer canal de sucesso em um dossiê acionável, em dezenas de pipelines canônicos executáveis e em uma máquina geradora de vídeos autorais de alta retenção.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** Análise de canais é manual e fragmentada. Vídeos virais não são decompostos em tokens estruturados nem geram pipelines canônicos reutilizáveis; cada novo vídeo exige roteirização e design do zero.
  * **Amanhã:** O comando `legoloop canal redc <canal-slug> [--amostra N]` recebe os metadados dos vídeos do canal, indexa os 100 vídeos virais, executa o Deep REDC nos Top 15 de pico e gera de forma 100% determinística a árvore autocontida `docs/channels/<slug>/` com Dossiê DNA (tom, design tokens, motion tokens, ganchos), Catálogo de Métodos desduplicados em pipelines canônicos e Pipeline Gerador Estrutural com injeção da marca do usuário.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Arquitetura autocontida `docs/channels/<canal>/` contendo:
  * `repc/` (Dossiê REDC: tom de voz, regras de copy, design tokens, motion tokens e matriz de ganchos virais).
  * `catalogo_videos.json` (Índice estruturado dos 100 vídeos virais com métricas e clusters).
  * `pipelines/` (Pipelines canônicos do LegoLoop derivados dos métodos identificados nos vídeos).
  * `gerador/` (Pipeline gerador estrutural com injeção parametrizável de tema e marca do usuário).
  * `outputs/` (Diretório isolado para roteiros, cartões e mídias renderizadas).
* [ ] Objetivo 2: Motor Go determinístico no pacote `internal/redc` e comandos CLI `legoloop canal redc` e `legoloop canal gerar`.
* [ ] Objetivo 3: Grafo DAG acíclico canônico `loops/legoloop-canal-redc/pipeline.dag.json` com validação estrita via `legoloop validate`.
* [ ] Objetivo 4: Skill canônica universal `skills/legoloop-canal-redc/SKILL.md` distribuída para todos os harnesses via `legoloop instalar-skills`.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não criaremos dependências externas novas além das já autorizadas no projeto (stdlib Go pura).
* ⛔ Scrapers de rede ao vivo que violem termos de serviço ou dependam de conexões instáveis: a ingestão aceita pastas locais pré-coletadas com metadados e transcrições em JSON/TXT/CSV estruturado.

---

## 3. Mapeamento Preliminar de Blocos e Subloops
Identificação das capacidades que a futura esteira precisará acionar. Tipos permitidos: `skill`, `script`, `subloop`, `unit_task`, `ui_task`, `stop_gate`.

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/legoloop-spec` | Elaborar especificação formal com critérios BDD e contratos JSON de REDC. |
| `BLOCO-02` | `unit_task` | `internal/redc` (Go) | Implementação dos parsers, clusterizador de métodos e sintetizador de tokens. |
| `BLOCO-03` | `unit_task` | `internal/cli` (Go) | Comandos CLI `legoloop canal redc` e `legoloop canal gerar`. |
| `BLOCO-04` | `script` | `go test -count=1 ./internal/redc/...` | Regressão mecânica, cobertura e integridade hermética. |
| `BLOCO-05` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação fresca obrigatória (AGENTS.md, regra 5). |
| `BLOCO-06` | `skill` | `skills/legoloop-fechar` | Canonização formal, manual operacional, livro didático e distribuição. |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via CLI `legoloop canal redc <canal>` ou Chat `/legoloop-canal-redc`).  
* **Nível de Risco & Blast Radius:** `Baixo (workspace local)`: escrita estritamente confinada a `docs/channels/<slug>/` e novos pacotes isolados.  
* **Política de Privacidade & Orçamento Cognitivo:** `Nuvem autorizada (9Router / code-fast)` | Teto de tokens por bloco: 2k-5k tokens. Zero segredos em disco ou tela.  
* **Contrato de Evidência Factual (Definition of Done Visual/Real):**  
  * Execução CLI `legoloop canal redc "testdata/fixtures/canal_sample" --json` com exit code 0 gerando a árvore completa em `docs/channels/canal_sample/`.  
  * Execução `legoloop validate loops/legoloop-canal-redc/pipeline.dag.json` com aprovação total de 0 ciclos e contratos de I/O válidos.  
  * Suíte completa verde: `go test -count=1 ./...`, `go vet ./...`, `funclen <= 35`.  
* **Premissas Assumidas (a confirmar no Stop-Gate 01):**  
  * Amostragem inteligente (Top 15 vídeos de pico para extração de estilo e clusterização temático-metodológica dos 100 vídeos) é suficiente e superior à análise cega de 100 vídeos.  
  * Os dados brutos entram via diretório de arquivos locais (`transcricoes/`, `metricas.json`).  
* **Riscos e Perguntas em Aberto:** Nenhum risco impeditivo; ambiente local isolado e hermético com testes em `testdata/fixtures`.

---

## 5. Critérios de Sucesso de Alto Nível
Cada critério precisa ter uma forma mecânica de prova.

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs e zero mocks provisórios. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura de testes mantida ou expandida (TDD Red-Green no pacote `internal/redc`). | Teste novo falhando antes do código | `go test -v ./internal/redc/...` |
| 3 | Árvore `docs/channels/<canal>/` gerada conforme contrato ISO | Existência mecânica dos arquivos e JSONs válidos | `legoloop canal redc` com verificação de filesystem |
| 4 | Grafo DAG acíclico aprovado | Validação do grafo pelo motor LegoLoop | `legoloop validate loops/legoloop-canal-redc/pipeline.dag.json` |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-07_SPEC_CANAL_REDC.md`](../specs/2026-11-07_SPEC_CANAL_REDC.md)
* **Gate de Transição:** Este Intent só é considerado aprovado após a validação humana explícita (*Human-in-the-Loop Gate*) no chat e com confiança da entrevista ≥ 95%, antes de qualquer linha de especificação ser redigida.

---

🛑 **STOP-GATE 01 DE AUDITORIA HUMANA:**  
*A intenção consolidada com 98% de certeza factual atende à sua expectativa?*  
*Somente após a aprovação explícita do usuário, deve-se gerar a SPEC (`TEMPLATE_SPEC.md`).*