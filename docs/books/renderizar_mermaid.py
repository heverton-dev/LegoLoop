"""Varre arquivos Markdown, detecta blocos ```mermaid e os renderiza como SVG
usando o mmdc (mermaid-cli) com o Chrome local configurado.
Substitui o bloco de código pela inclusão da imagem: ![Diagrama](img/diagrama_X.svg){width=95%}
"""
import hashlib
import os
import re
import subprocess
import sys

DIR_BOOKS = os.path.dirname(os.path.abspath(__file__))
PUPPETEER_CFG = os.path.join(DIR_BOOKS, "puppeteer-config.json")

def render_diagram(code: str, out_svg: str) -> bool:
    if os.path.exists(out_svg) and os.path.getsize(out_svg) > 100:
        return True  # Já renderizado
    os.makedirs(os.path.dirname(out_svg), exist_ok=True)
    temp_mmd = out_svg + ".mmd"
    with open(temp_mmd, "w", encoding="utf-8") as f:
        f.write(code.strip())
    
    cmd = [
        "mmdc.cmd",
        "-p", PUPPETEER_CFG,
        "-i", temp_mmd,
        "-o", out_svg,
        "-b", "transparent"
    ]
    try:
        res = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=30)
        if os.path.exists(temp_mmd):
            os.remove(temp_mmd)
        if res.returncode == 0 and os.path.exists(out_svg):
            print(f"  [OK] Diagrama renderizado: {out_svg}")
            return True
        else:
            print(f"  [AVISO] Falha ao renderizar {out_svg}: {res.stderr}")
            return False
    except Exception as e:
        print(f"  [ERRO] Exceção ao executar mmdc: {e}")
        return False

def process_file(filepath: str):
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    pattern = re.compile(r"```mermaid\s*\n(.*?)\n```", re.DOTALL)
    matches = list(pattern.finditer(content))
    if not matches:
        return

    book_dir = os.path.dirname(os.path.dirname(os.path.abspath(filepath)))
    img_dir = os.path.join(book_dir, "img")
    os.makedirs(img_dir, exist_ok=True)

    new_content = content
    # Processar de trás para frente para manter índices
    for match in reversed(matches):
        code = match.group(1)
        h = hashlib.sha256(code.encode("utf-8")).hexdigest()[:8]
        svg_filename = f"diagrama_{h}.svg"
        svg_path = os.path.join(img_dir, svg_filename)
        rel_svg = f"/img/{svg_filename}"

        if render_diagram(code, svg_path):
            replacement = f"\n![Diagrama de Arquitetura]({rel_svg}){{width=95%}}\n"
            start, end = match.span()
            new_content = new_content[:start] + replacement + new_content[end:]

    if new_content != content:
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(new_content)
        print(f"[OK] Arquivo atualizado com diagramas SVG: {filepath}")

if __name__ == "__main__":
    for p in sys.argv[1:]:
        process_file(p)
