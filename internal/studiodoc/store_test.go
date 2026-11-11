package studiodoc_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"legoloop/internal/dag"
	"legoloop/internal/studiodoc"
)

func sampleValidLoop(id, name string) *studiodoc.Processo {
	return &studiodoc.Processo{
		ID:          id,
		Name:        name,
		Description: "Loop de teste",
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

	lp := sampleValidLoop("pesquisa-mercado", "Pesquisa de Mercado")
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

	// Verifica se gravou no disco como .loop.json
	expectedPath := filepath.Join(dir, "pesquisa-mercado.loop.json")
	if _, err := os.Stat(expectedPath); err != nil {
		t.Fatalf("arquivo .loop.json não foi criado em %s: %v", expectedPath, err)
	}
}

func TestStoreSaveInvalidDAG(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	// Loop com nó sem trigger e com dependência inexistente
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
		t.Fatal("esperava erro ao salvar loop com DAG inválido")
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
		lp := sampleValidLoop(badID, "Teste")
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

	l1 := sampleValidLoop("beta-loop", "Beta Loop")
	l2 := sampleValidLoop("alpha-loop", "Alpha Loop")

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
	if list[0].ID != "alpha-loop" || list[1].ID != "beta-loop" {
		t.Fatalf("ordenação inesperada: %+v", list)
	}
}

func TestStoreDuplicate(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore falhou: %v", err)
	}

	orig := sampleValidLoop("orig", "Original")
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

	lp := sampleValidLoop("apagar", "Para Apagar")
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

func TestStoreListIgnoresNonLoopFilesAndCorruptFiles(t *testing.T) {
	dir := t.TempDir()
	store, err := studiodoc.NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}

	// Arquivo comum não .loop.json
	_ = os.WriteFile(filepath.Join(dir, "README.md"), []byte("ignorar"), 0600)
	// Pasta interna
	_ = os.Mkdir(filepath.Join(dir, "subfolder"), 0700)
	// Arquivo .loop.json corrompido
	_ = os.WriteFile(filepath.Join(dir, "corrompido.loop.json"), []byte("{invalid json"), 0600)

	lp := sampleValidLoop("valido", "Válido")
	if err := store.Save(lp); err != nil {
		t.Fatal(err)
	}

	list, err := store.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "valido" {
		t.Fatalf("esperava apenas 1 loop válido, obteve %+v", list)
	}
}
