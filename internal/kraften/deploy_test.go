package kraften

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDeployStreamAcompanhaJobAteDone(t *testing.T) {
	consultas := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps/app01/update-zip", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"job":"job-1"}`))
	})
	mux.HandleFunc("/api/admin/jobs/job-1", func(w http.ResponseWriter, r *http.Request) {
		consultas++
		if consultas < 3 {
			w.Write([]byte(`{"status":"running"}`))
			return
		}
		w.Write([]byte(`{"status":"done"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	l := NovoLauncherClient(srv.URL, "user", "pass", srv.Client())
	l.PollInterval = time.Millisecond
	cfg := DeployConfig{AppID: "app01", ZipPath: criarZip(t, entradasSaudaveis()), Subdomain: "sub01", TargetHost: "server.vpsconexao.org"}
	status, err := DeployStream(cfg, l, "")
	if err != nil {
		t.Fatalf("DeployStream: %v", err)
	}
	if status != "done" {
		t.Fatalf("status final esperado done, obteve %s", status)
	}
	if consultas != 3 {
		t.Fatalf("esperava 3 consultas ao job, obteve %d", consultas)
	}
}

func TestDeployStreamFalhaQuandoJobFalha(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps/app01/update-zip", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"job":"job-2"}`))
	})
	mux.HandleFunc("/api/admin/jobs/job-2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"failed"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	l := NovoLauncherClient(srv.URL, "user", "pass", srv.Client())
	l.PollInterval = time.Millisecond
	cfg := DeployConfig{AppID: "app01", ZipPath: criarZip(t, entradasSaudaveis()), Subdomain: "sub01", TargetHost: "server.vpsconexao.org"}
	if _, err := DeployStream(cfg, l, ""); err == nil {
		t.Fatal("job failed deve interromper o Processo")
	}
}

func TestDeployStreamRejeitaConfigInvalida(t *testing.T) {
	l := NovoLauncherClient("http://localhost", "u", "p", nil)
	if _, err := DeployStream(DeployConfig{}, l, ""); err == nil {
		t.Fatal("configuração vazia deve ser rejeitada")
	}
}

func TestDeployStreamEnviaZipParaInstallQuandoNovo(t *testing.T) {
	var caminho string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps/app01/install-zip", func(w http.ResponseWriter, r *http.Request) {
		caminho = r.URL.Path
		w.Write([]byte(`{"job":"job-3"}`))
	})
	mux.HandleFunc("/api/admin/jobs/job-3", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"done"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	l := NovoLauncherClient(srv.URL, "user", "pass", srv.Client())
	l.PollInterval = time.Millisecond
	cfg := DeployConfig{AppID: "app01", ZipPath: criarZip(t, entradasSaudaveis()), Subdomain: "sub01", TargetHost: "server.vpsconexao.org"}
	if _, err := DeployStream(cfg, l, "install-zip"); err != nil {
		t.Fatalf("DeployStream install: %v", err)
	}
	if !strings.HasSuffix(caminho, "/install-zip") {
		t.Fatalf("esperava install-zip, obteve %s", caminho)
	}
}
