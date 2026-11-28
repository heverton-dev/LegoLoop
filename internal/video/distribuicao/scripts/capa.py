#!/usr/bin/env python3
"""Capa de video 1280x720: fundo sem texto + texto e logo pelo codigo (Pillow).

Uso: capa.py --fundo F --linha1 T --saida S [--linha2 T] [--logo L] [--fonte TTF]
            [--ouro #RRGGBB] [--texto #RRGGBB] [--texto-direita]
O canto inferior direito fica livre para a duracao que a rede desenha por cima.
"""
import argparse
import os
import sys

from PIL import Image, ImageDraw, ImageFont

W, H = 1280, 720
FONTES = [
    'C:/Windows/Fonts/arialbd.ttf',
    '/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf',
    '/Library/Fonts/Arial Bold.ttf',
]


def cobrir(im, w, h):
    r = max(w / im.width, h / im.height)
    im = im.resize((int(im.width * r) + 1, int(im.height * r) + 1), Image.LANCZOS)
    x, y = (im.width - w) // 2, (im.height - h) // 2
    return im.crop((x, y, x + w, y + h))


def achar_fonte(pedida):
    for f in ([pedida] if pedida else []) + FONTES:
        if f and os.path.exists(f):
            return f
    sys.exit('nenhuma fonte TrueType encontrada; passe --fonte')


def ajustar(d, texto, caminho, tam_max, largura):
    tam = tam_max
    while tam > 30:
        f = ImageFont.truetype(caminho, tam)
        if d.textlength(texto, font=f) <= largura:
            return f
        tam -= 4
    return ImageFont.truetype(caminho, 30)


def degrade(im, lado_dir):
    grad = Image.new('RGBA', (W, H), (0, 0, 0, 0))
    gd = ImageDraw.Draw(grad)
    for x in range(W):
        t = (x / W) if lado_dir else (1 - x / W)
        gd.line((x, 0, x, H), fill=(11, 14, 28, int(max(0, (t - 0.35) / 0.65) * 190)))
    return Image.alpha_composite(im.convert('RGBA'), grad).convert('RGB')


def main():
    a = argparse.ArgumentParser()
    a.add_argument('--fundo', required=True)
    a.add_argument('--linha1', required=True)
    a.add_argument('--linha2', default='')
    a.add_argument('--saida', required=True)
    a.add_argument('--logo', default='')
    a.add_argument('--fonte', default='')
    a.add_argument('--ouro', default='#FFD24A')
    a.add_argument('--texto', default='#FFFFFF')
    a.add_argument('--texto-direita', action='store_true')
    o = a.parse_args()
    im = degrade(cobrir(Image.open(o.fundo).convert('RGB'), W, H), o.texto_direita)
    d = ImageDraw.Draw(im)
    zona = 620
    x0 = (W - zona - 50) if o.texto_direita else 50
    fonte = achar_fonte(o.fonte)
    f1 = ajustar(d, o.linha1.upper(), fonte, 250, zona)
    d.text((x0, 150), o.linha1.upper(), font=f1, fill=o.ouro, stroke_width=9, stroke_fill='#000000')
    b1 = d.textbbox((x0, 150), o.linha1.upper(), font=f1, stroke_width=9)
    if o.linha2:
        f2 = ajustar(d, o.linha2.upper(), fonte, 190, zona)
        d.text((x0, b1[3] + 10), o.linha2.upper(), font=f2, fill=o.texto, stroke_width=8, stroke_fill='#000000')
    if o.logo and os.path.exists(o.logo):
        logo = Image.open(o.logo).convert('RGB').resize((110, 110), Image.LANCZOS)
        m = Image.new('L', (110, 110), 0)
        ImageDraw.Draw(m).ellipse((0, 0, 109, 109), fill=255)
        im.paste(logo, (30 if o.texto_direita else W - 140, 30), m)
    os.makedirs(os.path.dirname(os.path.abspath(o.saida)), exist_ok=True)
    im.save(o.saida, optimize=True)
    print('ok', o.saida, im.size)


if __name__ == '__main__':
    main()
