package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/assets"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

// geradaComDuasPecas gera a matriz de alfa e beta (sem declaração) com um cérebro que as infere.
func geradaComDuasPecas(t *testing.T) (*ambiente, *MatrizDeps, *brainFake) {
	t.Helper()
	a := novoAmbiente(t)
	b := &brainFake{resp: respMatriz}
	md := itemsFixos(a.item("alfa", ""), a.item("beta", ""))
	md.Brain, md.BrainName = b, "fake/x"
	if code, out, errb := a.rodar(Deps{Matriz: md}, "matriz", "gerar", "--vault", a.vault); code != 0 {
		t.Fatalf("exit %d: %s %s", code, out, errb)
	}
	return a, md, b
}

func origemDe(a *ambiente, nome string) string {
	for _, e := range a.registro().Entradas {
		if e.Nome == nome && e.Perfil != nil {
			return e.Perfil.Origem
		}
	}
	return ""
}

// AC-E01.9: confirmar e definir, e a confirmação sobrevive ao --forcar sem chamar o cérebro para a peça.
func TestMatrizPerfilCLI_ConfirmarEDefinirSobrevivemAoForcar(t *testing.T) {
	a, md, b := geradaComDuasPecas(t)
	if origemDe(a, "alfa") != "inferido" {
		t.Fatalf("pré-condição: alfa inferido, é %q", origemDe(a, "alfa"))
	}
	code, out, errb := a.rodar(Deps{Matriz: md}, "matriz", "perfil", "confirmar", "alfa", "--vault", a.vault)
	if code != 0 || !strings.Contains(out, "Perfil de alfa confirmado") {
		t.Fatalf("confirmar: exit %d: %s %s", code, out, errb)
	}
	code, out, errb = a.rodar(Deps{Matriz: md}, "matriz", "perfil", "definir", "beta", "--etapa", "construir", "--natureza", "software", "--caso", "tdd", "--vault", a.vault)
	if code != 0 || !strings.Contains(out, "Perfil de beta definido") {
		t.Fatalf("definir: exit %d: %s %s", code, out, errb)
	}
	if origemDe(a, "alfa") != "confirmado" || origemDe(a, "beta") != "confirmado" {
		t.Fatalf("origens: alfa=%q beta=%q", origemDe(a, "alfa"), origemDe(a, "beta"))
	}
	antes := len(b.prompts)
	if code, out, errb = a.rodar(Deps{Matriz: md}, "matriz", "gerar", "--forcar", "--vault", a.vault); code != 0 {
		t.Fatalf("forcar: exit %d: %s %s", code, out, errb)
	}
	if len(b.prompts) != antes {
		t.Errorf("o cérebro foi chamado %d vez(es) para peças confirmadas", len(b.prompts)-antes)
	}
	if origemDe(a, "alfa") != "confirmado" || origemDe(a, "beta") != "confirmado" {
		t.Errorf("depois do --forcar: alfa=%q beta=%q", origemDe(a, "alfa"), origemDe(a, "beta"))
	}
	if _, err := os.Stat(filepath.Join(a.vault, "perfil.confirmado.json")); err != nil {
		t.Errorf("a confirmação fica em arquivo próprio do cofre: %v", err)
	}
}

func TestMatrizPerfilCLI_FalhasNaoMudamNada(t *testing.T) {
	a, md, _ := geradaComDuasPecas(t)
	antes, _ := os.ReadFile(filepath.Join(a.vault, "registry_matrix.json"))
	casos := []struct {
		args []string
		code int
		msg  string
	}{
		{[]string{"perfil", "definir", "alfa", "--etapa", "inventada"}, 2, "inventada"},
		{[]string{"perfil", "definir", "alfa"}, 2, "ao menos um termo"},
		{[]string{"perfil", "confirmar", "nao-existe"}, 1, "nao-existe"},
	}
	for _, c := range casos {
		args := append(append([]string{"matriz"}, c.args...), "--vault", a.vault)
		code, _, errb := a.rodar(Deps{Matriz: md}, args...)
		if code != c.code || !strings.Contains(errb, c.msg) {
			t.Errorf("%v: exit %d (queria %d), stderr %q", c.args, code, c.code, errb)
		}
	}
	if depois, _ := os.ReadFile(filepath.Join(a.vault, "registry_matrix.json")); string(depois) != string(antes) {
		t.Error("a matriz mudou apesar dos erros")
	}
	if _, err := os.Stat(filepath.Join(a.vault, "perfil.confirmado.json")); err == nil {
		t.Error("nenhuma confirmação deveria ter sido gravada")
	}
}

