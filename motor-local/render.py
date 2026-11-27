#!/usr/bin/env python3
"""Paper Billions Animate Video engine renderer: timeline.json -> mp4 (every shot is cached, so a re-run only redraws what changed).

  python render.py <project>/timeline.json --out <project>/renders/final.mp4 [--draft] [--shots 3-9] [--fps 30]

--draft      1280x720, fast preview
--shots a-b  render only shots a..b (1-based) into the cache, no final file (use for checking one scene)
Needs: pip install playwright, python -m playwright install chromium, ffmpeg on PATH.
Paths inside timeline.json (bg, sprite closed/open) are relative to the timeline's folder.
"""
import argparse, array, hashlib, json, math, os, shutil, subprocess, sys
_here = os.path.dirname(os.path.abspath(__file__))
ENGINE_HASH = hashlib.sha1(b''.join(open(os.path.join(_here, f), 'rb').read() for f in ('engine.html', 'layers.js', 'puppet.js', 'sets.js', 'formato.py', 'legendas.py') if os.path.exists(os.path.join(_here, f)))).hexdigest()[:10]   # engine edits must invalidate the shot cache

from playwright.sync_api import sync_playwright
sys.path.insert(0, _here)
from formato import FORMATOS, janela, reenquadrar  # noqa: E402
import legendas as LEG  # noqa: E402

HERE = os.path.dirname(os.path.abspath(__file__))
IMG_KEYS = {'bg', 'closed', 'open', 'src'}


def load_sprite(ref, base):
    f = ref if os.path.isabs(ref) else os.path.join(base, ref)
    if not os.path.exists(f):
        sys.exit(f'MISSING SPRITE FILE: {f}  (run scripts/key_sprites.py first)')
    return json.load(open(f, encoding='utf-8'))


def absify(o, base):
    """turn every image path in a shot into a file:/// URL so the engine page can load it from anywhere"""
    if isinstance(o, dict):
        o = dict(o)
        if isinstance(o.get('sprite'), str):          # "sprite": "sprites/boss_cheer.json"
            sp = load_sprite(o['sprite'], base); o['sprite'] = sp
            if 'h' not in o and 'h' in sp: o['h'] = sp['h']
        return {k: (to_url(v, base) if k in IMG_KEYS and isinstance(v, str) else absify(v, base)) for k, v in o.items()}
    if isinstance(o, list):
        return [absify(x, base) for x in o]
    return o


def to_url(p, base):
    if p.startswith(('http', 'data:', 'file:')) or not p:
        return p
    f = p if os.path.isabs(p) else os.path.join(base, p)
    if not os.path.exists(f):
        sys.exit(f'MISSING FILE: {f}  (check the path in timeline.json)')
    return 'file:///' + os.path.abspath(f).replace('\\', '/').lstrip('/')


