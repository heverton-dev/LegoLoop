package homologacao

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const provChave = "chave-de-teste-xyz"

func provEscrever(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.WriteFile(caminho, []byte(conteudo), 0o600); err != nil {
		t.Fatal(err)
	}
}

func provLer(t *testing.T, caminho string) string {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func provServidor(t *testing.T, versao string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer "+provChave {
			http.Error(w, "sem chave", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/api/version":
			fmt.Fprintf(w, `{"version":%q}`, versao)
		case "/v1/models":
			fmt.Fprint(w, `{"data":[{"id":"combo-a"},{"id":"combo-b"}]}`)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func provPortaFechada(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := l.Addr().String()
	l.Close()
	return addr
}

func TestProviderVerificaECorrigeComBackup(t *testing.T) {
	dir := t.TempDir()
	const velho, novo = "velho.exemplo", "novo.exemplo"
	json := filepath.Join(dir, "opencode.json")
	jsonc := filepath.Join(dir, "mimocode.jsonc")
	yml := filepath.Join(dir, "models.yml")
	sem := filepath.Join(dir, "sem.json")
	origJSON := "{\n  \"provider\": {\n    \"x\": {\n      \"options\": {\n        \"baseURL\": \"https://" + velho + "/v1\"\n      }\n    }\n  }\n}\n"
	origJSONC := "{\n  // comentário\n  \"provider\": {\n    \"baseURL\": \"https://" + velho + ":8080/v1\"\n  }\n}\n"
	origYML := "providers:\n  nine:\n    baseUrl: '" + "https://" + velho + "/v1'\nother:\n  baseUrl: https://" + velho + "/v1\n"
	origSem := "{\"outro\": {\"baseURL\": \"https://" + velho + "/v1\"}}\n"
	provEscrever(t, json, origJSON)
	provEscrever(t, jsonc, origJSONC)
	provEscrever(t, yml, origYML)
	provEscrever(t, sem, origSem)

	ends, err := LerEnderecos(yml)
	if err != nil || len(ends) != 1 || ends[0].Host != velho || ends[0].Linha != 3 {
		t.Fatalf("yml: %+v %v", ends, err)
	}
	if ends, _ = LerEnderecos(jsonc); len(ends) != 1 || ends[0].Host != velho+":8080" {
		t.Fatalf("jsonc: %+v", ends)
	}
	if ends, err = LerEnderecos(sem); err != nil || len(ends) != 0 {
		t.Fatalf("sem seção deve dar lista vazia: %+v %v", ends, err)
	}

	for _, f := range []string{json, jsonc, yml, sem} {
		m, err := Corrigir(f, velho, novo, false)
		if err != nil || m.Aplicada {
			t.Fatalf("sem confirmação não aplica: %+v %v", m, err)
		}
		if _, err := os.Stat(f + SufixoBackup); err == nil {
			t.Fatalf("sem confirmação não cria backup: %s", f)
		}
	}
	if provLer(t, json) != origJSON {
		t.Fatal("arquivo mudou sem confirmação")
	}

	m, err := Corrigir(json, velho, novo, true)
	if err != nil || !m.Aplicada || !m.Conferida {
		t.Fatalf("correção: %+v %v", m, err)
	}
	if provLer(t, json+SufixoBackup) != origJSON {
		t.Fatal("backup deve ser cópia do original")
	}
	if got, want := provLer(t, json), strings.Replace(origJSON, velho, novo, 1); got != want {
		t.Fatalf("só o host pode mudar:\n%s", got)
	}

	if _, err := Corrigir(yml, velho, novo, true); err != nil {
		t.Fatal(err)
	}
	if got, want := provLer(t, yml), strings.Replace(origYML, velho, novo, 1); got != want {
		t.Fatalf("a seção seguinte não pode mudar:\n%s", got)
	}

	if m, err = Corrigir(sem, velho, novo, true); err != nil || m.Aplicada {
		t.Fatalf("sem seção: %+v %v", m, err)
	}
	if provLer(t, sem) != origSem {
		t.Fatal("arquivo sem seção foi alterado")
	}
	if _, err := os.Stat(sem + SufixoBackup); err == nil {
		t.Fatal("arquivo sem seção não ganha backup")
	}
}

func TestProviderSondaSemVazarChave(t *testing.T) {
	srv := provServidor(t, "0.5.95")
	host := strings.TrimPrefix(srv.URL, "http://")
	s := SondarHost(context.Background(), srv.Client(), "http", host, provChave)
	if !s.Resolve || s.Versao != "0.5.95" || len(s.Combos) != 2 || s.Combos[0] != "combo-a" {
		t.Fatalf("sonda: %+v", s)
	}
	s.Comparar("0.5.116")
	if !s.Antigo || s.VersaoDisponivel != "0.5.116" {
		t.Fatalf("0.5.95 deve ser antigo frente a 0.5.116: %+v", s)
	}
	s.Comparar("0.5.95")
	if s.Antigo {
		t.Fatal("versão igual não é antiga")
	}

	morto := SondarHost(context.Background(), srv.Client(), "http", provPortaFechada(t), provChave)
	if morto.Resolve || morto.Erro == "" {
		t.Fatalf("porta fechada não resolve: %+v", morto)
	}
	if strings.Contains(fmt.Sprintf("%+v %+v", s, morto), provChave) {
		t.Fatal("a chave vazou no resultado")
	}
}

func TestCorrigirDoisHostsMantemOBackupDoOriginal(t *testing.T) {
	arq := filepath.Join(t.TempDir(), "opencode.json")
	original := "{\n  \"provider\": {\n    \"a\": {\"options\": {\"baseURL\": \"https://antigo-um.test/v1\"}},\n    \"b\": {\"options\": {\"baseURL\": \"https://antigo-dois.test/v1\"}}\n  }\n}\n"
	if err := os.WriteFile(arq, []byte(original), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, antigo := range []string{"antigo-um.test", "antigo-dois.test"} {
		if _, err := Corrigir(arq, antigo, "novo.test", true); err != nil {
			t.Fatal(err)
		}
	}
	bak, err := os.ReadFile(arq + SufixoBackup)
	if err != nil || string(bak) != original {
		t.Fatalf("o backup deve guardar o original: %q (%v)", bak, err)
	}
}
