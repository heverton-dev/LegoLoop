package redc_test

import (
	"strings"
	"testing"

	"legoloop/internal/redc"
)

func TestGenerateStructuralVideo_Success(t *testing.T) {
	tempDir := t.TempDir()
	dna := redc.ChannelDNA{
		CanalSlug:    "canal-referencia",
		TomDeVoz:     "Eletrizante e persuasivo",
		RitmoCortes:  "Cortes a cada 2.5s",
		DesignTokens: map[string]string{"cor_primaria": "#000000"},
		MotionTokens: map[string]string{"transicao": "snap_zoom"},
		TopGanchos:   []string{"Pare tudo que voce esta fazendo"},
	}

	req := redc.GenerationRequest{
		CanalSlug: "canal-referencia",
		Tema:      "Como construir agentes autônomos com LegoLoop",
		Marca: redc.BrandTokens{
			NomeMarca:     "MinhaEmpresa IA",
			CorPrimaria:   "#3B82F6",
			CorSecundaria: "#10B981",
			FonteTitulo:   "Inter Black",
			FonteCorpo:    "Inter Regular",
			TomPreferido:  "Executivo e altamente técnico",
		},
	}

	outPath, err := redc.GenerateStructuralVideo(tempDir, dna, req)
	if err != nil {
		t.Fatalf("esperava sucesso ao gerar video com brand swap, obteve: %v", err)
	}

	if outPath == "" {
		t.Fatalf("caminho de saida nao deve ser vazio")
	}

	// Ler conteúdo gerado
	raw, err := redc.ReadFile(outPath)
	if err != nil {
		t.Fatalf("falha ao ler output gerado: %v", err)
	}
	content := string(raw)

	if !strings.Contains(content, "MinhaEmpresa IA") {
		t.Errorf("roteiro deve conter o nome da marca do usuario")
	}
	if !strings.Contains(content, "#3B82F6") {
		t.Errorf("roteiro deve conter a cor primaria da marca do usuario")
	}
	if !strings.Contains(content, "Pare tudo que voce esta fazendo") {
		t.Errorf("roteiro deve preservar a estrutura de gancho comprovada do canal")
	}
}

func TestGenerateStructuralVideo_MissingTema(t *testing.T) {
	tempDir := t.TempDir()
	dna := redc.ChannelDNA{CanalSlug: "canal"}
	req := redc.GenerationRequest{CanalSlug: "canal", Tema: ""}

	_, err := redc.GenerateStructuralVideo(tempDir, dna, req)
	if err != redc.ErrTemaObrigatorio {
		t.Errorf("esperava ErrTemaObrigatorio, obteve: %v", err)
	}
}
