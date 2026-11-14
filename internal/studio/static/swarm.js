let sel = "", page = 0, curParent = "", curProc = "", closedProcs = {};
const PER_PAGE = 20, palette = ['#8ab4f8', '#f28b82', '#fdd663', '#81c995', '#c58af9'];
const roleTools = { orquestrador: 'LegoLoop CLI', planejador: 'PLAN v2', arquiteto: 'Templates v2', 'engenheiro-backend': 'Go / REST', 'engenheiro-tdd': 'Go Test / Vitest', 'engenheiro-frontend': 'React / Tailwind', 'engenheiro-secops': 'Security Gates', 'dissecador-web': 'REA Web', 'auditor-qualidade': 'Taxlinter / Funclen', 'redator-tecnico': 'Manual Creator' };
function colorIdx(name) { let h = 0; for (let i = 0; i < (name || '').length; i++) h = (h * 31 + name.charCodeAt(i)) | 0; return Math.abs(h) % palette.length; }
function botSVG(st, name, kind, anim) {
  const c = palette[colorIdx(name)], cls = 'bot ' + (anim !== false ? st : '');
  if (kind === 'coord') return `<svg class="${cls}" viewBox="0 0 8 8"><rect x="2" y="0" width="4" height="3" fill="${c}"/><rect x="3" y="1" width="1" height="1" fill="#1e293b"/><rect x="5" y="1" width="1" height="1" fill="#1e293b"/><rect x="1" y="3" width="6" height="4" fill="#334155"/><rect x="3" y="3" width="2" height="4" fill="#fff"/><rect x="3" y="4" width="2" height="2" fill="#d9480f"/><rect x="2" y="7" width="1" height="1" fill="#334155"/><rect x="5" y="7" width="1" height="1" fill="#334155"/></svg>`;
  if (kind === 'ver') return `<svg class="${cls}" viewBox="0 0 8 8"><rect x="2" y="0" width="4" height="3" fill="${c}"/><rect x="2" y="1" width="2" height="2" fill="#bae6fd"/><rect x="5" y="1" width="1" height="1" fill="#1e293b"/><rect x="1" y="3" width="5" height="4" fill="${c}"/><rect x="6" y="3" width="2" height="3" fill="#e2e8f0"/><rect x="2" y="7" width="1" height="1" fill="${c}"/><rect x="4" y="7" width="1" height="1" fill="${c}"/></svg>`;
  return `<svg class="${cls}" viewBox="0 0 8 8"><rect x="1" y="0" width="6" height="1" fill="#facc15"/><rect x="2" y="1" width="4" height="2" fill="${c}"/><rect x="3" y="2" width="1" height="1" fill="#1e293b"/><rect x="5" y="2" width="1" height="1" fill="#1e293b"/><rect x="1" y="3" width="6" height="3" fill="${c}"/><rect x="1" y="5" width="6" height="1" fill="#78350f"/><rect x="2" y="6" width="1" height="2" fill="${c}"/><rect x="5" y="6" width="1" height="2" fill="${c}"/></svg>`;
}
function fmtDur(iso) { if (!iso) return '0s'; const s = Math.max(0, Math.floor((Date.now() - new Date(iso).getTime()) / 1000)), m = Math.floor(s / 60); return m > 0 ? `${m}m ${s % 60}s` : `${s}s`; }
function renderCrumb(ws) { const path = []; let p = curParent; while (p) { path.unshift(p); const w = ws.find(x => x.name === p); p = w ? (w.parent || "") : ""; } let html = '<a data-act="sub" data-val="">Macro</a>'; for (const n of path) html += ` › <a data-act="sub" data-val="${n}">${n}</a>`; document.getElementById('crumb').innerHTML = html; }
function renderGroup(title, list, allWs, kind) {
  if (!list.length) return '';
  return `<div class="grp"><span>${title}</span><span>${list.length}</span></div>` + list.map(w => {
    const subs = allWs.filter(x => x.parent === w.name).length, subBtn = subs ? `<button class="sub-btn" data-act="sub" data-val="${w.name}">${subs} sub ▶</button>` : '';
    const ci = colorIdx(w.name), isSel = sel === w.name;
    return `<div class="w ${w.status} c${ci} ${isSel ? 'sel' : ''}" data-act="sel" data-val="${w.name}">${botSVG(w.status, w.name, kind)}<div class="w-info"><div class="row"><b>${w.name}</b><span class="row-gap">${subBtn}<span class="badge ${w.status}">${w.status}</span></span></div><div class="meta">${w.role} · ${w.model || '-'} · ${w.workspace_mode}</div></div></div>`;
  }).join('');
}
function renderFicha(allWs, ms) {
  const el = document.getElementById('ficha'); if (!sel) { el.innerHTML = ''; return; }
  const idx = allWs.findIndex(w => w.name === sel), w = idx >= 0 ? allWs[idx] : { name: sel, role: 'orquestrador', model: 'maestro', workspace_mode: 'inherit', status: 'working' };
  const ord = ms.find(m => m.to === '@' + w.name || m.to === '@all'), out = ms.find(m => m.from === w.name), subs = allWs.filter(x => x.parent === w.name).map(x => x.name);
  const iso = (w.workspace_mode === 'share' || w.workspace_mode === 'branch') ? 'Isolado (git worktree)' : 'Conjunto (inherit)', pCls = w.status === 'killed' ? 'p-100' : w.status === 'waiting' ? 'p-15' : 'p-70';
  const kind = w.role === 'orquestrador' ? 'coord' : w.role === 'auditor-qualidade' ? 'ver' : 'trab';
  el.innerHTML = `<div class="ficha">${botSVG(w.status, w.name, kind, false)}<div class="ficha-body"><div class="row"><b>${w.name}</b><span class="meta">⏱ ${fmtDur(w.created_at)} · <span class="badge ${w.status}">${w.status}</span></span></div><div class="meta">Papel: ${w.role} · Cérebro: ${w.model || '-'} · Ferramentas: ${roleTools[w.role] || 'CLI'} · Modo: ${iso}</div><div class="meta">Input: ${ord ? ord.body : 'aguardando'} · Output: ${out ? out.body : '-'}</div>${subs.length ? `<div class="meta">Subprocessos: ${subs.join(', ')}</div>` : ''}<div class="prog"><div class="prog-bar ${pCls}"></div></div><div class="ficha-actions"><input id="cmd-in" class="btn cmd-in" placeholder="Enviar instrução para @${w.name}..."><button class="btn" data-act="post" data-val="${w.name}">Enviar</button><button class="btn btn-danger" data-act="kill" data-val="${w.name}">Kill</button></div></div></div>`;
}
function render(payload) {
  window._last = payload;
  const order = (payload.order || []).filter(p => !closedProcs[p]);
  if (!curProc || !payload.processes[curProc] || closedProcs[curProc]) curProc = order[0] || "";
  document.getElementById('tabs').innerHTML = order.map(p => `<button class="tab ${p === curProc ? 'active' : ''}" data-act="proc" data-val="${p}">${p}<span class="tab-x" data-act="close" data-val="${p}">×</span></button>`).join('');
  const d = (payload.processes && payload.processes[curProc]) || { workers: [], messages: [] }, allWs = d.workers || [], ms = (d.messages || []).slice().reverse();
  renderCrumb(allWs);
  const ws = allWs.filter(w => (w.parent || "") === curParent), ver = ws.filter(w => w.role === 'auditor-qualidade'), trab = ws.filter(w => w.role !== 'auditor-qualidade');
  document.getElementById('summary').textContent = `1 coordenador · ${ver.length} verificador(es) · ${trab.length} trabalhador(es) · ${ms.length} posts`;
  document.getElementById('counts').textContent = `${ws.length + 1} neste nível`;
  document.getElementById('filter').textContent = sel ? `(filtro: ${sel})` : '(todos)';
  const coord = [{ name: curParent || 'coordenador', role: 'orquestrador', model: 'maestro', workspace_mode: 'inherit', status: 'working' }];
  document.getElementById('workers').innerHTML = renderGroup('Coordenador', coord, allWs, 'coord') + renderGroup('Verificadores', ver, allWs, 'ver') + renderGroup('Trabalhadores', trab, allWs, 'trab');
  renderFicha(allWs, ms);
  const fms = ms.filter(m => !sel || m.from === sel || m.to === '@' + sel || m.to === '@all'), totalPages = Math.max(1, Math.ceil(fms.length / PER_PAGE));
  if (page >= totalPages) page = totalPages - 1;
  document.getElementById('msgs').innerHTML = fms.slice(page * PER_PAGE, (page + 1) * PER_PAGE).map(m => `<div class="msg"><div class="meta">#${m.seq} <b>${m.from}</b> → <span class="mention">${m.to || '@all'}</span></div><div class="msg-body">${m.body}</div></div>`).join('');
  document.getElementById('pager').innerHTML = `<button class="btn" data-act="page" data-val="${page - 1}" ${page <= 0 ? 'disabled' : ''}>Anterior</button><span>Página ${page + 1} de ${totalPages}</span><button class="btn" data-act="page" data-val="${page + 1}" ${page >= totalPages - 1 ? 'disabled' : ''}>Próxima</button>`;
}
document.addEventListener('click', async e => {
  const t = e.target.closest('[data-act]'); if (!t) return;
  const act = t.dataset.act, val = t.dataset.val;
  if (act === 'close') { e.stopPropagation(); closedProcs[val] = true; if (curProc === val) curProc = ""; render(window._last); }
  else if (act === 'proc') { curProc = val; curParent = ""; sel = ""; page = 0; render(window._last); }
  else if (act === 'sub') { e.stopPropagation(); curParent = val; sel = ""; page = 0; render(window._last); }
  else if (act === 'sel') { sel = (sel === val) ? "" : val; page = 0; render(window._last); }
  else if (act === 'page') { page = Number(val); render(window._last); }
  else if (act === 'post') { const inp = document.getElementById('cmd-in'); if (!inp || !inp.value.trim()) return; await fetch('/api/swarm/post', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ proc: curProc, to: '@' + val, body: inp.value.trim() }) }); inp.value = ''; }
  else if (act === 'kill') { await fetch('/api/swarm/kill', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ proc: curProc, name: val }) }); }
});
const es = new EventSource('/api/swarm/events'); es.onmessage = e => render(JSON.parse(e.data));
