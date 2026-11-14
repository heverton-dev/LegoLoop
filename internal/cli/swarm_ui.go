package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"time"

	"legoloop/internal/board"
)

func discoverBoards(root, singleDir, procName string) map[string]string {
	out := map[string]string{}
	if root == "" {
		out[procName] = singleDir
		return out
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return out
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		p := filepath.Join(root, e.Name())
		if _, err := os.Stat(filepath.Join(p, "swarm_board.json")); err == nil {
			out[e.Name()] = p
		}
	}
	return out
}

func allBoardsPayload(root, singleDir, procName string) ([]byte, error) {
	boards := discoverBoards(root, singleDir, procName)
	names := make([]string, 0, len(boards))
	for k := range boards {
		names = append(names, k)
	}
	sort.Strings(names)
	procs := map[string]any{}
	for _, name := range names {
		sb, err := board.OpenSwarmBoard(boards[name])
		if err != nil {
			continue
		}
		st, _ := sb.Status()
		msgs, _ := sb.ReadMessages(0, nil)
		procs[name] = map[string]any{"workers": st.Workers, "messages": msgs}
	}
	return json.Marshal(map[string]any{"order": names, "processes": procs})
}

func handleSwarmSSE(root, singleDir, procName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE nao suportado", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		var last string
		tick := time.NewTicker(250 * time.Millisecond)
		defer tick.Stop()
		for {
			if data, err := allBoardsPayload(root, singleDir, procName); err == nil && string(data) != last {
				last = string(data)
				_, _ = fmt.Fprintf(w, "data: %s\n\n", last)
				flusher.Flush()
			}
			select {
			case <-r.Context().Done():
				return
			case <-tick.C:
			}
		}
	}
}

func cmdSwarmUI(_ context.Context, args []string, d Deps) int {
	fs := newFlags("swarm ui", d)
	addr := fs.String("addr", "127.0.0.1:7879", "endereco HTTP do painel")
	dir := fs.String("dir", defaultBoardDir(), "diretorio do quadro")
	root := fs.String("root", "", "diretorio raiz com multiplos quadros")
	proc := fs.String("process", "swarm", "nome do processo em execucao")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swarmUIHTML))
	})
	mux.HandleFunc("/api/events", handleSwarmSSE(*root, *dir, *proc))
	fmt.Fprintf(d.Stdout, "Painel Swarm em http://%s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		return fail(d, err)
	}
	return 0
}

