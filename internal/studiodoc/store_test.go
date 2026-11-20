package studiodoc_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"legoloop/internal/caminhos"
	"legoloop/internal/dag"
	"legoloop/internal/paths"
	"legoloop/internal/studiodoc"
)

func sampleValidProcesso(id, name string) *studiodoc.Processo {
	return &studiodoc.Processo{
		ID:          id,
		Name:        name,
		Description: "Processo de teste",
		Version:     "2.0.0",
		Trigger:     &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:   "n1",
				Type: dag.TypeSkill,
				Ref:  "interview",
				Routing: dag.RoutingConfig{
					Mode:  dag.ModeManual,
					Model: "llama3.1:8b",
				},
				Output: dag.HandoffContract{
					Artifact: "BRIEFING.md",
					Gate:     "file_exists('BRIEFING.md')",
				},
			},
		},
	}
}

func TestStoreSaveAndGet(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	lp := sampleValidProcesso("pesquisa-mercado", "Pesquisa de Mercado")
	if err := store.Save(lp); err != nil {
		t.Fatalf("Save falhou: %v", err)
	}

	got, err := store.Get("pesquisa-mercado")
	if err != nil {
		t.Fatalf("Get falhou: %v", err)
	}
	if got.ID != "pesquisa-mercado" || got.Name != "Pesquisa de Mercado" {
		t.Fatalf("dados divergentes: got %+v", got)
	}
	if len(got.Nodes) != 1 || got.Nodes[0].Output.Artifact != "BRIEFING.md" {
		t.Fatalf("nós divergentes: got %+v", got.Nodes)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || !strings.HasPrefix(entries[0].Name(), "pesquisa-mercado.") {
		t.Fatalf("arquivo do Processo não foi criado em %s: %v", dir, entries)
	}
}

func TestStoreSaveInvalidDAG(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	bad := &studiodoc.Processo{
		ID:      "invalido",
		Name:    "Inválido",
		Version: "2.0.0",
		Trigger: &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{
			{
				ID:        "n1",
				Type:      dag.TypeSkill,
				DependsOn: []string{"nao_existe"},
			},
		},
	}

	if err := store.Save(bad); err == nil {
		t.Fatal("esperava erro ao salvar Processo com DAG inválido")
	}
}

func TestStoreSaveInvalidID(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	invalidIDs := []string{
		"",
		"../perigo",
		"..\\perigo",
		"sub/dir",
		"C:perigo",
		"nome com espacos",
		"teste$bad",
	}

	for _, badID := range invalidIDs {
		lp := sampleValidProcesso(badID, "Teste")
		if err := store.Save(lp); err == nil {
			t.Errorf("esperava erro ao salvar com ID inválido %q", badID)
		}
	}
}

func TestStoreList(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	l1 := sampleValidProcesso("beta-proc", "Beta Processo")
	l2 := sampleValidProcesso("alpha-proc", "Alpha Processo")

	if err := store.Save(l1); err != nil {
		t.Fatal(err)
	}
	if err := store.Save(l2); err != nil {
		t.Fatal(err)
	}

	list, err := store.List()
	if err != nil {
		t.Fatalf("List falhou: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("esperava 2 itens, obteve %d", len(list))
	}
	// Deve ordenar por Name alfabético
	if list[0].ID != "alpha-proc" || list[1].ID != "beta-proc" {
		t.Fatalf("ordenação inesperada: %+v", list)
	}
}

func TestStoreDuplicate(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	orig := sampleValidProcesso("orig", "Original")
	if err := store.Save(orig); err != nil {
		t.Fatal(err)
	}

	dup, err := store.Duplicate("orig", "copia", "Cópia")
	if err != nil {
		t.Fatalf("Duplicate falhou: %v", err)
	}
	if dup.ID != "copia" || dup.Name != "Cópia" {
		t.Fatalf("dup divergente: %+v", dup)
	}

	// Ambas devem existir no store
	gotOrig, err := store.Get("orig")
	if err != nil || gotOrig.Name != "Original" {
		t.Fatalf("orig corrompido: %v", err)
	}
	gotDup, err := store.Get("copia")
	if err != nil || gotDup.Name != "Cópia" {
		t.Fatalf("dup não encontrado: %v", err)
	}
}

func TestStoreRemove(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	lp := sampleValidProcesso("apagar", "Para Apagar")
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}

	if err := store.Remove("apagar"); err != nil {
		t.Fatalf("Remove falhou: %v", err)
	}

	_, err = store.Get("apagar")
	if !errors.Is(err, studiodoc.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound após remover, obteve %v", err)
	}

	list, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("esperava lista vazia, obteve %d", len(list))
	}

	// Remove de ID inexistente dá ErrNotFound
	if err := store.Remove("inexistente"); !errors.Is(err, studiodoc.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound, obteve %v", err)
	}
	// Remove com ID inválido dá ErrInvalidID
	if err := store.Remove("../bad"); !errors.Is(err, studiodoc.ErrInvalidID) {
		t.Fatalf("esperava ErrInvalidID, obteve %v", err)
	}
	// Get com ID inválido dá ErrInvalidID
	if _, err := store.Get("../bad"); !errors.Is(err, studiodoc.ErrInvalidID) {
		t.Fatalf("esperava ErrInvalidID, obteve %v", err)
	}
	// Duplicate inexistente dá erro
	if _, err := store.Duplicate("inexistente", "novo", "Novo"); !errors.Is(err, studiodoc.ErrNotFound) {
		t.Fatalf("esperava ErrNotFound, obteve %v", err)
	}
	// Duplicate com novo ID inválido dá erro
	if _, err := store.Duplicate("apagar", "../invalido", "Novo"); err == nil {
		t.Fatal("esperava erro ao duplicar com novo ID inválido")
	}
}

