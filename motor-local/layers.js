/* ---------- generic layers: any prop with any animation ----------
   A shot may carry "layers": [ {kind, x, y, ...} ]. Everything a "diorama" needs lives here:
   pictures (props cut out by key_sprites), text, shapes, draw-on lines/arrows, counters, stamps, documents, pins,
   magnifier insets, flip calendars, bar fields and the old objects (cash, fire, balloon, stopwatch, ekg).
   Animation of any layer: "in"/"out" presets, "loop" idle motion, "kf" keyframes, "attach" to a character.
*/
const EASE = {
  lin: x => x, out: x => ease(x), inout: x => x < .5 ? 2 * x * x : 1 - Math.pow(-2 * x + 2, 2) / 2,
  back: x => back(x),
  bounce: x => { const n = 7.5625, d = 2.75; if (x < 1 / d) return n * x * x; if (x < 2 / d) return n * (x -= 1.5 / d) * x + .75; if (x < 2.5 / d) return n * (x -= 2.25 / d) * x + .9375; return n * (x -= 2.625 / d) * x + .984375 },
  elastic: x => x <= 0 ? 0 : x >= 1 ? 1 : Math.pow(2, -10 * x) * Math.sin((x * 10 - .75) * 2.0944) + 1
};
const c01 = x => Math.max(0, Math.min(1, x));
const ANCH = { c: [.5, .5], bc: [.5, 1], bl: [0, 1], br: [1, 1], tl: [0, 0], tc: [.5, 0], tr: [1, 0], cl: [0, .5], cr: [1, .5] };
const IMGDIM = {};
const FONTCLS = { hand: 'hand7', hand6: 'hand', brush: 'brush', bold: 'bold', sans: 'sans' };
const esc = s => String(s).replace(/&/g, '&amp;').replace(/</g, '&lt;');

function preloadShot(sh) {
  const urls = new Set();
  (function walk(o) {
    if (typeof o === 'string') { if (/^(file:|http|data:image)/.test(o) && /(\.(png|jpe?g|webp|gif)(\?|$))|^data:image/i.test(o)) urls.add(o); }
    else if (Array.isArray(o)) o.forEach(walk);
    else if (o && typeof o === 'object') Object.values(o).forEach(walk);
  })(sh);
  return Promise.all([...urls].map(u => new Promise(res => {
    if (IMGDIM[u]) return res();
    const im = new Image(); const done = () => { if (im.naturalWidth) IMGDIM[u] = [im.naturalWidth, im.naturalHeight]; res() };
    im.onload = () => im.decode ? im.decode().then(done, done) : done(); im.onerror = done; im.src = u; setTimeout(res, 8000);
  })));
}

/* ----- characters: path, swap of pictures, auto flip ----- */
function resolveChar(c, t, hor) {
  let o = c;
  if (c.swap) {
    let e = null, pv = null, ix = -1; c.swap.forEach((s, i) => { if (t >= s.t) { e = s; ix = i } });
    if (e) {
      const prevS = ix > 0 ? Object.assign({}, c, c.swap[ix - 1]) : c; o = Object.assign({}, c, e); o._st = e.t;
      o._prevArms = prevS.arms || (typeof PP_POSE !== 'undefined' ? PP_POSE[prevS.pose || 'idle'] : null); delete o.t
    }
  }
  if (c.path && c.path.length) {
    const P = c.path; let i = 0; while (i < P.length - 1 && t >= P[i + 1].t) i++;
    const a = P[i], b = P[Math.min(i + 1, P.length - 1)];
    let p = b.t > a.t ? c01((t - a.t) / (b.t - a.t)) : 1; const e = (EASE[b.ease || 'inout'] || EASE.inout)(p);
    const x = a.x + (b.x - a.x) * e, y = a.y + (b.y - a.y) * e, s = (a.s == null ? (o.s || 1) : a.s) + ((b.s == null ? (o.s || 1) : b.s) - (a.s == null ? (o.s || 1) : a.s)) * e;
    const moving = p > 0 && p < 1 && (Math.abs(b.x - a.x) + Math.abs(b.y - a.y)) > 2;
    o = Object.assign({}, o, { x, y, s });
    if (moving && o.walkPose !== false) o.pose = 'walk';
    if (moving && o.autoflip !== false && Math.abs(b.x - a.x) > 2) o.flip = (b.x < a.x) !== !!c.faceLeft;
  }
  if (hor != null && o.depth !== false) { const k = (o.y - hor) / (1040 - hor); o = Object.assign({}, o, { s: (o.s == null ? 1 : o.s) * Math.max(.25, k) }) }
  return o;
}