def audio_env(path, fps):
    """loudness per video frame (0..1): drives the mouths"""
    if not path or not os.path.exists(path):
        return []
    r = subprocess.run(['ffmpeg', '-v', 'error', '-i', path, '-f', 's16le', '-ac', '1', '-ar', '8000', '-'], capture_output=True)
    a = array.array('h'); a.frombytes(r.stdout[: len(r.stdout) // 2 * 2])
    win = 8000 // fps; env = []
    for i in range(0, len(a) - win, win):
        seg = a[i:i + win]; env.append(math.sqrt(sum(x * x for x in seg) / win) / 32768.0)
    m = max(env) if env else 1
    return [min(1.0, v / (m * 0.6 + 1e-9)) for v in env]


def render_bucket(items, W, H, fps, env, total, fmt='16:9', brand=None):
    with sync_playwright() as p:
        br = p.chromium.launch()
        FW, FH = FORMATOS[fmt]['out']; pg = br.new_page(viewport={'width': FW, 'height': FH}, device_scale_factor=(W / FW))
        pg.goto('file:///' + os.path.join(HERE, 'engine.html').replace('\\', '/').lstrip('/')); pg.wait_for_function('window.__ready===true')
        pg.evaluate('([a,f])=>setEnv(a,f)', [env, fps])
        x0, cw = janela(fmt); pg.evaluate('([a,b,c,d])=>setFormat(a,b,c,d)', [x0, cw, FW, FH])
        if brand: print('  marca:', pg.evaluate('b=>setBrand(b)', brand), flush=True)
        for si, sh, clip, d, st in sorted(items, key=lambda x: x[0]):
            n = max(2, int(round(d * fps)))
            pg.evaluate('([s,t])=>setShot(s,t)', [sh, st]); pg.evaluate('s=>preloadShot(s)', sh)
            ff = subprocess.Popen(['ffmpeg', '-y', '-v', 'error', '-f', 'image2pipe', '-framerate', str(fps), '-c:v', 'mjpeg', '-i', '-', '-vf', f'scale={W}:{H}',
                                   '-c:v', 'libx264', '-preset', 'veryfast', '-crf', '19', '-pix_fmt', 'yuv420p', clip + '.part.mp4'], stdin=subprocess.PIPE)
            for f in range(n):
                pg.evaluate('t=>render(t)', f / fps)
                ff.stdin.write(pg.screenshot(type='jpeg', quality=92))
            ff.stdin.close(); ff.wait(); os.replace(clip + '.part.mp4', clip)
            print(f'  shot {si+1:3d}/{total} {sh.get("id","")} {d:.2f}s done', flush=True)
        br.close()


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument('timeline'); ap.add_argument('--out'); ap.add_argument('--fps', type=int)
    ap.add_argument('--draft', action='store_true'); ap.add_argument('--formato', choices=list(FORMATOS)); ap.add_argument('--marca'); ap.add_argument('--legendas', choices=['sim', 'nao']); ap.add_argument('--shots'); ap.add_argument('--jobs', type=int, help='parallel browsers (default: number of CPUs, max 6)')
    ap.add_argument('--sem-checagem', action='store_true', help='pula a checagem obrigatória (check_timeline.py); use só em cortes derivados')
    a = ap.parse_args()
    tl_path = os.path.abspath(a.timeline); base = os.path.dirname(tl_path)
    if not a.sem_checagem and not a.shots:
        # portão: não renderiza um vídeo que o verificador reprova (cenas sem cenário, ritmo fora da faixa, voz ausente)
        chk = subprocess.run([sys.executable, os.path.join(_here, '..', 'check_timeline.py'), base], capture_output=True, text=True)
        print(chk.stdout.strip()[-1500:])
        if chk.returncode != 0 or 'CHECK FAILED' in chk.stdout:
            sys.exit('RENDER RECUSADO: check_timeline.py reprovou a timeline. Corrija os PROBLEM acima (ou use --sem-checagem em cortes derivados).')
    tl = json.load(open(tl_path, encoding='utf-8'))
    fps = a.fps or tl.get('fps', 30)
    fmt = a.formato or tl.get('formato', '16:9')
    W, H = FORMATOS[fmt]['draft' if a.draft else 'out']
    brand = None
    for cand in ([a.marca] if a.marca else []) + ([os.path.join(base, tl['marca'])] if tl.get('marca') else []) + [os.path.join(base, 'marca.json'), os.path.join(os.path.dirname(base), 'marca.json')]:
        if cand and os.path.exists(cand): brand = json.load(open(cand, encoding='utf-8')); print('  marca:', cand); break
    out = a.out or os.path.join(base, 'renders', 'final.mp4'); os.makedirs(os.path.dirname(out), exist_ok=True)
    cache = os.path.join(base, '.render_cache', ('draft' if a.draft else 'full') + '_' + fmt.replace(':', 'x')); os.makedirs(cache, exist_ok=True)
    voice = tl.get('audio'); voice = (voice if os.path.isabs(voice) else os.path.join(base, voice)) if voice else None
    sound = tl.get('soundtrack'); sound = (sound if os.path.isabs(sound) else os.path.join(base, sound)) if sound else None
    shots = tl['shots']; lo, hi = 1, len(shots)
    if a.shots:
        p = a.shots.split('-'); lo = int(p[0]); hi = int(p[-1])
    # start time of every shot (the mouths need the global clock)
    starts = []; t = 0.0
    for sh in shots: starts.append(t); t += float(sh['duration'])
    leg_on = bool(brand and fmt in brand.get('legendas', {}).get('ativo_em', [])) if a.legendas is None else a.legendas == 'sim'
    linhas = LEG.carregar(base) if leg_on else None
    if leg_on and linhas:
        shots = LEG.aplicar(json.loads(json.dumps(shots)), starts, linhas, brand.get('legendas', {}) if brand else {}, brand.get('cores', {}) if brand else {})
        print('  legendas: ligadas,', sum(1 for sh in shots for l in sh.get('layers', []) if l.get('legenda')), 'blocos')
    elif leg_on: print('  legendas: sem data/lines.json, desligadas')
    env = audio_env(voice, fps)
    clips = []; todo = []
    for si, sh0 in enumerate(shots):
        if not (lo <= si + 1 <= hi): continue
        sh = absify(reenquadrar(dict(sh0, cast=tl.get('cast', {})), fmt, bool(leg_on and linhas)), base)
        key = hashlib.sha1((ENGINE_HASH + fmt + json.dumps(brand or {}, sort_keys=True) + json.dumps(sh, sort_keys=True) + f'{starts[si]:.3f}{fps}{W}' + (str(os.path.getmtime(voice)) if voice and os.path.exists(voice) else '')).encode()).hexdigest()[:16]
        clip = os.path.join(cache, f'{si+1:03d}_{key}.mp4'); clips.append(clip)
        if os.path.exists(clip): print(f'  shot {si+1:3d}/{len(shots)} cached', flush=True)
        else: todo.append((si, sh, clip, float(sh['duration']), starts[si]))
    jobs = max(1, min(a.jobs or int(os.environ.get('ROBBI_JOBS', 0)) or (os.cpu_count() or 2), 6, len(todo) or 1))
    if todo:
        print(f'  rendering {len(todo)} shots with {jobs} parallel browser(s)', flush=True)
        buckets = [[] for _ in range(jobs)]; load = [0.0] * jobs
        for it in sorted(todo, key=lambda x: -x[3]): k = load.index(min(load)); buckets[k].append(it); load[k] += it[3]
        if jobs == 1: render_bucket(buckets[0], W, H, fps, env, len(shots), fmt, brand)
        else:
            import multiprocessing as mp
            ps = [mp.Process(target=render_bucket, args=(bk, W, H, fps, env, len(shots), fmt, brand)) for bk in buckets if bk]
            [q.start() for q in ps]; [q.join() for q in ps]
            if any(q.exitcode for q in ps): sys.exit('a render worker failed')
    if a.shots:
        print('Shots rendered into cache (no final file for a partial render).'); return
    lst = os.path.join(cache, 'list.txt'); open(lst, 'w').write(''.join(f"file '{c.replace(chr(92), '/')}'\n" for c in clips))
    silent = os.path.join(cache, 'silent.mp4'); subprocess.run(['ffmpeg', '-y', '-v', 'error', '-f', 'concat', '-safe', '0', '-i', lst, '-c', 'copy', silent], check=True)
    track = sound if sound and os.path.exists(sound) else voice
    if track and os.path.exists(track):
        subprocess.run(['ffmpeg', '-y', '-v', 'error', '-i', silent, '-i', track, '-map', '0:v', '-map', '1:a', '-c:v', 'copy', '-c:a', 'aac', '-b:a', '192k', '-shortest', out], check=True)
    else:
        shutil.copy(silent, out); print('  (no audio file: silent video)')
    print('Done:', out)


if __name__ == '__main__':
    main()
