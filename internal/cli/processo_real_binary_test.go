package cli

import (
	"legoloop/internal/caminhos"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/processo"
)

// fabricarReal roda o binário de verdade (sem --dry-run) num HOME, cofre e projeto descartáveis.
func fabricarReal(t *testing.T, cer *cerebroPorPapel, extra ...string) (binRun, string, string) {
	t.Helper()
	bin := buildBinary(t)
	root := t.TempDir()
	vault, proj := filepath.Join(root, "vault"), filepath.Join(root, "proj")
	matrizDeTeste(t, vault)
	args := append([]string{"fabricar", "extrair audio", "--goal", "construir e testar um extrator", "--vault", vault,
		"--docs", proj, "--yes", "--sem-entrevista", "--no-link"}, extra...)
	return runBinaryEm(t, bin, cer.servidor(t), root, proj, args...), vault, proj
}

// TestProcessoRealBinarioGravaOProcessoEAsSeisVisoesCoerentes prova pelo binário que a Camada 4 nasce do
// Processo: o processo.json fica no cofre e valida, e DAG, PLAN e SPEC listam as mesmas etapas, na mesma
// ordem, com gates que verificam conteúdo. A SPEC não é pedida ao cérebro.
func TestProcessoRealBinarioGravaOProcessoEAsSeisVisoesCoerentes(t *testing.T) {
	cer := &cerebroPorPapel{}
	r, vault, proj := fabricarReal(t, cer)
	if r.code != 0 {
		t.Fatalf("exit %d\n%s\n%s", r.code, r.out, r.errb)
	}
	p, err := processo.Ler(processo.CaminhoNoCofre(vault, "legoloop-extrair-audio"))
	if err != nil {
		t.Fatalf("o processo.json deve existir e validar: %v", err)
	}
	if p.Natureza != "software" || p.Versao != "1.0.0" || len(p.Etapas) == 0 {
		t.Fatalf("processo gravado errado: %+v", p)
	}
	var ids []string
	for _, e := range p.Etapas {
		ids = append(ids, e.ID)
		if e.Peca == "" || !(strings.Contains(e.Gate, "word_count") || strings.Contains(e.Gate, "line_count") || strings.Contains(e.Gate, "schema_valid")) {
			t.Errorf("etapa %s sem peça ou com gate só de existência: %+v", e.ID, e)
		}
	}
	assertDAGComAsEtapas(t, filepath.Join(vault, caminhos.DirLoops, "legoloop-extrair-audio", "pipeline.dag.json"), ids)
	plan := firstGlob(t, filepath.Join(proj, "docs", "plans", "*_PLAN_EXTRAIR_AUDIO.md"))
	spec := firstGlob(t, filepath.Join(proj, "docs", "specs", "*_SPEC_EXTRAIR_AUDIO.md"))
	for _, id := range ids {
		if !strings.Contains(plan, "`"+id+"`") || !strings.Contains(spec, id) {
			t.Errorf("a etapa %s deve aparecer no PLAN e na SPEC", id)
		}
	}
	if strings.Contains(plan, "Aplicar a skill") {
		t.Error("o PLAN não deve usar o texto genérico da v1")
	}
	if cer.spec.Load() != 0 {
		t.Errorf("a SPEC nasce do Processo, não do cérebro: %d pedidos", cer.spec.Load())
	}
}

func assertDAGComAsEtapas(t *testing.T, path string, ids []string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	d, err := dag.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	var nos []string
	for _, n := range d.Nodes {
		nos = append(nos, n.ID)
	}
	if !slices.Equal(nos, ids) {
		t.Errorf("o DAG deve ter as etapas do processo, na ordem: %v != %v", nos, ids)
	}
}

// TestProcessoRealBinarioLacunaParaSemGravar: etapa sem peça apta não bloqueia o fabricar;
// segue com exit 0 e emite sugestões externas curadoria-externa-*.
func TestProcessoRealBinarioLacunaParaSemGravar(t *testing.T) {
	cer := &cerebroPorPapel{needsJSON: `{"natureza": "software", "etapas": ["construir", "publicar"], "casos": []}`}
	r, _, _ := fabricarReal(t, cer)
	if r.code != 0 || !strings.Contains(r.errb, "curadoria-externa-construir") || !strings.Contains(r.errb, "curadoria-externa-publicar") {
		t.Fatalf("esperava exit 0 com sugestões curadoria-externa-*; exit %d\n%s\n%s", r.code, r.out, r.errb)
	}
}

// TestProcessoRealBinarioCerebroCaidoParaSemGravar: o cérebro cai ao dar as necessidades; o fabricar para.
func TestProcessoRealBinarioCerebroCaidoParaSemGravar(t *testing.T) {
	cer := &cerebroPorPapel{needsFalha: true}
	r, vault, proj := fabricarReal(t, cer)
	if r.code != 1 || cer.needs.Load() == 0 {
		t.Fatalf("esperava exit 1 depois de pedir as necessidades; exit %d (pedidos %d)\n%s\n%s", r.code, cer.needs.Load(), r.out, r.errb)
	}
	assertSoAMatriz(t, vault, proj)
}

func assertSoAMatriz(t *testing.T, vault, proj string) {
	t.Helper()
	got := append(filesUnder(vault), filesUnder(proj)...)
	if len(got) != 1 || filepath.Base(got[0]) != "registry_matrix.json" {
		t.Errorf("nada além da matriz deveria existir: %v", got)
	}
}
