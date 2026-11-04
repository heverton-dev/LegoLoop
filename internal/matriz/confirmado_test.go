package matriz_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/matriz"
)

func porNome(r *matriz.Registro, nome string) matriz.Entrada {
	for _, e := range r.Entradas {
		if e.Nome == nome {
			return e
		}
	}
	return matriz.Entrada{}
}

func chamadasSobre(b *spyBrain, nome string) int {
	n := 0
	for _, p := range b.prompts {
		if strings.Contains(p, nome) {
			n++
		}
	}
	return n
}

func hashDaArvore(t *testing.T, raiz string) string {
	t.Helper()
	var sb strings.Builder
	_ = filepath.Walk(raiz, func(p string, i os.FileInfo, err error) error {
		if err == nil && !i.IsDir() {
			b, _ := os.ReadFile(p)
			sb.WriteString(p + "\n" + string(b) + "\n")
		}
		return nil
	})
	return sb.String()
}

// AC-E01.9: a confirmação grava o perfil na matriz e num arquivo próprio, por hash.
func TestConfirmado_ConfirmarMarcaOPerfilEGravaPorHash(t *testing.T) {
	c := novo12(t)
	b := &spyBrain{resp: respPadrao}
	c.gerar(b)
	if o := porNome(lerRegistro(t, c.vault), "skill-05").Perfil.Origem; o != matriz.OrigemInferido {
		t.Fatalf("pré-condição: skill-05 deveria ser inferido, é %q", o)
	}
	e, err := matriz.Confirmar(c.vault, "skill-05")
	if err != nil {
		t.Fatal(err)
	}
	if e.Perfil.Origem != matriz.OrigemConfirmado {
		t.Errorf("origem devolvida: %q", e.Perfil.Origem)
	}
	reg := lerRegistro(t, c.vault)
	if o := porNome(reg, "skill-05").Perfil.Origem; o != matriz.OrigemConfirmado {
		t.Errorf("registry_matrix.json: origem %q", o)
	}
	if !strings.Contains(arquivo(t, c, "matrix/skill.json"), `"origem": "confirmado"`) {
		t.Error("matrix/skill.json deve refletir a confirmação")
	}
	if !strings.Contains(arquivo(t, c, "perfil.confirmado.json"), porNome(reg, "skill-05").Hash) {
		t.Error("perfil.confirmado.json deve guardar a confirmação pelo hash do conteúdo")
	}
}

func TestConfirmado_SobreviveAoForcarSemChamarOCerebro(t *testing.T) {
	c := novo12(t)
	b := &spyBrain{resp: respPadrao}
	c.gerar(b)
	if _, err := matriz.Confirmar(c.vault, "skill-05"); err != nil {
		t.Fatal(err)
	}
	antes := chamadasSobre(b, "skill-05")
	if antes == 0 {
		t.Fatal("pré-condição: o cérebro deveria ter visto a skill-05 na 1ª geração")
	}
	s, err := matriz.Generate(t.Context(), matriz.Options{Vault: c.vault, Forcar: true}, c.deps(b))
	if err != nil {
		t.Fatal(err)
	}
	if o := porNome(lerRegistro(t, c.vault), "skill-05").Perfil.Origem; o != matriz.OrigemConfirmado {
		t.Errorf("depois do --forcar a origem é %q", o)
	}
	if got := chamadasSobre(b, "skill-05"); got != antes {
		t.Errorf("o cérebro foi chamado %d vez(es) para uma peça confirmada", got-antes)
	}
	if len(s.Invalidadas) != 0 {
		t.Errorf("nada foi invalidado: %v", s.Invalidadas)
	}
}

func TestConfirmado_DefinirGravaOsTermosDados(t *testing.T) {
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	p := matriz.PerfilDeUso{Etapas: []string{"construir"}, Naturezas: []string{"software"}, CasosDeUso: []string{"tdd"}}
	e, err := matriz.Definir(c.vault, "skill-06", p)
	if err != nil {
		t.Fatal(err)
	}
	got := porNome(lerRegistro(t, c.vault), "skill-06").Perfil
	if got.Origem != matriz.OrigemConfirmado || got.Etapas[0] != "construir" || got.Naturezas[0] != "software" || got.CasosDeUso[0] != "tdd" {
		t.Errorf("perfil gravado: %+v", got)
	}
	if e.Perfil.Origem != matriz.OrigemConfirmado {
		t.Errorf("devolvido: %+v", e.Perfil)
	}
	// e sobrevive à regeração forçada
	if _, err := matriz.Generate(t.Context(), matriz.Options{Vault: c.vault, Forcar: true}, c.deps(&spyBrain{resp: respPadrao})); err != nil {
		t.Fatal(err)
	}
	if got := porNome(lerRegistro(t, c.vault), "skill-06").Perfil; got.Naturezas[0] != "software" || got.Origem != matriz.OrigemConfirmado {
		t.Errorf("depois do --forcar: %+v", got)
	}
}

