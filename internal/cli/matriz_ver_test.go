package cli

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/assets"
)

type verJSON struct {
	TotalFound int `json:"total_found"`
	Entradas   []struct {
		Nome   string   `json:"nome"`
		Tipo   string   `json:"tipo"`
		Tier   string   `json:"tier"`
		Score  float64  `json:"total_score"`
		Origem string   `json:"origem"`
		Copias int      `json:"copias"`
		Refs   []string `json:"refs"`
	} `json:"entradas"`
	Ratings []map[string]any `json:"ratings"`
}

// matrizDeCincoCaminhos monta o cenário do AC-E01.10: 3 cópias idênticas da skill x (declarada), 1 script
// (inferido) e 1 skill y (definida pelo usuário: software, construir, tdd).
func matrizDeCincoCaminhos(t *testing.T) (*ambiente, *MatrizDeps) {
	t.Helper()
	a := novoAmbiente(t)
	var items []assets.Item
	for _, raiz := range a.tresCopias() {
		p := filepath.Join(raiz, "x", "SKILL.md")
		items = append(items, assets.Item{Kind: assets.KindSkill, Name: "x", Description: "Faz x.", Path: p, Meta: map[string]string{"sha256": "h-x"}})
	}
	items = append(items, a.item("y", ""))
	sp := filepath.Join(a.home, "scripts", "limpa.sh")
	_ = os.MkdirAll(filepath.Dir(sp), 0o755)
	_ = os.WriteFile(sp, []byte("echo limpa\n"), 0o644)
	items = append(items, assets.Item{Kind: assets.KindScript, Name: "limpa", Description: "Limpa.", Path: sp, Meta: map[string]string{}})
	md := itemsFixos(items...)
	md.Brain, md.BrainName = &brainFake{resp: respMatriz}, "fake/x"
	if code, _, errb := a.rodar(Deps{Matriz: md}, "matriz", "gerar", "--vault", a.vault); code != 0 {
		t.Fatalf("gerar: exit %d: %s", code, errb)
	}
	if code, _, errb := a.rodar(Deps{Matriz: md}, "matriz", "perfil", "definir", "y", "--natureza", "software", "--etapa", "construir", "--caso", "tdd", "--vault", a.vault); code != 0 {
		t.Fatalf("definir: exit %d: %s", code, errb)
	}
	return a, md
}

func (a *ambiente) ver(md *MatrizDeps, args ...string) (int, string, string) {
	a.t.Helper()
	return a.rodar(Deps{Matriz: md}, append([]string{"matriz", "ver", "--vault", a.vault}, args...)...)
}

func (a *ambiente) verJSON(md *MatrizDeps, args ...string) verJSON {
	a.t.Helper()
	code, out, errb := a.ver(md, append(args, "--json")...)
	if code != 0 {
		a.t.Fatalf("exit %d: %s", code, errb)
	}
	var v verJSON
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		a.t.Fatalf("o stdout do --json não é JSON puro: %v\n%s", err, out)
	}
	return v
}

func nomesDe(v verJSON) string {
	var ns []string
	for _, e := range v.Entradas {
		ns = append(ns, e.Nome)
	}
	return strings.Join(ns, ",")
}

// linhasDeDados devolve as linhas da tabela, sem cabeçalho, sem título e sem linhas vazias.
func linhasDeDados(out string) []string {
	var ls []string
	for _, l := range strings.Split(out, "\n") {
		if l = strings.TrimSpace(l); l != "" && !strings.HasPrefix(l, "===") && !strings.HasPrefix(l, "TIER") && !strings.HasPrefix(l, "----") {
			ls = append(ls, l)
		}
	}
	return ls
}

// AC-E01.10: as 3 cópias aparecem numa linha com a contagem; --copias mostra um caminho por linha.
func TestMatrizVer_AgrupaCopiasEMostraUmaLinhaPorCaminhoComCopias(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	code, out, errb := a.ver(md)
	if code != 0 {
		t.Fatalf("exit %d: %s", code, errb)
	}
	ls := linhasDeDados(out)
	if len(ls) != 3 {
		t.Fatalf("esperava 3 linhas (x agrupada, y, limpa), vieram %d:\n%s", len(ls), out)
	}
	if !strings.Contains(out, "(3 ENCONTRADAS)") {
		t.Errorf("o título conta as linhas mostradas:\n%s", out)
	}
	for _, l := range ls {
		if strings.Fields(l)[2] == "x" && !strings.HasSuffix(l, " 3") {
			t.Errorf("a linha de x deve terminar com a contagem 3: %q", l)
		}
	}
	code, out, _ = a.ver(md, "--copias")
	if ls = linhasDeDados(out); code != 0 || len(ls) != 5 {
		t.Fatalf("--copias: exit %d, %d linhas:\n%s", code, len(ls), out)
	}
	if !strings.Contains(out, filepath.ToSlash(".codex/skills/x/SKILL.md")) {
		t.Errorf("--copias mostra o caminho de cada peça:\n%s", out)
	}
}

