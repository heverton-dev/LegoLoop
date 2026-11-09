package redc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GenerateStructuralVideo gera um roteiro estrutural e especificações de vídeo com Brand Swap.
func GenerateStructuralVideo(rootDir string, dna ChannelDNA, req GenerationRequest) (string, error) {
	if strings.TrimSpace(req.Tema) == "" {
		return "", ErrTemaObrigatorio
	}
	outDir := filepath.Join(rootDir, "docs", "channels", dna.CanalSlug, "outputs")
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return "", fmt.Errorf("redc: falha ao criar pasta de outputs: %w", err)
	}

	slug := slugify(req.Tema)
	filename := fmt.Sprintf("roteiro_%s_%d.md", slug, time.Now().Unix())
	outPath := filepath.Join(outDir, filename)

	content := formatStructuralVideoContent(dna, req)
	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("redc: erro ao gravar roteiro: %w", err)
	}
	return outPath, nil
}

func formatStructuralVideoContent(dna ChannelDNA, req GenerationRequest) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# 🎬 Roteiro Estrutural & Brand Swap\n\n"))
	sb.WriteString(fmt.Sprintf("* **Canal Clonado (Referência):** %s\n", dna.CanalSlug))
	sb.WriteString(fmt.Sprintf("* **Tema do Vídeo:** %s\n", req.Tema))
	sb.WriteString(fmt.Sprintf("* **Marca Autoral:** %s\n", req.Marca.NomeMarca))
	sb.WriteString(fmt.Sprintf("* **Cores da Marca:** Primária %s | Secundária %s\n", req.Marca.CorPrimaria, req.Marca.CorSecundaria))
	sb.WriteString(fmt.Sprintf("* **Tipografia Autoral:** Títulos: %s | Corpo: %s\n", req.Marca.FonteTitulo, req.Marca.FonteCorpo))
	sb.WriteString(fmt.Sprintf("* **Tom de Voz Solicitado:** %s\n\n", req.Marca.TomPreferido))

	sb.WriteString("## ⏱️ Estrutura Viral Preservada\n\n")
	sb.WriteString(fmt.Sprintf("1. **[00:00 - 00:03] Gancho Viral:** \"%s: %s\"\n", getPrimaryHook(dna), req.Tema))
	sb.WriteString(fmt.Sprintf("2. **[00:03 - 00:15] Retenção & Problema:** Cadência dinâmica (%s)\n", dna.RitmoCortes))
	sb.WriteString("3. **[00:15 - 00:45] Revelação da Metodologia:** Exibição com cartões 3D no estilo clonado\n")
	sb.WriteString(fmt.Sprintf("4. **[00:45 - 00:60] Chamada para Ação (CTA):** Assinatura da marca %s\n", req.Marca.NomeMarca))
	return sb.String()
}

func getPrimaryHook(dna ChannelDNA) string {
	if len(dna.TopGanchos) > 0 {
		return dna.TopGanchos[0]
	}
	return "Descubra o segredo agora"
}

// ReadFile utilitário hermético de leitura de arquivos em disco.
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
