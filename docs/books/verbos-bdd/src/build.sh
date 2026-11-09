#!/usr/bin/env bash
# Atalho: compila o livro Verbos/BDD e audita a geometria do PDF.
set -euo pipefail
cd "$(dirname "$0")/../.."
bash compilar.sh verbos-bdd
legoloop legopress auditar "verbos-bdd/$(python -c "import json;print(json.load(open('verbos-bdd/livro.json',encoding='utf-8'))['nome_base'])").pdf"
