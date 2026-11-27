#!/usr/bin/env python3
"""Validate a project before rendering and measure its pacing.
  python check_timeline.py <project>     -> problems, warnings, pacing numbers, then CHECK OK or CHECK FAILED
Pacing targets (taken from the reference format): 12-21 cuts per minute, median shot 2.5-4 s, 4+ moving things per scene shot,
a typed black card every 30-60 s, no plate shown in more than 3 shots in a row, a new picture (plate) almost every shot."""
import json, os, sys, collections, subprocess
if len(sys.argv) > 1 and sys.argv[1] in ('-h', '--help'): print(__doc__); sys.exit(0)
P = sys.argv[1] if len(sys.argv) > 1 else '.'
bad, warn = [], []
try: tl = json.load(open(os.path.join(P, 'timeline.json'), encoding='utf-8'))
except Exception as e: print('CHECK FAILED: timeline.json:', e); sys.exit(1)
CAST = tl.get('cast', {})
TYPES = {'scene', 'numbers', 'chart', 'linechart', 'market', 'triptych', 'card'}
OVK = {'tag', 'label', 'text', 'lower', 'ticker', 'clock'}
FX = {'dust', 'bokeh', 'rain', 'bills', 'confetti', 'sparks', 'steam'}
OBJ = {'cash', 'fire', 'balloon', 'stopwatch', 'ekg'}
LK = {'img', 'text', 'rect', 'circle', 'poly', 'path', 'stamp', 'doc', 'pin', 'inset', 'calendar', 'bars', 'cash', 'fire', 'balloon', 'stopwatch', 'ekg'}
IN = {'pop', 'slideL', 'slideR', 'slideU', 'slideD', 'drop', 'fade', 'grow', 'stamp', 'swing', 'flipX'}
LOOP = {'bob', 'shake', 'swing', 'pulse', 'spin', 'float', 'flicker'}
TR = {'whipL', 'whipR', 'whipU', 'whipD', 'fog', 'zoom', 'focus', 'fade', 'flash'}
EXPR = {'neutral', 'smile', 'grin', 'shock', 'worry', 'angry', 'sad', 'smug'}
shots = tl.get('shots', [])
if not shots: bad.append('no shots')
exists = lambda p: (p.startswith(('http', 'data:')) or os.path.exists(os.path.join(P, p)))

SETKINDS = {'shop','office','street','city','boardroom','stage','trading','map','vault','server','capitol','living'}
def check_char(n, c):
    if c.get('puppet') is not None:
        p = c['puppet']
        if isinstance(p, str) and p not in CAST: bad.append(f'{n}: puppet "{p}" is not in timeline "cast"')
        return
    check_sprite(n, c.get('sprite'))
def check_sprite(n, sp):
    if isinstance(sp, str):
        f = os.path.join(P, sp)
        if not os.path.exists(f): bad.append(f'{n}: missing sprite json {sp} (run key_sprites.py)'); return
        j = json.load(open(f))
        for key in ('closed', 'open'):
            if j.get(key) and not exists(j[key]): bad.append(f'{n}: missing sprite image {j[key]}')
    elif not isinstance(sp, dict): bad.append(f'{n}: "sprite" must be a json path')

def moving(sh):
    """how many things move in a shot"""
    m = 0
    for l in sh.get('layers', []):
        if l.get('in') or l.get('loop') or l.get('kf') or l.get('draw') or l.get('typed') or l.get('counter') or l.get('kind') in ('doc', 'pin', 'calendar', 'bars', 'cash', 'fire', 'balloon', 'stopwatch', 'ekg', 'stamp') or l.get('wipe'): m += 1
    for c in sh.get('chars', []):
        m += 1 if (c.get('pose') not in (None, 'idle') or c.get('path') or c.get('swap') or c.get('talk') or c.get('enter')) else 0.5
    m += len(sh.get('fx', [])) + len(sh.get('overlays', [])) * .5 + len(sh.get('items', [])) + len(sh.get('objects', []))
    if sh.get('cam'): m += 1
    if sh.get('shake'): m += .5
    if sh.get('grade') and isinstance(sh['grade'], list) or (isinstance(sh.get('grade'), dict) and sh['grade'].get('a1') is not None): m += .5
    for pn in sh.get('panels', []): m += moving(pn) * .6 + 1
    if sh.get('crowd'): m += 1
    return m