/* ----- one layer's animated state ----- */
function layerState(l, t, chars) {
  const a = l.start || 0, b = l.end == null ? 1e9 : l.end; if (t < a) return null;
  const lt = t - a; let x = l.x || 0, y = l.y || 0, s = l.s == null ? 1 : l.s, rot = l.rot || 0, op = l.opacity == null ? 1 : l.opacity, sx = 1;
  if (l.kf) {
    const base = { x, y, s, rot, op };
    ['x', 'y', 's', 'rot', 'op'].forEach(k => {
      const pts = [[0, base[k], 'inout']].concat(l.kf.filter(q => q[k] != null).map(q => [q.t, q[k], q.ease || 'inout'])).sort((u, v) => u[0] - v[0]);
      let v = pts[0][1]; for (let i = 0; i < pts.length - 1; i++) { if (lt >= pts[i][0]) { const p = c01((lt - pts[i][0]) / Math.max(1e-6, pts[i + 1][0] - pts[i][0])); v = pts[i][1] + (pts[i + 1][1] - pts[i][1]) * (EASE[pts[i + 1][2]] || EASE.inout)(p) } }
      if (lt >= pts[pts.length - 1][0]) v = pts[pts.length - 1][1];
      if (k === 'x') x = v; else if (k === 'y') y = v; else if (k === 's') s = v; else if (k === 'rot') rot = v; else op = v;
    });
  }
  const inn = l.in; if (inn) {
    const d = l.inDur || .45, p = c01(lt / d), e = ease(p), bk = back(p);
    if (inn === 'pop') { s *= Math.max(0.001, bk); op *= c01(p * 3) }
    else if (inn === 'slideL') x -= (1 - e) * 1100; else if (inn === 'slideR') x += (1 - e) * 1100;
    else if (inn === 'slideU') y -= (1 - e) * 900; else if (inn === 'slideD') y += (1 - e) * 900;
    else if (inn === 'drop') y -= (1 - EASE.bounce(p)) * 900;
    else if (inn === 'fade') op *= p;
    else if (inn === 'grow') s *= Math.max(0.001, e);
    else if (inn === 'stamp') { s *= 1 + (1 - e) * 1.9; rot += (1 - e) * -12; op *= c01(p * 5) }
    else if (inn === 'swing') rot += (1 - EASE.elastic(p)) * 50;
    else if (inn === 'flipX') { sx = Math.max(0.001, e) }
  }
  if (l.out && t > b) {
    const p = c01((t - b) / (l.outDur || .3)); if (p >= 1) return null; const e = ease(p);
    if (l.out === 'fade') op *= 1 - p; else if (l.out === 'pop') s *= 1 - e; else if (l.out === 'slideL') x -= e * 1100; else if (l.out === 'slideR') x += e * 1100; else if (l.out === 'slideU') y -= e * 900; else if (l.out === 'slideD') y += e * 900;
  } else if (t > b) return null;
  const lp = l.loop; if (lp) {
    if (lp === 'bob') y += 8 * Math.sin(lt * 5); else if (lp === 'shake') { x += 5 * Math.sin(lt * 41); rot += 1.6 * Math.sin(lt * 33) } else if (lp === 'swing') rot += 8 * Math.sin(lt * 2.1);
    else if (lp === 'pulse') s *= 1 + .05 * Math.sin(lt * 5.5); else if (lp === 'spin') rot += lt * 120; else if (lp === 'float') { y += 14 * Math.sin(lt * 1.3); x += 6 * Math.sin(lt * .8) }
    else if (lp === 'flicker') op *= .86 + .14 * Math.sin(lt * 23);
  }
  if (l.attach != null && chars && chars[l.attach]) { x += chars[l.attach].x; y += chars[l.attach].y }
  return { x, y, s, rot, op, sx, lt };
}

