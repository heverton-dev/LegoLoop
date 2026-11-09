"""Verifica resolução, fps, duração e áudio de um vídeo com ffprobe."""
import argparse
import json
import subprocess
import sys


def _sondar(caminho):
    r = subprocess.run(
        ["ffprobe", "-v", "error", "-show_entries",
         "stream=codec_type,width,height,r_frame_rate:format=duration",
         "-of", "json", caminho], capture_output=True, text=True)
    if r.returncode != 0:
        raise OSError(r.stderr.strip() or "ffprobe falhou")
    return json.loads(r.stdout)


def _fps(txt):
    num, _, den = txt.partition("/")
    return float(num) / float(den or 1)


def verificar(dados, largura, altura, fps, duracao, tolerancia):
    """Retorna lista de falhas (vazia se conforme)."""
    falhas = []
    v = next((s for s in dados["streams"] if s["codec_type"] == "video"), None)
    if v is None:
        return ["sem stream de video"]
    if (v["width"], v["height"]) != (largura, altura):
        falhas.append(f"resolucao {v['width']}x{v['height']} != {largura}x{altura}")
    if abs(_fps(v["r_frame_rate"]) - fps) > 0.01:
        falhas.append(f"fps {_fps(v['r_frame_rate']):.2f} != {fps}")
    dur = float(dados["format"]["duration"])
    if abs(dur - duracao) > tolerancia:
        falhas.append(f"duracao {dur:.2f}s fora de {duracao}s +- {tolerancia}")
    if not any(s["codec_type"] == "audio" for s in dados["streams"]):
        falhas.append("sem stream de audio")
    return falhas


def _main(argv=None):
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--video", required=True)
    ap.add_argument("--largura", type=int, default=1920)
    ap.add_argument("--altura", type=int, default=1080)
    ap.add_argument("--fps", type=float, default=60)
    ap.add_argument("--duracao", type=float, required=True)
    ap.add_argument("--tolerancia", type=float, default=0.1)
    ap.add_argument("--saida")
    a = ap.parse_args(argv)
    try:
        falhas = verificar(_sondar(a.video), a.largura, a.altura, a.fps, a.duracao, a.tolerancia)
    except OSError as e:
        falhas = [f"erro ao ler video: {e}"]
    for f in falhas:
        print(f"erro: {f}", file=sys.stderr)
    if a.saida:
        with open(a.saida, "w", encoding="utf8") as fh:
            json.dump({"conforme": not falhas, "falhas": falhas}, fh, ensure_ascii=False)
    return 1 if falhas else 0


if __name__ == "__main__":
    raise SystemExit(_main())
