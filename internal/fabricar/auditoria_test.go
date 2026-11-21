package fabricar_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"legoloop/internal/caminhos"
	"legoloop/internal/cli"
	"legoloop/internal/fabricar"
	"legoloop/internal/gate"
	"legoloop/internal/triagem"
)

// dataFixa dá nomes ISO estáveis aos artefatos do scaffold.
var dataFixa = time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)

// stubExit devolve o critério binário fixo para provar o contrato sem tocar o sistema real.
func stubExit(comando string, exit int) fabricar.TesteFunc {
	return func(fabricar.Alvo) fabricar.ResultadoTeste {
		return fabricar.ResultadoTeste{Comando: comando, Saida: "evidencia do stub", Exit: exit}
	}
}

// escrever grava o arquivo e falha o teste se não conseguir.
func escrever(t *testing.T, caminho, conteudo string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(caminho), 0o755); err != nil {
		t.Fatalf("criar %s: %v", filepath.Dir(caminho), err)
	}
	if err := os.WriteFile(caminho, []byte(conteudo), 0o644); err != nil {
		t.Fatalf("gravar %s: %v", caminho, err)
	}
}

// TestAlvoArbitrario prova que a auditoria universal é agnóstica: audita peça, função, documento,
// binário e verbo com o mesmo critério binário, e reprova com ticket e refatoração quando falta critério.
func TestAlvoArbitrario(t *testing.T) {
	dir := t.TempDir()
	escrever(t, filepath.Join(dir, "script.py"), "print('auditoria')\n")
	escrever(t, filepath.Join(dir, "util.go"), "package util\n\nfunc Analisar() int { return 1 }\n")
	escrever(t, filepath.Join(dir, "notas.md"), "# Notas\n\nConteudo de estudo da auditoria universal.\n")
	escrever(t, filepath.Join(dir, "dados.bin"), "\x01\x02\x03\x04")
	verboBom := filepath.Join(dir, "internal", "auditor")
	escrever(t, filepath.Join(verboBom, "auditor.go"), "package auditor\n")
	escrever(t, filepath.Join(verboBom, "auditor_test.go"), "package auditor\n\nfunc TestBase(t *testing.T) {}\n")

	casos := []struct {
		nome string
		alvo fabricar.Alvo
	}{
		{"peca", fabricar.Alvo{Tipo: fabricar.AlvoPeca, Nome: "script", Caminho: filepath.Join(dir, "script.py")}},
		{"funcao", fabricar.Alvo{Tipo: fabricar.AlvoFuncao, Nome: "Analisar", Caminho: filepath.Join(dir, "util.go")}},
		{"documento", fabricar.Alvo{Tipo: fabricar.AlvoDocumento, Nome: "notas", Caminho: filepath.Join(dir, "notas.md")}},
		{"binario", fabricar.Alvo{Tipo: fabricar.AlvoBinario, Nome: "dados", Caminho: filepath.Join(dir, "dados.bin")}},
		{"verbo", fabricar.Alvo{Tipo: fabricar.AlvoVerbo, Nome: "auditor", Caminho: verboBom}},
	}
	for _, c := range casos {
		rel, code, err := fabricar.Auditar(c.alvo, stubExit("go test -count=1", 0))
		if err != nil {
			t.Fatalf("%s: auditar: %v", c.nome, err)
		}
		if code != 0 {
			t.Errorf("%s: exit %d, esperado 0", c.nome, code)
		}
		if rel.Veredito != fabricar.VereditoAprovado {
			t.Errorf("%s: veredito %q, esperado %q", c.nome, rel.Veredito, fabricar.VereditoAprovado)
		}
		if len(rel.Dissecao) == 0 {
			t.Errorf("%s: disseção vazia", c.nome)
		}
		if len(rel.Evidencias) == 0 || rel.Evidencias[0].Comando == "" {
			t.Errorf("%s: sem evidência mecânica registrada", c.nome)
		}
		texto := fabricar.Renderizar(rel)
		for _, secao := range []string{"## Disseção", "## Lacunas", "## Tickets", "## Refatoração", "## Veredito"} {
			if !strings.Contains(texto, secao) {
				t.Errorf("%s: render sem a seção %q", c.nome, secao)
			}
		}
	}

	// verbo sem teste unitário e sem critério mecânico reprova, com ticket e refatoração.
	falho := filepath.Join(dir, "internal", "falho")
	escrever(t, filepath.Join(falho, "falho.go"), "package falho\n")
	rel, code, err := fabricar.Auditar(
		fabricar.Alvo{Tipo: fabricar.AlvoVerbo, Nome: "falho", Caminho: falho},
		stubExit("go test -count=1", 1),
	)
	if err != nil {
		t.Fatalf("auditar verbo falho: %v", err)
	}
	if code != 1 {
		t.Errorf("exit %d do verbo reprovado, esperado 1", code)
	}
	if rel.Veredito != fabricar.VereditoReprovado {
		t.Errorf("veredito %q do verbo reprovado", rel.Veredito)
	}
	if len(rel.Tickets) == 0 || rel.Tickets[0].ID != "TICK-AUD-01" {
		t.Errorf("tickets %v do verbo reprovado, esperado TICK-AUD-01", rel.Tickets)
	}
	if strings.TrimSpace(rel.Refatoracao) == "" {
		t.Error("refatoração vazia no verbo reprovado")
	}

	// documento sem conteúdo reprova pelo critério interno.
	vazio := filepath.Join(dir, "vazio.md")
	escrever(t, vazio, "   \n")
	rel, code, err = fabricar.Auditar(
		fabricar.Alvo{Tipo: fabricar.AlvoDocumento, Nome: "vazio", Caminho: vazio},
		stubExit("criterio interno", 0),
	)
	if err != nil {
		t.Fatalf("auditar documento vazio: %v", err)
	}
	if code != 1 || rel.Veredito != fabricar.VereditoReprovado {
		t.Errorf("documento vazio: exit %d veredito %q, esperado 1 e %q", code, rel.Veredito, fabricar.VereditoReprovado)
	}
}