func TestConfirmado_DefinirComTermoForaDoVocabularioRecusaESemMudarNada(t *testing.T) {
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	antes := hashDaArvore(t, c.vault)
	_, err := matriz.Definir(c.vault, "skill-06", matriz.PerfilDeUso{Etapas: []string{"inventada"}, Naturezas: []string{"software"}})
	if !errors.Is(err, matriz.ErrVocabulario) {
		t.Fatalf("esperava ErrVocabulario, veio %v", err)
	}
	if hashDaArvore(t, c.vault) != antes {
		t.Error("o cofre mudou apesar do erro")
	}
	// termo de outra categoria no lugar errado também é recusado
	if _, err := matriz.Definir(c.vault, "skill-06", matriz.PerfilDeUso{Etapas: []string{"software"}}); !errors.Is(err, matriz.ErrVocabulario) {
		t.Errorf("'software' não é etapa: %v", err)
	}
}

func TestConfirmado_PecaInexistenteOuAmbiguaOuSemPerfil(t *testing.T) {
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	antes := hashDaArvore(t, c.vault)
	if _, err := matriz.Confirmar(c.vault, "nao-existe"); !errors.Is(err, matriz.ErrPecaNaoEncontrada) {
		t.Errorf("inexistente: %v", err)
	}
	if hashDaArvore(t, c.vault) != antes {
		t.Error("o cofre mudou apesar dos erros")
	}
	// matriz gerada sem cérebro: a peça sem declaração fica sem perfil, e sem perfil não há o que confirmar
	s := novo12(t)
	s.gerar(nil)
	if porNome(lerRegistro(t, s.vault), "skill-05").Perfil != nil {
		t.Fatal("pré-condição: skill-05 sem cérebro deveria ficar sem perfil")
	}
	if _, err := matriz.Confirmar(s.vault, "skill-05"); !errors.Is(err, matriz.ErrSemPerfil) {
		t.Errorf("peça sem perfil não se confirma, define-se: %v", err)
	}
	if _, err := matriz.Confirmar(t.TempDir(), "skill-05"); !errors.Is(err, matriz.ErrMatrizAusente) {
		t.Errorf("sem matriz: %v", err)
	}
}

func TestConfirmado_ConteudoAlteradoInvalidaAConfirmacaoEAvisa(t *testing.T) {
	c := novo12(t)
	b := &spyBrain{resp: respPadrao}
	c.gerar(b)
	if _, err := matriz.Confirmar(c.vault, "skill-05"); err != nil {
		t.Fatal(err)
	}
	c.alterar("skill-05", "\nmais uma linha\n")
	s, err := matriz.Generate(t.Context(), matriz.Options{Vault: c.vault}, c.deps(b))
	if err != nil {
		t.Fatal(err)
	}
	if o := porNome(lerRegistro(t, c.vault), "skill-05").Perfil.Origem; o != matriz.OrigemInferido {
		t.Errorf("com o conteúdo alterado o perfil volta a inferido, é %q", o)
	}
	if len(s.Invalidadas) != 1 || !strings.Contains(s.Invalidadas[0], "skill-05") {
		t.Errorf("o resumo deve avisar da confirmação perdida: %v", s.Invalidadas)
	}
}

func TestConfirmado_CopiasIdenticasHerdamAConfirmacao(t *testing.T) {
	c := novoCofre(t)
	a := c.skill("par", "")
	// segunda cópia: mesmo conteúdo, outro caminho
	b := a
	b.Path = filepath.Join(c.dir, "par-copia", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(b.Path), 0o755)
	dados, _ := os.ReadFile(a.Path)
	_ = os.WriteFile(b.Path, dados, 0o644)
	c.items = append(c.items, b)
	c.gerar(&spyBrain{resp: respPadrao})
	if _, err := matriz.Confirmar(c.vault, "par"); err != nil {
		t.Fatal(err)
	}
	for _, e := range lerRegistro(t, c.vault).Entradas {
		if e.Perfil.Origem != matriz.OrigemConfirmado {
			t.Errorf("%s: a cópia idêntica deve ficar confirmada junto (%q)", e.Ref, e.Perfil.Origem)
		}
	}
}

