#!/usr/bin/env bash
# Compila o livro a partir de livro.json, audita a geometria e gera previews.
set -euo pipefail
cd "$(dirname "$0")/../../.."
B=docs/livro-verbos-bdd
N=$(python -c "import json;print(json.load(open('$B/livro.json',encoding='utf-8'))['nome_base'])")
H=$(sha256sum docs/manuals/2026-11-07_MANUAL_DIDATICO_VERBOS_BDD.md | cut -d' ' -f1)
python -c "
import json;p='$B/livro.json';m=json.load(open(p,encoding='utf-8'));m['hash_fonte']='$H'
json.dump(m,open(p,'w',encoding='utf-8'),ensure_ascii=False,indent=2)"
: > "$B/$N.md"
for p in $(python -c "import json;print(' '.join(json.load(open('$B/livro.json',encoding='utf-8'))['partes']))"); do
  cat "$B/partes/$p" >> "$B/$N.md"; printf '\n\n' >> "$B/$N.md"
done
pandoc "$B/$N.md" --template "$B/ativos/livro.typst" --pdf-engine typst --columns=1000 -o "$B/$N.pdf"
legoloop legopress auditar "$B/$N.pdf"
mkdir -p "$B/preview"; rm -f "$B"/preview/*.png
python -c "
import fitz
d=fitz.open('$B/$N.pdf')
for i,p in enumerate(d): p.get_pixmap(dpi=70).save('$B/preview/p%02d.png'%(i+1))
print('paginas:',len(d))"
