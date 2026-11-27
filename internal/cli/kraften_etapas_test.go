package cli

import (
	"archive/zip"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// zipSaudavelKraften grava um .zip com os 3 artefatos universais válidos.
func zipSaudavelKraften(t *testing.T) string {
	t.Helper()
	caminho := filepath.Join(t.TempDir(), "pacote.zip")
	arq, err := os.Create(caminho)
	if err != nil {
		t.Fatalf("criar zip: %v", err)
	}
	defer arq.Close()
	w := zip.NewWriter(arq)
	entradas := map[string]string{
		"app/Dockerfile":       "FROM alpine:3.20\nHEALTHCHECK CMD wget -q -O- http://localhost:8080/health || exit 1\n",
		"app/docker-stack.yml": "services:\n  web:\n    networks:\n      - network_conexao\n",
		"app/.env.example":     "ADMIN_USER=admin\nADMIN_PASS=segredo\n",
	}
	for nome, corpo := range entradas {
		f, err := w.Create(nome)
		if err != nil {
			t.Fatalf("entrada %s: %v", nome, err)
		}
		if _, err := f.Write([]byte(corpo)); err != nil {
			t.Fatalf("gravar %s: %v", nome, err)
		}
	}
	if err := w.Close(); err != nil {
		t.Fatalf("fechar zip: %v", err)
	}
	return caminho
}

func depsKraften() Deps {
	return Deps{Stdin: os.Stdin, Stdout: os.Stdout, Stderr: os.Stderr}
}

func TestKraftenEtapaPreflight(t *testing.T) {
	dest := t.TempDir()
	code := cmdKraftenDeploy(context.Background(), []string{
		"preflight", "--zip", zipSaudavelKraften(t), "--dest", dest,
	}, depsKraften())
	if code != 0 {
		t.Fatalf("preflight válido rejeitado, exit %d", code)
	}
	if _, err := os.Stat(filepath.Join(dest, "preflight.json")); err != nil {
		t.Fatalf("artefato preflight.json ausente: %v", err)
	}
}

func TestKraftenEtapaPreflightRejeitaZipRuim(t *testing.T) {
	dest := t.TempDir()
	code := cmdKraftenDeploy(context.Background(), []string{
		"preflight", "--zip", filepath.Join(t.TempDir(), "inexistente.zip"), "--dest", dest,
	}, depsKraften())
	if code == 0 {
		t.Fatal("zip inválido deve falhar com exit != 0")
	}
}

func TestKraftenEtapaSync(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"id":"app01"}`))
	})
	srvCF := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"success":true}`))
			return
		}
		w.Write([]byte(`{"result":[],"success":true}`))
	}))
	defer srvCF.Close()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	dest := t.TempDir()
	code := cmdKraftenDeploy(context.Background(), []string{
		"sync", "--app", "app01", "--subdominio", "sub01", "--alvo", "server.vpsconexao.org",
		"--dominio-base", "vpsconexao.org",
		"--launcher", srv.URL, "--user", "u", "--pass", "p",
		"--cloudflare", srvCF.URL, "--zone", "zid", "--token", "tok", "--dest", dest,
	}, depsKraften())
	if code != 0 {
		t.Fatalf("sync rejeitado, exit %d", code)
	}
	corpo, err := os.ReadFile(filepath.Join(dest, "sync.json"))
	if err != nil {
		t.Fatalf("artefato sync.json ausente: %v", err)
	}
	var carga map[string]any
	if err := json.Unmarshal(corpo, &carga); err != nil {
		t.Fatalf("sync.json inválido: %v", err)
	}
	if carga["fqdn"] != "sub01.vpsconexao.org" {
		t.Fatalf("fqdn inesperado: %v", carga["fqdn"])
	}
}

func TestKraftenEtapaDeploy(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps/app01/update-zip", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"job":"job-1"}`))
	})
	mux.HandleFunc("/api/admin/jobs/job-1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"done"}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()
	dest := t.TempDir()
	code := cmdKraftenDeploy(context.Background(), []string{
		"deploy", "--app", "app01", "--zip", zipSaudavelKraften(t),
		"--launcher", srv.URL, "--user", "u", "--pass", "p", "--dest", dest,
	}, depsKraften())
	if code != 0 {
		t.Fatalf("deploy rejeitado, exit %d", code)
	}
	corpo, err := os.ReadFile(filepath.Join(dest, "deploy.json"))
	if err != nil {
		t.Fatalf("artefato deploy.json ausente: %v", err)
	}
	var carga map[string]any
	if err := json.Unmarshal(corpo, &carga); err != nil || carga["status"] != "done" {
		t.Fatalf("deploy.json inesperado: %s (%v)", corpo, err)
	}
}

func TestKraftenEtapaGatesReprovaSemTLS(t *testing.T) {
	original := fronteiraTLSKraften
	originalRede := redeAnexadaKraften
	fronteiraTLSKraften = func(string) int { return 1 }
	redeAnexadaKraften = func(string, string) int { return 1 }
	defer func() { fronteiraTLSKraften = original; redeAnexadaKraften = originalRede }()
	dest := t.TempDir()
	code := cmdKraftenDeploy(context.Background(), []string{
		"gates", "--health", "http://localhost/health", "--host-novo", "sub01.vpsconexao.org",
		"--rede", "rede_x", "--servico", "svc_x",
		"--cloudflare", "http://localhost", "--zone", "z", "--token", "t", "--dest", dest,
	}, depsKraften())
	if code == 0 {
		t.Fatal("gate TLS=1 deve reprovar com exit != 0")
	}
	if _, err := os.Stat(filepath.Join(dest, "gates.json")); err != nil {
		t.Fatalf("gates.json deve existir mesmo reprovado: %v", err)
	}
}

func TestKraftenEtapaReceipt(t *testing.T) {
	dest := t.TempDir()
	escrever := func(nome, corpo string) {
		if err := os.WriteFile(filepath.Join(dest, nome), []byte(corpo), 0o644); err != nil {
			t.Fatalf("gravar %s: %v", nome, err)
		}
	}
	escrever("deploy.json", `{"status":"done"}`)
	escrever("gates.json", `{"cf_old":0,"cf_new":1,"tls_verdict":0,"http_code":200}`)
	code := cmdKraftenDeploy(context.Background(), []string{
		"receipt", "--app", "app01", "--dest", dest,
	}, depsKraften())
	if code != 0 {
		t.Fatalf("receipt rejeitado, exit %d", code)
	}
	corpo, err := os.ReadFile(filepath.Join(dest, "DEPLOY_RECEIPT.json"))
	if err != nil {
		t.Fatalf("recibo ausente: %v", err)
	}
	var r map[string]any
	if err := json.Unmarshal(corpo, &r); err != nil {
		t.Fatalf("recibo inválido: %v", err)
	}
	if r["status"] != "done" || r["http_code"] != float64(200) || r["cf_status"] != float64(1) {
		t.Fatalf("recibo inesperado: %v", r)
	}
}
