import pathlib
import subprocess
import sys

AQUI = pathlib.Path(__file__).parent

def gerar(tmp, size, fps, dur):
    out = tmp / "v.mp4"
    subprocess.run(["ffmpeg", "-y", "-v", "error", "-f", "lavfi", "-i",
                    f"color=c=black:s={size}:r={fps}:d={dur}", "-f", "lavfi", "-i",
                    f"anullsrc=r=48000:cl=stereo", "-t", str(dur), "-shortest",
                    "-pix_fmt", "yuv420p", str(out)], check=True)
    return out

def rodar(v, dur):
    return subprocess.run([sys.executable, str(AQUI / "verificar_video.py"), "--video", str(v),
                           "--largura", "1920", "--altura", "1080", "--fps", "60",
                           "--duracao", str(dur), "--tolerancia", "0.1"],
                          capture_output=True, text=True)

def test_video_conforme_passa(tmp_path):
    assert rodar(gerar(tmp_path, "1920x1080", 60, 3), 3).returncode == 0

def test_fps_errado_falha(tmp_path):
    r = rodar(gerar(tmp_path, "1920x1080", 30, 3), 3)
    assert r.returncode == 1 and "fps" in r.stderr

def test_duracao_fora_da_tolerancia_falha(tmp_path):
    r = rodar(gerar(tmp_path, "1920x1080", 60, 5), 3)
    assert r.returncode == 1 and "duracao" in r.stderr

def test_arquivo_inexistente_falha(tmp_path):
    assert rodar(tmp_path / "nao.mp4", 3).returncode == 1

def test_relatorio_json_gravado_quando_solicitado(tmp_path):
    import json
    v = gerar(tmp_path, "1920x1080", 60, 3)
    out = tmp_path / "verificacao.json"
    r = subprocess.run([sys.executable, str(AQUI / "verificar_video.py"), "--video", str(v),
                        "--duracao", "3", "--saida", str(out)], capture_output=True, text=True)
    assert r.returncode == 0, r.stderr
    rel = json.loads(out.read_text(encoding="utf8"))
    assert rel["conforme"] is True and rel["falhas"] == []
