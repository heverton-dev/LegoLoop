# Backlog - legoloop-video-animado

## BL-01: faixa etaria "bebe" - pose de engatinhar
- Status: FUTURO, fora do escopo atual do epico E07.
- Decisao do usuario (2026-11-26): bebe nao entra nos personagens de teste agora. Vocabulario aceita o valor "bebe" (ValidarFaixa, internal/video/inventario/genero.go), mas nenhum item do inventario usa esse valor hoje.
- Motivo tecnico: puppetSVG (motor original) nao tem pose de engatinhar; corpo so fica em pe. Pose deitada/engatinhando exige geometria nova, mudanca maior no motor.
- Prova da tentativa: projetos/video-diorama/molde-de-pilula-preview/bebe-teste.png (bebe sentado/em pe, sem gravata via dress onesie).
- Quando retomar: abrir ticket de motor dedicado (pose "engatinhar"/"sentado-no-chao"), com previa visual antes de qualquer gate humano.

## BL-02: estilo palito
- Status: FUTURO, por decisao do usuario (2026-11-26). Fica so o molde-de-pilula por enquanto.
- O que foi feito e desfeito: inventario/video/estilos/palito/estilo.json foi criado (rascunho, herda molde-de-pilula) e removido a pedido do usuario antes da validacao humana.
- Prova visual gerada e preservada para quando retomar: projetos/video-diorama/palito-preview/validacao-estilo-palito.png (3 personagens com as cores do estilo removido).
- Quando retomar: recriar inventario/video/estilos/palito/estilo.json (mesmos dados, ja registrados em docs/progress/2026-11-25_plan-legoloop-diorama-render/achados.md, camada 23) e seguir TICK-V01-08.2.
