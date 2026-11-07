#!/usr/bin/env bash
# Compila uma obra da pasta docs/livro/obras/<slug> em PDF 170x240mm.
# Uso: bash docs/livro/obras/compilar.sh <slug>
set -euo pipefail
cd "$(dirname "$0")"
slug="$1"
manifesto="$slug/livro.json"
nome=$(python -c "import json;print(json.load(open('$manifesto',encoding='utf-8'))['nome_base'])")
partes=$(python -c "import json;m=json.load(open('$manifesto',encoding='utf-8'));print(' '.join('$slug/partes/'+p for p in m['partes']))")
mkdir -p "$slug/build"
# shellcheck disable=SC2086
pandoc $partes -f markdown -t typst --standalone --wrap=none \
  --template _ativos/livro.typst -o "$slug/build/livro.typ"
typst compile --root . "$slug/build/livro.typ" "../$nome.pdf"
echo "PDF: docs/livro/$nome.pdf"
