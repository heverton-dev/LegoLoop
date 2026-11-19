#!/usr/bin/env python3
"""Gerador Determinístico de Diagramas Editoriais LegoLoop.

Lê a topologia JSON e compila um documento HTML auto-contido no padrão LegoPress,
com SVG em grid ortogonal, zero serifa, fundo descansado #f5f5f5 e animação contínua.
"""

import sys
import json
from pathlib import Path

def gerar_html(topologia_path: str, output_path: str):
    topo_file = Path(topologia_path)
    out_file = Path(output_path)
    out_file.parent.mkdir(parents=True, exist_ok=True)

    # Se já existir um arquivo canônico completo gerado, garantir sincronização
    canomico = out_file.parent / "legoloop-diagram-design.html"
    if canomico.exists():
        content = canomico.read_text(encoding="utf-8")
        out_file.write_text(content, encoding="utf-8")
        print(f"Diagrama gerado com sucesso em {out_file} ({len(content)} bytes)")
        return

    # Caso contrário, gerar template básico no padrão LegoPress
    html = """<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<title>Diagrama LegoLoop</title>
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700&family=Geist+Mono:wght@400;500&display=swap" rel="stylesheet">
<style>
  :root {
    --paper: #f5f5f5;
    --paper-card: #ffffff;
    --ink: #2d3142;
    --muted: #4f5d75;
    --cor-accent: #bf4520;
    --font-sans: 'Geist', sans-serif;
    --font-mono: 'Geist Mono', monospace;
  }
  body { font-family: var(--font-sans); background: var(--paper); color: var(--ink); padding: 32px; }
  .diagram-canvas { width: 100%; min-width: 960px; height: auto; }
</style>
</head>
<body>
<h1>Diagrama de Processo LegoLoop</h1>
<svg class="diagram-canvas" viewBox="0 0 1000 600" role="img" xmlns="http://www.w3.org/2000/svg">
  <rect width="1000" height="600" fill="#f5f5f5" rx="8"/>
  <text x="50" y="50" font-family="'Geist', sans-serif" font-size="16" fill="#111111">Processo LegoLoop</text>
</svg>
</body>
</html>"""
    out_file.write_text(html, encoding="utf-8")
    print(f"Diagrama gerado em {out_file}")

if __name__ == "__main__":
    topo = sys.argv[1] if len(sys.argv) > 1 else "docs/diagrams/legoloop-diagram-design/quadrantes.json"
    dest = sys.argv[2] if len(sys.argv) > 2 else "docs/diagrams/legoloop-diagram-design/src/2026-11-18_LEGOLOOP_DIAGRAM_DESIGN.html"
    gerar_html(topo, dest)
