import json
import pathlib
from whisper_palavras import palavras_de_whisper
from alinhar_roteiro import alinhar

AMOSTRA = pathlib.Path(__file__).parent / "fixtures" / "whisper_amostra.json"

def carregar():
    return json.loads(AMOSTRA.read_text(encoding="utf8"))

def test_junta_subpalavras_e_descarta_tokens_especiais():
    ps = palavras_de_whisper(carregar())
    texto = " ".join(p["word"] for p in ps)
    assert texto.startswith("Eu sempre tive uma grande limitação")
    assert not any(p["word"].startswith("[_") for p in ps)

def test_tempos_monotonicos_e_positivos():
    ps = palavras_de_whisper(carregar())
    assert len(ps) > 20
    for a, b in zip(ps, ps[1:]):
        assert a["start_ms"] <= b["start_ms"]
    assert all(p["end_ms"] >= p["start_ms"] for p in ps)

def test_alinha_roteiro_real_com_saida_do_whisper_cli():
    ps = palavras_de_whisper(carregar())
    roteiro = ["Eu sempre tive uma grande limitação no meu negócio",
               "Mas hoje eu consigo"]
    segs = alinhar(roteiro, ps)
    assert len(segs) == 2
    assert segs[0]["inicio_ms"] == ps[0]["start_ms"]
    assert segs[0]["fim_ms"] <= segs[1]["inicio_ms"]
