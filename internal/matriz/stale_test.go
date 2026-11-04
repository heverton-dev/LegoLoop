package matriz_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"legoloop/internal/assets"
	"legoloop/internal/matriz"
	"legoloop/internal/ste"
)

func staleDe(t *testing.T, c *cofre, reg *matriz.Registro) matriz.Stale {
	t.Helper()
	s, err := matriz.Staleness(context.Background(), reg, c.items, c.deps(nil))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func geradoNovo12(t *testing.T) (*cofre, *matriz.Registro) {
	t.Helper()
	c := novo12(t)
	c.gerar(&spyBrain{resp: respPadrao})
	return c, lerRegistro(t, c.vault)
}

// AC-E01.5: 1 peça nova, 1 alterada e 1 removida.
func TestStaleness_NovaAlteradaERemovida(t *testing.T) {
	c, reg := geradoNovo12(t)
	c.skill("skill-nova", "")
	c.alterar("skill-03", "\nmudou\n")
	c.remover("skill-09")
	s := staleDe(t, c, reg)
	if len(s.Novas) != 1 || len(s.Alteradas) != 1 || len(s.Removidas) != 1 || !s.Defasada() {
		t.Fatalf("stale = %+v", s)
	}
	if !strings.HasSuffix(s.Novas[0], "skill-nova/SKILL.md") || !strings.HasSuffix(s.Alteradas[0], "skill-03/SKILL.md") || !strings.HasSuffix(s.Removidas[0], "skill-09/SKILL.md") {
		t.Errorf("as listas devem trazer a Ref de cada peça: %+v", s)
	}
}

func TestStaleness_SemMudancaNaoEstaDefasada(t *testing.T) {
	c, reg := geradoNovo12(t)
	s := staleDe(t, c, reg)
	if s.Defasada() || s.Ausente || s.FormatoAntigo || len(s.Novas)+len(s.Alteradas)+len(s.Removidas)+len(s.Retentar) != 0 {
		t.Errorf("stale = %+v", s)
	}
}

func TestStaleness_SemArquivoDeMatrizEstaAusente(t *testing.T) {
	c := novo12(t)
	s := staleDe(t, c, nil)
	if !s.Ausente || !s.Defasada() {
		t.Errorf("stale = %+v", s)
	}
}

func TestStaleness_MatrizV1TemFormatoAntigo(t *testing.T) {
	c := novo12(t)
	v1 := `{"generated_at":"2026-11-02","matrix_version":"1.0.0","total_pieces":1,"ratings":[{"ref":"~/a/SKILL.md","name":"a","total_score":0.9,"tier":"ouro"}]}`
	_ = os.WriteFile(filepath.Join(c.vault, "registry_matrix.json"), []byte(v1), 0o644)
	s := staleDe(t, c, lerRegistro(t, c.vault))
	if !s.FormatoAntigo || !s.Defasada() || s.Ausente {
		t.Errorf("stale = %+v", s)
	}
}

func TestStaleness_PecaComFalhaAnteriorContaComoDefasagem(t *testing.T) {
	c := novoCofre(t)
	c.skill("boa", "")
	c.skill("quebrada", "")
	c.gerar(&spyBrain{resp: respPadrao, falhaEm: "quebrada"})
	s := staleDe(t, c, lerRegistro(t, c.vault))
	if len(s.Retentar) != 1 || !strings.HasSuffix(s.Retentar[0], "quebrada/SKILL.md") || !s.Defasada() {
		t.Errorf("stale = %+v", s)
	}
}

func TestStaleness_ScriptAlteradoEHookPorIdentidade(t *testing.T) {
	c := novoCofre(t)
	sc := c.outro(assets.KindScript, "sc", ".ps1", "Write-Host um")
	hk := c.outro(assets.KindHook, "hk", ".json", `{"a":1}`)
	c.gerar(&spyBrain{resp: respPadrao})
	reg := lerRegistro(t, c.vault)

	_ = os.WriteFile(hk.Path, []byte(`{"a":2,"segredo":"x"}`), 0o644) // o arquivo de configuração muda por motivos alheios
	if s := staleDe(t, c, reg); s.Defasada() {
		t.Errorf("mudar o arquivo de configuração do hook não muda a peça: %+v", s)
	}
	for i, it := range c.items {
		if it.Name == "hk" {
			it.Description = "Outra descrição"
			c.items[i] = it
		}
	}
	if s := staleDe(t, c, reg); len(s.Alteradas) != 1 {
		t.Errorf("a identidade do hook mudou: %+v", s)
	}
	_ = os.WriteFile(sc.Path, []byte("Write-Host dois"), 0o644)
	if s := staleDe(t, c, reg); len(s.Alteradas) != 2 {
		t.Errorf("o script alterado também: %+v", s)
	}
}

func TestStaleness_NaoReleSkillNemChamaCerebro(t *testing.T) {
	c, reg := geradoNovo12(t)
	var lidos []string
	d := c.deps(nil) // nil: qualquer chamada ao cérebro seria pânico
	d.ReadFile = func(p string) ([]byte, error) { lidos = append(lidos, p); return os.ReadFile(p) }
	if _, err := matriz.Staleness(context.Background(), reg, c.items, d); err != nil {
		t.Fatal(err)
	}
	for _, p := range lidos {
		if strings.HasSuffix(p, "SKILL.md") {
			t.Errorf("a defasagem não relê skill (a varredura traz o sha256): %s", p)
		}
	}
}

func TestStaleness_ListasOrdenadasEDeterministicas(t *testing.T) {
	c, reg := geradoNovo12(t)
	for _, n := range []string{"z-nova", "a-nova", "m-nova"} {
		c.skill(n, "")
	}
	s := staleDe(t, c, reg)
	if !slices.IsSorted(s.Novas) || len(s.Novas) != 3 {
		t.Errorf("Novas deve estar ordenada: %v", s.Novas)
	}
}

func TestStaleness_ContextoCancelado(t *testing.T) {
	c, reg := geradoNovo12(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := matriz.Staleness(ctx, reg, c.items, c.deps(nil)); err == nil {
		t.Fatal("contexto cancelado deve voltar erro")
	}
}

func TestStaleness_MensagensSeguemOPerfilSTEPT(t *testing.T) {
	casos := []matriz.Stale{
		{Ausente: true},
		{FormatoAntigo: true},
		{Novas: []string{"a"}, Alteradas: []string{"b"}, Removidas: []string{"c"}, Retentar: []string{"d"}},
	}
	for _, s := range casos {
		m := s.Mensagem()
		if m == "" || len(ste.Lint(m)) != 0 || !strings.Contains(m, "matriz gerar") {
			t.Errorf("mensagem %q (violações: %v)", m, ste.Lint(m))
		}
	}
	if (matriz.Stale{}).Mensagem() != "" {
		t.Error("matriz em dia não gera mensagem")
	}
	if m := casos[2].Mensagem(); !strings.Contains(m, "1 nova") || !strings.Contains(m, "1 alterada") || !strings.Contains(m, "1 removida") {
		t.Errorf("a mensagem deve dizer os números: %q", m)
	}
}

// AC-NF.4: 1.600 peças; defasagem em até 2 s e geração sem mudança em até 5 s.
func TestStaleness_DesempenhoCom1600Pecas(t *testing.T) {
	c := novoCofre(t)
	for i := 0; i < 1600; i++ {
		nome := fmt.Sprintf("skill-%04d", i)
		c.items = append(c.items, assets.Item{Kind: assets.KindSkill, Name: nome, Description: "peça peculiar " + nome,
			Path: filepath.Join(c.dir, nome, "SKILL.md"), Meta: map[string]string{"sha256": fmt.Sprintf("h%04d", i)}})
	}
	d := c.deps(&spyBrain{resp: respPadrao})
	d.ReadFile = func(string) ([]byte, error) { return []byte("# peça"), nil }
	t1 := time.Now()
	if _, err := matriz.Generate(context.Background(), matriz.Options{Vault: c.vault}, d); err != nil {
		t.Fatal(err)
	}
	if dt := time.Since(t1); dt > 10*time.Second {
		t.Errorf("primeira geração de 1.600 peças levou %v (guarda: 10 s); inferência ou nota lenta demais", dt)
	}
	reg := lerRegistro(t, c.vault)

	t0 := time.Now()
	s, err := matriz.Staleness(context.Background(), reg, c.items, d)
	if dt := time.Since(t0); err != nil || dt > 2*time.Second || s.Defasada() {
		t.Errorf("defasagem de 1.600 peças: %v (limite 2 s), stale=%+v, err=%v", dt, s, err)
	}
	t0 = time.Now()
	sum, err := matriz.Generate(context.Background(), matriz.Options{Vault: c.vault}, d)
	if dt := time.Since(t0); err != nil || dt > 5*time.Second || sum.Reprocessadas != 0 {
		t.Errorf("geração sem mudança com 1.600 peças: %v (limite 5 s), reprocessadas=%d, err=%v", dt, sum.Reprocessadas, err)
	}
}

func TestStaleness_RemovidasVemOrdenadasMesmoSaindoDeUmMapa(t *testing.T) {
	c, reg := geradoNovo12(t)
	for _, n := range []string{"skill-02", "skill-11", "skill-05", "skill-08", "skill-01"} {
		c.remover(n)
	}
	for i := 0; i < 20; i++ { // a ordem de um mapa muda de uma volta para outra
		if s := staleDe(t, c, reg); len(s.Removidas) != 5 || !slices.IsSorted(s.Removidas) {
			t.Fatalf("Removidas deve vir ordenada: %v", s.Removidas)
		}
	}
}

func TestStaleness_MensagemDizAsPecasATentarDeNovo(t *testing.T) {
	m := matriz.Stale{Retentar: []string{"a", "b"}}.Mensagem()
	if !strings.Contains(m, "2 esperando o cérebro") || len(ste.Lint(m)) != 0 {
		t.Errorf("mensagem = %q", m)
	}
	if m := (matriz.Stale{Novas: []string{"a", "b", "c"}}).Mensagem(); !strings.Contains(m, "3 novas") {
		t.Errorf("plural: %q", m)
	}
}

func TestStaleness_MatrizV2VaziaNaoEFormatoAntigo(t *testing.T) {
	c := novoCofre(t) // nenhuma peça
	c.gerar(&spyBrain{resp: respPadrao})
	s := staleDe(t, c, lerRegistro(t, c.vault))
	if s.FormatoAntigo || s.Defasada() {
		t.Errorf("uma matriz v2 sem peças está em dia: %+v", s)
	}
}

func TestStaleness_MensagemExataComNumerosEmSingularEPlural(t *testing.T) {
	um := matriz.Stale{Novas: []string{"a"}, Alteradas: []string{"b"}, Removidas: []string{"c"}}.Mensagem()
	if quer := "A matriz de peças está desatualizada: 1 nova, 1 alterada, 1 removida. Rode legoloop matriz gerar."; um != quer {
		t.Errorf("singular:\n got %q\nquer %q", um, quer)
	}
	varios := matriz.Stale{Novas: []string{"a", "b"}, Alteradas: []string{"c", "d", "e"}, Removidas: []string{"f", "g"}}.Mensagem()
	if quer := "A matriz de peças está desatualizada: 2 novas, 3 alteradas, 2 removidas. Rode legoloop matriz gerar."; varios != quer {
		t.Errorf("plural:\n got %q\nquer %q", varios, quer)
	}
}

func TestPrevia_QuantasPecasSeriamRecatalogadas(t *testing.T) {
	c := novo12(t)
	ctx := context.Background()
	if n, err := matriz.Previa(ctx, nil, c.items, c.deps(nil)); err != nil || n != 15 {
		t.Errorf("sem matriz todas seriam: %d %v", n, err)
	}
	c.gerar(&spyBrain{resp: respPadrao})
	reg := lerRegistro(t, c.vault)
	if n, _ := matriz.Previa(ctx, reg, c.items, c.deps(nil)); n != 0 {
		t.Errorf("em dia: %d", n)
	}
	c.skill("nova", "")
	c.alterar("skill-02", "\nx\n")
	c.remover("skill-05") // remoção não gera trabalho de catalogação
	if n, _ := matriz.Previa(ctx, reg, c.items, c.deps(nil)); n != 2 {
		t.Errorf("1 nova + 1 alterada = 2, veio %d", n)
	}
	v1 := &matriz.Registro{Versao: "1.0.0", Ratings: lerRegistro(t, c.vault).Ratings[:1]}
	if n, _ := matriz.Previa(ctx, v1, c.items, c.deps(nil)); n != 15 {
		t.Errorf("a v1 não tem perfil: todas seriam (%d)", n)
	}
}

func TestPrevia_CopiasIdenticasContamUmaVezSo(t *testing.T) {
	c := novoCofre(t)
	base := c.skill("copiada", "")
	for _, r := range []string{"b", "c"} {
		p := filepath.Join(c.dir, r, "copiada", "SKILL.md")
		_ = os.MkdirAll(filepath.Dir(p), 0o755)
		bb, _ := os.ReadFile(base.Path)
		_ = os.WriteFile(p, bb, 0o644)
		cp := base
		cp.Path = p
		c.items = append(c.items, cp)
	}
	c.skill("outra", "")
	ctx := context.Background()
	if n, _ := matriz.Previa(ctx, nil, c.items, c.deps(nil)); n != 2 {
		t.Errorf("sem matriz: 3 cópias idênticas e 1 peça distinta = 2 conteúdos, veio %d", n)
	}
	c.gerar(&spyBrain{resp: respPadrao})
	reg := lerRegistro(t, c.vault)
	c.alterar("copiada", "\nmudou\n") // só o 1º caminho muda
	// o caminho alterado vira conteúdo novo; as outras duas cópias seguem iguais
	if n, _ := matriz.Previa(ctx, reg, c.items, c.deps(nil)); n != 1 {
		t.Errorf("1 conteúdo novo = 1, veio %d", n)
	}
}

func TestPrevia_PecaEsperandoOCerebroEntraNaConta(t *testing.T) {
	c := novoCofre(t)
	c.skill("sem-pista", "")
	c.gerarCom(nil, &contador{}, matriz.Options{}) // sem cérebro: a peça fica esperando
	reg := lerRegistro(t, c.vault)
	if n, _ := matriz.Previa(context.Background(), reg, c.items, c.deps(nil)); n != 1 {
		t.Errorf("a peça esperando o cérebro seria recatalogada: %d", n)
	}
	if n, _ := matriz.Previa(context.Background(), lerRegistro(t, c.vault), nil, c.deps(nil)); n != 0 {
		t.Errorf("sem peças não há o que recatalogar: %d", n)
	}
}

func TestPrevia_TiposDiferentesComMesmoConteudoSaoDoisConteudos(t *testing.T) {
	c := novoCofre(t)
	sk := c.skill("igual", "")
	texto, _ := os.ReadFile(sk.Path)
	c.outro(assets.KindScript, "igual-script", ".ps1", string(texto))
	c.items[0].Meta["sha256"] = "" // o hash da skill passa a ser o do texto, igual ao do script
	if n, _ := matriz.Previa(context.Background(), nil, c.items, c.deps(nil)); n != 2 {
		t.Errorf("skill e script com o mesmo texto são 2 conteúdos: %d", n)
	}
}