/* ----- helpers for shapes ----- */
function partial(pts, p) {
  if (p >= 1) return { pts, ang: pts.length > 1 ? Math.atan2(pts[pts.length - 1][1] - pts[pts.length - 2][1], pts[pts.length - 1][0] - pts[pts.length - 2][0]) : 0 };
  let tot = 0; const L = []; for (let i = 1; i < pts.length; i++) { const d = Math.hypot(pts[i][0] - pts[i - 1][0], pts[i][1] - pts[i - 1][1]); L.push(d); tot += d }
  let rem = tot * p; const out = [pts[0]]; let ang = 0;
  for (let i = 1; i < pts.length; i++) { const d = L[i - 1]; ang = Math.atan2(pts[i][1] - pts[i - 1][1], pts[i][0] - pts[i - 1][0]); if (rem >= d) { out.push(pts[i]); rem -= d } else { const f = d ? rem / d : 0; out.push([pts[i - 1][0] + (pts[i][0] - pts[i - 1][0]) * f, pts[i - 1][1] + (pts[i][1] - pts[i - 1][1]) * f]); break } }
  return { pts: out, ang };
}
function arrowHead(x, y, ang, sz, col) { const a1 = ang + 2.7, a2 = ang - 2.7; return `<path d="M${x.toFixed(1)} ${y.toFixed(1)} L${(x + sz * Math.cos(a1)).toFixed(1)} ${(y + sz * Math.sin(a1)).toFixed(1)} M${x.toFixed(1)} ${y.toFixed(1)} L${(x + sz * Math.cos(a2)).toFixed(1)} ${(y + sz * Math.sin(a2)).toFixed(1)}" stroke="${col}" stroke-width="${Math.max(5, sz / 4)}" stroke-linecap="round" fill="none"/>` }
function pageSVG(x, y, w, h, m, d, cur) { return `<g transform="translate(${x} ${y})"><rect x="${-w / 2}" y="${-h / 2}" width="${w}" height="${h}" rx="8" fill="#fbfaf6" stroke="#222" stroke-width="6"/><rect x="${-w / 2}" y="${-h / 2}" width="${w}" height="${h * .3}" rx="8" fill="${cur || '#d9322e'}" stroke="#222" stroke-width="6"/><text y="${-h / 2 + h * .22}" text-anchor="middle" font-size="${h * .2}" class="hand7" fill="#fff">${esc(m)}</text><text y="${h * .27}" text-anchor="middle" font-size="${h * .5}" class="hand7" fill="#222">${esc(d)}</text></g>` }

