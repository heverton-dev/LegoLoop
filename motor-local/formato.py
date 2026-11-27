"""Reenquadra um corte desenhado em 1920x1080 para 9:16 ou 1:1.

O mundo continua 1920x1080; a saida mostra uma janela dele (recorte), centrada no personagem principal.
Textos e objetos fixos na tela (hud) sao reposicionados dentro da janela, maiores e dentro da area segura.
"""
import copy

# formato: (largura e altura de saida cheia, largura do recorte no mundo, fator texto, fator objeto, area segura y)
FORMATOS = {
    '16:9': {'out': (1920, 1080), 'draft': (1280, 720), 'cw': 1920.0, 'kt': 1.0, 'ko': 1.0, 'safe': (0.0, 1.0), 'wrap': None},
    '9:16': {'out': (1080, 1920), 'draft': (720, 1280), 'cw': 1080 * 1080 / 1920, 'kt': 1.9, 'ko': 1.5, 'safe': (0.15, 0.78), 'wrap': (13, 64)},
    '1:1': {'out': (1080, 1080), 'draft': (720, 720), 'cw': 1080.0, 'kt': 1.3, 'ko': 1.15, 'safe': (0.08, 0.90), 'wrap': (20, 84)},
}
NAO_MOVER = {'pin'}  # posicoes geograficas no mapa


def janela(fmt):
    f = FORMATOS[fmt]
    return 960 - f['cw'] / 2, f['cw']


def _foco(sh):
    xs = []
    for c in sh.get('chars', []):
        p = c.get('path')
        xs.append(p[-1]['x'] if p else c.get('x', 960))
    return sum(xs) / len(xs) if xs else 960.0


def reenquadrar(sh, fmt, legendas=False):
    if fmt == '16:9':
        return sh
    f = FORMATOS[fmt]
    cw, x0 = f['cw'], 960 - f['cw'] / 2
    sh = copy.deepcopy(sh)
    if sh.get('type') == 'card':
        sh['wrap'], sh['size'] = f['wrap']
        return sh
    chars = sh.get('chars', [])
    # personagens: aproxima uns dos outros para caberem na janela
    if len(chars) > 1:
        xs = [c.get('x', 960) for c in chars]
        span, mid = max(xs) - min(xs), sum(xs) / len(xs)
        lim = cw * 0.62
        if span > lim:
            k = lim / span
            for c in chars:
                c['x'] = mid + (c.get('x', 960) - mid) * k
                for kf in c.get('path', []) or []:
                    kf['x'] = mid + (kf['x'] - mid) * k
    foco = _foco(sh)
    # camera: centra no foco, sem mostrar fora do mundo
    cam = sh.get('cam') or [{'t': 0, 'x': 960, 'y': 540, 'z': 1.0}]
    for kf in cam:
        z = kf.get('z', 1.0)
        meia = cw / (2 * z)
        kf['x'] = min(1920 - meia, max(meia, kf.get('x', 960) + (foco - 960)))
    sh['cam'] = cam
    lo, hi = f['safe'][0] * 1080, f['safe'][1] * 1080
    hi_txt = min(hi, 0.58 * 1080) if legendas else hi  # com legendas, os textos sobem para nao colidir
    for l in sh.get('layers', []):
        if l.get('kind') in NAO_MOVER:
            continue
        if l.get('hud'):
            # posicao relativa na tela 16:9 -> mesma posicao relativa na janela
            l['x'] = x0 + (l.get('x', 960) / 1920.0) * cw
            l['y'] = min(hi if l.get('legenda') else hi_txt, max(lo, l.get('y', 540)))
            if l.get('kind') == 'text':
                txt = l.get('text') or ''
                size0 = l.get('size', 84)
                size = size0 * (cw / 1920.0) * f['kt']
                if txt:
                    size = min(size, 0.92 * cw / max(1, len(txt) * 0.55))
                l['size'] = round(size, 1)
                if 'sw' in l:
                    l['sw'] = max(3, round(l['sw'] * size / max(1, size0), 1))
            else:
                l['s'] = round((l.get('s') or 1.0) * (cw / 1920.0) * f['ko'], 3)
        else:
            # objetos do mundo: puxa para perto do foco
            if 'x' in l:
                l['x'] = foco + (l['x'] - foco) * (0.45 if fmt == '9:16' else 0.75)
    return sh
