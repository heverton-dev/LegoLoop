# Backlog - legoloop-video-animado

## BL-01: faixa etaria "bebe" - pose de engatinhar
- Status: FUTURO, fora do escopo atual do epico E07.
- Decisao do usuario (2026-11-26): bebe nao entra nos personagens de teste agora. Vocabulario aceita o valor "bebe" (ValidarFaixa, internal/video/inventario/genero.go), mas nenhum item do inventario usa esse valor hoje.
- Motivo tecnico: puppetSVG (motor original) nao tem pose de engatinhar; corpo so fica em pe. Pose deitada/engatinhando exige geometria nova, mudanca maior no motor.
- Prova da tentativa: projetos/video-diorama/molde-de-pilula-preview/bebe-teste.png (bebe sentado/em pe, sem gravata via dress onesie).
- Quando retomar: abrir ticket de motor dedicado (pose "engatinhar"/"sentado-no-chao"), com previa visual antes de qualquer gate humano.