/* ----- draw one layer in local coordinates (origin = its x,y) ----- */
function layerBody(l, st, idx, t) {
  const lt = st.lt, k = l.kind || 'img';
  const col = l.color || '#111';
  if (k === 'img') {
    let src = l.src, w = l.w, h = l.h;
    if (l.sprite) { const sp = l.sprite, e = (l.expr && sp[l.expr]) ? sp[l.expr] : sp; src = (l.open && e.open) ? e.open : e.closed; h = h || sp.h || 400; w = w || h * (e.ar || 1) }
    const dm = IMGDIM[src]; if (dm) { if (!w && !h) { w = dm[0]; h = dm[1] } else if (!h) h = w * dm[1] / dm[0]; else if (!w) w = h * dm[0] / dm[1] }
    w = w || 300; h = h || 300; const an = ANCH[l.anchor || 'c'];
    let g = `<image href="${src}" x="${(-w * an[0]).toFixed(1)}" y="${(-h * an[1]).toFixed(1)}" width="${w.toFixed(1)}" height="${h.toFixed(1)}" preserveAspectRatio="none"/>`;
    if (l.sprite && l.ground !== false && (l.anchor === 'bc' || l.ground)) g = `<ellipse cx="${(w * .08).toFixed(1)}" cy="${(h * (1 - an[1]) - 2).toFixed(1)}" rx="${(w * .46).toFixed(1)}" ry="${Math.max(6, h * .045).toFixed(1)}" fill="#000" opacity=".38" style="filter:blur(5px)"/>` + g;
    if (l.wipe) { const p = c01((lt - (l.wipe.at || 0)) / (l.wipe.dur || .6)); const cid = 'wp' + idx; g = `<clipPath id="${cid}"><rect x="${-w * an[0]}" y="${-h * an[1]}" width="${w * p}" height="${h}"/></clipPath><g clip-path="url(#${cid})">${g}</g>` }
    return g;
  }
  if (k === 'text') {
    let str = l.text == null ? '' : String(l.text);
    if (l.counter) { const c = l.counter, pp = ease(c01((lt - (c.at || 0)) / (c.dur || 1.2))); str = fmtNum(c.from + (c.to - c.from) * pp, c) }
    else if (l.typed) { const n = Math.floor(Math.max(0, lt - (l.typeAt || 0)) * (l.cps || 18)) + 1; str = str.slice(0, Math.min(str.length, n)) }
    const size = l.size || 64, lines = str.split('\n'), cls = FONTCLS[l.font || 'hand'] || 'hand7', anc = l.align === 'left' ? 'start' : l.align === 'right' ? 'end' : 'middle';
    const lh = size * 1.18; let s = '';
    if (l.box) { const mx = Math.max(...String(l.text).split('\n').map(z => z.length), 1), pad = l.box.pad == null ? 22 : l.box.pad, bw = mx * size * .5 + pad * 2, bh = lines.length * lh + pad, bx = anc === 'start' ? -pad : anc === 'end' ? -bw + pad : -bw / 2;
      s += `<rect x="${bx}" y="${-size * .85 - pad / 2}" width="${bw}" height="${bh}" rx="${l.box.rx == null ? 8 : l.box.rx}" fill="${l.box.fill || '#f4f3ee'}" stroke="${l.box.stroke || '#111'}" stroke-width="${l.box.sw || 5}"/>` }
    const glow = l.glow ? `filter:drop-shadow(0 0 ${l.glow}px ${col})` : '';
    lines.forEach((z, i) => { s += `<text x="0" y="${(i * lh).toFixed(1)}" text-anchor="${anc}" font-size="${size}" class="${cls}" fill="${col}" ${l.stroke ? `stroke="${l.stroke}" stroke-width="${l.sw || 8}" paint-order="stroke"` : ''} ${l.italic ? 'font-style="italic"' : ''} xml:space="preserve" style="${glow}">${esc(z)}</text>` });
    return s;
  }
  if (k === 'rect' || k === 'circle') {
    const an = ANCH[l.anchor || 'c'], f = l.fill || 'none', sk = l.stroke || '#111', sw = l.sw == null ? 6 : l.sw;
    if (k === 'circle') return `<circle r="${l.r || 60}" fill="${f}" stroke="${sk}" stroke-width="${sw}"/>`;
    const w = l.w || 200, h = l.h || 120; return `<rect x="${-w * an[0]}" y="${-h * an[1]}" width="${w}" height="${h}" rx="${l.rx || 0}" fill="${f}" stroke="${sk}" stroke-width="${sw}"/>`;
  }
  if (k === 'poly') {
    const pts = l.pts || [[0, 0], [100, 0]], p = l.draw ? ease(c01((lt - (l.draw.at || 0)) / (l.draw.dur || 1))) : 1, pp = partial(pts, p);
    const d = pp.pts.map((q, i) => (i ? 'L' : 'M') + q[0].toFixed(1) + ' ' + q[1].toFixed(1)).join(' ') + (l.close && p >= 1 ? ' Z' : '');
    let s = `<path d="${d}" fill="${l.fill || 'none'}" stroke="${l.stroke || '#d9322e'}" stroke-width="${l.sw || 8}" stroke-linecap="round" stroke-linejoin="round" ${l.dash ? `stroke-dasharray="${l.dash}"` : ''}/>`;
    const last = pp.pts[pp.pts.length - 1];
    if (l.head === 'arrow') s += arrowHead(last[0], last[1], pp.ang, l.headSize || 34, l.stroke || '#d9322e');
    if (l.head === 'dot') s += `<circle cx="${last[0]}" cy="${last[1]}" r="${l.headSize || 14}" fill="${l.stroke || '#d9322e'}"/>`;
    return s;
  }
  if (k === 'path') {
    const p = l.draw ? ease(c01((lt - (l.draw.at || 0)) / (l.draw.dur || 1))) : 1;
    return `<path d="${l.d}" pathLength="1" fill="${l.fill || 'none'}" stroke="${l.stroke || '#d9322e'}" stroke-width="${l.sw || 8}" stroke-linecap="round" stroke-dasharray="1 1" stroke-dashoffset="${(1 - p).toFixed(4)}"/>`;
  }
  if (k === 'stamp') {
    const size = l.size || 64, txt = String(l.text || 'STAMP'), w = Math.max(160, txt.length * size * .56 + 50), h = size * 1.5, c = l.color || '#c92b27';
    return `<g opacity=".97"><rect x="${-w / 2}" y="${-h / 2}" width="${w}" height="${h}" rx="10" fill="${l.fill || 'rgba(250,246,232,.93)'}" stroke="${c}" stroke-width="9"/><rect x="${-w / 2 + 11}" y="${-h / 2 + 11}" width="${w - 22}" height="${h - 22}" rx="6" fill="none" stroke="${c}" stroke-width="3"/><text y="${size * .34}" text-anchor="middle" font-size="${size}" class="bold" fill="${c}">${esc(txt)}</text></g>`;
  }
  if (k === 'doc') {
    const w = l.w || 760, size = l.size || 38, rh = size * 1.55, rows = l.rows || [], hdr = l.header, hh = hdr ? 52 : 0, pad = 26, av = l.avatar ? 110 : 0;
    const h = hh + pad * 2 + rows.length * rh + 10;
    let s = `<g transform="translate(${-w / 2} ${-h / 2})"><rect width="${w}" height="${h}" rx="10" fill="${l.paper || '#f7f6f2'}" stroke="#222" stroke-width="6"/>`;
    if (hdr) s += `<rect width="${w}" height="${hh}" rx="10" fill="${l.hcolor || '#1b2a41'}" stroke="#222" stroke-width="6"/><text x="22" y="${hh * .68}" font-size="26" class="sans" fill="#fff" font-style="italic">${esc(hdr)}</text>`;
    if (l.avatar) s += `<rect x="${pad}" y="${hh + pad}" width="86" height="86" rx="6" fill="#1b2a41"/><circle cx="${pad + 43}" cy="${hh + pad + 34}" r="18" fill="#9fb0c8"/><path d="M${pad + 10} ${hh + pad + 86} Q${pad + 43} ${hh + pad + 40} ${pad + 76} ${hh + pad + 86} Z" fill="#9fb0c8"/>`;
    let tcur = (l.typeAt || .2);
    rows.forEach((r, i) => {
      const y = hh + pad + (i + .78) * rh, x0 = pad + av, rw = w - pad * 2 - av;
      if (r.gray) { s += `<rect x="${x0}" y="${y - size * .35}" width="${rw * (r.w || .8)}" height="12" rx="6" fill="#c9c9c4"/>`; return }
      const txt = String(r.text || ''), cps = l.cps || 26, dur = txt.length / cps, st0 = r.at == null ? tcur : r.at; tcur = st0 + dur + .15;
      const n = Math.floor(Math.max(0, lt - st0) * cps) + 1, shown = txt.slice(0, Math.min(txt.length, n)), fin = lt >= st0 + dur;
      if (r.hl && lt >= st0) s += `<rect x="${x0 - 6}" y="${y - size * .8}" width="${Math.min(rw, shown.length * size * .5 + 12)}" height="${size * 1.1}" fill="#ffe66b" opacity=".85" rx="4"/>`;
      s += `<text x="${x0}" y="${y}" font-size="${size}" class="${r.font === 'sans' ? 'sans' : 'hand'}" fill="${r.color || '#222'}" xml:space="preserve">${esc(shown)}</text>`;
      if (r.circle && fin) { const p = ease(c01((lt - (st0 + dur)) / .5)), ew = Math.min(rw, txt.length * size * .5) / 2 + 16; s += `<ellipse cx="${x0 + ew - 8}" cy="${y - size * .3}" rx="${ew}" ry="${size * .78}" fill="none" stroke="#d9322e" stroke-width="6" pathLength="1" stroke-dasharray="1 1" stroke-dashoffset="${(1 - p).toFixed(3)}" transform="rotate(-3 ${x0} ${y})"/>` }
    });
    return s + '</g>';
  }
  if (k === 'pin') {
    const p = c01(lt / .5), dy = -(1 - EASE.bounce(p)) * 260, c = l.color || '#d9322e';
    let s = `<g transform="translate(0 ${dy})"><path d="M0 0 C-30 -40 -34 -68 -34 -80 A34 34 0 1 1 34 -80 C34 -68 30 -40 0 0 Z" fill="${c}" stroke="#111" stroke-width="5"/><circle cy="-80" r="13" fill="#fff"/></g>`;
    if (l.text && lt > .35) { const w = String(l.text).length * 17 + 36; s += `<g transform="translate(0 -150) scale(${back(c01((lt - .35) / .3))})"><rect x="${-w / 2}" y="-34" width="${w}" height="56" rx="8" fill="#fff" stroke="#111" stroke-width="4"/><text y="8" text-anchor="middle" font-size="32" class="hand7" fill="#111">${esc(l.text)}</text></g>` }
    return s;
  }
  if (k === 'inset') {
    const r = l.r || 190, z = l.zoom || 3, fx = l.fx == null ? .5 : l.fx, fy = l.fy == null ? .5 : l.fy, cid = 'in' + idx, iw = 1920 * z, ih = 1080 * z;
    return `<clipPath id="${cid}"><circle r="${r}"/></clipPath><circle r="${r + 8}" fill="#fff"/><g clip-path="url(#${cid})"><image href="${l.src}" x="${-fx * iw}" y="${-fy * ih}" width="${iw}" height="${ih}" preserveAspectRatio="none"/></g><circle r="${r}" fill="none" stroke="#111" stroke-width="7" stroke-dasharray="${l.dashed === false ? '' : '16 10'}"/>`;
  }
  if (k === 'calendar') {
    const steps = l.steps || [{ at: 0, m: 'JAN', d: '1' }], w = l.w || 250, h = l.h || 280; let i = 0; steps.forEach((q, j) => { if (lt >= q.at) i = j });
    const cur = steps[i], prev = steps[i - 1], fp = c01((lt - cur.at) / .22); let s = '';
    if (prev && fp < 1) { s += pageSVG(0, 0, w, h, cur.m, cur.d, l.color) ; s += `<g transform="translate(0 ${-h / 2}) scale(1 ${Math.max(0.02, 1 - fp)}) translate(0 ${h / 2})">${pageSVG(0, 0, w, h, prev.m, prev.d, l.color)}</g>` }
    else s += pageSVG(0, 0, w, h, cur.m, cur.d, l.color);
    if (l.circleAt != null && lt >= l.circleAt) { const p = ease(c01((lt - l.circleAt) / .5)); s += `<ellipse cx="0" cy="${h * .22}" rx="${w * .34}" ry="${h * .3}" fill="none" stroke="#d9322e" stroke-width="9" pathLength="1" stroke-dasharray="1 1" stroke-dashoffset="${(1 - p).toFixed(3)}" transform="rotate(-8 0 ${h * .22})"/>` }
    return s;
  }
  if (k === 'bars') {
    const n = l.n || 120, W = l.w || 1300, H = l.h || 260, r = rnd(l.seed || 3), cols = l.colors || ['#a9bad0', '#8da4c0'], ap = l.appear || 1.2, sh = l.shrink; let s = '';
    const sc = sh ? 1 - (1 - (sh.to == null ? .1 : sh.to)) * ease(c01((lt - (sh.at || 0)) / (sh.dur || 1.5))) : 1, bw = W / n;
    for (let i = 0; i < n; i++) { const hh = (H * .45 + r() * H * .55) * sc * ease(c01((lt - (i / n) * ap) / .3)); s += `<rect x="${(i * bw).toFixed(1)}" y="${(-hh).toFixed(1)}" width="${(bw * 1.15).toFixed(1)}" height="${hh.toFixed(1)}" fill="${cols[i % cols.length]}" stroke="#3b4a60" stroke-width="2"/>` }
    return `<g transform="translate(${-W / 2} 0)">${s}</g>`;
  }
  if (k === 'cash') return cashSVG(0, 0, 1, l.rows || 2);
  if (k === 'fire') return fireSVG(lt, 0, 0, 1);
  if (k === 'balloon') return balloonSVG(lt, 0, 0, 1);
  if (k === 'stopwatch') return stopSVG(lt, Object.assign({}, l, { x: 0, y: 0, s: 1, start: l.at0 || 0 }));
  if (k === 'ekg') return ekgSVG(lt, Object.assign({}, l, { x: 0, y: 0, s: 1 }));
  return '';
}

