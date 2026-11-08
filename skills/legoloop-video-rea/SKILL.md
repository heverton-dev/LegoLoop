---
name: legoloop-video-rea
description: Use quando o usuário quiser FAZER ENGENHARIA REVERSA de um vídeo (YouTube, Instagram Reels, TikTok ou arquivo local), extraindo o método, ferramentas, prompts explícitos/ocultos e gerando a governança canônica e o pipeline do LegoLoop, via /legoloop-video-rea. NÃO use para rodar pipelines já existentes ou compilar livros avulsos.
consumes: Artifact:VideoURL
produces: Artifact:Pipeline
artifact: loops/<slug>/pipeline.dag.json
gate: exit_code == 0
suggested_model: sonnet
---

# 🎥 legoloop-video-rea — Pipeline Universal de Engenharia Reversa Audiovisual

> **Autor:** LegoLoop  
> **Versão:** 1.2.0  
> **Template:** TEMPLATE_PIECE_SKILL.md v2.0.0  

---

## 1. Contexto & Contratos de I/O

**Quando usar:** O usuário fornece uma URL de vídeo (YouTube, Instagram Reels, TikTok) ou vídeo local e deseja extrair todo o conhecimento prático, decodificar prompts explícitos e ocultos, fazer OCR de trechos de código, curar as ferramentas contra peças soberanas locais Tier Ouro do cofre, conduzir a entrevista socrática de afunilamento, emitir a governança ISO fatiada e fabricar um novo pipeline canônico (`legoloop-*` ou `<slug>`) com seu Livro Didático LegoPress.  
**Quando NÃO usar:** Executar pipelines que já foram construídos (use `legoloop run` ou `/legoloop-video-motion`) ou compilar livros de textos arbitrários fora da esteira (use `/legoloop-legopress`).

* **Consumes (Entradas Obrigatórias):**
  - URL do vídeo ou caminho local de arquivo multimídia.
* **Produces (Saídas Mecânicas):**
  - Target Workspace isolado (`docs/videorea/<slug>/`) com camada volátil (`scratch/`), forense (`dossier.json`), matriz original (`original_matrix.json`), curadoria (`curation.json`), governança local (`governance/`), livro didático (`book/`), pipeline (`loop/`) e skill (`skill/`).
  - Projeção canônica na raiz (`loops/<slug>/pipeline.dag.json` e `skills/<slug>/SKILL.md`).
  - Golden Fixture determinística de 3s em MP4 1080p60 (`loop/fixture_3s.mp4`).
  - Livro Didático compilado pelo LegoPress (`docs/books/<slug>/`).

---

## 2. Invariantes & Regras Negativas (Negative Guardrails)

> [!CAUTION]
> **REGRAS NEGATIVAS INVIOLÁVEIS:**
> 1. **NUNCA utilizar o protocolo STE-PT durante a Entrevista Socrática.** O STE-PT engessa o diálogo de afunilamento. Durante a entrevista socrática, o diálogo humano-agente deve ser fluido, conversacional e apresentar estritamente **uma única pergunta por turno**.
> 2. **NUNCA avançar de fase sem aprovação explícita no Stop-Gate humano correspondente.** Nenhuma etapa de governança, código ou projeção nasce como fato consumado (`APPROVED` ou `COMPLETED`).
> 3. **NUNCA pular a Entrevista Socrática.** O rascunho da INTENT deve atingir confiança $\ge 95\%$ via motor oficial antes de avançar para a SPEC.
> 4. **NUNCA criar tickets em ROUTE.json com status prematuro.** Todo ticket nasce como `PENDING` e só progride após validação mecânica real (TDD).
> 5. **NUNCA poluir a raiz do repositório com arquivos intermediários.** Todo processamento, áudio extraído e mídias brutas residem obrigatoriamente dentro de `docs/videorea/<slug>/scratch/`.

---

## 3. Protocolo de Operação Fatiado em 7 Fases

### Fase 1: Forense & Matriz Original
* Extração bruta dos metadados, transcrição e OCR em `dossier.json`.
* Mapeamento da stack tecnológica do autor em `original_matrix.json` (custos mensais, licenças e timestamps de evidência).

### Fase 2: Curadoria Soberana com Diagnóstico de Host
* Mapeamento de alternativas abertas a custo R$ 0,00/mês em `curation.json`.
* Diagnóstico de encaixe Lego no SO do operador (`installed_on_host`, `host_location`, `fit_status`, `fit_solution`, `gold_standard`, `install_command`).
* 🛑 **STOP-GATE 01 (`SG-CURADORIA`):** O operador avalia o raio-X das ferramentas e aprova o plano de soberania.

### Fase 3: Governança Fatiada & Entrevista Socrática
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
   * Materialização de `YYYY-MM-DD_PLAN_<SLUG>.md` e `YYYY-MM-DD_ROUTE_<SLUG>.json` (espelhos 1:1 isomórficos v2.2.0 com tickets em `PENDING` e hashes SHA-256).
   * 🛑 **STOP-GATE 04 (`SG-ROUTE`):** O operador aprova o loteamento e roteamento dos tickets.

### Fase 4: Materialização do Motor Executável
* Construção do grafo acíclico em `loop/pipeline.dag.json` (validado com `legoloop validate`).
* Síntese do Golden Fixture de 3 segundos em `loop/fixture_3s.mp4`.
* 🛑 **STOP-GATE 05 (`SG-DAG-FIXTURE`):** O operador valida a topologia do DAG e a mídia de teste.

### Fase 5: Skill Canônica Universal
* Redação de `skill/SKILL.md` com contratos de I/O e guardrails negativos.
* Distribuição para todos os harnesses via `legoloop instalar-skills`.

### Fase 6: Livro Didático LegoPress
* Compilação do livro em `book/` (`livro.json`, `partes/`, `.md` e `.pdf`).
* Auditoria geométrica de páginas órfãs com `legoloop legopress auditar`.

### Fase 7: Canonização e Fechamento com Telemetria
* Execução de `legoloop fechar "<slug>"` para emissão do relatório final com balanço financeiro, orçado vs realizado e recompilação do binário no PATH.
