import json
import pathlib
import subprocess

RAIZ = pathlib.Path(__file__).resolve().parent.parent
DAG = RAIZ / "docs/videorea/video-001/loop/pipeline.dag.json"

def carregar():
    return json.loads(DAG.read_text(encoding="utf8"))

def nos():
    return {n["id"]: n for n in carregar()["nodes"]}

def test_tres_primeiros_nos_em_cadeia():
    n = nos()
    assert n["n01_extrair_audio"]["depends_on"] == [] if "depends_on" in n["n01_extrair_audio"] else True
    assert n["n02_transcrever"]["depends_on"] == ["n01_extrair_audio"]
    assert n["n03_alinhar_roteiro"]["depends_on"] == ["n02_transcrever"]

def test_transcricao_usa_whisper_cli_local_com_json_completo():
    cmd = nos()["n02_transcrever"]["command"]
    assert "whisper-cli" in cmd[0].lower()
    assert "-ojf" in cmd

def test_audio_em_16khz_mono_no_scratch():
    cmd = nos()["n01_extrair_audio"]["command"]
    assert cmd[0] == "ffmpeg"
    assert "16000" in cmd and "1" in cmd
    assert all("scratch/" in c for c in cmd if c.endswith(".wav"))

def test_cada_no_tem_gate_de_artefato():
    for n in nos().values():
        assert n["output"]["gate"].startswith("file_exists(")

def test_sem_dominio_externo_pago():
    texto = DAG.read_text(encoding="utf8").lower()
    for proibido in ("heygen.com", "elevenlabs", "higgsfield.ai", "api.openai.com"):
        assert proibido not in texto

def test_dag_valido_para_o_legoloop():
    r = subprocess.run(["legoloop", "validate", str(DAG)], capture_output=True, text=True)
    assert r.returncode == 0, r.stdout + r.stderr
