package videomotion

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var (
	// ErrInvalidBookConfig indica parâmetros incompletos para a geração do livro.
	ErrInvalidBookConfig = errors.New("videomotion: configuração do livro didático inválida")
)

// DidacticBookConfig parametriza a criação do manual em formato de livro.
type DidacticBookConfig struct {
	OutputDir    string
	PipelineName string
	Version      string
}

// LegoPressBookManifest representa a estrutura canônica de livro.json no LegoPress.
type LegoPressBookManifest struct {
	VersaoManifesto  string   `json:"versao_manifesto"`
	Titulo           string   `json:"titulo"`
	NomeBase         string   `json:"nome_base"`
	PastaPartes      string   `json:"pasta_partes"`
	Template         string   `json:"template"`
	Partes           []string `json:"partes"`
	Revisoes         []string `json:"revisoes"`
	HashFonte        string   `json:"hash_fonte"`
	UltimaCompilacao string   `json:"ultima_compilacao"`
}

// GenerateDidacticBook cria o livro didático canônico na pasta docs/books/<nome>.
func GenerateDidacticBook(cfg DidacticBookConfig) (string, error) {
	if cfg.OutputDir == "" || cfg.PipelineName == "" {
		return "", ErrInvalidBookConfig
	}
	bookDir := filepath.Join(cfg.OutputDir, cfg.PipelineName)
	partesDir := filepath.Join(bookDir, "partes")
	if err := os.MkdirAll(partesDir, 0755); err != nil {
		return "", err
	}
	manifest, chapters := buildBookChapters(cfg)
	if err := writeChapters(partesDir, chapters); err != nil {
		return "", err
	}
	return writeManifest(bookDir, manifest)
}

func writeChapters(partesDir string, chapters map[string]string) error {
	for file, content := range chapters {
		path := filepath.Join(partesDir, file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

func writeManifest(bookDir string, manifest LegoPressBookManifest) (string, error) {
	manifestPath := filepath.Join(bookDir, "livro.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", err
	}
	return manifestPath, os.WriteFile(manifestPath, data, 0644)
}

func buildBookChapters(cfg DidacticBookConfig) (LegoPressBookManifest, map[string]string) {
	manifest := LegoPressBookManifest{
		VersaoManifesto: "1.0",
		Titulo:          "Pipeline Canônico de Motion Video",
		NomeBase:        "2026-11-07_PIPELINE-MOTION-VIDEO",
		PastaPartes:     "partes",
		Template:        "ativos/livro.typst",
		Partes: []string{
			"00-frontmatter.md",
			"01-macro-fundacao.md",
			"02-meso-arquitetura-4-pilares.md",
			"03-micro-passo-a-passo.md",
			"04-pratica-receitas-multi-aspecto.md",
			"05-rotina-diff-updater.md",
		},
		Revisoes:         []string{},
		HashFonte:        cfg.PipelineName + "-v" + cfg.Version,
		UltimaCompilacao: "2026-11-07",
	}
	chapters := map[string]string{
		"00-frontmatter.md":                    cap0Content,
		"01-macro-fundacao.md":                 cap1Content,
		"02-meso-arquitetura-4-pilares.md":     cap2Content,
		"03-micro-passo-a-passo.md":            cap3Content,
		"04-pratica-receitas-multi-aspecto.md": cap4Content,
		"05-rotina-diff-updater.md":            cap5Content,
	}
	return manifest, chapters
}

const cap0Content = `---
title: "Pipeline Canônico de Motion Video"
subtitle: "Engenharia Reversa de Motion Design, Sincronia Fonética e Renderização por Código"
author:
  - Ecossistema LegoLoop
date: "07/11/2026"
lang: pt-BR
institute: "LegoLoop Foundation"
eyebrow: "MANUAL OPERACIONAL DIDÁTICO · VIDEO AS CODE"
cover-title: "Motion Video"
cover-accent: "Pipeline Canônico por Código"
cover-caption: "OS SEIS NÓS DO PIPELINE video-motion"
cover-stages:
  - "Silêncio & Ataque +40ms"
  - "Matriz Words & Punch RMS"
  - "Vidro 3D & Tipografia Anton"
  - "Cache SHA-256 & Bridge"
  - "Multi-Aspecto (9:16, 16:9, 1:1)"
  - "Broadcast Bundle & LegoPress"
header-mark: "MOTION VIDEO"
accent: "2563eb"
tagline: |
  Do tutorial ao código executável: pipeline determinístico, multi-aspecto e com zero dependência proprietária
toc: true
toc-depth: 2
abstract: |
  Este livro é o manual técnico e didático oficial do pipeline video-motion do LegoLoop.
---

# Como ler este livro

## Estrutura da obra

| Parte | Nível | Conteúdo |
| :--- | :--- | :--- |
| **I** | Macro | Fundação filosófica e determinismo estrito |
| **II** | Meso | Arquitetura dos quatro pilares técnicos |
| **III** | Micro | Guia prático de execução e stop-gates de QA |
| **IV** | Prática | Receitas multi-aspecto e broadcast bundle |
| **V** | Governança | Rotina de atualização contínua por diff |
`

const cap1Content = `# Parte I — Fundação e Princípio Determinístico

O pipeline Video Motion foi concebido sob a Lei Fundamental de Engenharia do LegoLoop:
**Tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de uma IA generativa.**

Cortes, alinhamento fonético, curvas de aceleração visual (springs e easings), cálculo de decibéis RMS e
normalização broadcast a -14 LUFS são executados por matemática exata e algoritmos locais.
`

const cap2Content = `# Parte II — A Arquitetura dos 4 Pilares

1. **Relógio Mestre Fonético (Speech-Sync):** A voz dita a linha do tempo. Cortes ocorrem 40ms antes do ataque vocal.
2. **Motion Reativo em Código (Remotion/Web-Canvas):** Cartões de vidro 3D e kinetic typography sem After Effects.
3. **Cache Criptográfico SHA-256:** Prompts sintéticos cacheados localmente com custo zero de rede em reexecuções.
4. **Bundle de Entrega Broadcast:** Geração simultânea de MP4, SRT, VTT, Copywriting e FCPXML para DaVinci Resolve.
`

const cap3Content = `# Parte III — Guia Prático de Execução

Execute a esteira com o comando canônico:
` + "```bash" + `
legoloop video-motion run --input video.mp4 --aspect 9:16 --output dist/
` + "```" + `
A folha de contato qa_sheet.jpg será gerada em menos de 2 segundos para aprovação humana no Stop-Gate.
`

const cap4Content = `# Parte IV — Receitas Multi-Aspect Ratio

- **9:16 (1080x1920):** Vertical para TikTok, Instagram Reels e YouTube Shorts.
- **16:9 (1920x1080):** Widescreen para YouTube desktop, aulas e documentários.
- **1:1 (1080x1080):** Quadrado para Feed do LinkedIn e Instagram carrossel.
`

const cap5Content = `# Parte V — Rotina de Atualização por Diff (LegoPress)

Sempre que o código ou o DAG do pipeline sofrer alterações no repositório, o comando LegoPress sincroniza
o livro didático de forma incremental:
` + "```bash" + `
legoloop legopress atualizar docs/books/video-motion/
` + "```" + `
`
