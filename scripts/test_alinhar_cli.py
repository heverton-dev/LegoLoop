import json
import pathlib
import subprocess
import sys

AQUI = pathlib.Path(__file__).parent

def test_cli_gera_segmentos_json(tmp_path):
    roteiro = tmp_path / "roteiro.txt"
    roteiro.write_text("Eu sempre tive uma grande limitação\nMas hoje eu consigo\n", encoding="utf8")
    saida = tmp_path / "segmentos.json"
    r = subprocess.run(
        [sys.executable, str(AQUI / "alinhar_roteiro.py"),
         "--roteiro", str(roteiro),
         "--whisper", str(AQUI / "fixtures" / "whisper_amostra.json"),
         "--out", str(saida)],
        capture_output=True, text=True)
    assert r.returncode == 0, r.stderr
    segs = json.loads(saida.read_text(encoding="utf8"))
    assert [s["id"] for s in segs] == [1, 2]

def test_cli_falha_com_exit_1_se_fala_ausente(tmp_path):
    roteiro = tmp_path / "roteiro.txt"
    roteiro.write_text("frase que nunca foi dita\n", encoding="utf8")
    r = subprocess.run(
        [sys.executable, str(AQUI / "alinhar_roteiro.py"),
         "--roteiro", str(roteiro),
         "--whisper", str(AQUI / "fixtures" / "whisper_amostra.json"),
         "--out", str(tmp_path / "x.json")],
        capture_output=True, text=True)
    assert r.returncode == 1
