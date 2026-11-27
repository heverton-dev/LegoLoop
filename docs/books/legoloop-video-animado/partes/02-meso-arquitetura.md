# 2. Arquitetura Meso

```
internal/video/
  inventario/   -> estilo.json, itens/, herança, validação humana (E07)
  ativos/       -> pré-voo v2, elenco vetorial, cenários (código/Flow/híbrido), entrega suspeita (E08/E09)
  motor/        -> ponte cenas->timeline, pré-voo do motor, render de referência (E10)
  render/       -> render Go legado, agora também liga ao motor novo (E10.4)
  roteiro/      -> (ex-diorama/roteiro, só renomeado, E11)

motor-local/    -> engine.html, puppet.js, sets.js, render.py (licença MIT, fora do rascunho pendente)
```

Fluxo: `video-inventario` valida o estilo -> `video-ativos` monta cenas.json do inventário -> `video-render timeline` chama a ponte do motor -> `video-render executar` chama o motor real (Playwright + ffmpeg) e mede o MP4.

Alias de compatibilidade: `diorama-roteiro`/`diorama-ativos` continuam funcionando, com aviso em stderr e o mesmo exit code dos comandos novos (`video-roteiro`/`video-ativos`).
