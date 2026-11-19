#!/usr/bin/env python3
"""Gerador da Demonstração de Pipeline Híbrido com 10 Nós.

Demonstra como o algoritmo resolve processos extensos (10 nós)
em um grid horizontal multilinhas (4 + 4 + 2 cards) com conectores perimetrais.
"""

import html
from pathlib import Path

NODES = [
    {"id": "e01_concepcao", "nome": "Concepção & Intenção", "badge": "SKILL", "cor": "#2563eb", "gate": "confianca >= 95%", "fase": "FASE 1"},
    {"id": "e02_especificacao", "nome": "Especificação EARS", "badge": "SKILL", "cor": "#2563eb", "gate": "spec.v1.json", "fase": "FASE 1"},
    {"id": "e03_planejamento", "nome": "Decomposição Tickets", "badge": "SKILL", "cor": "#2563eb", "gate": "plan.v1.json", "fase": "FASE 1"},
    {"id": "e04_roteamento", "nome": "Mapeamento DAG & Nós", "badge": "SCRIPT", "cor": "#0f172a", "gate": "pipeline.dag.json", "fase": "FASE 2"},
    {"id": "e05_tdd_mutacao", "nome": "TDD e Mutação", "badge": "SCRIPT", "cor": "#0f172a", "gate": "mutation_fails == true", "fase": "FASE 2"},
    {"id": "e06_orchestrix", "nome": "Execução Estéril", "badge": "SCRIPT", "cor": "#0f172a", "gate": "exit_code == 0", "fase": "FASE 2"},
    {"id": "e07_qa_inspecao", "nome": "Inspeção & Linter", "badge": "SCRIPT", "cor": "#0f172a", "gate": "zero_desvios == true", "fase": "FASE 3"},
    {"id": "e08_stop_gate", "nome": "Stop-Gate Humano", "badge": "STOP-GATE", "cor": "#d97706", "gate": "prompt(Aprovar?)", "fase": "FASE 3"},
    {"id": "e09_canonizacao", "nome": "Deduplicação Cofre", "badge": "SCRIPT", "cor": "#0f172a", "gate": "sha256_match == true", "fase": "FASE 4"},
    {"id": "e10_fechamento", "nome": "Fechamento 9/9", "badge": "SCRIPT", "cor": "#0f172a", "gate": "report_audit_9_de_9", "fase": "FASE 4"}
]

def gerar_svg_10_nos():
    cols = 4
    card_w = 216
    card_h = 136
    gap_x = 40
    gap_y = 56
    margin_x = 24
    margin_y = 32

    total_w = 1040
    total_h = 580

    svg_cards = []
    svg_edges = []

    posicoes = []
    for i, node in enumerate(NODES):
        r = i // cols
        c = i % cols
        x = margin_x + c * (card_w + gap_x)
        y = margin_y + r * (card_h + gap_y)
        posicoes.append((x, y, r, c))

        cor = node["cor"]
        badge = node["badge"]
        nid = node["id"]
        nnome = html.escape(node["nome"])
        ngate = html.escape(node["gate"])

        border_stroke = cor if badge in ("STOP-GATE", "UNIT-NÓ") else "#bfc0c0"
        border_w = "1.4" if badge in ("STOP-GATE", "UNIT-NÓ") else "1"

        card = f"""
    <!-- CARD {i+1}: {nid} -->
    <g transform="translate({x}, {y})" filter="url(#shadow-10)">
      <rect width="{card_w}" height="{card_h}" rx="8" fill="#ffffff" stroke="{border_stroke}" stroke-width="{border_w}"/>
      <rect x="12" y="12" width="64" height="18" rx="3" fill="{cor}"/>
      <text x="44" y="24" text-anchor="middle" font-family="'Geist Mono',monospace" font-size="8.5" font-weight="700" fill="#ffffff">{badge}</text>
      <text x="82" y="24" font-family="'Geist Mono',monospace" font-size="9" font-weight="600" fill="#5a6580">{node['fase']}</text>
      <text x="12" y="50" font-family="'Geist',sans-serif" font-size="12" font-weight="700" fill="#0f172a">{nnome}</text>
      <text x="12" y="68" font-family="'Geist Mono',monospace" font-size="10" fill="#4f5d75">{nid}</text>
      <line x1="12" y1="84" x2="{card_w - 12}" y2="84" stroke="#bfc0c0" stroke-width="0.8"/>
      <rect x="12" y="96" width="{card_w - 24}" height="24" rx="4" fill="#f5f5f5" stroke="rgba(45,49,66,0.08)"/>
      <text x="18" y="112" font-family="'Geist Mono',monospace" font-size="8.5" fill="#2d3142">Gate: {ngate}</text>
    </g>"""
        svg_cards.append(card)

    # Conectores entre nós
    for i in range(len(NODES) - 1):
        x1, y1, r1, c1 = posicoes[i]
        x2, y2, r2, c2 = posicoes[i+1]

        # Se for na mesma linha
        if r1 == r2:
            start_x = x1 + card_w
            start_y = y1 + card_h // 2
            end_x = x2
            end_y = y2 + card_h // 2
            edge = f"""  <path d="M{start_x},{start_y} L{end_x},{end_y}" fill="none" stroke="#4f5d75" stroke-width="2" class="flow-edge" marker-end="url(#arr-10)"/>"""
            svg_edges.append(edge)
        else:
            # Transição de final de linha para início da linha seguinte (Perimetral em U)
            # Sai do lado direito do último card da linha 1
            start_x = x1 + card_w
            start_y = y1 + card_h // 2
            mid_x = start_x + 16
            gutter_y = y1 + card_h + gap_y // 2
            entry_x = x2 - 14
            end_y = y2 + card_h // 2

            # Traçado ortogonal limpo contornando por baixo da linha
            path_d = f"M{start_x},{start_y} L{mid_x},{start_y} L{mid_x},{gutter_y} L{entry_x},{gutter_y} L{entry_x},{end_y} L{x2},{end_y}"
            edge = f"""  <!-- Conexão entre linhas: Nó {i+1} -> Nó {i+2} -->
  <path d="{path_d}" fill="none" stroke="#2563eb" stroke-width="2" stroke-linejoin="round" class="flow-edge" marker-end="url(#arr-10)"/>"""
            svg_edges.append(edge)

    svg_content = f"""<svg class="diagram-canvas" viewBox="0 0 {total_w} {total_h}" role="img" xmlns="http://www.w3.org/2000/svg">
  <defs>
    <marker id="arr-10" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto">
      <path d="M0,1 L8,4 L0,7 L1.5,4 Z" fill="#4f5d75"/>
    </marker>
    <filter id="shadow-10" x="-5%" y="-8%" width="110%" height="120%">
      <feDropShadow dx="0" dy="2" stdDeviation="3" flood-opacity="0.05"/>
    </filter>
  </defs>

  <rect width="{total_w}" height="{total_h}" fill="#f5f5f5" rx="8"/>

  <!-- CONECTORES ENTRE NÓS -->
{''.join(svg_edges)}

  <!-- CARDS DOS NÓS -->
{''.join(svg_cards)}
</svg>"""
    return svg_content

