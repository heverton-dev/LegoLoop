import json
import pathlib
import subprocess

RAIZ = pathlib.Path(__file__).resolve().parent.parent
DAG = RAIZ / "docs/videorea/video-001/loop/pipeline.dag.json"

def nos():
    return {n["id"]: n for n in json.loads(DAG.read_text(encoding="utf8"))["nodes"]}

def antecessores(id_, n):
    vistos, fila = set(), list(n[id_].get("depends_on", []))
    while fila:
        x = fila.pop()
        if x not in vistos:
            vistos.add(x)
            fila += n[x].get("depends_on", [])
    return vistos

def test_oito_nos_em_ordem_linear():
    n = nos()
    ordem = ["n01_extrair_audio", "n02_transcrever", "n03_alinhar_roteiro", "n04_plano_cenas",
             "n05_gate_plano", "n06_renderizar_cenas", "n07_montar", "n08_verificar"]
    assert list(n) == ordem
    for a, b in zip(ordem, ordem[1:]):
        assert n[b]["depends_on"] == [a]

def test_gate_humano_do_plano_precede_o_render():
    n = nos()
    assert n["n05_gate_plano"]["type"] == "stop_gate"
    assert "prompt(" in n["n05_gate_plano"]["output"]["gate"]
    assert "n05_gate_plano" in antecessores("n06_renderizar_cenas", n)

def test_render_usa_hyperframes_local_a_60fps():
    cmd = nos()["n06_renderizar_cenas"]["command"]
    assert cmd[:2] == ["hyperframes", "render"]
    assert cmd[cmd.index("-f") + 1] == "60"

def test_montagem_ffmpeg_1080p60():
    cmd = " ".join(nos()["n07_montar"]["command"])
    assert nos()["n07_montar"]["command"][0] == "ffmpeg"
    assert "1920:1080" in cmd and "fps=60" in cmd

def test_verificacao_final_checa_final_mp4():
    n = nos()["n08_verificar"]
    assert "verificar_video.py" in " ".join(n["command"])
    assert n["inputs"]["video"].endswith("final.mp4')")

def test_dag_continua_valido():
    r = subprocess.run(["legoloop", "validate", str(DAG)], capture_output=True, text=True)
    assert r.returncode == 0, r.stdout + r.stderr
