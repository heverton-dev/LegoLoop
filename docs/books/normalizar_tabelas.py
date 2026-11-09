"""Reescreve tabelas pipe: alinha colunas e distribui a largura relativa (traços) de forma equilibrada.
Garante pisos mínimos estritos para evitar que colunas de texto curto (ex: Nível, Parte, Autor)
colapsem e sobreponham texto no Typst compilado pelo Pandoc.
"""
import re
import sys
import unicodedata

def clean_txt(s):
    # Remove formatações markdown e links para medição real do texto
    t = re.sub(r'\[([^\]]+)\]\([^)]+\)', r'\1', s)
    t = re.sub(r'[*_`]', '', t)
    return t.strip()

def w(s):
    t = clean_txt(s)
    return sum(2 if unicodedata.east_asian_width(c) in 'WF' else 1 for c in t)

def cells(l):
    return [c.strip() for c in l.strip().strip('|').split('|')]

def is_numeric_col(rows, i):
    # Verifica se a coluna inteira (exceto cabeçalho) é numérica ou de índice curto
    for r in rows[1:]:
        if i >= len(r):
            continue
        val = clean_txt(r[i])
        if not (val.isdigit() or val in ('-', '#', 'Nº', 'No', '', 'OK', 'KO')):
            return False
    return True

def fmt(rows, aligns):
    n = len(aligns)
    raw_ws = []
    max_words = []
    for i in range(n):
        col_cells = [r[i] for r in rows if i < len(r)]
        raw_w = max([w(c) for c in col_cells] + [3])
        words = []
        for c in col_cells:
            words.extend(clean_txt(c).split())
        max_word = max([w(wd) for wd in words] + [3])
        raw_ws.append(raw_w)
        max_words.append(max_word)
    
    # Calcular pesos equilibrados com piso mínimo estrito
    weights = []
    for i in range(n):
        mw = max_words[i]
        rw = raw_ws[i]
        if is_numeric_col(rows, i) and mw <= 4:
            piso = max(6, mw + 2)
            peso = min(piso, rw)
        else:
            # Colunas textuais: piso mínimo de 14 traços ou palavra mais longa + 4
            piso = max(14, mw + 4)
            if rw > 40:
                # Compressão suave para textos longos (amortecedor não-linear)
                peso = 35 + int((rw - 40) ** 0.65)
            else:
                peso = max(piso, rw)
            peso = max(piso, peso)
        weights.append(peso)
    
    total_w = sum(weights)
    target_total = max(80, min(120, total_w))
    scaled_ws = [max(weights[i], int(weights[i] * target_total / total_w)) for i in range(n)]
    
    out = []
    for k, r in enumerate(rows):
        padded_cells = []
        for i in range(n):
            val = r[i] if i < len(r) else ''
            pad = max(0, scaled_ws[i] - w(val))
            padded_cells.append(val + ' ' * pad)
        out.append('| ' + ' | '.join(padded_cells) + ' |')
        if k == 0:
            sep_cells = []
            for i, a in enumerate(aligns):
                sw = max(scaled_ws[i], 5)
                left = ':' if a in 'lc' else '-'
                right = ':' if a in 'rc' else '-'
                middle = '-' * max(1, sw - 2)
                sep_cells.append(left + middle + right)
            out.append('| ' + ' | '.join(sep_cells) + ' |')
    return out

def process_file(path):
    with open(path, 'r', encoding='utf-8') as f:
        L = f.read().split('\n')
    o = []
    i = 0
    changed = False
    while i < len(L):
        if L[i].lstrip().startswith('|') and i + 1 < len(L) and re.match(r'^\s*\|?\s*:?-{2,}', L[i + 1]):
            hdr = cells(L[i])
            seps = cells(L[i + 1])
            al = ['c' if s.startswith(':') and s.endswith(':') else 'r' if s.endswith(':') else 'l' for s in seps]
            rows = [hdr]
            i += 2
            while i < len(L) and L[i].lstrip().startswith('|'):
                rows.append(cells(L[i]))
                i += 1
            formatted = fmt(rows, al)
            o += formatted
            changed = True
        else:
            o.append(L[i])
            i += 1
    if changed:
        with open(path, 'w', encoding='utf-8') as f:
            f.write('\n'.join(o))
        print(f"[OK] Tabelas normalizadas: {path}")

if __name__ == '__main__':
    for p in sys.argv[1:]:
        process_file(p)
