import pytest
from alinhar_roteiro import alinhar

def palavras(texto, passo=400, ini=0):
    out, t = [], ini
    for w in texto.split():
        out.append({"word": w, "start_ms": t, "end_ms": t + passo - 50})
        t += passo
    return out

ROTEIRO = [
    "Eu sempre tive uma limitação",
    "hoje consigo gravar vídeos editados",
    "vou mostrar todo o sistema",
    "você precisa do Claude Code",
    "se inscreva no canal",
]

def test_cinco_falas_geram_cinco_segmentos_ordenados():
    segs = alinhar(ROTEIRO, palavras(" ".join(ROTEIRO)))
    assert len(segs) == 5
    assert [s["id"] for s in segs] == [1, 2, 3, 4, 5]
    for a, b in zip(segs, segs[1:]):
        assert a["fim_ms"] <= b["inicio_ms"]
    assert all(s["fim_ms"] > s["inicio_ms"] for s in segs)

def test_ignora_acentos_pontuacao_e_caixa():
    segs = alinhar(["Olá, mundo!"], palavras("ola MUNDO"))
    assert segs[0]["inicio_ms"] == 0 and segs[0]["fim_ms"] == 750

def test_texto_do_segmento_preserva_roteiro():
    segs = alinhar(ROTEIRO, palavras(" ".join(ROTEIRO)))
    assert segs[2]["texto"] == ROTEIRO[2]

def test_fala_do_roteiro_ausente_levanta_erro():
    with pytest.raises(ValueError):
        alinhar(["frase que nunca foi dita"], palavras("outra coisa totalmente"))
