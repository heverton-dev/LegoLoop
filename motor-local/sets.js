/* RobbiSkill v2: code-drawn sets (backgrounds) with their own light sources.
   Shot field:  "set":{"kind":"office","time":"night","calendar":"NOV 17","sign":"CAFE"}   (replaces "bg")
   Each set returns {svg, lights:[{x,y,r,a}], glows:[{x,y,r,c,a}], cones:[{x,y,w0,w1,len,lean,a}]}; the engine darkens the picture and cuts the lights out of it. */
const SINK = '#18161B';
function sR(x, y, w, h, f, o) { o = o || {}; return `<rect x="${x}" y="${y}" width="${w}" height="${h}" rx="${o.r || 0}" fill="${f}" ${o.sw === 0 ? '' : `stroke="${SINK}" stroke-width="${o.sw || 5}"`} ${o.op ? `opacity="${o.op}"` : ''}/>`; }
function sC(x, y, r, f, o) { o = o || {}; return `<circle cx="${x}" cy="${y}" r="${r}" fill="${f}" ${o.sw === 0 ? '' : `stroke="${SINK}" stroke-width="${o.sw || 5}"`}/>`; }
function sP(d, f, o) { o = o || {}; return `<path d="${d}" fill="${f}" ${o.sw === 0 ? '' : `stroke="${SINK}" stroke-width="${o.sw || 5}" stroke-linejoin="round"`}/>`; }
function sL(x1, y1, x2, y2, w, c) { return `<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}" stroke="${c || SINK}" stroke-width="${w || 5}" stroke-linecap="round"/>`; }
function sT(x, y, s, size, f, o) { o = o || {}; return `<text x="${x}" y="${y}" font-size="${size}" class="${o.cls || 'brush'}" fill="${f}" text-anchor="${o.a || 'middle'}" ${o.rot ? `transform="rotate(${o.rot} ${x} ${y})"` : ''} ${o.st ? `stroke="${SINK}" stroke-width="${o.st}" paint-order="stroke"` : ''}>${s}</text>`; }
function sRng(seed) { let s = (seed >>> 0) || 1; return () => { s = (s * 1664525 + 1013904223) >>> 0; return s / 4294967296 } }
function sky(time, h) {
  const c = { night: ['#0d1230', '#2b3573'], dusk: ['#3b2a6b', '#e8806a'], day: ['#7ec4ee', '#cfeaf7'], red: ['#2a0a12', '#8a1f2b'] }[time] || ['#0d1230', '#2b3573'];
  return `<defs><linearGradient id="sky${time}" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="${c[0]}"/><stop offset="1" stop-color="${c[1]}"/></linearGradient></defs><rect width="1920" height="${h}" fill="url(#sky${time})"/>`;
}
function stars(seed, n, h) { const r = sRng(seed); let s = ''; for (let i = 0; i < n; i++) s += `<circle cx="${(r() * 1920).toFixed(0)}" cy="${(r() * h).toFixed(0)}" r="${(0.8 + r() * 1.8).toFixed(1)}" fill="#fff" opacity="${(0.4 + r() * .5).toFixed(2)}"/>`; return s; }
function moon(x, y, r) { return `<circle cx="${x}" cy="${y}" r="${r * 2.2}" fill="#dfe6ff" opacity=".10"/><circle cx="${x}" cy="${y}" r="${r * 1.5}" fill="#dfe6ff" opacity=".16"/>${sC(x, y, r, '#F4F0D6', { sw: 4 })}<circle cx="${x - r * .3}" cy="${y - r * .2}" r="${r * .18}" fill="#E0DBB8"/><circle cx="${x + r * .35}" cy="${y + r * .3}" r="${r * .12}" fill="#E0DBB8"/>`; }
/* a row of buildings with lit windows; returns svg + the lit window centres (for glow) */
function skyline(x0, x1, base, hmin, hmax, seed, o) {
  o = o || {}; const r = sRng(seed); let s = '', x = x0; const cols = o.cols || ['#1d2548', '#232c55', '#19203f', '#2a3260'], wins = [];
  while (x < x1) {
    const w = 90 + r() * 120, h = hmin + r() * (hmax - hmin), c = cols[Math.floor(r() * cols.length)];
    s += sR(x, base - h, w, h, c, { sw: o.sw == null ? 4 : o.sw });
    const nx = Math.floor((w - 20) / 34), ny = Math.floor((h - 30) / 46);
    for (let i = 0; i < nx; i++) for (let j = 0; j < ny; j++) {
      if (r() < (o.lit == null ? .45 : o.lit)) { const wx = x + 14 + i * 34, wy = base - h + 22 + j * 46; s += `<rect x="${wx}" y="${wy}" width="20" height="28" fill="${r() < .75 ? '#F6C65A' : '#9FD6F5'}" opacity=".95"/>`; if (wins.length < 40 && r() < .2) wins.push([wx + 10, wy + 14]) }
    }
    x += w + (o.gap == null ? 6 : o.gap);
  }
  return { svg: s, wins };
}
function winGlow(w, c, a, r) { return w.map(p => ({ x: p[0], y: p[1], r: r || 52, c: c || '#FFC85A', a: a || .16 })) }

