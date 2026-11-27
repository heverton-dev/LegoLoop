package kraften

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// novoServidorGates responde /health com o código informado e simula TLS.
func novoServidorGates(t *testing.T, codigo int) *httptest.Server {
	t.Helper()
	return httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(codigo)
	}))
}

func TestVerifyGatesPassaComProvasBinarias(t *testing.T) {
	gates := Gates{
		CFOld:      0,
		CFNew:      1,
		TLSVerdict: 0,
		HTTPCode:   http.StatusOK,
		RedeCount:  1,
	}
	if err := VerifyGates(gates); err != nil {
		t.Fatalf("gates válidos reprovados: %v", err)
	}
}

func TestVerifyGatesReprovaHTTP500(t *testing.T) {
	gates := Gates{CFOld: 0, CFNew: 1, TLSVerdict: 0, HTTPCode: 500}
	if err := VerifyGates(gates); err == nil {
		t.Fatal("HTTP 500 deve reprovar os gates")
	}
}

func TestVerifyGatesReprovaCNAMEAntigoVivo(t *testing.T) {
	gates := Gates{CFOld: 1, CFNew: 1, TLSVerdict: 0, HTTPCode: 200}
	if err := VerifyGates(gates); err == nil {
		t.Fatal("CF_OLD=1 deve reprovar os gates")
	}
}

func TestVerifyGatesReprovaCNAMENovoAusente(t *testing.T) {
	gates := Gates{CFOld: 0, CFNew: 0, TLSVerdict: 0, HTTPCode: 200}
	if err := VerifyGates(gates); err == nil {
		t.Fatal("CF_NEW=0 deve reprovar os gates")
	}
}

func TestVerifyGatesReprovaTLS(t *testing.T) {
	gates := Gates{CFOld: 0, CFNew: 1, TLSVerdict: 1, HTTPCode: 200}
	if err := VerifyGates(gates); err == nil {
		t.Fatal("TLS=1 deve reprovar os gates")
	}
}

func TestVerifyGatesReprovaRedeSemAnexo(t *testing.T) {
	gates := Gates{CFOld: 0, CFNew: 1, TLSVerdict: 0, HTTPCode: 200, RedeCount: 0}
	if err := VerifyGates(gates); err == nil {
		t.Fatal("REDE=0 deve reprovar os gates")
	}
}

func TestVerifyGatesPassaComRedeAnexada(t *testing.T) {
	gates := Gates{CFOld: 0, CFNew: 1, TLSVerdict: 0, HTTPCode: 200, RedeCount: 1}
	if err := VerifyGates(gates); err != nil {
		t.Fatalf("REDE=1 deve aprovar os gates: %v", err)
	}
}

func TestColetarGatesLeContagensTLSHTTP(t *testing.T) {
	srv := novoServidorGates(t, http.StatusOK)
	defer srv.Close()
	contagens := func(nome string) int {
		if nome == "novo" {
			return 1
		}
		return 0
	}
	tlsVerify := func(string) int { return 0 }
	httpCheck := func(url string) (int, error) {
		resp, err := srv.Client().Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		return resp.StatusCode, nil
	}
	gates, err := ColetarGates(contagens, tlsVerify, httpCheck, func() int { return 1 }, srv.URL+"/health", "antigo", "novo")
	if err != nil {
		t.Fatalf("ColetarGates: %v", err)
	}
	if gates.CFOld != 0 || gates.CFNew != 1 || gates.TLSVerdict != 0 || gates.HTTPCode != 200 {
		t.Fatalf("gates coletados incorretos: %+v", gates)
	}
}

func TestColetarGatesPropagaFalhaDeHTTP(t *testing.T) {
	contagens := func(string) int { return 0 }
	tlsVerify := func(string) int { return 0 }
	httpCheck := func(string) (int, error) { return 0, errors.New("sem rota") }
	if _, err := ColetarGates(contagens, tlsVerify, httpCheck, func() int { return 1 }, "https://x/health", "antigo", "novo"); err == nil {
		t.Fatal("falha de HTTP deve propagar erro")
	}
}
