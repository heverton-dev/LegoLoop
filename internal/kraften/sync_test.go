package kraften

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

const respostaListaCNAME = `{"result":[],"success":true}`

// novoServidorCFArmazena simula a API v4 com contagem de registros e chamadas.
func novoServidorCFArmazena(t *testing.T, registros map[string]int, chamadas *[]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*chamadas = append(*chamadas, r.Method+" "+r.URL.Path)
		if r.URL.Path == "/zones/zid/dns_records" && r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]any{"success": true})
			return
		}
		if r.Method == http.MethodGet {
			q := r.URL.Query().Get("name")
			n := 0
			if registros != nil {
				n = registros[q]
			}
			if n == 0 {
				w.Write([]byte(respostaListaCNAME))
				return
			}
			w.Write([]byte(`{"result":[{"id":"rec1"}],"success":true}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
}

func TestFqdnDeComDominioBaseExplicito(t *testing.T) {
	if got := FqdnDe("sub01", "vpsconexao.org"); got != "sub01.vpsconexao.org" {
		t.Fatalf("fqdn inesperado: %s", got)
	}
	if got := FqdnDe("sub01.outro.com", "ignorado.org"); got != "sub01.outro.com" {
		t.Fatalf("fqdn completo deve passar direto: %s", got)
	}
}

func TestSyncCloudflareCriaCNAME(t *testing.T) {
	var chamadas []string
	srv := novoServidorCFArmazena(t, nil, &chamadas)
	defer srv.Close()
	cf := NovoCloudflareClient(srv.URL, "zid", "token", srv.Client())
	if err := cf.EnsureCNAME("sub01.vpsconexao.org", "server.vpsconexao.org"); err != nil {
		t.Fatalf("EnsureCNAME: %v", err)
	}
	encontrouPOST := false
	for _, c := range chamadas {
		if c == "POST /zones/zid/dns_records" {
			encontrouPOST = true
		}
	}
	if !encontrouPOST {
		t.Fatalf("POST de criação ausente: %v", chamadas)
	}
}

func TestSyncCloudflareNaoDuplicaCNAME(t *testing.T) {
	var chamadas []string
	srv := novoServidorCFArmazena(t, map[string]int{"sub01.vpsconexao.org": 1}, &chamadas)
	defer srv.Close()
	cf := NovoCloudflareClient(srv.URL, "zid", "token", srv.Client())
	if err := cf.EnsureCNAME("sub01.vpsconexao.org", "server.vpsconexao.org"); err != nil {
		t.Fatalf("EnsureCNAME: %v", err)
	}
	for _, c := range chamadas {
		if c == "POST /zones/zid/dns_records" {
			t.Fatal("CNAME existente não pode ser recriado")
		}
	}
}

func TestSyncCloudflareRecusaHostVazio(t *testing.T) {
	srv := novoServidorCFArmazena(t, nil, &[]string{})
	defer srv.Close()
	cf := NovoCloudflareClient(srv.URL, "zid", "token", srv.Client())
	if err := cf.EnsureCNAME("", "server.vpsconexao.org"); !errors.Is(err, ErrInvalidDeployInput) {
		t.Fatalf("host vazio deve ser rejeitado, obteve %v", err)
	}
}

func TestSyncMetadataRegistraESincronizaDNS(t *testing.T) {
	var ordem []string
	mux := http.NewServeMux()
	mux.HandleFunc("/api/admin/apps", func(w http.ResponseWriter, r *http.Request) {
		ordem = append(ordem, "apps")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"id":"app01"}`))
	})
	srvCF := novoServidorCFArmazena(t, nil, &ordem)
	defer srvCF.Close()
	srv := httptest.NewServer(mux)
	defer srv.Close()
	cf := NovoCloudflareClient(srvCF.URL, "zid", "token", srv.Client())
	l := NovoLauncherClient(srv.URL, "user", "pass", srv.Client())
	cfg := DeployConfig{AppID: "app01", ZipPath: "pacote.zip", Subdomain: "sub01", TargetHost: "server.vpsconexao.org", BaseDomain: "vpsconexao.org"}
	if err := SyncMetadata(cfg, l, cf); err != nil {
		t.Fatalf("SyncMetadata: %v", err)
	}
	if len(ordem) < 2 || ordem[0] != "GET /zones/zid/dns_records" || ordem[1] != "POST /zones/zid/dns_records" {
		t.Fatalf("CNAME deve nascer antes de qualquer consulta DNS do navegador; ordem: %v", ordem)
	}
	encontrouApps := false
	for _, etapa := range ordem {
		if etapa == "apps" {
			encontrouApps = true
		}
	}
	if !encontrouApps {
		t.Fatalf("registro de metadados ausente: %v", ordem)
	}
}