/* ---------------- sets ---------------- */
const SETS = {
  office(o) {
    const night = (o.time || 'night') !== 'day', cal = (o.calendar || 'MAY 1925').split(' ');
    let s = sR(0, 0, 1920, 1080, '#2d5a52', { sw: 0 });
    for (let i = 0; i < 24; i++) s += `<rect x="${i * 80}" y="0" width="40" height="760" fill="#356a60" opacity=".55"/>`;
    s += sR(0, 640, 1920, 120, '#6b4a30', { sw: 0 }) + sL(0, 640, 1920, 640, 6) + sR(0, 760, 1920, 320, '#4a3424', { sw: 0 }) + sL(0, 760, 1920, 760, 6);
    for (let i = 0; i < 13; i++) s += sL(i * 160, 760, i * 160 - 120 + i * 20, 1080, 3, '#2c1f16');
    // window
    s += sR(1220, 150, 460, 470, '#0d1230', { sw: 8 });
    const cityW = skyline(1226, 1676, 616, 90, 300, 11, { lit: .5 }); s += `<clipPath id="ow"><rect x="1224" y="154" width="452" height="462"/></clipPath><g clip-path="url(#ow)">${sky(night ? 'night' : 'day', 616)}${night ? stars(3, 24, 260) : ''}${night ? moon(1590, 230, 32) : ''}${cityW.svg}</g>`;
    s += sL(1450, 150, 1450, 620, 8) + sL(1220, 385, 1680, 385, 8) + sR(1208, 614, 484, 22, '#6b4a30');
    s += sP('M1180 120 L1260 120 L1230 640 L1170 640 Z', '#8a1f26') + sP('M1740 120 L1660 120 L1690 640 L1750 640 Z', '#8a1f26') + sL(1150, 120, 1780, 120, 10, '#6b4a30');
    // clock, calendar, picture
    s += sC(560, 290, 70, '#F4EFE2', { sw: 7 }) + sL(560, 290, 560, 240, 6) + sL(560, 290, 598, 304, 6) + sC(560, 290, 6, SINK, { sw: 0 });
    const calDia = cal[1] || '', calTam = calDia.length > 2 ? Math.round(64 * 2 / calDia.length) : 64; // AC-E12.1: evita a data transbordar o cartao (achado do usuario, 2026-11-26)
    s += sR(900, 190, 150, 190, '#F4EFE2', { r: 6 }) + sR(900, 190, 150, 50, '#B3283A', { r: 6 }) + sT(975, 228, cal[0], 36, '#fff', { cls: 'sans' }) + sT(975, 330, calDia, calTam, SINK, { cls: 'brush' });
    s += sR(230, 200, 190, 150, '#d9b66a', { sw: 8 }) + sR(246, 216, 158, 118, '#6ea6b8', { sw: 4 }) + sP('M246 334 L300 270 L340 310 L372 280 L404 334 Z', '#4f7f62', { sw: 4 });
    // desk
    s += sR(120, 640, 700, 60, '#7c5233', { r: 8, sw: 6 }) + sR(150, 700, 640, 190, '#6b4627', { r: 6, sw: 6 });
    [0, 1, 2].forEach(i => s += sR(170 + i * 205, 720, 185, 70, '#7f5836', { r: 4 }) + sC(262 + i * 205, 755, 8, '#d9b66a'));
    s += sR(180, 640, 260, 36, '#e8e2cc', { sw: 4 }) + sR(480, 610, 120, 30, '#2b2b33', { r: 6 }) + sC(560, 600, 12, '#2b2b33', { sw: 0 });
    // banker lamp
    s += sR(650, 600, 60, 40, '#b8962e', { r: 6 }) + sL(680, 600, 680, 520, 10) + sP('M590 530 Q680 470 770 530 L750 560 L610 560 Z', '#2F7A52', { sw: 6 });
    // extras: plant, filing cabinet, phone, paper stack, waste bin
    s += sR(1010, 640, 70, 140, '#8a6a4a', { r: 6 }) + sC(1045, 600, 58, '#2f7a52') + sC(1010, 640, 38, '#3a9a62') + sC(1082, 644, 36, '#3a9a62');
    s += sR(40, 560, 150, 280, '#5a6270', { r: 6 }) + [0, 1, 2].map(i => sR(52, 576 + i * 88, 126, 70, '#6c7482', { r: 4 }) + sR(100, 602 + i * 88, 30, 8, '#c9c9d6', { r: 3, sw: 3 })).join('');
    s += sR(470, 626, 110, 30, '#e8e2cc', { sw: 4 }) + sR(496, 604, 70, 26, '#cfc8ae', { sw: 4 }) + sR(240, 618, 90, 16, '#d8d2ba', { sw: 4 });
    s += sR(1700, 760, 60, 90, '#6b6b78', { r: 8 });
    // rug
    s += `<ellipse cx="960" cy="960" rx="640" ry="92" fill="#7E1C24" stroke="${SINK}" stroke-width="6"/><ellipse cx="960" cy="960" rx="560" ry="72" fill="none" stroke="#c9a04a" stroke-width="5"/>`;
    return { svg: s, floorY: 1010, lights: [{ x: 680, y: 640, r: 560, a: .95 }, { x: 1450, y: 400, r: 420, a: .55 }, { x: 960, y: 560, r: 760, a: .22 }], glows: [{ x: 680, y: 585, r: 260, c: '#FFD27A', a: .35 }, ...winGlow(cityW.wins, '#FFC85A', .12, 40)], cones: [{ x: 680, y: 560, w0: 120, w1: 520, len: 420, lean: 0, a: .16, c: '#FFD27A' }] };
  },
  street(o) {
    const t = o.time || 'night'; let s = sky(t === 'day' ? 'day' : 'night', 700) + (t === 'day' ? '' : stars(5, 40, 330) + moon(260, 170, 56));
    const far = skyline(-20, 1940, 700, 160, 360, 21, { cols: ['#161c3a', '#1a2142', '#131833'], lit: .35, sw: 0 }); s += `<g opacity=".7" style="filter:blur(3px)">${far.svg}</g>`;
    s += sR(0, 690, 1920, 110, '#2c2a40', { sw: 0 });
    const L = skyline(-30, 520, 790, 330, 620, 31, { cols: ['#3a2c4a', '#4a3a3a', '#2f3a52'], lit: .5 }), Rr = skyline(1500, 1950, 790, 300, 560, 41, { cols: ['#2f3a52', '#3a2c4a', '#42303a'], lit: .5 });
    s += L.svg + Rr.svg;
    s += sR(560, 560, 640, 230, '#6a3a2c', { sw: 6 }) + sP('M540 560 L1220 560 L1190 500 L570 500 Z', '#B3283A', { sw: 6 });
    for (let i = 0; i < 8; i++) s += `<rect x="${570 + i * 80}" y="502" width="40" height="56" fill="#F4EFE2" opacity=".9"/>`;
    s += sR(610, 600, 230, 170, '#F7D88A', { sw: 6 }) + sR(900, 600, 230, 170, '#F7D88A', { sw: 6 }) + sT(880, 548, o.sign || 'CAFE', 46, '#fff', { cls: 'brush', st: 8 });
    s += sR(0, 790, 1920, 290, '#34344a', { sw: 0 }) + sL(0, 790, 1920, 790, 6) + sR(0, 790, 1920, 70, '#5a5a70', { sw: 0 }) + sL(0, 860, 1920, 860, 5);
    for (let i = 0; i < 20; i++) s += sL(i * 110, 860, i * 110 - 60 + i * 8, 1080, 3, '#26263a');
    // extras: parked taxi, trash can, bus-stop sign, crosswalk stripes, pigeons
    s += sR(1320, 705, 300, 90, '#F2C230', { r: 18 }) + sR(1370, 660, 190, 60, '#F2C230', { r: 16 }) + sR(1388, 670, 70, 40, '#9fd6f5', { r: 6 }) + sR(1470, 670, 70, 40, '#9fd6f5', { r: 6 }) + sC(1380, 800, 28, '#1c1c24') + sC(1560, 800, 28, '#1c1c24') + sR(1425, 640, 70, 20, '#fff', { r: 4 });
    s += sR(40, 722, 54, 76, '#4f6a58', { r: 6 }) + sR(32, 712, 70, 16, '#3a5043', { r: 4 });
    s += sR(300, 600, 12, 200, '#2a2a33') + sR(250, 580, 112, 54, '#2b5aa8', { r: 6 }) + sT(306, 618, 'BUS', 36, '#fff', { cls: 'sans' });
    for (let i = 0; i < 9; i++) s += `<rect x="${200 + i * 160}" y="930" width="90" height="24" fill="#c9c9d6" opacity=".35"/>`;
    [[980, 886], [1034, 902], [1010, 940]].forEach(([x, y]) => s += `<ellipse cx="${x}" cy="${y}" rx="17" ry="11" fill="#8a8a99" stroke="${SINK}" stroke-width="3"/><circle cx="${x + 12}" cy="${y - 8}" r="7" fill="#8a8a99" stroke="${SINK}" stroke-width="3"/><path d="M${x + 19} ${y - 8} l9 3 l-9 2 z" fill="#e8a23a"/>`);
    // street lamp
    s += sR(1318, 280, 16, 520, '#2a2a33', { sw: 5 }) + sP('M1326 280 Q1326 250 1250 250', 'none', { sw: 9 }) + sP('M1210 250 L1290 250 L1272 286 L1228 286 Z', '#F7D88A', { sw: 5 });
    return { svg: s, floorY: 1000, lights: [{ x: 1250, y: 560, r: 560, a: .95 }, { x: 880, y: 700, r: 420, a: .8 }, { x: 960, y: 700, r: 1100, a: .2 }], glows: [{ x: 1250, y: 270, r: 200, c: '#FFD27A', a: .5 }, { x: 880, y: 700, r: 280, c: '#FFD27A', a: .25 }, ...winGlow(L.wins.concat(Rr.wins), '#FFC85A', .12, 40)], cones: [{ x: 1250, y: 288, w0: 60, w1: 560, len: 560, lean: 0, a: .22, c: '#FFD27A' }] };
  },
  city(o) {
    const t = o.time || 'night'; let s = sky(t === 'day' ? 'day' : t === 'dusk' ? 'dusk' : 'night', 800) + (t === 'night' ? stars(9, 70, 380) : '') + (t === 'day' ? '' : moon(1480, 190, 62));
    const far = skyline(-20, 1940, 760, 120, 330, 5, { cols: ['#1a2145', '#1f2850'], lit: .35, sw: 0 }); s += `<g opacity=".8">${far.svg}</g>`;
    if (o.landmark === 'tower') s += sP('M960 140 L992 330 L1030 560 L1096 780 L824 780 L890 560 L928 330 Z', '#7a5a3a', { sw: 6 }) + sP('M900 520 Q960 470 1020 520', 'none', { sw: 8 }) + sP('M850 700 Q960 640 1070 700', 'none', { sw: 8 }) + sL(960, 140, 960, 90, 6) + sC(960, 84, 8, '#F7D88A', { sw: 0 });
    else if (o.landmark === 'skyscraper') s += sR(840, 150, 240, 640, '#2a3568', { sw: 6 }) + sR(900, 100, 120, 60, '#2a3568', { sw: 6 }) + sL(960, 100, 960, 50, 6);
    const near = skyline(-20, 1940, 800, 200, 460, 17, { cols: ['#222a52', '#2a3260', '#1c2346'], lit: .5 }); s += near.svg;
    s += sR(0, 800, 1920, 280, '#101a3a', { sw: 0 }) + sL(0, 800, 1920, 800, 5);
    for (let i = 0; i < 30; i++) s += `<rect x="${i * 70}" y="${820 + (i % 4) * 40}" width="${40 + (i % 3) * 20}" height="6" fill="#F7D88A" opacity=".28"/>`;
    return { svg: s, floorY: 1010, lights: [{ x: 960, y: 700, r: 900, a: .5 }, { x: 1480, y: 190, r: 380, a: .5 }], glows: [{ x: 1480, y: 190, r: 300, c: '#CFD8FF', a: .25 }, ...winGlow(near.wins, '#FFC85A', .14, 44)], cones: [] };
  },
  boardroom(o) {
    let s = sR(0, 0, 1920, 1080, '#262a3a', { sw: 0 }); for (let i = 0; i < 12; i++) s += `<rect x="${i * 170}" y="0" width="12" height="640" fill="#2f3446"/>`;
    s += sR(0, 640, 1920, 440, '#3a2a22', { sw: 0 }) + sL(0, 640, 1920, 640, 6);
    s += sR(260, 120, 1400, 440, '#0d1230', { sw: 8 }); const w = skyline(266, 1654, 556, 120, 330, 8, { lit: .5 }); s += `<clipPath id="bw"><rect x="264" y="124" width="1392" height="432"/></clipPath><g clip-path="url(#bw)">${sky('night', 556)}${stars(4, 30, 200)}${w.svg}</g>`;
    for (let i = 1; i < 4; i++) s += sL(260 + i * 350, 120, 260 + i * 350, 560, 8);
    s += sP('M180 860 L1740 860 L1560 700 L360 700 Z', '#5a3a28', { sw: 7 }) + sR(180, 860, 1560, 50, '#43291b', { sw: 7 });
    for (let i = 0; i < 6; i++) { const x = 440 + i * 210; s += sR(x - 55, 640, 110, 90, '#1b1d28', { r: 22 }) + sR(x - 55, 690, 110, 20, '#2a2d3c', { sw: 4 }) }
    [520, 960, 1400].forEach(x => s += sL(x, 0, x, 60, 4) + sP(`M${x - 70} 130 L${x + 70} 130 L${x + 40} 60 L${x - 40} 60 Z`, '#d9b66a', { sw: 5 }));
    return { svg: s, floorY: 1010, lights: [{ x: 520, y: 400, r: 420, a: .7 }, { x: 960, y: 400, r: 480, a: .75 }, { x: 1400, y: 400, r: 420, a: .7 }, { x: 960, y: 840, r: 900, a: .3 }], glows: [{ x: 960, y: 130, r: 400, c: '#FFD27A', a: .3 }, ...winGlow(w.wins, '#FFC85A', .12, 40)], cones: [520, 960, 1400].map(x => ({ x, y: 128, w0: 120, w1: 520, len: 560, lean: 0, a: .14, c: '#FFD27A' })) };
  },
  stage(o) {
    const sc = o.spot || '#FFD27A'; let s = sR(0, 0, 1920, 1080, '#2a1f4a', { sw: 0 });
    s += sP('M0 0 L330 0 Q260 400 330 760 L0 760 Z', '#6a1f3a', { sw: 6 }) + sP('M1920 0 L1590 0 Q1660 400 1590 760 L1920 760 Z', '#6a1f3a', { sw: 6 });
    for (let i = 0; i < 6; i++) s += sL(i * 56, 0, i * 52 + 20, 760, 3, '#4a132a');
    s += sR(470, 150, 980, 520, '#14173a', { r: 8, sw: 8 }) + `<rect x="486" y="166" width="948" height="488" fill="${o.screen || '#1f2a66'}" opacity=".9"/>`;
    s += sR(0, 760, 1920, 320, '#3a2a66', { sw: 0 }) + sL(0, 760, 1920, 760, 6);
    for (let i = 0; i < 10; i++) s += sL(i * 210, 760, i * 210 - 160, 1080, 3, '#2b1f4e');
    for (let i = 0; i < 9; i++) s += sC(120 + i * 210, 60 + (i % 2) * 14, 12, '#fff', { sw: 3 });
    return { svg: s, floorY: 1000, lights: [{ x: o.spotX || 960, y: 880, r: 520, a: .95 }, { x: 960, y: 400, r: 520, a: .5 }, { x: 960, y: 600, r: 1100, a: .2 }], glows: [{ x: 960, y: 400, r: 460, c: '#5878ff', a: .22 }], cones: [{ x: o.spotX || 960, y: 0, w0: 80, w1: 560, len: 900, lean: 0, a: .24, c: sc }] };
  },
  trading(o) {
    const red = (o.mood || 'red') === 'red'; let s = sR(0, 0, 1920, 1080, '#14182e', { sw: 0 }); const r = sRng(12), gl = [];
    for (let j = 0; j < 4; j++) for (let i = 0; i < 9; i++) {
      const x = 90 + i * 200, y = 80 + j * 150, up = r() < (red ? .12 : .8); s += sR(x, y, 176, 124, '#0b0e1c', { r: 6, sw: 5 }) + `<rect x="${x + 8}" y="${y + 8}" width="160" height="108" fill="${up ? '#123d2a' : '#4a1018'}"/>`;
      let d = `M${x + 14} ${y + 90}`; for (let k = 1; k <= 8; k++) d += ` L${x + 14 + k * 18} ${y + 90 - (up ? k * 8 : -k * 2) + (r() * 20 - 10)}`; s += `<path d="${d}" fill="none" stroke="${up ? '#3dff9a' : '#ff4a5a'}" stroke-width="5"/>`; gl.push({ x: x + 88, y: y + 62, r: 120, c: up ? '#3dff9a' : '#ff4a5a', a: .10 });
    }
    s += sR(0, 700, 1920, 380, '#1e2236', { sw: 0 }) + sL(0, 700, 1920, 700, 6);
    for (let i = 0; i < 5; i++) { const x = 100 + i * 380; s += sR(x, 780, 330, 40, '#3a3f58', { r: 6 }) + sR(x + 40, 740, 120, 40, '#0b0e1c', { r: 4 }) + sR(x + 180, 750, 90, 30, '#0b0e1c', { r: 4 }) }
    return { svg: s, floorY: 1030, lights: [{ x: 960, y: 300, r: 900, a: .8 }, { x: 960, y: 800, r: 700, a: .45 }], glows: gl, cones: [] };
  },
  map(o) {
    let s = sR(0, 0, 1920, 1080, o.sea || '#16294e', { sw: 0 }); for (let i = 1; i < 12; i++) s += sL(0, i * 90, 1920, i * 90, 1.5, '#2c4678') + sL(i * 160, 0, i * 160, 1080, 1.5, '#2c4678');
    const land = '#C9A878', ls = (d) => sP(d, land, { sw: 5 });
    s += ls('M120 190 L330 140 L540 170 L600 280 L520 420 L440 520 L330 460 L250 360 L150 300 Z');            // north america
    s += ls('M470 580 L580 560 L650 650 L620 800 L550 940 L500 800 L460 680 Z');                                // south america
    s += ls('M900 210 L1010 170 L1120 190 L1090 280 L1000 320 L920 300 Z');                                   // europe
    s += ls('M920 350 L1060 340 L1130 420 L1120 600 L1040 740 L960 680 L920 520 Z');                           // africa
    s += ls('M1130 180 L1400 140 L1700 200 L1790 330 L1650 440 L1500 500 L1360 440 L1230 350 L1140 290 Z');   // asia
    s += ls('M1560 640 L1740 620 L1790 720 L1680 800 L1560 760 Z');                                            // australia
    s += ls('M560 70 L680 60 L720 110 L640 150 L580 120 Z');                                                   // greenland
    (o.cities || []).forEach(c => s += `<circle cx="${c[0]}" cy="${c[1]}" r="9" fill="#FFD27A"/><circle cx="${c[0]}" cy="${c[1]}" r="26" fill="#FFD27A" opacity=".25"/>`);
    return { svg: s, floorY: 1010, lights: [{ x: 960, y: 540, r: 1200, a: .6 }], glows: [], cones: [] };
  },
  vault(o) {
    let s = sR(0, 0, 1920, 1080, '#2a2d3a', { sw: 0 }); for (let i = 0; i < 16; i++) s += sR(i * 130, 0, 126, 700, i % 2 ? '#32364a' : '#2d3144', { sw: 3 });
    s += sR(0, 740, 1920, 340, '#1b1d28', { sw: 0 }) + sL(0, 740, 1920, 740, 6);
    s += sC(960, 420, 330, '#6a7080') + sC(960, 420, 290, '#8b93a6') + sC(960, 420, 210, '#6a7080') + sC(960, 420, 60, '#d9b66a');
    for (let i = 0; i < 12; i++) { const a = i * Math.PI / 6; s += sC(960 + 250 * Math.cos(a), 420 + 250 * Math.sin(a), 14, '#d9b66a', { sw: 4 }) }
    for (let i = 0; i < 4; i++) { const a = i * Math.PI / 2 + .4; s += sL(960, 420, 960 + 170 * Math.cos(a), 420 + 170 * Math.sin(a), 16, '#d9b66a') + sC(960 + 170 * Math.cos(a), 420 + 170 * Math.sin(a), 24, '#d9b66a') }
    return { svg: s, floorY: 1020, lights: [{ x: 960, y: 420, r: 620, a: .9 }, { x: 960, y: 900, r: 700, a: .4 }], glows: [{ x: 960, y: 420, r: 420, c: '#FFD27A', a: .25 }], cones: [] };
  },
  server(o) {
    let s = sR(0, 0, 1920, 1080, '#08202b', { sw: 0 }); const r = sRng(6), gl = [];
    s += sP('M0 1080 L760 560 L1160 560 L1920 1080 Z', '#10333f', { sw: 0 });
    for (let k = 0; k < 6; k++) {
      const x = 80 + k * 110, w = 90; s += sP(`M${x} ${1000 - k * 60} L${x + w} ${960 - k * 60} L${x + w} ${340 + k * 30} L${x} ${250 + k * 40} Z`, '#1a4a5a', { sw: 4 });
      const x2 = 1840 - k * 110; s += sP(`M${x2} ${1000 - k * 60} L${x2 - w} ${960 - k * 60} L${x2 - w} ${340 + k * 30} L${x2} ${250 + k * 40} Z`, '#1a4a5a', { sw: 4 });
      for (let j = 0; j < 10; j++) { const on = r() < .6; s += `<rect x="${x + 14}" y="${300 + k * 36 + j * 56}" width="${w - 28}" height="10" fill="${on ? '#3dffe0' : '#124a52'}"/>`; if (on && gl.length < 24) gl.push({ x: x + 40, y: 305 + k * 36 + j * 56, r: 70, c: '#3dffe0', a: .12 }) }
    }
    s += sR(760, 560, 400, 20, '#0b2a33', { sw: 0 });
    return { svg: s, floorY: 1030, lights: [{ x: 960, y: 700, r: 800, a: .85 }], glows: gl, cones: [] };
  },
  capitol(o) {
    const t = o.time || 'dusk'; let s = sky(t, 760);
    s += sR(0, 760, 1920, 320, '#3b3350', { sw: 0 }) + sL(0, 760, 1920, 760, 5);
    s += sR(380, 520, 1160, 240, '#EDE6D4', { sw: 6 }) + sP('M380 520 L960 400 L1540 520 Z', '#E2DAC4', { sw: 6 });
    for (let i = 0; i < 12; i++) s += sR(420 + i * 94, 540, 40, 210, '#F8F3E4', { sw: 4 });
    s += sR(800, 330, 320, 70, '#EDE6D4', { sw: 6 }) + sP('M820 330 Q960 130 1100 330 Z', '#F8F3E4', { sw: 6 }) + sL(960, 150, 960, 100, 6) + sP('M960 100 L992 108 L960 118 Z', '#B3283A', { sw: 3 });
    for (let i = 0; i < 6; i++) s += sR(360 - i * 12, 760 + i * 26, 1200 + i * 24, 26, i % 2 ? '#d8d0bc' : '#e4dcc6', { sw: 4 });
    [[200, 800], [1720, 800]].forEach(([x, y]) => s += sR(x - 12, y - 60, 24, 160, '#6a4a30', { sw: 5 }) + sC(x, y - 90, 70, '#2d6a4a'));
    return { svg: s, floorY: 1030, lights: [{ x: 960, y: 560, r: 700, a: .5 }, { x: 960, y: 900, r: 900, a: .25 }], glows: [{ x: 960, y: 520, r: 500, c: '#FFB38A', a: .22 }], cones: [] };
  },
  shop(o) {
    let s = sR(0, 0, 1920, 1080, '#cdbd9c', { sw: 0 }); const r = sRng(77), gl = [];
    for (let i = 0; i < 16; i++) s += `<rect x="${i * 120}" y="0" width="60" height="760" fill="#c2b08c" opacity=".5"/>`;
    s += sR(0, 760, 1920, 320, '#7a6a58', { sw: 0 }) + sL(0, 760, 1920, 760, 6);
    for (let i = 0; i < 14; i++) s += sL(i * 150, 760, i * 150 - 110 + i * 14, 1080, 3, '#5a4c3d');
    for (let j = 1; j < 4; j++) s += sL(0, 760 + j * 80, 1920, 760 + j * 80, 2, '#5a4c3d');
    // shelf units with boxes of filters and spark plugs
    const cols = ['#2b5aa8', '#e8812a', '#d9322e', '#f2c230', '#4a9a5a', '#e8e2d0'];
    [[70, 520], [640, 520], [1210, 520]].forEach(([x0, w]) => {
      s += sR(x0, 200, w, 540, '#6a4a30', { sw: 6 });
      for (let j = 0; j < 4; j++) {
        const y = 250 + j * 120; s += sR(x0 + 10, y + 100, w - 20, 14, '#8a6a46', { sw: 4 });
        let x = x0 + 26;
        while (x < x0 + w - 70) {
          const bw = 44 + r() * 34, bh = 56 + r() * 36, c = cols[Math.floor(r() * cols.length)];
          if (j % 2 === 0) s += sR(x, y + 100 - bh, bw, bh, c, { sw: 4 }) + sR(x + 6, y + 100 - bh + 10, bw - 12, 14, '#fff', { sw: 0, op: .8 });
          else s += sC(x + bw / 2, y + 100 - bw / 2 - 4, bw / 2, c, { sw: 4 }) + sR(x + bw / 2 - 6, y + 100 - bw - 18, 12, 24, '#c9c9d6', { sw: 3 });
          x += bw + 12;
        }
      }
    });
    // neon sign
    s += sR(560, 80, 800, 110, '#12121c', { r: 16, sw: 7 }) + sT(960, 162, o.sign || 'AUTOPEÇAS', 80, '#ff5a6a', { cls: 'brush', st: 6 });
    // ceiling lamps
    [380, 960, 1540].forEach(x => s += sR(x - 120, 0, 240, 26, '#e8e2d0', { sw: 5 }));
    // counter with register, front right
    s += sR(1340, 700, 500, 200, '#6a3a2c', { sw: 6 }) + sR(1320, 676, 540, 34, '#d9b66a', { sw: 6 }) + sR(1500, 596, 130, 80, '#9a9aa8', { r: 8 }) + sR(1520, 610, 90, 30, '#bfeaff', { sw: 3 });
    // floor mat and a stack of tyres
    s += `<ellipse cx="760" cy="960" rx="360" ry="58" fill="#3a3340" stroke="${SINK}" stroke-width="5"/>`;
    for (let k = 0; k < 3; k++) s += `<ellipse cx="200" cy="${900 - k * 40}" rx="96" ry="30" fill="#1c1c24" stroke="${SINK}" stroke-width="5"/><ellipse cx="200" cy="${888 - k * 40}" rx="96" ry="30" fill="#2a2a33" stroke="${SINK}" stroke-width="5"/>`;
    gl.push({ x: 960, y: 135, r: 380, c: '#ff5a6a', a: .28 });
    return { svg: s, floorY: 1010, lights: [{ x: 380, y: 160, r: 600, a: .85 }, { x: 960, y: 160, r: 700, a: .9 }, { x: 1540, y: 160, r: 600, a: .85 }, { x: 960, y: 800, r: 1000, a: .3 }], glows: gl, cones: [380, 960, 1540].map(x => ({ x, y: 26, w0: 140, w1: 480, len: 520, lean: 0, a: .16, c: '#FFF2C8' })) };
  },
  living(o) {
    let s = sR(0, 0, 1920, 1080, '#4a3a5a', { sw: 0 }) + sR(0, 700, 1920, 380, '#5a3f2e', { sw: 0 }) + sL(0, 700, 1920, 700, 6);
    s += sR(560, 180, 800, 450, '#12121c', { r: 14, sw: 8 }) + `<rect x="580" y="200" width="760" height="410" fill="${o.screen || '#2a4a9a'}"/>` + sR(840, 630, 240, 24, '#12121c', { sw: 5 });
    s += sR(120, 520, 480, 220, '#8a3a4a', { r: 40, sw: 6 }) + sR(100, 450, 120, 290, '#7a2f3f', { r: 40, sw: 6 }) + sR(1300, 520, 480, 220, '#3a5a8a', { r: 40, sw: 6 });
    s += sR(1500, 300, 60, 220, '#6a4a30', { sw: 5 }) + sP('M1470 300 L1590 300 L1560 220 L1500 220 Z', '#F7D88A', { sw: 5 });
    return { svg: s, floorY: 1020, lights: [{ x: 960, y: 400, r: 620, a: .8 }, { x: 1530, y: 280, r: 380, a: .7 }, { x: 960, y: 800, r: 900, a: .3 }], glows: [{ x: 960, y: 400, r: 500, c: o.screen || '#5a8aff', a: .3 }], cones: [] };
  }
};
function setSVG(o, t) { const f = SETS[o.kind]; return f ? f(o, t) : { svg: sR(0, 0, 1920, 1080, '#333', { sw: 0 }), floorY: 1000, lights: [], glows: [], cones: [] } }
/* darkness with soft holes at the light sources + additive glow + cones: one mood for the whole picture */
const NIGHT = { night: '#34407f', dusk: '#7a5a8a', day: null, red: '#7a2030' };
function nightSVG(o, S, t) {
  const time = o.time || 'night', amb = o.ambient || NIGHT[time]; if (!amb || o.lighting === false) return '';
  const fl = 1 + (o.flicker ? .06 * Math.sin(t * 17) : 0); let defs = '', holes = '', gl = '', cn = '';
  (S.lights || []).concat(o.lights || []).forEach((L, i) => {
    const a = Math.min(1, L.a * fl);
    defs += `<radialGradient id="hl${i}"><stop offset="0" stop-color="#000" stop-opacity="${a}"/><stop offset=".55" stop-color="#000" stop-opacity="${(a * .72).toFixed(3)}"/><stop offset="1" stop-color="#000" stop-opacity="0"/></radialGradient>`;
    holes += `<circle cx="${L.x}" cy="${L.y}" r="${L.r}" fill="url(#hl${i})"/>`;
  });
  (S.glows || []).concat(o.glows || []).forEach((G, i) => {
    defs += `<radialGradient id="gw${i}"><stop offset="0" stop-color="${G.c}" stop-opacity="${G.a}"/><stop offset=".5" stop-color="${G.c}" stop-opacity="${(G.a * .4).toFixed(3)}"/><stop offset="1" stop-color="${G.c}" stop-opacity="0"/></radialGradient>`;
    gl += `<circle cx="${G.x}" cy="${G.y}" r="${G.r}" fill="url(#gw${i})"/>`;
  });
  (S.cones || []).concat(o.cones || []).forEach((C, i) => {
    defs += `<linearGradient id="cn${i}" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="${C.c || '#FFD27A'}" stop-opacity="${C.a}"/><stop offset="1" stop-color="${C.c || '#FFD27A'}" stop-opacity="0"/></linearGradient>`;
    cn += `<path d="M${C.x - C.w0 / 2} ${C.y} L${C.x + C.w0 / 2} ${C.y} L${C.x + (C.lean || 0) + C.w1 / 2} ${C.y + C.len} L${C.x + (C.lean || 0) - C.w1 / 2} ${C.y + C.len} Z" fill="url(#cn${i})"/>`;
  });
  return `<defs>${defs}<mask id="nmask" maskUnits="userSpaceOnUse" x="0" y="0" width="1920" height="1080"><rect width="1920" height="1080" fill="#fff"/>${holes}</mask></defs><rect width="1920" height="1080" fill="${amb}" mask="url(#nmask)" style="mix-blend-mode:multiply"/><g style="mix-blend-mode:screen">${cn}${gl}</g>`;
}
