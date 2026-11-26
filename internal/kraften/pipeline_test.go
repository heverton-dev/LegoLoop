package kraften

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// novoStudioSimulado cobre install/update, jobs, apps e health.
func novoStudioSimulado(t *testing.T, httpCode int) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"app01"}`))
	})
	mux.HandleFunc("/api/admin/apps/app01/update-zip", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"job":"job-1"}`))
	})
	mux.HandleFunc("/api/admin/jobs/job-1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"done"}`))
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(httpCode)
	})
	return httptest.NewServer(mux)
}

func novoCFSimulado(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success":true}`))
			return
		}
		w.Write([]byte(`{"result":[],"success":true}`))
	}))
}

func configDeTeste(t *testing.T) DeployConfig {
	t.Helper()
	return DeployConfig{
		AppID:      "app01",
		ZipPath:    criarZip(t, entradasSaudaveis()),
		Subdomain:  "sub01",
		TargetHost: "server.vpsconexao.org",
	}
}

func TestExecutePipeline5NosEmSequencia(t *testing.T) {
	studio := novoStudioSimulado(t, http.StatusOK)
	defer studio.Close()
	cf := novoCFSimulado(t)
	defer cf.Close()
	dest := t.TempDir()
	cfg := configDeTeste(t)
	rcp, err := ExecutePipeline(cfg, PipelineDeps{
		Launcher:     NovoLauncherClient(studio.URL, "u", "p", studio.Client()),
		Cloudflare:   NovoCloudflareClient(cf.URL, "zid", "token", studio.Client()),
		Dest:         dest,
		HealthURL:    studio.URL + "/health",
		Contagens:    contagensDeTeste(),
		TLSVerify:    func(string) int { return 0 },
		HTTPCheck:    checkHTTP(studio.Client()),
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("ExecutePipeline: %v", err)
	}
	if rcp.Status != "done" || rcp.HTTPCode != 200 || rcp.CFStatus != 1 {
		t.Fatalf("recibo inesperado: %+v", rcp)
	}
	if _, err := os.Stat(filepath.Join(dest, NomeReceipt)); err != nil {
		t.Fatalf("recibo não gravado: %v", err)
	}
}

func TestExecutePipelineSemReciboQuandoGatesFalham(t *testing.T) {
	studio := novoStudioSimulado(t, http.StatusInternalServerError)
	defer studio.Close()
	cf := novoCFSimulado(t)
	defer cf.Close()
	dest := t.TempDir()
	_, err := ExecutePipeline(configDeTeste(t), PipelineDeps{
		Launcher:     NovoLauncherClient(studio.URL, "u", "p", studio.Client()),
		Cloudflare:   NovoCloudflareClient(cf.URL, "zid", "token", studio.Client()),
		Dest:         dest,
		HealthURL:    studio.URL + "/health",
		Contagens:    contagensDeTeste(),
		TLSVerify:    func(string) int { return 0 },
		HTTPCheck:    checkHTTP(studio.Client()),
		PollInterval: time.Millisecond,
	})
	if err == nil {
		t.Fatal("gates reprovados devem falhar o pipeline")
	}
	if _, statErr := os.Stat(filepath.Join(dest, NomeReceipt)); statErr == nil {
		t.Fatal("recibo de sucesso não pode existir quando os gates falham")
	}
}

func TestExecutePipelineRejeitaConfigInvalida(t *testing.T) {
	_, err := ExecutePipeline(DeployConfig{}, PipelineDeps{Dest: t.TempDir()})
	if err == nil {
		t.Fatal("configuração vazia deve ser rejeitada")
	}
}

// contagensDeTeste devolve 0 para o CNAME antigo e 1 para o novo.
func contagensDeTeste() func(string) int {
	return func(nome string) int {
		if nome == "" || nome == "antigo" {
			return 0
		}
		return 1
	}
}

func checkHTTP(hc *http.Client) func(string) (int, error) {
	return func(url string) (int, error) {
		resp, err := hc.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		return resp.StatusCode, nil
	}
}
