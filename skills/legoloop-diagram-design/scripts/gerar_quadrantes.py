#!/usr/bin/env python3
"""Ingestão e Estruturação Determinística de Quadrantes (Nó e02).

Lê o pipeline.dag.json ou manifest.json do processo alvo e emite
docs/diagrams/<slug>/quadrantes.json com categorização taxonômica e geométrica.
"""

import sys
import json
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
        quadrantes_nodes.append({
            "id": nid,
            "name": nname,
            "type": ntype,
            "badge": estilo["badge"],
            "color": estilo["cor"],
            "layer": i + 1,
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