// TestScaffoldAuditoriaAtomica prova que o Processo auditoria-atomica nasce com o nascimento 5/5
// gerado pelo scaffold: manifestos válidos, skill, Quinteto ISO e DAG que passa no legoloop validate.
func TestScaffoldAuditoriaAtomica(t *testing.T) {
	raiz := t.TempDir()
	gravados, err := fabricar.MaterializarAuditoriaAtomica(fabricar.Roots{Vault: raiz, Docs: raiz}, dataFixa)
	if err != nil {
		t.Fatalf("materializar: %v", err)
	}
	if len(gravados) == 0 {
		t.Fatal("nenhum artefato gravado")
	}
	dir := filepath.Join(raiz, caminhos.DirProcessos, "auditoria-atomica")
	for _, manifesto := range []string{
		caminhos.ArquivoDAG,
		"auditoria-atomica" + caminhos.SufixoStudio,
		triagem.ArquivoTriagem,
		triagem.ArquivoMemoria,
		"anticorpos.json",
	} {
		if _, err := os.Stat(filepath.Join(dir, manifesto)); err != nil {
			t.Errorf("manifesto ausente %s: %v", manifesto, err)
		}
	}
	if _, err := os.Stat(filepath.Join(raiz, "skills", "legoloop-auditoria-atomica", "SKILL.md")); err != nil {
		t.Errorf("skill ausente: %v", err)
	}
	for _, iso := range []string{
		filepath.Join("docs", "intents", "2026-11-17_INTENT_AUDITORIA_ATOMICA.md"),
		filepath.Join("docs", "specs", "2026-11-17_SPEC_AUDITORIA_ATOMICA.md"),
		filepath.Join("docs", "dags", "2026-11-17_DAG_AUDITORIA_ATOMICA.json"),
	} {
		if _, err := os.Stat(filepath.Join(raiz, iso)); err != nil {
			t.Errorf("artefato ISO ausente %s: %v", iso, err)
		}
	}
	for _, par := range [][2]string{{"plans", ".md"}, {"routes", ".json"}} {
		padrao := filepath.Join(raiz, "docs", par[0], "*AUDITORIA_ATOMICA*"+par[1])
		if achados, _ := filepath.Glob(padrao); len(achados) == 0 {
			t.Errorf("nenhum artefato ISO em docs/%s (%s)", par[0], padrao)
		}
	}
	conferirSchemas(t, dir)
	validatePassa(t, filepath.Join(dir, caminhos.ArquivoDAG))
	var d struct {
		Nodes []json.RawMessage `json:"nodes"`
	}
	bruto, err := os.ReadFile(filepath.Join(dir, caminhos.ArquivoDAG))
	if err != nil {
		t.Fatalf("ler DAG: %v", err)
	}
	if err := json.Unmarshal(bruto, &d); err != nil {
		t.Fatalf("DAG não é JSON: %v", err)
	}
	if len(d.Nodes) != 5 {
		t.Errorf("DAG com %d nós, esperado 5", len(d.Nodes))
	}
}