function layersSVG(list, t, chars, filt) {
  let out = '';
  (list || []).forEach((l, i) => {
    if (filt && !filt(l)) return;
    const st = layerState(l, t, chars); if (!st || st.op <= 0.003) return;
    const body = layerBody(l, st, i, t); if (!body) return;
    const flip = (l.flip ? -1 : 1) * st.sx;
    const fl = [(l.sprite && l.kind !== 'text' && l.lit !== false && typeof LIGHTSTYLE !== 'undefined' && LIGHTSTYLE) ? LIGHTSTYLE : '', l.glow && l.kind !== 'text' ? `drop-shadow(0 0 ${l.glow}px ${l.glowColor || '#8bff3a'})` : '', l.shadow ? 'drop-shadow(0 10px 12px rgba(0,0,0,.45))' : '', l.filter || ''].filter(Boolean).join(' ');
    out += `<g transform="translate(${st.x.toFixed(1)} ${st.y.toFixed(1)}) rotate(${st.rot.toFixed(2)}) scale(${(st.s * flip).toFixed(4)} ${st.s.toFixed(4)})" opacity="${st.op.toFixed(3)}" ${fl ? `style="filter:${fl}"` : ''}>${body}</g>`;
  });
  return out;
}

/* ----- colour grading: one mood per shot (multiply wash, glows) ----- */
function gradeSVG(sh, t) {
  const g = sh.grade; if (!g) return ''; const list = Array.isArray(g) ? g : [g]; let s = '';
  list.forEach((q, i) => {
    const a0 = q.a0 != null ? q.a0 : (q.alpha == null ? .25 : q.alpha), a1 = q.a1 != null ? q.a1 : a0;
    const p = q.dur ? c01((t - (q.at || 0)) / q.dur) : c01(t / sh.duration), a = a0 + (a1 - a0) * p, bl = q.blend || (q.radial ? 'screen' : 'multiply');
    if (q.radial) { const id = 'gr' + i; s += `<radialGradient id="${id}" cx="${q.radial[0]}" cy="${q.radial[1]}" r="${q.radial[2]}" gradientUnits="userSpaceOnUse"><stop offset="0" stop-color="${q.color}" stop-opacity="1"/><stop offset="1" stop-color="${q.color}" stop-opacity="0"/></radialGradient><rect width="1920" height="1080" fill="url(#${id})" opacity="${a.toFixed(3)}" style="mix-blend-mode:${bl}"/>` }
    else s += `<rect width="1920" height="1080" fill="${q.color}" opacity="${a.toFixed(3)}" style="mix-blend-mode:${bl}"/>`;
  });
  return s;
}

