"""Reescreve tabelas pipe: alinha colunas e faz a largura (traços) proporcional ao maior conteúdo.
Pandoc usa os traços do separador como largura relativa; isso evita quebra de linha nas células."""
import re, sys, unicodedata
def w(s): return sum(2 if unicodedata.east_asian_width(c) in 'WF' else 1 for c in s)
def cells(l): return [c.strip() for c in l.strip().strip('|').split('|')]
def fmt(rows, aligns):
    n=len(aligns); ws=[max(3, *(w(r[i]) for r in rows)) for i in range(n)]
    out=[]
    for k,r in enumerate(rows):
        out.append('| '+' | '.join(r[i]+' '*(ws[i]-w(r[i])) for i in range(n))+' |')
        if k==0:
            out.append('| '+' | '.join((':' if a in 'lc' else '-')+'-'*(ws[i]-2)+(':' if a in 'rc' else '-') for i,a in enumerate(aligns))+' |')
    return out
for path in sys.argv[1:]:
    L=open(path,encoding='utf-8').read().split('\n'); o=[]; i=0
    while i<len(L):
        if L[i].lstrip().startswith('|') and i+1<len(L) and re.match(r'^\s*\|?\s*:?-{2,}', L[i+1]):
            hdr=cells(L[i]); seps=cells(L[i+1])
            al=['c' if s.startswith(':') and s.endswith(':') else 'r' if s.endswith(':') else 'l' for s in seps]
            rows=[hdr]; i+=2
            while i<len(L) and L[i].lstrip().startswith('|'): rows.append(cells(L[i])); i+=1
            o+=fmt(rows,al)
        else: o.append(L[i]); i+=1
    open(path,'w',encoding='utf-8').write('\n'.join(o))