func TestStoreListIgnoresForeignFilesAndCorruptFiles(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte("ignorar"), 0600)
	_ = os.Mkdir(filepath.Join(dir, "subfolder"), 0700)
	_ = os.WriteFile(filepath.Join(dir, "corrompido"+caminhos.SufixoStudio), []byte("{invalid json"), 0600)

	lp := sampleValidProcesso("valido", "Válido")
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}

	list, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "valido" {
		t.Fatalf("esperava apenas 1 Processo válido, obteve %+v", list)
	}
}

func TestStoreLegacySuffixCompatibility(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	raw := []byte(`{"id":"legado","name":"Legado","version":"2.0.0","trigger":{"type":"manual"},"nodes":[{"id":"n1","type":"skill","name":"N1"}]}`)
	if err := os.WriteFile(filepath.Join(dir, "legado"+paths.SufixoAntigo), raw, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("legado")
	if err != nil {
		t.Fatalf("Get legado falhou: %v", err)
	}
	if got.ID != "legado" || got.Name != "Legado" {
		t.Fatalf("processo inesperado: %+v", got)
	}
	list, err := store.List()
	if err != nil || len(list) != 1 || list[0].ID != "legado" {
		t.Fatalf("List legado falhou: %+v, err=%v", list, err)
	}
}

func makeProcessoV2Fixture() *studiodoc.Processo {
	return &studiodoc.Processo{
		ID: "proc-v2", Name: "Processo V2", Description: "Descricao completa",
		Version: "1.0.0", SchemaVersion: 2, Date: "2026-11-14",
		ExecutionStrategy:  map[string]any{"supervision_level": "epic"},
		UpstreamSignatures: map[string]any{"route": map[string]any{"sha256": "deadbeef"}},
		Trigger:            &dag.Trigger{Type: dag.TriggerManual},
		Nodes: []dag.Node{{
			ID: "n1", Name: "No 1", Type: dag.TypeSkill, Nature: "micro", Role: "engenheiro-backend", Ref: "skills/exemplo",
			Routing:          dag.RoutingConfig{Mode: dag.ModeSuggested, Harness: "9router", Via: "9router", Model: "t-codigo-128k", Effort: "high"},
			Isolation:        map[string]any{"write_boundary": "strict"},
			Limits:           map[string]any{"max_output_tokens": 300},
			RollbackStrategy: "abort_processo", Tooling: map[string]any{"runtime": "legoloop"},
			Verification: map[string]any{"gate": "exit_code == 0"}, SpecRef: "AC-1", Status: "PENDING",
			ContextIn: []string{"in.txt"}, ContextOut: []string{"out.txt"},
			Output: dag.HandoffContract{Artifact: "out.txt", Gate: "file_exists('out.txt')"},
		}},
	}
}

func TestStorePreservaMetadadosV2(t *testing.T) {
	store, err := studiodoc.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(makeProcessoV2Fixture()); err != nil {
		t.Fatalf("Save v2 falhou: %v", err)
	}
	got, err := store.Get("proc-v2")
	if err != nil {
		t.Fatalf("Get v2 falhou: %v", err)
	}
	if got.SchemaVersion != 2 || got.Date != "2026-11-14" || got.Description != "Descricao completa" {
		t.Fatalf("metadados basicos perdidos: %+v", got)
	}
	if got.ExecutionStrategy["supervision_level"] != "epic" || got.UpstreamSignatures["route"] == nil {
		t.Fatalf("metadados compostos perdidos: exec=%+v sig=%+v", got.ExecutionStrategy, got.UpstreamSignatures)
	}
}

func TestStoreSubpastaEDAG(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	subDAG := filepath.Join(dir, "proc-sub")
	if err := os.MkdirAll(subDAG, 0o700); err != nil {
		t.Fatal(err)
	}
	dagRaw := []byte(`{"name":"Processo DAG","nodes":[{"id":"a","type":"skill","name":"A"}]}`)
	if err := os.WriteFile(filepath.Join(subDAG, caminhos.ArquivoDAG), dagRaw, 0o600); err != nil {
		t.Fatal(err)
	}
	got, err := store.Get("proc-sub")
	if err != nil {
		t.Fatalf("Get %s em subpasta falhou: %v", caminhos.ArquivoDAG, err)
	}
	if got.ID != "proc-sub" || got.Name != "Processo DAG" || len(got.Nodes) != 1 {
		t.Fatalf("processo de subpasta divergente: %+v", got)
	}
	list, err := store.List()
	if err != nil || len(list) != 1 || list[0].ID != "proc-sub" {
		t.Fatalf("List subpasta falhou: %+v err=%v", list, err)
	}
}
