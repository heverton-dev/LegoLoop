import json
import pathlib
import re
import subprocess
import sys

RAIZ = pathlib.Path(__file__).resolve().parent.parent
WS = RAIZ / "docs/videorea/video-001"
AQUI = pathlib.Path(__file__).parent

def test_projeto_hyperframes_nao_depende_de_cdn():
    html = (WS / "scratch/projeto_hf/index.html").read_text(encoding="utf8")
    assert not re.search(r"""src=["']https?://""", html)
    assert (WS / "scratch/projeto_hf/gsap.min.js").exists()

def test_fixture_3s_1080p60_conforme():
    v = WS / "loop/fixture_3s.mp4"
    assert v.exists()
    r = subprocess.run([sys.executable, str(AQUI / "verificar_video.py"), "--video", str(v),
                        "--duracao", "3"], capture_output=True, text=True)
    assert r.returncode == 0, r.stderr

def test_dag_nao_referencia_dominio_pago():
    txt = (WS / "loop/pipeline.dag.json").read_text(encoding="utf8").lower()
    assert not any(d in txt for d in ("heygen.com", "elevenlabs", "higgsfield.ai"))
