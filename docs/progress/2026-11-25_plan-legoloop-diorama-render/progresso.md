---
schema: progresso.v1
iniciativa: "D03"
plano: "docs/plans/2026-11-25_PLAN_LEGOLOOP_DIORAMA_RENDER.md"
atualizado_em: "2026-11-25T00:00:00Z"
---

# Progresso Vivo: D03

> Gerado automaticamente a partir do estado validado por Gate. Não edite à mão.

## Resumo

| Total | Aprovados (PASSED) | Em andamento (IN_PROGRESS) | Reprovados (FAILED) | Bloqueados (BLOCKED) | Pendentes (PENDING) |
|---:|---:|---:|---:|---:|---:|
| 15 | 0 | 0 | 0 | 0 | 15 |

## Ponto de Retomada (Em Andamento)

| Ticket | Épico | Responsável | Último registro |
|---|---|---|---|
| nenhum | - | - | - |

## Reprovados ou Bloqueados

| Ticket | Épico | Estado | Exit | Prova |
|---|---|---|---:|---|
| nenhum | - | - | - | - |

## Pendentes

| Ticket | Épico | Título | Critérios |
|---|---|---|---|
| `TICK-DR-01.1` | E01 | Criar contratos estritos de ConfigRender, Timings, Timeline, Batida, EventoSFX e RelatorioRender | AC-E01.1 |
| `TICK-DR-01.2` | E01 | Implementar pre-voo do Ciclo 3: cobertura de cenas, sprites, provedor registrado, ferramentas e render.json | AC-E01.1, AC-E01.2, AC-NF.2 |
| `TICK-DR-02.1` | E02 | Implementar Etapa 6 (voz): interface Sintetizador, registro de provedores, adaptador edge-tts e gate de creditos | AC-E02.1, AC-E02.5 |
| `TICK-DR-02.2` | E02 | Implementar Etapa 6 (Whisper): Transcritor, WER por linha (15%) e global (5%), timings.json e regravacao seletiva com parada em audio identico | AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.7 |
| `TICK-DR-02.3` | E02 | Implementar mapeamento deterministico de fonemas PT-BR para visemas | AC-E02.6 |
| `TICK-DR-03.1` | E03 | Implementar Etapa 8: leitura estrita, camadas 2.5D por profundidade, camera e troca de pose/expressao | AC-E03.1, AC-E03.2, AC-E03.3 |
| `TICK-DR-03.2` | E03 | Implementar segmentacao em batidas e gate check_timeline (12-21 cortes/min, mediana 2,5-4 s) | AC-E03.4 |
| `TICK-DR-03.3` | E03 | Implementar legendas karaoke palavra a palavra e legenda.srt | AC-E03.5 |
| `TICK-DR-04.1` | E04 | Corrigir teto de True Peak para -1,5 dBTP em videomotion.ValidateLoudnessConformity (achado F2) | AC-E04.4 |
| `TICK-DR-04.2` | E04 | Implementar Etapa 7: eventos de SFX, BGM com ducking, normalizacao em duas passadas com margem de pico e verificacao independente | AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.5 |
| `TICK-DR-05.1` | E05 | Implementar sonda e escolha de encoder NVENC/QSV/AMF com fallback libx264 | AC-E05.1 |
| `TICK-DR-05.2` | E05 | Implementar render master 16:9 em 720p/1080p/4K, versoes numeradas e gate de integridade | AC-E05.2, AC-E05.3 |
| `TICK-DR-05.3` | E05 | Implementar cortes 9:16 por re-render da timeline com area segura e legenda fora do narrador | AC-E05.4, AC-E05.5 |
| `TICK-DR-06.1` | E06 | Orquestrar Processo completo, E2E hermetico no diorama-demo, determinismo e cobertura >= 95% | AC-E06.2, AC-NF.1, AC-NF.2, AC-NF.3 |
| `TICK-DR-06.2` | E06 | Integrar CLI legoloop diorama-render, DAG v2, Skill e fechamento 5 de 5 via legoloop fechar | AC-E06.1 |

## Aprovados com Prova

| Ticket | Épico | Exit | Prova | Atualizado em |
|---|---|---:|---|---|
| nenhum | - | - | - | - |
