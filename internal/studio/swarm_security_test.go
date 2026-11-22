package studio

import (
	"context"
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/heverton-dev/LegoLoop/internal/board"
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
		for _, legacy := range []string{"static/swarm.html", "static/swarm.js", "static/swarm.css"} {
			if _, err := fs.Stat(staticFS, legacy); err == nil {
				t.Fatalf("arquivo legado %s deveria ter sido removido do staticFS", legacy)
			}
		}
		html := readStatic(t, "sessao.html")
		for _, want := range []string{"/components/sessao-app.js", "/components/sessao.css"} {
			if !strings.Contains(html, want) {
				t.Fatalf("sessao.html deve carregar %s", want)
			}
		}
		appJS := readStatic(t, "components/sessao-app.js")
		for _, want := range []string{"/components/a2ui-bundle.js", "/components/legoloop-catalog.js"} {
			if !strings.Contains(appJS, want) {
				t.Fatalf("sessao-app.js deve carregar %s", want)
			}
		}
		for _, jsFile := range []string{"components/sessao-app.js", "components/legoloop-catalog.js", "components/a2ui-bundle.js"} {
			js := readStatic(t, jsFile)
			for i, line := range strings.Split(js, "\n") {
				trimmed := strings.TrimSpace(line)
				if strings.HasPrefix(trimmed, "//") {
					continue
				}
				if strings.Contains(line, "innerHTML") || strings.Contains(line, "insertAdjacentHTML") {
					t.Fatalf("%s:%d usa HTML dinamico proibido: %s", jsFile, i+1, trimmed)
				}
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

		// 1. GET /swarm cria sessao A2UI e redireciona 303 para /sessao/{id}
		recPage := do(t, h, http.MethodGet, "/swarm", "")
		loc := recPage.Header().Get("Location")
		if recPage.Code != http.StatusSeeOther || !strings.HasPrefix(loc, "/sessao/") {
			t.Fatalf("GET /swarm status=%d location=%q body=%s", recPage.Code, loc, recPage.Body.String())
		}
		sessID := strings.TrimPrefix(loc, "/sessao/")
		sess, err := s.A2UIHub().GetSession(sessID)
		if err != nil || sess == nil {
			t.Fatalf("sessao A2UI %q nao encontrada no hub: %v", sessID, err)
		}
		recSess := do(t, h, http.MethodGet, loc, "")
		if recSess.Code != http.StatusOK || !strings.Contains(recSess.Body.String(), "LegoLoop A2UI") {
			t.Fatalf("GET %s status=%d body=%s", loc, recSess.Code, recSess.Body.String())
		}
		for _, removed := range []string{"/swarm.js", "/swarm.css"} {
			if recOld := do(t, h, http.MethodGet, removed, ""); recOld.Code != http.StatusNotFound {
				t.Fatalf("GET %s deveria retornar 404 apos remocao, obteve %d", removed, recOld.Code)
			}
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
const fs = require("fs"), vm = require("vm");
function makeEl(tag, id="") {
  return {
    tagName: tag.toUpperCase(), id, className: "", textContent: "", value: "", disabled: false, hidden: false,
    type: "", max: 0, htmlFor: "", listeners: {}, dataset: {}, attributes: {}, children: [], style: {},
    setAttribute(k, v) { this.attributes[k] = String(v); },
    getAttribute(k) { return this.attributes[k] || null; },
    append(...nodes) { for (const n of nodes) if (n) this.children.push(typeof n === "string" ? { tagName: "#TEXT", textContent: n, children: [] } : n); },
    appendChild(n) { if (n) this.children.push(n); return n; },
    replaceChildren(...nodes) { this.children = []; this.append(...nodes); },
    addEventListener(ev, fn) { (this.listeners[ev] ||= []).push(fn); },
    dispatch(ev, payload={}) { for (const fn of (this.listeners[ev] || [])) fn(payload); }
  };
}
function findById(root, id) {
  if (root.id === id) return root;
  for (const c of (root.children || [])) { const f = findById(c, id); if (f) return f; }
  return null;
}
function findByTag(root, tag) {
  if (root.tagName === tag) return root;
  for (const c of (root.children || [])) { const f = findByTag(c, tag); if (f) return f; }
  return null;
}
function collectTags(root, out=[]) {
  if (root.tagName) out.push(root.tagName);
  for (const c of (root.children || [])) collectTags(c, out);
  return out;
}
const bundleCode = fs.readFileSync("static/components/a2ui-bundle.js", "utf8").replace(/\bexport\s+/g, "");
const catalogCode = fs.readFileSync("static/components/legoloop-catalog.js", "utf8")
  .replace(/^\s*import\s+[^\n]+;\s*$/gm, "")
  .replace(/\bexport\s+/g, "");
const doc = {
  createElement(tag) { return makeEl(tag); },
  createTextNode(txt) { return { tagName: "#TEXT", textContent: String(txt), children: [] }; }
};
const ctx = {
  document: doc, console, Date, Math, Number, String, Boolean, Array, Object, JSON, Map, Set, RegExp, Error,
  fetch: async () => ({ ok: true, status: 200, text: async () => "" })
};
vm.createContext(ctx);
vm.runInContext(bundleCode + "\n" + catalogCode + "\nthis.MessageProcessor = MessageProcessor; this.A2UISurfaceRenderer = A2UISurfaceRenderer;", ctx);

const processor = new ctx.MessageProcessor();
let xssRejected = false;
try {
  processor.process({
    session_id: "sess_1", surface_id: "proc:swarm", message_type: "createSurface",
    components: [{ id: "swarm-board", kind: "LegoSwarmBoard", props: { messages: [{ seq: 1, agent: "w1", text: "<img src=x onerror='window.__xssProbe=1'>" }] } }]
  });
} catch (_e) { xssRejected = true; }
if (!xssRejected) throw new Error("payload XSS nao foi rejeitado pelo MessageProcessor");

const env = {
  session_id: "sess_1", surface_id: "proc:swarm", message_type: "createSurface",
  components: [{
    id: "swarm-board", kind: "LegoSwarmBoard",
    props: {
      title: "LegoLoop Swarm", progress: 50, action: "swarm_post",
      workers: [{ name: "w1", role: "engenheiro-backend", state: "done" }],
      messages: [{ seq: 1, from: "w1", text: "<img src=x>" }]
    }
  }]
};
processor.process(env);
const actions = [];
const root = makeEl("div", "a2ui-root");
const renderer = new ctx.A2UISurfaceRenderer(root, (act) => actions.push(act));
renderer.render(processor.getSurface("proc:swarm"));

if (collectTags(root).includes("IMG")) throw new Error("tag IMG injetada no DOM do LegoSwarmBoard");
const draftIn = findById(root, "swarm-draft-in");
if (!draftIn) throw new Error("#swarm-draft-in nao encontrado");
draftIn.value = "rascunho em andamento";
draftIn.dispatch("input");

// Re-render SSE deve preservar o rascunho digitado
processor.process(env);
renderer.render(processor.getSurface("proc:swarm"));
const draftAfter = findById(root, "swarm-draft-in");
if (!draftAfter || draftAfter.value !== "rascunho em andamento") {
  throw new Error("SSE destruiu rascunho em #swarm-draft-in: " + (draftAfter && draftAfter.value));
}

// Submissao envia acao swarm_post e limpa rascunho
const form = findByTag(root, "FORM");
if (!form) throw new Error("formulario do LegoSwarmBoard nao encontrado");
form.dispatch("submit", { preventDefault() {} });
if (actions.length !== 1 || actions[0].name !== "swarm_post" || actions[0].payload.text !== "rascunho em andamento") {
  throw new Error("acao swarm_post inesperada: " + JSON.stringify(actions));
}
if (draftAfter.value !== "") throw new Error("rascunho nao foi limpo apos envio");
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
