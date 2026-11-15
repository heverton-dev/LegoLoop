package studio

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func createHardeningLink(t *testing.T, target, link string) {
	t.Helper()
	if err := os.Symlink(target, link); err == nil {
		return
	}
	if out, err := exec.Command("cmd", "/c", "mklink", "/J", link, target).CombinedOutput(); err != nil {
		t.Skipf("symlink/junction indisponível no ambiente: %v (%s)", err, string(out))
	}
}

type sseSpyWriter struct {
	mu             sync.Mutex
	header         http.Header
	code           int
	buf            strings.Builder
	activeWrites   atomic.Int32
	raceDetected   atomic.Bool
	deadlineSet    atomic.Bool
	lastDeadline   time.Time
	deadlineMu     sync.Mutex
	onWriteOrFlush func()
}

func newSSESpyWriter() *sseSpyWriter {
	return &sseSpyWriter{header: make(http.Header)}
}

func (w *sseSpyWriter) Header() http.Header { return w.header }

func (w *sseSpyWriter) WriteHeader(statusCode int) { w.code = statusCode }

func (w *sseSpyWriter) Write(p []byte) (int, error) {
	if cur := w.activeWrites.Add(1); cur > 1 {
		w.raceDetected.Store(true)
	}
	defer w.activeWrites.Add(-1)
	if w.onWriteOrFlush != nil {
		w.onWriteOrFlush()
	}
	time.Sleep(2 * time.Millisecond)
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.Write(p)
}

func (w *sseSpyWriter) Flush() {
	if cur := w.activeWrites.Add(1); cur > 1 {
		w.raceDetected.Store(true)
	}
	defer w.activeWrites.Add(-1)
	if w.onWriteOrFlush != nil {
		w.onWriteOrFlush()
	}
}

func (w *sseSpyWriter) SetWriteDeadline(t time.Time) error {
	w.deadlineMu.Lock()
	defer w.deadlineMu.Unlock()
	w.deadlineSet.Store(true)
	w.lastDeadline = t
	return nil
}

