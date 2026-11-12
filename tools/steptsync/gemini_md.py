"""Regenera o bloco gerenciado do ~/.gemini/GEMINI.md a partir do stept-rules.md (fonte única).

Uso: python tools/steptsync/gemini_md.py [GEMINI.md] [stept-rules.md]
Cria GEMINI.md.bak-stept na primeira execução e não o sobrescreve depois. Remove a seção "Respostas ao Usuário",
que contradiz o modo prova (manda usar tabelas), e troca o bloco entre os marcadores.
"""
import os
import re
import shutil
import sys

INI = "<!-- stept:inicio (gerado por tools/steptsync/gemini_md.py; não edite) -->"
FIM = "<!-- stept:fim -->"


def main() -> int:
    home = os.path.expanduser("~")
    gem = sys.argv[1] if len(sys.argv) > 1 else os.path.join(home, ".gemini", "GEMINI.md")
    regra = sys.argv[2] if len(sys.argv) > 2 else os.path.join(home, ".claude", "rules", "stept-rules.md")
    with open(gem, "rb") as f:
        texto = f.read().decode("utf-8-sig")
    with open(regra, "rb") as f:
        corpo = f.read().decode("utf-8").strip()
    if not os.path.exists(gem + ".bak-stept"):
        shutil.copyfile(gem, gem + ".bak-stept")
    texto = re.sub(re.escape(INI) + r".*?" + re.escape(FIM) + r"\n?", "", texto, flags=re.S)
    texto = re.sub(r"### Respostas ao Usuário\n.*?(?=\n### |\Z)", "", texto, flags=re.S)
    texto = texto.rstrip("\n") + "\n\n" + INI + "\n" + corpo + "\n" + FIM + "\n"
    with open(gem, "wb") as f:
        f.write(texto.encode("utf-8"))
    return 0


if __name__ == "__main__":
    sys.exit(main())
