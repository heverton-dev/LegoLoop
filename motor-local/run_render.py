#!/usr/bin/env python3
"""Ponte entre o contrato do motor Go (ExecutorPython: python run_render.py <dir> [args]) e o
render.py real do motor (que espera <timeline.json> --out <mp4> [args]).

Grava <dir>/render.mp4 (ArquivoRender) e copia os clipes por cena do cache real do render.py
para <dir>/cache/ (DirCache), um arquivo por cena, para a prova de determinismo (AC-E10.5).
"""
import os
import shutil
import subprocess
import sys

HERE = os.path.dirname(os.path.abspath(__file__))


def main():
    if len(sys.argv) < 2:
        sys.exit('uso: run_render.py <dir> [--formato 9:16] [--shots 1-1] [...]')
    projeto = os.path.abspath(sys.argv[1])
    extra = sys.argv[2:]
    fmt = '16:9'
    if '--formato' in extra:
        fmt = extra[extra.index('--formato') + 1]
    timeline = os.path.join(projeto, 'timeline.json')
    saida = os.path.join(projeto, 'render.mp4')
    if '--out' in extra:
        idx = extra.index('--out')
        saida = os.path.abspath(extra[idx + 1])
        del extra[idx:idx + 2]
    cmd = [sys.executable, os.path.join(HERE, 'render.py'), timeline, '--out', saida, '--sem-checagem'] + extra
    r = subprocess.run(cmd, cwd=HERE, capture_output=True, text=True)
    sys.stdout.write(r.stdout)
    sys.stderr.write(r.stderr)
    if r.returncode != 0:
        sys.exit(r.returncode)
    cache_real = os.path.join(projeto, '.render_cache', 'full_' + fmt.replace(':', 'x'))
    cache_go = os.path.join(projeto, 'cache')
    os.makedirs(cache_go, exist_ok=True)
    if os.path.isdir(cache_real):
        for nome in sorted(os.listdir(cache_real)):
            if nome.endswith('.mp4'):
                shutil.copyfile(os.path.join(cache_real, nome), os.path.join(cache_go, nome))
    print('run_render: ok, render.mp4 e cache/ gravados')


if __name__ == '__main__':
    main()