func TestMatrizVer_TabelaTrazTipoPerfilEOrigem(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	_, out, _ := a.ver(md)
	for _, quer := range []string{"TIPO", "NATUREZAS", "ETAPAS", "CASOS", "ORIGEM", "CÓPIAS", "software", "construir", "tdd", "confirmado", "declarado", "inferido", "script"} {
		if !strings.Contains(out, quer) {
			t.Errorf("falta %q na tabela:\n%s", quer, out)
		}
	}
}

func TestMatrizVer_FiltrosReduzemOResultado(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	casos := []struct {
		args []string
		quer string
	}{
		{[]string{"--tipo", "script"}, "limpa"},
		{[]string{"--tipo", "skill"}, "x,y"},
		{[]string{"--natureza", "software"}, "y"},
		{[]string{"--natureza", "dados"}, "x,limpa"},
		{[]string{"--etapa", "construir"}, "y"},
		{[]string{"--etapa", "coletar"}, "x,limpa"},
		{[]string{"--caso", "tdd"}, "y"},
		{[]string{"--origem", "confirmado"}, "y"},
		{[]string{"--origem", "declarado"}, "x"},
		{[]string{"--origem", "inferido"}, "limpa"},
		{[]string{"--tipo", "skill", "--natureza", "dados"}, "x"},
	}
	for _, c := range casos {
		got := a.verJSON(md, c.args...)
		if nomesDe(got) != c.quer {
			t.Errorf("%v: veio %q, esperava %q", c.args, nomesDe(got), c.quer)
		}
		if got.TotalFound != len(got.Entradas) {
			t.Errorf("%v: total_found=%d com %d entradas", c.args, got.TotalFound, len(got.Entradas))
		}
	}
}

func TestMatrizVer_TierEMinFiltramPelaNota(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	todas := a.verJSON(md)
	if len(todas.Entradas) != 3 {
		t.Fatalf("pré-condição: 3 entradas, vieram %d", len(todas.Entradas))
	}
	alvo := todas.Entradas[0]
	var querTier []string
	maximo := 0.0
	for _, e := range todas.Entradas {
		if e.Tier == alvo.Tier {
			querTier = append(querTier, e.Nome)
		}
		maximo = max(maximo, e.Score)
	}
	if got := nomesDe(a.verJSON(md, "--tier", alvo.Tier)); got != strings.Join(querTier, ",") {
		t.Errorf("--tier %s: veio %q, esperava %q", alvo.Tier, got, strings.Join(querTier, ","))
	}
	if got := a.verJSON(md, "--min", "0.0001"); len(got.Entradas) != 3 {
		t.Errorf("--min baixo não corta nada: %d", len(got.Entradas))
	}
	if got := a.verJSON(md, "--min", jsonNum(maximo+0.001)); len(got.Entradas) != 0 {
		t.Errorf("--min acima da maior nota não deixa nada: %d", len(got.Entradas))
	}
}

func jsonNum(f float64) string { b, _ := json.Marshal(f); return string(b) }

func TestMatrizVer_JSONTemCopiasERefsEMantemRatings(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	v := a.verJSON(md)
	if len(v.Ratings) != 3 {
		t.Errorf("o campo ratings (compatível com o formato antigo) traz uma nota por linha: %d", len(v.Ratings))
	}
	for _, e := range v.Entradas {
		if e.Nome == "x" && (e.Copias != 3 || len(e.Refs) != 3) {
			t.Errorf("x: copias=%d refs=%v", e.Copias, e.Refs)
		}
		if e.Nome != "x" && e.Copias != 1 {
			t.Errorf("%s: copias=%d", e.Nome, e.Copias)
		}
	}
	if sem := a.verJSON(md, "--copias"); sem.TotalFound != 5 {
		t.Errorf("--copias --json: total_found=%d", sem.TotalFound)
	}
}