def main():
    out_dir = Path("docs/diagrams/legoloop-concepcao/src")
    out_dir.mkdir(parents=True, exist_ok=True)
    html_file = out_dir / "demonstracao_10_nos_hibrido.html"

    svg_markup = gerar_svg_10_nos()

    page_html = f"""<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Demonstração Híbrida com 10 Nós — LegoLoop</title>
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700;800&family=Geist+Mono:wght@400;500;600;700&display=swap" rel="stylesheet">
<style>
  :root {{
    --paper:        #f5f5f5;
    --paper-card:   #ffffff;
    --ink:          #2d3142;
    --ink-strong:   #0f172a;
    --muted:        #4f5d75;
    --soft:         #5a6580;
    --rule-solid:   #bfc0c0;
    --accent:       #bf4520;
    --font-sans:    'Geist', system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    --font-mono:    'Geist Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
  }}
  * {{ box-sizing: border-box; margin: 0; padding: 0; }}
  body {{
    font-family: var(--font-sans);
    background: var(--paper);
    color: var(--ink);
    line-height: 1.5;
    padding: 32px;
    max-width: 1140px;
    margin: 0 auto;
  }}
  .page-header {{
    border-bottom: 2px solid var(--accent);
    padding-bottom: 16px;
    margin-bottom: 24px;
  }}
  .page-eyebrow {{
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--accent);
    margin-bottom: 6px;
  }}
  .page-title {{
    font-size: 1.8rem;
    font-weight: 800;
    letter-spacing: -0.03em;
    color: var(--ink-strong);
    margin-bottom: 8px;
  }}
  .page-desc {{
    font-size: 14px;
    color: var(--muted);
  }}
  .diagram-container {{
    background: var(--paper-card);
    border: 1px solid var(--rule-solid);
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.03);
    margin-bottom: 24px;
    overflow-x: auto;
  }}
  .diagram-canvas {{
    display: block;
    width: 100%;
    height: auto;
    max-width: 1040px;
    margin: 0 auto;
  }}
  .flow-edge {{
    stroke-dasharray: 6 3;
    animation: flow-dash 1.2s linear infinite;
  }}
  @keyframes flow-dash {{
    to {{ stroke-dashoffset: -18; }}
  }}
  .explaining-box {{
    background: #ffffff;
    border: 1px solid var(--rule-solid);
    border-left: 4px solid var(--accent);
    padding: 16px 20px;
    border-radius: 6px;
    font-size: 13.5px;
    color: var(--ink);
    line-height: 1.6;
  }}
</style>
</head>
<body>
  <header class="page-header">
    <div class="page-eyebrow">ARQUITETURA DE DESIGN SYSTEM · ESCALABILIDADE DE NÓS</div>
    <h1 class="page-title">Demonstração da Abordagem Híbrida: Processo com 10 Nós</h1>
    <p class="page-desc">Combinação da elegância horizontal (Opção 2) com o roteamento matricial perimetral (Opção 3) em 3 linhas com largura constante de 1040px.</p>
  </header>

  <div class="diagram-container">
    {svg_markup}
  </div>

  <div class="explaining-box">
    <strong>Como o algoritmo se comporta:</strong>
    <br>• Linha 1 (4 nós): Da esquerda para a direita (e01 → e02 → e03 → e04).
    <br>• Conexão perimetral: Desce da lateral direita da Linha 1 e entra pela esquerda da Linha 2.
    <br>• Linha 2 (4 nós): Da esquerda para a direita (e05 → e06 → e07 → e08).
    <br>• Linha 3 (2 nós): Da esquerda para a direita (e09 → e10).
    <br>• <strong>Vantagem:</strong> Mantém a largura perfeita de 1040px sem rolagem horizontal nem vertical excessiva.
  </div>
</body>
</html>"""

    html_file.write_text(page_html, encoding="utf-8")
    print(f"[OK] Demonstração gerada em: {html_file}")

if __name__ == "__main__":
    main()
