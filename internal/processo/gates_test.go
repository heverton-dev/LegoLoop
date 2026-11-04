package processo_test

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/gate"
	"legoloop/internal/processo"
	"legoloop/internal/ste"
)

func comPromessa(tipo, caminho, g string) processo.Processo {
	p := valido()
	p.Etapas = []processo.Etapa{etapa("e01")}
	p.Etapas[0].Produz = processo.Artefato{Caminho: caminho, Tipo: tipo, Promete: true}
	p.Etapas[0].Gate = g
	return p
}

func TestGateTrivialExistenciaIsoladaComPromessaEhRecusada(t *testing.T) {
	err := comPromessa("markdown", "out/x.md", "file_exists('out/x.md')").Validate()
	if !errors.Is(err, processo.ErrGateTrivial) || !errors.Is(err, processo.ErrInvalido) {
		t.Fatalf("esperava ErrGateTrivial e ErrInvalido, veio %v", err)
	}
	msg := err.Error()
	for _, want := range []string{"e01", "word_count('out/x.md') >= 30", "has_pattern"} {
		if !strings.Contains(msg, want) {
			t.Errorf("a mensagem deve citar %q e sugerir o gate da biblioteca:\n%s", want, msg)
		}
	}
}

func TestGateTrivialOutrosGatesSemConteudoTambemSaoRecusados(t *testing.T) {
	for nome, g := range map[string]string{
		"exit code":      "exit_code == 0",
		"existe ou":      "file_exists('out/x.md') || word_count('out/x.md') >= 30",
		"existe e saída": "file_exists('out/x.md') && exit_code == 0",
	} {
		t.Run(nome, func(t *testing.T) {
			if err := comPromessa("markdown", "out/x.md", g).Validate(); !errors.Is(err, processo.ErrGateTrivial) {
				t.Errorf("gate %q deveria ser trivial: %v", g, err)
			}
		})
	}
}

func TestGateTrivialSemPromessaOuMarcadorAceitaExistencia(t *testing.T) {
	semPromessa := comPromessa("markdown", "out/x.md", "file_exists('out/x.md')")
	semPromessa.Etapas[0].Produz.Promete = false
	if err := semPromessa.Validate(); err != nil {
		t.Errorf("sem promessa o gate de existência vale: %v", err)
	}
	marcador := comPromessa("marcador", "out/ok.flag", "file_exists('out/ok.flag')")
	if err := marcador.Validate(); err != nil {
		t.Errorf("marcador só existe: %v", err)
	}
}

func TestGateTrivialGateHumanoNaoEhTrivial(t *testing.T) {
	if err := comPromessa("markdown", "out/x.md", "prompt('O texto está bom?')").Validate(); err != nil {
		t.Errorf("a aprovação de uma pessoa confere o conteúdo: %v", err)
	}
}

func TestGateTrivialBibliotecaGeraOGateEsperado(t *testing.T) {
	for _, c := range []struct {
		art  processo.Artefato
		want string
	}{
		{processo.Artefato{Caminho: "out/x.md", Tipo: "markdown"},
			`file_exists('out/x.md') && word_count('out/x.md') >= 30 && !has_pattern('out/x.md', '\{\{')`},
		{processo.Artefato{Caminho: "out/x.json", Tipo: "json", Esquema: "schemas/x.json"},
			`file_exists('out/x.json') && schema_valid('schemas/x.json', 'out/x.json')`},
		{processo.Artefato{Caminho: "out/x.json", Tipo: "json"},
			`file_exists('out/x.json') && has_pattern('out/x.json', '^\s*[\{\[]')`},
		{processo.Artefato{Caminho: "src/x.go", Tipo: "codigo"}, `file_exists('src/x.go') && word_count('src/x.go') >= 1`},
		{processo.Artefato{Caminho: "out/x.csv", Tipo: "dados"}, `file_exists('out/x.csv') && line_count('out/x.csv') >= 2`},
		{processo.Artefato{Caminho: "out/ok.flag", Tipo: "marcador"}, `file_exists('out/ok.flag')`},
	} {
		got, err := processo.GateParaArtefato(c.art)
		if err != nil || got != c.want {
			t.Errorf("%s: gate = %q (%v), queria %q", c.art.Tipo, got, err, c.want)
		}
		if _, err := gate.Compile(got); err != nil {
			t.Errorf("o gate da biblioteca deve compilar em internal/gate: %q → %v", got, err)
		}
	}
}

