---
name: legoloop-video-rea
description: Use quando o usuário quiser FAZER ENGENHARIA REVERSA de um vídeo (YouTube, Instagram Reels, TikTok ou arquivo local), extraindo o método, ferramentas, prompts explícitos/ocultos e gerando a governança canônica e o Processo do LegoLoop, via /legoloop-video-rea. NÃO use para rodar Processos já existentes ou compilar livros avulsos.
consumes: Artifact:VideoURL
produces: Artifact:Processo
artifact: processos/<slug>/processo.dag.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🎥 legoloop-video-rea — Processo Universal de Engenharia Reversa Audiovisual

> **Autor:** LegoLoop  
> **Versão:** 1.3.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  
> **Combustível Cognitivo:** Suporte nativo a `--model <id>` e `--effort <low|medium|high>` com seleção inteligente em cascata.  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário fornece uma URL de vídeo (YouTube, Instagram Reels, TikTok) ou vídeo local e deseja extrair todo o conhecimento prático, decodificar prompts explícitos e ocultos, fazer OCR de trechos de código, curar as ferramentas contra Peças soberanas locais Tier Ouro do Cofre, conduzir a entrevista socrática de afunilamento, emitir a governança ISO fatiada e fabricar um novo Processo canônico (`legoloop-*` ou `<slug>`) com seu Livro Didático LegoPress.  
**Quando NÃO usar:** Executar Processos que já foram construídos (use `legoloop run` ou `/legoloop-video-motion`) ou compilar livros de textos arbitrários fora do Processo (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - URL do vídeo ou caminho local de arquivo multimídia.
* **Produces (Saídas Mecânicas):**
  - Target Workspace isolado (`docs/videorea/<slug>/`) com camada volátil (`scratch/`), forense (`dossier.json`), matriz original (`original_matrix.json`), curadoria (`curation.json`), governança local (`governance/`), livro didático (`book/`), Processo (`processo/`) e skill (`skill/`).
  - Projeção canônica na raiz (`processos/<slug>/processo.dag.json` e `skills/<slug>/SKILL.md`).
  - Golden Fixture determinística de 3s em MP4 1080p60 (`processo/fixture_3s.mp4`).
  - Livro Didático compilado pelo LegoPress (`docs/books/<slug>/`).

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA utilizar o protocolo STE-PT durante a Entrevista Socrática.** O STE-PT engessa o diálogo de afunilamento. Durante a entrevista socrática, o diálogo humano-agente deve ser fluido, conversacional e apresentar estritamente **uma única pergunta por turno**.
> 2. **NUNCA avançar de Etapa sem aprovação explícita no Stop-Gate humano correspondente.** Nenhuma Etapa de governança, código ou projeção nasce como fato consumado (`APPROVED` ou `COMPLETED`).
> 3. **NUNCA pular a Entrevista Socrática.** O rascunho da INTENT deve atingir confiança $\ge 95\%$ via motor oficial antes de avançar para a SPEC.
> 4. **NUNCA criar Tickets em ROUTE.json com status prematuro.** Todo Ticket nasce como `PENDING` e só progride após validação mecânica real (TDD).
> 5. **NUNCA poluir a raiz do repositório com arquivos intermediários.** Todo processamento, áudio extraído e mídias brutas residem obrigatoriamente dentro de `docs/videorea/<slug>/scratch/`.

---

## 3. Protocolo de Operação Fatiado em 7 Etapas

### Etapa 1: Forense & Matriz Original
* Extração bruta dos metadados, transcrição e OCR em `dossier.json`.
* Mapeamento da stack tecnológica do autor em `original_matrix.json` (custos mensais, licenças e timestamps de evidência).

### Etapa 2: Curadoria Soberana com Diagnóstico de Host
* Mapeamento de alternativas abertas a custo R$ 0,00/mês em `curation.json`.
* Diagnóstico de encaixe Lego no SO do operador (`installed_on_host`, `host_location`, `fit_status`, `fit_solution`, `gold_standard`, `install_command`).
* 🛑 **STOP-GATE 01 (`SG-CURADORIA`):** O operador avalia o raio-X das ferramentas e aprova o plano de soberania.

### Etapa 3: Governança Fatiada & Entrevista Socrática
1. **Entrevista Socrática (Sem STE-PT):**
   * Disparo via `legoloop fabricar entrevista iniciar "<ideia>"`.
   * Diálogo limpo e direto, **uma pergunta por vez**, até a confiança atingir a meta de $\ge 95\%$.
2. **Geração da INTENT:**
   * Materialização de `YYYY-MM-DD_INTENT_<SLUG>.md` com status `UNDER_REVIEW`.
   * 🛑 **STOP-GATE 02 (`SG-INTENT`):** O operador aprova o escopo inegociável da intenção.
3. **Geração da SPEC Técnica:**
   * Materialização de `YYYY-MM-DD_SPEC_<SLUG>.md` estruturada com base nos 13 guias autorais e contratos formais BDD.
   * 🛑 **STOP-GATE 03 (`SG-SPEC`):** O operador aprova a especificação técnica.
4. **Geração do PLAN e ROUTE Executável:**
   * Materialização de `YYYY-MM-DD_PLAN_<SLUG>.md` e `YYYY-MM-DD_ROUTE_<SLUG>.json` (espelhos 1:1 isomórficos v2.2.0 com Tickets em `PENDING` e hashes SHA-256).
   * 🛑 **STOP-GATE 04 (`SG-ROUTE`):** O operador aprova o loteamento e roteamento dos Tickets.

### Etapa 4: Materialização do Motor Executável
* Construção do DAG em `processo/processo.dag.json` (validado com `legoloop validate`).
* Síntese do Golden Fixture de 3 segundos em `processo/fixture_3s.mp4`.
* 🛑 **STOP-GATE 05 (`SG-DAG-FIXTURE`):** O operador valida a topologia do DAG e a mídia de teste.

### Etapa 5: Skill Canônica Universal
* Redação de `skill/SKILL.md` com contratos de I/O e guardrails negativos.
* Distribuição para todos os harnesses via `legoloop instalar-skills`.

### Etapa 6: Livro Didático LegoPress
* Compilação do livro em `book/` (`livro.json`, `partes/`, `.md` e `.pdf`).
* Auditoria geométrica de páginas órfãs com `legoloop legopress auditar`.

### Etapa 7: Canonização e Fechamento com Telemetria
* Execução de `legoloop fechar "<slug>"` para emissão do relatório final com balanço financeiro, orçado vs realizado e recompilação do binário no PATH.
