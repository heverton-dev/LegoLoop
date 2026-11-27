# 1. Fundação Macro

O macroprocesso `legoloop-video-animado` substitui o estilo único por um **inventário global por estilo** (`internal/video/inventario`): cada estilo (`molde-de-pilula`) guarda referências, paleta, luz de contorno e itens aprovados (pose, humor, cenário), com herança entre estilos.

O `projeto.json` ganha o contrato v2 (`modo_producao` + `estilo`), com o contrato legado (`estilo_visual`) ainda lido com aviso de depreciação.

A Peça 2 (elenco e cenários) passa a montar o elenco direto do inventário, sem depender do Flow para estilos já aprovados, mantendo o caminho de recorte (rembg + despill) para estilos que ainda usam fichas.

O motor de render original (`puppet.js`, `sets.js`, SVG 100% autoral) foi portado para `motor-local/` sob licença MIT do próprio autor, ligado ao render Go por uma ponte (`internal/video/motor`) que traduz `cenas.json` + elenco + timings + inventário em `timeline.json`.
