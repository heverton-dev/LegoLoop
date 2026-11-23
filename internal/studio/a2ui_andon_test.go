package studio

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/heverton-dev/LegoLoop/internal/a2ui"
	"github.com/heverton-dev/LegoLoop/internal/triagem"
)

func TestAndonCordPauseResumeCancel(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	sess := s.A2UIHub().EnsureSession("sess_teste_andon", "legoloop-fabricar")
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava estado running, obteve %s", sess.GetState())
	}

	actionURL := "/api/sessao/" + sess.ID + "/action"
	pRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_pause"}`)
	if pRec.Code != http.StatusAccepted {
		t.Fatalf("POST pause code=%d body=%s", pRec.Code, pRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionPaused {
		t.Fatalf("esperava estado paused, obteve %s", sess.GetState())
	}

	rRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_resume"}`)
	if rRec.Code != http.StatusAccepted {
		t.Fatalf("POST resume code=%d body=%s", rRec.Code, rRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava estado running, obteve %s", sess.GetState())
	}

	cRec := do(t, h, http.MethodPost, actionURL, `{"surface_id":"andon","origin":"ui","name":"andon_cord_cancel"}`)
	if cRec.Code != http.StatusAccepted {
		t.Fatalf("POST cancel code=%d body=%s", cRec.Code, cRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionCancelled {
		t.Fatalf("esperava estado cancelled, obteve %s", sess.GetState())
	}
}

func TestProcessosAtivosAPI(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	sess := s.A2UIHub().EnsureSession("sess_proc_api", "legoloop-fabricar")

	pageRec := do(t, h, http.MethodGet, "/processos", "")
	if pageRec.Code != http.StatusOK || !strings.Contains(pageRec.Body.String(), "Gestao de Processos") {
		t.Fatalf("esperava 200 na página /processos, obteve %d", pageRec.Code)
	}

	listRec := do(t, h, http.MethodGet, "/api/processos/ativos", "")
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET /api/processos/ativos code=%d", listRec.Code)
	}
	var itens []ProcessoAtivoItem
	if err := json.Unmarshal(listRec.Body.Bytes(), &itens); err != nil {
		t.Fatalf("unmarshal itens: %v", err)
	}
	achou := false
	for _, it := range itens {
		if it.ID == sess.ID && it.Processo == "legoloop-fabricar" {
			achou = true
			break
		}
	}
	if !achou {
		t.Fatalf("sessão %s não listada nos processos ativos", sess.ID)
	}

	andRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/andon", "{}")
	if andRec.Code != http.StatusOK {
		t.Fatalf("POST andon code=%d body=%s", andRec.Code, andRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionPaused {
		t.Fatalf("esperava paused após andon, obteve %s", sess.GetState())
	}

	resRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/resume", "{}")
	if resRec.Code != http.StatusOK {
		t.Fatalf("POST resume code=%d body=%s", resRec.Code, resRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionRunning {
		t.Fatalf("esperava running após resume, obteve %s", sess.GetState())
	}

	canRec := do(t, h, http.MethodPost, "/api/processos/"+sess.ID+"/cancel", "{}")
	if canRec.Code != http.StatusOK {
		t.Fatalf("POST cancel code=%d body=%s", canRec.Code, canRec.Body.String())
	}
	if sess.GetState() != a2ui.SessionCancelled {
		t.Fatalf("esperava cancelled após cancel, obteve %s", sess.GetState())
	}
}

func TestProcessosAtivosDiscoEAndon(t *testing.T) {
	s, _ := newServer(t)
	h := s.Handler()

	procDir := filepath.Join(s.ProjectDir, "processos", "proc-disco")
	if err := os.MkdirAll(procDir, 0o755); err != nil {
		t.Fatalf("mkdir procDir: %v", err)
	}

	memDisco := triagem.MemoriaExecucao{
		Processo:     "proc-disco",
		SessaoID:     "sess_disco_1",
		PID:          os.Getpid(),
		EtapaAtual:   "executando",
		ProximoPasso: "etapa_2",
	}
	if err := triagem.SalvarMemoria(procDir, memDisco); err != nil {
		t.Fatalf("salvar memoria: %v", err)
	}

	s.A2UIHub().EnsureSession("sess_hub_1", "legoloop-fabricar")

	testarListagemDiscoEHub(t, h)
	testarAndonDiscoPause(t, h, procDir, "sess_disco_1")
	testarAndonDiscoResume(t, h, procDir, "sess_disco_1")
	testarAndonDiscoCancel(t, h, procDir, "sess_disco_1")
}

func testarListagemDiscoEHub(t *testing.T, h http.Handler) {
	rec := do(t, h, http.MethodGet, "/api/processos/ativos", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/processos/ativos code=%d", rec.Code)
	}
	var itens []ProcessoAtivoItem
	if err := json.Unmarshal(rec.Body.Bytes(), &itens); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	vistos := make(map[string]int)
	for _, it := range itens {
		vistos[it.ID]++
	}
	if vistos["sess_disco_1"] != 1 {
		t.Fatalf("esperava 1 ocorrencia de sess_disco_1, obteve %d", vistos["sess_disco_1"])
	}
	if vistos["sess_hub_1"] != 1 {
		t.Fatalf("esperava 1 ocorrencia de sess_hub_1, obteve %d", vistos["sess_hub_1"])
	}
}

func testarAndonDiscoPause(t *testing.T, h http.Handler, dir, sessID string) {
	rec := do(t, h, http.MethodPost, "/api/processos/"+sessID+"/andon", "{}")
	if rec.Code != http.StatusOK {
		t.Fatalf("POST andon code=%d body=%s", rec.Code, rec.Body.String())
	}
	mem, err := triagem.LerMemoria(dir, sessID)
	if err != nil {
		t.Fatalf("ler memoria pause: %v", err)
	}
	if mem.EtapaAtual != "pausado_andon_cord" {
		t.Fatalf("etapa esperada pausado_andon_cord, obteve %s", mem.EtapaAtual)
	}
}

func testarAndonDiscoResume(t *testing.T, h http.Handler, dir, sessID string) {
	rec := do(t, h, http.MethodPost, "/api/processos/"+sessID+"/resume", "{}")
	if rec.Code != http.StatusOK {
		t.Fatalf("POST resume code=%d body=%s", rec.Code, rec.Body.String())
	}
	mem, err := triagem.LerMemoria(dir, sessID)
	if err != nil {
		t.Fatalf("ler memoria resume: %v", err)
	}
	if mem.EtapaAtual != "executando" {
		t.Fatalf("etapa esperada executando, obteve %s", mem.EtapaAtual)
	}
}

func testarAndonDiscoCancel(t *testing.T, h http.Handler, dir, sessID string) {
	rec := do(t, h, http.MethodPost, "/api/processos/"+sessID+"/cancel", "{}")
	if rec.Code != http.StatusOK {
		t.Fatalf("POST cancel code=%d body=%s", rec.Code, rec.Body.String())
	}
	if _, err := triagem.LerMemoria(dir, sessID); err == nil {
		t.Fatalf("esperava sessão removida do disco após cancel")
	}
}