/* ----- camera keyframes ----- */
var VIEW_HW = 960; // half width of the visible window (set by setFormat for 9:16 and 1:1)
function camState(sh, t) {
  const K = sh.cam; if (!K || !K.length) return null;
  let i = 0; while (i < K.length - 1 && t >= K[i + 1].t) i++;
  const a = K[i], b = K[Math.min(i + 1, K.length - 1)], p = b.t > a.t ? c01((t - a.t) / (b.t - a.t)) : 1, e = (EASE[b.ease || 'inout'] || EASE.inout)(p), g = (u, v, d) => (u == null ? d : u) + ((v == null ? d : v) - (u == null ? d : u)) * e;
  // keep the view inside the picture: never zoom out past 1 and never pan past the plate edge (no black bands)
  const z = Math.max(1, g(a.z, b.z, 1)), hx = VIEW_HW / z, hy = 540 / z;
  return { x: Math.min(1920 - hx, Math.max(hx, g(a.x, b.x, 960))), y: Math.min(1080 - hy, Math.max(hy, g(a.y, b.y, 540))), z };
}
function shakeOf(sh, t) {
  let dx = 0, dy = 0; (sh.shake || []).forEach(q => { const p = (t - q.at) / (q.dur || .35); if (p >= 0 && p < 1) { const a = (q.amp || 14) * (1 - p); dx += a * Math.sin(t * 61); dy += a * Math.cos(t * 53) } }); return [dx, dy];
}

