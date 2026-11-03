package specgen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"legoloop/docs/templates"
	"legoloop/internal/graph"
	"legoloop/internal/interview"
)

// Generator orquestra a criação determinística de especificações técnicas formais.
type Generator struct {
	llm   interview.LLMClient
	index *graph.Index
	// Now fornece a data do cabeçalho; nil usa time.Now. Injetável para geração determinística.
	Now func() time.Time
}

func (g *Generator) now() time.Time {
	if g.Now != nil {
		return g.Now()
	}
	return time.Now()
}

// NewGenerator inicializa o gerador híbrido de SPEC.
func NewGenerator(llm interview.LLMClient, idx *graph.Index) *Generator {
	return &Generator{
		llm:   llm,
		index: idx,
	}
}

// SpecResult contém os dados estruturados extraídos da especificação gerada.
type SpecResult struct {
	Theme       string   `json:"theme"`
	FilePath    string   `json:"file_path"`
	Criteria    []string `json:"criteria"`
	RawMarkdown string   `json:"raw_markdown"`
}

// specRole é a instrução do motor. A estrutura da SPEC não é descrita aqui: vem do TEMPLATE_SPEC.md
// embutido (docs/templates), a fonte única de verdade.
const specRole = `Você é o Arquiteto de Especificações Formais do LegoLoop (fusão Addy Osmani + Matt Pocock + AIDD-Spec).
Sua missão é transformar um documento de INTENÇÃO e o GRAFO de código existente em uma ESPECIFICAÇÃO TÉCNICA DE ENGENHARIA AUDITÁVEL.

REGRAS INVIOLÁVEIS:
1. Siga EXATAMENTE a estrutura do TEMPLATE CANÔNICO abaixo: mesmos títulos de seção, mesma numeração, mesma ordem. Não renomeie, não remova e não acrescente seções numeradas.
2. Substitua TODO {{PLACEHOLDER}} por conteúdo real. Nenhum "{{" pode sobrar. Apague os exemplos que não se aplicam.
3. Todo critério AC-XX deve ser verificável mecanicamente, com valores concretos. Cada épico precisa de ao menos um critério de caminho de falha, e requisitos não funcionais viram AC-NF.
4. Use somente comandos de verificação que existem no repositório (go build, go test, go vet, gofmt). golangci-lint só se estiver instalado.
5. O cabeçalho (status, versão, referência da intenção) é preenchido pelo motor; mantenha apenas o título "# 📋 SPEC: <tema>".
6. Encerre com o STOP-GATE 02 do template.

Data de hoje: %s. Referência da intenção: %s.

TEMPLATE CANÔNICO (TEMPLATE_SPEC.md v%s):

%s`

// SystemPrompt monta o prompt do sistema a partir do template canônico embutido.
func SystemPrompt(today, intentRef string) string {
	return fmt.Sprintf(specRole, today, intentRef, templates.Version, templates.StripGuidance(templates.Spec()))
}

// Generate consome o caminho do INTENT, consulta o grafo e gera a SPEC.
// ErrContract: a SPEC devolvida pelo LLM não segue o contrato do template. Diferente de falha de transporte,
// vale tentar de novo com o motivo da rejeição.
var ErrContract = errors.New("especificação gerada violou o contrato")

func (g *Generator) Generate(ctx context.Context, baseDir, intentPath string) (*SpecResult, error) {
	intentContent, err := os.ReadFile(intentPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao ler artefato de intenção: %w", err)
	}
	today := g.now().Format("2006-01-02")
	rawMarkdown, err := g.llm.Chat(ctx, specMessages(today, intentPath, string(intentContent)))
	if err != nil {
		return nil, fmt.Errorf("falha ao gerar spec via LLM: %w", err)
	}
	criteria, err := validateSpecMarkdown(rawMarkdown)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrContract, err)
	}
	specsDir, specFileName := specLocation(baseDir, intentPath, today)
	doc, err := stampHeader(rawMarkdown, newHeader(today, intentPath, specsDir, string(intentContent)))
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrContract, err)
	}
	specFilePath := filepath.Join(specsDir, specFileName)
	if err := os.WriteFile(specFilePath, []byte(doc), 0644); err != nil {
		return nil, fmt.Errorf("falha ao gravar spec no disco: %w", err)
	}
	return &SpecResult{Theme: specFileName, FilePath: specFilePath, Criteria: criteria, RawMarkdown: doc}, nil
}

func specMessages(today, intentPath, intent string) []interview.Message {
	user := fmt.Sprintf(`Gere a SPEC formal para a seguinte INTENÇÃO:

%s

Atenção: Garanta que todas as structs propostas sejam compatíveis com a arquitetura Go e os critérios AC-XX usem asserções binárias verificáveis mecanicamente.`, intent)
	return []interview.Message{
		{Role: "system", Content: SystemPrompt(today, filepath.ToSlash(intentPath))},
		{Role: "user", Content: user},
	}
}

