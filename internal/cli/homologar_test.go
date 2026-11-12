package cli

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"legoloop/internal/roteamento"
	"legoloop/internal/roteamento/homologacao"
)

type harnessContador struct {
	chamadas *int
	saida    *bytes.Buffer
	antes    *string
}

func (h harnessContador) Executar(_ context.Context, _, prompt string, _ time.Duration) (homologacao.Saida, error) {
	if *h.chamadas == 0 {
		*h.antes = h.saida.String()
	}
	*h.chamadas++
	s := homologacao.Saida{CustoUSD: 0.01, Tokens: 10, PicoMB: 420}
	if prompt == homologacao.PromptB {
		s.ExitCode = 1
	}
	return s, nil
}

func ambienteHomologar(t *testing.T, livre uint64) (Deps, *bytes.Buffer, *int, *string) {
	t.Helper()
	out := &bytes.Buffer{}
	chamadas, antes := new(int), new(string)
	inv := []roteamento.ItemInv{{
		Par:     roteamento.Par{Harness: "claude", Via: roteamento.ViaOficial},
		Binario: "/bin/claude", Modelos: []string{"sonnet"}, Estado: roteamento.EstadoPendente,
	}}
	d := Deps{
		Stdin: strings.NewReader(""), Stdout: out, Stderr: &bytes.Buffer{}, Home: t.TempDir(),
		Rota: &RotaDeps{Inventario: func(context.Context) ([]roteamento.ItemInv, error) { return inv, nil }},
		Homologar: &HomologarDeps{
			Harness: func(_, _, _ string) homologacao.Harness {
				return harnessContador{chamadas: chamadas, saida: out, antes: antes}
			},
			MemoriaLivre: func() (uint64, error) { return livre, nil },
		},
	}
	return d, out, chamadas, antes
}

func TestHomologarConfirmacaoEPendente(t *testing.T) {
	const chave = "claude|oficial"
	for _, args := range [][]string{{}, {"--non-interactive"}} {
		d, out, chamadas, _ := ambienteHomologar(t, 4<<30)
		if code := Run(context.Background(), append([]string{"homologar"}, args...), d); code != 0 {
			t.Fatalf("exit %d: %s", code, out)
		}
		if *chamadas != 0 {
			t.Fatalf("sem confirmação não pode chamar: %d", *chamadas)
		}
		if !strings.Contains(out.String(), "Teto total estimado: 3 chamadas") {
			t.Fatalf("falta o teto total: %s", out)
		}
		st, err := homologacao.Ler(homologacao.Caminho(d.Home))
		if err != nil || st[chave].Estado != homologacao.EstadoPendente {
			t.Fatalf("esperava pendente: %v %v", st, err)
		}
	}

	d, out, chamadas, antes := ambienteHomologar(t, 4<<30)
	if code := Run(context.Background(), []string{"homologar", "--yes", "--teto", "1s"}, d); code != 0 {
		t.Fatalf("exit %d: %s", code, out)
	}
	if *chamadas != 3 || !strings.Contains(*antes, "Teto total estimado") {
		t.Fatalf("chamadas=%d, teto antes da primeira chamada: %q", *chamadas, *antes)
	}
	st, err := homologacao.Ler(homologacao.Caminho(d.Home))
	if err != nil || st[chave].Estado != homologacao.EstadoHomologado || st[chave].PesoMB != 420 {
		t.Fatalf("esperava homologado com peso: %v %v", st, err)
	}
	lido := rotaLerHomologacao(filepath.Join(d.Home, ".aidd", "homologacao.json"))
	if lido[chave].Estado != roteamento.EstadoHomologado || lido[chave].PesoMB != 420 {
		t.Fatalf("o inventário não lê o resultado: %v", lido)
	}

	d, out, chamadas, _ = ambienteHomologar(t, 100<<20)
	if code := Run(context.Background(), []string{"homologar", "--yes"}, d); code != 1 {
		t.Fatalf("memória baixa deve dar exit 1, veio %d", code)
	}
	st, _ = homologacao.Ler(homologacao.Caminho(d.Home))
	if *chamadas != 0 || len(st) != 0 {
		t.Fatalf("memória baixa não pode chamar nem gravar: %d %v", *chamadas, st)
	}
}

const chaveProvTeste = "chave-de-teste-prov"

