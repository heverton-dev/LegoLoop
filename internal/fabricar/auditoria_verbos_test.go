package fabricar_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"legoloop/internal/fabricar"
)

// stubPorVerbo devolve o critério binário simulado com exit fixo por verbo.
func stubPorVerbo(exits map[string]int) fabricar.TesteFunc {
	return func(a fabricar.Alvo) fabricar.ResultadoTeste {
		return fabricar.ResultadoTeste{
			Comando: "go test -count=1 ./internal/" + a.Nome + "/...",
			Saida:   "simulado para " + a.Nome,
			Exit:    exits[a.Nome],
		}
	}
}

// arvoreVerbos cria os seis verbos com fonte e teste unitário sob raiz/internal.
func arvoreVerbos(t *testing.T, raiz string) {
	t.Helper()
	for _, v := range fabricar.VerbosAuditados {
		dir := filepath.Join(raiz, "internal", v)
		escrever(t, filepath.Join(dir, v+".go"), "package "+v+"\n")
		escrever(t, filepath.Join(dir, v+"_test.go"), "package "+v+"\n\nfunc TestBase(t *testing.T) {}\n")
	}
}

// TestAuditoriaSeisVerbosCriterioBinario prova que cada verbo audita com critério
// binário (exit 0/1) e registra relatório atômico com lacunas, tickets, veredito e evidência.
func TestAuditoriaSeisVerbosCriterioBinario(t *testing.T) {
	raiz := t.TempDir()
	arvoreVerbos(t, raiz)
	exits := map[string]int{"matriz": 0, "listar": 0, "refinar": 0, "remover": 0, "compor": 1, "fechar": 0}
	dir := t.TempDir()
	rels, codigos, err := fabricar.AuditarVerbosEGravar(raiz, stubPorVerbo(exits), dir)
	if err != nil {
		t.Fatal(err)
	}
	esperado := []int{0, 0, 0, 0, 1, 0}
	if !reflect.DeepEqual(codigos, esperado) {
		t.Fatalf("códigos = %v, esperava %v", codigos, esperado)
	}
	if len(rels) != 6 {
		t.Fatalf("esperava 6 relatórios, obteve %d", len(rels))
	}
	for i, v := range fabricar.VerbosAuditados {
		b, err := os.ReadFile(filepath.Join(dir, v+".md"))
		if err != nil {
			t.Fatalf("relatório atômico do verbo %s não foi gravado: %v", v, err)
		}
		conteudo := string(b)
		if !strings.Contains(conteudo, "Item: "+v) {
			t.Errorf("relatório de %s não nomeia o alvo: %q", v, conteudo[:min(200, len(conteudo))])
		}
		for _, secao := range []string{"## Disseção", "## Lacunas", "## Tickets", "## Refatoração", "## Veredito", "### Evidências"} {
			if !strings.Contains(conteudo, secao) {
				t.Errorf("relatório de %s sem seção %s", v, secao)
			}
		}
		esperava := fabricar.VereditoAprovado
		if exits[v] != 0 {
			esperava = fabricar.VereditoReprovado
		}
		if !strings.Contains(conteudo, "**"+esperava+"**") {
			t.Errorf("relatório de %s deveria terminar em %s", v, esperava)
		}
		if !strings.Contains(conteudo, "go test -count=1 ./internal/"+v+"/...") {
			t.Errorf("relatório de %s não registra a evidência do critério binário", v)
		}
		if exits[v] != 0 && !strings.Contains(conteudo, "TICK-AUD-") {
			t.Errorf("verbo %s reprovado sem ticket de lacuna no relatório", v)
		}
		if got := fabricar.CodigoDe(rels[i]); got != codigos[i] {
			t.Errorf("CodigoDe(%s) = %d, coerência exigia %d", v, got, codigos[i])
		}
	}
}

// TestAuditoriaVerbosReaisCobremOsSeis prova que os seis verbos reais do repositório
// (incluindo compor e fechar, intocados hoje) cobrem o critério estrutural e saem
// aprovados quando o critério binário mecânico devolve exit 0.
func TestAuditoriaVerbosReaisCobremOsSeis(t *testing.T) {
	raiz := filepath.Join("..", "..")
	exits := map[string]int{}
	for _, v := range fabricar.VerbosAuditados {
		exits[v] = 0
	}
	rels, codigos, err := fabricar.AuditarVerbosEGravar(raiz, stubPorVerbo(exits), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range fabricar.VerbosAuditados {
		if rels[i].Alvo.Nome != v {
			t.Fatalf("posição %d trouxe %s, esperava %s", i, rels[i].Alvo.Nome, v)
		}
		if codigos[i] != fabricar.ExitAprovado {
			t.Errorf("verbo real %s reprovado no critério estrutural: lacunas %v", v, rels[i].Lacunas)
		}
		if len(rels[i].Evidencias) != 1 || rels[i].Evidencias[0].Exit != 0 {
			t.Errorf("verbo %s sem evidência mecânica aprovada: %+v", v, rels[i].Evidencias)
		}
	}
}