// specLocation devolve a pasta docs/specs (criada) e o nome do arquivo derivado do nome da INTENT.
func specLocation(baseDir, intentPath, today string) (dir, name string) {
	dir = filepath.Join(baseDir, "docs", "specs")
	if err := os.MkdirAll(dir, 0755); err != nil {
		dir = baseDir
	}
	base := filepath.Base(intentPath)
	name = strings.Replace(base, "INTENT", "SPEC", 1)
	if name == base {
		name = fmt.Sprintf("%s_SPEC_%s.md", today, "GERADA")
	}
	return dir, name
}

// header reúne os campos do cabeçalho do TEMPLATE_SPEC que o motor conhece com certeza.
type header struct {
	Today, IntentDisplay, IntentLink, Target string
}

var targetRe = regexp.MustCompile(`(?m)^>\s*\*\*Alvo:\*\*\s*(.+?)\s*$`)

func newHeader(today, intentPath, specsDir, intent string) header {
	abs, err := filepath.Abs(intentPath)
	if err != nil {
		abs = intentPath
	}
	link := filepath.ToSlash(abs)
	if rel, err := filepath.Rel(specsDir, abs); err == nil {
		link = filepath.ToSlash(rel)
	}
	display := filepath.ToSlash(filepath.Join("docs", "intents", filepath.Base(intentPath)))
	if filepath.Base(filepath.Dir(abs)) != "intents" {
		display = filepath.ToSlash(filepath.Base(intentPath))
	}
	target := "A definir"
	if m := targetRe.FindStringSubmatch(intent); len(m) > 1 {
		target = m[1]
	}
	return header{Today: today, IntentDisplay: display, IntentLink: link, Target: target}
}

func (h header) render() string {
	var b strings.Builder
	b.WriteString("> **Status:** Proposta para Revisão  \n")
	fmt.Fprintf(&b, "> **Data:** %s  \n", h.Today)
	b.WriteString("> **Versão:** 1.0.0  \n")
	b.WriteString("> **Autor:** Antigravity AI + Usuário (Fusão Addy Osmani + Matt Pocock + AIDD-Spec)  \n")
	fmt.Fprintf(&b, "> **Referência da Intenção:** [%s](%s)  \n", h.IntentDisplay, h.IntentLink)
	fmt.Fprintf(&b, "> **Alvo Primário:** %s  \n", h.Target)
	b.WriteString("> **Consome:** INTENT aprovada · **Produz:** PLAN e ROUTE  \n")
	fmt.Fprintf(&b, "> **Template:** TEMPLATE_SPEC.md v%s  \n", templates.Version)
	return b.String()
}

var (
	titleRe   = regexp.MustCompile(`(?m)^# 📋 SPEC:.*$`)
	sectionRe = regexp.MustCompile(`(?m)^## `)
	// placeholderRe casa {{NOME}} do template, mas não a sintaxe de templates Go/Handlebars ({{.Campo}}, {{#if}}).
	placeholderRe = regexp.MustCompile(`\{\{\s*[A-Za-zÀ-ÿ_][^}\n]*\}\}`)
)

// stampHeader substitui tudo entre o título e a primeira seção pelo cabeçalho do template.
// O corpo escrito pelo modelo é preservado; o cabeçalho nunca depende do modelo.
func stampHeader(md string, h header) (string, error) {
	title := titleRe.FindStringIndex(md)
	if title == nil {
		return "", errors.New("a SPEC não tem o título '# 📋 SPEC: <tema>'")
	}
	rest := md[title[1]:]
	sec := sectionRe.FindStringIndex(rest)
	if sec == nil {
		return "", errors.New("a SPEC não tem nenhuma seção numerada após o título")
	}
	return md[:title[1]] + "\n\n" + h.render() + "\n---\n\n" + rest[sec[0]:], nil
}

// validateSpecMarkdown assegura a integridade das 6 seções obrigatórias e extrai critérios AC-XX.
func validateSpecMarkdown(content string) ([]string, error) {
	requiredHeadings := []string{
		"1. Premissas Técnicas Assumidas",
		"2. Fronteiras de Ação & Regras de Governança",
		"3. Visão Geral, Arquitetura & Costuras de Teste",
		"4. Contratos de Dados & Schemas Formais",
		"5. Critérios Binários de Aceite por Nó",
		"6. Comandos Operacionais & Critério de Conclusão Global",
		"STOP-GATE 02",
	}

	for _, h := range requiredHeadings {
		if !strings.Contains(content, h) {
			return nil, fmt.Errorf("seção obrigatória ausente na SPEC: %q", h)
		}
	}
	if m := placeholderRe.FindString(content); m != "" {
		return nil, fmt.Errorf("a SPEC deixou placeholder sem preencher: %s", m)
	}

	reAC := regexp.MustCompile(`(?i)\*\s*\[\s*\]\s*[\x60*]*(AC-[A-Za-z0-9_.-]+)`)
	matches := reAC.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return nil, errors.New("a SPEC não contém nenhum critério binário de aceite 'AC-XX'")
	}

	var criteria []string
	for _, m := range matches {
		if len(m) > 1 {
			criteria = append(criteria, m[1])
		}
	}

	return criteria, nil
}
