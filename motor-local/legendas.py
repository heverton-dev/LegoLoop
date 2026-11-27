"""Legendas queimadas: blocos de N palavras sincronizados com a voz, como camadas de texto fixas na tela."""
import json
import os
import re

NUM = re.compile(r'\d|bilh|milh|mil\b|US\$', re.I)


def blocos(lines, n):
    """lines: [{'text','start','end'}] -> [(inicio, fim, texto)] com tempo proporcional ao tamanho das palavras."""
    out = []
    for l in lines:
        ws = l['text'].split()
        if not ws:
            continue
        t0 = float(l['start']); t1 = max(t0 + 0.5, float(l['end']) - 0.3)
        peso = [len(w) + 1 for w in ws]; tot = float(sum(peso)); acc = 0.0
        tw = []
        for p in peso:
            tw.append(t0 + (t1 - t0) * acc / tot); acc += p
        tw.append(t1)
        for i in range(0, len(ws), n):
            out.append((tw[i], tw[min(i + n, len(ws))], ' '.join(ws[i:i + n])))
    return out


def carregar(base):
    for f in (os.path.join(base, 'data', 'lines.json'),):
        if os.path.exists(f):
            return json.load(open(f, encoding='utf-8'))['lines']
    return None


def aplicar(shots, starts, lines, cfg, cores):
    """Adiciona camadas de legenda a cada corte de cena (cartoes pretos ficam sem legenda)."""
    n = int(cfg.get('palavras_por_bloco', 3)); y = cfg.get('posicao_y_pct', 68) / 100.0 * 1080
    cor = cfg.get('cor', '#FFFFFF'); dest = cores.get(cfg.get('destaque', 'foco'), cfg.get('destaque', '#FFD24A'))
    bl = blocos(lines, n)
    for sh, s0 in zip(shots, starts):
        if sh.get('type') == 'card':
            continue
        d = float(sh['duration'])
        for (a, b, txt) in bl:
            if b <= s0 or a >= s0 + d:
                continue
            sh.setdefault('layers', []).append({
                'kind': 'text', 'text': txt, 'x': 960, 'y': y, 'size': 84, 'font': 'sans',
                'color': dest if NUM.search(txt) else cor, 'stroke': '#000000', 'sw': 10,
                'start': round(max(0.0, a - s0), 3), 'end': round(min(d, b - s0), 3),
                'in': 'pop', 'inDur': 0.12, 'hud': True, 'legenda': True})
    return shots