func TestMatrizVer_FalhasDeUsoSaem2EListamOAceito(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	casos := []struct {
		args []string
		quer string
	}{
		{[]string{"--tier", "platina"}, "ouro, prata, bronze, toxica"},
		{[]string{"--tipo", "xyz"}, "skill, mcp, script, hook, gate, processo"},
		{[]string{"--origem", "xyz"}, "declarado, inferido, confirmado, nenhuma"},
		{[]string{"--natureza", "inventada"}, "legoloop matriz vocab listar"},
		{[]string{"--etapa", "software"}, "legoloop matriz vocab listar"},
		{[]string{"--min", "2"}, "entre 0 e 1"},
	}
	for _, c := range casos {
		code, out, errb := a.ver(md, c.args...)
		if code != 2 || !strings.Contains(errb, c.quer) || out != "" {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", c.args, code, out, errb)
		}
	}
}

func TestMatrizVer_SemMatrizSai1DizendoParaGerar(t *testing.T) {
	a := novoAmbiente(t)
	code, out, errb := a.rodar(Deps{}, "matriz", "ver", "--vault", a.vault)
	if code != 1 || out != "" || !strings.Contains(errb, "legoloop matriz gerar") {
		t.Errorf("exit %d, stdout %q, stderr %q", code, out, errb)
	}
}

func TestMatrizVer_UsoIncorretoComArgumentoSobrando(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	if code, _, errb := a.ver(md, "sobrando"); code != 2 || !strings.Contains(errb, "Uso: legoloop matriz") {
		t.Errorf("exit %d: %s", code, errb)
	}
	if code, _, _ := a.ver(md, "--flag-que-nao-existe"); code != 2 {
		t.Errorf("flag desconhecida: exit %d", code)
	}
}

// matrix é alias descontinuado de matriz ver: mesmo resultado, aviso no stderr.
func TestMatrizVer_MatrixEAliasDescontinuado(t *testing.T) {
	a, md := matrizDeCincoCaminhos(t)
	reg := filepath.Join(a.vault, "registry_matrix.json")
	_, quer, _ := a.rodar(Deps{Matriz: md}, "matriz", "ver", "--registro", reg, "--tipo", "skill")
	code, out, errb := a.rodar(Deps{Matriz: md}, "matrix", "--path", reg, "--kind", "skill")
	if code != 0 || out != quer {
		t.Fatalf("exit %d; o alias deve imprimir o mesmo que matriz ver.\nalias:\n%s\nver:\n%s", code, out, quer)
	}
	if !strings.Contains(errb, "descontinuado") || !strings.Contains(errb, "legoloop matriz ver") {
		t.Errorf("o alias avisa no stderr: %q", errb)
	}
	_, out, _ = a.rodar(Deps{Matriz: md}, "matrix", "--path="+reg, "--kind=script", "--json")
	var v verJSON
	if err := json.Unmarshal([]byte(out), &v); err != nil || nomesDe(v) != "limpa" {
		t.Errorf("a forma --flag=valor também é traduzida: %v %q", err, out)
	}
}

// Uma matriz da v1 (só `ratings`) continua legível: o tipo vem do caminho, como antes.
func TestMatrizVer_LeMatrizV1SemEntradas(t *testing.T) {
	a := novoAmbiente(t)
	p := filepath.Join(a.home, "v1.json")
	v1 := `{"generated_at":"2026-11-03","matrix_version":"1.0","total_pieces":2,"ratings":[
	 {"ref":"skills/ouro","name":"ouro","total_score":0.95,"tier":"ouro"},
	 {"ref":"skills/prata","name":"prata","total_score":0.7,"tier":"prata"},
	 {"ref":"scripts/ruim.ps1","name":"ruim","total_score":0.25,"tier":"toxica"}]}`
	_ = os.WriteFile(p, []byte(v1), 0o644)
	code, out, errb := a.rodar(Deps{}, "matriz", "ver", "--registro", p, "--tipo", "skill", "--json")
	var v verJSON
	if err := json.Unmarshal([]byte(out), &v); code != 0 || err != nil || nomesDe(v) != "ouro,prata" {
		t.Errorf("exit %d err %v: %s %s", code, err, out, errb)
	}
}

