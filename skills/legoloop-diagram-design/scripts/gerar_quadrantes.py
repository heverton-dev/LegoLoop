#!/usr/bin/env python3
"""Ingestão e Estruturação Determinística de Quadrantes (Nó e02).

Lê o pipeline.dag.json ou manifest.json do processo alvo e emite
docs/diagrams/<slug>/quadrantes.json com categorização taxonômica e geométrica.
"""

import sys
import re
import json
from pathlib import Path

def extrair_nomes_arquivos_ou_insumos(valor) -> str:
    if not valor or valor == "—":
        return "—"
    itens = []
    def processar_item(v):
        if not v:
            return
        if isinstance(v, (list, tuple)):
            for sub in v:
                processar_item(sub)
            return
        if isinstance(v, dict):
            for k, val in v.items():
                processar_item(val)
            return
        s = str(v).strip()
        file_matches = re.findall(r"file\(['\"]([^'\"]+)['\"]\)", s)
        if file_matches:
            for f in file_matches:
                p_name = Path(f).name
                if p_name and p_name not in itens:
                    itens.append(p_name)
            return
        art_matches = re.findall(r"artifact\(['\"]([^'\"]+)['\"]\)", s)
        if art_matches:
            for a in art_matches:
                p_name = Path(a).name
                if p_name and p_name not in itens:
                    itens.append(p_name)
            return
        if any(s.endswith(ext) for ext in ('.json', '.md', '.html', '.svg', '.tsx', '.css', '.go', '.py', '.txt', '.yml', '.yaml')):
            p_name = Path(s).name
            if p_name and p_name not in itens:
                itens.append(p_name)
            return
        if s.startswith("var.") or s in ("alvo", "target"):
            if s not in itens:
                itens.append(s)
            return
        limpa = s.strip("'\"")
        limpa = re.sub(r"^(?:file|artifact)\(['\"]?([^'\"]+)['\"]?\)$", r"\1", limpa)
        if limpa and limpa not in itens:
            p_cand = Path(limpa).name if any(limpa.endswith(e) for e in ('.json', '.md', '.html', '.svg')) else limpa
            itens.append(p_cand)
    processar_item(valor)
    return ", ".join(itens) if itens else "—"

def carregar_descricoes_da_skill(slug: str) -> dict:
    descricoes = {}
    candidatos = [
        Path(f"skills/{slug}/SKILL.md"),
        Path(f"skills/legoloop-{slug}/SKILL.md"),
        Path(f"skills/{slug.replace('legoloop-', '')}/SKILL.md")
    ]
    for c in candidatos:
        if c.exists():
            texto = c.read_text(encoding="utf-8")
            for m in re.finditer(r"\*\*[^*]+?\(`([^`]+)`\):\*\*\s*([^\n\r]+)", texto):
                descricoes[m.group(1).strip()] = m.group(2).strip()
            for m in re.finditer(r"\*\*`([^`]+)`:\*\*\s*([^\n\r]+)", texto):
                descricoes[m.group(1).strip()] = m.group(2).strip()
            if descricoes:
                break
    return descricoes

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

def resolver_dag(alvo: str):
    p = Path(alvo)
    if p.is_file():
        if p.name == "manifest.json":
            slug = p.parent.name
            dag = Path("processos") / slug / "pipeline.dag.json"
            return slug, dag if dag.exists() else p
        if p.name == "pipeline.dag.json":
            return p.parent.name, p
        return p.stem, p
    if (Path("processos") / alvo / "pipeline.dag.json").exists():
        return Path(alvo).name, Path("processos") / alvo / "pipeline.dag.json"
    if (Path("processos") / f"legoloop-{alvo}" / "pipeline.dag.json").exists():
        slug = f"legoloop-{Path(alvo).name}"
        return slug, Path("processos") / slug / "pipeline.dag.json"
    return Path(alvo).name, Path("processos") / Path(alvo).name / "pipeline.dag.json"

def main():
    if len(sys.argv) < 2:
        alvo = "legoloop-diagram-design"
    else:
        alvo = sys.argv[1]

    slug, dag_path = resolver_dag(alvo)
    slug = Path(slug).name
    descs_skill = carregar_descricoes_da_skill(slug)

    nodes = []
    if dag_path.exists():
        try:
            with open(dag_path, "r", encoding="utf-8") as f:
                dag_data = json.load(f)
                nodes = dag_data.get("nodes", [])
        except Exception as e:
            print(f"[WARN] Erro ao ler DAG: {e}")

    quadrantes_nodes = []
    for i, n in enumerate(nodes):
        nid = n.get("id", f"e0{i+1}")
        ntype = n.get("type", "skill")
        nname = n.get("name", nid)
        estilo = obter_estilo_peca(ntype)
        cin_limpo = extrair_nomes_arquivos_ou_insumos(n.get("context_in") or n.get("inputs"))
        cout_limpo = extrair_nomes_arquivos_ou_insumos(n.get("context_out") or n.get("output", {}).get("artifact"))
        desc_no = descs_skill.get(nid) or n.get("description") or f"Executa a etapa {nname} com isolamento estéril."
        quadrantes_nodes.append({
            "id": nid,
            "name": nname,
            "description": desc_no,
            "type": ntype,
            "badge": estilo["badge"],
            "color": estilo["cor"],
            "layer": i + 1,
            "inputs_clean": cin_limpo,
            "output_clean": cout_limpo,
            "gate": n.get("output", {}).get("gate", "exit_code == 0"),
            "depends_on": n.get("depends_on", [])
        })

    out_file = Path("docs/diagrams") / slug / "quadrantes.json"
    out_file.parent.mkdir(parents=True, exist_ok=True)
    
    quadrantes_data = {
        "slug": slug,
        "nodes_count": len(quadrantes_nodes),
        "layers_count": len(quadrantes_nodes),
        "nodes": quadrantes_nodes
    }
    out_file.write_text(json.dumps(quadrantes_data, indent=2), encoding="utf-8")
    print(f"[OK] e02 quadrantes gerado com sucesso: {out_file}")

if __name__ == "__main__":
    main()
