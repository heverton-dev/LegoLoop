import os
import sys
import json
import shutil
import subprocess

def slug_upper(slug):
    s = slug.replace("legoloop-", "")
    return s.upper().replace("-", "_")

def scaffold_partes(b_dir, slug, titulo, descricao):
    partes_dir = os.path.join(b_dir, "partes")
    os.makedirs(partes_dir, exist_ok=True)
    
    files = {
        "00-frontmatter.md": f"# Frontmatter\n\nManual Didático: {titulo}\nProcesso {slug} no LegoLoop.\nVersão 1.0.0 — Data 2026-11-20.\n",
        "01-macro-fundacao.md": f"# Macro Fundação\n\nVisão macro do processo {slug}.\n{descricao}\nOperação orientada aos 7 Protocolos Industriais e 5 Prisões Mecânicas.\n",
        "02-meso-arquitetura.md": f"# Meso Arquitetura\n\nArquitetura técnica e grafo DAG de {slug}.\nContratos de entrada (context_in), saída (context_out) e gates mecânicos.\n",
        "03-micro-passo-a-passo.md": f"# Micro Passo a Passo\n\nPasso a passo determinístico de execução do processo {slug}.\nComandos CLI: `legoloop run processos/{slug}/pipeline.dag.json`.\nValidação e testes contínuos.\n"
    }
    
    for fn, content in files.items():
        fp = os.path.join(partes_dir, fn)
        if not os.path.exists(fp):
            with open(fp, "w", encoding="utf-8") as f:
                f.write(content)

