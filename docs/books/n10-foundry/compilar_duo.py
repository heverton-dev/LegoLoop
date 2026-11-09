import json
import os
import shutil
import subprocess
import sys

REPO_ROOT = os.path.abspath(".")
N10_DIR = os.path.join(REPO_ROOT, "docs", "books", "n10-foundry")
GLOBAL_ATIVOS = os.path.join(REPO_ROOT, "docs", "books", "ativos")

def preparar_ativos(target_dir):
    ativos_dest = os.path.join(target_dir, "ativos")
    os.makedirs(ativos_dest, exist_ok=True)
    if os.path.exists(GLOBAL_ATIVOS):
        for item in os.listdir(GLOBAL_ATIVOS):
            s = os.path.join(GLOBAL_ATIVOS, item)
            d = os.path.join(ativos_dest, item)
            if os.path.isfile(s):
                shutil.copy2(s, d)
    modelo_typst = os.path.join(REPO_ROOT, "docs", "books", "_modelo_livro.typst")
    if os.path.exists(modelo_typst):
        shutil.copy2(modelo_typst, os.path.join(ativos_dest, "livro.typst"))

def compilar_edicao(subpasta, nome_edicao):
    edicao_dir = os.path.join(N10_DIR, subpasta)
    preparar_ativos(edicao_dir)

    livro_json_path = os.path.join(edicao_dir, "livro.json")
    if not os.path.exists(livro_json_path):
        print(f"[ERRO] {livro_json_path} não encontrado.")
        return False

    with open(livro_json_path, "r", encoding="utf-8") as f:
        manifest = json.load(f)

    nome_base = manifest.get("nome_base", subpasta)
    partes = manifest.get("partes", [])
    pasta_partes = os.path.join(edicao_dir, manifest.get("pasta_partes", "partes"))

    conteudo_unificado = []
    for p in partes:
        p_path = os.path.join(pasta_partes, p)
        if os.path.exists(p_path):
            with open(p_path, "r", encoding="utf-8") as pf:
                conteudo_unificado.append(pf.read())
                conteudo_unificado.append("\n\n")

    md_final_name = f"{nome_base}.md"
    md_final_path = os.path.join(edicao_dir, md_final_name)
    with open(md_final_path, "w", encoding="utf-8") as f:
        f.write("".join(conteudo_unificado))

    # Também espelha o markdown na raiz de n10-foundry
    shutil.copy2(md_final_path, os.path.join(N10_DIR, md_final_name))

    pdf_final_name = f"{nome_base}.pdf"
    pdf_in_edicao = os.path.join(edicao_dir, pdf_final_name)
    template_path = os.path.join("ativos", "livro.typst")

    cmd = [
        "pandoc",
        md_final_name,
        "--template", template_path,
        "--pdf-engine", "typst",
        "--pdf-engine-opt=--root=.",
        "-o", pdf_final_name
    ]

    print(f"\n[LegoPress Duo] Compilando {nome_edicao} com Pandoc + Typst...")
    res = subprocess.run(cmd, cwd=edicao_dir, capture_output=True, text=True)
    if res.returncode == 0:
        pdf_dest_root = os.path.join(N10_DIR, pdf_final_name)
        shutil.copy2(pdf_in_edicao, pdf_dest_root)
        size = os.path.getsize(pdf_dest_root)
        print(f"✓ PDF GERADO: {pdf_dest_root} ({size:,} bytes)")
        return True
    else:
        print(f"✗ ERRO ao compilar {nome_edicao}:\n{res.stderr}")
        return False

def main():
    print("=" * 70)
    print("  LEGOPRESS DUO // COMPILAÇÃO DUPLA (ENGENHARIA vs EXECUTIVO)")
    print("=" * 70)

    ok_eng = compilar_edicao("engenharia", "Edição de Engenharia (Cockpit PhD)")
    ok_exec = compilar_edicao("executivo", "Edição Executiva (Cabine Leigo)")

    if ok_eng and ok_exec:
        print("\n" + "=" * 70)
        print("✓ AMBOS OS LIVROS DIDÁTICOS FORAM COMPILADOS COM SUCESSO!")
        print(f"  [1] Engenharia: docs/books/n10-foundry/2026-11-09_N10-FOUNDRY-ENGENHARIA.pdf")
        print(f"  [2] Executivo:  docs/books/n10-foundry/2026-11-09_N10-FOUNDRY-EXECUTIVO.pdf")
        print("=" * 70)
        sys.exit(0)
    else:
        print("\n[FALHA] Compilação incompleta.")
        sys.exit(1)

if __name__ == "__main__":
    main()
