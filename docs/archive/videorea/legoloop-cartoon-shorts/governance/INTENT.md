# 🎯 INTENT: Pipeline Canônico LegoLoop Cartoon Shorts (legoloop-cartoon-shorts)

> **Status:** APPROVED  
> **Data:** 2026-11-09  
> **Autor:** LegoLoop (Video-REA Engine)  
> **Iniciativa:** Soberania Audiovisual Local Tier Ouro  
> **Alvo:** `internal/videorea/`, `internal/cli/`, `loops/legoloop-cartoon-shorts/`, `skills/legoloop-cartoon-shorts/`, `web/studio/`  
> **Sessão:** intv_1794238593016059900  
> **Confiança da Entrevista:** 100%  
> **Template:** TEMPLATE_INTENT.md v2.0.0  

---

## 1. Contexto e Motivação
* **O Problema Real:** Criadores de conteúdo gastam entre R$ 180,00 e R$ 250,00/mês com Claude Pro e ElevenLabs ou esgotam cotas de tokens tentando gerar shorts animados. O autor do vídeo (@ideias_artificiais) adaptou o projeto open-source `mromidev/claude-animated-shorts` e forneceu a skill `cartoon-audio-short` e o kit `cartoon-short-kit.zip`, mas manteve dependência do Claude e ElevenLabs pagos.
* **Quem Sofre:** Criadores de shorts/reels, operadores do LegoLoop e agentes que precisam produzir vídeos verticais animados (1080x1920 a 30fps) sem custos recorrentes de API e com execução determinística local.
* **O Porquê:** Adoção soberana completa (Tier Ouro): incorporar o motor Three.js r160 (`cartoon-short-kit`), substituir o ElevenLabs pelo **Edge-TTS Neural Local (`pt-BR-AntonioNeural`)** interpretando tags de emoção entre colchetes (`[reflexivo]`, `[animado]`), substituir o Claude pelo **9Router Local**, e renderizar frame a frame via Playwright Chromium (SwiftShader CPU WebGL) com mixagem determinística de áudio a **50% (-6 dB Gated RMS)** via FFmpeg a custo R$ 0,00/mês.
* **Estado Atual vs. Estado Desejado:**
  * **Hoje:** O ecossistema não possui o fluxo canônico de Cartoon Shorts empacotado para a Tríade Completa.
  * **Amanhã:** Execução 100% soberana do **LegoLoop Cartoon Shorts** via CLI (`legoloop run legoloop-cartoon-shorts`), via Studio Web (com player de preview) e via chat (`/legoloop-cartoon-shorts`), gerando o MP4 final 1080x1920 a 30fps e artefatos intermediários isolados no workspace.

---

## 2. Escopo Inegociável

### O que ESTÁ no Escopo (Objetivos Primários):
* [ ] Objetivo 1: Integrar o kit oficial `cartoon-short-kit` (Three.js r160, `paint.js`, `motion.js`, `roman.js`, `timing.py`, `render.mjs`) no motor `internal/videorea/`.
* [ ] Objetivo 2: Implementar a geração de roteiro com 10 ideias de curiosidade (gancho paradoxal + twist + checagem verdade/lenda) no ritmo de 3.7 palavras/s via 9Router.
* [ ] Objetivo 3: Implementar a síntese de voz soberana via Edge-TTS neural (`pt-BR-AntonioNeural`) com parser de emoções e mixagem com trava de -6 dB Gated RMS via FFmpeg.
* [ ] Objetivo 4: Disponibilizar a Tríade Completa de interfaces: CLI (`legoloop run legoloop-cartoon-shorts`), painel do Studio Web com preview player e comando de chat (`/legoloop-cartoon-shorts`).
* [ ] Objetivo 5: Garantir conformidade com as 14 leis de AGENTS.md, zero stubs, cobertura TDD estrita e recompilação do binário no PATH.

### O que NÃO ESTÁ no Escopo (Não-Objetivos / Out-of-Scope):
* ⛔ Não faremos refatorações oportunistas fora dos módulos declarados em `Alvo`.
* ⛔ Não alteraremos contratos públicos de API/JSON sem spec formal prévia.
* ⛔ Não utilizaremos APIs comerciais pagas (ElevenLabs, Claude Pro pago, Runway).
* ⛔ Não utilizaremos subagentes de harness (`invoke_subagent`).
* ⛔ Não gravaremos arquivos fora de `docs/videorea/legoloop-cartoon-shorts/scratch/`.

---

## 3. Mapeamento Preliminar de Blocos e Subloops

| Identificador | Tipo da Peça | Especialidade / Skill / Subloop | Papel na Esteira |
| :--- | :--- | :--- | :--- |
| `BLOCO-01` | `skill` | `skills/legoloop-cartoon-shorts` | Intake, 10 ideias de curiosidade e roteiro com fórmula de 3.7 palavras/s |
| `BLOCO-02` | `unit_task` | `internal/videorea/tts_edge.go` | Síntese de voz neural local (`pt-BR-AntonioNeural`) com tags de emoção |
| `BLOCO-03` | `unit_task` | `internal/videorea/timing_aligner.go` | Alinhamento determinístico fonético de tomadas via `timing.py` |
| `BLOCO-04` | `unit_task` | `internal/videorea/threejs_renderer.go` | Renderização 1080x1920 30fps Three.js r160 via Playwright headless |
| `BLOCO-05` | `unit_task` | `internal/videorea/ffmpeg_mux.go` | Mixagem Gated RMS a -6 dB, limitador -1 dBFS e encode H.264 CRF 17 |
| `BLOCO-06` | `script` | `go test ./...` | Suíte de testes unitários herméticos com zero stubs |
| `BLOCO-07` | `stop_gate` | `go build -o bin/legoloop.exe ./cmd/legoloop` | Recompilação obrigatória e instalação no PATH |

---

## 4. Governança Operacional e Restrições Inegociáveis
* **Gatilho Operacional de Disparo (Trigger):** `manual` (via CLI `legoloop run legoloop-cartoon-shorts`, Studio ou chat `/legoloop-cartoon-shorts`).
* **Nível de Risco & Blast Radius:** Baixo (workspace local isolado em `docs/videorea/legoloop-cartoon-shorts/`).
* **Política de Privacidade & Orçamento Cognitivo:** Custo R$ 0,00 de inferência e execução 100% local.
* **Contrato de Evidência Factual:** Produção de `loop/fixture_3s.mp4` e vídeo completo no disco com player de preview funcional no Studio.

---

## 5. Critérios de Sucesso de Alto Nível

| # | Critério | Como se prova | Comando / Evidência |
| :--- | :--- | :--- | :--- |
| 1 | Todas as invariantes da esteira passam com zero stubs. | Suíte verde e busca por stubs vazia | `go test ./...` com exit code 0 |
| 2 | Cobertura TDD mantida com testes reais. | Teste unitário de renderização falha antes da implementação | `go test -count=1 ./internal/videorea/...` |
| 3 | Tríade Completa operacional (CLI, Studio e Chat). | Execução ponta a ponta gera MP4 válido | `bin\legoloop.exe run legoloop-cartoon-shorts` |
| 4 | Custo operacional mensal nulo. | Telemetria de custo auditada | `R$ 0,00/mes` comprovado |

---

## 6. Handoff de Governança
* **Próximo Documento Mandatório:** [`docs/specs/2026-11-09_SPEC_LEGOLOOP_CARTOON_SHORTS.md`](../specs/2026-11-09_SPEC_LEGOLOOP_CARTOON_SHORTS.md)
* **Gate de Transição:** Intent aprovado por auditoria humana e confiança de 100%.
