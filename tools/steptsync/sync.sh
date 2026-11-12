#!/bin/sh
# Espelha o motor internal/ste (fonte unica) para o stept-cli avulso e recompila os binarios.
set -e
SRC="$(cd "$(dirname "$0")/../.." && pwd)"
DST="${1:-$HOME/Desktop/stept-cli}"
for f in prompt.go prompt_prova.go switch.go manual.go scan.go sync.go; do
  cp "$SRC/internal/ste/$f" "$DST/internal/ste/$f"
done
cd "$DST" && go test ./... && go build -o bin/stept.exe ./cmd/stept