func TestMatrizPerfilCLI_SemMatrizDizParaGerar(t *testing.T) {
	a := novoAmbiente(t)
	code, _, errb := a.rodar(Deps{}, "matriz", "perfil", "confirmar", "alfa", "--vault", a.vault)
	if code != 1 || !strings.Contains(errb, "legoloop matriz gerar") {
		t.Errorf("exit %d: %s", code, errb)
	}
}

func TestMatrizPerfilCLI_UsoIncorreto(t *testing.T) {
	a := novoAmbiente(t)
	for _, args := range [][]string{{"matriz", "perfil"}, {"matriz", "perfil", "inexistente"}, {"matriz", "perfil", "confirmar"}, {"matriz", "perfil", "definir"}} {
		if code, _, errb := a.rodar(Deps{}, args...); code != 2 || !strings.Contains(errb, "Uso: legoloop matriz") {
			t.Errorf("%v: exit %d: %q", args, code, errb)
		}
	}
}

func TestMatrizPerfilCLI_GerarAvisaQuandoAConfirmacaoPerdeAValidade(t *testing.T) {
	a := novoAmbiente(t)
	it := a.item("alfa", "")
	md := itemsFixos(it)
	md.Brain, md.BrainName = &brainFake{resp: respMatriz}, "fake/x"
	a.rodar(Deps{Matriz: md}, "matriz", "gerar", "--vault", a.vault)
	a.rodar(Deps{Matriz: md}, "matriz", "perfil", "confirmar", "alfa", "--vault", a.vault)
	it.Meta = map[string]string{"sha256": "h-novo"}
	md2 := itemsFixos(it)
	md2.Brain, md2.BrainName = md.Brain, "fake/x"
	code, _, errb := a.rodar(Deps{Matriz: md2}, "matriz", "gerar", "--vault", a.vault)
	if code != 0 || !strings.Contains(errb, "aviso: alfa: a confirmação do perfil perdeu a validade") {
		t.Errorf("exit %d, stderr:\n%s", code, errb)
	}
	if origemDe(a, "alfa") != "inferido" {
		t.Errorf("o perfil volta a inferido, é %q", origemDe(a, "alfa"))
	}
}

func TestMatrizPerfilCLI_MensagensSeguemOPerfilSTEPT(t *testing.T) {
	if vs := ste.Default.Violations(); len(vs) != 0 {
		t.Fatalf("catálogo com violações: %v", vs)
	}
	for _, id := range []string{"matriz.perfil_confirmado", "matriz.perfil_definido"} {
		if m := ste.Msg(id); strings.HasPrefix(m, "[mensagem desconhecida") {
			t.Errorf("mensagem %s não está no catálogo", id)
		}
	}
	var _ = matriz.OrigemConfirmado
	var _ assets.Item
}

func TestMatrizPerfilCLI_DefinirGravaCasoEtapaENatureza(t *testing.T) {
	a, md, _ := geradaComDuasPecas(t)
	if code, _, errb := a.rodar(Deps{Matriz: md}, "matriz", "perfil", "definir", "beta", "--etapa", "construir", "--natureza", "software", "--caso", "tdd", "--vault", a.vault); code != 0 {
		t.Fatalf("exit %d: %s", code, errb)
	}
	for _, e := range a.registro().Entradas {
		if e.Nome != "beta" {
			continue
		}
		p := e.Perfil
		if len(p.CasosDeUso) != 1 || p.CasosDeUso[0] != "tdd" || len(p.Etapas) != 1 || p.Etapas[0] != "construir" || len(p.Naturezas) != 1 || p.Naturezas[0] != "software" {
			t.Errorf("perfil gravado: %+v", p)
		}
	}
}

func TestMatrizPerfilCLI_SubcomandoDesconhecidoComPecaEUso(t *testing.T) {
	a, md, _ := geradaComDuasPecas(t)
	code, _, errb := a.rodar(Deps{Matriz: md}, "matriz", "perfil", "inexistente", "alfa", "--vault", a.vault)
	if code != 2 || !strings.Contains(errb, "Uso: legoloop matriz") {
		t.Errorf("exit %d: %s", code, errb)
	}
}