ids = set(); total = 0.0; run = 0; last = None; plate_run = 0; lastplate = None
durs = []; cards = []; people_run = 0; plates = collections.Counter(); starts = []; t0 = 0.0; poor = []
for i, sh in enumerate(shots, 1):
    n = f'shot {i} ({sh.get("id","?")})'; ty = sh.get('type', 'scene'); d = sh.get('duration', 0); starts.append(t0)
    if isinstance(d, (int, float)): total += d; t0 += d; durs.append(d)
    if sh.get('id') in ids: bad.append(f'{n}: duplicate id')
    ids.add(sh.get('id'))
    if ty not in TYPES: bad.append(f'{n}: unknown type "{ty}" (allowed: {sorted(TYPES)})')
    if not isinstance(d, (int, float)) or d <= 0: bad.append(f'{n}: duration must be > 0'); continue
    if d > 7 and ty == 'scene': warn.append(f'{n}: {d:.1f}s is long for a scene: split it into two diorama shots')
    elif d < 0.7: warn.append(f'{n}: only {d:.2f}s: too short to read')
    if ty == 'card': cards.append(t0 - d)
    pk = sh.get('bg') or (sh['set'].get('kind', '?') + '/' + sh['set'].get('time', 'night') if sh.get('set') else None)
    if sh.get('bg') and not exists(sh['bg']): bad.append(f'{n}: missing file {sh["bg"]}')
    if pk:
        plates[pk] += 1
        plate_run = plate_run + 1 if pk == lastplate else 1; lastplate = pk
        if plate_run > 3: warn.append(f'{n}: the same plate {plate_run} shots in a row; new place or new angle needed')
    if ty == 'scene' and not sh.get('bg') and not sh.get('set'): bad.append(f'{n}: scene needs "set" (code set) or "bg" (a plate)')
    if sh.get('set') and sh['set'].get('kind') not in SETKINDS: bad.append(f'{n}: unknown set kind {sh["set"].get("kind")} (allowed: {sorted(SETKINDS)})')
    if ty == 'scene':
        people_run = people_run + 1 if sh.get('chars') else 0
        if people_run > 3: warn.append(f'{n}: more than 3 shots with people in a row; add an object/number/map/card shot')
        mv = moving(sh)
        if d >= 2 and mv < 3.5: poor.append((n, mv))
    for pn in sh.get('panels', []):
        if pn.get('bg') and not exists(pn['bg']): bad.append(f'{n}: missing panel file {pn["bg"]}')
        for c in pn.get('chars', []): check_char(n, c)
        for l in pn.get('layers', []):
            if l.get('src') and not exists(l['src']): bad.append(f'{n}: missing panel layer file {l["src"]}')
    for c in sh.get('chars', []):
        check_char(n, c)
        for s2 in c.get('swap', []):
            if s2.get('sprite'): check_sprite(n, s2['sprite'])
        if not (0 <= c.get('x', 960) <= 1920): bad.append(f'{n}: character x={c.get("x")} is off-screen')
        if c.get('y', 1000) < 700 or c.get('y', 1000) > 1200: warn.append(f'{n}: character feet y={c.get("y")} (floor is usually 950-1040)')
        if c.get('expr') and c['expr'] not in EXPR: warn.append(f'{n}: unknown expr {c["expr"]}')
    for j, l in enumerate(sh.get('layers', [])):
        ln = f'{n} layer {j}'; k = l.get('kind', 'img')
        if k not in LK: bad.append(f'{ln}: unknown kind "{k}" (allowed: {sorted(LK)})'); continue
        if k == 'img':
            if l.get('src') and not exists(l['src']): bad.append(f'{ln}: missing file {l["src"]}')
            elif l.get('sprite'): check_sprite(ln, l['sprite'])
            elif not l.get('src'): bad.append(f'{ln}: img needs "src" or "sprite"')
        if k == 'inset' and not (l.get('src') and exists(l['src'])): bad.append(f'{ln}: inset needs an existing "src" plate')
        if l.get('in') and l['in'] not in IN: bad.append(f'{ln}: unknown "in" {l["in"]} (allowed: {sorted(IN)})')
        if l.get('loop') and l['loop'] not in LOOP: bad.append(f'{ln}: unknown "loop" {l["loop"]}')
        if l.get('start', 0) >= d: bad.append(f'{ln}: starts after the shot ends')
        if l.get('attach') is not None and l['attach'] >= len(sh.get('chars', [])): bad.append(f'{ln}: attach index has no character')
        if not l.get('hud') and k not in ('stamp',) and not (-300 <= l.get('x', 0) <= 2220 and -300 <= l.get('y', 0) <= 1380): warn.append(f'{ln}: position {l.get("x")},{l.get("y")} is outside the picture')
    for key in ('tin', 'tout'):
        if sh.get(key) and sh[key] not in TR: bad.append(f'{n}: unknown {key} "{sh[key]}" (allowed: {sorted(TR)})')
    for f in sh.get('fx', []):
        k = f if isinstance(f, str) else f.get('kind')
        if k not in FX: bad.append(f'{n}: unknown fx "{k}"')
    for o in sh.get('overlays', []):
        if o.get('kind') not in OVK: bad.append(f'{n}: unknown overlay kind "{o.get("kind")}"')
        if o.get('start', 0) > d: bad.append(f'{n}: overlay starts after the shot ends')
    for o in sh.get('objects', []):
        if o.get('kind') not in OBJ: bad.append(f'{n}: unknown object "{o.get("kind")}"')
    if ty in ('numbers', 'triptych') and not (sh.get('items') or sh.get('panels') or sh.get('objects') or sh.get('layers')): warn.append(f'{n}: empty {ty} shot')

