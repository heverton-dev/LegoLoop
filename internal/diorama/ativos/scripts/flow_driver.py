#!/usr/bin/env python3
"""Driver autoral LegoLoop para Google Flow via Chrome CDP (porta 9333)."""
import argparse
import base64
import json
import os
import re
import sys
import time
from io import BytesIO


def anexar_ref(pg, termo):
    pg.keyboard.press("Escape")
    pg.locator("button:has-text('add')").last.click()
    busca = pg.get_by_placeholder(re.compile("Pesquisar|Search", re.I))
    busca.wait_for(state="visible", timeout=8000)
    busca.fill(termo)
    pg.wait_for_timeout(1200)
    pg.locator("span.asset-title >> text=/" + re.escape(termo) + "/i >> visible=true").first.click(timeout=8000)
    pg.wait_for_timeout(800)
    inc = pg.locator("text=/Incluir no comando|Add to prompt/i >> visible=true")
    if inc.count():
        inc.first.click(timeout=5000)
    pg.wait_for_timeout(1000)


def criar_projeto(p, cdp_url):
    browser = p.chromium.connect_over_cdp(cdp_url)
    ctx = browser.contexts[0]
    pg = next((x for x in ctx.pages if "flow.google.com" in x.url), None) or ctx.new_page()
    pg.goto("https://flow.google.com/", wait_until="domcontentloaded")
    pg.wait_for_timeout(2500)
    if "/about" in pg.url or "accounts.google.com" in pg.url or "signin" in pg.url:
        sys.exit("Chrome do Flow nao esta logado: entre na conta Google antes de criar projeto")
    pg.locator("text=/Novo projeto|New project/i >> visible=true").first.click(timeout=15000)
    pid = ""
    for _ in range(25):
        pg.wait_for_timeout(1000)
        m = re.search(r"/project/([0-9a-f-]{36})", pg.url)
        if m:
            pid = m.group(1)
            break
    if not pid:
        sys.exit("nao foi possivel criar ou capturar o UUID do projeto no Google Flow")
    print(json.dumps({"ok": True, "flow_projeto": pid}))


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--acao", choices=["gerar", "criar"], default="gerar")
    ap.add_argument("--cdp", default="http://127.0.0.1:9222")
    ap.add_argument("--flow-projeto", default="")
    ap.add_argument("--prompt", default="")
    ap.add_argument("--saida-arquivo", default="")
    ap.add_argument("--ref", action="append", default=[])
    ap.add_argument("--espera", type=int, default=150)
    ap.add_argument("--baixar-todas", action="store_true", default=True)
    args = ap.parse_args()

    from PIL import Image
    from playwright.sync_api import sync_playwright

    with sync_playwright() as p:
        if args.acao == "criar":
            criar_projeto(p, args.cdp)
            return
        if not args.flow_projeto or not args.prompt or not args.saida_arquivo:
            sys.exit("--flow-projeto, --prompt e --saida-arquivo sao obrigatorios para gerar")
        browser = p.chromium.connect_over_cdp(args.cdp)
        ctx = browser.contexts[0]
        pg = next((x for x in ctx.pages if args.flow_projeto in x.url), None)
        if not pg:
            pg = ctx.new_page()
            pg.goto(f"https://flow.google.com/project/{args.flow_projeto}", wait_until="domcontentloaded")
        pg.wait_for_timeout(2000)
        if args.flow_projeto not in pg.url:
            sys.exit(f"trava flow_projeto violada: {pg.url}")
        caixa = pg.locator("div.ProseMirror, [contenteditable='true']").first
        if not caixa.is_visible():
            sys.exit("Chrome do Flow sem caixa de prompt visivel (faca login na conta Google)")

        antes = set(pg.evaluate("[...document.querySelectorAll('img')].filter(i=>i.naturalWidth>=1000&&i.src.includes('/asb/')).map(i=>i.src)"))
        for r in args.ref:
            anexar_ref(pg, r)

        caixa.click()
        pg.keyboard.press("Control+A")
        pg.keyboard.press("Backspace")
        pg.keyboard.insert_text(args.prompt)
        pg.wait_for_timeout(600)
        pg.keyboard.press("Enter")

        deadline = time.time() + args.espera
        novas = []
        quietos = 0
        while time.time() < deadline:
            atuais = [s for s in pg.evaluate("[...document.querySelectorAll('img')].filter(i=>i.naturalWidth>=1000&&i.src.includes('/asb/')).map(i=>i.src)") if s not in antes]
            if atuais and len(atuais) == len(novas):
                quietos += 1
            else:
                quietos = 0
            novas = atuais
            if novas and (not args.baixar_todas or quietos >= 4 or len(novas) >= 2):
                break
            time.sleep(2)
        if not novas:
            sys.exit("timeout aguardando imagem gerada no Google Flow")

        os.makedirs(os.path.dirname(os.path.abspath(args.saida_arquivo)), exist_ok=True)
        base_no_ext, ext = os.path.splitext(args.saida_arquivo)
        arquivos = []
        larg, alt = 0, 0
        alvos = novas if args.baixar_todas else [novas[-1]]
        for idx, src in enumerate(alvos, start=1):
            dados = pg.evaluate(
                "async s=>{const b=await (await fetch(s)).blob();const r=new FileReader();return await new Promise(ok=>{r.onload=()=>ok(r.result);r.readAsDataURL(b)})}",
                src,
            )
            destino = args.saida_arquivo if idx == 1 else f"{base_no_ext}_var{idx}{ext}"
            im = Image.open(BytesIO(base64.b64decode(dados.split(",", 1)[1]))).convert("RGBA")
            im.save(destino, format="PNG")
            if idx == 1:
                larg, alt = im.size
            arquivos.append(destino)
        print(json.dumps({"ok": True, "arquivo": args.saida_arquivo, "arquivos": arquivos, "largura": larg, "altura": alt}))


if __name__ == "__main__":
    main()