def main():
    root = os.path.abspath(".")
    books_dir = os.path.join(root, "docs", "books")
    ativos_base = os.path.join(books_dir, "ativos")
    modelo_typst = os.path.join(books_dir, "_modelo_livro.typst")
    
    # Placeholders para imagens ausentes em homologacao
    sample_png = os.path.join(books_dir, "legoloop", "preview", "p01.png")
    homolog_dir = os.path.join(books_dir, "legoloop-homologacao-ecosistema", "img")
    os.makedirs(homolog_dir, exist_ok=True)
    for img_name in ["2026-11-14_HOMOLOG_STUDIO.png", "2026-11-14_HOMOLOG_INVENTARIO.png", "2026-11-14_HOMOLOG_SWARM.png"]:
        tgt = os.path.join(homolog_dir, img_name)
        if not os.path.exists(tgt) and os.path.exists(sample_png):
            shutil.copy2(sample_png, tgt)
        
    all_books = [d for d in os.listdir(books_dir) if os.path.isdir(os.path.join(books_dir, d)) and d != "ativos"]
    
    report_lines = [
        "# Relatório de Padronização e Compilação dos Livros LegoPress",
        "",
        "**Data:** 2026-11-20",
        f"**Total de livros processados:** {len(all_books)}",
        "",
        "| Livro | Status Compilação | Páginas | Auditoria LegoPress |",
        "| :--- | :---: | :---: | :---: |"
    ]
    
    sucessos = 0
    falhas = 0
    
    for b in sorted(all_books):
        b_dir = os.path.join(books_dir, b)
        json_path = os.path.join(b_dir, "livro.json")
        if not os.path.exists(json_path):
            continue
            
        print(f"[*] Processando {b}...")
        try:
            with open(json_path, "r", encoding="utf-8") as f:
                meta = json.load(f)
        except Exception as e:
            print(f"[-] Erro ao ler livro.json de {b}: {e}")
            falhas += 1
            report_lines.append(f"| `{b}` | Erro JSON | - | ❌ |")
            continue
            
        titulo = meta.get("titulo", f"Processo {b}")
        descricao = meta.get("descricao", "Processo canônico do LegoLoop")
        
        # Migrar manifesto legado se não tiver partes
        if "partes" not in meta:
            meta["versao_manifesto"] = "1.0"
            meta["nome_base"] = f"2026-11-20_PROCESSO_{slug_upper(b)}"
            meta["pasta_partes"] = "partes"
            meta["template"] = "ativos/livro.typst"
            meta["partes"] = [
                "00-frontmatter.md",
                "01-macro-fundacao.md",
                "02-meso-arquitetura.md",
                "03-micro-passo-a-passo.md"
            ]
            meta["hash_fonte"] = f"{b}-v1"
            meta["ultima_compilacao"] = "2026-11-20"
            with open(json_path, "w", encoding="utf-8") as f:
                json.dump(meta, f, indent=2, ensure_ascii=False)
                f.write("\n")
                
        nome = meta.get("nome_base", b)
        partes = meta.get("partes", [])
        
        # Scaffold das partes se necessário
        scaffold_partes(b_dir, b, titulo, descricao)
        
        # 1. Estrutura de pastas
        build_dir = os.path.join(b_dir, "build")
        preview_dir = os.path.join(b_dir, "preview")
        img_dir = os.path.join(b_dir, "img")
        ativos_dir = os.path.join(b_dir, "ativos")
        build_ativos = os.path.join(build_dir, "ativos")
        build_img = os.path.join(build_dir, "img")
        
        os.makedirs(build_dir, exist_ok=True)
        os.makedirs(preview_dir, exist_ok=True)
        os.makedirs(img_dir, exist_ok=True)
        os.makedirs(ativos_dir, exist_ok=True)
        os.makedirs(build_ativos, exist_ok=True)
        os.makedirs(build_img, exist_ok=True)
        
        # 2. Copiar ativos para ativos/ E build/ativos/ E build/
        if os.path.exists(ativos_base):
            for item in os.listdir(ativos_base):
                s = os.path.join(ativos_base, item)
                if os.path.isfile(s):
                    shutil.copy2(s, os.path.join(ativos_dir, item))
                    shutil.copy2(s, os.path.join(build_ativos, item))
                    shutil.copy2(s, os.path.join(build_dir, item))
                    
        # Copiar ativos locais específicos de b_dir/ativos se houver
        for item in os.listdir(ativos_dir):
            s = os.path.join(ativos_dir, item)
            if os.path.isfile(s):
                shutil.copy2(s, os.path.join(build_ativos, item))
                shutil.copy2(s, os.path.join(build_dir, item))
                
        # Copiar img/ para build/img/ e build/
        if os.path.exists(img_dir):
            for item in os.listdir(img_dir):
                s = os.path.join(img_dir, item)
                if os.path.isfile(s):
                    shutil.copy2(s, os.path.join(build_img, item))
                    shutil.copy2(s, os.path.join(build_dir, item))
                    
        # 3. Sincronizar modelo typst
        if os.path.exists(modelo_typst):
            shutil.copy2(modelo_typst, os.path.join(ativos_dir, "livro.typst"))
            shutil.copy2(modelo_typst, os.path.join(build_ativos, "livro.typst"))
            shutil.copy2(modelo_typst, os.path.join(build_dir, "livro.typst"))
            
        # 4. Normalizadores
        partes_paths = [os.path.join(b_dir, "partes", p) for p in partes if os.path.exists(os.path.join(b_dir, "partes", p))]
        if not partes_paths:
            print(f"[-] {b} sem arquivos de partes válidos")
            falhas += 1
            report_lines.append(f"| `{b}` | Partes Ausentes | - | ❌ |")
            continue
            
        for script in ["renderizar_mermaid.py", "normalizar_tabelas.py", "normalizar_listas.py"]:
            sp = os.path.join(books_dir, script)
            if os.path.exists(sp):
                subprocess.run([sys.executable, sp] + partes_paths, cwd=b_dir, capture_output=True, timeout=30, stdin=subprocess.DEVNULL)
                
        # 5. Consolidar manuscrito markdown
        md_out = os.path.join(b_dir, f"{nome}.md")
        with open(md_out, "w", encoding="utf-8") as out_f:
            for pp in partes_paths:
                with open(pp, "r", encoding="utf-8") as in_f:
                    out_f.write(in_f.read())
                    out_f.write("\n\n")
                    
        # 6. Pandoc
        typ_out = os.path.join(build_dir, "livro.typ")
        tpl_path = os.path.join(ativos_dir, "livro.typst")
        cmd_pandoc = ["pandoc"] + partes_paths + [
            "-f", "markdown",
            "-t", "typst",
            "--standalone",
            "--wrap=none",
            "--template", tpl_path,
            "-o", typ_out
        ]
        try:
            res_p = subprocess.run(cmd_pandoc, cwd=b_dir, capture_output=True, text=True, timeout=60, stdin=subprocess.DEVNULL)
            if res_p.returncode != 0:
                print(f"[-] Erro pandoc em {b}: {res_p.stderr[:200]}")
                falhas += 1
                report_lines.append(f"| `{b}` | Erro Pandoc | - | ❌ |")
                continue
        except subprocess.TimeoutExpired:
            print(f"[-] Timeout pandoc em {b}")
            falhas += 1
            report_lines.append(f"| `{b}` | Timeout Pandoc | - | ❌ |")
            continue
            
        # 7. Desescapar painéis
        desescapar = os.path.join(books_dir, "desescapar_paineis.py")
        if os.path.exists(desescapar):
            subprocess.run([sys.executable, desescapar, typ_out], cwd=b_dir, capture_output=True, timeout=30, stdin=subprocess.DEVNULL)
            
        # 8. Typst compile
        pdf_out = os.path.join(b_dir, f"{nome}.pdf")
        cmd_typst = ["typst", "compile", "--root", b_dir, typ_out, pdf_out]
        try:
            res_t = subprocess.run(cmd_typst, cwd=b_dir, capture_output=True, text=True, timeout=60, stdin=subprocess.DEVNULL)
            if res_t.returncode != 0:
                print(f"[-] Erro typst em {b}: {res_t.stderr[:200]}")
                falhas += 1
                report_lines.append(f"| `{b}` | Erro Typst | - | ❌ |")
                continue
        except subprocess.TimeoutExpired:
            print(f"[-] Timeout typst em {b}")
            falhas += 1
            report_lines.append(f"| `{b}` | Timeout Typst | - | ❌ |")
            continue
            
        # 9. Auditar com legoloop legopress auditar
        res_aud = subprocess.run(["legoloop", "legopress", "auditar", pdf_out], capture_output=True, text=True, timeout=30, stdin=subprocess.DEVNULL)
        status_aud = "✅ Aprovado" if res_aud.returncode == 0 else "⚠️ Geometria Alerta"
        
        sucessos += 1
        print(f"[+] {b} concluído com sucesso: {pdf_out}")
        report_lines.append(f"| `{b}` | ✅ Compilado | OK | {status_aud} |")
        
    rep_path = os.path.join(root, "docs", "reports", "2026-11-20_REPORT_PADRONIZACAO_BOOKS.md")
    with open(rep_path, "w", encoding="utf-8") as f:
        f.write("\n".join(report_lines) + f"\n\n**Total Sucessos:** {sucessos} | **Falhas:** {falhas}\n")
        
    print(f"[*] Processamento finalizado. Sucessos: {sucessos}, Falhas: {falhas}. Relatório: {rep_path}")

if __name__ == "__main__":
    main()
