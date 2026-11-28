package cli

import (
	"strings"
	"testing"
)

// As 4 etapas novas aparecem na ajuda do video-render.
func TestVideoRenderAjudaListaEtapasNovas(t *testing.T) {
	ajuda := rodarCLI(t, "video-render")
	for _, sub := range []string{"prevoo", "voz", "mix", "motor"} {
		if !strings.Contains(ajuda.stdout, "  "+sub+" ") {
			t.Errorf("ajuda sem o subcomando %q: %q", sub, ajuda.stdout)
		}
	}
}

// Sem --projeto, cada etapa nova sai com exit 2 e nomeia a flag.
func TestVideoRenderEtapasExigemProjeto(t *testing.T) {
	for _, sub := range []string{"prevoo", "voz", "mix", "motor"} {
		r := rodarCLI(t, "video-render", sub)
		if r.code != 2 || !strings.Contains(r.stderr, "--projeto") {
			t.Errorf("%s sem --projeto: exit %d, stderr %q", sub, r.code, r.stderr)
		}
	}
}

// Flag desconhecida no motor sai com exit 2.
func TestVideoRenderMotorFlagInvalida(t *testing.T) {
	if r := rodarCLI(t, "video-render", "motor", "--nao-existe"); r.code != 2 {
		t.Fatalf("exit %d, quer 2", r.code)
	}
}

// Projeto inexistente: cada etapa falha com exit 1 e a mensagem leva o nome da etapa.
func TestVideoRenderEtapasFalhamEmProjetoVazio(t *testing.T) {
	vazio := t.TempDir()
	for _, sub := range []string{"prevoo", "voz", "mix", "motor"} {
		args := []string{"video-render", sub, "--projeto", vazio}
		if sub == "motor" {
			args = append(args, "--raiz", t.TempDir())
		}
		r := rodarCLI(t, args...)
		if r.code != 1 || !strings.Contains(r.stderr, "video-render "+sub+":") {
			t.Errorf("%s em projeto vazio: exit %d, stderr %q", sub, r.code, r.stderr)
		}
	}
}
