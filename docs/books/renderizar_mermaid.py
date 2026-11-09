"""Varre arquivos Markdown, detecta blocos ```mermaid e os renderiza como PNG de alta resolução (300 DPI)
usando o mmdc (mermaid-cli) com o Chrome local configurado.
Também detecta referências legadas a /img/diagrama_*.svg geradas com foreignObject e as converte
para PNG transparente de alta definição via Chromium/Playwright, eliminando o problema de caixas vazias no Typst.
Substitui as referências nos Markdown por: ![Diagrama de Arquitetura](/img/diagrama_X.png){width=85%}
"""
import asyncio
import hashlib
import os
import re
import subprocess
import sys

DIR_BOOKS = os.path.dirname(os.path.abspath(__file__))
PUPPETEER_CFG = os.path.join(DIR_BOOKS, "puppeteer-config.json")
CHROME_BIN = r"C:\Program Files\Google\Chrome\Application\chrome.exe"

def render_diagram_png(code: str, out_png: str) -> bool:
    if os.path.exists(out_png) and os.path.getsize(out_png) > 500:
        return True  # Já renderizado
    os.makedirs(os.path.dirname(out_png), exist_ok=True)
    temp_mmd = out_png + ".mmd"
    with open(temp_mmd, "w", encoding="utf-8") as f:
        f.write(code.strip())
    
    cmd = [
        "mmdc.cmd",
        "-p", PUPPETEER_CFG,
        "-i", temp_mmd,
        "-o", out_png,
        "-b", "transparent",
        "-s", "3"
    ]
    try:
        res = subprocess.run(cmd, shell=True, capture_output=True, text=True, timeout=45)
        if os.path.exists(temp_mmd):
            os.remove(temp_mmd)
        if res.returncode == 0 and os.path.exists(out_png):
            print(f"  [OK] Diagrama PNG renderizado: {out_png}")
            return True
        else:
            print(f"  [AVISO] Falha ao renderizar {out_png}: {res.stderr}")
            return False
    except Exception as e:
        print(f"  [ERRO] Exceção ao executar mmdc: {e}")
        return False

async def _convert_svg_to_png_async(svg_path: str, png_path: str, scale=3):
    from playwright.async_api import async_playwright
    svg_content = open(svg_path, encoding='utf-8').read()
    vb = re.search(r'viewBox="([^"]+)"', svg_content)
    w_val, h_val = 1100, 700
    if vb:
        parts = vb.group(1).split()
        if len(parts) == 4:
            w_val = float(parts[2])
            h_val = float(parts[3])

    async with async_playwright() as p:
        browser = await p.chromium.launch(
            executable_path=CHROME_BIN,
            args=["--no-sandbox", "--disable-gpu"]
        )
        page = await browser.new_page(device_scale_factor=scale)
        html_content = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  * {{ margin: 0; padding: 0; }}
  body {{ background: transparent; display: inline-block; }}
  svg {{ display: block; width: {w_val}px !important; height: {h_val}px !important; }}
</style>
</head>
<body>
  {svg_content}
</body>
</html>"""
        await page.set_content(html_content)
        await page.wait_for_timeout(400)
        svg_elem = await page.wait_for_selector("svg")
        await svg_elem.screenshot(path=png_path, omit_background=True)
        await browser.close()
    return os.path.exists(png_path) and os.path.getsize(png_path) > 500

def convert_svg_to_png(svg_path: str, png_path: str) -> bool:
    if os.path.exists(png_path) and os.path.getsize(png_path) > 500:
        return True
    try:
        return asyncio.run(_convert_svg_to_png_async(svg_path, png_path))
    except Exception as e:
        print(f"  [ERRO] Falha ao converter SVG {svg_path} para PNG: {e}")
        return False

def process_file(filepath: str):
    if not os.path.exists(filepath):
        return
    with open(filepath, "r", encoding="utf-8") as f:
        content = f.read()

    book_dir = os.path.dirname(os.path.dirname(os.path.abspath(filepath)))
    img_dir = os.path.join(book_dir, "img")
    os.makedirs(img_dir, exist_ok=True)
    new_content = content

    # 1. Tratar blocos ```mermaid
    pattern_mmd = re.compile(r"```mermaid\s*\n(.*?)\n```", re.DOTALL)
    matches_mmd = list(pattern_mmd.finditer(new_content))
    if matches_mmd:
        for match in reversed(matches_mmd):
            code = match.group(1)
            h = hashlib.sha256(code.encode("utf-8")).hexdigest()[:8]
            png_filename = f"diagrama_{h}.png"
            png_path = os.path.join(img_dir, png_filename)
            rel_png = f"/img/{png_filename}"

            if render_diagram_png(code, png_path):
                replacement = f"\n![Diagrama de Arquitetura]({rel_png}){{width=85%}}\n"
                start, end = match.span()
                new_content = new_content[:start] + replacement + new_content[end:]

    # 2. Tratar referências legadas a diagrama_*.svg que possam ter foreignObject
    pattern_svg = re.compile(r"!\[(.*?)\]\((/?img/diagrama_[a-zA-Z0-9_-]+)\.svg\)(\{[^}]*\})?")
    for match in pattern_svg.finditer(new_content):
        alt_txt = match.group(1)
        rel_base = match.group(2).lstrip("/")
        svg_filename = os.path.basename(rel_base) + ".svg"
        svg_path = os.path.join(book_dir, "img", svg_filename)
        
        # Só converte se o SVG existir e contiver foreignObject (diagramas mermaid sem texto no Typst)
        if os.path.exists(svg_path):
            with open(svg_path, "r", encoding="utf-8") as sf:
                svg_data = sf.read()
            if "foreignObject" in svg_data:
                png_filename = os.path.basename(rel_base) + ".png"
                png_path = os.path.join(book_dir, "img", png_filename)
                if convert_svg_to_png(svg_path, png_path):
                    old_str = match.group(0)
                    new_str = f"![{alt_txt}](/img/{png_filename}){{width=85%}}"
                    new_content = new_content.replace(old_str, new_str)
                    print(f"  [CONVERTIDO] SVG foreignObject -> PNG: {png_filename}")

    if new_content != content:
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(new_content)
        print(f"[OK] Arquivo atualizado com diagramas PNG: {filepath}")

if __name__ == "__main__":
    for p in sys.argv[1:]:
        process_file(p)