// registroManual grava uma matriz v2 feita à mão: notas e tiers variados, uma peça sem perfil e o mesmo
// hash em dois tipos (que não são cópias uma da outra).
func registroManual(t *testing.T, a *ambiente) {
	t.Helper()
	reg := `{"generated_at":"2026-11-04","matrix_version":"2.0.0","total_pieces":5,
	"ratings":[
	 {"ref":"~/a","name":"a","total_score":0.95,"tier":"ouro"},
	 {"ref":"~/b","name":"b","total_score":0.20,"tier":"toxica"},
	 {"ref":"~/c","name":"c","total_score":0.95,"tier":"ouro"},
	 {"ref":"~/d","name":"d","total_score":0.7,"tier":"prata"},
	 {"ref":"~/e","name":"e","total_score":0.40,"tier":"bronze"}],
	"entradas":[
	 {"ref":"~/a","nome":"a","tipo":"skill","hash":"h1","perfil":{"casos_de_uso":["tdd"],"etapas":["construir"],"naturezas":["software"],"origem":"inferido"}},
	 {"ref":"~/b","nome":"b","tipo":"script","hash":"h1"},
	 {"ref":"~/c","nome":"a","tipo":"skill","hash":"h1","perfil":{"casos_de_uso":["tdd"],"etapas":["construir"],"naturezas":["software"],"origem":"inferido"}},
	 {"ref":"~/d","nome":"d","tipo":"skill","hash":"h2","perfil":{"casos_de_uso":[],"etapas":["coletar"],"naturezas":["dados"],"origem":"declarado"}},
	 {"ref":"~/e","nome":"e","tipo":"skill","hash":"h3"}]}`
	if err := os.MkdirAll(a.vault, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(a.vault, "registry_matrix.json"), []byte(reg), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMatrizVer_RegistroManualAgrupaPorTipoEHashEFiltraPelaNota(t *testing.T) {
	a := novoAmbiente(t)
	registroManual(t, a)
	casos := []struct {
		args []string
		quer string
	}{
		{nil, "a,b,d,e"}, // a e c são cópias; b tem o mesmo hash, mas é outro tipo
		{[]string{"--tier", "prata"}, "d"},
		{[]string{"--tier", "OURO"}, "a"}, // maiúsculas valem
		{[]string{"--tipo", "SKILL"}, "a,d,e"},
		{[]string{"--min", "0.7"}, "a,d"}, // a nota mínima é inclusiva
		{[]string{"--min", "0.71"}, "a"},
		{[]string{"--natureza", "software"}, "a"}, // peça sem perfil nunca passa num filtro de termo
		{[]string{"--origem", "nenhuma"}, "b,e"},
		{[]string{"--tier", "bronze", "--origem", "nenhuma"}, "e"},
	}
	for _, c := range casos {
		if got := nomesDe(a.verJSON(nil, c.args...)); got != c.quer {
			t.Errorf("%v: veio %q, esperava %q", c.args, got, c.quer)
		}
	}
	if v := a.verJSON(nil); v.Entradas[0].Copias != 2 {
		t.Errorf("a e c são cópias: %d", v.Entradas[0].Copias)
	}
}

func TestMatrizVer_CabecalhoDaUltimaColuna(t *testing.T) {
	a := novoAmbiente(t)
	registroManual(t, a)
	_, out, _ := a.ver(nil)
	if !strings.Contains(out, "CÓPIAS") || strings.Contains(out, "CAMINHO") {
		t.Errorf("agrupado: coluna CÓPIAS:\n%s", out)
	}
	_, out, _ = a.ver(nil, "--copias")
	if !strings.Contains(out, "CAMINHO") || strings.Contains(out, "CÓPIAS") {
		t.Errorf("--copias: coluna CAMINHO:\n%s", out)
	}
}

func registroComPapeis(t *testing.T, a *ambiente) {
	t.Helper()
	reg := `{"generated_at":"2026-11-10","matrix_version":"2.0.0","total_pieces":4,
	"ratings":[
	 {"ref":"~/a","name":"peca-tdd","total_score":0.95,"tier":"ouro"},
	 {"ref":"~/b","name":"peca-audio","total_score":0.85,"tier":"prata"},
	 {"ref":"~/c","name":"peca-multi","total_score":0.90,"tier":"ouro"},
	 {"ref":"~/d","name":"peca-secops","total_score":0.70,"tier":"prata"}],
	"entradas":[
	 {"ref":"~/a","nome":"peca-tdd","tipo":"skill","hash":"h1","papel_primario":"engenheiro-tdd","perfil":{"casos_de_uso":["tdd"],"etapas":["construir"],"naturezas":["software"],"origem":"declarado"}},
	 {"ref":"~/b","nome":"peca-audio","tipo":"skill","hash":"h2","papel_primario":"animador-motion","perfil":{"casos_de_uso":["motion"],"etapas":["renderizar"],"naturezas":["midia"],"origem":"declarado"}},
	 {"ref":"~/c","nome":"peca-multi","tipo":"skill","hash":"h3","papel_primario":"engenheiro-backend","papeis_secundarios":["engenheiro-tdd"],"perfil":{"casos_de_uso":["api"],"etapas":["construir"],"naturezas":["software"],"origem":"declarado"}},
	 {"ref":"~/d","nome":"peca-secops","tipo":"skill","hash":"h4","papel_primario":"engenheiro-secops","perfil":{"casos_de_uso":["seguranca"],"etapas":["verificar"],"naturezas":["infra"],"origem":"declarado"}}]}
`
	if err := os.MkdirAll(a.vault, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(a.vault, "registry_matrix.json"), []byte(reg), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestMatrizVerFiltrosPapelGuilda(t *testing.T) {
	t.Run("filtro por papel", testMatrizVerFiltroPapel)
	t.Run("filtro por guilda", testMatrizVerFiltroGuilda)
	t.Run("erro papel e guilda desconhecidos", testMatrizVerErroDesconhecidos)
	t.Run("filtro por frontmatter no disco", testMatrizVerFiltroFrontmatterDisco)
}

func testMatrizVerFiltroPapel(t *testing.T) {
	a := novoAmbiente(t)
	registroComPapeis(t, a)
	casos := []struct {
		args []string
		quer string
	}{
		{[]string{"--papel", "engenheiro-tdd"}, "peca-tdd,peca-multi"},
		{[]string{"--papel", "animador-motion"}, "peca-audio"},
		{[]string{"--papel", "ENGENHEIRO-SECOPS"}, "peca-secops"},
	}
	for _, c := range casos {
		if got := nomesDe(a.verJSON(nil, c.args...)); got != c.quer {
			t.Errorf("%v: veio %q, esperava %q", c.args, got, c.quer)
		}
	}
}

func testMatrizVerFiltroGuilda(t *testing.T) {
	a := novoAmbiente(t)
	registroComPapeis(t, a)
	casos := []struct {
		args []string
		quer string
	}{
		{[]string{"--guilda", "audiovisual"}, "peca-audio"},
		{[]string{"--guilda", "construcao"}, "peca-tdd,peca-multi"},
		{[]string{"--guilda", "INFRA"}, "peca-secops"},
	}
	for _, c := range casos {
		if got := nomesDe(a.verJSON(nil, c.args...)); got != c.quer {
			t.Errorf("%v: veio %q, esperava %q", c.args, got, c.quer)
		}
	}
}

func testMatrizVerErroDesconhecidos(t *testing.T) {
	a := novoAmbiente(t)
	registroComPapeis(t, a)
	casos := []struct {
		args []string
		quer string
	}{
		{[]string{"--papel", "ninja-react"}, "catálogo"},
		{[]string{"--guilda", "guilda-inexistente"}, "catálogo"},
	}
	for _, c := range casos {
		code, out, errb := a.ver(nil, c.args...)
		if code != 1 || !strings.Contains(errb, c.quer) || out != "" {
			t.Errorf("%v: exit %d, stdout %q, stderr %q", c.args, code, out, errb)
		}
	}
}

func testMatrizVerFiltroFrontmatterDisco(t *testing.T) {
	a := novoAmbiente(t)
	skillPath := filepath.Join(a.home, "skills", "minha-skill", "SKILL.md")
	_ = os.MkdirAll(filepath.Dir(skillPath), 0o755)
	_ = os.WriteFile(skillPath, []byte("---\npapel: dissecador-web\npapeis: [engenheiro-tdd]\n---\n# Minha Skill\n"), 0o644)
	reg := `{"generated_at":"2026-11-10","matrix_version":"2.0.0","total_pieces":1,
	"ratings":[{"ref":"` + filepath.ToSlash(skillPath) + `","name":"minha-skill","total_score":0.90,"tier":"ouro"}],
	"entradas":[{"ref":"` + filepath.ToSlash(skillPath) + `","nome":"minha-skill","tipo":"skill","hash":"h1"}]}`
	_ = os.MkdirAll(a.vault, 0o755)
	_ = os.WriteFile(filepath.Join(a.vault, "registry_matrix.json"), []byte(reg), 0o644)
	if got := nomesDe(a.verJSON(nil, "--papel", "dissecador-web")); got != "minha-skill" {
		t.Fatalf("esperava minha-skill por papel primario no frontmatter, obteve %q", got)
	}
	if got := nomesDe(a.verJSON(nil, "--papel", "engenheiro-tdd")); got != "minha-skill" {
		t.Fatalf("esperava minha-skill por papel secundario no frontmatter, obteve %q", got)
	}
	if got := nomesDe(a.verJSON(nil, "--guilda", "reversa")); got != "minha-skill" {
		t.Fatalf("esperava minha-skill por guilda reversa, obteve %q", got)
	}
}
