#!/usr/bin/env python3
"""Recorte transparente (rembg + despill verde) e fatiamento de sprites LegoLoop."""
import argparse
import json
import os
import sys

import numpy as np
from PIL import Image
from rembg import remove


def despill_verde(rgba_img):
    arr = np.array(rgba_img).astype(float)
    r, g, b, a = arr[..., 0], arr[..., 1], arr[..., 2], arr[..., 3]
    gd = g - np.maximum(r, b)
    mask_verde = gd > 75
    arr[..., 1] = np.where(gd > 0, g - gd * 0.85, g)
    arr[..., 3] = np.where(mask_verde, 0, a)
    return Image.fromarray(arr.clip(0, 255).astype("uint8"), "RGBA")


def montar_celulas(w, h, layout):
    cols, rows = 2, 1
    if layout == "3x2":
        cols, rows = 3, 2
    elif layout == "6x":
        cols, rows = 3, 2
    cw, ch = w // cols, h // rows
    celulas = []
    idx = 0
    for r in range(rows):
        for c in range(cols):
            idx += 1
            tipo = "expressao" if layout == "3x2" else ("pose" if layout == "6x" or idx == 1 else "expressao")
            celulas.append((tipo, (c * cw, r * ch, (c + 1) * cw, (r + 1) * ch)))
    return celulas


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("--projeto", required=True)
    ap.add_argument("--personagem-id", required=True)
    ap.add_argument("--entrada", required=True)
    ap.add_argument("--layout", choices=["1x2", "3x2", "6x"], default="1x2")
    args = ap.parse_args()

    if not os.path.exists(args.entrada):
        sys.exit(f"ficha nao encontrada: {args.entrada}")

    im = Image.open(args.entrada).convert("RGBA")
    recortada = despill_verde(remove(im))
    w, h = recortada.size
    celulas = montar_celulas(w, h, args.layout)

    outdir = os.path.join(args.projeto, "sprites")
    os.makedirs(outdir, exist_ok=True)

    saida = []
    for idx, (tipo, box) in enumerate(celulas, start=1):
        crop = recortada.crop(box)
        bb = crop.getchannel("A").point(lambda v: 255 if v > 20 else 0).getbbox()
        if not bb:
            continue
        crop = crop.crop(bb)
        sp_id = f"{args.personagem_id}_{tipo}_{idx}"
        rel_png = f"sprites/{sp_id}.png"
        crop.save(os.path.join(args.projeto, rel_png), format="PNG")
        cw, ch = crop.size
        saida.append({
            "id": sp_id,
            "personagem_id": args.personagem_id,
            "tipo": tipo,
            "arquivo_png": rel_png,
            "largura": cw,
            "altura": ch,
            "alfa_valido": True,
            "despill_verde": True,
            "bounding_box": {"x": box[0] + bb[0], "y": box[1] + bb[1], "w": cw, "h": ch},
        })
    if not saida:
        sys.exit(f"nenhum sprite valido encontrado em {args.entrada}")
    print(json.dumps(saida))


if __name__ == "__main__":
    main()
