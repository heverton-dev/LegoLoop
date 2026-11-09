# 🗺️ PLAN: Edição Automatizada de Vídeo Soberana (video-edition)

> **Status:** UNDER_REVIEW  
> **Data:** 2026-11-08  
> **Versão:** 1.0.0  
> **Autor:** LegoLoop Video-REA Engine  
> **Iniciativa:** Mega-Pipeline Audiovisual Soberano  
> **Referência da Intenção:** [2026-11-08_INTENT_VIDEO_EDITION.md](2026-11-08_INTENT_VIDEO_EDITION.md)  
> **Referência da Especificação:** [2026-11-08_SPEC_VIDEO_EDITION.md](2026-11-08_SPEC_VIDEO_EDITION.md)  
> **Espelho Executável JSON:** [2026-11-08_ROUTE_VIDEO_EDITION.json](2026-11-08_ROUTE_VIDEO_EDITION.json)  
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
| **E01** | Preparação, Notas de Edição e Cortes | 2 | `AC-E01.1`, `AC-E01.2` | `micro` |
| **E02** | Câmera YuNet Casada com Motion e SFX com Ducking | 2 | `AC-E02.1`, `AC-E02.2` | `micro` |
| **E03** | Entrega Híbrida e Validação Mecânica | 2 | `AC-E03.1`, `AC-NF.1` | `micro` |

---

## 3. Decomposição Estruturada de Tickets (Espelho 1:1 com ROUTE.json)

### Épico 01 — Preparação, Notas de Edição e Cortes

#### TICK-01.1: Ingestão, Sincronia e Extirpação de Notas de Edição
* **Spec Ref:** `AC-E01.1`
* **Natureza:** `micro`
* **Tipo:** `unit_task`
* **Executor:** `internal/videorea`
* **Comando de Verificação:** `go test -run TestSpokenNotesExtraction -v ./internal/videorea/...`
* **Critério de Saída (Gate):** `exit_code == 0`

#### TICK-01.2: Cálculo de Cortes com Folga Pré-Fala (+40ms) e Marcadores FCPXML
* **Spec Ref:** `AC-E01.2`
* **Natureza:** `micro`
* **Tipo:** `unit_task`
* **Executor:** `internal/videorea`
* **Comando de Verificação:** `go test -run TestPreSpeechMargin -v ./internal/videorea/...`
* **Critério de Saída (Gate):** `exit_code == 0`

---

### Épico 02 — Câmera YuNet Casada com Motion e SFX com Ducking

#### TICK-02.1: Rastreamento Facial YuNet com Olhos a 35% e Casamento com Motion
* **Spec Ref:** `AC-E02.1`
* **Natureza:** `micro`
* **Tipo:** `unit_task`
* **Executor:** `internal/videorea`
* **Comando de Verificação:** `go test -run TestCameraFramingEyesHeight -v ./internal/videorea/...`
* **Critério de Saída (Gate):** `exit_code == 0`

#### TICK-02.2: Alinhamento de SFX nos Frames de Frenagem e Ducking a -18dB
* **Spec Ref:** `AC-E02.2`
* **Natureza:** `micro`
* **Tipo:** `unit_task`
* **Executor:** `internal/videorea`
* **Comando de Verificação:** `go test -run TestSFXDuckingAndAlignment -v ./internal/videorea/...`
* **Critério de Saída (Gate):** `exit_code == 0`

---

### Épico 03 — Entrega Híbrida e Validação Mecânica

#### TICK-03.1: Geração de Timeline FCPXML Aberta e Render Direto MP4
* **Spec Ref:** `AC-E03.1`
* **Natureza:** `micro`
* **Tipo:** `unit_task`
* **Executor:** `internal/videorea`
* **Comando de Verificação:** `go test -run TestHybridDeliveryOutput -v ./internal/videorea/...`
* **Critério de Saída (Gate):** `exit_code == 0`

#### TICK-03.2: Validação Topológica do Grafo DAG e Linter de Linhas
* **Spec Ref:** `AC-NF.1`
* **Natureza:** `micro`
* **Tipo:** `script`
* **Executor:** `legoloop validate` e `tools/funclen`
* **Comando de Verificação:** `go run ./cmd/legoloop validate loops/video-edition/pipeline.dag.json && go run ./tools/funclen internal/videorea`
* **Critério de Saída (Gate):** `exit_code == 0`

---

## 4. Matriz de Risco & Protocolo de Rollback

* **Risco:** Desalinhamento entre o áudio da câmera e gravação de tela externa.  
  * **Mitigação:** Correlação cruzada das formas de onda com fallback para offset manual guiado por marcador.
* **Risco:** Desvio de frames por taxa de quadros variável (VFR).  
  * **Mitigação:** Detecção no `ffprobe` e conversão prévia com FFmpeg para CFR (Constant Frame Rate) em 60fps.
* **Rollback:** Em caso de erro irrecuperável em qualquer nó, executar descarte de modificações não commitadas via `git checkout -- .` mantendo os artefatos de governança preservados.
