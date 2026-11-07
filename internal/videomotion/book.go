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

type bookManifest struct {
	Titulo    string        `json:"titulo"`
	Slug      string        `json:"slug"`
	Versao    string        `json:"versao"`
	Capitulos []bookChapter `json:"capitulos"`
}

type bookChapter struct {
	Ordem   int    `json:"ordem"`
	Titulo  string `json:"titulo"`
	Arquivo string `json:"arquivo"`
}

// GenerateDidacticBook cria os 5 capítulos didáticos e o manifesto livro.json.
func GenerateDidacticBook(cfg DidacticBookConfig) (string, error) {
	if cfg.OutputDir == "" || cfg.PipelineName == "" {
		return "", ErrInvalidBookConfig
	}
	bookDir := filepath.Join(cfg.OutputDir, "livro-"+cfg.PipelineName)
	if err := os.MkdirAll(bookDir, 0755); err != nil {
		return "", err
	}
	manifest, chapters := buildBookChapters(cfg)
	for file, content := range chapters {
		path := filepath.Join(bookDir, file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return "", err
		}
	}
	manifestPath := filepath.Join(bookDir, "livro.json")
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return "", err
	}
	return manifestPath, os.WriteFile(manifestPath, data, 0644)
}

func buildBookChapters(cfg DidacticBookConfig) (bookManifest, map[string]string) {
	manifest := bookManifest{
		Titulo: "Manual Operacional e Didático do Pipeline Video Motion",
		Slug:   cfg.PipelineName,
		Versao: cfg.Version,
		Capitulos: []bookChapter{
			{Ordem: 1, Titulo: "Fundação e Princípio Determinístico", Arquivo: "capitulo-01-fundacao.md"},
			{Ordem: 2, Titulo: "Arquitetura dos 4 Pilares", Arquivo: "capitulo-02-arquitetura-4-pilares.md"},
			{Ordem: 3, Titulo: "Guia Prático Passo a Passo", Arquivo: "capitulo-03-passo-a-passo.md"},
			{Ordem: 4, Titulo: "Receitas Multi-Aspecto (9:16, 16:9, 1:1)", Arquivo: "capitulo-04-receitas-multi-aspecto.md"},
			{Ordem: 5, Titulo: "Rotina de Atualização Contínua por Diff", Arquivo: "capitulo-05-rotina-diff-updater.md"},
		},
	}
	chapters := map[string]string{
		"capitulo-01-fundacao.md":               cap1Content,
		"capitulo-02-arquitetura-4-pilares.md":  cap2Content,
		"capitulo-03-passo-a-passo.md":          cap3Content,
		"capitulo-04-receitas-multi-aspecto.md": cap4Content,
		"capitulo-05-rotina-diff-updater.md":    cap5Content,
	}
	return manifest, chapters
}

const cap1Content = `# Capítulo 1 — Fundação e Princípio Determinístico

O pipeline Video Motion foi concebido sob a Lei Fundamental de Engenharia do LegoLoop:
**Tudo o que puder ser determinístico NUNCA deve ser deixado para a probabilidade de uma IA generativa.**

Cortes, alinhamento fonético, curvas de aceleração visual (springs e easings), cálculo de decibéis RMS e
normalização broadcast a -14 LUFS são executados por matemática exata e algoritmos locais.
`

const cap2Content = `# Capítulo 2 — A Arquitetura dos 4 Pilares

1. **Relógio Mestre Fonético (Speech-Sync):** A voz dita a linha do tempo. Cortes ocorrem 40ms antes do ataque vocal.
2. **Motion Reativo em Código (Remotion/Web-Canvas):** Cartões de vidro 3D e kinetic typography sem dependência de After Effects.
3. **Cache Criptográfico SHA-256:** Prompts sintéticos cacheados localmente com custo zero de rede em reexecuções.
4. **Bundle de Entrega Broadcast:** Geração simultânea de MP4, SRT, VTT, Copywriting e FCPXML para DaVinci Resolve.
`

const cap3Content = `# Capítulo 3 — Guia Prático de Execução

Execute a esteira com o comando canônico:
` + "```bash" + `
legoloop run loops/video-motion/pipeline.dag.json --var input_media="meu_video.mp4"
` + "```" + `
A folha de contato qa_sheet.jpg será gerada em menos de 3 segundos para aprovação humana no Stop-Gate.
`

const cap4Content = `# Capítulo 4 — Receitas Multi-Aspect Ratio

- **9:16 (1080x1920):** Vertical para TikTok, Instagram Reels e YouTube Shorts (com margens de segurança).
- **16:9 (1920x1080):** Widescreen para YouTube desktop, aulas e documentários.
- **1:1 (1080x1080):** Quadrado para Feed do LinkedIn e Instagram carrossel.
`

const cap5Content = `# Capítulo 5 — Rotina de Atualização por Diff (LegoPress)

Sempre que o código ou o DAG do pipeline sofrer alterações no repositório, o comando LegoPress sincroniza
o livro didático de forma incremental:
` + "```bash" + `
legoloop legopress atualizar docs/manuals/livro-video-motion/
` + "```" + `
`
