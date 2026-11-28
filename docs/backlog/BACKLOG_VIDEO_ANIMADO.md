# Backlog - legoloop-video-animado

## BL-01: faixa etaria "bebe" - pose de engatinhar
- Status: FUTURO, fora do escopo atual do epico E07.
- Decisao do usuario (2026-11-26): bebe nao entra nos personagens de teste agora. Vocabulario aceita o valor "bebe" (ValidarFaixa, internal/video/inventario/genero.go), mas nenhum item do inventario usa esse valor hoje.
- Motivo tecnico: puppetSVG (motor original) nao tem pose de engatinhar; corpo so fica em pe. Pose deitada/engatinhando exige geometria nova, mudanca maior no motor.
- Prova da tentativa (apagada na limpeza de 2026-11-27 a pedido do usuario; refazer a previa ao retomar): projetos/video-animado/molde-de-pilula-preview/bebe-teste.png (bebe sentado/em pe, sem gravata via dress onesie).
- Quando retomar: abrir ticket de motor dedicado (pose "engatinhar"/"sentado-no-chao"), com previa visual antes de qualquer gate humano.

## BL-02: estilo palito
- Status: FUTURO, por decisao do usuario (2026-11-26). Fica so o molde-de-pilula por enquanto.
- O que foi feito e desfeito: inventario/video/estilos/palito/estilo.json foi criado (rascunho, herda molde-de-pilula) e removido a pedido do usuario antes da validacao humana.
- Prova visual (apagada na limpeza de 2026-11-27 a pedido do usuario; os dados do estilo seguem no achados.md): projetos/video-animado/palito-preview/validacao-estilo-palito.png (3 personagens com as cores do estilo removido).
- Quando retomar: recriar inventario/video/estilos/palito/estilo.json (mesmos dados, ja registrados em docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md, camada 23) e seguir TICK-V01-08.2.

## BL-03: cortes verticais 9:16 a partir do master longo (Etapa 8b)
- Status: FUTURO. Registrado no ticket 2b (2026-11-27) ao canonizar o legoloop-video-render.
- O que falta: SPEC AC-E05.4 e AC-E05.5 (janela 9:16 que contem o narrador, legenda fora da caixa dele, saida de 15 a 59 s, sem short quando a timeline tem menos de 15 s). Nao existe funcao de corte vertical em internal/video/render; so o erro tipado ErrCorteVertical.
- Motivo do adiamento: o motor ja renderiza o video em 9:16 com voz; o corte so importa para um master longo 16:9, que ainda nao e o formato de producao.
- Quando retomar: ticket proprio com TDD, depois de existir um master longo 16:9 pelo motor.
