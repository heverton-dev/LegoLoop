package legopress_test

import (
	"testing"

	"legoloop/internal/legopress"
)

func TestPartitionHierarchy_Valid(t *testing.T) {
	rawTopics := []legopress.RawTopic{
		{Title: "A Crise dos Chatbots Monoliticos", ContentType: "concept"},
		{Title: "Os Cinco Verbos Canonicos", ContentType: "process"},
		{Title: "O Motor de Execucao em Go", ContentType: "code"},
		{Title: "Receitas de Esteiras Reais", ContentType: "recipe"},
	}

	partitioner := legopress.NewPartitioner()
	parts, err := partitioner.Partition(rawTopics)
	if err != nil {
		t.Fatalf("esperava sucesso no particionamento, obteve erro: %v", err)
	}

	if len(parts) != 4 {
		t.Fatalf("esperava 4 partes particionadas, obteve: %d", len(parts))
	}

	// Verificar classificacao hierarquica
	if parts[0].Level != 1 { // Macro
		t.Errorf("esperava conceito em nivel 1 (Macro), obteve: %d", parts[0].Level)
	}
	if parts[1].Level != 2 { // Meso
		t.Errorf("esperava processo em nivel 2 (Meso), obteve: %d", parts[1].Level)
	}
	if parts[2].Level != 3 { // Micro
		t.Errorf("esperava codigo em nivel 3 (Micro), obteve: %d", parts[2].Level)
	}
}

func TestPartitionHierarchy_Empty(t *testing.T) {
	partitioner := legopress.NewPartitioner()
	_, err := partitioner.Partition([]legopress.RawTopic{})
	if err == nil {
		t.Fatal("esperava erro ao particionar topicos vazios, obteve nil")
	}
}
