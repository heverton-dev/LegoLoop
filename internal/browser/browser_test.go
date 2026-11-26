package browser

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolverCDPURLOrdem9222e9333(t *testing.T) {
	srv9333 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/json/version" {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv9333.Close()

	srvOff := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srvOff.Close()

	got, err := ResolverCDPURL(context.Background(), srv9333.Client(), []string{srvOff.URL, srv9333.URL})
	if err != nil {
		t.Fatalf("esperado sucesso na segunda URL, obtido: %v", err)
	}
	if got != srv9333.URL {
		t.Fatalf("esperado %q, obtido %q", srv9333.URL, got)
	}

	if _, err := ResolverCDPURL(context.Background(), srvOff.Client(), []string{srvOff.URL}); !errors.Is(err, ErrNavegadorFechado) {
		t.Fatalf("esperado ErrNavegadorFechado, obtido: %v", err)
	}
}

func TestPerfilPadraoReaproveitaLegadoOuCriaLegoLoop(t *testing.T) {
	home1 := t.TempDir()
	legado := filepath.Join(home1, ".chrome-automation-profile")
	if err := os.MkdirAll(legado, 0o755); err != nil {
		t.Fatal(err)
	}
	if got := PerfilPadrao(home1); got != legado {
		t.Fatalf("deveria reaproveitar perfil legado %q, obtido %q", legado, got)
	}

	home2 := t.TempDir()
	esperado := filepath.Join(home2, ".legoloop", "chrome-profile")
	if got := PerfilPadrao(home2); got != esperado {
		t.Fatalf("esperado perfil padrao %q, obtido %q", esperado, got)
	}
}

func TestConstruirArgsEAbrirPersistente(t *testing.T) {
	args := ConstruirArgs("/tmp/perfil", 9222, "https://flow.google.com")
	joined := strings.Join(args, " ")
	for _, want := range []string{
		"--user-data-dir=/tmp/perfil",
		"--remote-debugging-port=9222",
		"--remote-allow-origins=*",
		"https://flow.google.com",
	} {
		if !strings.Contains(joined, want) {
			t.Fatalf("arg esperado %q ausente em %q", want, joined)
		}
	}

	var exeCap string
	var argsCap []string
	home := t.TempDir()
	err := AbrirPersistente(context.Background(), OpcoesAbertura{
		Home:       home,
		Porta:      9222,
		URLInicial: "https://flow.google.com",
		DetectarExe: func() (string, error) {
			return "chrome.exe", nil
		},
		IniciarProcesso: func(_ context.Context, exe string, a []string) error {
			exeCap = exe
			argsCap = a
			return nil
		},
	})
	if err != nil {
		t.Fatalf("AbrirPersistente falhou: %v", err)
	}
	if exeCap != "chrome.exe" || len(argsCap) < 3 {
		t.Fatalf("lancamento inesperado: exe=%q args=%v", exeCap, argsCap)
	}
}

func TestErrosBrowser(t *testing.T) {
	if _, err := ResolverCDPURL(context.Background(), nil, []string{"http://[::1]:namedport"}); !errors.Is(err, ErrNavegadorFechado) {
		t.Fatalf("esperado ErrNavegadorFechado para URL invalida, obtido: %v", err)
	}
	if got := PortasPadrao(); len(got) != 2 || got[0] != "http://127.0.0.1:9222" || got[1] != "http://127.0.0.1:9333" {
		t.Fatalf("PortasPadrao inesperadas: %v", got)
	}
	errSemExe := AbrirPersistente(context.Background(), OpcoesAbertura{
		Home:        t.TempDir(),
		DetectarExe: func() (string, error) { return "", ErrNavegadorNaoEncontrado },
	})
	if !errors.Is(errSemExe, ErrNavegadorNaoEncontrado) {
		t.Fatalf("esperado ErrNavegadorNaoEncontrado, obtido: %v", errSemExe)
	}
	if _, err := DetectarNavegador(func(string) (string, error) { return "", errors.New("nao") }, func(string) bool { return false }); !errors.Is(err, ErrNavegadorNaoEncontrado) {
		t.Fatalf("esperado ErrNavegadorNaoEncontrado quando ausente, obtido: %v", err)
	}
	if got, err := DetectarNavegador(func(n string) (string, error) {
		if n == "chrome" {
			return "/usr/bin/chrome", nil
		}
		return "", errors.New("nao")
	}, func(string) bool { return false }); err != nil || got != "/usr/bin/chrome" {
		t.Fatalf("esperado /usr/bin/chrome, obtido %q err=%v", got, err)
	}
	if got, err := DetectarNavegador(func(string) (string, error) { return "", errors.New("nao") }, func(p string) bool {
		return strings.Contains(p, "chrome.exe")
	}); err != nil || !strings.Contains(got, "chrome.exe") {
		t.Fatalf("esperado caminho fixo chrome.exe, obtido %q err=%v", got, err)
	}
	if PerfilPadrao("") == "" {
		t.Fatalf("PerfilPadrao('') nao deve ser vazio")
	}
	if a := ConstruirArgs("/p", 0, ""); len(a) != 3 || a[1] != "--remote-debugging-port=9222" {
		t.Fatalf("ConstruirArgs com porta 0 inesperado: %v", a)
	}
	if err := AbrirPersistente(context.Background(), OpcoesAbertura{
		Home:            t.TempDir(),
		DetectarExe:     func() (string, error) { return "nao_existe_bin_xyz", nil },
		IniciarProcesso: nil,
	}); err == nil {
		t.Fatalf("esperado erro ao iniciar binario inexistente")
	}
	if _, err := DetectarNavegador(nil, nil); err != nil && !errors.Is(err, ErrNavegadorNaoEncontrado) {
		t.Fatalf("DetectarNavegador(nil, nil) erro inesperado: %v", err)
	}
	srv500 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv500.Close()
	if err := sondarCDP(context.Background(), srv500.Client(), srv500.URL); err == nil {
		t.Fatalf("sondarCDP deveria falhar com status 500")
	}
	if _, err := ResolverCDPURL(context.Background(), nil, nil); err != nil && !errors.Is(err, ErrNavegadorFechado) {
		t.Fatalf("ResolverCDPURL(nil, nil) erro inesperado: %v", err)
	}
	if err := AbrirPersistente(context.Background(), OpcoesAbertura{
		Home: "\x00invalido",
		DetectarExe: func() (string, error) {
			return "chrome.exe", nil
		},
	}); err == nil {
		t.Fatalf("AbrirPersistente deveria falhar com home invalida")
	}
	if got := caminhosFixosNavegador(); len(got) == 0 {
		t.Fatalf("caminhosFixosNavegador nao deve ser vazio")
	}
	srvReadErr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srvReadErr.Close()
	if err := sondarCDP(context.Background(), srvReadErr.Client(), srvReadErr.URL); err != nil {
		t.Fatalf("sondarCDP deveria passar em 200 OK: %v", err)
	}
	homeEnv := t.TempDir()
	t.Setenv("USERPROFILE", homeEnv)
	t.Setenv("HOME", homeEnv)
	if got := PerfilPadrao(""); !strings.HasPrefix(got, homeEnv) {
		t.Fatalf("PerfilPadrao('') deveria usar home %q, obtido %q", homeEnv, got)
	}
	_ = AbrirPersistente(context.Background(), OpcoesAbertura{
		Home:            t.TempDir(),
		IniciarProcesso: func(context.Context, string, []string) error { return nil },
	})
}
