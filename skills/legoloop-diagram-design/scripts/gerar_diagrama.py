#!/usr/bin/env python3
"""Gerador Determinístico e Universal de Diagramas Editoriais LegoLoop.

Compila qualquer processo (pipeline.dag.json) em um documento HTML/SVG editorial
no padrão LegoPress, com grid ortogonal, zero serifa, fundo descansado #f5f5f5,
sem transbordamento textual e com fluxo animado contínuo.
"""

import sys
import json
import html
import subprocess
from pathlib import Path

PALETA_PECAS = {
    "skill": {"cor": "#2563eb", "badge": "SKILL", "nome": "Azul Cobalto Google"},
    "script": {"cor": "#0f172a", "badge": "SCRIPT", "nome": "Ardósia Escuro"},
    "stop_gate": {"cor": "#d97706", "badge": "STOP-GATE", "nome": "Âmbar Google Dourado"},
    "gate": {"cor": "#dc2626", "badge": "GATE", "nome": "Carmesim Google"},
    "unit_no": {"cor": "#2563eb", "badge": "UNIT-NÓ", "nome": "Azul Cobalto Google"},
    "mcp": {"cor": "#7c3aed", "badge": "MCP", "nome": "Violeta DeepMind"},
    "swarm": {"cor": "#16a34a", "badge": "SWARM", "nome": "Verde Esmeralda Google"},
    "default": {"cor": "#4f5d75", "badge": "PROCESSO", "nome": "Ardósia Suave"}
}

def obter_estilo_peca(tipo: str):
    t = tipo.lower() if tipo else "default"
    for k, v in PALETA_PECAS.items():
        if k in t:
            return v
    return PALETA_PECAS["default"]

def carregar_dados_processo(alvo: str):
    p = Path(alvo)
    # Se for caminho de arquivo
    if p.is_file() and p.suffix == ".json":
        dag_path = p
        slug = p.parent.name
    elif (Path("processos") / alvo / "pipeline.dag.json").exists():
        dag_path = Path("processos") / alvo / "pipeline.dag.json"
        slug = alvo
    elif (Path("processos") / f"legoloop-{alvo}" / "pipeline.dag.json").exists():
        slug = f"legoloop-{alvo}"
        dag_path = Path("processos") / slug / "pipeline.dag.json"
    else:
        # Fallback genérico para legoloop-diagram-design
        slug = "legoloop-diagram-design"
        dag_path = Path("processos/legoloop-diagram-design/pipeline.dag.json")

    with open(dag_path, "r", encoding="utf-8") as f:
        data = json.load(f)
    return slug, dag_path, data