/* ----- shot-to-shot transitions applied to the whole picture ----- */
function transitionStyle(sh, t) {
  const d = sh.duration; let tx = 0, ty = 0, k = 1, blur = 0, fog = 0, op = 1;
  const apply = (name, p, isIn) => { // p: 0..1 progress of the transition (0 = fully "away")
    const q = 1 - ease(p);
    if (name === 'whipL') { tx += (isIn ? -1 : 1) * q * 700; blur = Math.max(blur, 22 * q) }
    else if (name === 'whipR') { tx += (isIn ? 1 : -1) * q * 700; blur = Math.max(blur, 22 * q) }
    else if (name === 'whipU') { ty += (isIn ? -1 : 1) * q * 420; blur = Math.max(blur, 22 * q) }
    else if (name === 'whipD') { ty += (isIn ? 1 : -1) * q * 420; blur = Math.max(blur, 22 * q) }
    else if (name === 'fog') { fog = Math.max(fog, .96 * q); blur = Math.max(blur, 8 * q) }
    else if (name === 'zoom') { k *= 1 + q * .35; blur = Math.max(blur, 14 * q) }
    else if (name === 'focus') { blur = Math.max(blur, 22 * q) }
    else if (name === 'fade') { op = Math.min(op, 1 - q) }
  };
  const di = sh.tinDur || .3, dO = sh.toutDur || .3;
  if (sh.tin && t < di) apply(sh.tin, t / di, true);
  if (sh.tout && t > d - dO) apply(sh.tout, (d - t) / dO, false);
  return { tx, ty, k, blur, fog, op };
}
