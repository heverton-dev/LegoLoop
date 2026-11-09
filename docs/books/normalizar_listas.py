#!/usr/bin/env python3
"""
normalizar_listas.py — Padrão LegoPress
Garante que toda lista em Markdown (bullet points '*' ou '-' ou ordenada '1.')
possua obrigatoriamente uma linha em branco antes do seu primeiro item,
evitando que o Pandoc renderize os itens como parágrafo contínuo em prosa.
"""
import sys
import re

def normalizar_arquivo(caminho):
    with open(caminho, 'r', encoding='utf-8') as f:
        linhas = f.readlines()

    novas_linhas = []
    modificado = False

    for i, linha in enumerate(linhas):
        curr = linha.rstrip('\r\n')
        is_list_item = bool(re.match(r'^(\*|-|\d+\.)\s+', curr))

        if is_list_item and novas_linhas:
            prev = novas_linhas[-1].rstrip('\r\n')
            prev_is_list = bool(re.match(r'^(\*|-|\d+\.)\s+', prev))
            prev_is_blank = (prev.strip() == '')

            # Se o anterior não é em branco e não é outro item de lista, injeta linha em branco
            if not prev_is_blank and not prev_is_list:
                novas_linhas.append('\n')
                modificado = True

        novas_linhas.append(linha)

    if modificado:
        with open(caminho, 'w', encoding='utf-8') as f:
            f.writelines(novas_linhas)
        print(f"[OK] Listas normalizadas com linha em branco em: {caminho}")
    return modificado

if __name__ == '__main__':
    for arq in sys.argv[1:]:
        normalizar_arquivo(arq)