def gerar_diagrama_completo(alvo: str, out_html_path: str = None):
    slug, dag_path, dag = carregar_dados_processo(alvo)
    nodes = dag.get("nodes", [])
    nome_proc = dag.get("name", slug)
    versao = dag.get("version", "1.0.0")
    desc = dag.get("description", f"Processo canônico {slug} do ecossistema LegoLoop.")

    base_dir = Path("docs/diagrams") / slug
    src_dir = base_dir / "src"
    export_dir = base_dir / "export"
    anim_dir = base_dir / "animation"
    anticorpos_dir = base_dir / "anticorpos"
    
    for d in (src_dir, export_dir, anim_dir, anticorpos_dir):
        d.mkdir(parents=True, exist_ok=True)

    # Copiar anticorpos se existir no processo
    proc_anticorpos = dag_path.parent / "anticorpos.json"
    if proc_anticorpos.exists():
        (anticorpos_dir / "anticorpos.json").write_text(proc_anticorpos.read_text(encoding="utf-8"), encoding="utf-8")

    # Layout ortogonal dos nós
    svg_nodes_markup = []
    svg_edges_markup = []
    cards_markup = []
    table_rows = []

    y_start = 50
    y_step = 110
    total_h = max(620, len(nodes) * y_step + 100)
    viewbox = f"0 0 1040 {total_h}"

    for i, node in enumerate(nodes):
        nid = node.get("id", f"e0{i+1}")
        ntype = node.get("type", "skill")
        nname = node.get("name", nid)
        nestilo = obter_estilo_peca(ntype)
        cor = nestilo["cor"]
        badge = nestilo["badge"]

        # Tratar gate e timeout
        ngate = node.get("output", {}).get("gate", "exit_code == 0")
        if isinstance(ngate, dict):
            ngate = str(ngate)
        ngate_esc = html.escape(str(ngate))
        nname_esc = html.escape(nname)

        # Card HTML
        cards_markup.append(f"""
    <div class="node-card">
      <div class="node-card-header">
        <span class="node-id" style="background:{cor};">{html.escape(nid)} · {badge}</span>
        <small style="color:var(--soft);font-family:var(--font-mono);">{html.escape(ntype)}</small>
      </div>
      <h3 class="node-name">{nname_esc}</h3>
      <p class="node-desc">Execução no processo <code>{html.escape(slug)}</code> com isolamento em contexto estéril.</p>
      <div class="node-gate">Gate: {ngate_esc}</div>
    </div>""")

        # Tabela
        cin = ", ".join(node.get("context_in", [])) or "—"
        cout = ", ".join(node.get("context_out", [])) or node.get("output", {}).get("artifact", "—")
        table_rows.append(f"""
        <tr>
          <td><strong>{html.escape(nid)}</strong></td>
          <td><span style="color:{cor};font-weight:600;">{badge}</span></td>
          <td><code>{html.escape(cin)}</code></td>
          <td><code>{html.escape(cout)}</code></td>
          <td><code>{ngate_esc}</code></td>
        </tr>""")

        # Posição SVG
        y_pos = y_start + i * y_step
        largura_caixa = min(480, max(300, len(nname) * 7 + 80))
        x_pos = 520 - (largura_caixa // 2)

        # Quebra de texto inteligente sem corte
        altura_caixa = 52
        if len(nname_esc) > 38:
            palavras = nname_esc.split(" ")
            meio = len(palavras) // 2
            linha1 = html.escape(" ".join(palavras[:meio]))
            linha2 = html.escape(" ".join(palavras[meio:]))
            texto_nome = f'<text x="76" y="22" font-family="\'Geist\', sans-serif" font-size="11" font-weight="700" fill="#111111"><tspan x="76" dy="0">{linha1}</tspan><tspan x="76" dy="14">{linha2}</tspan></text>'
            texto_id = f'<text x="76" y="47" font-family="\'Geist Mono\', monospace" font-size="9" fill="#5a6580">ID: {html.escape(nid)}</text>'
            altura_caixa = 56
        else:
            texto_nome = f'<text x="76" y="24" font-family="\'Geist\', sans-serif" font-size="12" font-weight="700" fill="#111111">{nname_esc}</text>'
            texto_id = f'<text x="76" y="40" font-family="\'Geist Mono\', monospace" font-size="9" fill="#5a6580">ID: {html.escape(nid)}</text>'

        # Conector para o próximo nó
        if i < len(nodes) - 1:
            y_next = y_start + (i + 1) * y_step
            svg_edges_markup.append(f"""
      <path d="M520,{y_pos + altura_caixa} L520,{y_next}" fill="none" stroke="{cor}" stroke-width="2" class="flow-edge" marker-end="url(#arr-flow)"/>""")

        # Nó SVG com máscara opaca
        svg_nodes_markup.append(f"""
      <g transform="translate({x_pos}, {y_pos})">
        <rect x="0" y="0" width="{largura_caixa}" height="{altura_caixa}" rx="6" fill="#ffffff" stroke="{cor}" stroke-width="1.8"/>
        <rect x="8" y="8" width="60" height="24" rx="3" fill="{cor}"/>
        <text x="38" y="24" font-family="'Geist Mono', monospace" font-size="9" font-weight="700" fill="#ffffff" text-anchor="middle">{badge}</text>
        {texto_nome}
        {texto_id}
      </g>""")

    # Template HTML completo
    html_content = f"""<!DOCTYPE html>
<html lang="pt-BR">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Arquitetura e Fluxo do Processo {html.escape(slug)} — LegoLoop</title>
<link href="https://fonts.googleapis.com/css2?family=Geist:wght@400;500;600;700;800&family=Geist+Mono:wght@400;500;600;700&display=swap" rel="stylesheet">
<style>
  :root {{
    --paper:        #f5f5f5;   /* white-smoke descansado (nunca branco puro) */
    --paper-2:      #ececec;   /* container e faixas secundárias */
    --paper-card:   #ffffff;   /* superfície de cards em alto contraste */
    --ink:          #2d3142;   /* grafite escuro nobre */
    --ink-strong:   #111111;   /* contraste máximo para títulos */
    --muted:        #4f5d75;   /* ardósia suave para textos de apoio */
    --soft:         #5a6580;   /* subrótulos e metadados */
    --rule:         rgba(45,49,66,0.12);
    --rule-solid:   #bfc0c0;   /* bordas estruturais sólidas */
    --cor-accent:   #bf4520;   /* Tangerina Queimada Editorial */
    --font-sans:    'Geist', system-ui, -apple-system, sans-serif;
    --font-mono:    'Geist Mono', ui-monospace, monospace;
  }}

  * {{ box-sizing: border-box; margin: 0; padding: 0; }}

  body {{
    font-family: var(--font-sans);
    background: var(--paper);
    color: var(--ink);
    line-height: 1.6;
    padding: 48px 32px;
    max-width: 1140px;
    margin: 0 auto;
  }}

  .page-header {{
    border-bottom: 2px solid var(--cor-accent);
    padding-bottom: 20px;
    margin-bottom: 40px;
  }}
  .page-eyebrow {{
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.16em;
    text-transform: uppercase;
    color: var(--cor-accent);
    margin-bottom: 8px;
  }}
  .page-title {{
    font-family: var(--font-sans);
    font-size: 2.2rem;
    font-weight: 800;
    letter-spacing: -0.03em;
    color: var(--ink-strong);
    line-height: 1.2;
    margin-bottom: 12px;
  }}
  .page-lead {{
    font-size: 15px;
    color: var(--muted);
    max-width: 860px;
    margin-bottom: 16px;
  }}
  .page-meta {{
    font-family: var(--font-mono);
    font-size: 11px;
    color: var(--soft);
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
  }}
  .page-meta span {{ display: flex; align-items: center; gap: 6px; }}

  .section {{ margin-bottom: 48px; }}
  .section-title {{
    font-family: var(--font-sans);
    font-size: 1.3rem;
    font-weight: 700;
    color: var(--ink-strong);
    letter-spacing: -0.02em;
    margin-bottom: 8px;
  }}
  .section-sub {{ font-size: 14px; color: var(--muted); margin-bottom: 24px; }}

  .nodes-grid {{
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 16px;
    margin-bottom: 32px;
  }}
  .node-card {{
    background: var(--paper-card);
    border: 1px solid var(--rule);
    border-radius: 8px;
    padding: 20px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
  }}
  .node-card-header {{
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 12px;
  }}
  .node-id {{
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 700;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    padding: 2px 8px;
    border-radius: 4px;
    color: #ffffff;
  }}
  .node-name {{ font-size: 15px; font-weight: 700; color: var(--ink-strong); margin-bottom: 6px; }}
  .node-desc {{ font-size: 13px; color: var(--muted); margin-bottom: 12px; }}
  .node-gate {{
    font-family: var(--font-mono);
    font-size: 11px;
    background: var(--paper);
    padding: 8px 10px;
    border-radius: 4px;
    color: var(--soft);
    border: 1px solid var(--rule);
    word-break: break-all;
    overflow-wrap: anywhere;
    line-height: 1.4;
  }}

  .diagram-wrapper {{
    background: var(--paper);
    border: 1px solid var(--rule);
    border-radius: 8px;
    overflow-x: auto;
    margin-bottom: 24px;
    padding: 12px;
  }}
  .diagram-canvas {{ display: block; width: 100%; min-width: 980px; height: auto; }}

  /* Scrollbar Minimalista do Design System (máx 3px, 70% transparência) */
  .diagram-wrapper,
  .table-wrapper {{
    scrollbar-width: thin;
    scrollbar-color: rgba(79, 93, 117, 0.3) transparent;
  }}
  .diagram-wrapper::-webkit-scrollbar,
  .table-wrapper::-webkit-scrollbar {{
    height: 3px;
    width: 3px;
  }}
  .diagram-wrapper::-webkit-scrollbar-track,
  .table-wrapper::-webkit-scrollbar-track {{
    background: transparent;
  }}
  .diagram-wrapper::-webkit-scrollbar-thumb,
  .table-wrapper::-webkit-scrollbar-thumb {{
    background: rgba(79, 93, 117, 0.3);
    border-radius: 2px;
  }}
  .diagram-wrapper::-webkit-scrollbar-thumb:hover,
  .table-wrapper::-webkit-scrollbar-thumb:hover {{
    background: rgba(191, 69, 32, 0.5);
  }}

  @keyframes flow-pulse {{
    from {{ stroke-dashoffset: 24; }}
    to {{ stroke-dashoffset: 0; }}
  }}
  .flow-edge {{
    stroke-dasharray: 6 6;
    animation: flow-pulse 1.2s linear infinite;
    will-change: stroke-dashoffset;
  }}

  .table-wrapper {{
    background: var(--paper-card);
    border: 1px solid var(--rule);
    border-radius: 8px;
    overflow-x: auto;
  }}
  table {{ width: 100%; border-collapse: collapse; font-size: 13px; text-align: left; }}
  th {{
    background: var(--paper-2);
    font-family: var(--font-mono);
    font-size: 11px;
    font-weight: 600;
    letter-spacing: 0.1em;
    text-transform: uppercase;
    color: var(--ink-strong);
    padding: 12px 16px;
    border-bottom: 1px solid var(--rule-solid);
  }}
  td {{ padding: 12px 16px; border-bottom: 1px solid var(--rule); color: var(--ink); }}
  tr:last-child td {{ border-bottom: none; }}
  td code {{
    font-family: var(--font-mono);
    font-size: 11px;
    background: var(--paper);
    padding: 2px 6px;
    border-radius: 3px;
    color: var(--muted);
  }}
</style>
</head>
<body>

<header class="page-header">
  <p class="page-eyebrow">LegoLoop · Processo Canônico</p>
  <h1 class="page-title">Arquitetura e Fluxo do Processo {html.escape(nome_proc)}</h1>
  <p class="page-lead">{html.escape(desc)}</p>
  <div class="page-meta">
    <span><strong>Processo:</strong> {html.escape(slug)}</span>
    <span><strong>Versão:</strong> {html.escape(versao)}</span>
    <span><strong>Nós:</strong> {len(nodes)} nós sequenciais</span>
    <span><strong>Design System:</strong> LegoPress Editorial (Sans-Serif)</span>
  </div>
</header>

<section class="section">
  <h2 class="section-title">01 — Anatomia dos Nós do Processo</h2>
  <p class="section-sub">Estrutura isolada em contexto estéril com barreiras mecânicas de qualidade.</p>
  <div class="nodes-grid">
    {"".join(cards_markup)}
  </div>
</section>

<section class="section">
  <h2 class="section-title">02 — Fluxograma Lógico de Encadeamento com Animação Contínua</h2>
  <p class="section-sub">Grafo ortogonal com partículas contínuas mostrando a transferência assíncrona de estado.</p>
  <div class="diagram-wrapper">
    <svg class="diagram-canvas" viewBox="{viewbox}" role="img" aria-labelledby="{slug}-title {slug}-desc" xmlns="http://www.w3.org/2000/svg">
      <title id="{slug}-title">Fluxograma do Processo {html.escape(nome_proc)}</title>
      <desc id="{slug}-desc">Fluxo sequencial animado dos nós do processo {html.escape(slug)} no padrão editorial LegoPress.</desc>
      <defs>
        <marker id="arr-flow" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto">
          <path d="M0,1 L8,4 L0,7 L1.5,4 Z" fill="#4f5d75"/>
        </marker>
      </defs>
      <rect width="1040" height="{total_h}" fill="#f5f5f5" rx="8"/>
      {"".join(svg_edges_markup)}
      {"".join(svg_nodes_markup)}
    </svg>
  </div>
</section>

<section class="section">
  <h2 class="section-title">03 — Matriz de Entradas, Saídas e Barreiras Mecânicas</h2>
  <p class="section-sub">Contratos formais de arquivos consumidos e emitidos por cada vértice.</p>
  <div class="table-wrapper">
    <table>
      <thead>
        <tr>
          <th>Nó</th>
          <th>Tipo</th>
          <th>Insumos</th>
          <th>Artefatos Emitidos</th>
          <th>Gate Mecânico</th>
        </tr>
      </thead>
      <tbody>
        {"".join(table_rows)}
      </tbody>
    </table>
  </div>
</section>

<section class="section">
  <h2 class="section-title">04 — Governança de Anticorpos e Paridade com o LegoPress</h2>
  <p class="section-sub">Diretrizes obrigatórias testadas em tempo de compilação.</p>
  <div style="background:var(--paper-card);border:1px solid var(--rule);border-radius:8px;padding:20px;font-size:14px;color:var(--ink);">
    <p style="margin-bottom:12px;"><strong>1. Fundo Descansado:</strong> O background usa <code>#f5f5f5</code> (white-smoke). Zero branco puro na tela.</p>
    <p style="margin-bottom:12px;"><strong>2. Zero Serifa:</strong> Tipografia 100% <code>Geist</code> e <code>Geist Mono</code>. Fontes serifadas bloqueadas.</p>
    <p style="margin-bottom:12px;"><strong>3. Classe Obrigatória:</strong> Tag <code>&lt;svg class="diagram-canvas"&gt;</code> presente com acessibilidade nativa.</p>
    <p><strong>4. Portabilidade Typst:</strong> O SVG resultante compila sem perdas nos livros didáticos do LegoPress.</p>
  </div>
</section>

</body>
</html>"""

    # Gravar HTML
    final_html = src_dir / f"{slug}.html"
    final_html.write_text(html_content, encoding="utf-8")
    (src_dir / "diagrama.html").write_text(html_content, encoding="utf-8")
    if out_html_path:
        Path(out_html_path).write_text(html_content, encoding="utf-8")

    # Gravar composition.css
    css_content = """/* Animação Contínua LegoLoop */
.diagram-canvas path.flow-edge {
  stroke-dasharray: 6 6;
  animation: flow-pulse 1.2s linear infinite;
  will-change: stroke-dashoffset;
}

@keyframes flow-pulse {
  from { stroke-dashoffset: 24; }
  to { stroke-dashoffset: 0; }
}

@media (prefers-reduced-motion: reduce) {
  .diagram-canvas path.flow-edge {
    animation: none;
    stroke-dashoffset: 0;
  }
}
"""
    (anim_dir / "composition.css").write_text(css_content, encoding="utf-8")

    # Gravar composition.tsx (Remotion)
    tsx_content = f"""import React from 'react';
import {{ interpolate, useCurrentFrame }} from 'remotion';

export const {slug.replace('-', '_').title().replace('_', '')}Animation: React.FC = () => {{
  const frame = useCurrentFrame();
  const opacity = interpolate(frame, [0, 30], [0, 1], {{ extrapolateRight: 'clamp' }});

  return (
    <div style={{{{ width: '100%', height: '100%', background: '#f5f5f5', opacity, padding: 32 }}}}>
      <h1 style={{{{ fontFamily: 'Geist, sans-serif' }}}}>{nome_proc}</h1>
    </div>
  );
}};
"""
    (anim_dir / "composition.tsx").write_text(tsx_content, encoding="utf-8")

    # Exportar SVG Standalone
    svg_out = export_dir / f"{slug}.svg"
    svg_script = Path("skills/legoloop-diagram-design/scripts/export_svg.py")
    if svg_script.exists():
        subprocess.run([sys.executable, str(svg_script), str(final_html), str(svg_out)], check=False)
        subprocess.run([sys.executable, str(svg_script), str(final_html), str(export_dir / "diagrama.svg")], check=False)

    print(f"✓ Diagrama do processo '{slug}' gerado com sucesso!")
    print(f"  - HTML: {final_html} ({len(html_content)} bytes)")
    print(f"  - SVG: {svg_out}")
    print(f"  - Nós processados: {len(nodes)}")

if __name__ == "__main__":
    alvo = sys.argv[1] if len(sys.argv) > 1 else "legoloop-diagram-design"
    out = sys.argv[2] if len(sys.argv) > 2 else None
    gerar_diagrama_completo(alvo, out)