func TestPathJailDeleteOriginAndSSERace(t *testing.T) {
	t.Run("RevealPathJail", func(t *testing.T) {
		proj := t.TempDir()
		vault := t.TempDir()
		outside := t.TempDir()

		inProj := filepath.Join(proj, "sub", "peca.md")
		inVault := filepath.Join(vault, "skills", "s1", "SKILL.md")
		outFile := filepath.Join(outside, "segredo.txt")
		for _, p := range []string{inProj, inVault, outFile} {
			if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(p, []byte("conteudo"), 0o644); err != nil {
				t.Fatal(err)
			}
		}

		s := &Server{ProjectDir: proj, Vault: vault, Host: host}

		if got, isDir, err := s.resolveSafePath("<projeto>/sub/peca.md"); err != nil || isDir || got != inProj {
			t.Fatalf("esperava sucesso para <projeto>/sub/peca.md: got=%q isDir=%v err=%v", got, isDir, err)
		}
		if got, isDir, err := s.resolveSafePath(inVault); err != nil || isDir || got != inVault {
			t.Fatalf("esperava sucesso para caminho no cofre: got=%q isDir=%v err=%v", got, isDir, err)
		}
		if _, _, err := s.resolveSafePath(outFile); err == nil {
			t.Fatalf("path jail deveria bloquear caminho fora do projeto/cofre: %s", outFile)
		}

		linkPath := filepath.Join(proj, "escape-link")
		createHardeningLink(t, outside, linkPath)
		if _, _, err := s.resolveSafePath(filepath.Join(linkPath, "segredo.txt")); err == nil {
			t.Fatalf("path jail deveria bloquear symlink/junction apontando para fora")
		}

		s.RevealLauncher = func(string, bool) error { return nil }
		h := s.Handler()
		rec := do(t, h, http.MethodPost, "/api/assets/reveal", `{"path":"<projeto>/sub/peca.md"}`)
		if rec.Code != http.StatusOK {
			t.Fatalf("POST /api/assets/reveal esperava 200, obteve %d: %s", rec.Code, rec.Body.String())
		}
		var revealResp map[string]any
		decode(t, rec, &revealResp)
		respPath, _ := revealResp["path"].(string)
		if respPath == "" || filepath.IsAbs(respPath) || strings.Contains(respPath, proj) || strings.Contains(respPath, vault) {
			t.Fatalf("POST /api/assets/reveal vazou caminho absoluto ou inválido: %q (body=%s)", respPath, rec.Body.String())
		}
	})

	t.Run("EliminatePiecePathJailAndToggleFavorite", func(t *testing.T) {
		proj := t.TempDir()
		vault := t.TempDir()
		outside := t.TempDir()
		prefixEvil := proj + "-evil"
		if err := os.MkdirAll(prefixEvil, 0o755); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.RemoveAll(prefixEvil) })

		outsideFile := filepath.Join(outside, "protegido.txt")
		prefixFile := filepath.Join(prefixEvil, "alvo.md")
		projRootCheck := t.TempDir()
		rootLevelSkill := filepath.Join(projRootCheck, "skill.md")
		if err := os.WriteFile(rootLevelSkill, []byte("dados"), 0o644); err != nil {
			t.Fatal(err)
		}
		sRoot, _ := newServer(t)
		sRoot.ProjectDir = projRootCheck
		sRoot.Vault = vault
		if st, _ := sRoot.eliminatePiece("full", rootLevelSkill, "outro-id"); st == 0 {
			t.Fatalf("eliminatePiece full deveria impedir remoção da raiz do projeto")
		}
		if _, err := os.Stat(projRootCheck); err != nil {
			t.Fatalf("raiz do projeto foi removida: %v", err)
		}

		for _, p := range []string{outsideFile, prefixFile} {
			if err := os.WriteFile(p, []byte("dados"), 0o644); err != nil {
				t.Fatal(err)
			}
		}

		s, _ := newServer(t)
		s.ProjectDir = proj
		s.Vault = vault

		if st, _ := s.eliminatePiece("full", outsideFile, "protegido.txt"); st == 0 {
			t.Fatalf("eliminatePiece full deveria bloquear caminho fora do projeto/cofre")
		}
		if _, err := os.Stat(outsideFile); err != nil {
			t.Fatalf("arquivo fora do projeto/cofre foi apagado: %v", err)
		}

		if st, _ := s.eliminatePiece("local", prefixFile, "alvo.md"); st == 0 {
			t.Fatalf("eliminatePiece local deveria bloquear bypass por prefixo (%s)", prefixFile)
		}
		if _, err := os.Stat(prefixFile); err != nil {
			t.Fatalf("arquivo fora do projeto com prefixo semelhante foi apagado: %v", err)
		}

		linkPath := filepath.Join(proj, "atalho")
		createHardeningLink(t, outside, linkPath)
		if st, _ := s.eliminatePiece("full", filepath.Join(linkPath, "protegido.txt"), "protegido.txt"); st == 0 {
			t.Fatalf("eliminatePiece full deveria bloquear symlink/junction")
		}
		if _, err := os.Stat(outsideFile); err != nil {
			t.Fatalf("alvo de symlink fora do projeto foi apagado: %v", err)
		}

		h := s.Handler()
		const workers = 16
		var wg sync.WaitGroup
		for i := range workers {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				path := "/api/assets/skill/peca-" + string(rune('a'+idx)) + "/favorite"
				rec := do(t, h, http.MethodPost, path, "{}")
				if rec.Code != http.StatusOK {
					t.Errorf("favorite %d retornou %d: %s", idx, rec.Code, rec.Body.String())
				}
			}(i)
		}
		wg.Wait()

		favRaw, err := os.ReadFile(filepath.Join(vault, "favorites.json"))
		if err != nil {
			t.Fatalf("falha ao ler favorites.json: %v", err)
		}
		var favMap map[string]bool
		if err := json.Unmarshal(favRaw, &favMap); err != nil {
			t.Fatalf("favorites.json corrompido após concorrência: %v (%s)", err, string(favRaw))
		}
		if len(favMap) != workers {
			t.Fatalf("lost update em toggleFavorite: esperava %d favoritos, obteve %d (%v)", workers, len(favMap), favMap)
		}

		corruptPath := filepath.Join(vault, "favorites.json")
		if err := os.WriteFile(corruptPath, []byte("{corrompido"), 0o644); err != nil {
			t.Fatal(err)
		}
		recErr := do(t, h, http.MethodPost, "/api/assets/skill/nova/favorite", "{}")
		if recErr.Code == http.StatusOK {
			t.Fatalf("toggleFavorite deveria retornar erro com favorites.json corrompido, retornou 200")
		}
	})

	t.Run("MutatingMethodsRequireOriginAndSecFetchSite", func(t *testing.T) {
		s, _ := newServer(t)
		h := s.Handler()
		methods := []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
		for _, method := range methods {
			recOrigin := do(t, h, method, "/api/processos/qualquer", "", func(r *http.Request) {
				r.Header.Set("Origin", "https://evil.example.com")
			})
			if recOrigin.Code != http.StatusForbidden {
				t.Errorf("%s com Origin cruzado esperava 403, obteve %d", method, recOrigin.Code)
			}
			recSite := do(t, h, method, "/api/processos/qualquer", "", func(r *http.Request) {
				r.Header.Set("Sec-Fetch-Site", "cross-site")
			})
			if recSite.Code != http.StatusForbidden {
				t.Errorf("%s com Sec-Fetch-Site=cross-site esperava 403, obteve %d", method, recSite.Code)
			}
		}
	})

	t.Run("SSEMutexDeadlineAndDAGStatus", func(t *testing.T) {
		spy := newSSESpyWriter()
		startSSE(spy, spy)
		if !spy.deadlineSet.Load() || !spy.lastDeadline.IsZero() {
			t.Fatalf("startSSE deveria chamar ResponseController.SetWriteDeadline(time.Time{}), set=%v deadline=%v",
				spy.deadlineSet.Load(), spy.lastDeadline)
		}

		sender := newSSESender(spy, spy)
		var wg sync.WaitGroup
		for i := range 12 {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				sender("node_start", map[string]any{"idx": id})
			}(i)
		}
		wg.Wait()
		if spy.raceDetected.Load() {
			t.Fatalf("data race detectado: múltiplas goroutines escreveram/flusharam SSE simultaneamente")
		}

		scanner := bufio.NewScanner(strings.NewReader(spy.buf.String()))
		dataLines := 0
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				dataLines++
				var ev map[string]any
				if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "data: ")), &ev); err != nil {
					t.Fatalf("linha SSE corrompida por concorrência: %q (%v)", line, err)
				}
			}
		}
		if dataLines != 12 {
			t.Fatalf("esperava 12 eventos SSE válidos, obteve %d", dataLines)
		}

		s, _ := newServer(t)
		h := s.Handler()
		body := `{
			"dag": {
				"name": "status-dag",
				"version": "2.0.0",
				"trigger": {"type": "manual"},
				"nodes": [{
					"id": "n1",
					"type": "script",
					"command": ["go", "version"],
					"routing": {"mode": "manual", "model": "none"},
					"output": {"artifact": "out.txt", "gate": "file_exists('out.txt')"}
				}]
			},
			"allow_scripts": true
		}`
		recRun := do(t, h, http.MethodPost, "/api/run", body)
		if recRun.Code != http.StatusOK {
			t.Fatalf("POST /api/run falhou: %d %s", recRun.Code, recRun.Body.String())
		}

		recStatus := do(t, h, http.MethodGet, "/api/dag", "")
		if recStatus.Code != http.StatusOK {
			t.Fatalf("GET /api/dag falhou: %d", recStatus.Code)
		}
		var st map[string]any
		decode(t, recStatus, &st)
		if st["status"] == "idle" {
			t.Fatalf("GET /api/dag ainda retorna mock idle após execução real: %v", st)
		}
		if st["dag"] != "status-dag" {
			t.Fatalf("GET /api/dag deveria rastrear dag='status-dag', obteve: %v", st)
		}
		if strings.Contains(recStatus.Body.String(), s.ProjectDir) {
			t.Fatalf("GET /api/dag vazou caminho absoluto do projeto: %s", recStatus.Body.String())
		}

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		srv := &http.Server{
			WriteTimeout: 40 * time.Millisecond,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fl := w.(http.Flusher)
				startSSE(w, fl)
				send := newSSESender(w, fl)
				time.Sleep(80 * time.Millisecond)
				send("dag_finish", map[string]any{"ok": true})
			}),
		}
		go func() { _ = srv.Serve(ln) }()
		defer srv.Shutdown(context.Background())

		resp, err := http.Get("http://" + ln.Addr().String() + "/sse")
		if err != nil {
			t.Fatalf("falha na conexão SSE com WriteTimeout curto: %v", err)
		}
		defer resp.Body.Close()
		var gotFinish bool
		sc := bufio.NewScanner(resp.Body)
		for sc.Scan() {
			if strings.Contains(sc.Text(), `"dag_finish"`) {
				gotFinish = true
			}
		}
		if err := sc.Err(); err != nil || !gotFinish {
			t.Fatalf("SSE foi interrompido pelo WriteTimeout (err=%v, gotFinish=%v)", err, gotFinish)
		}
	})
}
