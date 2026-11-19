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

def encurtar_caminho(val, max_len=24):
    if not val or val == "—":
        return "—"
    if isinstance(val, list):
        val = ", ".join(val)
    val = str(val).strip()
    if len(val) <= max_len:
        return val
    p = Path(val)
    if len(p.name) > 0 and len(p.name) <= max_len:
        return p.name
    return val[:max_len-3] + "..."

def extrair_subtitulos(nname, nid, max_len=24):
    if not nname or nname == nid:
        return nid.replace("_", " ").title(), ""
    palavras = nname.split()
    l1 = []
    l2 = []
    for w in palavras:
        if not l2 and len(" ".join(l1 + [w])) <= max_len:
            l1.append(w)
        elif len(" ".join(l2 + [w])) <= max_len:
            l2.append(w)
        else:
            break
    sub1 = " ".join(l1) if l1 else nname[:max_len]
    sub2 = " ".join(l2) if l2 else ""
    return sub1, sub2

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
        if p.name == "pipeline.dag.json":
            dag_path = p
            slug = p.parent.name
        elif p.name in ("quadrantes.json", "manifest.json"):
            slug = p.parent.name
            dag_candidate = Path("processos") / slug / "pipeline.dag.json"
            dag_path = dag_candidate if dag_candidate.exists() else p
        else:
            dag_path = p
            slug = p.stem
    elif p.is_dir() and (p / "pipeline.dag.json").exists():
        dag_path = p / "pipeline.dag.json"
        slug = p.name
    elif (Path("processos") / alvo / "pipeline.dag.json").exists():
        dag_path = Path("processos") / alvo / "pipeline.dag.json"
        slug = Path(alvo).name
    elif (Path("processos") / f"legoloop-{alvo}" / "pipeline.dag.json").exists():
        slug = f"legoloop-{Path(alvo).name}"
        dag_path = Path("processos") / slug / "pipeline.dag.json"
    else:
        # Fallback genérico para legoloop-diagram-design
        slug = "legoloop-diagram-design"
        dag_path = Path("processos/legoloop-diagram-design/pipeline.dag.json")

    # Higienização estrita de slug para isolamento de diretório
    slug = Path(slug).name

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
    assets_dir = base_dir / "assets"
    export_dir = base_dir / "export"
    anim_dir = base_dir / "animation"
    anticorpos_dir = base_dir / "anticorpos"
    
    for d in (src_dir, assets_dir, export_dir, anim_dir, anticorpos_dir):
        d.mkdir(parents=True, exist_ok=True)

    # Copiar anticorpos se existir no processo, ou inicializar vazio
    proc_anticorpos = dag_path.parent / "anticorpos.json"
    anticorpos_file = anticorpos_dir / "anticorpos.json"
    if proc_anticorpos.exists():
        anticorpos_file.write_text(proc_anticorpos.read_text(encoding="utf-8"), encoding="utf-8")
    elif not anticorpos_file.exists():
        anticorpos_file.write_text(json.dumps({"anticorpos": []}, indent=2), encoding="utf-8")

    # Gravar manifest.json e quadrantes.json de governança
    manifest_file = base_dir / "manifest.json"
    manifest_data = {
        "slug": slug,
        "name": nome_proc,
        "version": versao,
        "nodes_count": len(nodes),
        "visual_type": "process",
        "semantic_pattern": "stage-framework",
        "generated_at": "2026-11-19T07:00:00Z"
    }
    manifest_file.write_text(json.dumps(manifest_data, indent=2), encoding="utf-8")

    quadrantes_file = base_dir / "quadrantes.json"
    quadrantes_data = {
        "slug": slug,
        "nodes": [
            {
                "id": n.get("id"),
                "name": n.get("name"),
                "type": n.get("type", "skill"),
                "badge": obter_estilo_peca(n.get("type", "skill"))["badge"],
                "color": obter_estilo_peca(n.get("type", "skill"))["cor"]
            }
            for n in nodes
        ]
    }
    quadrantes_file.write_text(json.dumps(quadrantes_data, indent=2), encoding="utf-8")

    # Layout Híbrido Horizontal e Multilinhas dos nós
    svg_nodes_markup = []
    svg_edges_markup = []
    cards_markup = []
    table_rows = []

    num_nodes = len(nodes)
    if num_nodes <= 1:
        cols = 1
    elif num_nodes <= 4:
        cols = num_nodes
    elif num_nodes in (5, 6, 9):
        cols = 3
    else:
        cols = 4

    rows = (num_nodes + cols - 1) // cols if cols > 0 else 1

    if cols == 4:
        card_w = 216
        gap_x = 40
        margin_x = 24
    elif cols == 3:
        card_w = 296
        gap_x = 44
        margin_x = 32
    elif cols == 2:
        card_w = 440
        gap_x = 64
        margin_x = 48
    else:
        card_w = 560
        gap_x = 40
        margin_x = 240

    card_h = 240
    gap_y = 60
    margin_y = 36
    total_w = 1040
    total_h = max(320, margin_y * 2 + rows * card_h + (rows - 1) * gap_y)
    viewbox = f"0 0 {total_w} {total_h}"

    posicoes = []
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

        # Insumos e Saídas tratados
        cin_raw = node.get("context_in") or node.get("inputs") or "—"
        cout_raw = node.get("context_out") or node.get("output", {}).get("artifact") or "—"
        cin_curto = encurtar_caminho(cin_raw, max_len=20)
        cout_curto = encurtar_caminho(cout_raw, max_len=20)
        ngate_str = str(ngate)
        gate_curto = ngate_str[:22] + ("..." if len(ngate_str) > 22 else "")
        gate_curto_esc = html.escape(gate_curto)

        sub1, sub2 = extrair_subtitulos(nname, nid, max_len=24)

        # Card HTML na grade editorial
        cards_markup.append(f"""
    <div class="node-card">
      <div class="node-card-header">
        <span class="node-id" style="background:{cor};">{html.escape(nid)} · {badge}</span>
        <small style="color:var(--soft);font-family:var(--font-mono);">{html.escape(ntype)}</small>
      </div>
      <h3 class="node-name">{nname_esc}</h3>
      <p class="node-desc">Execução no processo <code>{html.escape(slug)}</code> com isolamento em contexto estéril.</p>
      <div style="font-family:var(--font-mono);font-size:11px;color:var(--soft);margin-bottom:10px;">
        <div><strong>In:</strong> <code>{html.escape(cin_curto)}</code></div>
        <div><strong>Out:</strong> <code>{html.escape(cout_curto)}</code></div>
      </div>
      <div class="node-gate">Gate: {ngate_esc}</div>
    </div>""")

        # Tabela
        cin_tab = ", ".join(node.get("context_in", [])) or "—"
        cout_tab = ", ".join(node.get("context_out", [])) or node.get("output", {}).get("artifact", "—")
        table_rows.append(f"""
        <tr>
          <td><strong>{html.escape(nid)}</strong></td>
          <td><span style="color:{cor};font-weight:600;">{badge}</span></td>
          <td><code>{html.escape(cin_tab)}</code></td>
          <td><code>{html.escape(cout_tab)}</code></td>
          <td><code>{ngate_esc}</code></td>
        </tr>""")

        # Posição na grade SVG
        r = i // cols
        c = i % cols
        x_pos = margin_x + c * (card_w + gap_x)
        y_pos = margin_y + r * (card_h + gap_y)
        posicoes.append((x_pos, y_pos, r, c))

        border_stroke = cor if badge in ("STOP-GATE", "UNIT-NÓ") else "#bfc0c0"
        border_w = "1.5" if badge in ("STOP-GATE", "UNIT-NÓ") else "1"

        # Nó SVG em card envolvente de 240px de altura idêntico ao print
        svg_nodes_markup.append(f"""
      <g class="node" id="node-{html.escape(nid)}" transform="translate({x_pos}, {y_pos})">
        <rect width="{card_w}" height="{card_h}" rx="8" fill="#ffffff" stroke="{border_stroke}" stroke-width="{border_w}"/>
        <rect x="14" y="14" width="68" height="20" rx="3" fill="{cor}"/>
        <text x="48" y="28" text-anchor="middle" font-family="'Geist Mono', monospace" font-size="9" font-weight="700" fill="#ffffff">{badge}</text>
        <text x="14" y="56" font-family="'Geist', sans-serif" font-size="13" font-weight="700" fill="#0f172a">{html.escape(nid)}</text>
        <text x="14" y="78" font-family="'Geist', sans-serif" font-size="11" fill="#4f5d75">{html.escape(sub1)}</text>
        <text x="14" y="96" font-family="'Geist', sans-serif" font-size="10" fill="#5a6580">{html.escape(sub2)}</text>
        <line x1="14" y1="116" x2="{card_w - 14}" y2="116" stroke="#bfc0c0" stroke-width="0.8"/>
        <text x="14" y="138" font-family="'Geist Mono', monospace" font-size="9" fill="#5a6580">In: {html.escape(cin_curto)}</text>
        <text x="14" y="156" font-family="'Geist Mono', monospace" font-size="9" fill="#5a6580">Out: {html.escape(cout_curto)}</text>
        <rect x="14" y="196" width="{card_w - 28}" height="30" rx="4" fill="#f5f5f5" stroke="rgba(45,49,66,0.1)"/>
        <text x="22" y="215" font-family="'Geist Mono', monospace" font-size="8.5" fill="#2d3142">Gate: {gate_curto_esc}</text>
      </g>""")

    # Conectores entre nós
    for i in range(len(nodes) - 1):
        x1, y1, r1, c1 = posicoes[i]
        x2, y2, r2, c2 = posicoes[i+1]
        ntype1 = nodes[i].get("type", "skill")
        cor_edge = obter_estilo_peca(ntype1)["cor"]

        if r1 == r2:
            start_x = x1 + card_w
            start_y = y1 + card_h // 2
            end_x = x2
            end_y = y2 + card_h // 2
            svg_edges_markup.append(f"""
      <path d="M{start_x},{start_y} L{end_x},{end_y}" fill="none" stroke="{cor_edge}" stroke-width="2" class="flow-edge" marker-end="url(#{slug}-arr-flow)"/>""")
        else:
            start_x = x1 + card_w
            start_y = y1 + card_h // 2
            mid_x = start_x + 16
            gutter_y = y1 + card_h + gap_y // 2
            entry_x = x2 - 14
            end_y = y2 + card_h // 2
            path_d = f"M{start_x},{start_y} L{mid_x},{start_y} L{mid_x},{gutter_y} L{entry_x},{gutter_y} L{entry_x},{end_y} L{x2},{end_y}"
            svg_edges_markup.append(f"""
      <path d="{path_d}" fill="none" stroke="{cor_edge}" stroke-width="2" stroke-linejoin="round" class="flow-edge" marker-end="url(#{slug}-arr-flow)"/>""")

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
    min-width: 0;
    overflow: hidden;
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
  .node-name {{ font-size: 15px; font-weight: 700; color: var(--ink-strong); margin-bottom: 6px; word-break: break-word; overflow-wrap: anywhere; }}
  .node-desc {{ font-size: 13px; color: var(--muted); margin-bottom: 12px; word-break: break-word; }}
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