func servidorProvTeste(t *testing.T, versao string, viuChave *atomic.Bool) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if viuChave != nil && r.Header.Get("Authorization") == "Bearer "+chaveProvTeste {
			viuChave.Store(true)
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

func escreverProvTeste(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.WriteFile(caminho, []byte(conteudo), 0o600); err != nil {
		t.Fatal(err)
	}
}

func lerProvTeste(t *testing.T, caminho string) string {
	t.Helper()
	b, err := os.ReadFile(caminho)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func existeProvTeste(caminho string) bool {
	_, err := os.Stat(caminho)
	return err == nil
}

func TestSondaECorrecaoDeProvider(t *testing.T) {
	t.Setenv("NINEROUTER_API_KEY", "")
	t.Setenv("NINEROUTER_KEY", chaveProvTeste)
	var viu atomic.Bool
	correto := servidorProvTeste(t, "0.5.116", &viu)
	velho := servidorProvTeste(t, "0.5.95", nil)
	t.Setenv("NINEROUTER_URL", correto.URL)
	hostCorreto := strings.TrimPrefix(correto.URL, "http://")
	hostVelho := strings.TrimPrefix(velho.URL, "http://")
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	hostMorto := l.Addr().String()
	l.Close()

	dir := t.TempDir()
	oc, mimo, yml, sem := filepath.Join(dir, "opencode.json"), filepath.Join(dir, "mimocode.jsonc"), filepath.Join(dir, "models.yml"), filepath.Join(dir, "sem.json")
	origOC := fmt.Sprintf("{\n  \"provider\": {\n    \"x\": {\n      \"options\": {\n        \"baseURL\": \"http://%s/v1\"\n      }\n    }\n  }\n}\n", hostMorto)
	origMimo := fmt.Sprintf("{\n  // comentário\n  \"provider\": {\n    \"baseURL\": \"http://%s/v1\"\n  }\n}\n", hostVelho)
	origYML := fmt.Sprintf("providers:\n  nine:\n    baseUrl: 'http://%s/v1'\nother:\n  baseUrl: http://%s/v1\n", hostCorreto, hostMorto)
	origSem := fmt.Sprintf("{\"outro\": {\"baseURL\": \"http://%s/v1\"}}\n", hostMorto)
	escreverProvTeste(t, oc, origOC)
	escreverProvTeste(t, mimo, origMimo)
	escreverProvTeste(t, yml, origYML)
	escreverProvTeste(t, sem, origSem)

	out := &bytes.Buffer{}
	novoDeps := func() Deps {
		return Deps{
			Stdin: strings.NewReader(""), Stdout: out, Stderr: &bytes.Buffer{}, Home: t.TempDir(),
			Provider: &ProviderDeps{HTTP: correto.Client(), Arquivos: []string{oc, mimo, yml, sem}, Esquema: "http", Canonico: hostCorreto},
		}
	}
	rodar := func(args ...string) {
		t.Helper()
		if code := Run(context.Background(), append([]string{"homologar"}, args...), novoDeps()); code != 0 {
			t.Fatalf("exit %d: %s", code, out)
		}
	}

	rodar("--provider")
	sonda := out.String()
	for _, want := range []string{hostMorto, "false", hostVelho, "0.5.95", "0.5.116", "combo-a,combo-b"} {
		if !strings.Contains(sonda, want) {
			t.Fatalf("a sonda não traz %q:\n%s", want, sonda)
		}
	}
	if !viu.Load() {
		t.Fatal("o host canônico confiável deve receber a chave")
	}

	rodar("--corrigir-provider", "--non-interactive")
	for f, orig := range map[string]string{oc: origOC, mimo: origMimo, yml: origYML, sem: origSem} {
		if lerProvTeste(t, f) != orig || existeProvTeste(f+homologacao.SufixoBackup) {
			t.Fatalf("sem confirmação nada muda: %s", f)
		}
	}

	rodar("--corrigir-provider", "--yes")
	if got, want := lerProvTeste(t, oc), strings.Replace(origOC, hostMorto, hostCorreto, 1); got != want {
		t.Fatalf("opencode.json: %s", got)
	}
	if got, want := lerProvTeste(t, mimo), strings.Replace(origMimo, hostVelho, hostCorreto, 1); got != want {
		t.Fatalf("mimocode.jsonc: %s", got)
	}
	if lerProvTeste(t, oc+homologacao.SufixoBackup) != origOC || lerProvTeste(t, mimo+homologacao.SufixoBackup) != origMimo {
		t.Fatal("o backup deve ser cópia do original")
	}
	for f, orig := range map[string]string{yml: origYML, sem: origSem} {
		if lerProvTeste(t, f) != orig || existeProvTeste(f+homologacao.SufixoBackup) {
			t.Fatalf("host correto e arquivo sem seção ficam intactos: %s", f)
		}
	}
	if strings.Contains(out.String(), chaveProvTeste) {
		t.Fatal("a saída não pode conter a chave")
	}
}
