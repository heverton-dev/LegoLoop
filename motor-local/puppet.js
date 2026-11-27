/* RobbiSkill v2: code-drawn puppet characters (original design: round head, pill body, one ink outline, hard-edged shadow).
   A puppet is drawn in SVG from a small spec; no images needed. Feet are at (x, y); s = 1 is a full-body wide-shot size (head radius 88).
   chars:[{"puppet":"hero","x":900,"y":1000,"s":1,"face":.3,"mood":"smug","pose":"gesture","talk":true}]
   "puppet" is a key into timeline.cast or an inline spec. */
const PINK = '#18161B';
const PP_POSE = {
  idle: ['down', 'down'], walk: ['down', 'down'], cheer: ['cheer', 'cheer'], panic: ['panic', 'panic'], clap: ['clap', 'clap'], think: ['down', 'chin'],
  point: ['down', 'point'], present: ['down', 'present'], gesture: ['hip', 'gesture'], shrug: ['shrug', 'shrug'], facepalm: ['down', 'facepalm'],
  hip: ['hip', 'hip'], pocket: ['pocket', 'pocket'], hold: ['hold', 'hold'], wave: ['down', 'wave'], push: ['down', 'push'], type: ['push', 'push']
};
/* left arm geometry [elbow, hand, showHand]; the right arm is mirrored */
const PP_ARM = {
  down: [[-92, -190], [-98, -130], 1], hip: [[-126, -195], [-84, -152], 1], pocket: [[-96, -188], [-70, -140], 0], present: [[-126, -255], [-206, -292], 1],
  gesture: [[-134, -228], [-150, -312], 1], chin: [[-122, -226], [-62, -286], 1], hold: [[-92, -190], [-46, -176], 1], push: [[-118, -226], [-190, -222], 1],
  point: [[-140, -246], [-232, -252], 1], cheer: [[-132, -300], [-152, -392], 1], panic: [[-112, -302], [-96, -396], 1], wave: [[-122, -292], [-150, -356], 1],
  clap: [[-62, -202], [-22, -232], 1], shrug: [[-124, -214], [-112, -252], 1], facepalm: [[-102, -252], [-34, -346], 1]
};
const PP_MOOD = { // eyes, mouth, brow lift, brow tilt, sweat
  neutral: ['dot', 'flat', 0, 0, 0], happy: ['dot', 'smile', 1, 0, 0], grin: ['dot', 'open', 1, 0, 0], shock: ['wide', 'o', 2, 0, 0], panic: ['wide', 'open', 2, 1, 1],
  worry: ['wide', 'flat', 1, 1, 1], smug: ['half', 'smirk', 0, 0, 0], sad: ['dot', 'frown', 0, 1, 0], angry: ['dot', 'flat', -1, -1, 0], sleep: ['closed', 'flat', 0, 0, 0],
  think: ['half', 'flat', 1, 0, 0], laugh: ['closed', 'open', 1, 0, 0]
};
function ppSpec(c) {
  let s = c.puppet; if (typeof s === 'string') s = (SH && SH.cast && SH.cast[s]) || {};
  const o = Object.assign({ skin: '#E6BE8E', suit: '#3B4660', shirt: '#F4EFE2', tie: '#A8232B', hair: 'short', hair_col: '#1C1B22' }, s || {});
  o.skin_d = o.skin_d || shade(o.skin, .14); o.suit_d = o.suit_d || shade(o.suit, .26); return o;
}
function ppArmGeo(name, t, side) {
  const g = PP_ARM[name] || PP_ARM.down; let e = g[0].slice(), h = g[1].slice();
  if (name === 'panic') { h[0] += 7 * Math.sin(t * 24 + side); h[1] += 5 * Math.sin(t * 29) }
  if (name === 'wave') { h[0] += 24 * Math.sin(t * 9 + 1) }
  if (name === 'clap') { const k = (Math.sin(t * 12) + 1) / 2; h[0] = -8 - 40 * (1 - k); e[0] = -50 - 18 * (1 - k) }
  if (name === 'cheer') { h[1] -= 12 * Math.abs(Math.sin(t * 5)) }
  if (name === 'push') { h[0] += 6 * Math.sin(t * 14 + side * 2) }
  return { e, h, show: g[2] };
}
function ppHair(st, col, R, fs, sw) {
  const f = `fill="${col}" stroke="${PINK}" stroke-width="${sw}" stroke-linejoin="round"`;
  const hl = (d) => `<path d="${d}" fill="#fff" opacity=".22"/>`;
  const cap = `M${-R - 6} 4 A${R + 6} ${R + 6} 0 0 1 ${R + 6} 4 C${R + 4} -30 ${58 + fs} -70 ${10 + fs} -66 C${-30 + fs} -62 ${-R - 2} -34 ${-R - 6} 4 Z`;
  const strands = (n) => { return ""; let q = ''; for (let i = 0; i < n; i++) { const x = -52 + i * (104 / (n - 1)) + fs * .5, y0 = -30 - 22 * (1 - Math.pow((x - fs * .5) / 70, 2)); q += `<path d="M${x} ${y0} Q${x + 7} ${y0 - 9} ${x + 15} ${y0 - 2}" fill="none" stroke="${PINK}" stroke-width="${sw * .5}" stroke-linecap="round" opacity=".5"/>`; } return q; };
  switch (st) {
    case 'bald': return '';
    case 'slick': return `<path d="${cap}" ${f}/>${hl(`M${-50 + fs} -62 Q${-10 + fs} -78 ${34 + fs} -64 Q${0 + fs} -70 ${-50 + fs} -62 Z`)}${strands(6)}`;
    case 'short': return `<path d="M${-R - 4} 0 A${R + 4} ${R + 4} 0 0 1 ${R + 4} 0 L${R - 4} -38 Q${fs} -66 ${-R + 4} -38 Z" ${f}/>${hl(`M${-48 + fs} -60 Q${-6 + fs} -76 ${38 + fs} -62 Q${0 + fs} -68 ${-48 + fs} -60 Z`)}${strands(7)}`;
    case 'curly': return [[-66, -46, 34], [-36, -72, 36], [4, -82, 36], [44, -70, 34], [68, -42, 30], [-76, -14, 26], [80, -10, 24], [-14, -56, 30], [26, -50, 28]].map(([x, y, r]) => `<circle cx="${x + fs * .5}" cy="${y}" r="${r}" ${f}/>`).join('') + [[-30, -78], [10, -88], [48, -74]].map(([x, y]) => `<circle cx="${x + fs * .5}" cy="${y}" r="9" fill="#fff" opacity=".2"/>`).join('');
    case 'bob': return `<path d="M${-R - 6} 40 L${-R - 6} -8 A${R + 6} ${R + 6} 0 0 1 ${R + 6} -8 L${R + 6} 40 Q${R - 6} 46 ${R - 12} 10 L${R - 18} -34 Q${fs} -62 ${-R + 18} -34 L${-R + 12} 10 Q${-R + 6} 46 ${-R - 6} 40 Z" ${f}/>${hl(`M${-50 + fs} -66 Q${0 + fs} -82 ${48 + fs} -66 Q${0 + fs} -72 ${-50 + fs} -66 Z`)}`;
    case 'long': return `<path d="M${-R - 8} 130 L${-R - 8} -8 A${R + 8} ${R + 8} 0 0 1 ${R + 8} -8 L${R + 8} 130 Q${R - 6} 138 ${R - 14} 60 L${R - 18} -34 Q${fs} -62 ${-R + 18} -34 L${-R + 14} 60 Q${-R + 6} 138 ${-R - 8} 130 Z" ${f}/>${hl(`M${-50 + fs} -66 Q${0 + fs} -82 ${48 + fs} -66 Q${0 + fs} -72 ${-50 + fs} -66 Z`)}`;
    case 'bun': return `<path d="${cap}" ${f}/><circle cx="${fs * .6}" cy="${-R - 22}" r="30" ${f}/>${hl(`M${-50 + fs} -62 Q${-10 + fs} -78 ${34 + fs} -64 Q${0 + fs} -70 ${-50 + fs} -62 Z`)}`;
    case 'grey': return `<path d="${cap}" fill="${col && col !== '#1C1B22' ? col : '#B9BCC6'}" stroke="${PINK}" stroke-width="${sw}" stroke-linejoin="round"/>${hl(`M${-50 + fs} -62 Q${-10 + fs} -78 ${34 + fs} -64 Q${0 + fs} -70 ${-50 + fs} -62 Z`)}${strands(6)}`;
    default: return '';
  }
}
function puppetSVG(c, t, so) {
  const sp = ppSpec(c), s = c.s || 1, fl = c.flip ? -1 : 1, lw = Math.max(5.2, 4.2 / s), R = 88, HY = -372, BT = -300, BB = -125;
  const tt = t + (c.phase || 0), pose = c.pose || 'idle', key = c._prevPose, st0 = c._st;
  const md = PP_MOOD[c.mood || 'neutral'] || PP_MOOD.neutral; let eyes = c.eyes || md[0], mouth = c.mouth || md[1];
  const talk = c.speaking || c.talk;
  if (talk && (!c.mouth || c.mouth === 'auto')) {
    const tq = Math.floor((t + (T0 || 0)) * FPS / 3) * 3 / FPS - (T0 || 0), v = speakEnv(tq); mouth = v > .6 ? 'open' : v > .2 ? 'mid' : md[1] === 'open' ? 'flat' : md[1];
  } else if (mouth === 'auto') mouth = md[1];
  const blinking = ((tt + (c.phase || 0) * 3) % 2.7) < .1; if (blinking && (eyes === 'dot' || eyes === 'wide')) eyes = 'closed';
  const face = c.face == null ? 0 : c.face, light = c.lightSide == null ? -1 : c.lightSide, brow = (c.brow != null ? c.brow : md[2]), bt = (c.browTilt != null ? c.browTilt : md[3]);
  // body motion
  let bob = 3 * Math.abs(Math.sin(tt * 2.4)), tilt = 1.2 * Math.sin(tt * 1.1), dx = 0, squash = 0, lift = 0;
  if (pose === 'cheer') { lift = 26 * Math.abs(Math.sin(tt * 5.5)); tilt = 3 * Math.sin(tt * 5.5) }
  if (pose === 'panic') { dx = 6 * Math.sin(tt * 30); tilt = 3 * Math.sin(tt * 24); bob = 4 * Math.abs(Math.sin(tt * 18)) }
  if (pose === 'clap') bob = 8 * Math.abs(Math.sin(tt * 6));
  if (pose === 'think') tilt = -3 + Math.sin(tt) * 1.2;
  if (pose === 'shrug') bob = 6 * Math.max(0, Math.sin(tt * 2.2));
  const wk = pose === 'walk' ? Math.sin(tt * 11) : 0; if (pose === 'walk') { bob = 8 * Math.abs(wk); tilt = 2 * wk }
  if (st0 != null && t - st0 < .3 && t >= st0) { const k = (t - st0) / .3; squash = .07 * Math.sin(k * Math.PI) * (1 - k * .5) }
  // enter slide
  let ex = 0; if (c.enter) { const e = ease((t - (c.enterAt || 0)) / .6); ex = c.enter === 'left' ? -(1 - e) * 520 : (1 - e) * 520 }
  // arms (snap with overshoot when the pose changed)
  const names = c.arms || PP_POSE[pose] || PP_POSE.idle; let aL = ppArmGeo(names[0], tt, -1), aR = ppArmGeo(names[1], tt, 1);
  if (c._prevArms && st0 != null && t - st0 < .3 && t >= st0) {
    const k = back((t - st0) / .3), pL = ppArmGeo(c._prevArms[0], tt, -1), pR = ppArmGeo(c._prevArms[1], tt, 1), mix = (p, q) => ({ e: [p.e[0] + (q.e[0] - p.e[0]) * k, p.e[1] + (q.e[1] - p.e[1]) * k], h: [p.h[0] + (q.h[0] - p.h[0]) * k, p.h[1] + (q.h[1] - p.h[1]) * k], show: k > .5 ? q.show : p.show });
    aL = mix(pL, aL); aR = mix(pR, aR);
  }
  const SH0 = [-56, BT + 52];
  let g = `<ellipse cx="${(c.x + ex + 12 * s).toFixed(0)}" cy="${c.y + 3}" rx="${(86 * s).toFixed(0)}" ry="${(15 * s).toFixed(0)}" fill="#000" opacity=".28" style="filter:blur(${(8 * s).toFixed(0)}px)"/>`;
  g += `<g transform="translate(${(c.x + ex + dx).toFixed(1)} ${(c.y - lift).toFixed(1)}) scale(${(s * fl).toFixed(4)} ${s})">`;
  const clip = c.cut != null ? `<clipPath id="pcut${Math.round(c.x)}"><rect x="-400" y="-900" width="800" height="${900 + c.cut}"/></clipPath>` : '';
  g += clip + (c.cut != null ? `<g clip-path="url(#pcut${Math.round(c.x)})">` : '<g>');
  // legs
  const legs = [[-24, wk * 30, Math.max(0, Math.cos(tt * 11)) * 16], [24, -wk * 30, Math.max(0, -Math.cos(tt * 11)) * 16]];
  legs.forEach(([x0, fx, lf]) => { g += `<line x1="${x0}" y1="${BB - 8}" x2="${x0 + fx}" y2="${-lf - 6}" stroke="${PINK}" stroke-width="${lw}" stroke-linecap="round"/><ellipse cx="${x0 + fx + 8 * Math.sign(face || 1)}" cy="${-lf - 2}" rx="${19}" ry="9" fill="${sp.shoes || '#24202a'}" stroke="${PINK}" stroke-width="${lw * .8}"/>`; });
  g += `<g transform="translate(0 ${-bob.toFixed(1)}) rotate(${tilt.toFixed(2)} 0 ${BB}) scale(${(1 + squash * .6).toFixed(3)} ${(1 - squash).toFixed(3)})">`;
  // back hair (long/bob) behind the body
  if (sp.hair === 'long' || sp.hair === 'bob') g += `<g transform="translate(${face * 10} ${HY})">${ppHair(sp.hair, sp.hair_col, R, face * 14, lw)}</g>`;
  // body
  const bodyP = `M-84 ${BT + 70} Q-84 ${BT} -14 ${BT} L14 ${BT} Q84 ${BT} 84 ${BT + 70} L84 ${BB - 50} Q84 ${BB} 30 ${BB} L-30 ${BB} Q-84 ${BB} -84 ${BB - 50} Z`;
  const dress = sp.dress;
  g += `<path d="${bodyP}" fill="${dress || sp.suit_d}"/>`;
  g += `<clipPath id="bc${Math.round(c.x)}"><path d="${bodyP}"/></clipPath><g clip-path="url(#bc${Math.round(c.x)})">`;
  g += `<path d="${bodyP}" transform="translate(${-light * 18} -4)" fill="${dress || sp.suit}"/>`;
  if (!dress) {
    g += `<path d="M-36 ${BT - 4} L0 ${BT + 96} L36 ${BT - 4} Z" fill="${sp.shirt}"/>`;
    if (sp.bowtie) g += `<path d="M0 ${BT + 38} L-26 ${BT + 24} L-26 ${BT + 52} Z M0 ${BT + 38} L26 ${BT + 24} L26 ${BT + 52} Z" fill="${sp.tie}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
    else g += `<path d="M-10 ${BT + 24} L10 ${BT + 24} L14 ${BT + 84} L0 ${BT + 100} L-14 ${BT + 84} Z" fill="${sp.tie}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
    [-1, 1].forEach(m => g += `<path d="M${m * 38} ${BT - 2} L${m * 2} ${BT + 100} L${m * 2} ${BB}" fill="none" stroke="${PINK}" stroke-width="${lw * .7}"/>`);
    // detail: collar points, buttons, pocket + pocket square, shoulder seams, jacket hem shade
    g += `<path d="M-14 ${BT + 2} L-30 ${BT + 26} L-6 ${BT + 22} Z M14 ${BT + 2} L30 ${BT + 26} L6 ${BT + 22} Z" fill="#fff" stroke="${PINK}" stroke-width="${lw * .55}" stroke-linejoin="round"/>`;
    [BT + 130, BT + 158].forEach(y => g += `<circle cx="${light > 0 ? -10 : 10}" cy="${y}" r="4.2" fill="${PINK}" opacity=".75"/>`);
    g += `<path d="M${-light * -50} ${BT + 96} l24 0" stroke="${PINK}" stroke-width="${lw * .5}" opacity=".6"/><path d="M${-52} ${BT + 92} l22 -2 l-1 -12 l-9 3 l-9 -4 z" fill="${sp.pocket || '#F4EFE2'}" stroke="${PINK}" stroke-width="${lw * .45}" opacity=".95"/>`;
    g += `<rect x="-84" y="${BB - 34}" width="168" height="40" fill="#000" opacity=".12"/>`;
  }
  g += `</g><path d="${bodyP}" fill="none" stroke="${PINK}" stroke-width="${lw}"/>`;
  // arms
  [[aL, -1], [aR, 1]].forEach(([a, m]) => {
    const sx = SH0[0] * (-m), sy = SH0[1], ex_ = a.e[0] * (-m), ey = a.e[1], hx = a.h[0] * (-m), hy = a.h[1];
    const cx = 1.35 * ex_ - .175 * (sx + hx), cy = 1.35 * ey - .175 * (sy + hy);
    g += `<path d="M${sx} ${sy} Q${cx.toFixed(1)} ${cy.toFixed(1)} ${hx.toFixed(1)} ${hy.toFixed(1)}" fill="none" stroke="${PINK}" stroke-width="${lw}" stroke-linecap="round"/>`;
    if (a.show) g += `<circle cx="${hx.toFixed(1)}" cy="${hy.toFixed(1)}" r="14" fill="${sp.skin}" stroke="${PINK}" stroke-width="${lw * .8}"/>`;
  });
  // head
  g += `<g transform="translate(${(face * 6).toFixed(1)} ${HY}) rotate(${(tilt * 1.3 + (c.tilt || 0)).toFixed(2)})">`;
  if (sp.hair === 'bald') [-1, 1].forEach(m => g += `<path d="M${m * 84} -40 L${m * 108} -34 L${m * 98} -18 L${m * 116} -4 L${m * 100} 10 L${m * 110} 26 L${m * 86} 34 Z" fill="${sp.hair_col === '#1C1B22' ? '#9A9AA4' : sp.hair_col}" stroke="${PINK}" stroke-width="${lw * .8}" stroke-linejoin="round"/>`);
  g += `<circle r="${R}" fill="${sp.skin_d}"/><clipPath id="hc${Math.round(c.x)}"><circle r="${R}"/></clipPath>`;
  g += `<circle cx="${light * 8}" cy="-18" r="${R}" fill="${sp.skin}" clip-path="url(#hc${Math.round(c.x)})"/>`;
  if (c.rim) g += `<circle cx="${-light * 10}" cy="4" r="${R}" fill="none" stroke="${c.rim}" stroke-width="7" opacity=".7" clip-path="url(#hc${Math.round(c.x)})"/>`;
  g += `<circle r="${R}" fill="none" stroke="${PINK}" stroke-width="${lw}"/>`;
  if (sp.hair !== 'long' && sp.hair !== 'bob') g += ppHair(sp.hair, sp.hair_col, R, face * 14, lw); else g += ppHair('short', sp.hair_col, R, face * 14, lw);
  if (sp.hat) { const hc = sp.hat_col || '#2E2E38'; g += `<path d="M${-R + 8} -54 C${-R} -${R + 78} ${R} -${R + 78} ${R - 8} -54 Z" fill="${hc}" stroke="${PINK}" stroke-width="${lw}"/><rect x="${-R + 8}" y="-82" width="${2 * R - 16}" height="26" fill="${sp.hat_band || '#7E1C24'}" stroke="${PINK}" stroke-width="${lw * .7}"/><ellipse cx="0" cy="-54" rx="${R + 38}" ry="15" fill="${hc}" stroke="${PINK}" stroke-width="${lw}"/>`; }
  // ears + blush + nose (drawn under the features)
  [-1, 1].forEach(m => { if (Math.abs(face) < .85 || (face > 0) === (m < 0)) g += `<ellipse cx="${m * (R - 2)}" cy="6" rx="11" ry="17" fill="${sp.skin_d}" stroke="${PINK}" stroke-width="${lw * .7}" transform="translate(${-face * 4 * m} 0)"/>`; });
  [-1, 1].forEach(m => g += `<ellipse cx="${m * 52 + face * 30}" cy="28" rx="${15 - Math.abs(face) * 4}" ry="9" fill="#E8727A" opacity="${md[4] || c.sweat ? .1 : (c.blush == null ? .32 : c.blush)}"/>`);
  g += `<path d="M${face * 36 - 4} 18 Q${face * 36 - 11} 30 ${face * 36 + 2} 32" fill="none" stroke="${PINK}" stroke-width="${lw * .55}" stroke-linecap="round" opacity=".7"/>`;
  const fx = face * 38, ey0 = -4;
  [-30, 30].forEach(e0 => {
    const x = e0 + fx, far = (face > 0 && e0 > 0) || (face < 0 && e0 < 0), q = far ? 1 - .3 * Math.abs(face) : 1;
    if (eyes === 'wide') g += `<ellipse cx="${x}" cy="${ey0}" rx="${17 * q}" ry="17" fill="#FBFBF7" stroke="${PINK}" stroke-width="${lw * .7}"/><ellipse cx="${x + (c.look || 0) * 6}" cy="${ey0}" rx="${5.5 * q}" ry="5.5" fill="${PINK}"/>`;
    else if (eyes === 'closed') g += `<path d="M${x - 11 * q} ${ey0} Q${x} ${ey0 + 9} ${x + 11 * q} ${ey0}" fill="none" stroke="${PINK}" stroke-width="${lw * .8}" stroke-linecap="round"/>`;
    else if (eyes === 'half') g += `<clipPath id="hf${Math.round(c.x)}${e0}"><rect x="${x - 20}" y="${ey0 - 3}" width="40" height="24"/></clipPath><ellipse cx="${x}" cy="${ey0}" rx="${6.5 * q}" ry="9" fill="${PINK}" clip-path="url(#hf${Math.round(c.x)}${e0})"/><line x1="${x - 13 * q}" y1="${ey0 - 4}" x2="${x + 13 * q}" y2="${ey0 - 4}" stroke="${PINK}" stroke-width="${lw * .8}" stroke-linecap="round"/>`;
    else g += `<ellipse cx="${x + (c.look || 0) * 5}" cy="${ey0}" rx="${5.8 * q}" ry="8.4" fill="${PINK}"/><circle cx="${x + 1.5 + (c.look || 0) * 5}" cy="${ey0 - 3}" r="1.9" fill="#fff" opacity=".85"/>`;
    const inner = e0 < 0 ? 1 : -1, by = ey0 - 29 - brow * 9, tl = bt * 7;
    g += `<line x1="${x - 13 * q}" y1="${by + (inner > 0 ? tl : -tl)}" x2="${x + 13 * q}" y2="${by + (inner > 0 ? -tl : tl)}" stroke="${PINK}" stroke-width="${lw * 1.15}" stroke-linecap="round"/>`;
    if (sp.glasses) g += `<ellipse cx="${x}" cy="${ey0}" rx="${25 * q}" ry="25" fill="#e1f0ff" fill-opacity=".22" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
  });
  if (sp.glasses) g += `<line x1="${fx - 6}" y1="${ey0 - 2}" x2="${fx + 6}" y2="${ey0 - 2}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
  const mx = face * 32, my = 42;
  if (sp.moustache) g += `<path d="M${mx} ${my - 20} C${mx - 12} ${my - 28} ${mx - 30} ${my - 22} ${mx - 40} ${my - 30} C${mx - 30} ${my - 10} ${mx - 12} ${my - 12} ${mx} ${my - 15} C${mx + 12} ${my - 12} ${mx + 30} ${my - 10} ${mx + 40} ${my - 30} C${mx + 30} ${my - 22} ${mx + 12} ${my - 28} ${mx} ${my - 20} Z" fill="${sp.hair_col}"/>`;
  if (sp.beard) g += `<path d="M${-R + 6} 10 Q${-R + 10} ${R + 6} 0 ${R + 4} Q${R - 10} ${R + 6} ${R - 6} 10 Q${R - 30} ${my + 18} 0 ${my + 14} Q${-R + 30} ${my + 18} ${-R + 6} 10 Z" fill="${sp.hair_col}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
  const mw = lw * .8, dk = '#5A1A1E';
  if (mouth === 'smile') g += `<path d="M${mx - 17} ${my - 2} Q${mx} ${my + 13} ${mx + 17} ${my - 2}" fill="none" stroke="${PINK}" stroke-width="${mw}" stroke-linecap="round"/>`;
  else if (mouth === 'smirk') { const d = face >= 0 ? 1 : -1; g += `<path d="M${mx - 15 * d} ${my + 3} Q${mx + 4 * d} ${my + 12} ${mx + 20 * d} ${my - 6}" fill="none" stroke="${PINK}" stroke-width="${mw}" stroke-linecap="round"/>`; }
  else if (mouth === 'frown') g += `<path d="M${mx - 15} ${my + 8} Q${mx} ${my - 6} ${mx + 15} ${my + 8}" fill="none" stroke="${PINK}" stroke-width="${mw}" stroke-linecap="round"/>`;
  else if (mouth === 'o') g += `<ellipse cx="${mx}" cy="${my + 2}" rx="8" ry="10" fill="${dk}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
  else if (mouth === 'mid') g += `<ellipse cx="${mx}" cy="${my + 3}" rx="13" ry="7" fill="${dk}" stroke="${PINK}" stroke-width="${lw * .7}"/>`;
  else if (mouth === 'open') g += `<path d="M${mx - 20} ${my - 5} L${mx + 20} ${my - 5} C${mx + 19} ${my + 24} ${mx - 19} ${my + 24} ${mx - 20} ${my - 5} Z" fill="${dk}" stroke="${PINK}" stroke-width="${lw * .7}"/><path d="M${mx - 18} ${my - 4} L${mx + 18} ${my - 4} L${mx + 16} ${my + 3} L${mx - 16} ${my + 3} Z" fill="#FAFAF4"/><ellipse cx="${mx}" cy="${my + 17}" rx="10" ry="6" fill="#D05860"/>`;
  else g += `<line x1="${mx - 12}" y1="${my + 2}" x2="${mx + 12}" y2="${my + 2}" stroke="${PINK}" stroke-width="${mw}" stroke-linecap="round"/>`;
  if (c.mood === 'sad' || c.tear) { const k = (tt * .6) % 1, tx = (face >= 0 ? -1 : 1) * 32 + face * 30, ty = 14 + 34 * k; g += `<path d="M${tx} ${ty - 12} C${tx + 8} ${ty} ${tx + 7} ${ty + 8} ${tx} ${ty + 8} C${tx - 7} ${ty + 8} ${tx - 8} ${ty} ${tx} ${ty - 12} Z" fill="#9AD2F2" stroke="${PINK}" stroke-width="${lw * .5}" opacity="${(1 - k * .6).toFixed(2)}"/>`; }
  if (md[4] || c.sweat) { const k = (tt * .8) % 1, sx0 = (face > 0 ? -1 : 1) * 62, sy0 = -60 + 30 * k; g += `<path d="M${sx0} ${sy0 - 19} C${sx0 + 12} ${sy0 - 3} ${sx0 + 11} ${sy0 + 10} ${sx0} ${sy0 + 10} C${sx0 - 11} ${sy0 + 10} ${sx0 - 12} ${sy0 - 3} ${sx0} ${sy0 - 19} Z" fill="#8FC8EE" stroke="${PINK}" stroke-width="${lw * .7}" opacity="${(1 - k * .4).toFixed(2)}"/>`; }
  g += '</g></g></g></g>';
  return g;
}
