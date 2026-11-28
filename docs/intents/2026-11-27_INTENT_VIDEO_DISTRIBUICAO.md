# INTENT — legoloop-video-distribuicao (P4 do processo composto legoloop-video-animado)

- **Data:** 2026-11-27
- **Origem:** auditoria `docs/progress/2026-11-27_p4_video_distribuicao/AUDITORIA_P4_VIDEO_DISTRIBUICAO.md` (35 itens, REA do projeto `proj--historias-palitos`).
- **Decisao do usuario:** nome correto `legoloop-video-distribuicao` (antes `diorama-distribuicao`); implementacao autonoma.

## Problema

Depois do render, o video ainda precisa de revisao, vitrine, agenda e publicacao. No original isso era manual: `publicar.py`, `miniatura.py`, `sheet.py` soltos e o `agenda.py` nunca foi escrito.

## Objetivo

Um Processo canonizado que leva o `render.mp4` aprovado ate o canal: revisao com parecer, vitrine com nota, acervo com agenda quente/perene, publicacao pelo Postiz (privada por padrao), comentarios com aprovacao humana e aprendizado semanal.

## Fora de escopo (declarado)

- Cortes verticais 9:16 de um master longo (BL-03, pertence ao P3).
- Kit de canal (ja existe `legoloop-canal-redc`).
- Instalar ou operar o Postiz e o VPS (pre-requisito do usuario).
- Coleta automatica de comentarios e de metricas pelas APIs das redes: o processo le arquivos JSON de entrada.

## Regras inegociaveis

1. Nada vira publico sem `--confirmo-publico` e autorizacao gravada no acervo.
2. Rotulo de conteudo gerado por IA e obrigatorio em todo envio.
3. Segredos so por ambiente ou `.env*`; nunca em artefato, log ou chat.
4. Resposta a comentario so sai depois de aprovada pelo usuario.
5. Testes nunca chamam redes reais.
