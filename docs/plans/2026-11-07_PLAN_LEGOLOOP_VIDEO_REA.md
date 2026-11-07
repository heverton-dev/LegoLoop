# 🗺️ PLAN: Pipeline Universal de Engenharia Reversa Audiovisual (legoloop-video-rea)

> **Status:** Aprovado pelo Usuário  
> **Data:** 2026-11-07  
> **Versão:** 1.0.0  
> **Autor:** Heverton & Ecossistema LegoLoop  
> **Iniciativa:** LegoLoop Autonomous Evolution  
> **Referência da Intenção:** [docs/intents/2026-11-07_INTENT_LEGOLOOP_VIDEO_REA.md](../intents/2026-11-07_INTENT_LEGOLOOP_VIDEO_REA.md)  
> **Referência da Especificação:** [docs/specs/2026-11-07_SPEC_LEGOLOOP_VIDEO_REA.md](../specs/2026-11-07_SPEC_LEGOLOOP_VIDEO_REA.md)  
> **Espelho Executável JSON:** [docs/routes/2026-11-07_ROUTE_LEGOLOOP_VIDEO_REA.json](../routes/2026-11-07_ROUTE_LEGOLOOP_VIDEO_REA.json)  
> **Template:** TEMPLATE_PLAN.md v2.2.0 (espelha TEMPLATE_ROUTE.json v2.2.0)  

---

## 1. Estratégia de Execução & Governança de Loteamento

* **Engine:** `orchestrix_hermetic` (execução isolada por bloco, sem context bloat)
* **Nível de Supervisão:** `micro`
* **Permissão de Subagentes:** `false` (subprocessos estéreis)
* **TDD Estrito:** `true` (Red ➔ Green ➔ Refactor obrigatório antes de qualquer código de produção)
* **Limite de Workers Paralelos:** `2`
* **Timeout Padrão por Bloco:** `180s`
* **Política de Rollback:** global `abort_pipeline`; por nó `stash_changes`

---

## 2. Mapa Geral de Épicos & Rastreabilidade de Critérios

| Épico ID | Nome do Épico | Total Tickets | Critérios de Aceite Cobertos | Natureza Dominante |
|---|---|---|---|---|
| **E01** | Ingestão Stealth & Crawler Multi-Plataforma | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Dissecação Forense, Prompts Ocultos e OCR | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Curadoria Matriz LegoLoop e Economia de Soberania | 2 | `AC-E03.1`, `AC-E03.2` | `micro` |
| **E04** | Stop-Gate Humano e Fábrica de Governança | 2 | `AC-E04.1`, `AC-E04.2` | `micro` |
| **E05** | Orquestrador DAG, CLI, Livro LegoPress e Skill | 2 | `AC-E05.1`, `AC-E05.2` | `macro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Ingestão Stealth & Crawler Multi-Plataforma

#### `TICK-01.1` — Detector de Plataforma e Ingestão
* **Nome:** `videorea_crawler_detect`
* **Natureza:** `micro` | **Complexidade Cognitiva:** `simples`
* **Alvo Principal:** `internal/videorea/crawler.go`
* **Depende De:** `[]`
* **Spec Ref:** `AC-E01.1`
* **Contrato de Teste (Highest Seam):** `TestDetectPlatform`
* **Critério de Sucesso:** Detecta corretamente YouTube, Instagram e TikTok.

#### `TICK-01.2` — Sanitizador de URL e Tratamento de Falha
* **Nome:** `videorea_url_sanitizer`
* **Natureza:** `nano` | **Complexidade Cognitiva:** `simples`
* **Alvo Principal:** `internal/videorea/crawler.go`
* **Depende De:** `["TICK-01.1"]`
* **Spec Ref:** `AC-E01.2`
* **Contrato de Teste (Highest Seam):** `TestInvalidURLHandling`
* **Critério de Sucesso:** Retorna `ErrInvalidURL` para URLs malformadas.

---

### Épico 02 — Dissecação Forense, Prompts Ocultos e OCR

#### `TICK-02.1` — Dissecador Forense de Prompts
* **Nome:** `videorea_forensics_prompts`
* **Natureza:** `micro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `internal/videorea/forensics.go`
* **Depende De:** `["TICK-01.2"]`
* **Spec Ref:** `AC-E02.1`
* **Contrato de Teste (Highest Seam):** `TestExtractForensicsPrompts`
* **Critério de Sucesso:** Separa prompts explícitos e reconstrói prompts ocultos.