const swarmUIHTML = `<!doctype html>
<html lang="pt-BR"><head><meta charset="utf-8"><title>LegoLoop Swarm</title>
<style>
:root{color-scheme:light dark;--font-ui:"Segoe UI Variable","Segoe UI",system-ui,sans-serif;--font-mono:"Cascadia Code","JetBrains Mono",ui-monospace,monospace;--bg:#f4f3ef;--surface:#fff;--surface-2:#f0eee8;--border:#dedad0;--text:#1d1c19;--muted:#6b675d;--accent:#d9480f;--accent-soft:#fde4d7;--ok:#2b8a3e;--ok-soft:#e3f4e6;--fail:#c92a2a;--fail-soft:#fde8e8;--shadow:0 1px 2px rgb(0 0 0/.06),0 4px 12px rgb(0 0 0/.06)}
@media(prefers-color-scheme:dark){:root{--bg:#131312;--surface:#1c1c1a;--surface-2:#252522;--border:#34332f;--text:#edece8;--muted:#9c988e;--accent:#f76707;--accent-soft:#3d2012;--ok:#40c057;--ok-soft:#16301d;--fail:#fa5252;--fail-soft:#3b1818}}
*{box-sizing:border-box}
html,body{height:100dvh;margin:0;overflow:hidden;font-family:var(--font-ui);background:var(--bg);color:var(--text)}
.app{display:flex;flex-direction:column;height:100dvh;padding:16px;gap:10px}
.top{flex-shrink:0;display:flex;align-items:center;justify-content:space-between;background:var(--surface);border:1px solid var(--border);border-radius:10px;padding:10px 16px;box-shadow:var(--shadow)}
.tabs{flex-shrink:0;display:flex;gap:8px;overflow-x:auto;padding-bottom:2px}
.tab{background:var(--surface);border:1px solid var(--border);color:var(--muted);padding:6px 12px;border-radius:8px;cursor:pointer;font-family:var(--font-mono);font-size:12px;white-space:nowrap}
.tab.active{border-color:var(--accent);background:var(--accent-soft);color:var(--accent);font-weight:700}
.brand{display:flex;align-items:center;gap:10px;font-weight:600;font-size:15px}
.dot{width:10px;height:10px;border-radius:50%;background:var(--ok);animation:pulse 1.6s infinite}
@keyframes pulse{0%,100%{transform:scale(1);opacity:1}50%{transform:scale(1.4);opacity:.45}}
.grid{flex:1;min-height:0;display:grid;grid-template-columns:380px 1fr;gap:12px}
.card{display:flex;flex-direction:column;min-height:0;background:var(--surface);border:1px solid var(--border);border-radius:10px;padding:14px;box-shadow:var(--shadow)}
.card-hd{flex-shrink:0;display:flex;justify-content:space-between;align-items:center;margin-bottom:10px}
.card-hd h3{margin:0;font-size:12px;text-transform:uppercase;letter-spacing:.05em;color:var(--muted)}
.scroll{flex:1;min-height:0;overflow-y:auto;overflow-x:hidden;padding-right:4px}
.grp{font-size:11px;font-weight:700;color:var(--muted);text-transform:uppercase;margin:10px 0 6px;display:flex;justify-content:space-between}
.grp:first-child{margin-top:0}
.w{display:flex;align-items:center;gap:10px;padding:8px 10px;border-radius:8px;margin-bottom:6px;background:var(--surface-2);cursor:pointer;border:1.5px solid var(--border);transition:border-color .15s,opacity .15s,filter .15s}
.w.waiting{filter:grayscale(1);opacity:.65}
.w:hover{border-color:var(--accent)}.w.sel{border-color:var(--accent);background:var(--accent-soft);opacity:1;filter:none}
.bot{width:26px;height:26px;flex-shrink:0;shape-rendering:crispEdges}
.bot.working{animation:work .55s steps(2) infinite}.bot.killed{opacity:.35;filter:grayscale(1)}
@keyframes work{0%{transform:translate(0,0)}50%{transform:translate(2px,-3px)}100%{transform:translate(-1px,0)}}
.w-info{flex:1;min-width:0}.row{display:flex;justify-content:space-between;align-items:center;gap:6px}.row b{white-space:nowrap;overflow:hidden;text-overflow:ellipsis;min-width:0}
.badge{font-family:var(--font-mono);font-size:10px;padding:2px 7px;border-radius:999px;font-weight:600;white-space:nowrap;flex-shrink:0}
.waiting{background:var(--accent-soft);color:var(--accent)}.working{background:var(--ok-soft);color:var(--ok)}.killed{background:var(--fail-soft);color:var(--fail)}
.meta{color:var(--muted);font-size:11px;font-family:var(--font-mono);margin-top:2px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.msg{padding:10px;border-bottom:1px solid var(--border);animation:fadeIn .2s ease}
@keyframes fadeIn{from{opacity:0;transform:translateY(-4px)}to{opacity:1;transform:none}}
.mention{color:var(--accent);font-weight:600;font-family:var(--font-mono)}
.crumb{display:flex;align-items:center;gap:6px;font-size:12px;font-family:var(--font-mono);color:var(--muted)}
.crumb a{color:var(--accent);cursor:pointer;text-decoration:none;font-weight:600}
.sub-btn{font-family:var(--font-mono);font-size:10px;padding:2px 7px;border-radius:6px;background:var(--surface);border:1px solid var(--border);color:var(--accent);cursor:pointer;white-space:nowrap;flex-shrink:0}
.pager{flex-shrink:0;display:flex;justify-content:space-between;align-items:center;padding-top:10px;border-top:1px solid var(--border);font-size:12px;color:var(--muted)}
.btn{background:var(--surface-2);border:1px solid var(--border);color:var(--text);padding:4px 10px;border-radius:6px;cursor:pointer;font-size:12px}
.btn:disabled{opacity:.4;cursor:default}
</style></head><body>
<div class="app">
  <div class="top">
    <div class="brand"><span class="dot"></span><span>LegoLoop Swarm</span><span class="crumb" id="crumb"></span></div>
    <div class="meta" id="summary">Conectando SSE...</div>
  </div>
  <div class="tabs" id="tabs"></div>
  <div class="grid">
    <div class="card"><div class="card-hd"><h3>Agentes</h3><span class="meta" id="counts"></span></div><div class="scroll" id="workers"></div></div>
    <div class="card"><div class="card-hd"><h3>Execução <span id="filter"></span></h3><span class="meta">Mais recente no topo</span></div><div class="scroll" id="msgs"></div><div class="pager" id="pager"></div></div>
  </div>
</div>
<script>
let sel="", page=0, curParent="", curProc=""; const PER_PAGE=20;
const palette=['#8ab4f8','#f28b82','#fdd663','#81c995','#c58af9'];
function botSVG(st,idx){const c=palette[idx%palette.length];return '<svg class="bot '+st+'" viewBox="0 0 8 8"><rect x="3" y="0" width="2" height="1" fill="'+c+'"/><rect x="1" y="1" width="6" height="4" fill="'+c+'"/><rect x="2" y="2" width="1" height="1" fill="#1e293b"/><rect x="5" y="2" width="1" height="1" fill="#1e293b"/><rect x="2" y="5" width="4" height="2" fill="'+c+'"/><rect x="2" y="7" width="1" height="1" fill="'+c+'"/><rect x="5" y="7" width="1" height="1" fill="'+c+'"/></svg>';}
function setProc(p){curProc=p;curParent="";sel="";page=0;if(window._last)render(window._last);}
function setSel(n){sel=(sel===n)?"":n;page=0;if(window._last)render(window._last);}
function enterSub(n,e){e.stopPropagation();curParent=n;sel="";page=0;if(window._last)render(window._last);}
function goPage(p){page=p;if(window._last)render(window._last);}
function renderCrumb(ws){
  const path=[]; let p=curParent;
  while(p){path.unshift(p);const w=ws.find(x=>x.name===p);p=w?(w.parent||""):"";}
  let html='<a onclick="enterSub(\'\',event)">Macro</a>';
  for(const n of path) html+=' › <a onclick="enterSub(\''+n+'\',event)">'+n+'</a>';
  document.getElementById('crumb').innerHTML=html;
}
function renderGroup(title,list,offset,allWs){
  if(!list.length)return '';
  return '<div class="grp"><span>'+title+'</span><span>'+list.length+'</span></div>'+list.map((w,i)=>{
    const subs=allWs.filter(x=>x.parent===w.name).length;
    const subBtn=subs?'<button class="sub-btn" onclick="enterSub(\''+w.name+'\',event)">'+subs+' sub ▶</button>':'';
    const bStyle=w.status==='working'?' style="border-color:'+palette[(offset+i)%palette.length]+'"':'';
    return '<div class="w '+w.status+' '+(sel===w.name?'sel':'')+'"'+bStyle+' onclick="setSel(\''+w.name+'\')">'+botSVG(w.status,offset+i)+
    '<div class="w-info"><div class="row"><b>'+w.name+'</b><span style="display:flex;gap:4px;align-items:center">'+subBtn+'<span class="badge '+w.status+'">'+w.status+'</span></span></div>'+
    '<div class="meta">'+w.role+' · '+(w.model||'-')+' · '+w.workspace_mode+'</div></div></div>';
  }).join('');
}
function render(payload){
  window._last=payload;
  const order=payload.order||[];
  if(!curProc||!payload.processes[curProc])curProc=order[0]||"";
  document.getElementById('tabs').innerHTML=order.map(p=>
    '<button class="tab '+(p===curProc?'active':'')+'" onclick="setProc(\''+p+'\')">'+p+'</button>'
  ).join('');
  const d=(payload.processes&&payload.processes[curProc])||{workers:[],messages:[]};
  const allWs=d.workers||[], ms=(d.messages||[]).slice().reverse();
  renderCrumb(allWs);
  const ws=allWs.filter(w=>(w.parent||"")===curParent);
  const ver=ws.filter(w=>w.role==='auditor-qualidade');
  const trab=ws.filter(w=>w.role!=='auditor-qualidade');
  document.getElementById('summary').textContent='1 coordenador · '+ver.length+' verificador(es) · '+trab.length+' trabalhador(es) · '+ms.length+' posts';
  document.getElementById('counts').textContent=(ws.length+1)+' neste nível';
  document.getElementById('filter').textContent=sel?'(filtro: '+sel+')':'(todos)';
  const coord=[{name:curParent||'coordenador',role:'orquestrador',model:'maestro',workspace_mode:'inherit',status:'working'}];
  document.getElementById('workers').innerHTML=renderGroup('Coordenador',coord,0,allWs)+renderGroup('Verificadores',ver,1,allWs)+renderGroup('Trabalhadores',trab,1+ver.length,allWs);
  const fms=ms.filter(m=>!sel||m.from===sel||m.to==='@'+sel||m.to==='@all');
  const totalPages=Math.max(1,Math.ceil(fms.length/PER_PAGE));
  if(page>=totalPages)page=totalPages-1;
  const slice=fms.slice(page*PER_PAGE,(page+1)*PER_PAGE);
  document.getElementById('msgs').innerHTML=slice.map(m=>
    '<div class="msg"><div class="meta">#'+m.seq+' <b>'+m.from+'</b> → <span class="mention">'+(m.to||'@all')+'</span></div>'+
    '<div style="margin-top:4px">'+m.body+'</div></div>'
  ).join('');
  document.getElementById('pager').innerHTML='<button class="btn" onclick="goPage('+(page-1)+')" '+(page<=0?'disabled':'')+'>Anterior</button><span>Página '+(page+1)+' de '+totalPages+'</span><button class="btn" onclick="goPage('+(page+1)+')" '+(page>=totalPages-1?'disabled':'')+'>Próxima</button>';
}
const es=new EventSource('/api/events');
es.onmessage=e=>render(JSON.parse(e.data));
</script></body></html>`