func TestConfirmado_ConfirmarPeloCaminhoEPeloNome(t *testing.T) {
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	ref := porNome(lerRegistro(t, c.vault), "skill-07").Ref
	if !strings.HasPrefix(ref, "~/") {
		t.Fatalf("pré-condição: a Ref deve ser abreviada com ~/, é %q", ref)
	}
	if _, err := matriz.Confirmar(c.vault, ref); err != nil {
		t.Errorf("pelo caminho: %v", err)
	}
}

func TestConfirmado_NomeComConteudosDiferentesEAmbiguo(t *testing.T) {
	c := novoCofre(t)
	a := c.skill("igual", "")
	b := a
	b.Path = filepath.Join(c.dir, "outra", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(b.Path), 0o755)
	_ = os.WriteFile(b.Path, []byte("conteúdo diferente\n"), 0o644)
	b.Meta = map[string]string{"sha256": "outro-hash"}
	c.items = append(c.items, b)
	c.gerar(&spyBrain{resp: respPadrao})
	antes := hashDaArvore(t, c.vault)
	_, err := matriz.Confirmar(c.vault, "igual")
	if !errors.Is(err, matriz.ErrPecaAmbigua) || !strings.Contains(err.Error(), "outra") {
		t.Errorf("esperava ErrPecaAmbigua citando os caminhos, veio %v", err)
	}
	if hashDaArvore(t, c.vault) != antes {
		t.Error("o cofre mudou apesar do erro")
	}
}

func TestConfirmado_DefinirSemTermoNenhumERecusado(t *testing.T) {
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	if _, err := matriz.Definir(c.vault, "skill-06", matriz.PerfilDeUso{}); !errors.Is(err, matriz.ErrVocabulario) {
		t.Errorf("perfil vazio: %v", err)
	}
}

func TestConfirmado_DefinirGuardaConsomeProduzELimpaOEsperaPeloCerebro(t *testing.T) {
	c := novo12(t)
	c.gerar(nil) // sem cérebro: skill-06 fica sem perfil e à espera dele
	if e := porNome(lerRegistro(t, c.vault), "skill-06"); e.Perfil != nil || !e.Retentar {
		t.Fatalf("pré-condição: sem perfil e à espera do cérebro: %+v", e)
	}
	c.gerar(&spyBrain{resp: respPadrao}) // agora o cérebro dá um perfil inferido, com consome/produz da peça
	antes := porNome(lerRegistro(t, c.vault), "skill-07").Perfil
	if len(antes.Consome) == 0 || len(antes.Produz) == 0 {
		t.Fatalf("pré-condição: consome/produz vêm da peça: %+v", antes)
	}
	e, err := matriz.Definir(c.vault, "skill-07", matriz.PerfilDeUso{Etapas: []string{"verificar"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(e.Perfil.Consome) == 0 || len(e.Perfil.Produz) == 0 {
		t.Errorf("definir não pode apagar consome/produz: %+v", e.Perfil)
	}
}

func TestConfirmado_DefinirEmPecaSemPerfilTiraDaFilaDoCerebro(t *testing.T) {
	c := novo12(t)
	c.gerar(nil)
	if _, err := matriz.Definir(c.vault, "skill-06", matriz.PerfilDeUso{Etapas: []string{"verificar"}}); err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"registry_matrix.json", "matrix/skill.json"} {
		var texto = arquivo(t, c, rel)
		i := strings.Index(texto, `"nome": "skill-06"`)
		fim := i + strings.Index(texto[i:], "}\n")
		if i < 0 || strings.Contains(texto[i:fim+1], `"retentar"`) {
			t.Errorf("%s: a peça definida pelo usuário não espera mais o cérebro", rel)
		}
	}
}

func TestConfirmado_PecaRemovidaNaoGeraAviso(t *testing.T) {
	c := novo12(t)
	b := &spyBrain{resp: respPadrao}
	c.gerar(b)
	if _, err := matriz.Confirmar(c.vault, "skill-05"); err != nil {
		t.Fatal(err)
	}
	c.remover("skill-05")
	s, err := matriz.Generate(t.Context(), matriz.Options{Vault: c.vault}, c.deps(b))
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Invalidadas) != 0 {
		t.Errorf("peça removida não é conteúdo mudado: %v", s.Invalidadas)
	}
}