func TestGateTrivialGateDaBibliotecaPassaNaValidacao(t *testing.T) {
	for _, a := range []processo.Artefato{
		{Caminho: "out/x.md", Tipo: "markdown", Promete: true},
		{Caminho: "out/x.json", Tipo: "json", Promete: true, Esquema: "schemas/x.json"},
		{Caminho: "out/x.json", Tipo: "json", Promete: true},
		{Caminho: "src/x.go", Tipo: "codigo", Promete: true},
		{Caminho: "out/x.csv", Tipo: "dados", Promete: true},
	} {
		g, err := processo.GateParaArtefato(a)
		if err != nil {
			t.Fatal(err)
		}
		if err := comPromessa(a.Tipo, a.Caminho, g).Validate(); err != nil {
			t.Errorf("%s: o gate da biblioteca deve passar na validação: %v", a.Tipo, err)
		}
	}
}

// avalia roda o gate contra arquivos reais numa pasta descartável.
func avalia(t *testing.T, g string, arquivos map[string]string) bool {
	t.Helper()
	dir := t.TempDir()
	for nome, texto := range arquivos {
		p := filepath.Join(dir, filepath.FromSlash(nome))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(texto), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer root.Close()
	e, err := gate.Compile(g)
	if err != nil {
		t.Fatal(err)
	}
	ok, _ := e.Eval(gate.Env{Root: root})
	return ok
}

func TestGateTrivialGateMarkdownConfereOConteudoDeVerdade(t *testing.T) {
	g, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/x.md", Tipo: "markdown"})
	longo := strings.Repeat("palavra ", 40)
	for nome, c := range map[string]struct {
		arq  map[string]string
		want bool
	}{
		"texto suficiente":  {map[string]string{"out/x.md": longo}, true},
		"curto demais":      {map[string]string{"out/x.md": "só dez palavras aqui"}, false},
		"vazio":             {map[string]string{"out/x.md": ""}, false},
		"placeholder sobra": {map[string]string{"out/x.md": longo + " {{ NOME }}"}, false},
		"arquivo ausente":   {map[string]string{}, false},
		"exatamente 30":     {map[string]string{"out/x.md": strings.Repeat("a ", 30)}, true},
		"29 palavras":       {map[string]string{"out/x.md": strings.Repeat("a ", 29)}, false},
	} {
		if got := avalia(t, g, c.arq); got != c.want {
			t.Errorf("%s: gate = %v, queria %v", nome, got, c.want)
		}
	}
}

func TestGateTrivialGateJSONConfereOFormato(t *testing.T) {
	com, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/x.json", Tipo: "json", Esquema: "schemas/x.json"})
	esquema := `{"type":"object","required":["preco"],"properties":{"preco":{"type":"number"}}}`
	if !avalia(t, com, map[string]string{"schemas/x.json": esquema, "out/x.json": `{"preco": 9.9}`}) {
		t.Error("JSON que cumpre o esquema deve passar")
	}
	if avalia(t, com, map[string]string{"schemas/x.json": esquema, "out/x.json": `{"preco": "caro"}`}) {
		t.Error("JSON fora do esquema deve falhar")
	}
	sem, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/x.json", Tipo: "json"})
	if !avalia(t, sem, map[string]string{"out/x.json": "  \n[1,2]"}) || !avalia(t, sem, map[string]string{"out/x.json": `{"a":1}`}) {
		t.Error("sem esquema, objeto ou lista JSON passam")
	}
	if avalia(t, sem, map[string]string{"out/x.json": "isto não é json"}) {
		t.Error("texto solto não passa como JSON")
	}
}

func TestGateTrivialBibliotecaRecusaEntradaInsegura(t *testing.T) {
	for nome, a := range map[string]processo.Artefato{
		"aspa no caminho":   {Caminho: "out/a'b.md", Tipo: "markdown"},
		"barra invertida":   {Caminho: `out\x.md`, Tipo: "markdown"},
		"caminho absoluto":  {Caminho: "/etc/x.md", Tipo: "markdown"},
		"sobe pasta":        {Caminho: "../x.md", Tipo: "markdown"},
		"vazio":             {Caminho: "", Tipo: "markdown"},
		"tipo desconhecido": {Caminho: "out/x.md", Tipo: "planilha"},
		"tipo vazio":        {Caminho: "out/x.md"},
		"esquema inseguro":  {Caminho: "out/x.json", Tipo: "json", Esquema: "../s.json"},
		"esquema com aspa":  {Caminho: "out/x.json", Tipo: "json", Esquema: "s'.json"},
		"crase no caminho":  {Caminho: "out/a`b.md", Tipo: "markdown"},
	} {
		t.Run(nome, func(t *testing.T) {
			if g, err := processo.GateParaArtefato(a); err == nil {
				t.Errorf("devia recusar, gerou %q", g)
			}
		})
	}
}

func TestGateTrivialMensagemSegueOPerfilSTE(t *testing.T) {
	err := comPromessa("markdown", "out/x.md", "file_exists('out/x.md')").Validate()
	for _, frase := range strings.Split(strings.TrimPrefix(err.Error(), processo.ErrInvalido.Error()+": "), ". ") {
		if v := ste.Lint(frase); len(v) != 0 {
			t.Errorf("mensagem fora do perfil STE-PT: %q → %v", frase, v)
		}
	}
}

func TestGateTrivialGateDeDadosConfereRegistros(t *testing.T) {
	g, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/x.csv", Tipo: "dados"})
	for nome, c := range map[string]struct {
		texto string
		want  bool
	}{
		"cabeçalho e um registro":      {"loja,preco\nA,9.90\n", true},
		"cabeçalho e muitos":           {"loja,preco\nA,1\nB,2\nC,3\n", true},
		"só o cabeçalho":               {"loja,preco\n", false},
		"vazio":                        {"", false},
		"uma linha sem quebra":         {"loja,preco", false},
		"duas linhas sem quebra final": {"loja,preco\nA,9.90", true},
	} {
		if got := avalia(t, g, map[string]string{"out/x.csv": c.texto}); got != c.want {
			t.Errorf("%s: gate = %v, queria %v", nome, got, c.want)
		}
	}
	if avalia(t, g, map[string]string{}) {
		t.Error("arquivo ausente deve falhar")
	}
}

func TestGateTrivialRouteUsaOGateDaEtapa(t *testing.T) {
	p := quatro()
	g, _ := processo.GateParaArtefato(processo.Artefato{Caminho: "out/e02.md", Tipo: "markdown"})
	p.Etapas[1].Produz = processo.Artefato{Caminho: "out/e02.md", Tipo: "markdown", Promete: true}
	p.Etapas[1].Gate = g
	pr := projetar(t, p)
	var d struct {
		Nodes []struct {
			Type, Label string
			Output      struct{ Gate string }
		}
	}
	if err := json.Unmarshal(pr.Route, &d); err != nil {
		t.Fatal(err)
	}
	achou := false
	for _, n := range d.Nodes {
		if n.Type == "unit_task" && strings.Contains(n.Label, "`e02`") {
			achou = true
			if n.Output.Gate != g {
				t.Errorf("o ticket da etapa deve usar o gate da etapa:\n%s\nquando era\n%s", n.Output.Gate, g)
			}
		}
	}
	if !achou {
		t.Fatal("ticket da etapa e02 não encontrado na ROUTE")
	}
	if !strings.Contains(pr.Plan, g) {
		t.Error("o PLAN espelha o gate da ROUTE")
	}
}