#### `TICK-02.2` — Inspetor de Assets e OCR de Telas Rápidas
* **Nome:** `videorea_forensics_ocr`
* **Natureza:** `micro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `internal/videorea/forensics.go`
* **Depende De:** `["TICK-02.1"]`
* **Spec Ref:** `AC-E02.2`
* **Contrato de Teste (Highest Seam):** `TestInspectAssetsAndOCR`
* **Critério de Sucesso:** Identifica ferramentas ocultas e texto de telas de código.

---

### Épico 03 — Curadoria Matriz LegoLoop e Economia de Soberania

#### `TICK-03.1` — Curador Semântico com Peças Tier Ouro
* **Nome:** `videorea_matrix_curator`
* **Natureza:** `micro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `internal/videorea/curator.go`
* **Depende De:** `["TICK-02.2"]`
* **Spec Ref:** `AC-E03.1`
* **Contrato de Teste (Highest Seam):** `TestCurateAgainstMatrix`
* **Critério de Sucesso:** Substitui ferramentas pagas por equivalentes locais Tier Ouro a R$ 0,00.

#### `TICK-03.2` — Gerador de Nomenclatura Canônica `legoloop-*`
* **Nome:** `videorea_canonical_naming`
* **Natureza:** `nano` | **Complexidade Cognitiva:** `simples`
* **Alvo Principal:** `internal/videorea/curator.go`
* **Depende De:** `["TICK-03.1"]`
* **Spec Ref:** `AC-E03.2`
* **Contrato de Teste (Highest Seam):** `TestCanonicalNaming`
* **Critério de Sucesso:** Garante que todos os slugs sugeridos possuem o prefixo `legoloop-*`.

---

### Épico 04 — Stop-Gate Humano e Fábrica de Governança

#### `TICK-04.1` — Stop-Gate de Curadoria Interativa
* **Nome:** `videorea_stopgate_curation`
* **Natureza:** `nano` | **Complexidade Cognitiva:** `simples`
* **Alvo Principal:** `internal/videorea/gate.go`
* **Depende De:** `["TICK-03.2"]`
* **Spec Ref:** `AC-E04.1`
* **Contrato de Teste (Highest Seam):** `TestStopGateCuration`
* **Critério de Sucesso:** Valida aprovação ou recusa do operador.

#### `TICK-04.2` — Fábrica de Tríade ISO e Golden Fixture (3s)
* **Nome:** `videorea_triad_and_fixture`
* **Natureza:** `micro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `internal/videorea/factory_bridge.go`
* **Depende De:** `["TICK-04.1"]`
* **Spec Ref:** `AC-E04.2`
* **Contrato de Teste (Highest Seam):** `TestGenerateGovernanceAndFixture`
* **Critério de Sucesso:** Cria INTENT, SPEC, PLAN, ROUTE e clipe sintético de 3s.

---

### Épico 05 — Orquestrador DAG, CLI, Livro LegoPress e Skill

#### `TICK-05.1` — Grafo DAG e Subcomando CLI
* **Nome:** `videorea_dag_and_cli`
* **Natureza:** `macro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `loops/video-rea/pipeline.dag.json` e `internal/cli/video_rea.go`
* **Depende De:** `["TICK-04.2"]`
* **Spec Ref:** `AC-E05.1`
* **Contrato de Teste (Highest Seam):** `TestVideoReaCLI`
* **Critério de Sucesso:** DAG validado pelo LegoLoop e CLI operacional.

#### `TICK-05.2` — Livro Didático LegoPress e Skill Canônica
* **Nome:** `videorea_book_and_skill`
* **Natureza:** `macro` | **Complexidade Cognitiva:** `media`
* **Alvo Principal:** `docs/books/video-rea/` e `skills/legoloop-video-rea/SKILL.md`
* **Depende De:** `["TICK-05.1"]`
* **Spec Ref:** `AC-E05.2`
* **Contrato de Teste (Highest Seam):** `TestVideoReaBookCompilation`
* **Critério de Sucesso:** PDF gerado em `docs/books/video-rea/` com auditoria 100% válida.
