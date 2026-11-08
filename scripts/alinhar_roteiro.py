"""Alinha falas de um roteiro às palavras transcritas (função pura)."""
import re
import unicodedata

LOOKAHEAD = 80  # fala real tem trechos fora do roteiro


def _norm(texto):
    s = unicodedata.normalize("NFD", texto.lower())
    s = "".join(c for c in s if not unicodedata.combining(c))
    return re.sub(r"[^a-z0-9]", "", s)


def _tokens(frase):
    return [t for t in (_norm(w) for w in frase.split()) if t]


def _achar_inicio(toks, norm, cursor):
    """Primeira posição onde até 2 tokens iniciais da fala coincidem."""
    alvo = toks[:2]
    for i in range(cursor, min(cursor + LOOKAHEAD + 1, len(norm))):
        if norm[i:i + len(alvo)] == alvo:
            return i
    return None


def _achar_fim(toks, norm, inicio):
    previsto = inicio + len(toks) - 1
    for i in range(max(inicio, previsto - 3), min(previsto + 4, len(norm))):
        if norm[i] == toks[-1] and i >= inicio:
            return i
    return min(previsto, len(norm) - 1)


def alinhar(roteiro, palavras):
    """Retorna [{id, texto, inicio_ms, fim_ms}] ordenados e sem sobreposição."""
    norm = [_norm(p["word"]) for p in palavras]
    segs, cursor = [], 0
    for n, frase in enumerate(roteiro, 1):
        toks = _tokens(frase)
        ini = _achar_inicio(toks, norm, cursor) if toks else None
        if ini is None:
            raise ValueError(f"fala {n} não encontrada na transcrição: {frase!r}")
        fim = _achar_fim(toks, norm, ini)
        segs.append({"id": n, "texto": frase,
                     "inicio_ms": palavras[ini]["start_ms"],
                     "fim_ms": palavras[fim]["end_ms"]})
        cursor = fim + 1
    return segs


def _main(argv=None):
    import argparse
    import json
    import sys
    from whisper_palavras import palavras_de_whisper

    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--roteiro", required=True)
    ap.add_argument("--whisper", required=True)
    ap.add_argument("--out", required=True)
    a = ap.parse_args(argv)
    with open(a.roteiro, encoding="utf8") as f:
        falas = [l.strip() for l in f if l.strip()]
    with open(a.whisper, encoding="utf8") as f:
        palavras = palavras_de_whisper(json.load(f))
    try:
        segs = alinhar(falas, palavras)
    except ValueError as e:
        print(f"erro: {e}", file=sys.stderr)
        return 1
    with open(a.out, "w", encoding="utf8") as f:
        json.dump(segs, f, ensure_ascii=False, indent=2)
    return 0


if __name__ == "__main__":
    raise SystemExit(_main())
