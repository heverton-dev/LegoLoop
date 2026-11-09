import json, os, subprocess

base_dir = "docs/books/legoloop-cartoon-shorts"
livro_json_path = os.path.join(base_dir, "livro.json")

with open(livro_json_path, "r", encoding="utf-8") as f:
    manifest = json.load(f)

nome_base = manifest.get("nome_base", "livro")
partes = manifest.get("partes", [])
pasta_partes = os.path.join(base_dir, manifest.get("pasta_partes", "partes"))

conteudo_unificado = []
for p in partes:
    p_path = os.path.join(pasta_partes, p)
    if os.path.exists(p_path):
        with open(p_path, "r", encoding="utf-8") as pf:
            conteudo_unificado.append(pf.read())
            conteudo_unificado.append("\n\n")

md_final_path = os.path.join(base_dir, f"{nome_base}.md")
with open(md_final_path, "w", encoding="utf-8") as f:
    f.write("".join(conteudo_unificado))

pdf_final_path = os.path.join(base_dir, f"{nome_base}.pdf")
template_path = os.path.join(base_dir, manifest.get("template", "ativos/livro.typst"))

cmd = [
    "pandoc",
    md_final_path,
    "--template", template_path,
    "--pdf-engine", "typst",
    "-o", pdf_final_path
]

print("Compilando livro do cartoon-shorts com pandoc + typst...")
res = subprocess.run(cmd, capture_output=True, text=True)
if res.returncode == 0:
    print(f"PDF COMPILADO COM SUCESSO: {pdf_final_path} ({os.path.getsize(pdf_final_path)} bytes)")
else:
    print("ERRO:", res.stderr)
