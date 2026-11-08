#!/usr/bin/env python3
"""
desescapar_paineis.py — Pós-processador do LegoPress
Restaura chamadas de funções Typst (#painel-*, #frase-destaque, etc.)
que tenham sido indevidamente escapadas com contra-barra (\#) pelo Pandoc.
"""
import sys
import re

def desescapar_typst(caminho_arquivo):
    with open(caminho_arquivo, 'r', encoding='utf-8') as f:
        conteudo = f.read()

    # 1. Padrão para painéis e macros customizadas:
    # \#painel-(info|alerta|perigo|sucesso|teoria)("...", \[ ... \])
    def fix_painel(m):
        bloco = m.group(0)
        # remove a barra de \#painel- e \#frase-destaque
        bloco = re.sub(r'\\#(painel-(?:info|alerta|perigo|sucesso|teoria)|frase-destaque|chip|ficha)', r'#\1', bloco)
        # desescapa colchetes externos de typst \[ e \]
        bloco = bloco.replace(r'\[', '[').replace(r'\]', ']')
        # desescapa hashtags internas do typst (#emph, #strong, etc)
        bloco = re.sub(r'\\#([a-zA-Z0-9_\-]+)', r'#\1', bloco)
        return bloco

    # Regex abrange chamadas de painel em uma ou múltiplas linhas
    pattern = r'\\#(?:painel-(?:info|alerta|perigo|sucesso|teoria)|frase-destaque)\(.*?\n?\s*\\\]\)'
    conteudo_corrigido = re.sub(pattern, fix_painel, conteudo, flags=re.DOTALL)

    # 2. Correção de qualquer \#painel remanescente de linha única
    conteudo_corrigido = re.sub(r'\\#(painel-(?:info|alerta|perigo|sucesso|teoria)|frase-destaque)', r'#\1', conteudo_corrigido)

    with open(caminho_arquivo, 'w', encoding='utf-8') as f:
        f.write(conteudo_corrigido)
    print(f"[OK] Painéis Typst desescapados com sucesso em: {caminho_arquivo}")

if __name__ == '__main__':
    if len(sys.argv) < 2:
        print("Uso: desescapar_paineis.py <arquivo.typ>")
        sys.exit(1)
    desescapar_typst(sys.argv[1])
