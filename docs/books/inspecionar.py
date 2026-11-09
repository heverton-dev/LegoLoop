"""Relatório de geometria: linhas por página e células de tabela com mais de uma linha."""
import sys, fitz
d=fitz.open(sys.argv[1]); out=sys.argv[2] if len(sys.argv)>2 else None
for i,p in enumerate(d):
    n=len([l for l in p.get_text().splitlines() if l.strip()])
    # spans com fonte <=6.3pt = corpo de tabela; agrupa por bloco e conta linhas distintas por bloco
    quebras=0
    for b in p.get_text('dict')['blocks']:
        ls=[l for l in b.get('lines',[]) if l['spans'] and max(s['size'] for s in l['spans'])<=6.4 and l['bbox'][1]>60]
        ys=sorted({round(l['bbox'][1]) for l in ls})
        if len(ys)>1: quebras+=1
    flag=(' CURTA' if n<15 and i>1 else '')+(f' TAB-QUEBRA:{quebras}' if quebras else '')
    print(f'{i+1:3d} {n:4d}{flag}')
    if out: p.get_pixmap(dpi=70).save(f'{out}/p{i+1:02d}.png')
