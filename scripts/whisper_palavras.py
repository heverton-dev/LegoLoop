"""Converte o JSON completo do whisper-cli (-ojf) em lista de palavras."""


def _especial(texto):
    return texto.startswith("[_") and texto.endswith("]")


def _fechar(atual, saida):
    if atual:
        saida.append(atual)


def palavras_de_whisper(dados):
    """Retorna [{word, start_ms, end_ms}] juntando subpalavras em palavras."""
    saida, atual = [], None
    for seg in dados.get("transcription", []):
        for tok in seg.get("tokens", []):
            txt = tok["text"]
            if _especial(txt):
                continue
            ini, fim = tok["offsets"]["from"], tok["offsets"]["to"]
            if txt.startswith(" ") or atual is None:
                _fechar(atual, saida)
                atual = {"word": txt.strip(), "start_ms": ini, "end_ms": max(ini, fim)}
            else:
                atual["word"] += txt
                atual["end_ms"] = max(atual["end_ms"], fim)
    _fechar(atual, saida)
    return [p for p in saida if p["word"]]
