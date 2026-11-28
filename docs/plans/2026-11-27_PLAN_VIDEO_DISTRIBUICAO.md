# PLAN — legoloop-video-distribuicao

- **Data:** 2026-11-27
- **SPEC:** `docs/specs/2026-11-27_SPEC_VIDEO_DISTRIBUICAO.md`
- **Execucao:** direta (fora do swarm), TDD vermelho e verde, um ticket por vez.

| Ticket | Titulo | Criterios | Arquivos |
| :--- | :--- | :--- | :--- |
| TICK-DIST-01 | Acervo, originalidade e autorizacao | AC-E01.1, AC-E01.2, AC-E01.3 | `internal/video/distribuicao/acervo.go` |
| TICK-DIST-02 | Agenda quente/perene, estoque e reels | AC-E02.1, AC-E02.2, AC-E02.3, AC-E02.4, AC-E02.5 | `internal/video/distribuicao/agenda.go` |
| TICK-DIST-03 | Revisao: parecer, folhas e pacote do revisor | AC-E03.1, AC-E03.2, AC-E03.3, AC-E03.4, AC-E03.5 | `internal/video/distribuicao/revisao.go` |
| TICK-DIST-04 | Vitrine e capa | AC-E04.1, AC-E04.2, AC-E04.3, AC-E04.4 | `internal/video/distribuicao/vitrine.go`, `internal/video/distribuicao/capa.go`, `internal/video/distribuicao/scripts/capa.py` |
| TICK-DIST-05 | Publicacao pelo Postiz com travas | AC-E05.1, AC-E05.2, AC-E05.3, AC-E05.4, AC-E05.5 | `internal/video/distribuicao/publicar.go` |
| TICK-DIST-06 | Comentarios e aprendizado | AC-E06.1, AC-E06.2 | `internal/video/distribuicao/comentarios.go`, `internal/video/distribuicao/aprendizado.go` |
| TICK-DIST-07 | CLI video-distribuicao e registro na espinha | AC-E07.1, AC-NF.1 | `internal/cli/video_distribuicao.go` |
| TICK-DIST-08 | Canonizar: processo, skill, manual, livro, fechar 9 de 9 | AC-NF.1 | `processos/legoloop-video-distribuicao/pipeline.dag.json`, `skills/legoloop-video-distribuicao/SKILL.md` |

Ordem: 01, 02 (acervo e agenda), 03, 04 (revisao e vitrine), 05 (publicacao), 06, 07 (CLI), 08 (canonizacao).