for n, mv in poor[:12]: warn.append(f'{n}: only {mv:.1f} moving things; a diorama needs 4+ (prop animation, counter, camera key, fx, second layer, character move)')
if len(poor) > 12: warn.append(f'... and {len(poor) - 12} more shots with too little motion')
aud = tl.get('audio')
if aud and os.path.exists(os.path.join(P, aud)):
    r = subprocess.run(['ffprobe', '-v', 'error', '-show_entries', 'format=duration', '-of', 'csv=p=0', os.path.join(P, aud)], capture_output=True, text=True)
    try:
        ad = float(r.stdout.strip())
        if abs(ad - total) > 0.6: bad.append(f'timeline is {total:.1f}s but the voiceover is {ad:.1f}s: run align.py')
    except Exception: pass
elif aud: bad.append(f'voiceover file missing: {aud}')

mins = max(total / 60, 1e-9); srt = sorted(durs); med = srt[len(srt) // 2] if srt else 0
print(f'{len(shots)} shots, {total:.1f} s | cuts/min {len(shots) / mins:.1f} (target 12-21) | median shot {med:.1f}s (target 2.5-4) | typed cards {len(cards)} | unique plates {len(plates)}')
if total > 60:
    if len(shots) / mins < 10: bad.append('too few cuts per minute (rule: a cut or camera move every 3 s): split long lines into two shots')
    if med > 5: bad.append('median shot is too long (target 2.5-4 s): the video would be a slideshow')
    if len(cards) < total / 90: warn.append('few black text cards: add a typed card at each chapter start and for every twist question (about one per 30-60 s)')
    if len(plates) < len(shots) * .4: warn.append('plates are reused too much: the reference uses a new place for almost every sentence')
for w in warn: print('warning:', w)
for b in bad: print('PROBLEM:', b)
print('CHECK OK' if not bad else 'CHECK FAILED'); sys.exit(1 if bad else 0)
