package distribuicao

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func tocar(t *testing.T, dir, rel string) {
	t.Helper()
	p := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
}

func vitrineOK(t *testing.T) (string, Vitrine) {
	t.Helper()
	dir := t.TempDir()
	v := Vitrine{
		Video: "renders/v1.mp4", SRT: "renders/legenda.srt", RotuloIA: true,
		Resumo: "Como US$ 2,7 bi sumiram.", Fontes: []string{"https://exemplo.org/a"}, Tags: []string{"financas"},
		Capitulos: []Capitulo{{"00:00", "Gancho"}, {"01:10", "Numeros"}},
		Capas:     map[string][]string{"16:9": {"capas/a.png", "capas/b.png", "capas/c.png"}},
	}
	for _, t2 := range []string{"a", "b", "c", "d", "e"} {
		v.Titulos = append(v.Titulos, NotaTitulo("Titulo "+t2, "titulo"))
	}
	for _, rel := range []string{"renders/v1.mp4", "renders/legenda.srt", "capas/a.png", "capas/b.png", "capas/c.png"} {
		tocar(t, dir, rel)
	}
	return dir, v
}

func TestNotaTitulo(t *testing.T) {
	bom := NotaTitulo("Dolar: 3 motivos da alta", "dolar")
	if bom.Nota != 10 {
		t.Fatalf("titulo ideal deveria ter 10: %+v", bom)
	}
	longo := NotaTitulo(strings.Repeat("a", 70), "a")
	if longo.Nota >= bom.Nota {
		t.Fatalf("titulo longo nao pode empatar: %+v", longo)
	}
	fora := NotaTitulo(strings.Repeat("a", 101), "a")
	if fora.Nota != 0 || len(fora.Motivos) == 0 {
		t.Fatalf("acima de 100 deve zerar: %+v", fora)
	}
	if v := NotaTitulo("  ", "a"); v.Nota != 0 {
		t.Fatalf("vazio deve zerar: %+v", v)
	}
	grito := NotaTitulo("VEJA ISSO AGORA!!", "veja")
	if grito.Nota >= bom.Nota || gritado("abc") || !gritado("ABC") || gritado("123") {
		t.Fatalf("caixa alta e exclamacao repetida perdem ponto: %+v", grito)
	}
	if comecaCom("qualquer", "") || !comecaCom("Dolar hoje", " DOLAR ") {
		t.Fatal("comecaCom incorreto")
	}
}

func TestPontuarTitulosOrdena(t *testing.T) {
	ts := PontuarTitulos([]string{"sem nada", "Juros: 5 erros", "juros"}, "juros")
	if ts[0].Texto != "Juros: 5 erros" || len(ts) != 3 {
		t.Fatalf("ordem por nota: %+v", ts)
	}
}

func TestVitrineValidaOPacote(t *testing.T) {
	dir, v := vitrineOK(t)
	if err := v.Validar(dir); err != nil {
		t.Fatalf("pacote valido reprovado: %v", err)
	}
	casos := map[string]func(*Vitrine){
		"4 titulos":         func(v *Vitrine) { v.Titulos = v.Titulos[:4] },
		"titulo vazio":      func(v *Vitrine) { v.Titulos[0].Texto = "" },
		"sem fontes":        func(v *Vitrine) { v.Fontes = nil },
		"sem tags":          func(v *Vitrine) { v.Tags = nil },
		"sem resumo":        func(v *Vitrine) { v.Resumo = " " },
		"capitulo ruim":     func(v *Vitrine) { v.Capitulos[0].Inicio = "zero" },
		"capitulo sem nome": func(v *Vitrine) { v.Capitulos[1].Titulo = "" },
		"sem capas":         func(v *Vitrine) { v.Capas = nil },
		"2 capas":           func(v *Vitrine) { v.Capas["16:9"] = v.Capas["16:9"][:2] },
		"capa ausente":      func(v *Vitrine) { v.Capas["16:9"][0] = "capas/nao.png" },
		"video ausente":     func(v *Vitrine) { v.Video = "renders/v9.mp4" },
		"srt fora":          func(v *Vitrine) { v.SRT = "../srt" },
		"srt absoluto":      func(v *Vitrine) { v.SRT = "C:/x.srt" },
	}
	for nome, mudar := range casos {
		t.Run(nome, func(t *testing.T) {
			_, v := vitrineOK(t)
			mudar(&v)
			if err := v.Validar(dir); !errors.Is(err, ErrVitrineInvalida) {
				t.Fatalf("quer ErrVitrineInvalida, obteve %v", err)
			}
		})
	}
	v.RotuloIA = false
	if err := v.Validar(dir); !errors.Is(err, ErrSemRotuloIA) {
		t.Fatalf("sem rotulo de IA: %v", err)
	}
}

func TestGravarPublicacaoEDescricao(t *testing.T) {
	dir, v := vitrineOK(t)
	desc := v.Descricao()
	for _, quer := range []string{"Capitulos:", "01:10 Numeros", "Fontes:", "apoio de IA"} {
		if !strings.Contains(desc, quer) {
			t.Fatalf("descricao sem %q:\n%s", quer, desc)
		}
	}
	semCap := v
	semCap.Capitulos, semCap.RotuloIA = nil, false
	if d := semCap.Descricao(); strings.Contains(d, "Capitulos:") || strings.Contains(d, "apoio de IA") {
		t.Fatalf("descricao minima incorreta:\n%s", d)
	}
	_ = GravarParecer(dir, Parecer{Versao: 1, Veredito: VereditoEnviar, Revisor: "r"})
	if err := Liberar(dir, 1, t0); err != nil {
		t.Fatal(err)
	}
	if err := GravarPublicacao(dir, v, 2); err != nil {
		t.Fatal(err)
	}
	dados, err := os.ReadFile(filepath.Join(dir, "publicacao.md"))
	if err != nil || !strings.Contains(string(dados), "- Titulo: "+v.Titulos[2].Texto) {
		t.Fatalf("publicacao.md incorreto: %v\n%s", err, dados)
	}
	if err := GravarPublicacao(dir, v, 9); !errors.Is(err, ErrVitrineInvalida) {
		t.Fatalf("indice fora: %v", err)
	}
	v.Fontes = nil
	if err := GravarPublicacao(dir, v, 0); !errors.Is(err, ErrVitrineInvalida) {
		t.Fatalf("pacote invalido: %v", err)
	}
}

func TestValidarTextoCapa(t *testing.T) {
	casos := []struct {
		l1, l2 string
		curto  bool
		ok     bool
	}{
		{"US$ 2,7 BI", "SUMIRAM", false, true},
		{"US$ 2,7 BI", "SUMIRAM?", true, false},
		{"", "x", false, false},
		{"UM DOIS TRES", "", true, true},
		{"UM DOIS TRES QUATRO", "CINCO", false, false},
	}
	for _, c := range casos {
		err := ValidarTextoCapa(c.l1, c.l2, c.curto)
		if (err == nil) != c.ok {
			t.Errorf("ValidarTextoCapa(%q,%q,%v)=%v, ok esperado %v", c.l1, c.l2, c.curto, err, c.ok)
		}
	}
}
