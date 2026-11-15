package studio

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"legoloop/internal/board"
)

func TestSwarmUIUnifiedNoXSSAndRealKill(t *testing.T) {
	t.Run("CLIWithoutDuplicateHTMLAndUsesStudioSecure", func(t *testing.T) {
		b, err := os.ReadFile(filepath.Join("..", "cli", "swarm_ui.go"))
		if err != nil {
			t.Fatalf("ler internal/cli/swarm_ui.go: %v", err)
		}
		src := string(b)
		if strings.Contains(src, "<!doctype html>") || strings.Contains(src, "swarmUIHTML") {
			t.Fatalf("internal/cli/swarm_ui.go ainda duplica HTML embutido")
		}
		if !strings.Contains(src, "studio.NewSwarmServer") && !strings.Contains(src, "&studio.Server") {
			t.Fatalf("internal/cli/swarm_ui.go deve delegar ao servidor de internal/studio")
		}
	})

	t.Run("StaticFilesAndZeroInnerHTML", func(t *testing.T) {
		html := readStatic(t, "swarm.html")
		if !strings.Contains(html, "/swarm.js") {
			t.Fatalf("swarm.html deve carregar /swarm.js")
		}
		js := readStatic(t, "swarm.js")
		for i, line := range strings.Split(js, "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "//") {
				continue
			}
			if strings.Contains(line, "innerHTML") || strings.Contains(line, "insertAdjacentHTML") {
				t.Fatalf("swarm.js:%d usa HTML dinamico proibido: %s", i+1, trimmed)
			}
		}
	})

	t.Run("UnifiedBoardDiscoveryAddSpawnModelPostKillAndCSRF", func(t *testing.T) {
		proj := initSwarmGitRepo(t)
		sb, err := board.OpenSwarmBoard(proj)
		if err != nil {
			t.Fatalf("OpenSwarmBoard: %v", err)
		}
		s, _ := newServer(t)
		s.ProjectDir = proj
		h := s.Handler()

		// 1. GET /swarm serve swarm.html
		recPage := do(t, h, http.MethodGet, "/swarm", "")
		if recPage.Code != http.StatusOK || !strings.Contains(recPage.Body.String(), "LegoLoop Swarm") {
			t.Fatalf("GET /swarm status=%d body=%s", recPage.Code, recPage.Body.String())
		}

		// 2. CSRF bloqueado pelo middleware secure()
		recCSRF := do(t, h, http.MethodPost, "/api/swarm/post", `{"proc":"`+filepath.Base(proj)+`","to":"@all","body":"x"}`, func(r *http.Request) {
			r.Header.Set("Origin", "http://evil.example")
		})
		if recCSRF.Code != http.StatusForbidden {
			t.Fatalf("CSRF com Origin externo deveria retornar 403, obteve %d", recCSRF.Code)
		}

		// 3. Add worker via HTTP (worktree em proj)
		procName := filepath.Base(proj)
		recAdd := do(t, h, http.MethodPost, "/api/swarm/add", `{"process":"`+procName+`","name":"w-sec","role":"engenheiro-backend","model":"m-init","workspace_mode":"branch"}`)
		if recAdd.Code != http.StatusOK {
			t.Fatalf("POST /api/swarm/add falhou: %d %s", recAdd.Code, recAdd.Body.String())
		}
		var added board.Worker
		decode(t, recAdd, &added)
		if added.Name != "w-sec" || added.WorktreePath == "" {
			t.Fatalf("worker adicionado invalido: %+v", added)
		}
		// Limpa RepoDir persistido para provar que handleSwarmKill passa s.ProjectDir ao git worktree prune/remove
		swarmJSONPath := filepath.Join(proj, ".legoloop-swarm", "swarm.json")
		if raw, err := os.ReadFile(swarmJSONPath); err == nil {
			cleaned := strings.ReplaceAll(string(raw), `"repo_dir"`, `"_ignored_repo_dir"`)
			_ = os.WriteFile(swarmJSONPath, []byte(cleaned), 0o644)
		}

		// 4. Troca de modelo via HTTP
		recModel := do(t, h, http.MethodPost, "/api/swarm/model", `{"process":"`+procName+`","name":"w-sec","model":"m-updated"}`)
		if recModel.Code != http.StatusOK {
			t.Fatalf("POST /api/swarm/model falhou: %d %s", recModel.Code, recModel.Body.String())
		}

		// 5. Spawn de processo real via HTTP
		cmdJSON := `["sh","-c","sleep 5"]`
		if runtime.GOOS == "windows" {
			cmdJSON = `["cmd","/c","ping -n 6 127.0.0.1 >nul"]`
		}
		recSpawn := do(t, h, http.MethodPost, "/api/swarm/spawn", `{"process":"`+procName+`","name":"w-sec","command":`+cmdJSON+`}`)
		if recSpawn.Code != http.StatusOK {
			t.Fatalf("POST /api/swarm/spawn falhou: %d %s", recSpawn.Code, recSpawn.Body.String())
		}

		// 6. Post de mensagem com payload XSS via HTTP
		xssPayload := `<img src=x onerror='window.__xssProbe=1'>`
		recPost := do(t, h, http.MethodPost, "/api/swarm/post", `{"proc":"`+procName+`","to":"@w-sec","body":"`+xssPayload+`"}`)
		if recPost.Code != http.StatusOK {
			t.Fatalf("POST /api/swarm/post falhou: %d %s", recPost.Code, recPost.Body.String())
		}

		// 7. Kill via HTTP usa s.ProjectDir para remover o worktree e encerrar PID
		recKill := do(t, h, http.MethodPost, "/api/swarm/kill", `{"proc":"`+procName+`","name":"w-sec"}`)
		if recKill.Code != http.StatusOK {
			t.Fatalf("POST /api/swarm/kill falhou: %d %s", recKill.Code, recKill.Body.String())
		}
		if _, err := os.Stat(added.WorktreePath); !os.IsNotExist(err) {
			t.Fatalf("worktree %s deveria ter sido removido usando s.ProjectDir", added.WorktreePath)
		}
		wtList, err := exec.Command("git", "-C", proj, "worktree", "list").CombinedOutput()
		if err != nil || strings.Contains(string(wtList), "w-sec") {
			t.Fatalf("git worktree list ainda contem w-sec (repoDir incorreto no KillWorker): %s", string(wtList))
		}
		st, err := sb.Status()
		if err != nil || len(st.Workers) != 1 || st.Workers[0].Status != board.WorkerKilled || st.Workers[0].Model != "m-updated" {
			t.Fatalf("estado inesperado pos-kill: %+v err=%v", st, err)
		}
	})

	t.Run("DOMStatePreservationReopenTabAndNoXSSInNode", func(t *testing.T) {
		verifySwarmJSDOMAndState(t)
	})
}

func initSwarmGitRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	runGit := func(args ...string) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v (%s)", args, err, string(out))
		}
	}
	runGit("init")
	runGit("config", "user.email", "test@legoloop.local")
	runGit("config", "user.name", "Tester")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# repo\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	runGit("add", "README.md")
	runGit("commit", "-m", "init")
	return dir
}

func verifySwarmJSDOMAndState(t *testing.T) {
	t.Helper()
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node ausente no PATH")
	}
	script := `
const vm = require("vm"), fs = require("fs");
const code = fs.readFileSync("static/swarm.js", "utf8");
function makeEl(tag, id="") {
  return {
    tagName: tag.toUpperCase(), id, className: "", textContent: "", value: "", disabled: false, hidden: false,
    dataset: {}, attributes: {}, children: [], style: {},
    setAttribute(k, v) { this.attributes[k] = String(v); if (k.startsWith("data-")) { const camel = k.slice(5).replace(/-([a-z])/g, (_, c) => c.toUpperCase()); this.dataset[camel] = String(v); } },
    getAttribute(k) { return this.attributes[k] || null; },
    append(...nodes) { for (const n of nodes) if (n) this.children.push(typeof n === "string" ? { tagName: "#TEXT", textContent: n, children: [] } : n); },
    appendChild(n) { if (n) this.children.push(n); return n; },
    replaceChildren(...nodes) { this.children = []; this.append(...nodes); },
    addEventListener() {},
    querySelector(sel) { return findIn(this, sel); }
  };
}
function findById(root, id) {
  if (root.id === id) return root;
  for (const c of (root.children || [])) { const f = findById(c, id); if (f) return f; }
  return null;
}
function findIn(root, sel) {
  if (sel.startsWith("#")) return findById(root, sel.slice(1));
  return null;
}
function collectTags(root, out=[]) {
  if (root.tagName) out.push(root.tagName);
  for (const c of (root.children || [])) collectTags(c, out);
  return out;
}
function collectText(root) {
  let s = root.textContent || "";
  for (const c of (root.children || [])) s += " " + collectText(c);
  return s;
}
const roots = {};
for (const id of ["crumb","summary","tabs","counts","workers","filter","ficha","msgs","pager","swarm-error"]) {
  roots[id] = makeEl("div", id);
}
let docClick = null;
const doc = {
  getElementById(id) {
    if (roots[id]) return roots[id];
    for (const r of Object.values(roots)) { const f = findById(r, id); if (f) return f; }
    return null;
  },
  createElement(tag) { return makeEl(tag); },
  createElementNS(_, tag) { return makeEl(tag); },
  createTextNode(txt) { return { tagName: "#TEXT", textContent: String(txt), children: [] }; },
  addEventListener(ev, fn) { if (ev === "click") docClick = fn; }
};
const win = {};
const ctx = {
  window: win, document: doc, console, Date, Math, Number, String, Boolean, Array, Object, JSON,
  EventSource: function() { return {}; },
  fetch: async () => ({ ok: true, status: 200, json: async () => ({}), text: async () => "" })
};
vm.createContext(ctx);
vm.runInContext(code, ctx);
const payload = {
  order: ["procA", "procB"],
  processes: {
    procA: {
      workers: [{ name: "w1", role: "engenheiro-backend", model: "m1", workspace_mode: "share", status: "working" }],
      messages: [{ seq: 1, from: "w1", to: "@all", body: "<img src=x onerror='window.__xssProbe=1'>" }]
    },
    procB: {
      workers: [{ name: "w2", role: "auditor-qualidade", model: "m2", workspace_mode: "inherit", status: "waiting" }],
      messages: []
    }
  }
};
ctx.render(payload);
// Selecionar w1 em procA
docClick({ target: { closest: () => ({ dataset: { act: "sel", val: "w1" } }) }, stopPropagation(){} });
const cmdIn1 = doc.getElementById("cmd-in");
if (!cmdIn1) throw new Error("#cmd-in nao encontrado ao selecionar w1");
cmdIn1.value = "rascunho em andamento";
// Re-render SSE deve preservar #cmd-in.value
ctx.render(payload);
const cmdInAfterSSE = doc.getElementById("cmd-in");
if (!cmdInAfterSSE || cmdInAfterSSE.value !== "rascunho em andamento") {
  throw new Error("SSE destruiu #cmd-in: " + (cmdInAfterSSE && cmdInAfterSSE.value));
}
// Trocar para procB nao deve vazar selectedAgent w1 de procA
docClick({ target: { closest: () => ({ dataset: { act: "proc", val: "procB" } }) }, stopPropagation(){} });
if (doc.getElementById("cmd-in") !== null) {
  throw new Error("selectedAgent vazou para procB");
}
// Voltar para procA restaura selectedAgent w1 e rascunho
docClick({ target: { closest: () => ({ dataset: { act: "proc", val: "procA" } }) }, stopPropagation(){} });
const cmdInRestored = doc.getElementById("cmd-in");
if (!cmdInRestored || cmdInRestored.value !== "rascunho em andamento") {
  throw new Error("estado da aba procA nao foi restaurado");
}
// Fechar procB e reabrir
docClick({ target: { closest: () => ({ dataset: { act: "close", val: "procB" } }) }, stopPropagation(){} });
docClick({ target: { closest: () => ({ dataset: { act: "reopen", val: "procB" } }) }, stopPropagation(){} });
const tabText = collectText(roots.tabs);
if (!tabText.includes("procB")) throw new Error("aba fechada procB nao reabriu: " + tabText);
// Verificar zero tag IMG criada no DOM de mensagens
const msgTags = collectTags(roots.msgs);
if (msgTags.includes("IMG") || win.__xssProbe) throw new Error("XSS executou ou gerou IMG");
console.log(JSON.stringify({ ok: true }));
`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", script).CombinedOutput()
	if err != nil {
		t.Fatalf("node DOM/state verification failed: %v\n%s", err, string(out))
	}
	var res map[string]bool
	if err := json.Unmarshal(out, &res); err != nil || !res["ok"] {
		t.Fatalf("resultado inesperado do node: %s", string(out))
	}
}
