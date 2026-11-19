#!/usr/bin/env python3
"""Inicializador Determinístico de Espaço de Diagrama LegoLoop (Nó e01).

Cria a árvore canônica de diretórios em docs/diagrams/<slug>/, copia ou
inicializa anticorpos e emite o manifest.json com metadados do processo.
"""

import sys
import json
from pathlib import Path
from datetime import datetime, timezone

def resolver_processo(alvo: str):
    p = Path(alvo)
    if p.is_file() and p.suffix == ".json":
        if p.name == "pipeline.dag.json":
            return p.parent.name, p
        return p.stem, p
    if p.is_dir() and (p / "pipeline.dag.json").exists():
        return p.name, p / "pipeline.dag.json"
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

    slug, dag_path = resolver_processo(alvo)
    slug = Path(slug).name

    base_dir = Path("docs/diagrams") / slug
    subdirs = ["src", "assets", "export", "animation", "anticorpos"]
    for s in subdirs:
        (base_dir / s).mkdir(parents=True, exist_ok=True)

    # Copiar ou inicializar anticorpos
    proc_anticorpos = dag_path.parent / "anticorpos.json" if dag_path.exists() else Path("anticorpos.json")
    anticorpos_out = base_dir / "anticorpos" / "anticorpos.json"
    if proc_anticorpos.exists():
        anticorpos_out.write_text(proc_anticorpos.read_text(encoding="utf-8"), encoding="utf-8")
    elif not anticorpos_out.exists():
        anticorpos_out.write_text(json.dumps({"anticorpos": []}, indent=2), encoding="utf-8")

    # Ler dados do DAG se existir
    nodes_count = 0
    nome_proc = slug
    versao = "1.0.0"
    if dag_path.exists():
        try:
            with open(dag_path, "r", encoding="utf-8") as f:
                dag_data = json.load(f)
                nodes_count = len(dag_data.get("nodes", []))
                nome_proc = dag_data.get("name", slug)
                versao = dag_data.get("version", "1.0.0")
        except Exception:
            pass

    manifest_data = {
        "slug": slug,
        "name": nome_proc,
        "version": versao,
        "nodes_count": nodes_count,
        "visual_type": "process",
        "semantic_pattern": "stage-framework",
        "generated_at": datetime.now(timezone.utc).isoformat()
    }
    manifest_out = base_dir / "manifest.json"
    manifest_out.write_text(json.dumps(manifest_data, indent=2), encoding="utf-8")
    print(f"[OK] e01 inicializado com sucesso: {base_dir}")

if __name__ == "__main__":
    main()
