#!/usr/bin/env bash
# Compila um livro de docs/books/<slug>/ e gera TUDO dentro da pasta do livro:
#   <slug>/<nome_base>.md   manuscrito consolidado
#   <slug>/<nome_base>.pdf  PDF 170x240mm com coluna de anotações
#   <slug>/build/livro.typ  intermediario (ignorado pelo git)
#   <slug>/preview/pNN.png  previews (se PyMuPDF estiver instalado)
# Uso: bash docs/books/compilar.sh <slug> [--sem-preview]
set -euo pipefail
cd "$(dirname "$0")"
slug="${1:?uso: compilar.sh <slug> [--sem-preview]}"
slug="${slug%/}"
[ -f "$slug/livro.json" ] || { echo "livro.json nao encontrado em $slug/" >&2; exit 1; }
nome=$(python -c "import json;print(json.load(open('$slug/livro.json',encoding='utf-8'))['nome_base'])")
tpl=$(python -c "import json;print(json.load(open('$slug/livro.json',encoding='utf-8')).get('template','ativos/livro.typst'))")
partes=$(python -c "import json;m=json.load(open('$slug/livro.json',encoding='utf-8'));print(' '.join(m['partes']))")
cd "$slug"
mkdir -p build preview img ativos
cp -r ../ativos/* ativos/ 2>/dev/null || true
cp -r ativos build/ 2>/dev/null || true
cp -r img build/ 2>/dev/null || true

# Sincroniza o template com o modelo oficial atualizado
if [ -f "ativos/livro.typst" ] && [ -f "../_modelo_livro.typst" ]; then
  cp "../_modelo_livro.typst" "ativos/livro.typst"
fi

# shellcheck disable=SC2086
python ../renderizar_mermaid.py $(for p in $partes; do echo "partes/$p"; done)

# shellcheck disable=SC2086
python ../normalizar_tabelas.py $(for p in $partes; do echo "partes/$p"; done)

: > "$nome.md"
for p in $partes; do cat "partes/$p" >> "$nome.md"; printf '\n\n' >> "$nome.md"; done
# shellcheck disable=SC2086
pandoc $(for p in $partes; do echo "partes/$p"; done) -f markdown -t typst --standalone --wrap=none \
  --template "$tpl" -o build/livro.typ
python ../desescapar_paineis.py build/livro.typ
typst compile --root . build/livro.typ "$nome.pdf"
echo "PDF: docs/books/$slug/$nome.pdf"
if [ "${2:-}" = "--com-preview" ]; then
  python -c "
import fitz,glob,os
for f in glob.glob('preview/p*.png'): os.remove(f)
d=fitz.open('$nome.pdf')
for i,p in enumerate(d): p.get_pixmap(dpi=70).save('preview/p%02d.png'%(i+1))
print('paginas:',len(d))" || echo "preview ignorado (PyMuPDF ausente)"
fi