// conferirSchemas valida os manifestos do nascimento contra os schemas oficiais.
func conferirSchemas(t *testing.T, dir string) {
	t.Helper()
	for manifesto, esquema := range map[string]string{
		triagem.ArquivoTriagem: "triagem.v1.json",
		triagem.ArquivoMemoria: "memoria_execucao.v1.json",
		"anticorpos.json":      "anticorpos.v1.json",
	} {
		b, err := os.ReadFile(filepath.Join(dir, manifesto))
		if err != nil {
			t.Errorf("manifesto %s: %v", manifesto, err)
			continue
		}
		schema, err := os.ReadFile(filepath.Join("..", "..", "schemas", esquema))
		if err != nil {
			t.Fatalf("schema %s: %v", esquema, err)
		}
		if !gate.ValidateJSONSchema(schema, b) {
			t.Errorf("%s não valida contra %s", manifesto, esquema)
		}
	}
}

// validatePassa roda o legoloop validate do CLI real sobre o DAG do Processo.
func validatePassa(t *testing.T, dagPath string) {
	t.Helper()
	var out, errb bytes.Buffer
	if code := cli.Run(context.Background(), []string{"validate", dagPath}, cli.Deps{Stdout: &out, Stderr: &errb}); code != 0 {
		t.Errorf("legoloop validate falhou em %s (%d): %s", dagPath, code, errb.String())
	}
}

// TestAlvoInexistenteExitDois prova que o alvo inexistente devolve exit 2 com erro claro e zero artefato órfão.
func TestAlvoInexistenteExitDois(t *testing.T) {
	dir := t.TempDir()
	dirSaida := filepath.Join(dir, "saida")
	if err := os.MkdirAll(dirSaida, 0o755); err != nil {
		t.Fatalf("criar %s: %v", dirSaida, err)
	}
	saida := filepath.Join(dirSaida, "relatorio.md")
	alvo := fabricar.Alvo{Tipo: fabricar.AlvoVerbo, Nome: "fantasma", Caminho: filepath.Join(dir, "internal", "fantasma")}
	rel, code, err := fabricar.AuditarEGravar(alvo, stubExit("go test -count=1", 0), saida)
	if !errors.Is(err, fabricar.ErrAlvoInexistente) {
		t.Errorf("erro %v, esperado ErrAlvoInexistente", err)
	}
	if code != 2 {
		t.Errorf("exit %d do alvo inexistente, esperado 2", code)
	}
	if rel.Veredito != "" {
		t.Errorf("veredito %q para alvo inexistente, esperado vazio", rel.Veredito)
	}
	if _, err := os.Stat(saida); !os.IsNotExist(err) {
		t.Errorf("artefato órfão gravado em %s", saida)
	}
	entradas, err := os.ReadDir(dirSaida)
	if err != nil {
		t.Fatalf("ler %s: %v", dirSaida, err)
	}
	if len(entradas) != 0 {
		t.Errorf("diretório de saída com %d artefatos, esperado 0", len(entradas))
	}
}
